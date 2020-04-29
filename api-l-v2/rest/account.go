package rest

import (
	"time"

	"github.com/golang/protobuf/ptypes"

	"fmt"

	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	"github.com/kataras/iris/v12"
)

// AccountSignIn 账户登陆
type AccountSignIn struct {
	Account     string `json:"account"`      // 账户
	Password    string `json:"password"`     // 密码
	MachineUUID string `json:"machine_uuid"` // 机器的uuid
}

// AccountSignInResponse 账户登陆的响应
type AccountSignInResponse struct {
	AccessToken    string    `json:"access_token"` // access_token
	OrganizationID int32     `json:"organization_id"`
	ExpiredAt      time.Time `json:"expired_at"`
}

// JinmuLAccountSignIn 账户登录
func (h *v2Handler) JinmuLAccountSignIn(ctx iris.Context) {
	var accountSignIn AccountSignIn
	err := ctx.ReadJSON(&accountSignIn)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}
	req := new(proto.JinmuLAccountLoginRequest)
	req.Account = accountSignIn.Account
	req.Password = accountSignIn.Password
	req.MachineId = accountSignIn.MachineUUID
	resp, err := h.rpcSvc.JinmuLAccountLogin(
		newRPCContext(ctx), req,
	)
	if err != nil {
		writeRpcInternalError(ctx, err, false)
		return
	}
	expiredAt, errTimestamp := ptypes.Timestamp(resp.ExpiredTime)
	if errTimestamp != nil {

		writeError(ctx, wrapError(ErrRPCInternal, "", fmt.Errorf("failed to get timestamp of %v: %s", resp.ExpiredTime, errTimestamp.Error())), false)
		return
	}
	rest.WriteOkJSON(ctx, AccountSignInResponse{
		AccessToken:    resp.AccessToken,
		OrganizationID: resp.OrganizationId,
		ExpiredAt:      expiredAt,
	})
}

// JinmuLAccountSignOut 账户登出
func (h *v2Handler) JinmuLAccountSignOut(ctx iris.Context) {
	req := new(proto.JinmuLAccountSignOutRequest)
	_, err := h.rpcSvc.JinmuLAccountSignOut(
		newRPCContext(ctx), req,
	)
	if err != nil {
		writeRpcInternalError(ctx, err, false)
		return
	}
	rest.WriteOkJSON(ctx, nil)
}
