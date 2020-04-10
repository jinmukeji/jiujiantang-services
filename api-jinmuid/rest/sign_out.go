package rest

import (
	"github.com/jinmukeji/gf-api2/pkg/rest"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
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
