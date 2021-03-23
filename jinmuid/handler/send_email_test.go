package handler

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	jinmuidpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	generalpb "github.com/jinmukeji/proto/v3/gen/micro/idl/ptypes/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LoggedInEmailNotificationTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

// SetupSuite 设置测试环境
func (suite *LoggedInEmailNotificationTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestLoggedInEmailNotificationSet  测试已登录状态下的设置邮箱邮件通知
func (suite *LoggedInEmailNotificationTestSuite) TestLoggedInEmailNotificationSet() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(jinmuidpb.LoggedInEmailNotificationRequest)
	req.Email = suite.Account.email
	req.Action = jinmuidpb.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_SET_SECURE_EMAIL
	req.Language = generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE
	req.UserId = userID
	req.SendToNewIfModify = false
	resp := new(jinmuidpb.LoggedInEmailNotificationResponse)
	err = suite.JinmuIDService.LoggedInEmailNotification(ctx, req, resp)
	assert.NoError(t, err)
}

// TestLoggedInEmailNotificationModify 测试已登录状态下修改邮箱
func (suite *LoggedInEmailNotificationTestSuite) TestLoggedInEmailNotificationModify() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(jinmuidpb.LoggedInEmailNotificationRequest)
	req.Email = suite.Account.email
	req.Action = jinmuidpb.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_MODIFY_SECURE_EMAIL
	req.Language = generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE
	req.UserId = userID
	req.SendToNewIfModify = false
	resp := new(jinmuidpb.LoggedInEmailNotificationResponse)
	err = suite.JinmuIDService.LoggedInEmailNotification(ctx, req, resp)
	assert.NoError(t, err)
}

//TestLoggedInEmailNotificationUnset 测试已登录状态下解绑邮箱
func (suite *LoggedInEmailNotificationTestSuite) TestLoggedInEmailNotificationUnset() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(jinmuidpb.LoggedInEmailNotificationRequest)
	req.Email = suite.Account.email
	req.Action = jinmuidpb.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_UNSET_SECURE_EMAIL
	req.Language = generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE
	req.UserId = userID
	req.SendToNewIfModify = false
	resp := new(jinmuidpb.LoggedInEmailNotificationResponse)
	err = suite.JinmuIDService.LoggedInEmailNotification(ctx, req, resp)
	assert.NoError(t, err)
}

// TestNotLoggedInEmailNotificationReset  测试未登录状态下的重置邮箱
func (suite *LoggedInEmailNotificationTestSuite) TestNotLoggedInEmailNotificationReset() {
	t := suite.T()
	ctx := context.Background()

	req := new(jinmuidpb.NotLoggedInEmailNotificationRequest)
	req.Email = suite.Account.email
	req.Action = jinmuidpb.NotLoggedInSemTemplateAction_NOT_LOGGED_IN_SEM_TEMPLATE_ACTION_FIND_RESET_PASSWORD
	req.Language = generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE
	resp := new(jinmuidpb.NotLoggedInEmailNotificationResponse)
	err := suite.JinmuIDService.NotLoggedInEmailNotification(ctx, req, resp)
	assert.NoError(t, err)
}

// TestNotLoggedInEmailNotificationUsername  测试未登录状态下的重置邮箱
func (suite *LoggedInEmailNotificationTestSuite) TestNotLoggedInEmailNotificationUsername() {
	t := suite.T()
	ctx := context.Background()

	req := new(jinmuidpb.NotLoggedInEmailNotificationRequest)
	req.Email = suite.Account.email
	req.Action = jinmuidpb.NotLoggedInSemTemplateAction_NOT_LOGGED_IN_SEM_TEMPLATE_ACTION_FIND_USERNAME
	req.Language = generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE
	resp := new(jinmuidpb.NotLoggedInEmailNotificationResponse)
	err := suite.JinmuIDService.NotLoggedInEmailNotification(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:21000] nonexistent username"), err)
}

// TestNotLoggedInEmailNotificationIn1Minute 测试未登录状态下的邮件通知1分钟内请求多次
func (suite *LoggedInEmailNotificationTestSuite) TestNotLoggedInEmailNotificationIn1Minute() {
	t := suite.T()
	ctx := context.Background()

	req := new(jinmuidpb.NotLoggedInEmailNotificationRequest)
	req.Email = suite.Account.email
	req.Action = jinmuidpb.NotLoggedInSemTemplateAction_NOT_LOGGED_IN_SEM_TEMPLATE_ACTION_FIND_RESET_PASSWORD
	req.Language = generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE
	resp := new(jinmuidpb.NotLoggedInEmailNotificationResponse)
	err := suite.JinmuIDService.NotLoggedInEmailNotification(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:32000] invalid request counts in 1 minute"), err)
	//assert.NoError(t, err)
}

// TestLoggedInEmailNotificationNotSet  测试已登录状态下的设置邮箱邮件通知
func (suite *LoggedInEmailNotificationTestSuite) TestLoggedInEmailNotificationNotSet() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(jinmuidpb.LoggedInEmailNotificationRequest)
	req.Email = suite.Account.email
	req.Action = jinmuidpb.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_SET_SECURE_EMAIL
	req.Language = generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE
	req.UserId = userID
	req.SendToNewIfModify = false
	resp := new(jinmuidpb.LoggedInEmailNotificationResponse)
	err = suite.JinmuIDService.LoggedInEmailNotification(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:25000] secure email has been set"), err)
}

func (suite *LoggedInEmailNotificationTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}

func TestLoggedInEmailNotificationTestSuite(t *testing.T) {
	suite.Run(t, new(LoggedInEmailNotificationTestSuite))
}
