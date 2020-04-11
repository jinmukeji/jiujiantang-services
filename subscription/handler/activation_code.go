package handler

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	ac "github.com/jinmukeji/jiujiantang-services/subscription/activation-code"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/subscription/v1"
)

// UseSubscriptionActivationCode 使用激活码
func (j *SubscriptionService) UseSubscriptionActivationCode(ctx context.Context, req *proto.UseSubscriptionActivationCodeRequest, resp *proto.UseSubscriptionActivationCodeResponse) error {
	// 激活码转成大写
	code := strings.ToUpper(req.Code)
	activationCode, err := j.datastore.GetSubscriptionActivationCodeInfo(ctx, code)
	if err != nil {
		return NewError(ErrInvalidActivationCode, fmt.Errorf("failed to get subscription activation code %s info: %s", code, err.Error()))
	}
	// 校验
	if !j.checkActivationCode(code, activationCode.ContractYear, activationCode.MaxUserLimits, activationCode.Checksum) {
		return NewError(ErrActivationCodeWrongChecksum, fmt.Errorf("checksum of code %s is wrong", code))
	}
	// 是否过期
	if activationCode.ExpiredAt.Before(time.Now()) {
		return NewError(ErrExpiredActivationCode, fmt.Errorf("activation code %s is expired", code))
	}
	// 是否售出
	if !activationCode.Sold {
		return NewError(ErrNotSoldActivationCode, fmt.Errorf("activation code %s is not sold", code))
	}
	// 是否已经激活
	if activationCode.Activated {
		return NewError(ErrActivatedActivationCode, fmt.Errorf("activation code %s is activated", code))
	}
	uuid, _ := uuid.NewUUID()
	errUseSubscriptionActivationCode := j.datastore.UseSubscriptionActivationCode(ctx, req.UserId, activationCode, uuid.String())
	if errUseSubscriptionActivationCode != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to use subscription activation code %s: %s", code, errUseSubscriptionActivationCode.Error()))
	}
	return nil
}

// GetSubscriptionActivationCodeInfo 得到激活码内容
func (j *SubscriptionService) GetSubscriptionActivationCodeInfo(ctx context.Context, req *proto.GetSubscriptionActivationCodeInfoRequest, resp *proto.GetSubscriptionActivationCodeInfoResponse) error {
	// 激活码转成大写
	code := strings.ToUpper(req.Code)
	activationCode, err := j.datastore.GetSubscriptionActivationCodeInfo(ctx, code)
	if err != nil {
		return NewError(ErrInvalidActivationCode, fmt.Errorf("activation code %s is invalid: %s", code, err.Error()))
	}
	// 校验
	if !j.checkActivationCode(code, activationCode.ContractYear, activationCode.MaxUserLimits, activationCode.Checksum) {
		return NewError(ErrActivationCodeWrongChecksum, fmt.Errorf("checksum of code %s is wrong", code))
	}
	// 是否过期
	if activationCode.ExpiredAt.Before(time.Now()) {
		return NewError(ErrExpiredActivationCode, fmt.Errorf("activation code %s is expired", code))
	}
	// 是否售出
	if !activationCode.Sold {
		return NewError(ErrNotSoldActivationCode, fmt.Errorf("activation code %s is not sold", code))
	}
	// 是否已经激活
	if activationCode.Activated {
		return NewError(ErrActivatedActivationCode, fmt.Errorf("activation code %s is activated", code))
	}
	resp.ContractYear = activationCode.ContractYear
	resp.MaxUserCount = activationCode.MaxUserLimits
	resp.Type = mapDBSubscriptionTypeToProto((activationCode.SubscriptionType))
	return nil
}

// checkActivationCode 校验激活码
func (j *SubscriptionService) checkActivationCode(code string, contractYear, maxUserLimits int32, checksum string) bool {
	helper := ac.NewActivationCodeCipherHelper()
	return code == helper.Decrypt(checksum, j.activationCodeEntryptKey, contractYear, maxUserLimits)
}
