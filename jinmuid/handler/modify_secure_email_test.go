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

// ModifySecureEmailTestSuite 修改安全邮箱测试
type ModifySecureEmailTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

// ModifySecureEmailTestSuite 设置测试环境
func (suite *ModifySecureEmailTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
	// 绑定邮箱
	t := suite.T()
	ctx := context.Background()
	serialNumber := getSetEmailSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(proto.SetSecureEmailRequest)
	req.UserId = userID
	req.VerificationCode = mvc
	req.SerialNumber = serialNumber
	req.Email = suite.Account.email
	resp := new(proto.SetSecureEmailResponse)
	err = suite.JinmuIDService.SetSecureEmail(ctx, req, resp)
	assert.NoError(t, err)
}

// TesModifySecureEmailNewVCIsNull  新验证码为空
func (suite *ModifySecureEmailTestSuite) TestModifySecureEmailNewVCIsNull() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getModifyEmailSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(proto.ModifySecureEmailRequest)
	req.UserId = userID
	req.NewVerificationCode = suite.Account.mvcIsNull
	req.NewSerialNumber = serialNumber
	req.NewEmail = suite.Account.emailNull
	req.OldEmail = suite.Account.email
	req.OldVerificationNumber = mvc
	resp := new(proto.ModifySecureEmailResponse)
	err = suite.JinmuIDService.ModifySecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:45000] old secure email does not exist"), err)
}

// TesModifySecureEmailNewSerialNumberIsNull  serialNumber为空
func (suite *ModifySecureEmailTestSuite) TestModifySecureEmailNewSerialNumberIsNull() {
	t := suite.T()
	ctx := context.Background()
	// 获取最新验证码
	mvc := getEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	mvcNew := getNewEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(proto.ModifySecureEmailRequest)
	req.UserId = userID
	req.NewVerificationCode = mvcNew
	req.NewSerialNumber = suite.Account.serialNumberIsNull
	req.NewEmail = suite.Account.emailNew
	req.OldEmail = suite.Account.email
	req.OldVerificationNumber = mvc
	resp := new(proto.ModifySecureEmailResponse)
	err = suite.JinmuIDService.ModifySecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:45000] old secure email does not exist"), err)
}

// TesModifySecureEmailNewEmailIsnull  新邮箱为空
func (suite *ModifySecureEmailTestSuite) TestModifySecureEmailNewEmailIsnull() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getModifyEmailSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	mvcNew := getNewEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(proto.ModifySecureEmailRequest)
	req.UserId = userID
	req.NewVerificationCode = mvcNew
	req.NewSerialNumber = serialNumber
	req.NewEmail = suite.Account.emailNull
	req.OldEmail = suite.Account.email
	req.OldVerificationNumber = mvc
	resp := new(proto.ModifySecureEmailResponse)
	err = suite.JinmuIDService.ModifySecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:45000] old secure email does not exist"), err)
}

// TesModifySecureEmailNewEmailFormatError  新邮箱格式错误
func (suite *ModifySecureEmailTestSuite) TestModifySecureEmailNewEmailFormatError() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getModifyEmailSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	mvcNew := getNewEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(proto.ModifySecureEmailRequest)
	req.UserId = userID
	req.NewVerificationCode = mvcNew
	req.NewSerialNumber = serialNumber
	req.NewEmail = suite.Account.emailError
	req.OldEmail = suite.Account.email
	req.OldVerificationNumber = mvc
	resp := new(proto.ModifySecureEmailResponse)
	err = suite.JinmuIDService.ModifySecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:45000] old secure email does not exist"), err)
}

// TesModifySecureEmailNewEmailIsExist 新邮箱已被其他用户绑定
func (suite *ModifySecureEmailTestSuite) TestModifySecureEmailNewEmailIsExist() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getModifyEmailSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	mvcNew := getNewEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(proto.ModifySecureEmailRequest)
	req.UserId = userID
	req.NewVerificationCode = mvcNew
	req.NewSerialNumber = serialNumber
	req.NewEmail = suite.Account.email
	req.OldEmail = suite.Account.email
	req.OldVerificationNumber = mvc
	resp := new(proto.ModifySecureEmailResponse)
	err = suite.JinmuIDService.ModifySecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:45000] old secure email does not exist"), err)
}

