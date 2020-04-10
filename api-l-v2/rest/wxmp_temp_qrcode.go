package rest

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	"github.com/jinmukeji/jiujiantang-services/service/wechat"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	"github.com/kataras/iris/v12"
)

// GetWxmpTempQrCodeUrl 得到临时微信二维码
func (h *v2Handler) GetWxmpTempQrCodeUrl(ctx iris.Context) {
	req := new(proto.GetWxmpTempQrCodeUrlRequest)
	resp, err := h.rpcSvc.GetWxmpTempQrCodeUrl(
		newRPCContext(ctx), req,
	)
	if err != nil {
		writeRpcInternalError(ctx, err, false)
		return
	}

	expiredAt, _ := ptypes.Timestamp(resp.ExpiredTime)
	rest.WriteOkJSON(ctx, wechat.WxmpTempQrCodeURL{
		ImageURL:  resp.ImageUrl,
		RawURL:    resp.RawUrl,
		ExpiredAt: expiredAt,
		SceneID:   resp.SceneId,
	})
}
