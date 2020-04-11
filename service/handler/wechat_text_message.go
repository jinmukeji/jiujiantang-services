package handler

import (
	"context"

	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
)

// WechatSendTextMessage 发送文字消息
func (j *JinmuHealth) WechatSendTextMessage(ctx context.Context, req *proto.WechatSendTextMessageRequest, resp *proto.WechatSendTextMessageResponse) error {
	w := j.wechat
	return w.WechatSendTextMessage(req.OpenId, req.Content)
}
