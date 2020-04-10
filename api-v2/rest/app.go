package rest

import (
	"errors"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	r "github.com/jinmukeji/gf-api2/pkg/rest"
	jwtmiddleware "github.com/jinmukeji/gf-api2/pkg/rest/jwt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

// NewApp 创建一个实现了 http.Handler 接口的应用程序
func NewApp(base string, jwtSignInKey string) http.Handler {
	app := iris.New().
		Configure(iris.WithRemoteAddrHeader("X-Forwarded-For"))

	// 配置自定义日志中间件
	app.Logger().Install(log)
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
		ErrorHandler:  onError,
		Expiration:    true,
		JwtSignInKey:  jwtSignInKey,
	})
	// 设置路由
	h := newV2Handler(jwtHandler)
	v2API := app.Party("/" + base)
	v2API.Post("/users/signup", jwtHandler.Serve, h.parseJWTToken, h.UserSignUp)

	v2API.Get("/owner/{owner_id:int}/organizations/users", setDataArrayFormat, jwtHandler.Serve, h.parseJWTToken, h.OwnerGetOrganizationUsersByOwnerID)
	v2API.Get("/configuration/local_notification", h.GetLocalNotifications)
	v2API.Post("/client/auth", h.ClientAuth)
	v2API.Get("/version", h.GetVersion)
	v2API.Get("/tips", h.GetTips)
	v2API.Post("/feedback", jwtHandler.Serve, h.parseJWTToken, h.SubmitFeedback)
	v2API.Post("/users/signin", jwtHandler.Serve, h.parseJWTToken, h.UserSignIn)
	v2API.Post("/users/signout", jwtHandler.Serve, h.parseJWTToken, h.SignOut)
	v2API.Get("/owner/users/{user_id:int}/profile", jwtHandler.Serve, h.parseJWTToken, h.GetUserProfile)
	v2API.Put("/owner/users/{user_id:int}/profile", jwtHandler.Serve, h.parseJWTToken, h.ModifyUserProfile)
	v2API.Get("/owner/users/{user_id:int}/preferences", jwtHandler.Serve, h.parseJWTToken, h.OwnerGetUserPreferences)
	v2API.Post("/owner/users/signup", jwtHandler.Serve, h.parseJWTToken, h.OwnerSignUp)
	v2API.Post("/owner/organizations/{organization_id:int}/users", setDataArrayFormat, jwtHandler.Serve, h.parseJWTToken, jwtHandler.Serve, h.parseJWTToken, h.OwnerAddOrganizationUsers)
	v2API.Get("/owner/organizations/{organization_id:int}/users", setDataArrayFormat, jwtHandler.Serve, h.parseJWTToken, h.OwnerGetOrganizationUsers)
	v2API.Post("/owner/organizations/{organization_id:int}/users/delete", jwtHandler.Serve, h.parseJWTToken, h.OwnerDeleteOrganizationUsers)
	v2API.Post("/owner/organizations/{organization_id:int}/subscription", jwtHandler.Serve, h.parseJWTToken, h.GetOrganizationSubscription)
	v2API.Post("/owner/measurements", jwtHandler.Serve, h.parseJWTToken, h.setClientIP, h.SubmitMeasurementData)
	v2API.Get("/owner/measurements", setDataArrayFormat, jwtHandler.Serve, h.parseJWTToken, h.SearchHistory)
	v2API.Put("/owner/measurements/{record_id:int}/remark", jwtHandler.Serve, h.parseJWTToken, h.SubmitRemark)
	v2API.Post("/owner/measurements/{record_id:int}/analyze", jwtHandler.Serve, h.parseJWTToken, h.GetAnalyzeResult)
	v2API.Get("/owner/measurements/{record_id:int}/analyze", jwtHandler.Serve, h.parseJWTToken, h.GetAnalyzeReport)
	v2API.Post("/res/getUrl", h.GetJMResBaseURL)
	v2API.Get("/owner/organizations", setDataArrayFormat, jwtHandler.Serve, h.parseJWTToken, h.OwnerGetOrganizations)
	v2API.Get("/owner/measurements/{record_id:int}/token", jwtHandler.Serve, h.parseJWTToken, h.GetRecordToken)
	v2API.Get("/owner/measurements/token/{token:string}/analyze", h.GetAnalyzeReportByToken)
	v2API.Post("/owner/notification", jwtHandler.Serve, h.parseJWTToken, h.ReadPushNotification)
	v2API.Get("/owner/notification", setDataArrayFormat, jwtHandler.Serve, h.parseJWTToken, h.GetPushNotifications)
	v2API.Get("/bluetooth_name_prefix", jwtHandler.Serve, h.parseJWTToken, h.GetBluetoothNamePrefixes)
	v2API.Post("/activation_code", jwtHandler.Serve, h.parseJWTToken, h.GetActivationCodeInfo)
	v2API.Post("/owner/user/{user_id}/activation_code", jwtHandler.Serve, h.parseJWTToken, h.UseSubscriptionActivationCode)
	v2API.Get("/user/{user_id:int}/subscription", jwtHandler.Serve, h.parseJWTToken, h.GetUserSubscriptions)
	v2API.Get("/user/{user_id:int}/token", jwtHandler.Serve, h.parseJWTToken, h.GetLatestToken)
	v2API.Post("/user/measurements/{user_id:int}/delete", jwtHandler.Serve, h.parseJWTToken, h.DeleteRecords)
	v2API.Post("/owner/{owner_id:int}/users/sign_up", jwtHandler.Serve, h.parseJWTToken, h.OwnerUserSignUp)
	v2API.Post("/owner/{owner_id:int}/users/delete", jwtHandler.Serve, h.parseJWTToken, h.OwnerDeleteUsers)
	v2API.Post("/bind/{user_id:int}/old_user", jwtHandler.Serve, h.parseJWTToken, h.BindOldUser)

	v2API.Post("/owner/measurements/{record_id:int}/v2/analyze", jwtHandler.Serve, h.parseJWTToken, h.GetV2AnalyzeResult)
	v2API.Get("/owner/measurements/token/{token:string}/v2/analyze", h.GetV2AnalyzeReportByToken)
	v2API.Get("/owner/measurements/{record_id:int}/v2/analyze", jwtHandler.Serve, h.parseJWTToken, h.GetV2AnalyzeReportByRecordID)
	// 周报
	v2API.Post("/owner/{user_id:int}/measurements/v2/weekly_report", jwtHandler.Serve, h.parseJWTToken, h.GetWeeklyReport)
	// 月报
	v2API.Post("/owner/{user_id:int}/measurements/v2/monthly_report", jwtHandler.Serve, h.parseJWTToken, h.GetMonthlyReport)
	// 周趋势
	v2API.Get("/owner/{user_id:int}/week_measurements", setDataArrayFormat, jwtHandler.Serve, h.parseJWTToken, h.SearchWeekHistory)
	// 月趋势
	v2API.Get("/owner/{user_id:int}/month_measurements", setDataArrayFormat, jwtHandler.Serve, h.parseJWTToken, h.SearchMonthHistory)
	if err := app.Build(); err != nil {
		log.Fatal(err)
	}

	return app
}

