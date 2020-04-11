package handler

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// UserValidateSecureQuestionTestSuite 对密保问题验证测试
type UserValidateSecureQuestionTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

// UserValidateSecureQuestionTestSuite 设置测试环境
func (suite *UserValidateSecureQuestionTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestUserValidateSecureQuestionBeforeModifyPassword 根据密保问题修改密码前验证回答的密保问题是否正确手机号
func (suite *UserValidateSecureQuestionTestSuite) TestUserValidateSecureQuestionBeforeModifyPassword() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.UserValidateSecureQuestionsBeforeModifyPasswordRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_PHONE,
		Phone:          suite.Account.phone,
		NationCode:     suite.Account.nationCode,
		SecureQuestions: []*proto.SecureQuestion{
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1New,
				Answer:      suite.Account.answer1New,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey2New,
				Answer:      suite.Account.answer2New,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey3New,
				Answer:      suite.Account.answer3New,
			},
		},
	}
	resp := new(proto.UserValidateSecureQuestionsBeforeModifyPasswordResponse)
	err := suite.JinmuIDService.UserValidateSecureQuestionsBeforeModifyPassword(ctx, req, resp)
	assert.NoError(t, err)
	assert.Equal(t, true, resp.Result)
}

// TestUserValidateSecureQuestionBeforeModifyPasswordUsername 根据密保问题修改密码前验证回答的密保问题是否正确（用户名）
func (suite *UserValidateSecureQuestionTestSuite) TestUserValidateSecureQuestionBeforeModifyPasswordUsername() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.UserValidateSecureQuestionsBeforeModifyPasswordRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_USERNAME,
		Username:       suite.Account.username,
		SecureQuestions: []*proto.SecureQuestion{
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1,
				Answer:      suite.Account.answer1,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey2,
				Answer:      suite.Account.answer2,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey3,
				Answer:      suite.Account.answer3,
			},
		},
	}
	resp := new(proto.UserValidateSecureQuestionsBeforeModifyPasswordResponse)
	err := suite.JinmuIDService.UserValidateSecureQuestionsBeforeModifyPassword(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:49000] secure questions are not set"), err)
	assert.Equal(t, false, resp.Result)
}

// TestUserValidateSecureQuestionBeforeModifyQuestion 测试修改密保前验证密保
func (suite *UserValidateSecureQuestionTestSuite) TestUserValidateSecureQuestionBeforeModifyQuestion() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	req := &proto.UserValidateSecureQuestionsBeforeModifyQuestionsRequest{
		UserId: userID,
		SecureQuestions: []*proto.SecureQuestion{
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1,
				Answer:      suite.Account.answer1,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey2,
				Answer:      suite.Account.answer2,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey3,
				Answer:      suite.Account.answer3,
			},
		},
	}
	resp := new(proto.UserValidateSecureQuestionsBeforeModifyQuestionsResponse)
	err = suite.JinmuIDService.UserValidateSecureQuestionsBeforeModifyQuestions(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:17000] mismatch secure question"), err)
}

// TestUserValidateSecureQuestionBeforeModifyPasswordQuestionIsNull  问题为空
func (suite *UserValidateSecureQuestionTestSuite) TestUserValidateSecureQuestionBeforeModifyPasswordQuestionIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.UserValidateSecureQuestionsBeforeModifyPasswordRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_PHONE,
		Phone:          suite.Account.phone,
		NationCode:     suite.Account.nationCode,
		SecureQuestions: []*proto.SecureQuestion{
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1Null,
				Answer:      suite.Account.answer1,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey2Null,
				Answer:      suite.Account.answer2,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey3Null,
				Answer:      suite.Account.answer3,
			},
		},
	}
	resp := new(proto.UserValidateSecureQuestionsBeforeModifyPasswordResponse)
	err := suite.JinmuIDService.UserValidateSecureQuestionsBeforeModifyPassword(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:17000] mismatch secure question"), err)
}

// TestUserValidateSecureQuestionBeforeModifyPasswordAnswerIsNull 答案为空
func (suite *UserValidateSecureQuestionTestSuite) TestUserValidateSecureQuestionBeforeModifyPasswordAnswerIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.UserValidateSecureQuestionsBeforeModifyPasswordRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_PHONE,
		Phone:          suite.Account.phone,
		NationCode:     suite.Account.nationCode,
		SecureQuestions: []*proto.SecureQuestion{
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1,
				Answer:      suite.Account.answer1Null,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey2,
				Answer:      suite.Account.answer2Null,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey3,
				Answer:      suite.Account.answer3,
			},
		},
	}
	resp := new(proto.UserValidateSecureQuestionsBeforeModifyPasswordResponse)
	err := suite.JinmuIDService.UserValidateSecureQuestionsBeforeModifyPassword(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:17000] mismatch secure question"), err)
}

// TestUserValidateSecureQuestionBeforeModifyPasswordPhoneFormatError 根据密保问题修改密码前验证回答的密保问题是否正确手机号
func (suite *UserValidateSecureQuestionTestSuite) TestUserValidateSecureQuestionBeforeModifyPasswordPhoneFormatError() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.UserValidateSecureQuestionsBeforeModifyPasswordRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_PHONE,
		Phone:          suite.Account.PhoneError,
		NationCode:     suite.Account.nationCode,
		SecureQuestions: []*proto.SecureQuestion{
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1,
				Answer:      suite.Account.answer1,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey2,
				Answer:      suite.Account.answer2,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey3,
				Answer:      suite.Account.answer3,
			},
		},
	}
	resp := new(proto.UserValidateSecureQuestionsBeforeModifyPasswordResponse)
	err := suite.JinmuIDService.UserValidateSecureQuestionsBeforeModifyPassword(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:22000] nonexistent phone"), err)
	assert.Equal(t, false, resp.Result)
}

