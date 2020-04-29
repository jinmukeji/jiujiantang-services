package handler

import (
	"errors"

	"github.com/jinmukeji/jiujiantang-services/jinmuid/mysqldb"
	"github.com/micro/go-micro/v2/server"
	"golang.org/x/net/context"
)

// WhiteList 白名单服务，无需身份信息
type WhiteList map[string]bool

var (
	// 白名单，用户不需要经过用户登录认证
	tokenWhiteList = WhiteList{
		"UserManagerAPI.GetVersion":                                      true,
		"UserManagerAPI.ClientAuth":                                      true,
		"UserManagerAPI.UserSignInByPhonePassword":                       true,
		"UserManagerAPI.UserSignInByUsernamePassword":                    true,
		"UserManagerAPI.SmsNotification":                                 true,
		"UserManagerAPI.UserSignUpByPhone":                               true,
		"UserManagerAPI.EmailNotification":                               true,
		"UserManagerAPI.UserSignInByPhoneVC":                             true,
		"UserManagerAPI.UserValidateUsernameOrPhone":                     true,
		"UserManagerAPI.UserValidateSecureQuestionsBeforeModifyPassword": true,
		"UserManagerAPI.UserResetPasswordViaSecureQuestions":             true,
		"UserManagerAPI.NotLoggedInEmailNotification":                    true,
		"UserManagerAPI.GetLatestVerificationCodes":                      true,
		"UserManagerAPI.UserResetPassword":                               true,
		"UserManagerAPI.GerResourceList":                                 true,
		"UserManagerAPI.FindUsernameBySecureEmail":                       true,
		"UserManagerAPI.VerifyUserSigninPhone":                           true,
		"UserManagerAPI.ValidateEmailVerificationCode":                   true,
		"UserManagerAPI.GetSecureQuestionsByPhoneOrUsername":             true,
		"UserManagerAPI.GetUserProfileByRecordID":                        true,
		"UserManagerAPI.GetUserProfile":                                  true,
		"UserManagerAPI.ModifyUserInformation":                           true,
		"UserManagerAPI.GetUserAndProfileInformation":                    true,
		"UserManagerAPI.ModifyUserProfile":                               true,
		"UserManagerAPI.SignUpUserViaUsernamePassword":                   true,
		"UserManagerAPI.ValidatePhoneVerificationCode":                   true,
	}
)

// AuthenticationWrapper 是 HandleWrapper 的 factory
type AuthenticationWrapper struct {
	datastore mysqldb.Datastore
}

// SetDataStore 设置数据库
func (w *AuthenticationWrapper) SetDataStore(datastore mysqldb.Datastore) {
	w.datastore = datastore
}

// HandleWrapper 生成登录验证 wrapper
func (w *AuthenticationWrapper) HandleWrapper() server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			// 白名单内的 method 不需要认证
			if _, ok := tokenWhiteList[req.Method()]; ok {
				return fn(ctx, req, rsp)
			}
			token, ok := TokenFromContext(ctx)
			if !ok {
				return NewError(ErrUserUnauthorized, errors.New("token错误 认证失败"))
			}
			userID, err := w.datastore.FindUserIDByToken(ctx, token)
			if err != nil || userID == 0 {
				// TODO: 后台加入 token 清理逻辑
				return NewError(ErrUserUnauthorized, errors.New("failed to get user_id by token"))
			}
			// 交给下一个 handler
			return fn(ctx, req, rsp)
		}
	}
}
