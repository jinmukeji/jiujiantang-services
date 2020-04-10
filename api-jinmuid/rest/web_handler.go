package rest

import (
	jwtmiddleware "github.com/jinmukeji/gf-api2/pkg/rest/jwt"
	devicepb "github.com/jinmukeji/proto/gen/micro/idl/jm/device/v1"
	jinmuidpb "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	"github.com/micro/go-micro/client"
)

type webHandler struct {
	rpcSvc        jinmuidpb.UserManagerAPIService
	jwtMiddleware *jwtmiddleware.Middleware
	rpcDeviceSvc  devicepb.DeviceManagerAPIService
}

const (
	rpcServiceName       = "com.jinmuhealth.srv.svc-jinmuid"
	rpcDeviceServiceName = "com.jinmuhealth.srv.svc-device"
)

func newWebHandler(jwtMiddleware *jwtmiddleware.Middleware) *webHandler {
	return &webHandler{
		rpcSvc:        jinmuidpb.NewUserManagerAPIService(rpcServiceName, client.DefaultClient),
		jwtMiddleware: jwtMiddleware,
		rpcDeviceSvc:  devicepb.NewDeviceManagerAPIService(rpcDeviceServiceName, client.DefaultClient),
	}
}
