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

// UserGetPreferencesTestSuite 获取用户的首选项配置
type UserGetPreferencesTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

// UserGetPreferencesTestSuite 设置测试环境
func (suite *UserGetPreferencesTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestUserGetPreferences 获取用户的首选项配置
func (suite *UserGetPreferencesTestSuite) TestUserGetPreferences() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	req := &proto.GetUserPreferencesRequest{
		UserId: userID,
	}

	resp := new(proto.GetUserPreferencesResponse)
	err = suite.JinmuIDService.GetUserPreferences(ctx, req, resp)
	assert.NoError(t, err)
}

// TestUserGetPreferencesUserIdIsError 获取用户的首选项配置
func (suite *UserGetPreferencesTestSuite) TestUserGetPreferencesUserIdIsError() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.GetUserPreferencesRequest{
		UserId: suite.Account.userID,
	}
	resp := new(proto.GetUserPreferencesResponse)
	err := suite.JinmuIDService.GetUserPreferences(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:1600] invalid user context"), err)
}

func (suite *UserGetPreferencesTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}

func TestUserGetPreferencesTestSuite(t *testing.T) {
	suite.Run(t, new(UserGetPreferencesTestSuite))
}
