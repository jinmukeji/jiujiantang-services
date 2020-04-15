package handler

import (
	db "github.com/jinmukeji/jiujiantang-services/device/mysqldb"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/device/v1"
	"github.com/micro/go-micro/v2/client"
)

// DeviceManagerService 设备关联srv
type DeviceManagerService struct {
	database db.Datastore
	RPCSvc   proto.DeviceManagerAPIService
}

const (
	rpcServiceName = "com.himalife.srv.device"
)

// NewDeviceManagerService 构建DeviceManagerService
func NewDeviceManagerService(datastore db.Datastore) *DeviceManagerService {
	j := &DeviceManagerService{
		database: datastore,
		RPCSvc:   proto.NewDeviceManagerAPIService(rpcServiceName, client.DefaultClient),
	}
	return j
}
