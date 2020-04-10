package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jinmukeji/jiujiantang-services/jinmuid/mysqldb"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
)

// ModifySecureEmail 修改安全邮箱
func (j *JinmuIDService) ModifySecureEmail(ctx context.Context, req *proto.ModifySecureEmailRequest, resp *proto.ModifySecureEmailResponse) error {
	token, ok := TokenFromContext(ctx)
	if !ok {
		return NewError(ErrInvalidUser, errors.New("failed to get token from context"))
	}
	userID, err := j.datastore.FindUserIDByToken(ctx, token)
	if err != nil {
		return NewError(ErrUserUnauthorized, fmt.Errorf("failed to get userID by token: %s", err.Error()))
	}
	user, errFindUserByEmail := j.datastore.FindUserByEmail(ctx, req.OldEmail)
	if errFindUserByEmail != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find user by old email %s: %s", req.OldEmail, errFindUserByEmail.Error()))
	}
	// 验证传入的user_id,phone,nation_code,token 都属于同一个用户
	if userID != req.UserId || user.UserID != req.UserId {
		return NewError(ErrInvalidUser, fmt.Errorf("token or phone or userID from request do not belong to current user %d", userID))
	}
	// 验证旧邮件格式
	if !checkEmailFormat(req.OldEmail) {
		return NewError(ErrInvalidEmailAddress, fmt.Errorf("old email %s format is invalid", req.OldEmail))
	}
	// 验证新邮件格式
	if !checkEmailFormat(req.NewEmail) {
		return NewError(ErrInvalidEmailAddress, fmt.Errorf("new email %s format is invalid", req.NewEmail))
	}
	// 判断安全邮箱是否跟以前一样
	if req.NewEmail == req.OldEmail {
		return NewError(ErrSameEmail, fmt.Errorf("new email %s shouldn't be same as the old one", req.NewEmail))
	}
	// 判断旧邮箱当前绑定邮箱
	user, errFindUserBySecureEmail := j.datastore.FindUserBySecureEmail(ctx, req.OldEmail)
	if errFindUserBySecureEmail != nil {
		return NewError(ErrNoneExistSecureEmail, fmt.Errorf("failed to find user by old email %s: %s", req.OldEmail, errFindUserBySecureEmail.Error()))
	}
	if user.UserID != req.UserId {
		return NewError(ErrSecureEmailAddressNotMatched, fmt.Errorf("secure email doesn't belong to current user %d", req.UserId))
	}
	if !user.HasSetEmail {
		return NewError(ErrSecureEmailNotSet, fmt.Errorf("old email %s has not been set", req.OldEmail))
	}
	if user.SecureEmail == req.NewEmail {
		return NewError(ErrSameEmail, fmt.Errorf("new email %s shouldn't be same as the old one", req.NewEmail))
	}
	// 判断新邮箱是否已经被任何人设置
	hasSetSecureEmailByAnyone, _ := j.datastore.HasSetSecureEmailByAnyone(ctx, req.NewEmail)
	if hasSetSecureEmailByAnyone {
		return NewError(ErrSecureEmailUsedByOthers, fmt.Errorf("the secure email %s has been used", req.NewEmail))
	}
	// 判断验证码是否有效
	record, errVerifyMVC := j.datastore.VerifyMVCBySecureEmail(ctx, req.NewSerialNumber, req.NewVerificationCode, req.NewEmail)
	if errVerifyMVC != nil {
		return NewError(ErrInvalidVcRecord, fmt.Errorf("failed to verify MVC by secure email %s: %s", req.NewEmail, errVerifyMVC.Error()))
	}

	// 是否有效
	if record.ExpiredAt.Before(time.Now()) || record.HasUsed {
		return NewError(ErrExpiredVcRecord, errors.New("expired vc record"))
	}

	// 判断验证号是否有效
	isValid, errVerifyVerificationNumberByEmail := j.datastore.VerifyVerificationNumberByEmail(ctx, req.OldVerificationNumber, req.OldEmail)
	if errVerifyVerificationNumberByEmail != nil || !isValid {
		return NewError(ErrInvalidVerificationNumber, fmt.Errorf("failed to verify verification number %s by email %s", req.OldVerificationNumber, req.OldEmail))
	}
	// TODO: 修改安全邮箱,修改新邮箱的验证码状态,设置验证号为已经使用的状态要写成事务
	// 修改安全邮箱
	errSetSecureEmailByUserID := j.datastore.SetSecureEmailByUserID(ctx, user.UserID, req.NewEmail)
	if errSetSecureEmailByUserID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to set secure email by user %d: %s", user.UserID, errSetSecureEmailByUserID.Error()))
	}
	// 修改新邮箱的验证码状态
	errModifyVcRecordStatusByEmail := j.datastore.ModifyVcRecordStatusByEmail(ctx, req.NewEmail, req.NewVerificationCode, req.NewSerialNumber)
	if errModifyVcRecordStatusByEmail != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to modify vc record status by email %s: %s", req.NewEmail, errModifyVcRecordStatusByEmail.Error()))
	}
	// 设置验证号为已经使用的状态
	errSetVerificationNumberAsUsed := j.datastore.SetVerificationNumberAsUsed(ctx, mysqldb.VerificationEmail, req.OldVerificationNumber)
	if errSetVerificationNumberAsUsed != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to set verification number of email to the status of used %s: %s", req.OldVerificationNumber, errSetVerificationNumberAsUsed.Error()))
	}
	return nil
}
