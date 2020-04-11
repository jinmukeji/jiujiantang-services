package rest

import (
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	"github.com/kataras/iris/v12"
)

// ValidateEmailVerificationCodeBody 验证邮箱验证码是否正确请求
type ValidateEmailVerificationCodeBody struct {
	Email            string `json:"email"`             // 邮箱
	SerialNumber     string `json:"serial_number"`     // 流水号
	VerificationCode string `json:"verification_code"` // 验证码
	VerificationType string `json:"verification_type"` // 验证类型
}

// ValidateEmailVerificationCodeReply 验证邮箱验证码是否正确响应
type ValidateEmailVerificationCodeReply struct {
	VerificationNumber string `json:"verification_number"` // 验证号
	UserID             int32  `json:"user_id"`             // 用户ID
}

// ValidateEmailVerificationCode 验证邮箱验证码是否正确
func (h *webHandler) ValidateEmailVerificationCode(ctx iris.Context) {
	var reqValidateEmailVerificationCode ValidateEmailVerificationCodeBody
	errReadJSON := ctx.ReadJSON(&reqValidateEmailVerificationCode)
	if errReadJSON != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", errReadJSON), false)
		return
	}
	req := new(proto.ValidateEmailVerificationCodeRequest)
	req.VerificationCode = reqValidateEmailVerificationCode.VerificationCode
	req.SerialNumber = reqValidateEmailVerificationCode.SerialNumber
	req.Email = reqValidateEmailVerificationCode.Email
	req.VerificationType = reqValidateEmailVerificationCode.VerificationType
	repl, errValidateEmailVerificationCode := h.rpcSvc.ValidateEmailVerificationCode(newRPCContext(ctx), req)
	if errValidateEmailVerificationCode != nil {
		writeRpcInternalError(ctx, errValidateEmailVerificationCode, false)
		return
	}
	rest.WriteOkJSON(ctx, ValidateEmailVerificationCodeReply{
		VerificationNumber: repl.VerificationNumber,
		UserID:             repl.UserId,
	})
}
