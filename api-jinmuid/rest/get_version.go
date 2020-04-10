package rest

import (
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	"github.com/kataras/iris/v12"
)

// GetVersion 获取服务版本信息
func (h *webHandler) GetVersion(ctx iris.Context) {
	resp, err := h.rpcSvc.GetVersion(newRPCContext(ctx), &proto.GetVersionRequest{})
	if err != nil {
		writeRpcInternalError(ctx, err, false)
		return
	}
	rest.WriteOkJSON(ctx, resp)
}