// parseJWTToken 解析JWT Token
func (h *v2Handler) parseJWTToken(ctx iris.Context) {
	userToken := h.jwtMiddleware.Get(ctx)
	var clientID, zone, customizedCode, name string
	if claims, ok := userToken.Claims.(jwt.MapClaims); ok && userToken.Valid {
		clientID = claims["client_id"].(string)
		zone = claims["zone"].(string)
		customizedCode = claims["customized_code"].(string)
		name = claims["name"].(string)
	}
	ctx.Values().Set(ClientIDKey, clientID)
	ctx.Values().Set(ClientZoneKey, zone)
	ctx.Values().Set(ClientNameKey, name)
	ctx.Values().Set(ClientCustomizedCodeKey, customizedCode)
	ctx.Next()
}

// setClientIP 在context中设置Client IP
func (h *v2Handler) setClientIP(ctx iris.Context) {
	ctx.Values().Set(RemoteClientIPKey, ctx.RemoteAddr())
	ctx.Next()
}

// onErrorDataArrayFormat 错误data数组格式
const onErrorDataArrayFormat = "onErrorDataArrayFormat"

func setDataArrayFormat(ctx iris.Context) {
	ctx.Values().Set(onErrorDataArrayFormat, true)
	ctx.Next()
}

func onError(ctx context.Context, err string) {
	isArray := ctx.Values().GetBoolDefault(onErrorDataArrayFormat, false)
	writeError(ctx, wrapError(ErrClientUnauthorized, "", errors.New(err)), isArray)
}
