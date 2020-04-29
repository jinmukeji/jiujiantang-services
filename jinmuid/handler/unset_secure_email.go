package handler

import (
	"context"
	"errors"
	"time"

	"fmt"

	"github.com/jinmukeji/jiujiantang-services/jinmuid/mysqldb"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
)

// UnsetSecureEmail 解除设置安全邮箱
func (j *JinmuIDService) UnsetSecureEmail(ctx context.Context, req *proto.UnsetSecureEmailRequest, resp *proto.UnsetSecureEmailResponse) error {
	if !checkEmailFormat(req.Email) {
		return NewError(ErrInvalidEmailAddress, errors.New("invalid email format"))
	}
	// 查找可用的最新的验证码记录
	latestRecord, errFindVcRecord := j.datastore.FindLatestVcRecord(ctx, req.Email, mysqldb.UnsetSecureEmail)
	if errFindVcRecord != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find latest vc record of email %s of unset secure email: %s", req.Email, errFindVcRecord.Error()))
	}
	// 是否过期
	if latestRecord.ExpiredAt.Before(time.Now()) {
		return NewError(ErrExpiredVcRecord, errors.New("expired vc record"))
	}
	if !(req.VerificationCode == latestRecord.Code && req.SerialNumber == latestRecord.SN) {
		return NewError(ErrInvalidVcRecord, errors.New("invalid vc record"))
	}
	// 修改验证码状态
	err := j.datastore.ModifyVcRecordStatus(ctx, latestRecord.RecordID)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to ModifyVcRecordStatus of record %d: %s", latestRecord.RecordID, err.Error()))
	}
	secureEmailExisting, _ := j.datastore.SecureEmailExists(ctx, req.UserId)
	if !secureEmailExisting {
		return NewError(ErrSecureEmailNotSet, fmt.Errorf("user %d has been set email", req.UserId))
	}
	secureEmailMatching, _ := j.datastore.MatchSecureEmail(ctx, req.Email, req.UserId)
	if !secureEmailMatching {
		return NewError(ErrSecureEmailAddressNotMatched, fmt.Errorf("email %s has not been set by current user %d", req.Email, req.UserId))
	}
	// TODO: UnsetSecureEmail，CreateAuditUserCredentialUpdate要在同一个事务
	err = j.datastore.UnsetSecureEmail(ctx, req.UserId)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to unset secure email of user %d", req.UserId))
	}

	// 增加审计记录
	now := time.Now()
	clientID, _ := ClientIDFromContext(ctx)
	record := &mysqldb.AuditUserCredentialUpdate{
		UserID:            req.UserId,
		ClientID:          clientID,
		UpdatedRecordType: mysqldb.EmailUpdated,
		OldValue:          req.Email,
		NewValue:          "",
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	errCreateAuditUserCredentialUpdate := j.datastore.CreateAuditUserCredentialUpdate(ctx, record)
	if errCreateAuditUserCredentialUpdate != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to create audit user credential update record of user %d: %s", req.UserId, errCreateAuditUserCredentialUpdate.Error()))
	}
	return nil
}
