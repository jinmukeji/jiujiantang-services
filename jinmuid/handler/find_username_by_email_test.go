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

type FindUsernameBySecureEmailTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

//SetupSuite 设置测试环境
func (suite *FindUsernameBySecureEmailTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestFindUsernameBySecureEmail 测试根据邮箱找回用户名
func (suite *FindUsernameBySecureEmailTestSuite) TestFindUsernameBySecureEmail() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getFindUserNameSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getSmsVerificationCode(suite.JinmuIDService, *suite.Account)
	resp := new(proto.FindUsernameBySecureEmailResponse)
	req := new(proto.FindUsernameBySecureEmailRequest)
	req.VerificationCode = mvc
	req.SerialNumber = serialNumber
	req.Email = suite.Account.email
	err := suite.JinmuIDService.FindUsernameBySecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:21000] nonexistent username"), err)
}

// TestFindUsernameBySecureEmailIsNotSet 测试根据邮箱找回用户名,未设置邮箱和用户名
func (suite *FindUsernameBySecureEmailTestSuite) TestFindUsernameBySecureEmailIsNotSet() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getFindUserNameSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getSmsVerificationCode(suite.JinmuIDService, *suite.Account)
	resp := new(proto.FindUsernameBySecureEmailResponse)
	req := new(proto.FindUsernameBySecureEmailRequest)
	req.VerificationCode = mvc
	req.SerialNumber = serialNumber
	req.Email = suite.Account.email
	err := suite.JinmuIDService.FindUsernameBySecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:21000] nonexistent username"), err)
}

// TestFindUsernameBySecureEmailNotExist 测试根据邮箱找回用户名邮箱不存在
func (suite *FindUsernameBySecureEmailTestSuite) TestFindUsernameBySecureEmailNotExist() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getFindUserNameSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getSmsVerificationCode(suite.JinmuIDService, *suite.Account)
	resp := new(proto.FindUsernameBySecureEmailResponse)
	req := new(proto.FindUsernameBySecureEmailRequest)
	req.VerificationCode = mvc
	req.SerialNumber = serialNumber
	req.Email = suite.Account.emailNoExist
	err := suite.JinmuIDService.FindUsernameBySecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:10001] database error"), err)
}

// TestFindUsernameBySecureEmailVCIsNull  验证码为空
func (suite *FindUsernameBySecureEmailTestSuite) TestFindUsernameBySecureEmailVCIsNull() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getFindUserNameSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	// mvc := getSmsVerificationCode(suite.JinmuIDService,*suite.Account)
	resp := new(proto.FindUsernameBySecureEmailResponse)
	req := new(proto.FindUsernameBySecureEmailRequest)
	req.VerificationCode = suite.Account.mvcIsNull
	req.SerialNumber = serialNumber
	req.Email = suite.Account.emailNoExist
	err := suite.JinmuIDService.FindUsernameBySecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:10001] database error"), err)
}

// TestFindUsernameBySecureEmailVCError 验证码错误
func (suite *FindUsernameBySecureEmailTestSuite) TestFindUsernameBySecureEmailVCError() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getFindUserNameSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	// mvc := getSmsVerificationCode(suite.JinmuIDService,*suite.Account)
	resp := new(proto.FindUsernameBySecureEmailResponse)
	req := new(proto.FindUsernameBySecureEmailRequest)
	req.VerificationCode = suite.Account.mvcError
	req.SerialNumber = serialNumber
	req.Email = suite.Account.emailNoExist
	err := suite.JinmuIDService.FindUsernameBySecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:10001] database error"), err)
}

// TestFindUsernameBySecureEmailSerialNumberIsNull  serialNumber为空
func (suite *FindUsernameBySecureEmailTestSuite) TestFindUsernameBySecureEmailSerialNumberIsNull() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getFindUserNameSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getSmsVerificationCode(suite.JinmuIDService, *suite.Account)
	resp := new(proto.FindUsernameBySecureEmailResponse)
	req := new(proto.FindUsernameBySecureEmailRequest)
	req.VerificationCode = mvc
	req.SerialNumber = serialNumber
	req.Email = suite.Account.emailNoExist
	err := suite.JinmuIDService.FindUsernameBySecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:10001] database error"), err)
}

// TestFindUsernameBySecureEmailIsNull 邮箱地址为空
func (suite *FindUsernameBySecureEmailTestSuite) TestFindUsernameBySecureEmailIsNull() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getFindUserNameSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getSmsVerificationCode(suite.JinmuIDService, *suite.Account)
	resp := new(proto.FindUsernameBySecureEmailResponse)
	req := new(proto.FindUsernameBySecureEmailRequest)
	req.VerificationCode = mvc
	req.SerialNumber = serialNumber
	req.Email = suite.Account.emailNull
	err := suite.JinmuIDService.FindUsernameBySecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:10001] database error"), err)
}

// TestFindUsernameBySecureEmailError  邮箱地址不存在
func (suite *FindUsernameBySecureEmailTestSuite) TestFindUsernameBySecureEmailError() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getFindUserNameSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getSmsVerificationCode(suite.JinmuIDService, *suite.Account)
	resp := new(proto.FindUsernameBySecureEmailResponse)
	req := new(proto.FindUsernameBySecureEmailRequest)
	req.VerificationCode = mvc
	req.SerialNumber = serialNumber
	req.Email = suite.Account.emailError
	err := suite.JinmuIDService.FindUsernameBySecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:10001] database error"), err)
}

// TestFindUsernameBySecureEmailNoExist 邮箱地址不存在
func (suite *FindUsernameBySecureEmailTestSuite) TestFindUsernameBySecureEmailNoExist() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getFindUserNameSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getSmsVerificationCode(suite.JinmuIDService, *suite.Account)
	resp := new(proto.FindUsernameBySecureEmailResponse)
	req := new(proto.FindUsernameBySecureEmailRequest)
	req.VerificationCode = mvc
	req.SerialNumber = serialNumber
	req.Email = suite.Account.emailNoExist
	err := suite.JinmuIDService.FindUsernameBySecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:21000] nonexistent username"), err)
}

// TestFindUsernameBySecureEmailVCIsExpired  验证码过期
func (suite *FindUsernameBySecureEmailTestSuite) TestFindUsernameBySecureEmailVCIsExpired() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getFindUserNameSerialNumber(suite.JinmuIDService, *suite.Account)
	resp := new(proto.FindUsernameBySecureEmailResponse)
	req := new(proto.FindUsernameBySecureEmailRequest)
	req.VerificationCode = "126462"
	req.SerialNumber = serialNumber
	req.Email = suite.Account.email
	err := suite.JinmuIDService.FindUsernameBySecureEmail(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:21000] nonexistent username"), err)
}

func (suite *FindUsernameBySecureEmailTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}
func TestFindUsernameBySecureEmailTestSuite(t *testing.T) {
	suite.Run(t, new(FindUsernameBySecureEmailTestSuite))

}
