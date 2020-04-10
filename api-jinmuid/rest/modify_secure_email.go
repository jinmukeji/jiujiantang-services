package rest

import (
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	"github.com/kataras/iris/v12"
)

// ModifySecureEmailRequest 修改安全邮箱请求
type ModifySecureEmailRequest struct {
	NewEmail              string `json:"new_email"`               // 新的邮箱
	NewVerificationCode   string `json:"new_verification_code"`   // 验证码
	NewSerialNumber       string `json:"new_serial_number"`       // 序列号
	OldEmail              string `json:"old_email"`               // 旧的邮箱
	OldVerificationNumber string `json:"old_verification_number"` // 旧的邮箱验证号
}

// 修改安全邮箱
func (h *webHandler) ModifySecureEmail(ctx iris.Context) {

	// 用户登录后修改安全邮箱先要向旧邮箱发送验证码，此时调用发送邮件接口 Post("/notification/email/logged"),返回一个serial_number
	// 然后验证用户填写的验证码是否正确，此时调用 Post("/user/validate_email_verification_code")
	// 验证成功后返回一个verification_number
	// 然后向用户填写的新邮箱发送验证码，此时调用发送邮件接口 Post("/notification/email/logged"),返回一个serial_number
	// 然后把新邮箱NewEmail,新邮箱的验证码NewVerificationCode，新邮箱的serial_number为NewSerialNumber，
	// 旧邮箱OldEmail，旧邮箱的verification_number为OldVerificationNumber，此时调用当前接口Post("/user/{user_id}/modify_secure_email")

	userID, err := ctx.Params().GetInt("user_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	var reqModifySecureEmail ModifySecureEmailRequest
	errReadJSON := ctx.ReadJSON(&reqModifySecureEmail)
	if errReadJSON != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", errReadJSON), false)
		return
	}
	req := new(proto.ModifySecureEmailRequest)
	req.UserId = int32(userID)
	req.NewEmail = reqModifySecureEmail.NewEmail
	req.NewVerificationCode = reqModifySecureEmail.NewVerificationCode
	req.NewSerialNumber = reqModifySecureEmail.NewSerialNumber
	req.OldEmail = reqModifySecureEmail.OldEmail
	req.OldVerificationNumber = reqModifySecureEmail.OldVerificationNumber
	_, errSetModifyEmail := h.rpcSvc.ModifySecureEmail(newRPCContext(ctx), req)
	if errSetModifyEmail != nil {
		writeRpcInternalError(ctx, errSetModifyEmail, false)
		return
	}
	rest.WriteOkJSON(ctx, nil)
}
