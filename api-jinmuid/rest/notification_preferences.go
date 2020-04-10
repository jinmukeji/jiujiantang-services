package rest

import (
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	"github.com/kataras/iris/v12"
)

// ModifyNotificationPreferences 通知配置首选项
type ModifyNotificationPreferences struct {
	PhoneEnabled  bool `json:"phone_enabled"`
	WechatEnabled bool `json:"wechat_enabled"`
	WeiboEnabled  bool `json:"weibo_enabled"`
}

// ModifyNotificationPreferencesBody 修改通知配置首选项的请求
type ModifyNotificationPreferencesBody struct {
	ModifyNotificationPreferences ModifyNotificationPreferences `json:"notification_preferences"`
}

// ModifyNotificationPreferences 修改通知配置首选项
func (h *webHandler) ModifyNotificationPreferences(ctx iris.Context) {
	var bodyModifyNotificationPreferences ModifyNotificationPreferencesBody

	err := ctx.ReadJSON(&bodyModifyNotificationPreferences)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}
	userID, err := ctx.Params().GetInt("user_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	reqModifyNotificationPreferences := new(proto.ModifyNotificationPreferencesRequest)
	reqModifyNotificationPreferences.UserId = int32(userID)
	reqModifyNotificationPreferences.PhoneEnabled = bodyModifyNotificationPreferences.ModifyNotificationPreferences.PhoneEnabled
	reqModifyNotificationPreferences.WechatEnabled = bodyModifyNotificationPreferences.ModifyNotificationPreferences.WechatEnabled
	reqModifyNotificationPreferences.WeiboEnabled = bodyModifyNotificationPreferences.ModifyNotificationPreferences.WeiboEnabled
	_, errModifyNotificationPreferences := h.rpcSvc.ModifyNotificationPreferences(
		newRPCContext(ctx), reqModifyNotificationPreferences,
	)
	if errModifyNotificationPreferences != nil {
		writeRpcInternalError(ctx, errModifyNotificationPreferences, false)
		return
	}
	rest.WriteOkJSON(ctx, nil)
}

// GetNotificationPreferences 获取通知配置首选项
func (h *webHandler) GetNotificationPreferences(ctx iris.Context) {
	userID, err := ctx.Params().GetInt("user_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	reqGetNotificationPreferences := new(proto.GetNotificationPreferencesRequest)
	reqGetNotificationPreferences.UserId = int32(userID)

	respGetNotificationPreferences, errGetNotificationPreferences := h.rpcSvc.GetNotificationPreferences(
		newRPCContext(ctx), reqGetNotificationPreferences,
	)
	if errGetNotificationPreferences != nil {
		writeRpcInternalError(ctx, errGetNotificationPreferences, false)
		return
	}

	rest.WriteOkJSON(ctx, ModifyNotificationPreferencesBody{
		ModifyNotificationPreferences{
			PhoneEnabled:  respGetNotificationPreferences.PhoneEnabled,
			WechatEnabled: respGetNotificationPreferences.WechatEnabled,
			WeiboEnabled:  respGetNotificationPreferences.WeiboEnabled,
		},
	})
}
