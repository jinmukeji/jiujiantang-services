package rest

import (
	"fmt"
	"net/url"

	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	"github.com/kataras/iris/v12"
)

// SubmitRemarkReq 提交分析报告备注的 Request
type SubmitRemarkReq struct {
	Remark string `json:"remark"`
}

// SubmitRemark 提交分析报告备注
func (h *handler) SubmitRemark(ctx iris.Context) {
	// url params
	recordID, err := ctx.Params().GetInt("record_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}

	// Session
	session, err := h.getSession(ctx)
	// TODO: session 授权判断应该封装的 wrapper 之中处理
	if err != nil || !session.Authorized || session.IsExpired {
		path := fmt.Sprintf("%s/app.html#/analysisreport?record_id=%d", h.WxH5ServerBase, recordID)
		redirectURL := fmt.Sprintf("%s/wx/oauth?redirect=%s", h.WxCallbackServerBase, url.QueryEscape(path))
		writeSessionErrorJSON(ctx, redirectURL, err)
		return
	}

	// build rpc request
	req := new(proto.SubmitRemarkRequest)
	req.RecordId = int32(recordID)
	var remark SubmitRemarkReq
	err = ctx.ReadJSON(&remark)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}
	req.Remark = remark.Remark
	_, errResp := h.rpcSvc.SubmitRemark(newRPCContext(ctx), req)
	if errResp != nil {
		writeError(ctx, wrapError(ErrRPCInternal, "", errResp), false)
		return
	}

	rest.WriteOkJSON(ctx, nil)
}
