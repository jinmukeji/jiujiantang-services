package handler

import (
	"context"
	"errors"
	"time"

    "github.com/jinmukeji/jiujiantang-services/jinmuid/mysqldb"
    "github.com/jinmukeji/go-pkg/crypto/encrypt/legacy"

	"fmt"

	"github.com/jinmukeji/go-pkg/crypto/rand"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
)

// UserResetPassword 用户重置密码
func (j *JinmuIDService) UserResetPassword(ctx context.Context, req *proto.UserResetPasswordRequest, resp *proto.UserResetPasswordResponse) error {
	protoVerificationType, errmapProtoVerificationTypeToDB := mapProtoVerificationTypeToDB(req.VerificationType)
	if errmapProtoVerificationTypeToDB != nil {
		return NewError(ErrInvalidVerificationType, errmapProtoVerificationTypeToDB)
	}
	isValid, errVerifyVerificationNumber := j.datastore.VerifyVerificationNumber(ctx, protoVerificationType, req.VerificationNumber, req.UserId)
	if errVerifyVerificationNumber != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to verify if verification number %s is invalid", req.VerificationNumber))
	}
	if !isValid {
		return NewError(ErrInvalidVerificationNumber, fmt.Errorf("verification number %s is invalid", req.VerificationNumber))
	}
	user, errFindUserByUserID := j.datastore.FindUserByUserID(ctx, req.UserId)
	if errFindUserByUserID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find user by userID %d: %s", req.UserId, errFindUserByUserID.Error()))
	}
	if !user.HasSetPassword {
		return NewError(ErrNotExistOldPassword, fmt.Errorf("old password of user %d does not exist", req.UserId))
	}
	helper := legacy.NewPasswordCipherHelper()
	oldPassword := helper.Decrypt(user.EncryptedPassword, user.Seed, j.encryptKey)
	if oldPassword == req.PlainPassword {
		return NewError(ErrSamePassword, errors.New("new password cannot equals old password"))
	}
	seed, _ := rand.RandomStringWithMask(rand.MaskLetterDigits, 4)
	encryptedPassword := helper.Encrypt(req.PlainPassword, seed, j.encryptKey)
	now := time.Now()
	clientID, _ := ClientIDFromContext(ctx)
	record := &mysqldb.AuditUserCredentialUpdate{
		UserID:            req.UserId,
		ClientID:          clientID,
		UpdatedRecordType: mysqldb.PasswordUpdated,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	// TODO: CreateAuditUserCredentialUpdate,SetPasswordByUserID,SetVerificationNumberAsUsed,DeleteTokenByUserID写在同一个事务中
	errCreateAuditUserCredentialUpdate := j.datastore.CreateAuditUserCredentialUpdate(ctx, record)
	if errCreateAuditUserCredentialUpdate != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to create audit user credential update: %s", errCreateAuditUserCredentialUpdate.Error()))
	}
	errSetPasswordByUserID := j.datastore.SetPasswordByUserID(ctx, req.UserId, encryptedPassword, seed)
	if errSetPasswordByUserID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to set password for user %d: %s", req.UserId, errSetPasswordByUserID.Error()))
	}
	errSetVerificationNumberAsUsed := j.datastore.SetVerificationNumberAsUsed(ctx, protoVerificationType, req.VerificationNumber)
	if errSetVerificationNumberAsUsed != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to set verification_number %s as used: %s", req.VerificationNumber, errSetVerificationNumberAsUsed.Error()))
	}
	errDeleteTokenByUserID := j.datastore.DeleteTokenByUserID(ctx, req.UserId)
	if errDeleteTokenByUserID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to delete token by user %d: %s", req.UserId, errDeleteTokenByUserID.Error()))
	}
	return nil
}

// mapProtoVerificationTypeToDB 转化proto中的VerificationType
func mapProtoVerificationTypeToDB(verificationType proto.VerificationType) (mysqldb.VerificationType, error) {
	switch verificationType {
	case proto.VerificationType_VERIFICATION_TYPE_INVALID:
		return mysqldb.VerificationPhone, fmt.Errorf("invalid proto verification type %d", verificationType)
	case proto.VerificationType_VERIFICATION_TYPE_UNSET:
		return mysqldb.VerificationEmail, fmt.Errorf("invalid proto verification type %d", verificationType)
	case proto.VerificationType_VERIFICATION_TYPE_PHONE:
		return mysqldb.VerificationPhone, nil
	case proto.VerificationType_VERIFICATION_TYPE_EMAIL:
		return mysqldb.VerificationEmail, nil
	}
	return mysqldb.VerificationPhone, fmt.Errorf("invalid proto verification type %d", verificationType)
}
