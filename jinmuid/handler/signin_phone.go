package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jinmukeji/jiujiantang-services/jinmuid/mysqldb"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
)

// PhoneValificationExpiration 手机验证超时时间
const PhoneValificationExpiration = time.Minute * 10

// UserSetSigninPhone 用户设置登录手机
func (j *JinmuIDService) UserSetSigninPhone(ctx context.Context, req *proto.UserSetSigninPhoneRequest, resp *proto.UserSetSigninPhoneResponse) error {
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
	errValidatePhoneFormat := validatePhoneFormat(req.Phone, req.NationCode)
	if errValidatePhoneFormat != nil {
		return errValidatePhoneFormat
	}
	// 该手机是否注册过
	exsit, err := j.datastore.ExistPhone(ctx, req.Phone, req.NationCode)
	if exsit || err != nil {
		return NewError(ErrExistRegisteredPhone, fmt.Errorf("failed to check existence of phone %s%s", req.Phone, req.NationCode))
	}
	vcRecord, errFindVcRecord := j.datastore.FindVcRecord(ctx, req.SerialNumber, req.Mvc, req.Phone, mysqldb.SetPhoneNumber)
	if errFindVcRecord != nil {
		return NewError(ErrInValidMVC, fmt.Errorf("failed to find vc record of phone %s: %s", req.Phone, errFindVcRecord.Error()))
	}
	// 是否失效
	if vcRecord.HasUsed {
		return NewError(ErrUsedVcRecord, errors.New("vc record is used"))
	}
	// 是否过期
	if vcRecord.ExpiredAt.Before(time.Now()) {
		return NewError(ErrExpiredVcRecord, errors.New("expired vc record"))
	}
	// TODO: ModifyVcRecordStatus，SetSigninPhoneByUserID 要在同一个事务
	errModifyVcRecordStatus := j.datastore.ModifyVcRecordStatus(ctx, vcRecord.RecordID)
	if errModifyVcRecordStatus != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to modify vc record status %d: %s", vcRecord.RecordID, errModifyVcRecordStatus.Error()))
	}
	user, errFindUserByUserID := j.datastore.FindUserByUserID(ctx, userID)
	if errFindUserByUserID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find user by userId %d: %s", userID, errFindUserByUserID.Error()))
	}
	// 登录电话已经被设置就报这个错误
	if user.HasSetPhone {
		return NewError(ErrExsitSignInPhone, fmt.Errorf("signin phone has been set for userId %d", userID))
	}
	errSetSigninPhoneByUserID := j.datastore.SetSigninPhoneByUserID(ctx, userID, req.Phone, req.NationCode)
	if errSetSigninPhoneByUserID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed set signin phone %s%s to UserId %d: %s", req.NationCode, req.Phone, userID, errSetSigninPhoneByUserID.Error()))
	}
	return nil
}

// VerifyUserSigninPhone 验证用户登录手机号
func (j *JinmuIDService) VerifyUserSigninPhone(ctx context.Context, req *proto.VerifyUserSigninPhoneRequest, resp *proto.VerifyUserSigninPhoneResponse) error {
	// 验证手机号格式
	errValidatePhoneFormat := validatePhoneFormat(req.Phone, req.NationCode)
	if errValidatePhoneFormat != nil {
		return errValidatePhoneFormat
	}
	// 修改手机号码的时候应该在发送的时候就判断手机号码是否是用户设置的，
	// 手机号查找用户
	user, errFindUserByPhone := j.datastore.FindUserByPhone(ctx, req.Phone, req.NationCode)
	if errFindUserByPhone != nil {
		return NewError(ErrNotExistSigninPhone, fmt.Errorf("failed to find user by signin phone %s%s", req.NationCode, req.Phone))
	}
	// 修改手机号需要进行手机号用户验证
	if req.Action == proto.TemplateAction_TEMPLATE_ACTION_MODIFY_PHONE_NUMBER {
		token, ok := TokenFromContext(ctx)
		if !ok {
			return NewError(ErrInvalidUser, errors.New("failed to get token from context"))
		}
		userID, err := j.datastore.FindUserIDByToken(ctx, token)
		if err != nil {
			return NewError(ErrUserUnauthorized, fmt.Errorf("failed to get userID by token: %s", err.Error()))
		}
		// 判断手机号是否是该用户的
		if user.UserID != userID {
			return NewError(ErrSignInPhoneNotBelongsToUser, fmt.Errorf("signin phone %s%s doesn't belong to current user %d", req.NationCode, req.Phone, user.UserID))
		}
	}
	vcRecord, errFindVcRecord := j.datastore.SearchVcRecord(ctx, req.SerialNumber, req.Mvc, req.Phone, req.NationCode)
	// 是否有效
	if errFindVcRecord != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to search vc record of phone %s%s: %s", req.NationCode, req.Phone, errFindVcRecord.Error()))
	}
	if vcRecord == nil {
		return NewError(ErrInvalidVcRecord, errors.New("vc is invalid"))
	}
	// 判断是否过期
	if time.Now().After(*vcRecord.ExpiredAt) {
		return NewError(ErrExpiredVcRecord, errors.New("vc is expired"))
	}
	// 判断是否使用
	if vcRecord.HasUsed {
		return NewError(ErrUsedVcRecord, errors.New("vc has used"))
	}
	// TODO: ModifyVcRecordStatus，CreatePhoneOrEmailVerfication要在同一个事务
	errModifyVcRecordStatus := j.datastore.ModifyVcRecordStatus(ctx, vcRecord.RecordID)
	if errModifyVcRecordStatus != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to modify vc record status of record %d: %s", vcRecord.RecordID, errModifyVcRecordStatus.Error()))
	}
	Vn := uuid.New().String()
	now := time.Now()
	expiredAt := now.Add(PhoneValificationExpiration)
	record := &mysqldb.PhoneOrEmailVerfication{
		VerificationType:   mysqldb.VerificationPhone,
		VerificationNumber: Vn,
		UserID:             user.UserID,
		ExpiredAt:          &expiredAt,
		HasUsed:            false,
		NationCode:         req.NationCode,
		SendTo:             req.Phone,
		CreatedAt:          now.UTC(),
		UpdatedAt:          now.UTC(),
	}
	errCreatePhoneOrEmailVerfication := j.datastore.CreatePhoneOrEmailVerfication(ctx, record)
	if errCreatePhoneOrEmailVerfication != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to create phone or email verification: %s", errCreatePhoneOrEmailVerfication.Error()))
	}
	resp.VerificationNumber = record.VerificationNumber
	resp.UserId = user.UserID
	return nil
}

