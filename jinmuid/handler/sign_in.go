package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"

	"github.com/google/uuid"
	"github.com/jinmukeji/jiujiantang-services/jinmuid/mysqldb"
	"github.com/jinmukeji/go-pkg/crypto/encrypt/legacy"
	"github.com/jinmukeji/go-pkg/crypto/hash"
	jinmuidpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	generalpb "github.com/jinmukeji/proto/v3/gen/micro/idl/ptypes/v2"
)

const (
	TokenAvailableDuration = time.Hour * 24 * 2
)

// UserSignInByPhonePassword 手机号密码登录
func (j *JinmuIDService) UserSignInByPhonePassword(ctx context.Context, req *jinmuidpb.UserSignInByPhonePasswordRequest, resp *jinmuidpb.UserSignInByPhonePasswordResponse) error {
	exist, errExistSignInPhone := j.datastore.ExistSignInPhone(ctx, req.Phone, req.NationCode)
	if errExistSignInPhone != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to check existence of phone %s%s", req.NationCode, req.Phone))
	}
	if !exist {
		return NewError(ErrNoneExistentPhone, fmt.Errorf("phone %s%s doesn't exist", req.NationCode, req.Phone))
	}
	user, err := j.datastore.FindUserByPhone(ctx, req.Phone, req.NationCode)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find username by phone %s%s: %s", req.NationCode, req.Phone, err.Error()))
	}
	if !user.IsActivated {
		return NewError(ErrDeactivatedUser, fmt.Errorf("user %d is deactivated", user.UserID))
	}
	if !user.HasSetPassword {
		return NewError(ErrNonexistentPassword, fmt.Errorf("password of user %d does not exist", user.UserID))
	}
	// 如果没有设置手机号
	if !user.HasSetPhone {
		return NewError(ErrNoneExistentPhone, fmt.Errorf("user %d has not set signin phone", user.UserID))
	}

	helper := legacy.NewPasswordCipherHelper()
	password := helper.Decrypt(user.EncryptedPassword, user.Seed, j.encryptKey)
	if hash.HexString(hash.SHA256String(password+req.Seed)) != req.HashedPassword {
		return NewError(ErrPhonePasswordNotMatch, fmt.Errorf("password and phone of user %d does not match", user.UserID))
	}
	token := uuid.New().String()
	tk, err := j.datastore.CreateToken(ctx, token, int32(user.UserID), TokenAvailableDuration)
	if err != nil {
		return NewError(ErrGetAccessTokenFailure, fmt.Errorf("failed to create token of user %d: %s", user.UserID, err.Error()))
	}
	json, errMarshal := json.Marshal(req)
	if errMarshal != nil {
		return fmt.Errorf("failed to parse request: %s", errMarshal.Error())
	}
	// 创建登录审计记录
	errCreateAuditUserSignin := j.createAuditUserSignin(ctx, user.UserID, req.Ip, req.SignInMachine, string(json))
	if errCreateAuditUserSignin != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to create audit user signin or signout of user %d: %s", user.UserID, errCreateAuditUserSignin.Error()))
	}
	if user.IsProfileCompleted {
		language, errmapProtoLanguageToDB := mapProtoLanguageToDB(generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE)
		if errmapProtoLanguageToDB != nil {
			return NewError(ErrInvalidLanguage, errmapProtoLanguageToDB)
		}
		err := j.datastore.SetLanguageByUserID(ctx, user.UserID, language)
		if err != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to set language for user %d: %s", user.UserID, err.Error()))
		}
	}
	resp.UserId = user.UserID
	resp.AccessToken = tk.Token
	resp.HasSetUserProfile = user.HasSetUserProfile
	resp.HasSetRegion = user.HasSetRegion
	resp.HasSetLanguage = user.HasSetLanguage
	resp.HasSetPassword = user.HasSetPassword
	resp.HasSetPhone = user.HasSetPhone
	resp.IsProfileCompleted = user.IsProfileCompleted
	expiredAt, _ := ptypes.TimestampProto(tk.ExpiredAt)
	resp.ExpiredTime = expiredAt
	return nil
}

