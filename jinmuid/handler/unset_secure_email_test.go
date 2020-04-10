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

//  UnsetSecureEmailTestSuite 通过密保问题重置密码测试
type UnsetSecureEmailTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

// UnsetSecureEmailTestSuite 设置测试环境
func (suite *UnsetSecureEmailTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestUnsetSecureEmail 测试解绑安全邮箱
func (suite *UnsetSecureEmailTestSuite) TestUnsetSecureEmail() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getUnSetEmailSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(proto.UnsetSecureEmailRequest)
	req.UserId = userID
	req.VerificationCode = mvc
	req.SerialNumber = serialNumber
	req.Email = suite.Account.email
	resp := new(proto.UnsetSecureEmailResponse)
	err = suite.JinmuIDService.UnsetSecureEmail(ctx, req, resp)
	assert.NoError(t, err)
}

// TestUnsetSecureEmailUserIdIsNull  设置安全邮箱userid为空
func (suite *UnsetSecureEmailTestSuite) TestUnsetSecureEmailUserIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getUnSetEmailSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	resp := new(proto.UnsetSecureEmailResponse)
	req := new(proto.UnsetSecureEmailRequest)
	req.UserId = suite.Account.userID
	req.VerificationCode = mvc
	req.SerialNumber = serialNumber
	req.Email = suite.Account.email
	err := suite.JinmuIDService.UnsetSecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:25000] secure email has been set"), err)
}

// TestUnsetSecureEmailVCIsNull     verificationCode为空
func (suite *UnsetSecureEmailTestSuite) TestUnsetSecureEmailVCIsNull() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getUnSetEmailSerialNumber(suite.JinmuIDService, *suite.Account)
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(proto.UnsetSecureEmailRequest)
	req.UserId = userID
	req.VerificationCode = suite.Account.mvcIsNull
	req.SerialNumber = serialNumber
	req.Email = suite.Account.email
	resp := new(proto.UnsetSecureEmailResponse)
	err = suite.JinmuIDService.UnsetSecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:31000] invalid vc record"), err)
}

// TestUnsetSecureEmailVCError     verificationCode错误,过期
func (suite *UnsetSecureEmailTestSuite) TestUnsetSecureEmailVCError() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getUnSetEmailSerialNumber(suite.JinmuIDService, *suite.Account)
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(proto.UnsetSecureEmailRequest)
	req.UserId = userID
	req.VerificationCode = suite.Account.mvcError
	req.SerialNumber = serialNumber
	req.Email = suite.Account.email
	resp := new(proto.UnsetSecureEmailResponse)
	err = suite.JinmuIDService.UnsetSecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:31000] invalid vc record"), err)
}

// TestUnsetSecureEmailSerialNumberIsNull    serialnumber为空
func (suite *SetSecureEmailTestSuite) TestUnsetSecureEmailSerialNumberIsNull() {
	t := suite.T()
	ctx := context.Background()
	// 获取最新验证码
	mvc := getEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(proto.UnsetSecureEmailRequest)
	req.UserId = userID
	req.VerificationCode = mvc
	req.SerialNumber = suite.Account.serialNumberIsNull
	req.Email = suite.Account.email
	resp := new(proto.UnsetSecureEmailResponse)
	err = suite.JinmuIDService.UnsetSecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:31000] invalid vc record"), err)
}

// TestUnsetSecureEmailIsNull   email为空
func (suite *SetSecureEmailTestSuite) TestUnsetSecureEmailIsNull() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getUnSetEmailSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(proto.UnsetSecureEmailRequest)
	req.UserId = userID
	req.VerificationCode = mvc
	req.SerialNumber = serialNumber
	req.Email = suite.Account.emailNull
	resp := new(proto.UnsetSecureEmailResponse)
	err = suite.JinmuIDService.UnsetSecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:23000] invalid email format"), err)
}

// TestUnsetSecureEmailFormatError    email格式错误
func (suite *SetSecureEmailTestSuite) TestUnsetSecureEmailFormatError() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getUnSetEmailSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(proto.UnsetSecureEmailRequest)
	req.UserId = userID
	req.VerificationCode = mvc
	req.SerialNumber = serialNumber
	req.Email = suite.Account.emailError
	resp := new(proto.UnsetSecureEmailResponse)
	err = suite.JinmuIDService.UnsetSecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:23000] invalid email format"), err)
}

// TestUnsetSecureEmailExist 邮箱已被其他用户使用
func (suite *SetSecureEmailTestSuite) TestUnsetSecureEmailExist() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getUnSetEmailSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(proto.UnsetSecureEmailRequest)
	req.UserId = userID
	req.VerificationCode = mvc
	req.SerialNumber = serialNumber
	req.Email = suite.Account.email
	resp := new(proto.UnsetSecureEmailResponse)
	err = suite.JinmuIDService.UnsetSecureEmail(ctx, req, resp)
	assert.NoError(t, err)
	assert.Error(t, errors.New("[errcode:25000] secure email has been set"), err)
}

func (suite *UnsetSecureEmailTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}

func TestUnsetSecureEmailTestSuite(t *testing.T) {
	suite.Run(t, new(UnsetSecureEmailTestSuite))
}
