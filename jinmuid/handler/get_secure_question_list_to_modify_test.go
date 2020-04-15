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

// GetSecureQuestionListToModifyTestSuite 获取修改密保前获取已设置密保列表
type GetSecureQuestionListToModifyTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

// GetSecureQuestionListToModifyTestSuite 设置测试环境
func (suite *GetSecureQuestionListToModifyTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestGetSecureQuestionListToModify 测试修改密保前获取已设置密保列表
func (suite *GetSecureQuestionListToModifyTestSuite) TestGetSecureQuestionListToModify() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	req := &proto.GetSecureQuestionListToModifyRequest{
		UserId: userID,
	}
	resp := new(proto.GetSecureQuestionListToModifyResponse)
	err = suite.JinmuIDService.GetSecureQuestionListToModify(ctx, req, resp)
	assert.NoError(t, err)
}

// TestGetSecureQuestionListToModifyQuestionsIsNull 未设置密保问题
func (suite *GetSecureQuestionListToModifyTestSuite) TestGetSecureQuestionListToModifyQuestionsIsNull() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	req := &proto.GetSecureQuestionListToModifyRequest{
		UserId: userID,
	}
	resp := new(proto.GetSecureQuestionListToModifyResponse)
	err = suite.JinmuIDService.GetSecureQuestionListToModify(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:49000] nonexistent secure questions"), err)
}

func (suite *GetSecureQuestionListToModifyTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}
func TestGetSecureQuestionListToModifyTestSuite(t *testing.T) {
	suite.Run(t, new(GetSecureQuestionListToModifyTestSuite))
}
