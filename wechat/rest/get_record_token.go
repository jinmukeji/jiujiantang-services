package rest

import (
	"fmt"
	"net/url"

	"github.com/jinmukeji/gf-api2/pkg/rest"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	"github.com/kataras/iris/v12"
)

// RecordToken 分享记录的token
type RecordToken struct {
	Token string `json:"token"`
}

// GetRecordToken 得到recordToken
func (h *handler) GetRecordToken(ctx iris.Context) {
	recordID, err := ctx.Params().GetInt("record_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}

	// Session
	session, err := h.getSession(ctx)
	if err != nil || !session.Authorized || session.IsExpired {
		path := fmt.Sprintf("%s/app.html#/analysisreport?record_id=%d", h.WxH5ServerBase, recordID)
		redirectURL := fmt.Sprintf("%s/wx/oauth?redirect=%s", h.WxCallbackServerBase, url.QueryEscape(path))
		writeSessionErrorJSON(ctx, redirectURL, err)
		return
	}

	req := new(proto.CreateReportShareTokenRequest)
	req.RecordId = int32(recordID)
	resp, err := h.rpcSvc.CreateReportShareToken(
		newRPCContext(ctx), req,
	)
	if err != nil {
		writeError(ctx, wrapError(ErrRPCInternal, "", err), false)
		return
	}
	rest.WriteOkJSON(ctx, RecordToken{
		Token: resp.Token,
	})
}
