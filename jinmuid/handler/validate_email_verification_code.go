package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jinmukeji/jiujiantang-services/jinmuid/mysqldb"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
)

// 邮箱验证超时时间
const (
	EmailValificationExpiration = time.Minute * 10
	ResetPassword               = "reset_password"
	ModifySecureEmail           = "modify_secure_email"
)

// ValidateEmailVerificationCode 验证邮箱验证码是否正确
func (j *JinmuIDService) ValidateEmailVerificationCode(ctx context.Context, req *proto.ValidateEmailVerificationCodeRequest, resp *proto.ValidateEmailVerificationCodeResponse) error {
	verificationType, errMapEmailTypeToDB := mapEmailTypeToDB(req.VerificationType)
	if verificationType == mysqldb.Unknown || errMapEmailTypeToDB != nil {
		return NewError(ErrInvalidValidationType, errors.New("invalid validation type"))
	}

	user, errFindUserByUserID := j.datastore.FindUserBySecureEmail(ctx, req.Email)
	if errFindUserByUserID != nil {
		return NewError(ErrInvalidSigninEmail, fmt.Errorf("failed to find user by secure email %s: %s", req.Email, errFindUserByUserID.Error()))
	}

	// 查找可用的最新的验证码记录
	latestRecord, errFindVcRecord := j.datastore.FindLatestVcRecord(ctx, req.Email, verificationType)
	if errFindVcRecord != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find latest vc record by email %s: %s", req.Email, errFindVcRecord.Error()))
	}
	// 是否过期
	if latestRecord.ExpiredAt.Before(time.Now()) {
		return NewError(ErrExpiredVcRecord, errors.New("expired vc record"))
	}
	if !(req.VerificationCode == latestRecord.Code && req.SerialNumber == latestRecord.SN) {
		return NewError(ErrInvalidVcRecord, errors.New("invalid vc record"))
	}
	// TODO: ModifyVcRecordStatus,CreatePhoneOrEmailVerfication 要在同一个事务
	// 修改验证码状态
	err := j.datastore.ModifyVcRecordStatus(ctx, latestRecord.RecordID)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to modify vc record status of record %d: %s", latestRecord.RecordID, err.Error()))
	}

	Vn := uuid.New().String()
	now := time.Now()
	expiredAt := now.Add(EmailValificationExpiration)
	record := &mysqldb.PhoneOrEmailVerfication{
		VerificationType:   mysqldb.VerificationEmail,
		VerificationNumber: Vn,
		UserID:             user.UserID,
		ExpiredAt:          &expiredAt,
		HasUsed:            false,
		SendTo:             req.Email,
		CreatedAt:          now.UTC(),
		UpdatedAt:          now.UTC(),
	}
	errCreatePhoneOrEmailVerfication := j.datastore.CreatePhoneOrEmailVerfication(ctx, record)
	if errCreatePhoneOrEmailVerfication != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to create phone or email verification record of email %s: %s", req.Email, errCreatePhoneOrEmailVerfication.Error()))
	}
	resp.VerificationNumber = record.VerificationNumber
	resp.UserId = user.UserID
	return nil

}

func mapEmailTypeToDB(emailType string) (mysqldb.Usage, error) {
	switch emailType {
	case ResetPassword:
		return mysqldb.FindResetPassword, nil
	case ModifySecureEmail:
		return mysqldb.ModifySecureEmail, nil
	}
	return mysqldb.Unknown, fmt.Errorf("invalid string email type %s", emailType)
}
