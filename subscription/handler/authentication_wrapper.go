package handler

import (
	"errors"

	"github.com/jinmukeji/jiujiantang-services/subscription/mysqldb"
	"github.com/micro/go-micro/server"
	"golang.org/x/net/context"
)

// WhiteList 白名单服务，无需身份信息
type WhiteList map[string]bool

var (
	// 白名单，用户不需要经过用户登录认证
	tokenWhiteList = WhiteList{
		"SubscriptionManagerAPI.GetVersion":                        true,
		"SubscriptionManagerAPI.GetSubscriptionActivationCodeInfo": true,
		"SubscriptionManagerAPI.GetSelectedUserSubscription":       true,
		"SubscriptionManagerAPI.CheckUserHaveSameSubscription":     true,
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
				return NewError(ErrUserUnauthorized, errors.New("登录授权已失效，请重新登录"))
			}
			// 交给下一个 handler
			return fn(ctx, req, rsp)
		}
	}
}
