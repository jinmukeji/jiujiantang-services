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

// UserGetSecureQuestionListTestSuite 获取密保问题列表测试
type UserGetSecureQuestionListTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

// UserGetSecureQuestionListTestSuite 设置测试环境
func (suite *UserGetSecureQuestionListTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestUserGetSecureQuestionList 测试获取密保问题列表
func (suite *UserGetSecureQuestionListTestSuite) TestUserGetSecureQuestionList() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	req := &proto.UserGetSecureQuestionListRequest{
		UserId: userID,
	}
	resp := new(proto.UserGetSecureQuestionListResponse)
	err = suite.JinmuIDService.UserGetSecureQuestionList(ctx, req, resp)
	assert.NoError(t, err)
}

// TestUserGetSecureQuestionListIsError 未设置密保问题
func (suite *UserGetSecureQuestionListTestSuite) TestUserGetSecureQuestionListIsError() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	req := &proto.UserGetSecureQuestionListRequest{
		UserId: userID,
	}
	resp := new(proto.UserGetSecureQuestionListResponse)
	err = suite.JinmuIDService.UserGetSecureQuestionList(ctx, req, resp)
	assert.Error(t, errors.New("[[errcode:49000] secure questions are not set"), err)
}

// TestUserGetSecureQuestionListUserIDError 未设置密保问题UserID不存在
func (suite *UserGetSecureQuestionListTestSuite) TestUserGetSecureQuestionListUserIDError() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.UserGetSecureQuestionListRequest{
		UserId: suite.Account.userID,
	}
	resp := new(proto.UserGetSecureQuestionListResponse)
	err := suite.JinmuIDService.UserGetSecureQuestionList(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:1600] Invalid user 0"), err)
}

func (suite *UserGetSecureQuestionListTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}

func TestUserGetSecureQuestionListTestSuite(t *testing.T) {
	suite.Run(t, new(UserGetSecureQuestionListTestSuite))
}
