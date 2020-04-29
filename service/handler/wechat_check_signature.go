package handler

import (
	"context"

	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
)

// WechatCheckWxSignature 验证微信接入的 Signature
func (j *JinmuHealth) WechatCheckWxSignature(ctx context.Context, req *proto.WechatCheckWxSignatureRequest, resp *proto.WechatCheckWxSignatureResponse) error {
	w := j.wechat
	resp.Ok = w.CheckWxSignature(req.Signature, req.Timestamp, req.Nonce)
	return nil
}
