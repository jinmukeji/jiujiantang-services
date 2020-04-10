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

type SetSecureEmailTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

// SetupSuite 设置测试环境
func (suite *SetSecureEmailTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestSetSecureEmail 测试设置安全邮箱
func (suite *SetSecureEmailTestSuite) TestSetSecureEmail() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
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

// TestSetSecureEmailUserIdIsNull  设置安全邮箱userid为空
func (suite *SetSecureEmailTestSuite) TestSetSecureEmailUserIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getSetEmailSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getEmailVerificationCode(suite.JinmuIDService, *suite.Account)

	req := new(proto.SetSecureEmailRequest)
	req.UserId = suite.Account.userID
	req.VerificationCode = mvc
	req.SerialNumber = serialNumber
	req.Email = suite.Account.email
	resp := new(proto.SetSecureEmailResponse)
	err := suite.JinmuIDService.SetSecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:25000] secure email has been set"), err)
}

// TestSetSecureEmailVCIsNull     verificationCode为空
func (suite *SetSecureEmailTestSuite) TestSetSecureEmailVCIsNull() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getSetEmailSerialNumber(suite.JinmuIDService, *suite.Account)
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(proto.SetSecureEmailRequest)
	req.UserId = userID
	req.VerificationCode = suite.Account.mvcIsNull
	req.SerialNumber = serialNumber
	req.Email = suite.Account.email
	resp := new(proto.SetSecureEmailResponse)
	err = suite.JinmuIDService.SetSecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:31000] invalid vc record"), err)
}

// TestSetSecureEmailVCError     verificationCode错误,过期
func (suite *SetSecureEmailTestSuite) TestSetSecureEmailVCError() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getSetEmailSerialNumber(suite.JinmuIDService, *suite.Account)
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(proto.SetSecureEmailRequest)
	req.UserId = userID
	req.VerificationCode = suite.Account.mvcError
	req.SerialNumber = serialNumber
	req.Email = suite.Account.email
	resp := new(proto.SetSecureEmailResponse)
	err = suite.JinmuIDService.SetSecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:31000] invalid vc record"), err)
}

// TestSetSecureEmailSerialNumberIsNull    serialnumber为空
func (suite *SetSecureEmailTestSuite) TestSetSecureEmailSerialNumberIsNull() {
	t := suite.T()
	ctx := context.Background()
	// 获取最新验证码
	mvc := getEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(proto.SetSecureEmailRequest)
	req.UserId = userID
	req.VerificationCode = mvc
	req.SerialNumber = suite.Account.serialNumberIsNull
	req.Email = suite.Account.email
	resp := new(proto.SetSecureEmailResponse)
	err = suite.JinmuIDService.SetSecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:31000] invalid vc record"), err)
}

// TestSetSecureEmailIsNull   email为空
func (suite *SetSecureEmailTestSuite) TestSetSecureEmailIsNull() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getSetEmailSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(proto.SetSecureEmailRequest)
	req.UserId = userID
	req.VerificationCode = mvc
	req.SerialNumber = serialNumber
	req.Email = suite.Account.emailNull
	resp := new(proto.SetSecureEmailResponse)
	err = suite.JinmuIDService.SetSecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:23000] invalid email format"), err)
}

// TestSetSecureEmailFormatError    email格式错误
func (suite *SetSecureEmailTestSuite) TestSetSecureEmailFormatError() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getSetEmailSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(proto.SetSecureEmailRequest)
	req.UserId = userID
	req.VerificationCode = mvc
	req.SerialNumber = serialNumber
	req.Email = suite.Account.emailError
	resp := new(proto.SetSecureEmailResponse)
	err = suite.JinmuIDService.SetSecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:23000] invalid email format"), err)
}

// TestSetSecureEmailExist 邮箱已被其他用户使用
func (suite *SetSecureEmailTestSuite) TestSetSecureEmailExist() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
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
	assert.Error(t, errors.New("[errcode:25000] secure email has been set"), err)
}

func (suite *SetSecureEmailTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}

func TestSetSecureEmailTestSuite(t *testing.T) {
	suite.Run(t, new(SetSecureEmailTestSuite))
}
