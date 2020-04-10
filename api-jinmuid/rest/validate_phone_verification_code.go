package rest

import (
	"github.com/jinmukeji/gf-api2/pkg/rest"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	"github.com/kataras/iris/v12"
)

// ValidatePhoneVerificationCodeReply 验证短信验证码是否正确响应
type ValidatePhoneVerificationCodeReply struct {
	VerificationNumber string `json:"verification_number"` // 验证号
}

// ValidatePhoneVerificationCodeBody 验证手机验证码是否正确请求
type ValidatePhoneVerificationCodeBody struct {
	Phone        string `json:"phone"`         // 电话
	Mvc          string `json:"mvc"`           // 验证码
	SerialNumber string `json:"serial_number"` // 序列号
	NationCode   string `json:"nation_code"`   // 区号
}

// ValidatePhoneVerificationCode 验证手机验证码是否正确
func (h *webHandler) ValidatePhoneVerificationCode(ctx iris.Context) {
	var reqvalidatePhoneVerificationCodeBody ValidatePhoneVerificationCodeBody
	errReadJSON := ctx.ReadJSON(&reqvalidatePhoneVerificationCodeBody)
	if errReadJSON != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", errReadJSON), false)
		return
	}
	req := new(proto.ValidatePhoneVerificationCodeRequest)
	req.Phone = reqvalidatePhoneVerificationCodeBody.Phone
	req.Mvc = reqvalidatePhoneVerificationCodeBody.Mvc
	req.SerialNumber = reqvalidatePhoneVerificationCodeBody.SerialNumber
	req.NationCode = reqvalidatePhoneVerificationCodeBody.NationCode

	resp, errValidatePhoneVerificationCode := h.rpcSvc.ValidatePhoneVerificationCode(newRPCContext(ctx), req)
	if errValidatePhoneVerificationCode != nil {
		writeRpcInternalError(ctx, errValidatePhoneVerificationCode, false)
		return
	}
	rest.WriteOkJSON(ctx, ValidatePhoneVerificationCodeReply{
		VerificationNumber: resp.VerificationNumber,
	})
}
