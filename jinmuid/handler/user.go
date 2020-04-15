package handler

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/jinmukeji/go-pkg/v2/crypto/encrypt/legacy"
	crypto "github.com/jinmukeji/go-pkg/v2/crypto/encrypt/legacy"
	"github.com/jinmukeji/go-pkg/v2/crypto/hash"
	"github.com/jinmukeji/go-pkg/v2/crypto/rand"
	"github.com/jinmukeji/jiujiantang-services/jinmuid/mysqldb"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
)

var (
	// 喜马把脉ID密码格式
	validPassword = regexp.MustCompile(`^[A-Za-z0-9]{8,20}$`)
)

// UserSetPassword 设置密码
func (j *JinmuIDService) UserSetPassword(ctx context.Context, req *proto.UserSetPasswordRequest, resp *proto.UserSetPasswordResponse) error {
	token, ok := TokenFromContext(ctx)
	if !ok {
		return NewError(ErrInvalidUser, errors.New("failed to get token from context"))
	}
	userID, err := j.datastore.FindUserIDByToken(ctx, token)
	if err != nil {
		return NewError(ErrUserUnauthorized, fmt.Errorf("failed to get userID by token: %s", err.Error()))
	}
	if userID != req.UserId {
		return NewError(ErrInvalidUser, fmt.Errorf("user %d from request and user %d from token are inconsistent", req.UserId, userID))
	}
	if req.PlainPassword == "" {
		return NewError(ErrEmptyPassword, fmt.Errorf("password of user %d is empty", req.UserId))
	}
	if !checkPasswordFormat(req.PlainPassword) {
		return NewError(ErrWrongFormatOfPassword, errors.New("password format is invalid"))
	}
	exist, errExistUserByUserID := j.datastore.ExistUserByUserID(ctx, req.UserId)
	if errExistUserByUserID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to check existence of user by userID %d: %s", req.UserId, errExistUserByUserID.Error()))
	}
	if !exist {
		return NewError(ErrInvalidUser, fmt.Errorf("user %d doesn't exist", req.UserId))
	}
	existPassword, errExistPasswordByUserID := j.datastore.ExistPasswordByUserID(ctx, req.UserId)
	if errExistPasswordByUserID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to check existence of password by userID %d: %s", req.UserId, errExistPasswordByUserID.Error()))
	}
	if existPassword {
		return NewError(ErrExistPassword, fmt.Errorf("password of user %d already exists", req.UserId))
	}
	seed, _ := rand.RandomStringWithMask(rand.MaskLetterDigits, 4)
	helper := crypto.NewPasswordCipherHelper()
	encryptedPassword := helper.Encrypt(req.PlainPassword, seed, j.encryptKey)
	errSetPasswordByUserID := j.datastore.SetPasswordByUserID(ctx, req.UserId, encryptedPassword, seed)
	if errSetPasswordByUserID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to set password of user %d: %s", req.UserId, errSetPasswordByUserID.Error()))
	}
	return nil
}

