package rest

import (
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	"github.com/kataras/iris/v12"
)

// SetUserPasswordBody 设置用户密码body
type SetUserPasswordBody struct {
	PlainPassword string `json:"plain_password"`
}

// ModifyUserPasswordBody 修改用户密码body
type ModifyUserPasswordBody struct {
	OldHashedPassword string `json:"old_hashed_password"`
	Seed              string `json:"seed"`
	NewPlainPassword  string `json:"new_plain_password"`
}

// JinmuService 金姆服务
type JinmuService struct {
	Service           string `json:"service"`
	ServiceDescrption string `json:"service_descrption"`
}

// SetUserPassword 设置密码
func (h *webHandler) SetUserPassword(ctx iris.Context) {
	userID, err := ctx.Params().GetInt("user_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	req := new(proto.UserSetPasswordRequest)
	req.UserId = int32(userID)
	var body SetUserPasswordBody
	errReadJSON := ctx.ReadJSON(&body)
	if errReadJSON != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", errReadJSON), false)
		return
	}
	req.PlainPassword = body.PlainPassword
	_, errUserSetPassword := h.rpcSvc.UserSetPassword(newRPCContext(ctx), req)
	if errUserSetPassword != nil {
		writeRpcInternalError(ctx, errUserSetPassword, false)
		return
	}
	rest.WriteOkJSON(ctx, nil)
}

// ModifyUserPassword 修改密码
func (h *webHandler) ModifyUserPassword(ctx iris.Context) {
	userID, err := ctx.Params().GetInt("user_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	var body ModifyUserPasswordBody
	errReadJSON := ctx.ReadJSON(&body)
	if errReadJSON != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", errReadJSON), false)
		return
	}
	req := new(proto.UserModifyPasswordRequest)
	req.UserId = int32(userID)
	req.OldHashedPassword = body.OldHashedPassword
	req.NewPlainPassword = body.NewPlainPassword
	req.Seed = body.Seed
	_, errUserModifyPassword := h.rpcSvc.UserModifyPassword(newRPCContext(ctx), req)
	if errUserModifyPassword != nil {
		writeRpcInternalError(ctx, errUserModifyPassword, false)
		return
	}
	rest.WriteOkJSON(ctx, nil)
}

// GetUserUsingServices 得到用户使用的服务
func (h *webHandler) GetUserUsingServices(ctx iris.Context) {
	userID, err := ctx.Params().GetInt("user_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	req := new(proto.UserGetUsingServiceRequest)
	req.UserId = int32(userID)
	resp, errUserGetUsingService := h.rpcSvc.UserGetUsingService(newRPCContext(ctx), req)
	if errUserGetUsingService != nil {
		writeRpcInternalError(ctx, errUserGetUsingService, false)
		return
	}
	jinmuServices := make([]JinmuService, len(resp.Clients))
	for idx, client := range resp.Clients {
		jinmuServices[idx] = JinmuService{
			Service:           client.Remark,
			ServiceDescrption: client.Usage,
		}
	}
	rest.WriteOkJSON(ctx, jinmuServices)
}
