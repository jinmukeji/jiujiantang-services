package rest

import (
	"github.com/jinmukeji/gf-api2/pkg/rest"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	"github.com/kataras/iris/v12"
)

// SetSecureEmailRequest 设置安全邮箱请求
type SetSecureEmailRequest struct {
	Email            string `json:"email"`             // 邮箱
	SerialNumber     string `json:"serial_number"`     // 流水号
	VerificationCode string `json:"verification_code"` // 验证码
}

// 设置安全邮箱
func (h *webHandler) SetSecureEmail(ctx iris.Context) {
	userID, err := ctx.Params().GetInt("user_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	var reqSetSecureEmail SetSecureEmailRequest
	errReadJSON := ctx.ReadJSON(&reqSetSecureEmail)
	if errReadJSON != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", errReadJSON), false)
		return
	}
	req := new(proto.SetSecureEmailRequest)
	req.UserId = int32(userID)
	req.VerificationCode = reqSetSecureEmail.VerificationCode
	req.SerialNumber = reqSetSecureEmail.SerialNumber
	req.Email = reqSetSecureEmail.Email
	_, errSetModifyEmail := h.rpcSvc.SetSecureEmail(newRPCContext(ctx), req)
	if errSetModifyEmail != nil {
		writeRpcInternalError(ctx, errSetModifyEmail, false)
		return
	}
	rest.WriteOkJSON(ctx, nil)
}