// TesModifySecureEmailOldEmailIsNull  旧邮箱为空
func (suite *ModifySecureEmailTestSuite) TestModifySecureEmailOldEmailIsNull() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getModifyEmailSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	mvcNew := getNewEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(proto.ModifySecureEmailRequest)
	req.UserId = userID
	req.NewVerificationCode = mvcNew
	req.NewSerialNumber = serialNumber
	req.NewEmail = suite.Account.emailNew
	req.OldEmail = suite.Account.emailNull
	req.OldVerificationNumber = mvc
	resp := new(proto.ModifySecureEmailResponse)
	err = suite.JinmuIDService.ModifySecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:45000] old secure email does not exist"), err)
}

// TesModifySecureEmailOldEmailFormatError  旧邮箱格式错误
func (suite *ModifySecureEmailTestSuite) TestModifySecureEmailOldEmailFormatError() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getModifyEmailSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	mvcNew := getNewEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(proto.ModifySecureEmailRequest)
	req.UserId = userID
	req.NewVerificationCode = mvcNew
	req.NewSerialNumber = serialNumber
	req.NewEmail = suite.Account.emailNew
	req.OldEmail = suite.Account.emailError
	req.OldVerificationNumber = mvc
	resp := new(proto.ModifySecureEmailResponse)
	err = suite.JinmuIDService.ModifySecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:45000] old secure email does not exist"), err)
}

// TesModifySecureEmailOldEmailIsExist 旧邮箱未被其他用户绑定
func (suite *ModifySecureEmailTestSuite) TestModifySecureEmailOldEmailIsExist() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getModifyEmailSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	mvcNew := getNewEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(proto.ModifySecureEmailRequest)
	req.UserId = userID
	req.NewVerificationCode = mvcNew
	req.NewSerialNumber = serialNumber
	req.NewEmail = suite.Account.emailNew
	req.OldEmail = suite.Account.emailNoExist
	req.OldVerificationNumber = mvc
	resp := new(proto.ModifySecureEmailResponse)
	err = suite.JinmuIDService.ModifySecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:45000] old secure email does not exist"), err)
}

// TesModifySecureEmail  测试修改安全邮箱
func (suite *ModifySecureEmailTestSuite) TestModifySecureEmail() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getModifyEmailSerialNumberNewEmail(suite.JinmuIDService, *suite.Account)
	suite.T().Log("serialNumber:", serialNumber)
	// 获取最新验证码
	verificationNumber := getVerificationNumber(suite.JinmuIDService, *suite.Account)
	suite.T().Log("verificationNumber:", verificationNumber)
	mvcNew := getNewEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	suite.T().Log("mvcNew:", mvcNew)
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(proto.ModifySecureEmailRequest)
	req.UserId = userID
	req.NewVerificationCode = mvcNew
	req.NewSerialNumber = serialNumber
	req.NewEmail = suite.Account.emailNew
	req.OldEmail = suite.Account.email
	req.OldVerificationNumber = verificationNumber
	resp := new(proto.ModifySecureEmailResponse)
	err = suite.JinmuIDService.ModifySecureEmail(ctx, req, resp)
	assert.NoError(t, err)
}

//TestModifySecureEmailOldIsNotExist旧邮箱不存在
func (suite *ModifySecureEmailTestSuite) TestModifySecureEmailOldIsNotExist() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getModifyEmailSerialNumberNewEmail(suite.JinmuIDService, *suite.Account)
	suite.T().Log("serialNumber:", serialNumber)
	// 获取最新验证码
	verificationNumber := getVerificationNumber(suite.JinmuIDService, *suite.Account)
	suite.T().Log("verificationNumber:", verificationNumber)
	mvcNew := getNewEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	suite.T().Log("mvcNew:", mvcNew)
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(proto.ModifySecureEmailRequest)
	req.UserId = userID
	req.NewVerificationCode = mvcNew
	req.NewSerialNumber = serialNumber
	req.NewEmail = suite.Account.emailNew
	req.OldEmail = suite.Account.email
	req.OldVerificationNumber = verificationNumber
	resp := new(proto.ModifySecureEmailResponse)
	err = suite.JinmuIDService.ModifySecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:45000] old secure email does not exist"), err)
}

func (suite *ModifySecureEmailTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}

func TestModifySecureEmailTestSuite(t *testing.T) {
	suite.Run(t, new(ModifySecureEmailTestSuite))
}
