package rest

import (
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	"github.com/kataras/iris/v12"
)

// JsSdkSignConfigRequest js sdk 配置请求
type JsSdkSignConfigRequest struct {
	URL string `json:"url"`
}

// JsSdkSignConfig js SDK sdk 配置
type JsSdkSignConfig struct {
	AppID     string `json:"app_id"`
	Timestamp string `json:"timestamp"`
	NonceStr  string `json:"noncestr"`
	Signature string `json:"signature"`
}

// GetWxJsSdkConfig 获取微信js SDK 配置
func (h *handler) GetWxJsSdkConfig(ctx iris.Context) {
	var b JsSdkSignConfigRequest
	err := ctx.ReadJSON(&b)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}

	req := new(proto.WechatGetWxJsSdkConfigRequest)
	req.Url = b.URL

	resp, err := h.rpcSvc.WechatGetWxJsSdkConfig(
		newRPCContext(ctx), req,
	)
	if err != nil {
		writeError(ctx, wrapError(ErrRPCInternal, "", err), false)
		return
	}

	cfg := JsSdkSignConfig{
		AppID:     resp.Config.AppId,
		Timestamp: resp.Config.Timestamp,
		NonceStr:  resp.Config.Noncestr,
		Signature: resp.Config.Signature,
	}

	rest.WriteOkJSON(ctx, cfg)
}
