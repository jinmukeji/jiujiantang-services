package handler

import (
	"github.com/jinmukeji/ae/v2/biz"
	"github.com/jinmukeji/jiujiantang-services/analysis/aws"
	db "github.com/jinmukeji/jiujiantang-services/analysis/mysqldb"
	subscriptionpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/subscription/v1"
	jinmuidpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	"github.com/micro/go-micro/v2/client"
)

// AnalysisManagerService 设备关联srv
type AnalysisManagerService struct {
	database        db.Datastore
	biz             *biz.BizEngineManager
	presetsFilePath string
	awsClient       *aws.Client
	jinmuidSvc      jinmuidpb.UserManagerAPIService
	subscriptionSvc subscriptionpb.SubscriptionManagerAPIService
}

const (
	jinmuidServiceName      = "com.himalife.srv.svc-jinmuid"
	subscriptionServiceName = "com.himalife.srv.svc-subscription"
)

// NewAnalysisManagerService 构建AnalysisManagerService
func NewAnalysisManagerService(datastore db.Datastore, biz *biz.BizEngineManager, presetsFilePath string, awsClient *aws.Client) *AnalysisManagerService {
	j := &AnalysisManagerService{
		database:        datastore,
		biz:             biz,
		presetsFilePath: presetsFilePath,
		awsClient:       awsClient,
		jinmuidSvc:      jinmuidpb.NewUserManagerAPIService(jinmuidServiceName, client.DefaultClient),
		subscriptionSvc: subscriptionpb.NewSubscriptionManagerAPIService(subscriptionServiceName, client.DefaultClient),
	}
	return j
}
