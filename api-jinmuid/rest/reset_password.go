package rest

import (
	"fmt"

	"github.com/jinmukeji/gf-api2/pkg/rest"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	"github.com/kataras/iris/v12"
)

const (
	// PhoneVerificationType 手机号验证码的类型
	PhoneVerificationType = "phone"
	// EmailVerificationType 邮件验证码的类型
	EmailVerificationType = "email"
)

// ResetPasswordBody 重置密码的body
type ResetPasswordBody struct {
	PlainPassword      string `json:"plain_password"`      // 明文密码
	VerificationNumber string `json:"verification_number"` // 验证号
	VerificationType   string `json:"verification_type"`   // 验证类型
}

// ResetPassword 重置密码
func (h *webHandler) ResetPassword(ctx iris.Context) {
	var body ResetPasswordBody
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
	req := new(proto.UserResetPasswordRequest)
	req.UserId = int32(userID)
	req.VerificationNumber = body.VerificationNumber
	req.PlainPassword = body.PlainPassword
	protoVerificationType, errMapProtoVerificationType := mapRestVerificationTypeToProto(body.VerificationType)
	if errMapProtoVerificationType != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errMapProtoVerificationType), false)
		return
	}
	req.VerificationType = protoVerificationType
	_, errUserResetPassword := h.rpcSvc.UserResetPassword(newRPCContext(ctx), req)
	if errUserResetPassword != nil {
		writeRpcInternalError(ctx, errUserResetPassword, false)
		return
	}
	rest.WriteOkJSON(ctx, nil)
}

// mapRestVerificationTypeToProto 将 rest 的 string 类型验证类型转换为 proto 类型
func mapRestVerificationTypeToProto(verificationType string) (proto.VerificationType, error) {
	switch verificationType {
	case PhoneVerificationType:
		return proto.VerificationType_VERIFICATION_TYPE_PHONE, nil
	case EmailVerificationType:
		return proto.VerificationType_VERIFICATION_TYPE_EMAIL, nil
	}
	return proto.VerificationType_VERIFICATION_TYPE_INVALID, fmt.Errorf("invalid string verification type %s", verificationType)
}