// UserSignInByUsernamePassword 用户名密码登录
func (j *JinmuIDService) UserSignInByUsernamePassword(ctx context.Context, req *jinmuidpb.UserSignInByUsernamePasswordRequest, resp *jinmuidpb.UserSignInByUsernamePasswordResponse) error {
	user, err := j.datastore.FindUserByUsername(ctx, req.Username)
	if user == nil || err != nil {
		return NewError(ErrNonexistentUsername, fmt.Errorf("failed to find user by username %s", req.Username))
	}
	if !user.IsActivated {
		return NewError(ErrDeactivatedUser, fmt.Errorf("user %d is deactivated", user.UserID))
	}
	if !user.HasSetPassword {
		return NewError(ErrNonexistentPassword, fmt.Errorf("password of user %d does not exist", user.UserID))
	}
	helper := legacy.NewPasswordCipherHelper()
	password := helper.Decrypt(user.EncryptedPassword, user.Seed, j.encryptKey)
	if hash.HexString(hash.SHA256String(password+req.Seed)) != req.HashedPassword {
		return NewError(ErrUsernamePasswordNotMatch, fmt.Errorf("username and password of user %d does not match", user.UserID))
	}
	token := uuid.New().String()
	tk, err := j.datastore.CreateToken(ctx, token, int32(user.UserID), TokenAvailableDuration)
	if err != nil {
		return NewError(ErrGetAccessTokenFailure, fmt.Errorf("failed to create access token of user %d", user.UserID))
	}
	json, errMarshal := json.Marshal(req)
	if errMarshal != nil {
		return fmt.Errorf("failed to parse request: %s", errMarshal.Error())
	}
	// 创建登录审计记录
	errCreateAuditUserSignin := j.createAuditUserSignin(ctx, user.UserID, req.Ip, req.SignInMachine, string(json))
	if errCreateAuditUserSignin != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to create audit user signin or signout of user %d: %s", user.UserID, errCreateAuditUserSignin.Error()))
	}
	if user.IsProfileCompleted {
		language, ermapProtoLanguageToDB := mapProtoLanguageToDB(generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE)
		if ermapProtoLanguageToDB != nil {
			return NewError(ErrInvalidLanguage, ermapProtoLanguageToDB)
		}
		err := j.datastore.SetLanguageByUserID(ctx, user.UserID, language)
		if err != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to set the language of user %d to %s: %s", user.UserID, language, err.Error()))
		}
	}
	resp.UserId = user.UserID
	resp.AccessToken = tk.Token
	resp.HasSetUserProfile = user.HasSetUserProfile
	resp.HasSetRegion = user.HasSetRegion
	resp.HasSetLanguage = user.HasSetLanguage
	resp.HasSetPassword = user.HasSetPassword
	resp.HasSetPhone = user.HasSetPhone
	resp.IsProfileCompleted = user.IsProfileCompleted
	expiredAt, _ := ptypes.TimestampProto(tk.ExpiredAt)
	resp.ExpiredTime = expiredAt
	return nil
}

// UserSignInByPhoneVC 手机号验证码登录
func (j *JinmuIDService) UserSignInByPhoneVC(ctx context.Context, req *jinmuidpb.UserSignInByPhoneVCRequest, resp *jinmuidpb.UserSignInByPhoneVCResponse) error {
	exist, errExistSignInPhone := j.datastore.ExistSignInPhone(ctx, req.Phone, req.NationCode)
	if errExistSignInPhone != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to check if phone %s%s exists", req.NationCode, req.Phone))
	}
	if !exist {
		return NewError(ErrNoneExistentPhone, fmt.Errorf("phone %s%s doesn't exist", req.NationCode, req.Phone))
	}
	user, err := j.datastore.FindUserByPhone(ctx, req.Phone, req.NationCode)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find username by phone %s%s: %s", req.NationCode, req.Phone, err.Error()))
	}
	if !user.IsActivated {
		return NewError(ErrDeactivatedUser, fmt.Errorf("user %d is deactivated", user.UserID))
	}
	vcRecord, errFindVcRecord := j.datastore.FindVcRecord(ctx, req.SerialNumber, req.Mvc, req.Phone, mysqldb.SignIn)
	if errFindVcRecord != nil {
		return NewError(ErrInvalidVcRecord, fmt.Errorf("failed to find Vc record of phone %s%s: %s", req.NationCode, req.Phone, errFindVcRecord.Error()))
	}
	// 是否失效
	if vcRecord.HasUsed {
		return NewError(ErrUsedVcRecord, errors.New("vc record is used"))
	}
	// 是否过期
	if vcRecord.ExpiredAt.Before(time.Now()) {
		return NewError(ErrExpiredVcRecord, errors.New("expired vc record"))
	}
	// TODO: ModifyVcRecordStatus,CreateToken要在同一个事务
	errModifyVcRecordStatus := j.datastore.ModifyVcRecordStatus(ctx, vcRecord.RecordID)
	if errModifyVcRecordStatus != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to modify vc record status of record %d: %s", vcRecord.RecordID, errModifyVcRecordStatus.Error()))
	}
	token := uuid.New().String()
	tk, err := j.datastore.CreateToken(ctx, token, int32(user.UserID), TokenAvailableDuration)
	if err != nil {
		return NewError(ErrGetAccessTokenFailure, fmt.Errorf("failed to create access token of user %d: %s", user.UserID, err.Error()))
	}
	json, errMarshal := json.Marshal(req)
	if errMarshal != nil {
		return fmt.Errorf("failed to parse request: %s", errMarshal.Error())
	}
	// 创建登录审计记录
	errCreateAuditUserSignin := j.createAuditUserSignin(ctx, user.UserID, req.Ip, req.SignInMachine, string(json))
	if errCreateAuditUserSignin != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to create audit user signin or signout of user %d: %s", user.UserID, errCreateAuditUserSignin.Error()))
	}
	if user.IsProfileCompleted {
		language, ermapProtoLanguageToDB := mapProtoLanguageToDB(generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE)
		if ermapProtoLanguageToDB != nil {
			return NewError(ErrInvalidLanguage, ermapProtoLanguageToDB)
		}
		err := j.datastore.SetLanguageByUserID(ctx, user.UserID, language)
		if err != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to set the language of user %d to %s: %s", user.UserID, language, err.Error()))
		}
	}
	resp.UserId = user.UserID
	resp.AccessToken = tk.Token
	resp.HasSetUserProfile = user.HasSetUserProfile
	resp.HasSetRegion = user.HasSetRegion
	resp.HasSetLanguage = user.HasSetLanguage
	resp.HasSetPassword = user.HasSetPassword
	resp.HasSetPhone = user.HasSetPhone
	resp.IsProfileCompleted = user.IsProfileCompleted
	expiredAt, _ := ptypes.TimestampProto(tk.ExpiredAt)
	resp.ExpiredTime = expiredAt
	return nil
}