// UserModifyPassword 修改密码
func (j *JinmuIDService) UserModifyPassword(ctx context.Context, req *proto.UserModifyPasswordRequest, resp *proto.UserModifyPasswordResponse) error {
	if !checkPasswordFormat(req.NewPlainPassword) {
		return NewError(ErrWrongFormatOfPassword, errors.New("password format is invalid"))
	}
	token, ok := TokenFromContext(ctx)
	if !ok {
		return NewError(ErrInvalidUser, errors.New("failed to get token from context"))
	}
	userID, err := j.datastore.FindUserIDByToken(ctx, token)
	if err != nil {
		return NewError(ErrUserUnauthorized, fmt.Errorf("failed to get userID by token: %s", err.Error()))
	}
	if userID != req.UserId {
		return NewError(ErrInvalidUser, fmt.Errorf("user %d from request and user %d from token are inconsistent", req.UserId, userID))
	}
	exist, errExistUserByUserID := j.datastore.ExistUserByUserID(ctx, req.UserId)
	if errExistUserByUserID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to check existence of user by userID %d: %s", req.UserId, errExistUserByUserID.Error()))
	}
	if !exist {
		return NewError(ErrInvalidUser, fmt.Errorf("user %d doesn't exist", req.UserId))
	}
	user, err := j.datastore.FindUserByUserID(ctx, req.UserId)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find user by userID %d: %s", req.UserId, err.Error()))
	}
	if !user.HasSetPassword {
		return NewError(ErrEmptyPassword, fmt.Errorf("user %d has not set password", req.UserId))
	}
	helper := legacy.NewPasswordCipherHelper()
	password := helper.Decrypt(user.EncryptedPassword, user.Seed, j.encryptKey)
	if hash.HexString(hash.SHA256String(password+req.Seed)) != req.OldHashedPassword {
		return NewError(ErrIncorrectOldPassword, errors.New("old password is incorrect"))
	}
	// 判断密码是否跟以前一样
	if password == req.NewPlainPassword {
		return NewError(ErrSamePassword, errors.New("new password shouldn't be same as the old one"))
	}
	seed, _ := rand.RandomStringWithMask(rand.MaskLetterDigits, 4)
	encryptedPassword := helper.Encrypt(req.NewPlainPassword, seed, j.encryptKey)
	now := time.Now()
	clientID, _ := ClientIDFromContext(ctx)
	record := &mysqldb.AuditUserCredentialUpdate{
		UserID:            req.UserId,
		ClientID:          clientID,
		UpdatedRecordType: mysqldb.PasswordUpdated,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	// TODO: CreateAuditUserCredentialUpdate,SetPasswordByUserID,DeleteTokenByUserID要在同一个事务
	errCreateAuditUserCredentialUpdate := j.datastore.CreateAuditUserCredentialUpdate(ctx, record)
	if errCreateAuditUserCredentialUpdate != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to create audit user credential of user %d: %s", req.UserId, errCreateAuditUserCredentialUpdate.Error()))
	}
	errSetPasswordByUserID := j.datastore.SetPasswordByUserID(ctx, req.UserId, encryptedPassword, seed)
	if errSetPasswordByUserID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to set password to user %d: %s", req.UserId, errSetPasswordByUserID.Error()))
	}
	errDeleteTokenByUserID := j.datastore.DeleteTokenByUserID(ctx, userID)
	if errDeleteTokenByUserID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to delete token by userID %d: %s", userID, errDeleteTokenByUserID.Error()))
	}
	return nil
}

// UserGetUsingService 用户正在使用的服务
func (j *JinmuIDService) UserGetUsingService(ctx context.Context, req *proto.UserGetUsingServiceRequest, resp *proto.UserGetUsingServiceResponse) error {
	token, ok := TokenFromContext(ctx)
	if !ok {
		return NewError(ErrInvalidUser, errors.New("failed to get token from context"))
	}
	userID, err := j.datastore.FindUserIDByToken(ctx, token)
	if err != nil {
		return NewError(ErrUserUnauthorized, fmt.Errorf("failed to get userID by token: %s", err.Error()))
	}
	if userID != req.UserId {
		return NewError(ErrInvalidUser, fmt.Errorf("user %d from request and user %d from token are inconsistent", req.UserId, userID))
	}
	clients, errFindUsingClientIDs := j.datastore.FindUsingClients(ctx, req.UserId)
	if errFindUsingClientIDs != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find using clients by user %d: %s", req.UserId, errFindUsingClientIDs.Error()))
	}

	if len(clients) != 0 {
		cs := make([]*proto.Client, len(clients))
		for idx, client := range clients {
			cs[idx] = &proto.Client{
				ClientId: client.ClientID,
				Remark:   client.Remark,
				Usage:    client.Usage,
			}
		}
		resp.Clients = cs
	}
	return nil
}

// checkPasswordFormat 检查密码格式(长度限制8-20位)
func checkPasswordFormat(password string) bool {
	return validPassword.MatchString(password)
}
