package rest

import (
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"

	"github.com/kataras/iris/v12"
)

// GetSecureQuestionsByPhoneOrUsernameRequest 根据用户名或者手机号获取当前设置的密保问题请求
type GetSecureQuestionsByPhoneOrUsernameRequest struct {
	ValidationType string `json:"validation_type"` // 获取方式
	Username       string `json:"username"`        // 用户名
	Phone          string `json:"phone"`           // 手机号码
	NationCode     string `json:"nation_code"`     // 区号
}

// GetSecureQuestionsByPhoneOrUsername 根据用户名或者手机号获取当前设置的密保问题响应
func (h *webHandler) GetSecureQuestionsByPhoneOrUsername(ctx iris.Context) {
	var reqGetSecureQuestion GetSecureQuestionsByPhoneOrUsernameRequest
	err := ctx.ReadJSON(&reqGetSecureQuestion)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}
	req := new(proto.GetSecureQuestionsByPhoneOrUsernameRequest)
	validationType, errmapRestValidationTypeToProto := mapRestValidationTypeToProto(reqGetSecureQuestion.ValidationType)
	if errmapRestValidationTypeToProto != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errmapRestValidationTypeToProto), false)
		return
	}
	req.ValidationType = validationType
	req.Username = reqGetSecureQuestion.Username
	req.Phone = reqGetSecureQuestion.Phone
	req.NationCode = reqGetSecureQuestion.NationCode
	resp, errValidateUsernameAndPhone := h.rpcSvc.GetSecureQuestionsByPhoneOrUsername(newRPCContext(ctx), req)
	if errValidateUsernameAndPhone != nil {
		writeRpcInternalError(ctx, errValidateUsernameAndPhone, false)
		return
	}
	secureKeyAndQuestions := make([]SecureKeyAndQuestion, len(resp.SecureQuestions))
	for idx, item := range resp.SecureQuestions {
		secureKeyAndQuestions[idx].Key = item.Key
		secureKeyAndQuestions[idx].Question = item.Question
	}
	rest.WriteOkJSON(ctx, secureKeyAndQuestions)
}
