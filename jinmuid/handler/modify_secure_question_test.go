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

// UserModifySecureQuestionsTestSuite 修改密保测试
type UserModifySecureQuestionsTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

// UserModifySecureQuestionsTestSuite 设置测试环境
func (suite *UserModifySecureQuestionsTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestUserModifySecureQuestions 测试修改密保
func (suite *UserModifySecureQuestionsTestSuite) TestUserModifySecureQuestions() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	req := &proto.UserModifySecureQuestionsRequest{
		UserId: userID,
		OldSecureQuestions: []*proto.SecureQuestion{
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
		NewSecureQuestions: []*proto.SecureQuestion{
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

	resp := new(proto.UserModifySecureQuestionsResponse)
	err = suite.JinmuIDService.UserModifySecureQuestions(ctx, req, resp)
	assert.NoError(t, err)
	assert.Equal(t, true, resp.Result)
}

// TestUserModifySecureQuestionsKeyIsNull  questionKey为空
func (suite *UserModifySecureQuestionsTestSuite) TestUserModifySecureQuestionsKeyIsNull() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	req := &proto.UserModifySecureQuestionsRequest{
		UserId: userID,
		OldSecureQuestions: []*proto.SecureQuestion{
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1Null,
				Answer:      suite.Account.answer1New,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey2Null,
				Answer:      suite.Account.answer2New,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey3Null,
				Answer:      suite.Account.answer3New,
			},
		},
		NewSecureQuestions: []*proto.SecureQuestion{
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

	resp := new(proto.UserModifySecureQuestionsResponse)
	err = suite.JinmuIDService.UserModifySecureQuestions(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:17000] mismatch secure question"), err)
}

// TestUserModifySecureQuestionsAnswerIsNull   Answer为空
func (suite *UserModifySecureQuestionsTestSuite) TestUserModifySecureQuestionsAnswerIsNull() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	req := &proto.UserModifySecureQuestionsRequest{
		UserId: userID,
		OldSecureQuestions: []*proto.SecureQuestion{
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1New,
				Answer:      suite.Account.answer1Null,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey2New,
				Answer:      suite.Account.answer2Null,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey3New,
				Answer:      suite.Account.answer3Null,
			},
		},
		NewSecureQuestions: []*proto.SecureQuestion{
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1New,
				Answer:      suite.Account.answer1Null,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey2New,
				Answer:      suite.Account.answer2Null,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey3New,
				Answer:      suite.Account.answer3Null,
			},
		},
	}

	resp := new(proto.UserModifySecureQuestionsResponse)
	err = suite.JinmuIDService.UserModifySecureQuestions(ctx, req, resp)
	assert.NoError(t, err)
	assert.Equal(t, false, resp.Result)
}

// TestUserModifySecureQuestionsAnswerIsLong   答案超出长度
func (suite *UserModifySecureQuestionsTestSuite) TestUserModifySecureQuestionsAnswerIsLong() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	req := &proto.UserModifySecureQuestionsRequest{
		UserId: userID,
		OldSecureQuestions: []*proto.SecureQuestion{
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1New,
				Answer:      suite.Account.answer1,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1New,
				Answer:      suite.Account.answer2,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey3New,
				Answer:      suite.Account.answer3,
			},
		},
		NewSecureQuestions: []*proto.SecureQuestion{
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1,
				Answer:      suite.Account.answer1Long,
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

	resp := new(proto.UserModifySecureQuestionsResponse)
	err = suite.JinmuIDService.UserModifySecureQuestions(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:19000] repeated question"), err)
}

// TestUserModifySecureQuestionsSame   新旧密保问题一致
func (suite *UserModifySecureQuestionsTestSuite) TestUserModifySecureQuestionsSame() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	req := &proto.UserModifySecureQuestionsRequest{
		UserId: userID,
		OldSecureQuestions: []*proto.SecureQuestion{
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
		NewSecureQuestions: []*proto.SecureQuestion{
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1Same,
				Answer:      suite.Account.answer1,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey2Same,
				Answer:      suite.Account.answer2,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey3Same,
				Answer:      suite.Account.answer3,
			},
		},
	}

	resp := new(proto.UserModifySecureQuestionsResponse)
	err = suite.JinmuIDService.UserModifySecureQuestions(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:17000] mismatch secure question"), err)
}

// TestUserModifySecureQuestionsFormat    客案有特殊字符
func (suite *UserModifySecureQuestionsTestSuite) TestUserModifySecureQuestionsFormat() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	req := &proto.UserModifySecureQuestionsRequest{
		UserId: userID,
		OldSecureQuestions: []*proto.SecureQuestion{
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1New,
				Answer:      suite.Account.answer1,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey2New,
				Answer:      suite.Account.answer2,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey3New,
				Answer:      suite.Account.answer3,
			},
		},
		NewSecureQuestions: []*proto.SecureQuestion{
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1,
				Answer:      suite.Account.answer1Format,
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

	resp := new(proto.UserModifySecureQuestionsResponse)
	err = suite.JinmuIDService.UserModifySecureQuestions(ctx, req, resp)
	assert.NoError(t, err)
	assert.Equal(t, false, resp.Result)
}

// TestUserModifySecureQuestionsNotExist  问题不存在
func (suite *UserModifySecureQuestionsTestSuite) TestUserModifySecureQuestionsNotExist() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	req := &proto.UserModifySecureQuestionsRequest{
		UserId: userID,
		OldSecureQuestions: []*proto.SecureQuestion{
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
		NewSecureQuestions: []*proto.SecureQuestion{
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1New,
				Answer:      suite.Account.answer1New,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey2New,
				Answer:      suite.Account.answer2New,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1NoExist,
				Answer:      suite.Account.answer3New,
			},
		},
	}

	resp := new(proto.UserModifySecureQuestionsResponse)
	err = suite.JinmuIDService.UserModifySecureQuestions(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:18000] wrong format of question"), err)
}

// TestUserModifySecureQuestionsCount  问题 个数不对
func (suite *UserModifySecureQuestionsTestSuite) TestUserModifySecureQuestionsCount() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	req := &proto.UserModifySecureQuestionsRequest{
		UserId: userID,
		OldSecureQuestions: []*proto.SecureQuestion{
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1,
				Answer:      suite.Account.answer1,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey2,
				Answer:      suite.Account.answer2,
			},
		},
		NewSecureQuestions: []*proto.SecureQuestion{
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey1New,
				Answer:      suite.Account.answer1New,
			},
			&proto.SecureQuestion{
				QuestionKey: suite.Account.questionKey2New,
				Answer:      suite.Account.answer2New,
			},
		},
	}

	resp := new(proto.UserModifySecureQuestionsResponse)
	err = suite.JinmuIDService.UserModifySecureQuestions(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:12000] wrong count of secure questions"), err)
}

func (suite *UserModifySecureQuestionsTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}

func TestUserModifySecureQuestionsTestSuite(t *testing.T) {
	suite.Run(t, new(UserModifySecureQuestionsTestSuite))
}
