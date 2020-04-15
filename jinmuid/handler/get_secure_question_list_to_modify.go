package handler

import (
	"context"
	"fmt"

	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
)

// GetSecureQuestionListToModify 修改密保前获取已设置密保列表
func (j *JinmuIDService) GetSecureQuestionListToModify(ctx context.Context, req *proto.GetSecureQuestionListToModifyRequest, resp *proto.GetSecureQuestionListToModifyResponse) error {
	// 判断UserID是否存在
	exist, errExistUserByUserID := j.datastore.ExistUserByUserID(ctx, req.UserId)
	if errExistUserByUserID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to check existence of user %d: %s", req.UserId, errExistUserByUserID.Error()))
	}
	if !exist {
		return NewError(ErrInvalidUser, fmt.Errorf("userId %d doesn't exist", req.UserId))
	}
	secureQuestions := getSecureQuestionsByLanguage()
	setSecureQuestions, errGetSecureQuestionListToModifyByUserID := j.datastore.GetSecureQuestionListToModifyByUserID(ctx, req.UserId)
	if errGetSecureQuestionListToModifyByUserID != nil {
		return NewError(ErrCurrentSecureQuestionsNotSet, fmt.Errorf("failed to get secure question list to modify by user %d: %s", req.UserId, errGetSecureQuestionListToModifyByUserID.Error()))
	}
	protoSecureQuestions := make([]*proto.SecureQuestionKeyAndQuestion, len(setSecureQuestions))
	for idx, item := range setSecureQuestions {
		protoSecureQuestions[idx] = &proto.SecureQuestionKeyAndQuestion{
			Key:      item,
			Question: secureQuestions[item],
		}
	}
	resp.SecureQuestions = protoSecureQuestions
	return nil
}
