package rest

import (
	"errors"
	"net/http"
	"time"

	r "github.com/jinmukeji/jiujiantang-services/pkg/rest"
	jwtmiddleware "github.com/jinmukeji/jiujiantang-services/pkg/rest/jwt"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

// NewApp 创建一个实现了 http.Handler 接口的应用程序
func NewApp(base string, jwtSignInKey string, debug bool) http.Handler {
	app := iris.New().
		Configure(iris.WithRemoteAddrHeader("X-Forwarded-For"))
	crs := cors.AllowAll()

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
		JwtDuration:   time.Minute * 10,
	})
	// 设置路由
	h := newWebHandler(jwtHandler)
	webAPI := app.Party("/"+base, crs).AllowMethods(iris.MethodOptions)
	webAPI.Get("/version", h.GetVersion)
	webAPI.Post("/client/auth", h.ClientAuth)
	webAPI.Post("/signin", jwtHandler.Serve, h.parseJWTToken, h.SignIn)
	webAPI.Post("/user/{user_id:int}/language", jwtHandler.Serve, h.parseJWTToken, h.SetWebLanguage)
	webAPI.Get("/user/{user_id:int}/language", jwtHandler.Serve, h.parseJWTToken, h.GetWebLanguage)
	webAPI.Get("/user/{user_id:int}/profile", jwtHandler.Serve, h.parseJWTToken, h.GetUserProfile)
	webAPI.Put("/user/{user_id:int}/profile", jwtHandler.Serve, h.parseJWTToken, h.ModifyUserProfile)
	webAPI.Post("/user/{user_id:int}/password", jwtHandler.Serve, h.parseJWTToken, h.SetUserPassword)
	webAPI.Put("/user/{user_id:int}/password", jwtHandler.Serve, h.parseJWTToken, h.ModifyUserPassword)
	webAPI.Post("/notification/sms", jwtHandler.Serve, h.parseJWTToken, h.SmsNotification)
	webAPI.Post("/signup/verification_number", jwtHandler.Serve, h.parseJWTToken, h.SignUpByPhoneVerificationNumber)
	webAPI.Get("/user/{user_id:int}/preferences", jwtHandler.Serve, h.parseJWTToken, h.GetUserPreferences)
	webAPI.Post("/user/signout", jwtHandler.Serve, h.parseJWTToken, h.SignOut)
	webAPI.Post("/notification/email", jwtHandler.Serve, h.parseJWTToken, h.EmailNotification)
	webAPI.Post("/user/{user_id:int}/safe_email", jwtHandler.Serve, h.parseJWTToken, h.SetSecureEmail)
	webAPI.Post("/user/{user_id:int}/safe_email/delete", jwtHandler.Serve, h.parseJWTToken, h.UnsetSecureEmail)
	webAPI.Post("/user/{user_id:int}/signin_phone", jwtHandler.Serve, h.parseJWTToken, h.SetSigninPhone)
	webAPI.Post("/user/region", jwtHandler.Serve, h.parseJWTToken, h.SelectRegion)
	webAPI.Post("/user/{user_id:int}/reset_password", jwtHandler.Serve, h.parseJWTToken, h.ResetPassword)
	webAPI.Post("/validate_signin_phone", jwtHandler.Serve, h.parseJWTToken, h.VerifySigninPhone)
	webAPI.Post("/user/{user_id:int}/set_secure_question", jwtHandler.Serve, h.parseJWTToken, h.SetSecureQuestions)
	webAPI.Post("/user/{user_id:int}/validate_question_before_modify", jwtHandler.Serve, h.parseJWTToken, h.ValidateSecureQuestionsBeforeModifyQuestions)
	webAPI.Post("/user/{user_id:int}/modify_secure_question", jwtHandler.Serve, h.parseJWTToken, h.ModifySecureQuestions)
	webAPI.Post("/user/validate_username_or_phone", jwtHandler.Serve, h.parseJWTToken, h.ValidateUsernameOrPhone)
	webAPI.Post("/user/validate_question_before_modify_password", jwtHandler.Serve, h.parseJWTToken, h.ValidateSecureQuestionsBeforeModifyPassword)
	webAPI.Post("/user/modify_password_via_question", jwtHandler.Serve, h.parseJWTToken, h.ResetPasswordViaSecureQuestions)
	webAPI.Get("/user/{user_id:int}/secure_question_list", jwtHandler.Serve, h.parseJWTToken, h.GetSecureQuestionList)
	webAPI.Put("/user/{user_id:int}/signin_phone", jwtHandler.Serve, h.parseJWTToken, h.ModifySigninPhone)
	if debug {
		webAPI.Post("/_debug/user/latest_verification_code", jwtHandler.Serve, h.parseJWTToken, h.GetLatestVerificationCodes)
	}
	webAPI.Get("/resource", h.GerResourceList)
	webAPI.Post("/user/validate_email_verification_code", jwtHandler.Serve, h.parseJWTToken, h.ValidateEmailVerificationCode)
	webAPI.Post("/user/validate_phone_verification_code", jwtHandler.Serve, h.parseJWTToken, h.ValidatePhoneVerificationCode)
	webAPI.Get("/user/{user_id}/secure_question_to_modify", h.GetSecureQuestionListToModify)
	webAPI.Post("/user/find_username_by_email", jwtHandler.Serve, h.parseJWTToken, h.FindUsernameBySecureEmail)
	webAPI.Post("/user/secure_question", jwtHandler.Serve, h.parseJWTToken, h.GetSecureQuestionsByPhoneOrUsername)
	webAPI.Post("/user/{user_id}/modify_secure_email", jwtHandler.Serve, h.parseJWTToken, h.ModifySecureEmail)
	webAPI.Get("/user/{user_id}/devices", jwtHandler.Serve, h.parseJWTToken, h.UserUsedDevice)
	webAPI.Get("/user/{user_id:int}/services", jwtHandler.Serve, h.parseJWTToken, h.GetUserUsingServices)
	webAPI.Get("/user/{user_id:int}/sign_in_machines", jwtHandler.Serve, h.parseJWTToken, h.UserGetSignInMachines)
	webAPI.Post("/user/{user_id:int}/notification_preferences", jwtHandler.Serve, h.parseJWTToken, h.ModifyNotificationPreferences)
	webAPI.Get("/user/{user_id:int}/notification_preferences", jwtHandler.Serve, h.parseJWTToken, h.GetNotificationPreferences)
	if err := app.Build(); err != nil {
		log.Fatal(err)
	}
	return app
}

// parseJWTToken 解析JWT Token
func (h *webHandler) parseJWTToken(ctx iris.Context) {
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

func onError(ctx context.Context, err string) {
	isArray := ctx.Values().GetBoolDefault(onErrorDataArrayFormat, false)
	writeError(ctx, wrapError(ErrClientUnauthorized, "", errors.New(err)), isArray)
}
