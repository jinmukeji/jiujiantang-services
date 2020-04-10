package handler

import (
	"errors"

	"github.com/jinmukeji/gf-api2/service/auth"
	"github.com/jinmukeji/gf-api2/service/mysqldb"
	"github.com/micro/go-micro/server"
	"golang.org/x/net/context"
)

// AccessTokenTypeLValue 一体机
const AccessTokenTypeLValue = "JinmuL"

// AccessTokenTypeWeChatValue 微信
const AccessTokenTypeWeChatValue = "Wechat"

// contextKey 用于从获取上下文数据
type contextKey string

// clientKey 用于从 context 中获取和设置客户端授权信息 {clientID, name, zone, expiredAt}
const clientKey contextKey = "X-Client"

// WhiteList 白名单服务，无需身份信息
type WhiteList map[string]bool

var (
	// 白名单，客户端不需要设置client信息
	clientAuthWhiteList = WhiteList{
		"JinmuhealthAPI.ClientAuth":              true,
		"JinmuhealthAPI.Echo":                    true,
		"JinmuhealthAPI.GetVersion":              true,
		"JinmuhealthAPI.GetFAQBaseUrl":           true,
		"JinmuhealthAPI.ScanQRCode":              true,
		"JinmuhealthAPI.GetAnalyzeResultByToken": true,
		"JinmuhealthAPI.GetJMResBaseUrl":         true,
		"JinmuhealthAPI.GetLocalNotifications":   true,
	}
	// 白名单，用户不需要经过用户登录认证
	tokenWhiteList = WhiteList{
		"JinmuhealthAPI.Echo":                    true,
		"JinmuhealthAPI.GetVersion":              true,
		"JinmuhealthAPI.ClientAuth":              true,
		"JinmuhealthAPI.UserSignIn":              true,
		"JinmuhealthAPI.AccountLogin":            true,
		"JinmuhealthAPI.GetFAQBaseUrl":           true,
		"JinmuhealthAPI.ScanQRCode":              true,
		"JinmuhealthAPI.CreateWxUser":            true,
		"JinmuhealthAPI.JinmuLAccountLogin":      true,
		"JinmuhealthAPI.GetAnalyzeResultByToken": true,
		"JinmuhealthAPI.GetJMResBaseUrl":         true,
		"JinmuhealthAPI.GetLocalNotifications":   true,
	}
)

// metaClient 客户端信息
type metaClient struct {
	ClientID       string
	Zone           string
	CustomizedCode string
	Name           string
	RemoteClientIP string
}

// AuthenticationWrapper 是 HandleWrapper 的 factory
type AuthenticationWrapper struct {
	datastore mysqldb.Datastore
}

// SetDataStore 设置数据库
func (w *AuthenticationWrapper) SetDataStore(datastore mysqldb.Datastore) {
	w.datastore = datastore
}

// AuthWrapper 设置Client
func (w *AuthenticationWrapper) AuthWrapper() server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			accessTokenType, _ := auth.AccessTokenTypeFromContext(ctx)
			if accessTokenType == AccessTokenTypeWeChatValue {
				return fn(ctx, req, rsp)
			}
			// 白名单内的 method 不需要设置client信息
			if _, ok := clientAuthWhiteList[req.Method()]; ok {
				return fn(ctx, req, rsp)
			}
			clientID, _ := auth.ClientIDFromContext(ctx)
			zone, _ := auth.ZoneFromContext(ctx)
			name, _ := auth.NameFromContext(ctx)
			customizedCode, _ := auth.CustomizedCodeFromContext(ctx)
			remoteClientIP, _ := auth.RemoteClientIPFromContext(ctx)
			ctx = addContextClient(ctx, metaClient{
				ClientID:       clientID,
				Zone:           zone,
				Name:           name,
				CustomizedCode: customizedCode,
				RemoteClientIP: remoteClientIP,
			})
			// 交给下一个 handler
			return fn(ctx, req, rsp)
		}
	}
}

// HandleWrapper 生成登录验证 wrapper
func (w *AuthenticationWrapper) HandleWrapper() server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			// 白名单内的 method 不需要认证
			if _, ok := tokenWhiteList[req.Method()]; ok {
				return fn(ctx, req, rsp)
			}
			accessTokenType, _ := auth.AccessTokenTypeFromContext(ctx)
			switch accessTokenType {
			case AccessTokenTypeWeChatValue:
				return fn(ctx, req, rsp)
			case AccessTokenTypeLValue:
				token, ok := auth.TokenFromContext(ctx)
				if !ok {
					return NewError(ErrUserUnauthorized, errors.New("token错误 认证失败"))
				}
				account, err := w.datastore.FindJinmuLAccountByToken(ctx, token)
				if err != nil || account == "" {
					// TODO: 后台加入 token 清理逻辑
					return NewError(ErrUserUnauthorized, errors.New("登录授权已失效，请重新登录"))
				}
				ctx = auth.AddContextAccount(ctx, account)
				machineUUID, err := w.datastore.FindMachineUUIDByToken(ctx, token)
				if err != nil || machineUUID == "" {
					// TODO: 后台加入 token 清理逻辑
					return NewError(ErrUserUnauthorized, errors.New("登录授权已失效，请重新登录"))
				}
				// 把  uuid 放入 context
				ctx = auth.AddContextMachineUUID(ctx, machineUUID)
				return fn(ctx, req, rsp)
			default:
				token, ok := auth.TokenFromContext(ctx)
				if !ok {
					return NewError(ErrUserUnauthorized, errors.New("token错误 认证失败"))
				}
				userID, err := w.datastore.FindUserIDByToken(ctx, token)
				if err != nil || userID == 0 {
					// TODO: 后台加入 token 清理逻辑
					return NewError(ErrUserUnauthorized, errors.New("登录授权已失效，请重新登录"))
				}
				// 把  userID 放入 context
				ctx = auth.AddContextUserID(ctx, userID)
				// 交给下一个 handler
				return fn(ctx, req, rsp)
			}
		}
	}
}

// clientFromContext 从上下文获取客户端信息
func clientFromContext(ctx context.Context) (metaClient, bool) {
	client, ok := ctx.Value(clientKey).(metaClient)
	return client, ok
}

// addContextClient 往上下文加入客户端信息
func addContextClient(ctx context.Context, client metaClient) context.Context {
	return context.WithValue(ctx, clientKey, client)
}
