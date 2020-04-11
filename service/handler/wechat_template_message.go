package handler

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
)

// WechatSendViewReportTemplateMessage 发送模版
func (j *JinmuHealth) WechatSendViewReportTemplateMessage(ctx context.Context, req *proto.WechatSendViewReportTemplateMessageRequest, resp *proto.WechatSendViewReportTemplateMessageResponse) error {
	w := j.wechat
	reportedAt, _ := ptypes.Timestamp(req.ReportedTime)
	return w.SendViewReportTemplateMessage(req.OpenId, int(req.RecordId), req.Nickname, reportedAt)
}
