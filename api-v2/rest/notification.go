package rest

import (
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	"github.com/kataras/iris/v12"
)

// ReadPushNotification 阅读通知的body
type ReadPushNotification struct {
	PnID int32 `json:"pn_id"`
}

// PushNotifications  通知
type PushNotifications struct {
	PnID          int32  `json:"pn_id"`
	PnDisplayTime string `json:"pn_display_time"`
	PnTitle       string `json:"pn_title"`
	PnImageURL    string `json:"pn_image_url"`
	PnContentURL  string `json:"pn_content_url"`
}

// ReadPushNotification 阅读通知
func (h v2Handler) ReadPushNotification(ctx iris.Context) {
	var readPushNotification ReadPushNotification
	err := ctx.ReadJSON(&readPushNotification)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}
	req := new(proto.ReadPushNotificationRequest)
	req.PnId = readPushNotification.PnID
	_, errReadPushNotification := h.rpcSvc.ReadPushNotification(
		newRPCContext(ctx), req,
	)
	if errReadPushNotification != nil {
		writeRPCInternalError(ctx, errReadPushNotification, false)
		return
	}
	rest.WriteOkJSON(ctx, nil)
}

// GetPushNotifications 得到通知
func (h v2Handler) GetPushNotifications(ctx iris.Context) {
	req := new(proto.GetPushNotificationsRequest)
	resp, errGetPushNotifications := h.rpcSvc.GetPushNotifications(
		newRPCContext(ctx), req,
	)
	if errGetPushNotifications != nil {
		writeRPCInternalError(ctx, errGetPushNotifications, true)
		return
	}
	pns := make([]PushNotifications, 0)
	for _, notification := range resp.Notifications {
		pns = append(pns, PushNotifications{
			PnID:          notification.PnId,
			PnContentURL:  notification.PnContentUrl,
			PnDisplayTime: notification.PnDisplayTime,
			PnImageURL:    notification.PnImageUrl,
			PnTitle:       notification.PnTitle,
		})
	}
	rest.WriteOkJSON(ctx, pns)
}
