package handler

import (
	"context"
	"errors"
	"time"

	"fmt"

	"github.com/jinmukeji/jiujiantang-services/jinmuid/mysqldb"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
)

// FindUsernameBySecureEmail 根据邮箱找回用户名
func (j *JinmuIDService) FindUsernameBySecureEmail(ctx context.Context, req *proto.FindUsernameBySecureEmailRequest, resp *proto.FindUsernameBySecureEmailResponse) error {
	username, errFindUsernameBySecureEmail := j.datastore.FindUsernameBySecureEmail(ctx, req.Email)
	if errFindUsernameBySecureEmail != nil {
		return NewError(ErrNonexistentUsername, fmt.Errorf("failed to find username by secure email %s: %s", req.Email, errFindUsernameBySecureEmail.Error()))
	}

	// 查找可用的最新的验证码记录
	latestRecord, errFindVcRecord := j.datastore.FindLatestVcRecord(ctx, req.Email, mysqldb.FindUsername)
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
		return NewError(ErrDatabase, fmt.Errorf("failed to modify Vc record status of record %d: %s", latestRecord.RecordID, err.Error()))
	}
	resp.Username = username
	return nil
}
