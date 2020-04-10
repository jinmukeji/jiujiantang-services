package rest

import (
	"github.com/jinmukeji/gf-api2/pkg/rest"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	"github.com/kataras/iris/v12"
)

// GerResourceList 资源列表
type GerResourceList struct {
	Languages []string `json:"languages"` // 语言列表
	Regions   []string `json:"regions"`   // 地区列表
}

func (h *webHandler) GerResourceList(ctx iris.Context) {
	req := new(proto.GerResourceListRequest)
	resp, errGerResourceList := h.rpcSvc.GerResourceList(newRPCContext(ctx), req)
	if errGerResourceList != nil {
		writeRpcInternalError(ctx, errGerResourceList, false)
		return
	}
	rest.WriteOkJSON(ctx, GerResourceList{
		Languages: resp.Languages,
		Regions:   resp.Regions,
	})
}
