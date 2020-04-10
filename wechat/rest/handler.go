package rest

import (
	"context"
	"net/http"

	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	"github.com/kataras/iris/v12"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
)

const (
	// AccessTokenTypeKey Access-token的type
	AccessTokenTypeKey = "Access-Token-Type"
	// AccessTokenTypeValue Access-token的value
	AccessTokenTypeValue = "Wechat"
)

type handler struct {
	rpcSvc               proto.JinmuhealthAPIService
	WxCallbackServerBase string
	WxH5ServerBase       string
}

const (
	rpcServiceName = "com.jinmuhealth.srv.svc-biz-core"
)

func newHandler(ops *Options) *handler {
	return &handler{
		rpcSvc:               proto.NewJinmuhealthAPIService(rpcServiceName, client.DefaultClient),
		WxCallbackServerBase: ops.WxCallbackServerBase,
		WxH5ServerBase:       ops.WxH5ServerBase,
	}
}

// newRPCContext 得到 RPC 的 Context
func newRPCContext(ctx iris.Context) context.Context {
	return metadata.NewContext(ctx.Request().Context(), map[string]string{
		// go 底层源码里面对 Key 传递的时候做了 CanonicalMIMEHeaderKey 处理
		http.CanonicalHeaderKey(AccessTokenTypeKey): AccessTokenTypeValue,
	})
}
