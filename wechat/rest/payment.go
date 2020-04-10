package rest

import (
	"fmt"
	"net/url"

	"github.com/jinmukeji/gf-api2/pkg/rest"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	"github.com/kataras/iris/v12"
)

// Payment 支付
type Payment struct {
	RecordID int32 `json:"record_id"`
}

// MakePayment 支付
func (h *handler) MakePayment(ctx iris.Context) {
	var payment Payment
	err := ctx.ReadJSON(&payment)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}

	// Session
	session, err := h.getSession(ctx)
	if err != nil || !session.Authorized || session.IsExpired {
		path := fmt.Sprintf("%s/app.html#/analysisreport?record_id=%d", h.WxH5ServerBase, payment.RecordID)
		redirectURL := fmt.Sprintf("%s/wx/oauth?redirect=%s", h.WxCallbackServerBase, url.QueryEscape(path))
		writeSessionErrorJSON(ctx, redirectURL, err)
		return
	}
	req := new(proto.JinmuLMakePaymentRequest)
	req.RecordId = payment.RecordID
	_, err = h.rpcSvc.JinmuLMakePayment(
		newRPCContext(ctx), req,
	)
	if err != nil {
		writeError(ctx, wrapError(ErrRPCInternal, "", err), false)
		return
	}

	rest.WriteOkJSON(ctx, nil)
}
