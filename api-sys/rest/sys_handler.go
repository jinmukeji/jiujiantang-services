package rest

import (
	"path"

	"github.com/jinmukeji/jiujiantang-services/api-sys/preference"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	"github.com/micro/go-micro/v2/client"
)

type sysHandler struct {
	rpcSvc            proto.XimaAPIService
	clientPreferences preference.ClientPreferences
}

const (
	rpcServiceName = "com.himalife.srv.svc-biz-core"
)

func newSysHandler(configFile string) *sysHandler {
	return &sysHandler{
		rpcSvc:            proto.NewXimaAPIService(rpcServiceName, client.DefaultClient),
		clientPreferences: preference.NewClientPreferences(path.Join(configFile)),
	}
}
