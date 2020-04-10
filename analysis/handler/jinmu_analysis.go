package handler

import (
	"github.com/jinmukeji/ae/v2/biz"
	"github.com/jinmukeji/gf-api2/analysis/aws"
	db "github.com/jinmukeji/gf-api2/analysis/mysqldb"
	jinmuidpb "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	subscriptionpb "github.com/jinmukeji/proto/gen/micro/idl/jm/subscription/v1"
	"github.com/micro/go-micro/client"
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
	jinmuidServiceName      = "com.jinmuhealth.srv.svc-jinmuid"
	subscriptionServiceName = "com.jinmuhealth.srv.svc-subscription"
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
