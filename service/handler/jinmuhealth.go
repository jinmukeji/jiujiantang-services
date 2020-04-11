package handler

import (
	ae "github.com/jinmukeji/ae-v1/core"
	"github.com/jinmukeji/jiujiantang-services/analysis/aws"
	"github.com/jinmukeji/jiujiantang-services/pkg/blocker"
	"github.com/jinmukeji/jiujiantang-services/service/mail"
	db "github.com/jinmukeji/jiujiantang-services/service/mysqldb"
	"github.com/jinmukeji/jiujiantang-services/service/wechat"
	jinmuidpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	analysispb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/analysis/v1"
	corepb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	devicepb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/device/v1"
	subscriptionpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/subscription/v1"
	calcpb "github.com/jinmukeji/proto/v3/gen/micro/idl/platform/calc/v2"
	"github.com/micro/go-micro/client"
)

const (
	rpcServiceName          = "com.himalife.srv.svc-biz-core"
	jinmuidServiceName      = "com.himalife.srv.svc-jinmuid"
	subscriptionServiceName = "com.himalife.srv.svc-subscription"
	deviceServiceName       = "com.himalife.srv.svc-device"
	rpcAnalysisServiceName  = "com.himalife.srv.svc-analysis"
)

// JinmuHealth 实现了 JinmuhealthHandler
type JinmuHealth struct {
	datastore              db.Datastore
	mailClient             *mail.Client
	algorithmClient        calcpb.CalcAPIService
	algorithmServerAddress string
	awsClient              *aws.Client
	analysisEngine         *ae.Engine
	wechat                 *wechat.Wxmp
	rpcSvc                 corepb.JinmuhealthAPIService
	jinmuidSvc             jinmuidpb.UserManagerAPIService
	subscriptionSvc        subscriptionpb.SubscriptionManagerAPIService
	deviceSvc              devicepb.DeviceManagerAPIService
	rpcAnalysisSvc         analysispb.AnalysisManagerAPIService
	blockerPool            *blocker.BlockerPool
}

// NewJinmuHealth 创建一个 JinmuHealth
func NewJinmuHealth(
	datastore db.Datastore,
	mailClient *mail.Client,
	algorithmClient calcpb.CalcAPIService,
	awsClient *aws.Client,
	ae *ae.Engine,
	wechat *wechat.Wxmp,
	algorithmServerAddress string) *JinmuHealth {

	j := &JinmuHealth{
		datastore:              datastore,
		mailClient:             mailClient,
		algorithmClient:        algorithmClient,
		awsClient:              awsClient,
		analysisEngine:         ae,
		wechat:                 wechat,
		rpcSvc:                 corepb.NewJinmuhealthAPIService(rpcServiceName, client.DefaultClient),
		jinmuidSvc:             jinmuidpb.NewUserManagerAPIService(jinmuidServiceName, client.DefaultClient),
		subscriptionSvc:        subscriptionpb.NewSubscriptionManagerAPIService(subscriptionServiceName, client.DefaultClient),
		deviceSvc:              devicepb.NewDeviceManagerAPIService(deviceServiceName, client.DefaultClient),
		rpcAnalysisSvc:         analysispb.NewAnalysisManagerAPIService(rpcAnalysisServiceName, client.DefaultClient),
		algorithmServerAddress: algorithmServerAddress,
	}
	return j
}
