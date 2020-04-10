package rest

import (
	"net/http"

	"github.com/iris-contrib/middleware/cors"
	r "github.com/jinmukeji/jiujiantang-services/pkg/rest.v3"
	"github.com/kataras/iris/v12"
)

// NewApp 创建一个实现了 http.Handler 接口的应用程序
func NewApp(base, configFile string) http.Handler {
	app := iris.New().
		Configure(iris.WithRemoteAddrHeader("X-Forwarded-For"))
	crs := cors.AllowAll()

	// 配置自定义日志中间件
	app.Logger().Install(log)
	app.UseGlobal(r.CidMiddleware, r.LogMiddleware)
	app.OnErrorCode(iris.StatusNotFound, r.NotFound)
	app.OnErrorCode(iris.StatusInternalServerError, r.InternalServerError)
	// 设置路由
	h := newSysHandler(configFile)
	sysAPI := app.Party("/"+base, crs).AllowMethods(iris.MethodOptions)
	sysAPI.Post("/_sys/client_preferences", h.ClientPreferences)
	sysAPI.Get("/version", h.GetVersion)
	if err := app.Build(); err != nil {
		log.Fatal(err)
	}

	return app
}
