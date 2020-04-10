package handler

import (
	ae "github.com/jinmukeji/ae-v1/core"
	"github.com/jinmukeji/jiujiantang-services/analysis/aws"
	"github.com/jinmukeji/jiujiantang-services/pkg/blocker"
	"github.com/jinmukeji/jiujiantang-services/service/mail"
	db "github.com/jinmukeji/jiujiantang-services/service/mysqldb"
	"github.com/jinmukeji/jiujiantang-services/service/wechat"
	jinmuidpb "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	analysispb "github.com/jinmukeji/proto/gen/micro/idl/jm/analysis/v1"
	corepb "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	devicepb "github.com/jinmukeji/proto/gen/micro/idl/jm/device/v1"
	subscriptionpb "github.com/jinmukeji/proto/gen/micro/idl/jm/subscription/v1"
	calcpb "github.com/jinmukeji/proto/gen/micro/idl/platform/calc/v2"
	"github.com/micro/go-micro/client"
)

const (
	rpcServiceName          = "com.jinmuhealth.srv.svc-biz-core"
	jinmuidServiceName      = "com.jinmuhealth.srv.svc-jinmuid"
	subscriptionServiceName = "com.jinmuhealth.srv.svc-subscription"
	deviceServiceName       = "com.jinmuhealth.srv.svc-device"
	rpcAnalysisServiceName  = "com.jinmuhealth.srv.svc-analysis"
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
	blockerPool *blocker.BlockerPool,
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
		blockerPool:            blockerPool,
		algorithmServerAddress: algorithmServerAddress,
	}
	return j
}