// TestUserValidateSecureQuestionBeforeModifyPasswordPhoneIsNull 根据密保问题修改密码前验证回答的密保问题是否正确手机号
func (suite *UserValidateSecureQuestionTestSuite) TestUserValidateSecureQuestionBeforeModifyPasswordPhoneIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.UserValidateSecureQuestionsBeforeModifyPasswordRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_PHONE,
		Phone:          suite.Account.phoneIsNull,
		NationCode:     suite.Account.nationCode,
		SecureQuestions: []*proto.SecureQuestion{
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1,
				Answer:      suite.Account.answer1,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey2,
				Answer:      suite.Account.answer2,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey3,
				Answer:      suite.Account.answer3,
			},
		},
	}
	resp := new(proto.UserValidateSecureQuestionsBeforeModifyPasswordResponse)
	err := suite.JinmuIDService.UserValidateSecureQuestionsBeforeModifyPassword(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:22000] nonexistent phone"), err)
	assert.Equal(t, false, resp.Result)
}

// TestUserValidateSecureQuestionBeforeModifyPasswordValidationTypeIsNull
func (suite *UserValidateSecureQuestionTestSuite) TestUserValidateSecureQuestionBeforeModifyPasswordValidationTypeIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.UserValidateSecureQuestionsBeforeModifyPasswordRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_UNKNOWN,
		Phone:          suite.Account.phone,
		NationCode:     suite.Account.nationCode,
		SecureQuestions: []*proto.SecureQuestion{
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1,
				Answer:      suite.Account.answer1,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey2,
				Answer:      suite.Account.answer2,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey3,
				Answer:      suite.Account.answer3,
			},
		},
	}
	resp := new(proto.UserValidateSecureQuestionsBeforeModifyPasswordResponse)
	err := suite.JinmuIDService.UserValidateSecureQuestionsBeforeModifyPassword(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:11000] invalid secure queston validation type"), err)
	assert.Equal(t, false, resp.Result)
}

// TestUserValidateSecureQuestionBeforeModifyPasswordUsernameIsNull 根据密保问题修改密码前验证回答的密保问题是否正确（用户名）
func (suite *UserValidateSecureQuestionTestSuite) TestUserValidateSecureQuestionBeforeModifyPasswordUsernameIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.UserValidateSecureQuestionsBeforeModifyPasswordRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_USERNAME,
		Username:       suite.Account.usernameNull,
		SecureQuestions: []*proto.SecureQuestion{
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1,
				Answer:      suite.Account.answer1,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey2,
				Answer:      suite.Account.answer2,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey3,
				Answer:      suite.Account.answer3,
			},
		},
	}
	resp := new(proto.UserValidateSecureQuestionsBeforeModifyPasswordResponse)
	err := suite.JinmuIDService.UserValidateSecureQuestionsBeforeModifyPassword(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:21000] nonexistent username"), err)
	assert.Equal(t, false, resp.Result)
}

// TestUserValidateSecureQuestionBeforeModifyQuestionError  密保问题错误
func (suite *UserValidateSecureQuestionTestSuite) TestUserValidateSecureQuestionBeforeModifyQuestionError() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.UserValidateSecureQuestionsBeforeModifyPasswordRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_PHONE,
		Phone:          suite.Account.phone,
		NationCode:     suite.Account.nationCode,
		SecureQuestions: []*proto.SecureQuestion{
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1New,
				Answer:      suite.Account.answer1New,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey2New,
				Answer:      suite.Account.answer2New,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey3New,
				Answer:      suite.Account.answer3New,
			},
		},
	}
	resp := new(proto.UserValidateSecureQuestionsBeforeModifyPasswordResponse)
	err := suite.JinmuIDService.UserValidateSecureQuestionsBeforeModifyPassword(ctx, req, resp)
	assert.NoError(t, err)
	assert.Equal(t, true, resp.Result)
}

// TestUserValidateSecureQuestionPhoneAndNationCodeError 手机号和区号为空
func (suite *UserValidateSecureQuestionTestSuite) TestUserValidateSecureQuestionPhoneAndNationCodeError() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.UserValidateSecureQuestionsBeforeModifyPasswordRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_PHONE,
		Phone:          suite.Account.phoneIsNull,
		NationCode:     suite.Account.nationCodeIsNull,
		SecureQuestions: []*proto.SecureQuestion{
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1,
				Answer:      suite.Account.answer1,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey2,
				Answer:      suite.Account.answer2,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey3,
				Answer:      suite.Account.answer3,
			},
		},
	}
	resp := new(proto.UserValidateSecureQuestionsBeforeModifyPasswordResponse)
	err := suite.JinmuIDService.UserValidateSecureQuestionsBeforeModifyPassword(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:22000] nonexistent phone"), err)
}

func (suite *UserValidateSecureQuestionTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}

func TestUserValidateSecureQuestionBeforeModifyPasswordTestSuite(t *testing.T) {
	suite.Run(t, new(UserValidateSecureQuestionTestSuite))
}
