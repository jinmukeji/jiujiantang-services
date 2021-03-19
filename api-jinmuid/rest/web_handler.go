package rest

import (
	jwtmiddleware "github.com/jinmukeji/jiujiantang-services/pkg/rest/jwt"
	devicepb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/device/v1"
	jinmuidpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	"github.com/micro/go-micro/v2/client"
)

type webHandler struct {
	rpcSvc        jinmuidpb.UserManagerAPIService
	jwtMiddleware *jwtmiddleware.Middleware
	rpcDeviceSvc  devicepb.DeviceManagerAPIService
}

const (
	rpcServiceName       = "com.himalife.srv.svc-jinmuid"
	rpcDeviceServiceName = "com.himalife.srv.svc-device"
)

func newWebHandler(jwtMiddleware *jwtmiddleware.Middleware) *webHandler {
	return &webHandler{
		rpcSvc:        jinmuidpb.NewUserManagerAPIService(rpcServiceName, client.DefaultClient),
		jwtMiddleware: jwtMiddleware,
		rpcDeviceSvc:  devicepb.NewDeviceManagerAPIService(rpcDeviceServiceName, client.DefaultClient),
	}
}
