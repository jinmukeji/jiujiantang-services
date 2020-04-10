package rest

import (
	"fmt"

	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	"github.com/kataras/iris/v12"
)

const (
	// SetSecureEmailSemNotification 设置安全邮箱
	SetSecureEmailSemNotification string = "set_secure_email"
	// UnsetSecureEmailSemNotification 解绑安全邮箱
	UnsetSecureEmailSemNotification string = "unset_secure_email"
	// ModifySecureEmailSemNotification 修改安全邮箱
	ModifySecureEmailSemNotification string = "modify_secure_email"
	// FindUsernameSemNotification 找回用户名
	FindUsernameSemNotification string = "find_username"
	// ResetPasswordSemNotification  找回/重置密码
	ResetPasswordSemNotification string = "reset_password"
)

// EmailNotificationRequest 邮箱通知请求
type EmailNotificationRequest struct {
	Email             string `json:"email"`                 // 邮箱
	Type              string `json:"type"`                  // 通知类型
	Language          string `json:"language"`              // 语言
	UserID            int32  `json:"user_id"`               // UserID
	SendToNewIfModify bool   `json:"send_to_new_if_modify"` // 是否是修改安全邮箱时向新邮箱发送邮件
}

// EmailNotificationResponse 邮件通知返回
type EmailNotificationResponse struct {
	SerialNumber string `json:"serial_number"` // 序列号
	Acknowledged bool   `json:"acknowledged"`  // 是否受理发送请求
	Message      string `json:"message"`       // 错误消息内容
}

// EmailNotification 邮件通知
func (h *webHandler) EmailNotification(ctx iris.Context) {
	var reqEmailNotification EmailNotificationRequest
	errReadJSON := ctx.ReadJSON(&reqEmailNotification)
	if errReadJSON != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", errReadJSON), false)
		return
	}

	hasLogged := mapFromActionToLogStatus(reqEmailNotification.Type)
	if hasLogged {
		req := new(proto.LoggedInEmailNotificationRequest)
		req.Email = reqEmailNotification.Email
		protoLanguage, errmapRestLanguageToProto := mapRestLanguageToProto(reqEmailNotification.Language)
		if errmapRestLanguageToProto != nil {
			writeError(ctx, wrapError(ErrInvalidValue, "", errmapRestLanguageToProto), false)
			return
		}
		req.Language = protoLanguage
		action, errMapEmailNotificationTypeToSemLoggedAction := mapEmailNotificationTypeToSemLoggedAction(reqEmailNotification.Type)
		if errMapEmailNotificationTypeToSemLoggedAction != nil {
			writeError(ctx, wrapError(ErrInvalidValue, "", errMapEmailNotificationTypeToSemLoggedAction), false)
			return
		}
		req.Action = action
		req.UserId = reqEmailNotification.UserID
		req.SendToNewIfModify = reqEmailNotification.SendToNewIfModify
		resp, errSemNotification := h.rpcSvc.LoggedInEmailNotification(newRPCContext(ctx), req)
		if errSemNotification != nil {
			writeRpcInternalError(ctx, errSemNotification, false)
			return
		}
		rest.WriteOkJSON(ctx, EmailNotificationResponse{
			SerialNumber: resp.SerialNumber,
			Acknowledged: resp.Acknowledged,
			Message:      resp.Message,
		})
	} else {
		req := new(proto.NotLoggedInEmailNotificationRequest)
		req.Email = reqEmailNotification.Email
		protoLanguage, errmapRestLanguageToProto := mapRestLanguageToProto(reqEmailNotification.Language)
		if errmapRestLanguageToProto != nil {
			writeError(ctx, wrapError(ErrInvalidValue, "", errmapRestLanguageToProto), false)
			return
		}
		req.Language = protoLanguage
		action, errMapEmailNotificationTypeToSemNotloggedAction := mapEmailNotificationTypeToSemNotloggedAction(reqEmailNotification.Type)
		if errMapEmailNotificationTypeToSemNotloggedAction != nil {
			writeError(ctx, wrapError(ErrInvalidValue, "", errMapEmailNotificationTypeToSemNotloggedAction), false)
			return
		}
		req.Action = action
		resp, errSemNotification := h.rpcSvc.NotLoggedInEmailNotification(newRPCContext(ctx), req)
		if errSemNotification != nil {
			writeRpcInternalError(ctx, errSemNotification, false)
			return
		}
		rest.WriteOkJSON(ctx, EmailNotificationResponse{
			SerialNumber: resp.SerialNumber,
			Acknowledged: resp.Acknowledged,
			Message:      resp.Message,
		})
	}
}

// mapEmailNotificationTypeToSemLoggedAction 将邮件通知的类型转换为 proto 里面已登陆状态下邮件通知的类型
func mapEmailNotificationTypeToSemLoggedAction(emailNotificationType string) (proto.LoggedInSemTemplateAction, error) {
	switch emailNotificationType {
	case SetSecureEmailSemNotification:
		return proto.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_SET_SECURE_EMAIL, nil
	case UnsetSecureEmailSemNotification:
		return proto.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_UNSET_SECURE_EMAIL, nil
	case ModifySecureEmailSemNotification:
		return proto.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_MODIFY_SECURE_EMAIL, nil
	}
	return proto.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_INVALID, fmt.Errorf("invalid string email notification type %s", emailNotificationType)
}

// mapEmailNotificationTypeToSemNotloggedAction 将邮件通知的类型转换为 proto 里面未登陆状态下邮件通知的类型
func mapEmailNotificationTypeToSemNotloggedAction(emailNotificationType string) (proto.NotLoggedInSemTemplateAction, error) {
	switch emailNotificationType {
	case FindUsernameSemNotification:
		return proto.NotLoggedInSemTemplateAction_NOT_LOGGED_IN_SEM_TEMPLATE_ACTION_FIND_USERNAME, nil
	case ResetPasswordSemNotification:
		return proto.NotLoggedInSemTemplateAction_NOT_LOGGED_IN_SEM_TEMPLATE_ACTION_FIND_RESET_PASSWORD, nil
	}
	return proto.NotLoggedInSemTemplateAction_NOT_LOGGED_IN_SEM_TEMPLATE_ACTION_INVALID, fmt.Errorf("invalid string email notification type %s", emailNotificationType)
}

// 根据发送行为判断用户是否登录
func mapFromActionToLogStatus(action string) bool {
	switch action {
	case SetSecureEmailSemNotification:
		return true
	case UnsetSecureEmailSemNotification:
		return true
	case ModifySecureEmailSemNotification:
		return true
	case FindUsernameSemNotification:
		return false
	case ResetPasswordSemNotification:
		return false
	default:
		return false
	}
}
