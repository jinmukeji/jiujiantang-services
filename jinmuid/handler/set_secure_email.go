package handler

import (
	"context"
	"errors"
	"time"

	"fmt"

	"github.com/jinmukeji/gf-api2/jinmuid/mysqldb"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
)

// SetSecureEmail 设置安全邮箱
func (j *JinmuIDService) SetSecureEmail(ctx context.Context, req *proto.SetSecureEmailRequest, resp *proto.SetSecureEmailResponse) error {
	// 验证邮箱格式
	if !checkEmailFormat(req.Email) {
		return NewError(ErrInvalidEmailAddress, fmt.Errorf("wrong format of email address %s", req.Email))
	}
	// 查找可用的最新的验证码记录
	latestRecord, errFindVcRecord := j.datastore.FindLatestVcRecord(ctx, req.Email, mysqldb.SetSecureEmail)
	if errFindVcRecord != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find latest vc code of secure email %s: %s", req.Email, errFindVcRecord.Error()))
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
		return NewError(ErrDatabase, fmt.Errorf("failed to modify vc record status of record %d: %s", latestRecord.RecordID, err.Error()))
	}

	// 判断邮箱是否已经被任何人设置
	hasSetSecureEmailByAnyone, _ := j.datastore.HasSetSecureEmailByAnyone(ctx, req.Email)
	if hasSetSecureEmailByAnyone {
		return NewError(ErrSecureEmailUsedByOthers, fmt.Errorf("failed to check if secure email %s has been set by anyone", req.Email))
	}
	// 判断是否已经设置了安全邮箱
	secureEmailExisting, _ := j.datastore.SecureEmailExists(ctx, req.UserId)
	if secureEmailExisting {
		return NewError(ErrSecureEmailExists, fmt.Errorf("user %d has set the email before", req.UserId))
	}
	// TODO: SetSecureEmail,CreateAuditUserCredentialUpdate要在同一个事务
	// 设置安全邮箱
	errSetSecureEmail := j.datastore.SetSecureEmail(ctx, req.Email, req.UserId)
	if errSetSecureEmail != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to set secure email of user %d", req.UserId))
	}
	// 增加审计记录
	now := time.Now()
	clientID, _ := ClientIDFromContext(ctx)
	record := &mysqldb.AuditUserCredentialUpdate{
		UserID:            req.UserId,
		ClientID:          clientID,
		UpdatedRecordType: mysqldb.EmailUpdated,
		OldValue:          "",
		NewValue:          req.Email,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	errCreateAuditUserCredentialUpdate := j.datastore.CreateAuditUserCredentialUpdate(ctx, record)
	if errCreateAuditUserCredentialUpdate != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to create audit user credential update: %s", errCreateAuditUserCredentialUpdate.Error()))
	}
	return nil
}
