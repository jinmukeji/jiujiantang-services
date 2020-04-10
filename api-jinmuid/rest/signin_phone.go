package rest

import (
	"fmt"

	"github.com/jinmukeji/gf-api2/pkg/rest"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	"github.com/kataras/iris/v12"
)

// SetSigninPhoneBody 设置登录电话的body
type SetSigninPhoneBody struct {
	Phone               string              `json:"phone"`             // 手机号
	NationCode          string              `json:"nation_code"`       // 区号
	Mvc                 string              `json:"mvc"`               // 验证码
	SerialNumber        string              `json:"serial_number"`     // 序列号
	SmsNotificationType SmsNotificationType `json:"verification_type"` // 短信类型
}

// VerifySigninPhoneBody 验证登录手机号的body
type VerifySigninPhoneBody struct {
	SetSigninPhoneBody
}

// VerifyUserSigninPhoneResp 验证登录电话的返回
type VerifyUserSigninPhoneResp struct {
	VerificationNumber string `json:"verification_number"` // 验证号
	UserID             int32  `json:"user_id"`             // 用户ID
}

// ModifySigninPhoneBody 修改登录手机号
type ModifySigninPhoneBody struct {
	SetSigninPhoneBody
	VerificationNumber string `json:"verification_number"` // 验证号
	OldNationCode      string `json:"old_nation_code"`     // 旧手机区号
	OldPhone           string `json:"old_phone"`           // 旧手机号
}

// SetSigninPhone 设置登录手机
func (h *webHandler) SetSigninPhone(ctx iris.Context) {
	var body SetSigninPhoneBody
	errReadJSON := ctx.ReadJSON(&body)
	if errReadJSON != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", errReadJSON), false)
		return
	}
	if !checkNationCode(body.NationCode) {
		writeError(ctx, wrapError(ErrNationCode, "", fmt.Errorf("nation code %s is wrong", body.NationCode)), false)
		return
	}
	userID, err := ctx.Params().GetInt("user_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	req := new(proto.UserSetSigninPhoneRequest)
	req.Phone = body.Phone
	req.NationCode = body.NationCode
	req.Mvc = body.Mvc
	req.SerialNumber = body.SerialNumber
	req.UserId = int32(userID)
	_, errUserSetSigninPhone := h.rpcSvc.UserSetSigninPhone(newRPCContext(ctx), req)
	if errUserSetSigninPhone != nil {
		writeRpcInternalError(ctx, errUserSetSigninPhone, false)
		return
	}
	rest.WriteOkJSON(ctx, nil)
}

// VerifySigninPhone 验证登录手机号
func (h *webHandler) VerifySigninPhone(ctx iris.Context) {
	var body VerifySigninPhoneBody
	errReadJSON := ctx.ReadJSON(&body)
	if errReadJSON != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", errReadJSON), false)
		return
	}
	// 检查区号
	if !checkNationCode(body.NationCode) {
		writeError(ctx, wrapError(ErrNationCode, "", fmt.Errorf("nation code %s is wrong", body.NationCode)), false)
		return
	}
	// 短信类型不是重置密码和修改手机号
	if body.SmsNotificationType != ResetPasswordSmsNotification && body.SmsNotificationType != ModifyPhoneSmsNotification {
		writeError(ctx, wrapError(ErrWrongSmsNotificationType, "", fmt.Errorf("wrong sms notification type %s when verify signin phone", body.SmsNotificationType)), false)
		return
	}
	req := new(proto.VerifyUserSigninPhoneRequest)
	req.Phone = body.Phone
	req.NationCode = body.NationCode
	req.Mvc = body.Mvc
	req.SerialNumber = body.SerialNumber

	action, errmapSmsNotificationTypeToProtoTemplateAction := mapSmsNotificationTypeToProtoTemplateAction(body.SmsNotificationType)
	if errmapSmsNotificationTypeToProtoTemplateAction != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errmapSmsNotificationTypeToProtoTemplateAction), false)
		return
	}
	req.Action = action
	resp, errVerifyUserSigninPhone := h.rpcSvc.VerifyUserSigninPhone(newRPCContext(ctx), req)
	if errVerifyUserSigninPhone != nil {
		writeRpcInternalError(ctx, errVerifyUserSigninPhone, false)
		return
	}
	rest.WriteOkJSON(ctx, VerifyUserSigninPhoneResp{
		VerificationNumber: resp.VerificationNumber,
		UserID:             resp.UserId,
	})
}

// ModifySigninPhone 修改登录手机号
func (h *webHandler) ModifySigninPhone(ctx iris.Context) {
	var body ModifySigninPhoneBody
	errReadJSON := ctx.ReadJSON(&body)
	if errReadJSON != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", errReadJSON), false)
		return
	}
	userID, err := ctx.Params().GetInt("user_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	req := new(proto.UserModifyPhoneRequest)
	req.UserId = int32(userID)
	req.SerialNumber = body.SerialNumber
	req.VerificationNumber = body.VerificationNumber
	req.Phone = body.Phone
	req.OldNationCode = body.OldNationCode
	req.NationCode = body.NationCode
	req.Mvc = body.Mvc
	req.OldPhone = body.OldPhone
	_, errUserModifyPhone := h.rpcSvc.UserModifyPhone(newRPCContext(ctx), req)
	if errUserModifyPhone != nil {
		writeRpcInternalError(ctx, errUserModifyPhone, false)
		return
	}
	rest.WriteOkJSON(ctx, nil)
}
