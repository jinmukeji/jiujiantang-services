package rest

import (
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	"github.com/kataras/iris/v12"
)

// UnsetSecureEmailRequest 解除设置安全邮箱的请求
type UnsetSecureEmailRequest struct {
	Email            string `json:"email"`             // 邮箱
	SerialNumber     string `json:"serial_number"`     // 流水号
	VerificationCode string `json:"verification_code"` // 验证码
}

// UnsetSecureEmail 解除设置安全邮箱
func (h *webHandler) UnsetSecureEmail(ctx iris.Context) {
	userID, err := ctx.Params().GetInt("user_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	var reqUnsetSecureEmail UnsetSecureEmailRequest
	errReadJSON := ctx.ReadJSON(&reqUnsetSecureEmail)
	if errReadJSON != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", errReadJSON), false)
		return
	}
	req := new(proto.UnsetSecureEmailRequest)
	req.UserId = int32(userID)
	req.VerificationCode = reqUnsetSecureEmail.VerificationCode
	req.SerialNumber = reqUnsetSecureEmail.SerialNumber
	req.Email = reqUnsetSecureEmail.Email
	_, errUnsetSecureEmail := h.rpcSvc.UnsetSecureEmail(newRPCContext(ctx), req)
	if errUnsetSecureEmail != nil {
		writeRpcInternalError(ctx, errUnsetSecureEmail, false)
		return
	}
	rest.WriteOkJSON(ctx, nil)
}
