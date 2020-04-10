package rest

import (
	"github.com/jinmukeji/gf-api2/pkg/rest"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	"github.com/kataras/iris/v12"
)

// GetSecureQuestionListToModify 修改密保前获取已设置密保列表
func (h *webHandler) GetSecureQuestionListToModify(ctx iris.Context) {
	userID, err := ctx.Params().GetInt("user_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	req := new(proto.GetSecureQuestionListToModifyRequest)
	req.UserId = int32(userID)
	resp, errGetSecureQuestionListToModify := h.rpcSvc.GetSecureQuestionListToModify(newRPCContext(ctx), req)
	if errGetSecureQuestionListToModify != nil {
		writeRpcInternalError(ctx, errGetSecureQuestionListToModify, false)
		return
	}
	secureKeyAndQuestions := make([]SecureKeyAndQuestion, len(resp.SecureQuestions))
	for idx, item := range resp.SecureQuestions {
		secureKeyAndQuestions[idx].Key = item.Key
		secureKeyAndQuestions[idx].Question = item.Question
	}
	rest.WriteOkJSON(ctx, secureKeyAndQuestions)
}