// UserModifyPhone 用户修改登录手机号
func (j *JinmuIDService) UserModifyPhone(ctx context.Context, req *proto.UserModifyPhoneRequest, resp *proto.UserModifyPhoneResponse) error {
	token, ok := TokenFromContext(ctx)
	if !ok {
		return NewError(ErrInvalidUser, errors.New("failed to get token from context"))
	}
	userID, err := j.datastore.FindUserIDByToken(ctx, token)
	if err != nil {
		return NewError(ErrUserUnauthorized, fmt.Errorf("failed to get userID by token: %s", err.Error()))
	}
	user, errFindUserByPhone := j.datastore.FindUserByPhone(ctx, req.OldPhone, req.OldNationCode)
	if errFindUserByPhone != nil {
		return NewError(ErrNotExistSigninPhone, fmt.Errorf("failed to find username by phone %s%s: %s", req.OldPhone, req.OldNationCode, errFindUserByPhone.Error()))
	}
	// 验证传入的user_id,phone,nation_code,token 都属于同一个用户
	if userID != req.UserId || user.UserID != req.UserId {
		return NewError(ErrInvalidUser, fmt.Errorf("token or phone or userID from request do not belong to current user %d", userID))
	}
	// 验证旧手机号格式
	errValidateOldhoneFormat := validatePhoneFormat(req.OldPhone, req.OldNationCode)
	if errValidateOldhoneFormat != nil {
		return errValidateOldhoneFormat
	}
	// 验证新手机号格式
	errValidatePhoneFormat := validatePhoneFormat(req.Phone, req.NationCode)
	if errValidatePhoneFormat != nil {
		return errValidatePhoneFormat
	}
	// 新手机是否注册过
	exsit, err := j.datastore.ExistPhone(ctx, req.Phone, req.NationCode)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to check existence of phone %s%s: %s", req.NationCode, req.Phone, err.Error()))
	}
	if exsit {
		return NewError(ErrExistRegisteredPhone, fmt.Errorf("phone %s%s has been registered", req.NationCode, req.Phone))
	}
	// 判断手机号是否跟以前一样
	if req.Phone == req.OldPhone && req.NationCode == req.OldNationCode {
		return NewError(ErrSamePhone, fmt.Errorf("new phone %s%s shouldn't be same as the old one", req.NationCode, req.Phone))
	}
	vcRecord, errFindVcRecord := j.datastore.SearchVcRecord(ctx, req.SerialNumber, req.Mvc, req.Phone, req.NationCode)
	// 是否有效
	if errFindVcRecord != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to search vc record of phone %s%s: %s", req.NationCode, req.Phone, errFindVcRecord.Error()))
	}
	if vcRecord == nil {
		return NewError(ErrInvalidVcRecord, fmt.Errorf("vc record of phone %s%s doesn't exist", req.NationCode, req.Phone))
	}
	// 判断是否过期
	if time.Now().After(*vcRecord.ExpiredAt) {
		return NewError(ErrExpiredVcRecord, errors.New("vc is expired"))
	}
	// 判断是否使用
	if vcRecord.HasUsed {
		return NewError(ErrUsedVcRecord, errors.New("vc has been used"))
	}
	// TODO: SetVcAsUsed，VerifyVerificationNumberByPhone，SetSigninPhoneByUserID，SetVerificationNumberAsUsed要在同一个事务
	// 设置验证码是使用过的
	errSetVcAsUsed := j.datastore.SetVcAsUsed(ctx, req.SerialNumber, req.Mvc, req.Phone, req.NationCode)
	if errSetVcAsUsed != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to set vc record of phone %s%s as used: %s", req.NationCode, req.Phone, errSetVcAsUsed.Error()))
	}
	// 判断验证号是否有效
	isValid, errVerifyVerificationNumberByPhone := j.datastore.VerifyVerificationNumberByPhone(ctx, req.VerificationNumber, req.OldPhone, req.OldNationCode)
	if errVerifyVerificationNumberByPhone != nil || !isValid {
		return NewError(ErrInvalidVerificationNumber, fmt.Errorf("failed to verify verification number by phone %s%s", req.OldNationCode, req.OldPhone))
	}
	// 修改手机
	errSetSigninPhoneByUserID := j.datastore.SetSigninPhoneByUserID(ctx, user.UserID, req.Phone, req.NationCode)
	if errSetSigninPhoneByUserID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to set sign in phone %s%s to user %d: %s", req.NationCode, req.Phone, user.UserID, errSetSigninPhoneByUserID.Error()))
	}
	// 设置验证号为已经使用的状态
	return j.datastore.SetVerificationNumberAsUsed(ctx, mysqldb.VerificationPhone, req.VerificationNumber)
}
