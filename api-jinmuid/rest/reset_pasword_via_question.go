package rest

import (
	"fmt"

	"github.com/jinmukeji/gf-api2/pkg/rest"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	"github.com/kataras/iris/v12"
)

const (
	// ValidationTypePhone 验证方式为手机号
	ValidationTypePhone = "phone"
	// ValidationTypeUsername 验证方式为用户名
	ValidationTypeUsername = "username"
)

// ValidateUsernameOrPhoneRequest 验证手机号码和用户名是否存在的请求
type ValidateUsernameOrPhoneRequest struct {
	ValidationType string `json:"validation_type"` // 验证方式
	Username       string `json:"username"`        // 用户名
	Phone          string `json:"phone"`           // 手机号码
	NationCode     string `json:"nation_code"`     // 区号
}

// ValidateUsernameOrPhoneResponse 验证手机号码和用户名是否存在的响应
type ValidateUsernameOrPhoneResponse struct {
	SecureEmail       string `json:"secure_email"`         // 安全邮箱
	HasSetSecureEmail bool   `json:"has_set_secure_email"` // 是否设置安全邮箱
}

// ValidateUsernameOrPhone 验证用户名或者手机号码是否存在
func (h *webHandler) ValidateUsernameOrPhone(ctx iris.Context) {
	var reqValidateUsernameOrPhone ValidateUsernameOrPhoneRequest
	err := ctx.ReadJSON(&reqValidateUsernameOrPhone)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}
	req := new(proto.UserValidateUsernameOrPhoneRequest)
	validationType, errmapRestValidationTypeToProto := mapRestValidationTypeToProto(reqValidateUsernameOrPhone.ValidationType)
	if errmapRestValidationTypeToProto != nil {
		writeError(ctx, wrapError(ErrInvalidSecureQuestionValidationMethod, "", errmapRestValidationTypeToProto), false)
		return
	}
	req.ValidationType = validationType
	req.Username = reqValidateUsernameOrPhone.Username
	req.Phone = reqValidateUsernameOrPhone.Phone
	req.NationCode = reqValidateUsernameOrPhone.NationCode
	respUserValidateUsernameOrPhone, errValidateUsernameAndPhone := h.rpcSvc.UserValidateUsernameOrPhone(newRPCContext(ctx), req)
	if errValidateUsernameAndPhone != nil {
		writeRpcInternalError(ctx, errValidateUsernameAndPhone, false)
		return
	}
	rest.WriteOkJSON(ctx, ValidateUsernameOrPhoneResponse{
		SecureEmail:       respUserValidateUsernameOrPhone.SecureEmail,
		HasSetSecureEmail: respUserValidateUsernameOrPhone.HasSetEmail,
	})
}

// ValidateSecureQuestionsBeforeModifyPasswordRequest 根据密保重置密码前对密保问题验证的请求
type ValidateSecureQuestionsBeforeModifyPasswordRequest struct {
	ValidationType  string           `json:"validation_type"`
	Username        string           `json:"username"`
	Phone           string           `json:"phone"`
	NationCode      string           `json:"nation_code"`
	SecureQuestions []SecureQuestion `json:"secure_questions"`
}

// ValidateSecureQuestionsBeforeModifyPasswordReply 根据密保重置密码前对密保问题验证的响应
type ValidateSecureQuestionsBeforeModifyPasswordReply struct {
	Result            bool     `json:"result"`
	WrongQuestionKeys []string `json:"wrong_question_keys"`
}

