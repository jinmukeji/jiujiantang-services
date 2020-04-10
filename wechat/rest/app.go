package rest

import (
	"net/http"

	"github.com/iris-contrib/middleware/cors"
	r "github.com/jinmukeji/gf-api2/pkg/rest"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
)

// NewApp 创建一个实现了 http.Handler 接口的应用程序
func NewApp(ops *Options) http.Handler {
	app := iris.New().
		Configure(iris.WithRemoteAddrHeader("X-Forwarded-For"))
	crs := cors.AllowAll()

	// 配置自定义日志中间件
	requestLogger := logger.New(logger.Config{
		// Status displays status code
		Status: true,
		// IP displays request's remote address
		IP: true,
		// Method displays the http method
		Method: true,
		// Path displays the request path
		Path: true,

		LogFunc: logFunc,
	})
	app.Use(requestLogger)
	app.UseGlobal(r.CidMiddleware, r.LogMiddleware)
	app.OnErrorCode(iris.StatusNotFound, r.NotFound)
	app.OnErrorCode(iris.StatusInternalServerError, internalServerError)

	// 设置路由
	h := newHandler(ops)
	api := app.Party("/"+ops.APIBase, crs).AllowMethods(iris.MethodOptions)
	api.Get("/version", h.GetVersion)

	api.Get("/wx", h.WeChatAuth)
	api.Post("/wx", h.receiveMessage)

	api.Get("/wx/oauth", h.WeChatOAuth)
	api.Get("/wx/oauth/callback", h.WeChatOAuthCallback)

	api.Get("/wx/api/measurements", h.SearchHistory)
	api.Put("/wx/api/measurements/{record_id:int}/remark", h.SubmitRemark)
	api.Get("/wx/api/measurements/{record_id:int}/analyze", h.GetAnalyzeReport)
	api.Post("/wx/api/payment", h.MakePayment)
	api.Post("/wx/api/jssdk/config", h.GetWxJsSdkConfig)
	api.Get("/wx/api/measurements/{record_id:int}/token", h.GetRecordToken)
	api.Get("/wx/api/measurements/token/{token:string}/analyze", h.GetAnalyzeReportByToken)
	api.Post("/wx/api/res/getUrl", h.GetJMResBaseURL)
	if err := app.Build(); err != nil {
		log.Fatal(err)
	}

	return app
}
