package handler

import (
	"context"

	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
)

// WechatGetWxJsSdkConfig 得到微信 JS SDK 配置信息
func (j *JinmuHealth) WechatGetWxJsSdkConfig(ctx context.Context, req *proto.WechatGetWxJsSdkConfigRequest, resp *proto.WechatGetWxJsSdkConfigResponse) error {
	w := j.wechat
	cfg, err := w.GetWxJsSdkConfig(req.Url)
	if err != nil {
		return NewError(ErrGetWxJsSdkConfigFaliure, err)
	}

	resp.Config = new(proto.JsSdkSignConfig)
	resp.Config.AppId = cfg.AppID
	resp.Config.Timestamp = cfg.Timestamp
	resp.Config.Noncestr = cfg.NonceStr
	resp.Config.Signature = cfg.Signature

	return nil
}
