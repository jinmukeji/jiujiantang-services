package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/protobuf/ptypes"

	"github.com/jinmukeji/jiujiantang-services/service/auth"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
)

// JinmuLAccountLogin 用户登录
func (j *JinmuHealth) JinmuLAccountLogin(ctx context.Context, req *proto.JinmuLAccountLoginRequest, resp *proto.JinmuLAccountLoginResponse) error {
	jinmuLAccount, err := j.datastore.FindJinmuLAccount(ctx, req.Account)
	if err != nil {
		return err
	}
	if jinmuLAccount.Password != req.Password {
		return NewError(ErrUsernamePasswordNotMatch, errors.New("Username and password does not match"))
	}
	token := auth.GenerateToken()
	if req.MachineId == "" {
		return NewError(ErrNullMachineUUID, errors.New("machine_uuid is nil"))
	}
	tk, errCreateAccessToken := j.datastore.CreateAccessToken(ctx, token, req.Account, req.MachineId, TokenAvailableDuration)
	if errCreateAccessToken != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to create accesstoken: %s", errCreateAccessToken.Error()))
	}
	resp.AccessToken = tk.Token
	resp.OrganizationId = jinmuLAccount.OrganizationID
	expiredAt, errTimestampProto := ptypes.TimestampProto(tk.ExpiredAt)
	if errTimestampProto != nil {
		return NewError(ErrProtoConversionFailure, fmt.Errorf("failed to parse timestamp at %s", errTimestampProto.Error()))
	}
	resp.ExpiredTime = expiredAt
	return nil
}
