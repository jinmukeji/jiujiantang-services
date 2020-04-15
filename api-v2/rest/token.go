package rest

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	jinmuidpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	"github.com/kataras/iris/v12"
)

// LatestToken 最新token信息
type LatestToken struct {
	AccessToken string    `json:"access_token"`
	UserID      int32     `json:"user_id"`
	ExpiredAt   time.Time `json:"expired_at"`
}

// GetLatestToken 获取最新的token
func (h *v2Handler) GetLatestToken(ctx iris.Context) {
	userID, err := ctx.Params().GetInt("user_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	req := new(jinmuidpb.GetLatestTokenRequest)
	req.UserId = int32(userID)
	resp, err := h.rpcJinmuidSvc.GetLatestToken(
		newRPCContext(ctx), req,
	)
	if err != nil {
		writeRPCInternalError(ctx, err, false)
		return
	}
	expiredAt, _ := ptypes.Timestamp(resp.ExpiredTime)
	rest.WriteOkJSON(ctx, LatestToken{
		AccessToken: resp.AccessToken,
		UserID:      req.UserId,
		ExpiredAt:   expiredAt.UTC(),
	})
}
