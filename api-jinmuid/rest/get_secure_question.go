package rest

import (
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	"github.com/kataras/iris/v12"
)

// SecureKeyAndQuestion 序号和对应的密保问题
type SecureKeyAndQuestion struct {
	Key      string `json:"key"`
	Question string `json:"question"`
}

// GetSecureQuestionList 获取密保问题列表
func (h *webHandler) GetSecureQuestionList(ctx iris.Context) {
	userID, err := ctx.Params().GetInt("user_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	req := new(proto.UserGetSecureQuestionListRequest)
	req.UserId = int32(userID)
	resp, errGetSecureQuestionList := h.rpcSvc.UserGetSecureQuestionList(newRPCContext(ctx), req)
	if errGetSecureQuestionList != nil {
		writeRpcInternalError(ctx, errGetSecureQuestionList, false)
		return
	}
	secureKeyAndQuestions := make([]SecureKeyAndQuestion, len(resp.SecureQuestions))
	for idx, item := range resp.SecureQuestions {
		secureKeyAndQuestions[idx].Key = item.Key
		secureKeyAndQuestions[idx].Question = item.Question
	}
	rest.WriteOkJSON(ctx, secureKeyAndQuestions)
}
