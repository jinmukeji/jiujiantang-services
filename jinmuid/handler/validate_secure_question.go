package handler

import (
	"context"
	"errors"

	"fmt"

	"github.com/jinmukeji/gf-api2/jinmuid/mysqldb"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
)

// SecureQuestion 密保问题
type SecureQuestion struct {
	QuestionKey string // 密保问题的Key
	Answer      string // 密保问题的答案
}

// UserValidateSecureQuestionsBeforeModifyPassword 根据密保重置密码前对密保问题验证
func (j *JinmuIDService) UserValidateSecureQuestionsBeforeModifyPassword(ctx context.Context, req *proto.UserValidateSecureQuestionsBeforeModifyPasswordRequest, resp *proto.UserValidateSecureQuestionsBeforeModifyPasswordResponse) error {
	if req.ValidationType == proto.ValidationType_VALIDATION_TYPE_UNKNOWN {
		return NewError(ErrInvalidSecureQuestionValidationMethod, errors.New("invalid secure queston validation type"))
	}
	err := validateQuestionFormat(req.SecureQuestions)
	if err != nil {
		return err
	}
	var questions []mysqldb.SecureQuestion
	if req.ValidationType == proto.ValidationType_VALIDATION_TYPE_PHONE {
		if req.Username != "" {
			return NewError(ErrInvalidValidationValue, errors.New("email should be empty when validation type is phone"))
		}
		existSignInPhone, errExistSignInPhone := j.datastore.ExistSignInPhone(ctx, req.Phone, req.NationCode)
		if errExistSignInPhone != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to check existence of phone %s%s: %s", req.NationCode, req.Phone, errExistSignInPhone.Error()))
		}
		if !existSignInPhone {
			return NewError(ErrNoneExistentPhone, fmt.Errorf("phone %s%s doesn't exist", req.NationCode, req.Phone))
		}

		questions, err = j.datastore.FindSecureQuestionByPhone(ctx, req.Phone, req.NationCode)
		if err != nil {
			return NewError(ErrCurrentSecureQuestionsNotSet, fmt.Errorf("failed to find secure questions by phone %s%s: %s", req.NationCode, req.Phone, err.Error()))
		}
	}
	if req.ValidationType == proto.ValidationType_VALIDATION_TYPE_USERNAME {
		if req.Phone != "" || req.NationCode != "" {
			return NewError(ErrInvalidValidationValue, errors.New("phone should be empty when validation type is username"))
		}
		existUsername, errExistUsername := j.datastore.ExistUsername(ctx, req.Username)
		if errExistUsername != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to check existence of username %s: %s", req.Username, errExistUsername.Error()))
		}
		if !existUsername {
			return NewError(ErrNonexistentUsername, fmt.Errorf("username %s doesn't exist", req.Username))
		}

		questions, err = j.datastore.FindSecureQuestionByUsername(ctx, req.Username)
		if err != nil {
			return NewError(ErrCurrentSecureQuestionsNotSet, fmt.Errorf("failed to find secure questions by username %s: %s", req.Username, err.Error()))
		}
	}

	// 比较密保是否正确
	wrongQuestions, err := compareSecureQuestion(req.SecureQuestions, questions)
	if err != nil {
		return err
	}
	if len(wrongQuestions) == 0 {
		resp.Result = true
	}
	resp.WrongQuestionKeys = wrongQuestions
	return nil

}

// UserValidateSecureQuestionsBeforeModifyQuestions 修改密保前验证密保
func (j *JinmuIDService) UserValidateSecureQuestionsBeforeModifyQuestions(ctx context.Context, req *proto.UserValidateSecureQuestionsBeforeModifyQuestionsRequest, resp *proto.UserValidateSecureQuestionsBeforeModifyQuestionsResponse) error {
	err := validateQuestionFormat(req.SecureQuestions)
	if err != nil {
		return err
	}
	questions, err := j.datastore.FindSecureQuestionByUserID(ctx, req.UserId)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find secure questions by userID %d: %s", req.UserId, err.Error()))
	}
	// 比较密保是否正确
	wrongQuestion, err := compareSecureQuestion(req.SecureQuestions, questions)
	if err != nil {
		return err
	}
	if len(wrongQuestion) == 0 {
		resp.Result = true
	}
	resp.WrongSecureQuestionKeys = wrongQuestion
	return nil
}

// compareSecureQuestion 比较密保问题是否正确
func compareSecureQuestion(reqQuestionAndAnswers []*proto.SecureQuestion, questions []mysqldb.SecureQuestion) ([]string, error) {

	correctQuestionCount := 0
	wrongQuestion := []string{}
	for _, item := range reqQuestionAndAnswers {
		for _, value := range questions {
			if item.QuestionKey == value.SecureQuestionKey {
				correctQuestionCount++
				if item.Answer != value.SecureAnswer {
					wrongQuestion = append(wrongQuestion, item.QuestionKey)
				}
			}
		}
	}

	if correctQuestionCount < SecureQuestionCount {
		return nil, NewError(ErrMismatchQuestion, errors.New("mismatch secure question"))
	}

	return wrongQuestion, nil
}
