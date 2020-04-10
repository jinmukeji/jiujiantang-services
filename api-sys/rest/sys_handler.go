package rest

import (
	"path"

	"github.com/jinmukeji/jiujiantang-services/api-sys/preference"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	"github.com/micro/go-micro/client"
)

type sysHandler struct {
	rpcSvc            proto.JinmuhealthAPIService
	clientPreferences preference.ClientPreferences
}

const (
	rpcServiceName = "com.jinmuhealth.srv.svc-biz-core"
)

func newSysHandler(configFile string) *sysHandler {
	return &sysHandler{
		rpcSvc:            proto.NewJinmuhealthAPIService(rpcServiceName, client.DefaultClient),
		clientPreferences: preference.NewClientPreferences(path.Join(configFile)),
	}
}
