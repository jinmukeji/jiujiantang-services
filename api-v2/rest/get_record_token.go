package rest

import (
	"fmt"

	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	"github.com/kataras/iris/v12"
)

// RecordToken 分享记录的token
type RecordToken struct {
	Token string `json:"token"`
	Link  string `json:"link"`
}

// GetRecordToken 得到recordToken
func (h *v2Handler) GetRecordToken(ctx iris.Context) {
	if ctx.Values().GetString(ClientIDKey) == seamlessClient {
		writeError(
			ctx,
			wrapError(ErrDeniedToAccessAPI, "", fmt.Errorf("%s is denied to access this API", seamlessClient)),
			false,
		)
		return
	}
	recordID, err := ctx.Params().GetInt("record_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	req := new(proto.CreateReportShareTokenRequest)
	req.RecordId = int32(recordID)
	resp, err := h.rpcSvc.CreateReportShareToken(
		newRPCContext(ctx), req,
	)
	if err != nil {
		writeRPCInternalError(ctx, err, false)
		return
	}
	rest.WriteOkJSON(ctx, RecordToken{
		Token: resp.Token,
		Link:  resp.Link,
	})
}
