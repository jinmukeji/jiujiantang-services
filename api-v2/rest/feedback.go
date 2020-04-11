package rest

import (
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	"github.com/kataras/iris/v12"
)

// Feedback 意见反馈
type Feedback struct {
	ContactWay string `json:"contact_way"`
	Content    string `json:"content"`
}

// SubmitFeedback 提交意见反馈
func (h *v2Handler) SubmitFeedback(ctx iris.Context) {
	var feedback Feedback
	err := ctx.ReadJSON(&feedback)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}
	req := new(proto.SubmitFeedbackRequest)
	req.ContactWay = feedback.ContactWay
	req.Content = feedback.Content
	_, err = h.rpcSvc.SubmitFeedback(
		newRPCContext(ctx), req,
	)
	if err != nil {
		writeRPCInternalError(ctx, err, false)
		return
	}

	rest.WriteOkJSON(ctx, nil)
}
