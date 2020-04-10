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

// UserSetSecureQuestionsTestSuite 用户设置密保问题测试
type UserSetSecureQuestionsTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

// UserSetSecureQuestionsTestSuite 设置测试环境
func (suite *UserSetSecureQuestionsTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestUserGetSecureQuestionsUserIdIsNull  测试用户设置密保问题UserId为空
func (suite *UserSetSecureQuestionsTestSuite) TestUserGetSecureQuestionsUserIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.UserSetSecureQuestionsRequest{
		UserId: suite.Account.userID,
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
	resp := new(proto.UserSetSecureQuestionsResponse)
	err := suite.JinmuIDService.UserSetSecureQuestions(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:2000] userId is error"), err)
}

// TestUserGetSecureQuestionsIsNull  设置密保问题为空
func (suite *UserSetSecureQuestionsTestSuite) TestUserGetSecureQuestionsIsNull() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	req := &proto.UserSetSecureQuestionsRequest{
		UserId: userID,
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

	resp := new(proto.UserSetSecureQuestionsResponse)
	err = suite.JinmuIDService.UserSetSecureQuestions(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:18000]all questionkey should not be empty"), err)
}

// TestUserGetSecureQuestionsAnswerIsNull  设置密保问题答案为空
func (suite *UserSetSecureQuestionsTestSuite) TestUserGetSecureQuestionsAnswerIsNull() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	req := &proto.UserSetSecureQuestionsRequest{
		UserId: userID,
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
				Answer:      suite.Account.answer3Null,
			},
		},
	}

	resp := new(proto.UserSetSecureQuestionsResponse)
	err = suite.JinmuIDService.UserSetSecureQuestions(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:15000] all answers should not be empty"), err)
}

// TestUserSetSecureQuestions 测试用户设置密保问题
func (suite *UserSetSecureQuestionsTestSuite) TestUserSetSecureQuestions() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	req := &proto.UserSetSecureQuestionsRequest{
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

	resp := new(proto.UserSetSecureQuestionsResponse)
	err = suite.JinmuIDService.UserSetSecureQuestions(ctx, req, resp)
	assert.NoError(t, err)
	assert.Equal(t, true, resp.Result)
}

// TestUserGetSecureQuestionIsExist   密保问题已设置
func (suite *UserSetSecureQuestionsTestSuite) TestUserGetSecureQuestionIsExist() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	req := &proto.UserSetSecureQuestionsRequest{
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

	resp := new(proto.UserSetSecureQuestionsResponse)
	err = suite.JinmuIDService.UserSetSecureQuestions(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:13000] secure questions already exist"), err)
}

// TestUserGetSecureQuestionAnswerIsLong  密保问题答案超长
func (suite *UserSetSecureQuestionsTestSuite) TestUserGetSecureQuestionAnswerIsLong() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	req := &proto.UserSetSecureQuestionsRequest{
		UserId: userID,
		SecureQuestions: []*proto.SecureQuestion{
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1,
				Answer:      suite.Account.answer1Long,
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

	resp := new(proto.UserSetSecureQuestionsResponse)
	err = suite.JinmuIDService.UserSetSecureQuestions(ctx, req, resp)
	assert.Error(t, err)
	assert.Equal(t, false, resp.Result)
}

// TestUserGetSecureQuestionSame 密保问题一致
func (suite *UserSetSecureQuestionsTestSuite) TestUserGetSecureQuestionSame() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	req := &proto.UserSetSecureQuestionsRequest{
		UserId: userID,
		SecureQuestions: []*proto.SecureQuestion{
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1,
				Answer:      suite.Account.answer1,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1,
				Answer:      suite.Account.answer1,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1,
				Answer:      suite.Account.answer1,
			},
		},
	}

	resp := new(proto.UserSetSecureQuestionsResponse)
	err = suite.JinmuIDService.UserSetSecureQuestions(ctx, req, resp)
	assert.Error(t, err)
	assert.Equal(t, false, resp.Result)
}

// TestUserGetSecureQuestionCount 密保问题个数不正确
// To Do
func (suite *UserSetSecureQuestionsTestSuite) TestUserGetSecureQuestionCount() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	req := &proto.UserSetSecureQuestionsRequest{
		UserId: userID,
		SecureQuestions: []*proto.SecureQuestion{
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1,
				Answer:      suite.Account.answer1,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1,
				Answer:      suite.Account.answer1,
			},
		},
	}

	resp := new(proto.UserSetSecureQuestionsResponse)
	err = suite.JinmuIDService.UserSetSecureQuestions(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:12000] wrong count of secure questions"), err)
}

func (suite *UserSetSecureQuestionsTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}

func TestUserSetSecureQuestionsTestSuite(t *testing.T) {
	suite.Run(t, new(UserSetSecureQuestionsTestSuite))
}
