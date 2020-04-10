package rest

import (
	"fmt"

	"github.com/jinmukeji/gf-api2/pkg/rest"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	"github.com/kataras/iris/v12"
)

const (
	// EmptyAnswer 答案为空
	EmptyAnswer = "答案为空"
	// TooLongAnswer 答案过长
	TooLongAnswer = "答案过长"
	// IllegalCharacterInAnswer 答案限制为中英文或数字
	IllegalCharacterInAnswer = "答案限制为中英文或数字"
	// ReservedWordInAnswer 答案含有保留词
	ReservedWordInAnswer = "答案含有保留词"
	// MaskWordInAnswer 答案含有屏蔽词
	MaskWordInAnswer = "答案含有屏蔽词"
	// SensitiveWordInAnswer 答案含有敏感词
	SensitiveWordInAnswer = "答案含有敏感词"
	// UnknownDescription 未知的错误原因
	UnknownDescription = "未知的错误原因"
)

// SetSecureQuestionsBody 设置密保问题的请求
type SetSecureQuestionsBody struct {
	SecureQuestions []SecureQuestion `json:"secure_questions"`
}

// SetSecureQuestionsReply 设置密保问题的响应
type SetSecureQuestionsReply struct {
	Result           bool              `json:"result"`
	InvalidQuestions []InvalidQuestion `json:"invalid_questions"`
}

// InvalidQuestion 不符合规范的密保问题
type InvalidQuestion struct {
	QuestionKey string   `json:"question_key"`
	Reason      []string `json:"reasons"`
}

// SetSecureQuestions 设置密保问题
func (h *webHandler) SetSecureQuestions(ctx iris.Context) {
	var reqSecureQuestions SetSecureQuestionsBody

	err := ctx.ReadJSON(&reqSecureQuestions)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}
	userID, err := ctx.Params().GetInt("user_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	SecureQuestions := make([]*proto.SecureQuestion, len(reqSecureQuestions.SecureQuestions))
	for idx, item := range reqSecureQuestions.SecureQuestions {
		SecureQuestions[idx] = &proto.SecureQuestion{
			QuestionKey: item.QuestionKey,
			Answer:      item.Answer,
		}
	}
	req := new(proto.UserSetSecureQuestionsRequest)
	req.UserId = int32(userID)
	req.SecureQuestions = SecureQuestions

	resp, err := h.rpcSvc.UserSetSecureQuestions(
		newRPCContext(ctx), req,
	)
	if err != nil {
		writeRpcInternalError(ctx, err, false)
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
	rest.WriteOkJSON(ctx, SetSecureQuestionsReply{
		Result:           resp.Result,
		InvalidQuestions: invalidQuestions,
	})
}

// mapProtoInvalidQuestionDescriptionToRest 将错误问题的描述转换为 rest 层使用的 string 类型
func mapProtoInvalidQuestionDescriptionToRest(protoDescription proto.InvalidQuestionDescription) (string, error) {
	switch protoDescription {
	case proto.InvalidQuestionDescription_INVALID_QUESTION_DESCRIPTION_INVALID:
		return "", fmt.Errorf("invalid proto question description %d", protoDescription)
	case proto.InvalidQuestionDescription_INVALID_QUESTION_DESCRIPTION_UNSET:
		return "", fmt.Errorf("invalid proto question description %d", protoDescription)
	case proto.InvalidQuestionDescription_INVALID_QUESTION_DESCRIPTION_EMPTY_ANSWER:
		return EmptyAnswer, nil
	case proto.InvalidQuestionDescription_INVALID_QUESTION_DESCRIPTION_TOO_LONG_ANSWER:
		return TooLongAnswer, nil
	case proto.InvalidQuestionDescription_INVALID_QUESTION_DESCRIPTION_ILLEGAL_CHARACTER_IN_ANSWER:
		return IllegalCharacterInAnswer, nil
	case proto.InvalidQuestionDescription_INVALID_QUESTION_DESCRIPTION_RESERVED_WORD_IN_ANSWER:
		return ReservedWordInAnswer, nil
	case proto.InvalidQuestionDescription_INVALID_QUESTION_DESCRIPTION_MASK_WORD_IN_ANSWER:
		return MaskWordInAnswer, nil
	case proto.InvalidQuestionDescription_INVALID_QUESTION_DESCRIPTION_SENSITIVE_WORD_IN_ANSWER:
		return SensitiveWordInAnswer, nil
	}
	return "", fmt.Errorf("invalid proto question description %d", protoDescription)
}
