package rest

import (
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	"github.com/kataras/iris/v12"
)

// FindUsernameBySecureEmailBody 根据邮箱找回用户名请求
type FindUsernameBySecureEmailBody struct {
	Email            string `json:"email"`             // 邮箱
	SerialNumber     string `json:"serial_number"`     // 流水号
	VerificationCode string `json:"verification_code"` // 验证码
}

// Username 用户名
type Username struct {
	Username string `json:"username"` // 用户名
}

// FindUsernameBySecureEmail 根据邮箱找回用户名
func (h *webHandler) FindUsernameBySecureEmail(ctx iris.Context) {
	var reqFindUsernameBySecureEmail FindUsernameBySecureEmailBody
	errReadJSON := ctx.ReadJSON(&reqFindUsernameBySecureEmail)
	if errReadJSON != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", errReadJSON), false)
		return
	}
	req := new(proto.FindUsernameBySecureEmailRequest)
	req.VerificationCode = reqFindUsernameBySecureEmail.VerificationCode
	req.SerialNumber = reqFindUsernameBySecureEmail.SerialNumber
	req.Email = reqFindUsernameBySecureEmail.Email
	repl, errSetModifyEmail := h.rpcSvc.FindUsernameBySecureEmail(newRPCContext(ctx), req)
	if errSetModifyEmail != nil {
		writeRpcInternalError(ctx, errSetModifyEmail, false)
		return
	}
	rest.WriteOkJSON(ctx, Username{
		Username: repl.Username,
	})
}
