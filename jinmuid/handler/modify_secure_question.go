package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/jinmukeji/gf-api2/jinmuid/mysqldb"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
)

// UserModifySecureQuestions 修改密保
func (j *JinmuIDService) UserModifySecureQuestions(ctx context.Context, req *proto.UserModifySecureQuestionsRequest, resp *proto.UserModifySecureQuestionsResponse) error {
	// 判断UserID是否存在
	exist, errExistUserByUserID := j.datastore.ExistUserByUserID(ctx, req.UserId)
	if errExistUserByUserID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to check existence of userID %d: %s", req.UserId, errExistUserByUserID.Error()))
	}
	if !exist {
		return NewError(ErrInvalidUser, fmt.Errorf("userId %d doesn't exist", req.UserId))
	}
	err := validateQuestionFormat(req.OldSecureQuestions)
	if err != nil {
		return err
	}
	// 验证密保问题的格式
	err = validateQuestionFormat(req.NewSecureQuestions)
	if err != nil {
		return err
	}
	questions, err := j.datastore.FindSecureQuestionByUserID(ctx, req.UserId)
	if err != nil {
		return NewError(ErrCurrentSecureQuestionsNotSet, fmt.Errorf("failed to find secure question by userID %d: %s", req.UserId, err.Error()))
	}
	// 比较密保是否正确
	wrongQuestions, err := compareSecureQuestion(req.OldSecureQuestions, questions)
	if err != nil {
		return err
	}
	if len(wrongQuestions) != 0 {
		resp.WrongSecureQuestionKeys = wrongQuestions
		return nil
	}

	// 验证密保问题的个数
	if len(req.NewSecureQuestions) != SecureQuestionCount {
		return NewError(ErrWrongSecureQuestionCount, fmt.Errorf("wrong secure questions count %d. It should be %d", len(req.NewSecureQuestions), SecureQuestionCount))
	}
	// 判断答案是否满足要求以及问题是否含有敏感词等
	invalidQuestionAndAnswer := isValidSecureQuestionAndAnswer(req.NewSecureQuestions)
	if len(invalidQuestionAndAnswer) != 0 {
		resp.InvalidSecureQuestions = invalidQuestionAndAnswer
	}

	// 新旧密保问题不能一样
	sameQuestions := sameSecureQuestion(req.OldSecureQuestions, req.NewSecureQuestions)
	if !sameQuestions {
		return NewError(ErrSameSecureQuestion, errors.New("new secure questions and old secure questions are same"))
	}
	questionAndAnswers := make([]mysqldb.SecureQuestion, SecureQuestionCount)
	for idx, item := range req.NewSecureQuestions {
		questionAndAnswers[idx].SecureQuestionKey = item.QuestionKey
		questionAndAnswers[idx].SecureAnswer = item.Answer
	}
	if len(resp.InvalidSecureQuestions) == 0 && len(resp.WrongSecureQuestionKeys) == 0 {
		resp.Result = true
		errSetSecureQuestion := j.datastore.SetSecureQuestion(ctx, int32(req.UserId), questionAndAnswers)
		if errSetSecureQuestion != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to set secure questions for userId %d: %s", req.UserId, errSetSecureQuestion.Error()))
		}
	}
	return nil

}

// 判断新旧密保问题是否完全一样
func sameSecureQuestion(oldSecureQuestions []*proto.SecureQuestion, newSecureQuestions []*proto.SecureQuestion) bool {
	sameQuestionCount := 0
	for _, oldItem := range oldSecureQuestions {
		for _, newItem := range newSecureQuestions {
			if oldItem.QuestionKey == newItem.QuestionKey {
				if oldItem.Answer == newItem.Answer {
					sameQuestionCount++
				}
			}
		}
	}
	return sameQuestionCount < SecureQuestionCount
}
