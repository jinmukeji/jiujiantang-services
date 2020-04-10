package rest

import (
	"context"
	"net/http"

	"github.com/jinmukeji/gf-api2/pkg/rest"
	"github.com/kataras/iris/v12"
	"github.com/micro/go-micro/metadata"
)

const (
	// userTokenKey 用于从 context 中获取和设置用户登录凭证
	userTokenKey = "Access-Token"
	// ClientZoneKey 用于从 Context 的 Metadata中获取和设置Zone
	ClientZoneKey = "ClientZone"
	// ClientNameKey 用于从 Context 的 Metadata中获取和设置Name
	ClientNameKey = "ClientName"
	// ClientIDKey 用于从 Context 的 Metadata中获取和设置ClientID
	ClientIDKey = "ClientID"
	// ClientCustomizedCodeKey 用于从 Context 的 Metadata中获取和设置CustomizedCode
	ClientCustomizedCodeKey = "ClientCustomizedCode"
	// RemoteClientIPKey 用于从 Context 的 Metadata中获取和设置Client的IP地址
	RemoteClientIPKey = "RemoteClientIP"
)

// newRPCContext 得到 RPC 的 Context
func newRPCContext(ctx iris.Context) context.Context {
	return metadata.NewContext(ctx.Request().Context(), map[string]string{
		// go 底层源码里面对 Key 传递的时候做了 CanonicalMIMEHeaderKey 处理
		http.CanonicalHeaderKey(userTokenKey):            ctx.GetHeader("X-Access-Token"),
		http.CanonicalHeaderKey(ClientIDKey):             ctx.Values().GetString(ClientIDKey),
		http.CanonicalHeaderKey(ClientZoneKey):           ctx.Values().GetString(ClientZoneKey),
		http.CanonicalHeaderKey(ClientCustomizedCodeKey): ctx.Values().GetString(ClientCustomizedCodeKey),
		http.CanonicalHeaderKey(ClientNameKey):           ctx.Values().GetString(ClientNameKey),
		http.CanonicalHeaderKey(RemoteClientIPKey):       ctx.Values().GetString(RemoteClientIPKey),
		http.CanonicalHeaderKey(rest.ContextCidKey):      ctx.Values().GetString(rest.ContextCidKey),
	})
}