// ValidateSecureQuestionsBeforeModifyPassword 根据密保重置密码前对密保问题验证
func (h *webHandler) ValidateSecureQuestionsBeforeModifyPassword(ctx iris.Context) {
	var reqValidateSecureQuestionsBeforeModifyPassword ValidateSecureQuestionsBeforeModifyPasswordRequest
	err := ctx.ReadJSON(&reqValidateSecureQuestionsBeforeModifyPassword)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}
	req := new(proto.UserValidateSecureQuestionsBeforeModifyPasswordRequest)
	protoVerificationType, errmapRestValidationTypeToProto := mapRestValidationTypeToProto(reqValidateSecureQuestionsBeforeModifyPassword.ValidationType)
	if errmapRestValidationTypeToProto != nil {
		writeError(ctx, wrapError(ErrInvalidSecureQuestionValidationMethod, "", errmapRestValidationTypeToProto), false)
		return
	}
	req.ValidationType = protoVerificationType
	req.Username = reqValidateSecureQuestionsBeforeModifyPassword.Username
	req.Phone = reqValidateSecureQuestionsBeforeModifyPassword.Phone
	req.NationCode = reqValidateSecureQuestionsBeforeModifyPassword.NationCode
	secureQuestions := make([]*proto.SecureQuestion, len(reqValidateSecureQuestionsBeforeModifyPassword.SecureQuestions))
	for idx, item := range reqValidateSecureQuestionsBeforeModifyPassword.SecureQuestions {
		secureQuestions[idx] = &proto.SecureQuestion{
			QuestionKey: item.QuestionKey,
			Answer:      item.Answer,
		}
	}
	req.SecureQuestions = secureQuestions
	resp, errValidateSecureQuestion := h.rpcSvc.UserValidateSecureQuestionsBeforeModifyPassword(newRPCContext(ctx), req)
	if errValidateSecureQuestion != nil {
		writeRpcInternalError(ctx, errValidateSecureQuestion, false)
		return
	}
	if resp.WrongQuestionKeys == nil {
		rest.WriteOkJSON(ctx, ValidateSecureQuestionsBeforeModifyPasswordReply{
			Result:            resp.Result,
			WrongQuestionKeys: []string{},
		})
		return
	}
	rest.WriteOkJSON(ctx, ValidateSecureQuestionsBeforeModifyPasswordReply{
		Result:            resp.Result,
		WrongQuestionKeys: resp.WrongQuestionKeys,
	})
}

// mapRestValidationTypeToProto 将 rest 使用的 string 类型的 validation_type 转换为 proto 类型
func mapRestValidationTypeToProto(validationType string) (proto.ValidationType, error) {
	switch validationType {
	case ValidationTypePhone:
		return proto.ValidationType_VALIDATION_TYPE_PHONE, nil
	case ValidationTypeUsername:
		return proto.ValidationType_VALIDATION_TYPE_USERNAME, nil
	}
	return proto.ValidationType_VALIDATION_TYPE_INVALID, fmt.Errorf("invalid string validation type %s", validationType)
}

// ResetPasswordViaSecureQuestionsRequest 根据密保问题重置密码的请求
type ResetPasswordViaSecureQuestionsRequest struct {
	ValidationType  string           `json:"validation_type"`
	Username        string           `json:"username"`
	Phone           string           `json:"phone"`
	NationCode      string           `json:"nation_code"`
	Password        string           `json:"password"`
	SecureQuestions []SecureQuestion `json:"secure_questions"`
}

// ResetPasswordViaSecureQuestionsReply 根据密保问题重置密码的响应
type ResetPasswordViaSecureQuestionsReply struct {
	Result            bool     `json:"result"`
	WrongQuestionKeys []string `json:"wrong_question_keys"`
}

// ResetPasswordViaSecureQuestions 根据密保问题重置密码
func (h *webHandler) ResetPasswordViaSecureQuestions(ctx iris.Context) {
	var reqResetPasswordViaSecureQuestions ResetPasswordViaSecureQuestionsRequest
	err := ctx.ReadJSON(&reqResetPasswordViaSecureQuestions)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}
	req := new(proto.UserResetPasswordViaSecureQuestionsRequest)
	protoVerificationType, errmapRestValidationTypeToProto := mapRestValidationTypeToProto(reqResetPasswordViaSecureQuestions.ValidationType)
	if errmapRestValidationTypeToProto != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errmapRestValidationTypeToProto), false)
		return
	}
	req.ValidationType = protoVerificationType
	req.Username = reqResetPasswordViaSecureQuestions.Username
	req.Phone = reqResetPasswordViaSecureQuestions.Phone
	req.NationCode = reqResetPasswordViaSecureQuestions.NationCode
	req.Password = reqResetPasswordViaSecureQuestions.Password
	secureQuestions := make([]*proto.SecureQuestion, len(reqResetPasswordViaSecureQuestions.SecureQuestions))
	for idx, item := range reqResetPasswordViaSecureQuestions.SecureQuestions {
		secureQuestions[idx] = &proto.SecureQuestion{
			QuestionKey: item.QuestionKey,
			Answer:      item.Answer,
		}
	}
	req.SecureQuestions = secureQuestions
	resp, errValidateSecureQuestion := h.rpcSvc.UserResetPasswordViaSecureQuestions(newRPCContext(ctx), req)
	if errValidateSecureQuestion != nil {
		writeRpcInternalError(ctx, errValidateSecureQuestion, false)
		return
	}
	rest.WriteOkJSON(ctx, ResetPasswordViaSecureQuestionsReply{
		Result:            resp.Result,
		WrongQuestionKeys: resp.WrongSecureQuestionKeys,
	})
}
