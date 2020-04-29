package rest

import (
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	"github.com/kataras/iris/v12"
)

// 注销登录
func (h *webHandler) SignOut(ctx iris.Context) {
	req := new(proto.UserSignOutRequest)
	req.Ip = ctx.RemoteAddr()
	_, err := h.rpcSvc.UserSignOut(
		newRPCContext(ctx), req,
	)

	if err != nil {
		writeRpcInternalError(ctx, err, false)
		return
	}
	rest.WriteOkJSON(ctx, nil)
}
