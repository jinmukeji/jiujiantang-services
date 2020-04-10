package rest

import (
	"fmt"

	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	"github.com/kataras/iris/v12"
)

// SubmitRemarkReq 提交分析报告备注的 Request
type SubmitRemarkReq struct {
	UserID int    `json:"user_id"`
	Remark string `json:"remark"`
}

// SubmitRemark 提交分析报告备注
func (h *v2Handler) SubmitRemark(ctx iris.Context) {
	// url params
	recordID, _ := ctx.Params().GetInt("record_id")
	// build rpc request
	req := new(proto.SubmitRemarkRequest)
	req.RecordId = int32(recordID)
	var remark SubmitRemarkReq
	err := ctx.ReadJSON(&remark)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}
	req.Remark = remark.Remark
	userID := int32(remark.UserID)
	if userID == 0 {
		writeError(
			ctx,
			wrapError(ErrValueRequired, fmt.Sprintf("Invalid request data, userID is %d", userID), nil),
			false,
		)
		return
	}
	req.UserId = userID
	_, errResp := h.rpcSvc.SubmitRemark(newRPCContext(ctx), req)
	if errResp != nil {
		writeRPCInternalError(ctx, errResp, false)
		return
	}

	rest.WriteOkJSON(ctx, nil)
}
