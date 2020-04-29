package handler

import (
	"context"
	"errors"
	"time"

	"fmt"

	"github.com/google/uuid"
	"github.com/jinmukeji/jiujiantang-services/jinmuid/mysqldb"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
)

// ValidatePhoneVerificationCode 注册时验证手机验证码是否正确
func (j *JinmuIDService) ValidatePhoneVerificationCode(ctx context.Context, req *proto.ValidatePhoneVerificationCodeRequest, resp *proto.ValidatePhoneVerificationCodeResponse) error {
	// 该手机是否注册过
	exsit, err := j.datastore.ExistPhone(ctx, req.Phone, req.NationCode)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to check existence of phone %s%s: %s", req.NationCode, req.Phone, err.Error()))
	}
	if exsit {
		return NewError(ErrExistRegisteredPhone, fmt.Errorf("phone %s%s has been registered", req.NationCode, req.Phone))
	}
	vcRecord, errFindVcRecord := j.datastore.FindVcRecord(ctx, req.SerialNumber, req.Mvc, req.Phone, mysqldb.SignUp)
	if errFindVcRecord != nil {
		return NewError(ErrInvalidVcRecord, fmt.Errorf("failed to find vc record by phone %s%s: %s", req.NationCode, req.Phone, errFindVcRecord.Error()))
	}
	// 是否失效
	if vcRecord.HasUsed {
		return NewError(ErrUsedVcRecord, errors.New("vc is used"))
	}
	// 是否过期
	if vcRecord.ExpiredAt.Before(time.Now()) {
		return NewError(ErrExpiredVcRecord, errors.New("vc is expired"))
	}
	// TODO: ModifyVcRecordStatus,CreatePhoneOrEmailVerfication 要在同一个事务
	// 修改验证码的使用状态为已使用
	errModifyVcRecordStatus := j.datastore.ModifyVcRecordStatus(ctx, vcRecord.RecordID)
	if errModifyVcRecordStatus != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to modify vc record status of recordID %d: %s", vcRecord.RecordID, errModifyVcRecordStatus.Error()))
	}
	Vn := uuid.New().String()
	now := time.Now()
	expiredAt := now.Add(PhoneValificationExpiration)
	record := &mysqldb.PhoneOrEmailVerfication{
		VerificationType:   mysqldb.VerificationPhone,
		VerificationNumber: Vn,
		ExpiredAt:          &expiredAt,
		NationCode:         req.NationCode,
		HasUsed:            false,
		SendTo:             req.Phone,
		CreatedAt:          now.UTC(),
		UpdatedAt:          now.UTC(),
	}
	errCreatePhoneOrEmailVerfication := j.datastore.CreatePhoneOrEmailVerfication(ctx, record)
	if errCreatePhoneOrEmailVerfication != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to create phone verfication of phone %s%s: %s", req.NationCode, req.Phone, errCreatePhoneOrEmailVerfication.Error()))
	}
	resp.VerificationNumber = record.VerificationNumber
	return nil
}