// createAuditUserSignin 创建登录审计记录
func (j *JinmuIDService) createAuditUserSignin(ctx context.Context, userID int32, ip, signInMachine, extraParams string) error {
	clientID, _ := ClientIDFromContext(ctx)
	errCreateAuditUserSigninSignout := j.datastore.CreateAuditUserSigninSignout(ctx, &mysqldb.AuditUserSigninSignout{
		UserID:        userID,
		ClientID:      clientID,
		RecordType:    mysqldb.SigninRecordType,
		IP:            ip,
		ExtraParams:   extraParams,
		SignInMachine: signInMachine,
	})
	if errCreateAuditUserSigninSignout != nil {
		return NewError(ErrGetAccessTokenFailure, fmt.Errorf("failed to create audit user signin or signout of user %d and client %s: %s", userID, clientID, errCreateAuditUserSigninSignout.Error()))
	}
	return nil
}

// GetLatestToken 最新的token
func (j *JinmuIDService) GetLatestToken(ctx context.Context, req *jinmuidpb.GetLatestTokenRequest, resp *jinmuidpb.GetLatestTokenResponse) error {
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
	accessToken := uuid.New().String()
	tk, err := j.datastore.CreateToken(ctx, accessToken, req.UserId, TokenAvailableDuration)
	if err != nil {
		return NewError(ErrGetAccessTokenFailure, fmt.Errorf("failed to create token of user %d: %s", req.UserId, err.Error()))
	}
	resp.AccessToken = tk.Token
	expiredAt, _ := ptypes.TimestampProto(tk.ExpiredAt)
	resp.ExpiredTime = expiredAt
	return nil
}

// UserGetSignInMachines 获取登录的设备
func (j *JinmuIDService) UserGetSignInMachines(ctx context.Context, req *jinmuidpb.UserGetSignInMachinesRequest, resp *jinmuidpb.UserGetSignInMachinesResponse) error {
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
	auditUserSigninSignouts, errFindUserSigninRecord := j.datastore.FindUserSigninRecord(ctx, req.UserId)
	if errFindUserSigninRecord != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find user sign in record of user %d: %s", req.UserId, errFindUserSigninRecord.Error()))
	}
	if len(auditUserSigninSignouts) != 0 {
		signInMachines := make([]*jinmuidpb.SignInMachine, len(auditUserSigninSignouts))
		for idx, auditUserSigninSignout := range auditUserSigninSignouts {
			createdAt, _ := ptypes.TimestampProto(auditUserSigninSignout.CreatedAt)
			signInMachines[idx] = &jinmuidpb.SignInMachine{
				SignInMachine: auditUserSigninSignout.SignInMachine,
				SignInTime:    createdAt,
			}
		}
		resp.Machines = signInMachines
	}
	return nil
}
