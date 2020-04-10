package rest

import (
	jwtmiddleware "github.com/jinmukeji/gf-api2/pkg/rest/jwt"
	analysispb "github.com/jinmukeji/proto/gen/micro/idl/jm/analysis/v1"
	corepb "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	jinmuidpb "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	subscriptionpb "github.com/jinmukeji/proto/gen/micro/idl/jm/subscription/v1"
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
	rpcServiceName             = "com.jinmuhealth.srv.svc-biz-core"
	rpcSubscriptionServiceName = "com.jinmuhealth.srv.svc-subscription"
	rpcJinmuidServiceName      = "com.jinmuhealth.srv.svc-jinmuid"
	rpcAnalysisServiceName     = "com.jinmuhealth.srv.svc-analysis"
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
