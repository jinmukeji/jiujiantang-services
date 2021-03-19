package rest

import (
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	"github.com/kataras/iris/v12"
)

// GetVersion 获取服务版本信息
func (h *handler) GetVersion(ctx iris.Context) {
	rest.WriteOkJSON(ctx, iris.Map{
		"version": "2.0.0",
	})
}
