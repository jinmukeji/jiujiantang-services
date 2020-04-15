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

type GetNotificationPreferenceTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

// SetupSuite 设置测试环境
func (suite *GetNotificationPreferenceTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestGetNotificationPreferences 获取通知配置首选项
func (suite *GetNotificationPreferenceTestSuite) TestGetNotificationPreferences() {
	t := suite.T()
	ctx := context.Background()
	// 登录
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	req := &proto.GetNotificationPreferencesRequest{
		UserId: userID,
	}
	resp := new(proto.GetNotificationPreferencesResponse)
	err = suite.JinmuIDService.GetNotificationPreferences(ctx, req, resp)
	assert.NoError(t, err)
}

// TestGetNotificationPreferencesUserInvalide 获取通知配置首选项invalid user context
func (suite *GetNotificationPreferenceTestSuite) TestGetNotificationPreferencesUserInvalide() {
	t := suite.T()
	ctx := context.Background()
	resp := new(proto.GetNotificationPreferencesResponse)
	req := &proto.GetNotificationPreferencesRequest{
		UserId: suite.Account.userID,
	}
	err := suite.JinmuIDService.GetNotificationPreferences(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:1600] invalid user context"), err)
}

// TestGetNotificationPreferencesFailedToFindUser  获取通知配置首选项Failed To Find  User
func (suite *GetNotificationPreferenceTestSuite) TestGetNotificationPreferencesFailedToFindUser() {
	t := suite.T()
	ctx := context.Background()
	resp := new(proto.GetNotificationPreferencesResponse)
	req := &proto.GetNotificationPreferencesRequest{
		UserId: suite.Account.userIDNotExist,
	}
	err := suite.JinmuIDService.GetNotificationPreferences(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:1600] invalid user context"), err)
}

func (suite *GetNotificationPreferenceTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}

func TestGetNotificationPreferenceTestSuite(t *testing.T) {
	suite.Run(t, new(GetNotificationPreferenceTestSuite))
}
