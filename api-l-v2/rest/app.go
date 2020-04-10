package rest

import (
	"errors"
	"net/http"

	r "github.com/jinmukeji/jiujiantang-services/pkg/rest"
	"github.com/kataras/iris/v12/context"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/iris-contrib/middleware/cors"
	jwtmiddleware "github.com/jinmukeji/jiujiantang-services/pkg/rest/jwt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
)

// NewApp 创建一个实现了 http.Handler 接口的应用程序
func NewApp(base string, jwtSignInKey string) http.Handler {
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
	app.OnErrorCode(iris.StatusInternalServerError, r.InternalServerError)
	// jwt配置
	jwtHandler := jwtmiddleware.New(jwtmiddleware.Config{
		// 这个方法将验证jwt的token
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSignInKey), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
		ErrorHandler:  OnError,
		Expiration:    true,
		JwtSignInKey:  jwtSignInKey,
	})
	// 设置路由
	h := newV2Handler(jwtHandler)
	v2API := app.Party("/"+base, crs).AllowMethods(iris.MethodOptions)
	v2API.Post("/client/auth", h.ClientAuth)
	v2API.Get("/version", h.GetVersion)
	v2API.Post("/account/signin", jwtHandler.Serve, h.parseJWTToken, h.JinmuLAccountSignIn)
	v2API.Post("/res/get_url", h.GetJMResBaseURL)
	v2API.Get("/wxmp/qrcode", h.GetWxmpTempQrCodeUrl)
	v2API.Post("/organizations/{organization_id:int}/devices", jwtHandler.Serve, h.parseJWTToken, h.JinmuLOwnerBulkBindDevices)
	v2API.Post("/organizations/{organization_id:int}/devices/delete", jwtHandler.Serve, h.parseJWTToken, h.JinmuLOwnerBulkUnbindDevices)
	v2API.Get("/organizations/{organization_id:int}/devices", jwtHandler.Serve, h.parseJWTToken, h.JinmuLOwnerGetOrganizationDeviceList)
	v2API.Post("/organizations/{organization_id:int}/devices/is_bindable", jwtHandler.Serve, h.parseJWTToken, h.JinmuLOwnerCheckDeviceBindable)
	v2API.Put("/users/{user_id:int}/profile", jwtHandler.Serve, h.parseJWTToken, h.JinmuLModifyUserProfile)
	v2API.Get("/users/{user_id:int}/profile", jwtHandler.Serve, h.parseJWTToken, h.JinmuLGetUserProfile)
	v2API.Get("/account/{account:string}/signout", jwtHandler.Serve, h.parseJWTToken, h.JinmuLAccountSignOut)
	v2API.Post("/measurements", jwtHandler.Serve, h.parseJWTToken, h.setClientIP, h.SubmitMeasurementData)
	v2API.Post("/measurements/{record_id:int}/analyze", jwtHandler.Serve, h.parseJWTToken, h.GetAnalyzeResult)
	v2API.Get("/measurements/{record_id:int}/analyze", jwtHandler.Serve, h.parseJWTToken, h.GetAnalyzeReport)
	if err := app.Build(); err != nil {
		log.Fatal(err)
	}
	return app
}

// parseJWTToken 解析JWT Token
func (h *v2Handler) parseJWTToken(ctx iris.Context) {
	userToken := h.jwtMiddleware.Get(ctx)
	var clientID, clientZone, clientCustomizedCode, clientName string
	if claims, ok := userToken.Claims.(jwt.MapClaims); ok && userToken.Valid {
		clientID = claims["client_id"].(string)
		clientZone = claims["zone"].(string)
		clientCustomizedCode = claims["customized_code"].(string)
		clientName = claims["name"].(string)
	}
	ctx.Values().Set(ClientIDKey, clientID)
	ctx.Values().Set(ClientZoneKey, clientZone)
	ctx.Values().Set(ClientNameKey, clientName)
	ctx.Values().Set(ClientCustomizedCodeKey, clientCustomizedCode)
	ctx.Next()
}

// onErrorDataArrayFormat 错误data数组格式
const onErrorDataArrayFormat = "onErrorDataArrayFormat"

// OnError 错误的执行路径
func OnError(ctx context.Context, err string) {
	isArray := ctx.Values().GetBoolDefault(onErrorDataArrayFormat, false)
	writeError(ctx, wrapError(ErrClientUnauthorized, "", errors.New(err)), isArray)
}

// setClientIP 在context中设置Client IP
func (h *v2Handler) setClientIP(ctx iris.Context) {
	ctx.Values().Set(RemoteClientIPKey, ctx.RemoteAddr())
	ctx.Next()
}
