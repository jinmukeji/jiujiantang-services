package handler

import (
	db "github.com/jinmukeji/jiujiantang-services/sem/mysqldb"
	sem "github.com/jinmukeji/jiujiantang-services/sem/sem_client"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/sem/v1"
	"github.com/micro/go-micro/client"
)

// SEMGateway 邮件网关
type SEMGateway struct {
	datastore        db.Datastore
	rpcSvc           proto.SemAPIService
	aliyunSEMClient  *sem.AliyunSEMClient
	neteaseSEMClient *sem.NetEaseSEMClient
}

const (
	// rpcServiceName RPC服务名称
	rpcServiceName = "com.himalife.srv.svc-sem-gw"
)

// NewSEMGateway 构建SEMGateway
func NewSEMGateway(datastore db.Datastore, aliyunSEMClient *sem.AliyunSEMClient, neteaseSEMClient *sem.NetEaseSEMClient) *SEMGateway {
	j := &SEMGateway{
		datastore:        datastore,
		rpcSvc:           proto.NewSemAPIService(rpcServiceName, client.DefaultClient),
		aliyunSEMClient:  aliyunSEMClient,
		neteaseSEMClient: neteaseSEMClient,
	}
	return j
}
