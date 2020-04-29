package rest

import (
	"github.com/jinmukeji/jiujiantang-services/pkg/rest.v3"
	"github.com/jinmukeji/jiujiantang-services/web-go/config"
	"github.com/kataras/iris/v12"
)

// VersionResponse 版本的返回
type VersionResponse struct {
	ServiceName string `json:"service_name"`
	Version     string `json:"version"`
}

// version 版本
var version = VersionResponse{
	ServiceName: config.FullServiceName(),
	Version:     config.ProductVersion,
}

// GetVersion 获取服务版本信息
func (h *webHandler) GetVersion(ctx iris.Context) {
	rest.WriteOkJSON(ctx, version)
}
