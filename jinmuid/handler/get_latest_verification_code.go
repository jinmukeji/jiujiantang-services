package handler

import (
	"context"
	"errors"

	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
)

// GetLatestVerificationCodes 获取最新验证码
func (j *JinmuIDService) GetLatestVerificationCodes(ctx context.Context, req *proto.GetLatestVerificationCodesRequest, resp *proto.GetLatestVerificationCodesResponse) error {
	latestVerificationCodes := make([]proto.LatestVerificationCode, len(req.SendTo))
	for idx, item := range req.SendTo {
		switch item.SendVia {
		case proto.SendVia_SEND_VIA_INVALID:
			return NewError(ErrWrongSendVia, errors.New("wrong send via"))
		case proto.SendVia_SEND_VIA_UNSET:
			return NewError(ErrWrongSendVia, errors.New("wrong send via"))
		case proto.SendVia_SEND_VIA_PHONE_SEND_VIA:
			if item.Email != "" {
				return NewError(ErrInvalidSendValue, errors.New("email should be empty when sending media is phone"))
			}

			latestVerificationCodes[idx].Phone = item.Phone
			latestVerificationCodes[idx].NationCode = item.NationCode
			latestVerificationCode, errSearchEarliestPhoneVerificationCode := j.datastore.SearchLatestPhoneVerificationCode(ctx, item.Phone, item.NationCode)
			if errSearchEarliestPhoneVerificationCode != nil {
				latestVerificationCodes[idx].VerificationCode = ""
			} else {
				latestVerificationCodes[idx].VerificationCode = latestVerificationCode
			}
		case proto.SendVia_SEND_VIA_USERNAME_SEND_VIA:
			if item.Phone != "" || item.NationCode != "" {
				return NewError(ErrInvalidSendValue, errors.New("phone should be empty when sending media is email"))
			}
			latestVerificationCodes[idx].Email = item.Email
			latestVerificationCode, errSearchEarliestEmailVerificationCode := j.datastore.SearchLatestEmailVerificationCode(ctx, item.Email)
			if errSearchEarliestEmailVerificationCode != nil {
				latestVerificationCodes[idx].VerificationCode = ""
			} else {
				latestVerificationCodes[idx].VerificationCode = latestVerificationCode
			}
		}
	}
	latestVerifications := make([]*proto.LatestVerificationCode, len(req.SendTo))
	for idx, item := range latestVerificationCodes {
		latestVerifications[idx] = &item
	}
	resp.LatestVerificationCodes = latestVerifications
	return nil
}
