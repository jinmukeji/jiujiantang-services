package rest

import (
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	"github.com/kataras/iris/v12"
)

// SecureQuestion 密保问题
type SecureQuestion struct {
	QuestionKey string `json:"question_key"`
	Answer      string `json:"answer"`
}

// ValidateSecureQuestionsBeforeModifyQuestionsBody 修改密保前验证密保的请求
type ValidateSecureQuestionsBeforeModifyQuestionsBody struct {
	SecureQuestions []SecureQuestion `json:"secure_questions"`
}

// ValidateSecureQuestionsBeforeModifyQuestionsReply 修改密保前验证密保的响应
type ValidateSecureQuestionsBeforeModifyQuestionsReply struct {
	Result            bool     `json:"result"`
	WrongQuestionKeys []string `json:"wrong_question_keys"`
}

// ValidateSecureQuestionsBeforeModifyQuestions 修改密保前验证密保
func (h *webHandler) ValidateSecureQuestionsBeforeModifyQuestions(ctx iris.Context) {
	userID, err := ctx.Params().GetInt("user_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	var reqSecureQuestions ValidateSecureQuestionsBeforeModifyQuestionsBody

	err = ctx.ReadJSON(&reqSecureQuestions)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}
	req := new(proto.UserValidateSecureQuestionsBeforeModifyQuestionsRequest)
	secureQuestions := make([]*proto.SecureQuestion, len(reqSecureQuestions.SecureQuestions))
	for idx, item := range reqSecureQuestions.SecureQuestions {
		secureQuestions[idx] = &proto.SecureQuestion{
			QuestionKey: item.QuestionKey,
			Answer:      item.Answer,
		}
	}
	req.SecureQuestions = secureQuestions
	req.UserId = int32(userID)
	resp, errValidateSecureQuestions := h.rpcSvc.UserValidateSecureQuestionsBeforeModifyQuestions(newRPCContext(ctx), req)
	if errValidateSecureQuestions != nil {
		writeRpcInternalError(ctx, errValidateSecureQuestions, false)
		return
	}
	if resp.WrongSecureQuestionKeys == nil {
		rest.WriteOkJSON(ctx, ValidateSecureQuestionsBeforeModifyQuestionsReply{
			Result:            resp.Result,
			WrongQuestionKeys: []string{},
		})
		return
	}
	rest.WriteOkJSON(ctx, ValidateSecureQuestionsBeforeModifyQuestionsReply{
		Result:            resp.Result,
		WrongQuestionKeys: resp.WrongSecureQuestionKeys,
	})
}

// ModifySecureQuestionsRequest 修改密保的请求
type ModifySecureQuestionsRequest struct {
	OldSecureQuestions []SecureQuestion `json:"old_secure_questions"`
	NewSecureQuestions []SecureQuestion `json:"new_secure_questions"`
}

// ModifySecureQuestionsReply 修改密保的响应
type ModifySecureQuestionsReply struct {
	Result                  bool              `json:"result"`
	WrongSecureQuestionKeys []string          `json:"wrong_secure_question_keys"`
	Reasons                 []InvalidQuestion `json:"reasons"`
}

// ModifySecureQuestions 修改密保
func (h *webHandler) ModifySecureQuestions(ctx iris.Context) {
	userID, err := ctx.Params().GetInt("user_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	var reqModifySecureQuestions ModifySecureQuestionsRequest

	err = ctx.ReadJSON(&reqModifySecureQuestions)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}
	req := new(proto.UserModifySecureQuestionsRequest)
	req.UserId = int32(userID)
	oldSecureQuestions := make([]*proto.SecureQuestion, len(reqModifySecureQuestions.OldSecureQuestions))
	for idx, item := range reqModifySecureQuestions.OldSecureQuestions {
		oldSecureQuestions[idx] = &proto.SecureQuestion{
			QuestionKey: item.QuestionKey,
			Answer:      item.Answer,
		}
	}
	req.OldSecureQuestions = oldSecureQuestions
	newSecureQuestions := make([]*proto.SecureQuestion, len(reqModifySecureQuestions.NewSecureQuestions))
	for idx, item := range reqModifySecureQuestions.NewSecureQuestions {
		newSecureQuestions[idx] = &proto.SecureQuestion{
			QuestionKey: item.QuestionKey,
			Answer:      item.Answer,
		}
	}
	req.NewSecureQuestions = newSecureQuestions

	resp, errValidateSecureQuestions := h.rpcSvc.UserModifySecureQuestions(newRPCContext(ctx), req)
	if errValidateSecureQuestions != nil {
		writeRpcInternalError(ctx, errValidateSecureQuestions, false)
		return
	}
	invalidQuestions := make([]InvalidQuestion, len(resp.InvalidSecureQuestions))
	if len(resp.InvalidSecureQuestions) != 0 {
		for idx, item := range resp.InvalidSecureQuestions {
			invalidQuestions[idx].QuestionKey = item.QuestionKey
			invalidDescription := make([]string, len(item.Reason))
			for idxDescription, itemDescription := range item.Reason {
				stringDescription, errmapProtoInvalidQuestionDescriptionToRest := mapProtoInvalidQuestionDescriptionToRest(itemDescription)
				if errmapProtoInvalidQuestionDescriptionToRest != nil {
					writeError(ctx, wrapError(ErrInvalidValue, "", errmapProtoInvalidQuestionDescriptionToRest), false)
					return
				}
				invalidDescription[idxDescription] = stringDescription
			}
			invalidQuestions[idx].Reason = invalidDescription
		}
	}
	if resp.WrongSecureQuestionKeys == nil {
		rest.WriteOkJSON(ctx, ModifySecureQuestionsReply{
			Result:                  resp.Result,
			WrongSecureQuestionKeys: []string{},
			Reasons:                 invalidQuestions,
		})
		return
	}
	rest.WriteOkJSON(ctx, ModifySecureQuestionsReply{
		Result:                  resp.Result,
		WrongSecureQuestionKeys: resp.WrongSecureQuestionKeys,
		Reasons:                 invalidQuestions,
	})
}
