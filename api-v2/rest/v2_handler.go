package rest

import (
	jwtmiddleware "github.com/jinmukeji/jiujiantang-services/pkg/rest/jwt"
	jinmuidpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	analysispb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/analysis/v1"
	corepb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	subscriptionpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/subscription/v1"
	"github.com/micro/go-micro/client"
)

type v2Handler struct {
	rpcSvc                    corepb.JinmuhealthAPIService
	jwtMiddleware             *jwtmiddleware.Middleware
	rpcSubscriptionManagerSvc subscriptionpb.SubscriptionManagerAPIService
	rpcJinmuidSvc             jinmuidpb.UserManagerAPIService
	rpcAnalysisSvc            analysispb.AnalysisManagerAPIService
}

const (
	rpcServiceName             = "com.himalife.srv.svc-biz-core"
	rpcSubscriptionServiceName = "com.himalife.srv.svc-subscription"
	rpcJinmuidServiceName      = "com.himalife.srv.svc-jinmuid"
	rpcAnalysisServiceName     = "com.himalife.srv.svc-analysis"
)

func newV2Handler(jwtMiddleware *jwtmiddleware.Middleware) *v2Handler {
	return &v2Handler{
		rpcSvc:                    corepb.NewJinmuhealthAPIService(rpcServiceName, client.DefaultClient),
		jwtMiddleware:             jwtMiddleware,
		rpcSubscriptionManagerSvc: subscriptionpb.NewSubscriptionManagerAPIService(rpcSubscriptionServiceName, client.DefaultClient),
		rpcJinmuidSvc:             jinmuidpb.NewUserManagerAPIService(rpcJinmuidServiceName, client.DefaultClient),
		rpcAnalysisSvc:            analysispb.NewAnalysisManagerAPIService(rpcAnalysisServiceName, client.DefaultClient),
	}
}
