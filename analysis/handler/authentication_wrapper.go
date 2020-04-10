package handler

import (
	"errors"

	"github.com/jinmukeji/jiujiantang-services/analysis/mysqldb"
	"github.com/micro/go-micro/server"
	"golang.org/x/net/context"
)

// WhiteList 白名单服务，无需身份信息
type WhiteList map[string]bool

var (
	// 白名单，用户不需要经过用户登录认证
	tokenWhiteList = WhiteList{
		"AnalysisManagerAPI.GetVersion":                  true,
		"AnalysisManagerAPI.GetAnalyzeResultBodyByToken": true,
		"AnalysisManagerAPI.GetAnalyzeResult":            true,
		"AnalysisManagerAPI.GetAnalysisContent":          true,
		"AnalysisManagerAPI.UpdateAnalyzeRecord":         true,
		"AnalysisManagerAPI.UpdateAnalyzeStatus":         true,
		"AnalysisManagerAPI.GetAnalyzeResultByToken":     true,
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

// TODO: 在https://jinmuhealth.atlassian.net/browse/PLAT-274中处理
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
			userID, err := w.datastore.FindUserIDByToken(token)
			if err != nil || userID == 0 {
				// TODO: 后台加入 token 清理逻辑
				return NewError(ErrUserUnauthorized, errors.New("登录授权已失效，请重新登录"))
			}
			// 交给下一个 handler
			return fn(ctx, req, rsp)
		}
	}
}

// contextKey 用于从获取上下文数据
type contextKey string

// clientKey 用于从 context 中获取和设置客户端授权信息 {clientID, name, zone, expiredAt}
const clientKey contextKey = "X-Client"

// metaClient 客户端信息
type metaClient struct {
	ClientID       string
	Zone           string
	CustomizedCode string
	Name           string
	RemoteClientIP string
}

// clientFromContext 从上下文获取客户端信息
func clientFromContext(ctx context.Context) (metaClient, bool) {
	client, ok := ctx.Value(clientKey).(metaClient)
	return client, ok
}
