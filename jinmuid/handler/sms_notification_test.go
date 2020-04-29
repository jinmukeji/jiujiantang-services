package handler

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	generalpb "github.com/jinmukeji/proto/v3/gen/micro/idl/ptypes/v2"
	jinmuidpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SmsNotificationTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

// SetupSuite 设置测试环境
func (suite *SmsNotificationTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestSmsNotificationSignUp 测试注册短信通知
func (suite *SmsNotificationTestSuite) TestSmsNotificationSignUp() {
	t := suite.T()
	ctx := context.Background()

	req := new(jinmuidpb.SmsNotificationRequest)
	req.Phone = suite.Account.phone
	req.Action = jinmuidpb.TemplateAction_TEMPLATE_ACTION_SIGN_UP
	req.Language = generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE
	req.NationCode = suite.Account.nationCode
	resp := new(jinmuidpb.SmsNotificationResponse)
	err := suite.JinmuIDService.SmsNotification(ctx, req, resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.SerialNumber)
}

// TestSmsNotificationSignIn 测试登录短信通知
func (suite *SmsNotificationTestSuite) TestSmsNotificationSignIn() {
	t := suite.T()
	ctx := context.Background()

	req := new(jinmuidpb.SmsNotificationRequest)
	req.Phone = suite.Account.phone
	req.Action = jinmuidpb.TemplateAction_TEMPLATE_ACTION_SIGN_IN
	req.Language = generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE
	req.NationCode = suite.Account.nationCode
	resp := new(jinmuidpb.SmsNotificationResponse)
	err := suite.JinmuIDService.SmsNotification(ctx, req, resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.SerialNumber)
}

// TestSmsNotificationResetPassword  测试重置密码短信通知
func (suite *SmsNotificationTestSuite) TestSmsNotificationResetPassword() {
	t := suite.T()
	ctx := context.Background()
	resp := new(jinmuidpb.SmsNotificationResponse)
	req := new(jinmuidpb.SmsNotificationRequest)
	req.Phone = suite.Account.phone
	req.Action = jinmuidpb.TemplateAction_TEMPLATE_ACTION_RESET_PASSWORD
	req.Language = generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE
	req.NationCode = suite.Account.nationCode
	err := suite.JinmuIDService.SmsNotification(ctx, req, resp)
	assert.NoError(t, err)
}

// TestSmsNotificationModifyPhoneNumber  测试修改手机号码通知
func (suite *SmsNotificationTestSuite) TestSmsNotificationModifyPhoneNumber() {
	t := suite.T()
	ctx := context.Background()

	req := new(jinmuidpb.SmsNotificationRequest)
	req.Phone = suite.Account.phone
	req.Action = jinmuidpb.TemplateAction_TEMPLATE_ACTION_MODIFY_PHONE_NUMBER
	req.Language = generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE
	req.NationCode = suite.Account.nationCode
	req.SendToNewIfModify = false
	resp := new(jinmuidpb.SmsNotificationResponse)
	err := suite.JinmuIDService.SmsNotification(ctx, req, resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.SerialNumber)
}

// TestSmsNotificationSetPhoneNubmer 测试设置手机号短信通知
func (suite *SmsNotificationTestSuite) TestSmsNotificationSetPhoneNubmer() {
	t := suite.T()
	ctx := context.Background()
	resp := new(jinmuidpb.SmsNotificationResponse)
	req := new(jinmuidpb.SmsNotificationRequest)
	req.Phone = suite.Account.phone
	req.Action = jinmuidpb.TemplateAction_TEMPLATE_ACTION_SET_PHONE_NUMBER
	req.Language = generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE
	req.NationCode = suite.Account.nationCode
	err := suite.JinmuIDService.SmsNotification(ctx, req, resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.SerialNumber)
}

// TestSmsNotificationPhoneIsNull  手机号为空
func (suite *SmsNotificationTestSuite) TestSmsNotificationPhoneIsNull() {
	t := suite.T()
	ctx := context.Background()
	resp := new(jinmuidpb.SmsNotificationResponse)
	req := new(jinmuidpb.SmsNotificationRequest)
	req.Phone = suite.Account.phoneIsNull
	req.Action = jinmuidpb.TemplateAction_TEMPLATE_ACTION_SIGN_IN
	req.Language = generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE
	req.NationCode = suite.Account.nationCode
	err := suite.JinmuIDService.SmsNotification(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:20000] invalid phone format"), err)
}

// TestSmsNotificationPhoneFormatError  手机格式错误
func (suite *SmsNotificationTestSuite) TestSmsNotificationPhoneFormatError() {
	t := suite.T()
	ctx := context.Background()
	resp := new(jinmuidpb.SmsNotificationResponse)
	req := new(jinmuidpb.SmsNotificationRequest)
	req.Phone = suite.Account.PhoneError
	req.Action = jinmuidpb.TemplateAction_TEMPLATE_ACTION_SIGN_IN
	req.Language = generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE
	req.NationCode = suite.Account.nationCode
	err := suite.JinmuIDService.SmsNotification(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:20000] invalid phone format"), err)
}

// TestSmsNotificationNationCodeIsNull  区号为空
func (suite *SmsNotificationTestSuite) TestSmsNotificationNationCodeIsNull() {
	t := suite.T()
	ctx := context.Background()
	resp := new(jinmuidpb.SmsNotificationResponse)
	req := new(jinmuidpb.SmsNotificationRequest)
	req.Phone = suite.Account.PhoneError
	req.Action = jinmuidpb.TemplateAction_TEMPLATE_ACTION_SIGN_IN
	req.Language = generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE
	req.NationCode = suite.Account.nationCodeIsNull
	err := suite.JinmuIDService.SmsNotification(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:20000] invalid phone format"), err)
}

// TestSmsNotificationSignUpPhoneExist  注册手机号已存在
func (suite *SmsNotificationTestSuite) TestSmsNotificationSignUpPhoneExist() {
	t := suite.T()
	ctx := context.Background()
	resp := new(jinmuidpb.SmsNotificationResponse)
	req := new(jinmuidpb.SmsNotificationRequest)
	req.Phone = suite.Account.phone
	req.Action = jinmuidpb.TemplateAction_TEMPLATE_ACTION_SIGN_UP
	req.Language = generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE
	req.NationCode = suite.Account.nationCode
	err := suite.JinmuIDService.SmsNotification(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:2000] phone number has been registered"), err)
}

// TestSmsNotificationSetPhoneExist 设置手机号已存在
func (suite *SmsNotificationTestSuite) TestSmsNotificationSetPhoneExist() {
	t := suite.T()
	ctx := context.Background()
	resp := new(jinmuidpb.SmsNotificationResponse)
	req := new(jinmuidpb.SmsNotificationRequest)
	req.Phone = suite.Account.phone
	req.Action = jinmuidpb.TemplateAction_TEMPLATE_ACTION_MODIFY_PHONE_NUMBER
	req.Language = generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE
	req.NationCode = suite.Account.nationCode
	err := suite.JinmuIDService.SmsNotification(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:2000] phone number has been registered"), err)
}

// TestSmsNotificationSignInPhoneNotExist  登录手机号不存在
func (suite *SmsNotificationTestSuite) TestSmsNotificationSignInPhoneNotExist() {
	t := suite.T()
	ctx := context.Background()
	resp := new(jinmuidpb.SmsNotificationResponse)
	req := new(jinmuidpb.SmsNotificationRequest)
	req.Phone = suite.Account.phoneNotExist
	req.Action = jinmuidpb.TemplateAction_TEMPLATE_ACTION_SIGN_IN
	req.Language = generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE
	req.NationCode = suite.Account.nationCode
	err := suite.JinmuIDService.SmsNotification(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:22000] phone is nonexistent"), err)
}

// TestSmsNotificationResetPasswordNotExist 重置手机号不存在
func (suite *SmsNotificationTestSuite) TestSmsNotificationResetPasswordNotExist() {
	t := suite.T()
	ctx := context.Background()
	resp := new(jinmuidpb.SmsNotificationResponse)
	req := new(jinmuidpb.SmsNotificationRequest)
	req.Phone = suite.Account.phoneNotExist
	req.Action = jinmuidpb.TemplateAction_TEMPLATE_ACTION_RESET_PASSWORD
	req.Language = generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE
	req.NationCode = suite.Account.nationCode
	err := suite.JinmuIDService.SmsNotification(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:22000] phone is nonexistent"), err)
}

// TestSmsNotificationModifyPhoneNotExist 修改手机号不存在
func (suite *SmsNotificationTestSuite) TestSmsNotificationModifyPhoneNotExist() {
	t := suite.T()
	ctx := context.Background()
	resp := new(jinmuidpb.SmsNotificationResponse)
	req := new(jinmuidpb.SmsNotificationRequest)
	req.Phone = suite.Account.phoneNotExist
	req.Action = jinmuidpb.TemplateAction_TEMPLATE_ACTION_SET_PHONE_NUMBER
	req.Language = generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE
	req.NationCode = suite.Account.nationCode
	err := suite.JinmuIDService.SmsNotification(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:22000] phone is nonexistent"), err)
}

func (suite *SmsNotificationTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}

func TestSmsNotificationTestSuite(t *testing.T) {
	suite.Run(t, new(SmsNotificationTestSuite))
}
