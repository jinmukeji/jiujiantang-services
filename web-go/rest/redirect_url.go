package rest

import (
	"fmt"

	"github.com/kataras/iris/v12"
)

const (
	// CategoryE 词条种类
	CategoryE = "e"
	// EnvironmentProduction 正式环境
	EnvironmentProduction = "production"
)

// RedirectResURL 获取服务版本信息
func (h *webHandler) RedirectResURL(ctx iris.Context) {
	// url params
	category := ctx.Params().Get("category")
	if category == CategoryE {
		key := ctx.Params().Get("key")
		env := ctx.URLParam("env")
		if env == "" {
			env = EnvironmentProduction
		}
		if h.res.Entry[env] == nil {
			statusNotFound(ctx)
			return
		}
		url := h.res.Entry[env][key]
		if url != "" {
			ctx.Redirect(url, iris.StatusFound)
			return
		}
	}
	statusNotFound(ctx)
}

func statusNotFound(ctx iris.Context) {
	ctx.StatusCode(iris.StatusNotFound)
	_, errText := ctx.Text(fmt.Sprintf("%d: Not found", iris.StatusNotFound))
	if errText != nil {
		return
	}
}
