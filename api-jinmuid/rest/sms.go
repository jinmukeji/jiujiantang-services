package rest

import (
	"fmt"

	"github.com/jinmukeji/gf-api2/pkg/rest"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	"github.com/kataras/iris/v12"
)

var nationCodes = [18]string{"+86", "+852", "+853", "+886", "+1", "+81", "+44", "+673", "+855", "+62", "+856", "+60", "+95", "+63", "+65", "+66", "+84", "+82"}

// SmsNotificationType 短信通知类型
type SmsNotificationType string

const (
	SigninSmsNotification        SmsNotificationType = "sign_in"
	SignupSmsNotification        SmsNotificationType = "sign_up"
	ResetPasswordSmsNotification SmsNotificationType = "reset_password"
	ModifyPhoneSmsNotification   SmsNotificationType = "modify_phone"
	SetPhoneSmsNotification      SmsNotificationType = "set_phone"
)

// SmsNotificationBody 短信通知body
type SmsNotificationBody struct {
	Phone             string              `json:"phone"`
	NotificationType  SmsNotificationType `json:"sms_notification_type"`
	Language          string              `json:"language"`
	UserID            int32               `json:"user_id"`
	NationCode        string              `json:"nation_code"`
	SendToNewIfModify bool                `json:"send_to_new_if_modify"`
}

// SmsNotification 短信通知
type SmsNotification struct {
	SerialNumber string `json:"serial_number"`
	Acknowledged bool   `json:"acknowledged"`
	Message      string `json:"message"`
}

// SmsNotification 短信通知
func (h *webHandler) SmsNotification(ctx iris.Context) {
	var body SmsNotificationBody
	errReadJSON := ctx.ReadJSON(&body)
	if errReadJSON != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", errReadJSON), false)
		return
	}
	if !checkNationCode(body.NationCode) {
		writeError(ctx, wrapError(ErrNationCode, "", fmt.Errorf("nation code %s is wrong", body.NationCode)), false)
		return
	}
	req := new(proto.SmsNotificationRequest)
	protoLanguage, errmapRestLanguageToProto := mapRestLanguageToProto(body.Language)
	if errmapRestLanguageToProto != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errmapRestLanguageToProto), false)
		return
	}
	req.Language = protoLanguage
	req.NationCode = body.NationCode
	req.Phone = body.Phone
	action, errmapSmsNotificationTypeToProtoTemplateAction := mapSmsNotificationTypeToProtoTemplateAction(body.NotificationType)
	if errmapSmsNotificationTypeToProtoTemplateAction != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errmapSmsNotificationTypeToProtoTemplateAction), false)
		return
	}
	req.Action = action
	req.IsForced = false
	req.SendToNewIfModify = body.SendToNewIfModify
	req.UserId = body.UserID
	resp, errSmsNotification := h.rpcSvc.SmsNotification(newRPCContext(ctx), req)
	if errSmsNotification != nil {
		writeRpcInternalError(ctx, errSmsNotification, false)
		return
	}
	rest.WriteOkJSON(ctx, SmsNotification{
		SerialNumber: resp.SerialNumber,
		Message:      resp.Message,
		Acknowledged: resp.Acknowledged,
	})
}

// checkNationCode 检查 nation_code
func checkNationCode(nationCode string) bool {
	for _, code := range nationCodes {
		if code == nationCode {
			return true
		}
	}
	return false
}

func mapSmsNotificationTypeToProtoTemplateAction(smsNotificationType SmsNotificationType) (proto.TemplateAction, error) {
	switch smsNotificationType {
	case SigninSmsNotification:
		return proto.TemplateAction_TEMPLATE_ACTION_SIGN_IN, nil
	case SignupSmsNotification:
		return proto.TemplateAction_TEMPLATE_ACTION_SIGN_UP, nil
	case ResetPasswordSmsNotification:
		return proto.TemplateAction_TEMPLATE_ACTION_RESET_PASSWORD, nil
	case ModifyPhoneSmsNotification:
		return proto.TemplateAction_TEMPLATE_ACTION_MODIFY_PHONE_NUMBER, nil
	case SetPhoneSmsNotification:
		return proto.TemplateAction_TEMPLATE_ACTION_SET_PHONE_NUMBER, nil
	}
	return proto.TemplateAction_TEMPLATE_ACTION_INVALID, fmt.Errorf("invalid sms notification type %s", smsNotificationType)
}
