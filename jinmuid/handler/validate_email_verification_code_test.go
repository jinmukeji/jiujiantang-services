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

// ValidateEmailVerificationCodeTestSuite 验证邮箱验证码是否正确测试
type ValidateEmailVerificationCodeTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

// ValidateEmailVerificationCodeTestSuite 设置测试环境
func (suite *ValidateEmailVerificationCodeTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestValidateEmailVerificationCode 测试验证邮箱验证码是否正确
func (suite *ValidateEmailVerificationCodeTestSuite) TestValidateEmailVerificationCode() {
	t := suite.T()
	ctx := context.Background()
	serialNumber := getModifyEmailSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	resp := new(proto.ValidateEmailVerificationCodeResponse)
	req := new(proto.ValidateEmailVerificationCodeRequest)
	req.VerificationCode = mvc
	req.SerialNumber = serialNumber
	req.Email = suite.Account.email
	req.VerificationType = suite.Account.verificationType
	err := suite.JinmuIDService.ValidateEmailVerificationCode(ctx, req, resp)
	assert.NoError(t, err)
}

// TestValidateEmailVerificationCodeVCIsNull  验证码为空
func (suite *ValidateEmailVerificationCodeTestSuite) TestValidateEmailVerificationCodeVCIsNull() {
	t := suite.T()
	ctx := context.Background()
	serialNumber := getSignInSerialNumber(suite.JinmuIDService, *suite.Account)
	resp := new(proto.ValidateEmailVerificationCodeResponse)
	req := new(proto.ValidateEmailVerificationCodeRequest)
	req.VerificationCode = suite.Account.mvcIsNull
	req.SerialNumber = serialNumber
	req.Email = suite.Account.email
	req.VerificationType = suite.Account.verificationType
	err := suite.JinmuIDService.ValidateEmailVerificationCode(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:30000] expired vc record"), err)
}

// TestValidateEmailVerificationCodeSerialNumberIsNull  serialnumber为空
func (suite *ValidateEmailVerificationCodeTestSuite) TestValidateEmailVerificationCodeSerialNumberIsNull() {
	t := suite.T()
	ctx := context.Background()
	// 获取最新验证码
	mvc := getEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	resp := new(proto.ValidateEmailVerificationCodeResponse)
	req := new(proto.ValidateEmailVerificationCodeRequest)
	req.VerificationCode = mvc
	req.SerialNumber = suite.Account.serialNumberIsNull
	req.Email = suite.Account.email
	req.VerificationType = suite.Account.verificationType
	err := suite.JinmuIDService.ValidateEmailVerificationCode(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:30000] expired vc record"), err)
}

// TestValidateEmailVerificationCodeEmailIsNull
func (suite *ValidateEmailVerificationCodeTestSuite) TestValidateEmailVerificationCodeEmailIsNull() {
	t := suite.T()
	ctx := context.Background()
	serialNumber := getSignInSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	resp := new(proto.ValidateEmailVerificationCodeResponse)
	req := new(proto.ValidateEmailVerificationCodeRequest)
	req.VerificationCode = mvc
	req.SerialNumber = serialNumber
	req.Email = suite.Account.emailNull
	req.VerificationType = suite.Account.verificationType
	err := suite.JinmuIDService.ValidateEmailVerificationCode(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:41000] invalid email"), err)
}

// TestValidateEmailVerificationCodeEmailError
func (suite *ValidateEmailVerificationCodeTestSuite) TestValidateEmailVerificationCodeEmailError() {
	t := suite.T()
	ctx := context.Background()
	serialNumber := getSignInSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	resp := new(proto.ValidateEmailVerificationCodeResponse)
	req := new(proto.ValidateEmailVerificationCodeRequest)
	req.VerificationCode = mvc
	req.SerialNumber = serialNumber
	req.Email = suite.Account.emailError
	req.VerificationType = suite.Account.verificationType
	err := suite.JinmuIDService.ValidateEmailVerificationCode(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:41000] invalid email"), err)
}

//TestValidateEmailVerificationCodeEmailNotExist
func (suite *ValidateEmailVerificationCodeTestSuite) TestValidateEmailVerificationCodeEmailNotExist() {
	t := suite.T()
	ctx := context.Background()
	serialNumber := getSignInSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	resp := new(proto.ValidateEmailVerificationCodeResponse)
	req := new(proto.ValidateEmailVerificationCodeRequest)
	req.VerificationCode = mvc
	req.SerialNumber = serialNumber
	req.Email = suite.Account.emailNoExist
	req.VerificationType = suite.Account.verificationType
	err := suite.JinmuIDService.ValidateEmailVerificationCode(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:41000] invalid email"), err)
}

func (suite *ValidateEmailVerificationCodeTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}

func TestValidateEmailVerificationCodeTestSuite(t *testing.T) {
	suite.Run(t, new(ValidateEmailVerificationCodeTestSuite))
}
