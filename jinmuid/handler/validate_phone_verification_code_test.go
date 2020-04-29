package handler

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	"github.com/stretchr/testify/suite"
)

// ValidatePhoneVerficationCodeTestSuite 验证手机号验证码
type ValidatePhoneVerficationCodeTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

//  ValidatePhoneVerficationCodeTestSuite 设置测试环境
func (suite *ValidatePhoneVerficationCodeTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestValidatePhoneVerificationCodeTestSuite
func (suite *ValidatePhoneVerficationCodeTestSuite) TestValidatePhoneVerificationCodeTestSuite() {
	t := suite.T()
	ctx := context.Background()
	serialNumber := getSignUpSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getSmsVerificationCode(suite.JinmuIDService, *suite.Account)
	req := &proto.ValidatePhoneVerificationCodeRequest{
		Phone:        suite.Account.phone,
		Mvc:          mvc,
		SerialNumber: serialNumber,
		NationCode:   suite.Account.nationCode,
	}
	resp := new(proto.ValidatePhoneVerificationCodeResponse)
	err := suite.JinmuIDService.ValidatePhoneVerificationCode(ctx, req, resp)
	assert.NoError(t, err)
}

// TestValidatePhoneVerificationCodePhoneIsRegistered  手机号被注册
func (suite *ValidatePhoneVerficationCodeTestSuite) TestValidatePhoneVerificationCodePhoneIsRegistered() {
	t := suite.T()
	ctx := context.Background()
	serialNumber := getSignUpSerialNumber(suite.JinmuIDService, *suite.Account)
	mvc := getSmsVerificationCode(suite.JinmuIDService, *suite.Account)
	req := &proto.ValidatePhoneVerificationCodeRequest{
		Phone:        suite.Account.phone,
		Mvc:          mvc,
		SerialNumber: serialNumber,
		NationCode:   suite.Account.nationCode,
	}
	resp := new(proto.ValidatePhoneVerificationCodeResponse)
	err := suite.JinmuIDService.ValidatePhoneVerificationCode(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:2000] phone number has been registered"), err)
}

// TestValidatePhoneInvaliedCodePhone 手机号码未被注册，验证码无效
func (suite *ValidatePhoneVerficationCodeTestSuite) TestValidatePhoneInvaliedCodePhone() {
	t := suite.T()
	ctx := context.Background()
	serialNumber := getSignUpSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	req := &proto.ValidatePhoneVerificationCodeRequest{
		Phone:        suite.Account.phoneNotExist,
		Mvc:          suite.Account.mvcError,
		SerialNumber: serialNumber,
		NationCode:   suite.Account.nationCode,
	}
	resp := new(proto.ValidatePhoneVerificationCodeResponse)
	err := suite.JinmuIDService.ValidatePhoneVerificationCode(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:31000] vc is invalid"), err)
}

//  TestValidatePhoneErrorFormatCodePhone 手机号码未被注册时，验证手机验证码是否正确
func (suite *ValidatePhoneVerficationCodeTestSuite) TestValidatePhoneErrorFormatCodePhone() {
	t := suite.T()
	ctx := context.Background()
	serialNumber := getSignUpSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	req := &proto.ValidatePhoneVerificationCodeRequest{
		Phone:        suite.Account.phoneNotExist,
		Mvc:          suite.Account.mvcFormatError,
		SerialNumber: serialNumber,
		NationCode:   suite.Account.nationCode,
	}
	resp := new(proto.ValidatePhoneVerificationCodeResponse)
	err := suite.JinmuIDService.ValidatePhoneVerificationCode(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:31000] vc is invalid"), err)
}

// TestWrongPhoneNationCode 手机号码区号不正确
func (suite *ValidatePhoneVerficationCodeTestSuite) TestWrongPhoneNationCode() {
	t := suite.T()
	ctx := context.Background()
	serialNumber := getSignUpSerialNumber(suite.JinmuIDService, *suite.Account)
	mvc := getSmsVerificationCode(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	req := &proto.ValidatePhoneVerificationCodeRequest{
		Phone:        suite.Account.phone,
		Mvc:          mvc,
		SerialNumber: serialNumber,
		NationCode:   suite.Account.nationCodeUSA,
	}
	resp := new(proto.ValidatePhoneVerificationCodeResponse)
	err := suite.JinmuIDService.ValidatePhoneVerificationCode(ctx, req, resp)
	assert.Error(t, errors.New("The NationCodePhone is error"), err)
}

// TestErrorFormatePhone 手机号码格式不正确
func (suite *ValidatePhoneVerficationCodeTestSuite) TestErrorFormatePhone() {
	t := suite.T()
	ctx := context.Background()
	serialNumber := getSignUpSerialNumber(suite.JinmuIDService, *suite.Account)
	mvc := getSmsVerificationCode(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	req := &proto.ValidatePhoneVerificationCodeRequest{
		Phone:        suite.Account.phoneFormatError,
		Mvc:          mvc,
		SerialNumber: serialNumber,
		NationCode:   suite.Account.nationCode,
	}
	resp := new(proto.ValidatePhoneVerificationCodeResponse)
	err := suite.JinmuIDService.ValidatePhoneVerificationCode(ctx, req, resp)
	assert.NoError(t, err)
}

//  TestValidatePhoneErrorSerialNumber 手机号码未被注册序列号错误
func (suite *ValidatePhoneVerficationCodeTestSuite) TestValidatePhoneErrorSerialNumber() {
	t := suite.T()
	ctx := context.Background()
	mvc := getSmsVerificationCode(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	req := &proto.ValidatePhoneVerificationCodeRequest{
		Phone:        suite.Account.phoneNotExist,
		Mvc:          mvc,
		SerialNumber: suite.Account.serialNumberIsNull,
		NationCode:   suite.Account.nationCode,
	}
	resp := new(proto.ValidatePhoneVerificationCodeResponse)
	err := suite.JinmuIDService.ValidatePhoneVerificationCode(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:31000] vc is invalid"), err)
}

func (suite *ValidatePhoneVerficationCodeTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}
func TestValidatePhoneVerficationCodeTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatePhoneVerficationCodeTestSuite))
}

// TestValidateHKPhoneVerificationCodeTestSuite 验证香港手机号码是否被注册
func (suite *ValidatePhoneVerficationCodeTestSuite) TestValidateHKPhoneVerificationCodeTestSuite() {
	t := suite.T()
	ctx := context.Background()
	serialNumber := getSignUpSerialNumberHK(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getSmsVerificationCodeHK(suite.JinmuIDService, *suite.Account)
	req := &proto.ValidatePhoneVerificationCodeRequest{
		Phone:        suite.Account.phoneHK,
		Mvc:          mvc,
		SerialNumber: serialNumber,
		NationCode:   suite.Account.nationCodeHK,
	}
	resp := new(proto.ValidatePhoneVerificationCodeResponse)
	err := suite.JinmuIDService.ValidatePhoneVerificationCode(ctx, req, resp)
	assert.NoError(t, err)
}
