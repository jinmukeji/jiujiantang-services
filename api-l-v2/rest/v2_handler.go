package rest

import (
	jwtmiddleware "github.com/jinmukeji/jiujiantang-services/pkg/rest/jwt"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	"github.com/micro/go-micro/client"
)

type v2Handler struct {
	rpcSvc        proto.JinmuhealthAPIService
	jwtMiddleware *jwtmiddleware.Middleware
}

const (
	rpcServiceName = "com.himalife.srv.svc-biz-core"
)

func newV2Handler(jwtMiddleware *jwtmiddleware.Middleware) *v2Handler {
	return &v2Handler{
		rpcSvc:        proto.NewJinmuhealthAPIService(rpcServiceName, client.DefaultClient),
		jwtMiddleware: jwtMiddleware,
	}
}
