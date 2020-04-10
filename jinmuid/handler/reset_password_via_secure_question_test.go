package handler

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// UserResetPasswordViaSecureQuestionTestSuite 通过密保问题重置密码测试
type UserResetPasswordViaSecureQuestionTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

// UserResetPasswordViaSecureQuestionTestSuite 设置测试环境
func (suite *UserResetPasswordViaSecureQuestionTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestUserResetPasswordViaSecureQuestionByPhone 测试通过密保问题重置密码
func (suite *UserResetPasswordViaSecureQuestionTestSuite) TestUserResetPasswordViaSecureQuestionByPhone() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.UserResetPasswordViaSecureQuestionsRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_PHONE,
		Phone:          suite.Account.phone,
		NationCode:     suite.Account.nationCode,
		Password:       suite.Account.password,
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
	resp := new(proto.UserResetPasswordViaSecureQuestionsResponse)
	err := suite.JinmuIDService.UserResetPasswordViaSecureQuestions(ctx, req, resp)
	assert.NoError(t, err)
	assert.Equal(t, true, resp.Result)
}

// TestUserResetPasswordViaSecureQuestionByUsername 测试通过密保问题重置密码
func (suite *UserResetPasswordViaSecureQuestionTestSuite) TestUserResetPasswordViaSecureQuestionByUsername() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.UserResetPasswordViaSecureQuestionsRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_USERNAME,
		Username:       suite.Account.username,
		Password:       suite.Account.password,
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
	resp := new(proto.UserResetPasswordViaSecureQuestionsResponse)
	err := suite.JinmuIDService.UserResetPasswordViaSecureQuestions(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:49000] secure questions are not set"), err)
}

//TestUserResetPasswordViaSecureQuestionTypeIsNull   type为空
func (suite *UserResetPasswordViaSecureQuestionTestSuite) TestUserResetPasswordViaSecureQuestionTypeIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.UserResetPasswordViaSecureQuestionsRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_UNKNOWN,
		Phone:          suite.Account.phone,
		NationCode:     suite.Account.nationCode,
		Password:       suite.Account.password,
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
	resp := new(proto.UserResetPasswordViaSecureQuestionsResponse)
	err := suite.JinmuIDService.UserResetPasswordViaSecureQuestions(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:11000] invalid secure queston validation type"), err)
}

// TestUserResetPasswordViaSecureQuestionPhoneError   手机格式不正确
func (suite *UserResetPasswordViaSecureQuestionTestSuite) TestUserResetPasswordViaSecureQuestionPhoneError() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.UserResetPasswordViaSecureQuestionsRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_PHONE,
		Phone:          suite.Account.PhoneError,
		NationCode:     suite.Account.nationCode,
		Password:       suite.Account.password,
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
	resp := new(proto.UserResetPasswordViaSecureQuestionsResponse)
	err := suite.JinmuIDService.UserResetPasswordViaSecureQuestions(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:20000] wrong format of phone"), err)
}

// TestUserResetPasswordViaSecureQuestionPhoneNotExist  手机号不存在
func (suite *UserResetPasswordViaSecureQuestionTestSuite) TestUserResetPasswordViaSecureQuestionPhoneNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.UserResetPasswordViaSecureQuestionsRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_PHONE,
		Phone:          suite.Account.phoneNotExist,
		NationCode:     suite.Account.nationCode,
		Password:       suite.Account.password,
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
	resp := new(proto.UserResetPasswordViaSecureQuestionsResponse)
	err := suite.JinmuIDService.UserResetPasswordViaSecureQuestions(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:22000] nonexistent phone"), err)
}

// TestUserResetPasswordViaSecureQuestionUsernameNotExist  用户名不存在
func (suite *UserResetPasswordViaSecureQuestionTestSuite) TestUserResetPasswordViaSecureQuestionUsernameNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.UserResetPasswordViaSecureQuestionsRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_USERNAME,
		Username:       suite.Account.usernameNull,
		Password:       suite.Account.password,
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
	resp := new(proto.UserResetPasswordViaSecureQuestionsResponse)
	err := suite.JinmuIDService.UserResetPasswordViaSecureQuestions(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:21000] nonexistent username"), err)
}

// TestUserResetPasswordViaSecureQuestionAnswerError  答案错误
func (suite *UserResetPasswordViaSecureQuestionTestSuite) TestUserResetPasswordViaSecureQuestionAnswerError() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.UserResetPasswordViaSecureQuestionsRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_PHONE,
		Phone:          suite.Account.phone,
		NationCode:     suite.Account.nationCode,
		Password:       suite.Account.password,
		SecureQuestions: []*proto.SecureQuestion{
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1,
				Answer:      suite.Account.answer1New,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey2,
				Answer:      suite.Account.answer2New,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey3,
				Answer:      suite.Account.answer3New,
			},
		},
	}
	resp := new(proto.UserResetPasswordViaSecureQuestionsResponse)
	err := suite.JinmuIDService.UserResetPasswordViaSecureQuestions(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:17000] mismatch secure question"), err)
}

// TestUserResetPasswordViaSecureQuestionPasswordSame 新旧密码一致

func (suite *UserResetPasswordViaSecureQuestionTestSuite) TestUserResetPasswordViaSecureQuestionPasswordSame() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.UserResetPasswordViaSecureQuestionsRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_PHONE,
		Phone:          suite.Account.phone,
		NationCode:     suite.Account.nationCode,
		Password:       suite.Account.password,
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
	resp := new(proto.UserResetPasswordViaSecureQuestionsResponse)
	err := suite.JinmuIDService.UserResetPasswordViaSecureQuestions(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:16000] new password cannot equals old password"), err)
}

func (suite *UserResetPasswordViaSecureQuestionTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}

func TestUserResetPasswordViaSecureQuestionTestSuite(t *testing.T) {
	suite.Run(t, new(UserResetPasswordViaSecureQuestionTestSuite))
}
