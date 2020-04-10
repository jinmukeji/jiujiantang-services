package handler

import (
	"context"
	"errors"

	"fmt"

	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
)

// GetSecureQuestionsByPhoneOrUsername 根据用户名或者手机号获取当前设置的密保问题
func (j *JinmuIDService) GetSecureQuestionsByPhoneOrUsername(ctx context.Context, req *proto.GetSecureQuestionsByPhoneOrUsernameRequest, resp *proto.GetSecureQuestionsByPhoneOrUsernameResponse) error {
	if req.ValidationType == proto.ValidationType_VALIDATION_TYPE_UNKNOWN {
		return NewError(ErrInvalidValidationMethod, errors.New("invalid secure queston validation type"))
	}
	var secureQuestions []string
	if req.ValidationType == proto.ValidationType_VALIDATION_TYPE_PHONE {
		if req.Username != "" {
			return NewError(ErrInvalidValidationValue, errors.New("non-empty username when getting secure questions by phone"))
		}
		var errGetSecureQuestionsByPhone error

		secureQuestions, errGetSecureQuestionsByPhone = j.datastore.GetSecureQuestionsByPhone(ctx, req.NationCode, req.Phone)
		if errGetSecureQuestionsByPhone != nil {
			return NewError(ErrCurrentSecureQuestionsNotSet, fmt.Errorf("failed to get secure questions by phone %s%s: %s", req.NationCode, req.Phone, errGetSecureQuestionsByPhone.Error()))
		}
	}
	if req.ValidationType == proto.ValidationType_VALIDATION_TYPE_USERNAME {
		if req.Phone != "" || req.NationCode != "" {
			return NewError(ErrInvalidValidationValue, errors.New("non-empty phone when getting secure questions by username"))
		}
		var errGetSecureQuestionsByUsername error
		secureQuestions, errGetSecureQuestionsByUsername = j.datastore.GetSecureQuestionsByUsername(ctx, req.Username)
		if errGetSecureQuestionsByUsername != nil {
			return NewError(ErrCurrentSecureQuestionsNotSet, fmt.Errorf("there are no secure questions set for username %s: %s", req.Username, errGetSecureQuestionsByUsername.Error()))
		}
	}
	protoSecureQuestions := make([]*proto.SecureQuestionKeyAndQuestion, len(secureQuestions))
	for idx, item := range secureQuestions {
		protoSecureQuestions[idx] = &proto.SecureQuestionKeyAndQuestion{
			Key:      item,
			Question: SecureQuestions[item],
		}
	}
	resp.SecureQuestions = protoSecureQuestions
	return nil
}
