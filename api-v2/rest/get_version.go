package rest

import (
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	"github.com/kataras/iris/v12"
)

// GetVersion 获取服务版本信息
func (h *v2Handler) GetVersion(ctx iris.Context) {
	resp, err := h.rpcSvc.GetVersion(newRPCContext(ctx), &proto.GetVersionRequest{})
	if err != nil {
		writeRPCInternalError(ctx, err, false)
		return
	}

	rest.WriteOkJSON(ctx, resp)
}
