package handler

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/jinmukeji/jiujiantang-services/jinmuid/mysqldb"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
)

const (
	// SecureQuestionCount 密保问题的数量
	SecureQuestionCount = 3
)

// UserSetSecureQuestions 用户设置密保问题
func (j *JinmuIDService) UserSetSecureQuestions(ctx context.Context, req *proto.UserSetSecureQuestionsRequest, resp *proto.UserSetSecureQuestionsResponse) error {
	// 验证密保问题的格式
	err := validateQuestionFormat(req.SecureQuestions)
	if err != nil {
		return err
	}
	// 判断UserID是否存在
	exist, errExistUserByUserID := j.datastore.ExistUserByUserID(ctx, req.UserId)
	if errExistUserByUserID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to check user existence by userID %d", req.UserId))
	}
	if !exist {
		return NewError(ErrInvalidUser, fmt.Errorf("userId %d doesn't exist", req.UserId))
	}
	// 判断密保问题是否已经设置过
	ok, _ := j.datastore.ExistsSecureQuestion(ctx, int32(req.UserId))
	if ok {
		return NewError(ErrSecureQuestionExist, fmt.Errorf("user %d has set secure questions before", req.UserId))
	}
	// 验证密保问题的个数
	if len(req.SecureQuestions) != SecureQuestionCount {
		return NewError(ErrWrongSecureQuestionCount, fmt.Errorf("wrong secure questions count %d. It should be %d", len(req.SecureQuestions), SecureQuestionCount))
	}
	// 判断答案是否满足要求以及问题是否含有敏感词等
	invalidQuestionAndAnswer := isValidSecureQuestionAndAnswer(req.SecureQuestions)
	if len(invalidQuestionAndAnswer) != 0 {
		resp.InvalidSecureQuestions = invalidQuestionAndAnswer
		return nil
	}
	questionAndAnswers := make([]mysqldb.SecureQuestion, SecureQuestionCount)
	for idx, item := range req.SecureQuestions {
		questionAndAnswers[idx].SecureQuestionKey = item.QuestionKey
		questionAndAnswers[idx].SecureAnswer = item.Answer
	}
	errSetSecureQuestion := j.datastore.SetSecureQuestion(ctx, int32(req.UserId), questionAndAnswers)
	if errSetSecureQuestion != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to set secure questions for user %d: %s", req.UserId, errSetSecureQuestion.Error()))
	}
	resp.Result = true
	return nil

}

// isValidSecureQuestionAndAnswer 验证设置密保的答案
func isValidSecureQuestionAndAnswer(questionAndAnswer []*proto.SecureQuestion) []*proto.InvalidQuestion {
	invalidQuestions := []*proto.InvalidQuestion{}
	for _, item := range questionAndAnswer {
		invalidDescription := []proto.InvalidQuestionDescription{}
		if item.Answer == "" {
			invalidDescription = append(invalidDescription, proto.InvalidQuestionDescription_INVALID_QUESTION_DESCRIPTION_EMPTY_ANSWER)
		}
		if utf8.RuneCountInString(item.Answer) > 15 {
			invalidDescription = append(invalidDescription, proto.InvalidQuestionDescription_INVALID_QUESTION_DESCRIPTION_TOO_LONG_ANSWER)
		}
		// 判断答案是否包含敏感词
		sensitiveWords := GetSensitiveWords(item.Answer)
		if len(sensitiveWords) != 0 {
			invalidDescription = append(invalidDescription, proto.InvalidQuestionDescription_INVALID_QUESTION_DESCRIPTION_SENSITIVE_WORD_IN_ANSWER)
		}
		// 判断答案是否包含保留词
		reservedWords := GetReservedWords(item.Answer)
		if len(reservedWords) != 0 {
			invalidDescription = append(invalidDescription, proto.InvalidQuestionDescription_INVALID_QUESTION_DESCRIPTION_RESERVED_WORD_IN_ANSWER)
		}
		// 判断答案是否包含屏蔽词
		maskWords := GetMaskWords(item.Answer)
		if len(maskWords) != 0 {
			invalidDescription = append(invalidDescription, proto.InvalidQuestionDescription_INVALID_QUESTION_DESCRIPTION_MASK_WORD_IN_ANSWER)
		}
		// 答案限制为中英文或数字
		for _, r := range item.Answer {
			if !(unicode.Is(unicode.Scripts["Han"], r) || unicode.IsDigit(r) || unicode.IsLetter(r)) {
				invalidDescription = append(invalidDescription, proto.InvalidQuestionDescription_INVALID_QUESTION_DESCRIPTION_ILLEGAL_CHARACTER_IN_ANSWER)
				break
			}
		}
		if len(invalidDescription) != 0 {
			invalidQuestions = append(invalidQuestions, &proto.InvalidQuestion{
				QuestionKey: item.QuestionKey,
				Reason:      invalidDescription,
			})
		}
	}

	return invalidQuestions
}

// validateQuestionFormat 验证密保问题的格式
func validateQuestionFormat(questionAndAnswers []*proto.SecureQuestion) error {
	// 验证密保问题的个数
	if len(questionAndAnswers) != SecureQuestionCount {
		return NewError(ErrWrongSecureQuestionCount, fmt.Errorf("wrong secure questions count %d. It should be %d", len(questionAndAnswers), SecureQuestionCount))
	}

	for _, item := range questionAndAnswers {
		// 答案不可以为空
		if item.Answer == "" {
			return NewError(ErrEmptyAnswer, errors.New("all answers should not be empty"))
		}
		// 验证问题的格式
		intQuestionKey, errGetIntQuestionKey := strconv.Atoi(item.QuestionKey)
		if errGetIntQuestionKey != nil {
			return NewError(ErrWrongFormatQuestion, errors.New("wrong format of question"))
		}
		if intQuestionKey < 0 || intQuestionKey > len(SecureQuestions) {
			return NewError(ErrWrongFormatQuestion, errors.New("wrong format of question key"))
		}
	}

	for _, item := range questionAndAnswers {
		if item.QuestionKey == "" {
			return NewError(ErrEmptyAnswer, errors.New("all answer keys should not be empty"))
		}
	}
	// 所有的问题不能一样
	sameQuestion := false
	for idx1 := 0; idx1 < SecureQuestionCount; idx1++ {
		for idx2 := idx1 + 1; idx2 < SecureQuestionCount; idx2++ {
			if questionAndAnswers[idx1].QuestionKey == questionAndAnswers[idx2].QuestionKey {
				sameQuestion = true
			}
		}
	}
	if sameQuestion {
		return NewError(ErrRepeatedQuestion, errors.New("repeated question"))
	}
	return nil

}
