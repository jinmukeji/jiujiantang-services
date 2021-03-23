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

// SigninPhoneTestSuite 登录测试登录手机号
type SigninPhoneTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

// SigninPhoneTestSuite 设置测试环境
func (suite *SigninPhoneTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestUserSetSigninPhone  测试设置登录手机号
func (suite *SigninPhoneTestSuite) TestUserSetSigninPhone() {
	t := suite.T()
	ctx := context.Background()
	// 登录
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	// 发送通知
	serialNumber := getSignInSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getSmsVerificationCode(suite.JinmuIDService, *suite.Account)
	reqUserSetSigninPhone := new(jinmuidpb.UserSetSigninPhoneRequest)
	reqUserSetSigninPhone.Phone = suite.Account.phone
	reqUserSetSigninPhone.SerialNumber = serialNumber
	reqUserSetSigninPhone.UserId = userID
	reqUserSetSigninPhone.Mvc = mvc
	reqUserSetSigninPhone.NationCode = suite.Account.nationCode
	respUserSetSigninPhone := new(jinmuidpb.UserSetSigninPhoneResponse)
	// 设置手机号
	errUserSetSigninPhone := suite.JinmuIDService.UserSetSigninPhone(ctx, reqUserSetSigninPhone, respUserSetSigninPhone)
	assert.Error(t, errors.New("[errcode:2000] phone number has been registered"), errUserSetSigninPhone)
}

// TestVerifyUserSigninPhone 验证登录手机号单元测试
func (suite *SigninPhoneTestSuite) TestVerifyUserSigninPhone() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getSignInSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getSmsVerificationCode(suite.JinmuIDService, *suite.Account)
	reqVerifyUserSigninPhone := new(jinmuidpb.VerifyUserSigninPhoneRequest)
	reqVerifyUserSigninPhone.NationCode = suite.Account.nationCode
	reqVerifyUserSigninPhone.Phone = suite.Account.phone
	reqVerifyUserSigninPhone.SerialNumber = serialNumber
	reqVerifyUserSigninPhone.Mvc = mvc
	respVerifyUserSigninPhone := new(jinmuidpb.VerifyUserSigninPhoneResponse)
	err := suite.JinmuIDService.VerifyUserSigninPhone(ctx, reqVerifyUserSigninPhone, respVerifyUserSigninPhone)
	assert.NoError(t, err)
}

// TestUserModifyPhone 修改手机号单元测试
func (suite *SigninPhoneTestSuite) TestUserModifyPhone() {
	t := suite.T()
	ctx := context.Background()
	// 登录
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	// 发送通知
	serialNumber := getSignInSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getSmsVerificationCode(suite.JinmuIDService, *suite.Account)
	// 验证登录手机号
	reqVerifyUserSigninPhone := new(jinmuidpb.VerifyUserSigninPhoneRequest)
	reqVerifyUserSigninPhone.NationCode = suite.Account.nationCode
	reqVerifyUserSigninPhone.Phone = suite.Account.phone
	reqVerifyUserSigninPhone.SerialNumber = serialNumber
	reqVerifyUserSigninPhone.Mvc = mvc
	respVerifyUserSigninPhone := new(jinmuidpb.VerifyUserSigninPhoneResponse)
	err = suite.JinmuIDService.VerifyUserSigninPhone(ctx, reqVerifyUserSigninPhone, respVerifyUserSigninPhone)
	assert.NoError(t, err)
	// 新手机号发送通知和获取最新验证码（新旧手机号一样）
	resp2 := new(jinmuidpb.SmsNotificationResponse)
	req2 := new(jinmuidpb.SmsNotificationRequest)
	req2.Phone = suite.Account.phone
	req2.Action = jinmuidpb.TemplateAction_TEMPLATE_ACTION_SET_PHONE_NUMBER
	req2.Language = generalpb.Language_LANGUAGE_ENGLISH
	req2.NationCode = suite.Account.nationCode
	err = suite.JinmuIDService.SmsNotification(ctx, req2, resp2)
	assert.NoError(t, err)
	respGetLatestVerificationCodes2 := new(jinmuidpb.GetLatestVerificationCodesResponse)
	reqGetLatestVerificationCodes2 := new(jinmuidpb.GetLatestVerificationCodesRequest)
	reqGetLatestVerificationCodes2.SendTo = []*jinmuidpb.SingleGetLatestVerificationCode{
		&jinmuidpb.SingleGetLatestVerificationCode{
			SendVia:    jinmuidpb.SendVia_SEND_VIA_PHONE_SEND_VIA,
			Phone:      suite.Account.phone,
			NationCode: suite.Account.nationCode,
		},
	}
	err = suite.JinmuIDService.GetLatestVerificationCodes(ctx, reqGetLatestVerificationCodes2, respGetLatestVerificationCodes2)
	assert.NoError(t, err)
	newMvc := respGetLatestVerificationCodes2.LatestVerificationCodes[0].VerificationCode
	// 修改手机号
	reqUserModifyPhone := new(jinmuidpb.UserModifyPhoneRequest)
	reqUserModifyPhone.SerialNumber = resp2.SerialNumber
	reqUserModifyPhone.Mvc = newMvc
	reqUserModifyPhone.NationCode = suite.Account.nationCode
	reqUserModifyPhone.UserId = userID
	reqUserModifyPhone.Phone = suite.Account.phone
	reqUserModifyPhone.VerificationNumber = respVerifyUserSigninPhone.VerificationNumber
	reqUserModifyPhone.OldNationCode = suite.Account.nationCode
	reqUserModifyPhone.OldPhone = suite.Account.phone
	respUserModifyPhone := new(jinmuidpb.UserModifyPhoneResponse)
	err = suite.JinmuIDService.UserModifyPhone(ctx, reqUserModifyPhone, respUserModifyPhone)
	assert.Error(t, errors.New("[errcode:38000] new phone cannot equals old phone"), err)
}

func (suite *SigninPhoneTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}

func TestUserSetSigninPhoneTestSuite(t *testing.T) {
	suite.Run(t, new(SigninPhoneTestSuite))
}
