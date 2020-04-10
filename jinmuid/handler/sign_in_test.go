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

// UserSignInTestSuite 用户测试
type UserSignInTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

// SetupSuite 设置测试环境
func (suite *UserSignInTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestUserSignInByPhonePassword 测试电话密码登录
func (suite *UserSignInTestSuite) TestUserSignInByPhonePassword() {
	t := suite.T()
	ctx := context.Background()
	resp := new(proto.UserSignInByPhonePasswordResponse)
	err := suite.JinmuIDService.UserSignInByPhonePassword(ctx, &proto.UserSignInByPhonePasswordRequest{
		Phone:          suite.Account.phone,
		HashedPassword: suite.Account.phonePassword,
		Seed:           suite.Account.seed,
		NationCode:     suite.Account.nationCode,
	}, resp)
	assert.NoError(t, err)
	assert.Equal(t, 105586, int(resp.UserId))
}

// TestUserSignInByUsernamePassword 测试用户名密码登录
func (suite *UserSignInTestSuite) TestUserSignInByUsernamePassword() {
	t := suite.T()
	ctx := context.Background()
	resp := new(proto.UserSignInByUsernamePasswordResponse)
	err := suite.JinmuIDService.UserSignInByUsernamePassword(ctx, &proto.UserSignInByUsernamePasswordRequest{
		Username:       suite.Account.username,
		HashedPassword: suite.Account.hashedPassword,
		Seed:           suite.Account.seedNull,
	}, resp)
	assert.NoError(t, err)
	assert.Equal(t, 786, int(resp.UserId))
}

// TestUserSignInByPhoneVC 测试手机验证码登录
func (suite *UserSignInTestSuite) TestUserSignInByPhoneVC() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getSignInSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getSmsVerificationCode(suite.JinmuIDService, *suite.Account)
	resp := new(proto.UserSignInByPhoneVCResponse)
	req := new(proto.UserSignInByPhoneVCRequest)
	req.Phone = suite.Account.phone
	req.Mvc = mvc
	req.SerialNumber = serialNumber
	req.NationCode = suite.Account.nationCode
	err := suite.JinmuIDService.UserSignInByPhoneVC(ctx, req, resp)
	assert.NoError(t, err)
}

// mockSigninByPhonePassword模拟登录
func mockSigninByPhonePassword(ctx context.Context, j *JinmuIDService, phone string, hashedPassword, seed, nationCode string) (context.Context, int32, error) {
	resp := new(proto.UserSignInByPhonePasswordResponse)
	err := j.UserSignInByPhonePassword(ctx, &proto.UserSignInByPhonePasswordRequest{
		Phone:          phone,
		HashedPassword: hashedPassword,
		Seed:           seed,
		NationCode:     nationCode,
	}, resp)
	if err != nil {
		return nil, int32(0), err
	}
	return AddContextToken(ctx, resp.AccessToken), resp.UserId, nil
}

// mockSignin 模拟登录
func mockUsernameSignin(ctx context.Context, j *JinmuIDService, username string, passwordHash, seed string) (context.Context, int32, error) {
	resp := new(proto.UserSignInByUsernamePasswordResponse)
	err := j.UserSignInByUsernamePassword(ctx, &proto.UserSignInByUsernamePasswordRequest{
		Username:       username,
		HashedPassword: passwordHash,
		Seed:           seed,
	}, resp)
	if err != nil {
		return nil, int32(0), err
	}
	return AddContextToken(ctx, resp.AccessToken), resp.UserId, nil
}

// TestUserSignInByPhoneNotExist 手机未注册
func (suite *UserSignInTestSuite) TestUserSignInByPhoneNotExist() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getSignInSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getSmsVerificationCode(suite.JinmuIDService, *suite.Account)
	resp := new(proto.UserSignInByPhoneVCResponse)
	req := new(proto.UserSignInByPhoneVCRequest)
	req.Phone = suite.Account.phoneNotExist
	req.Mvc = mvc
	req.SerialNumber = serialNumber
	err := suite.JinmuIDService.UserSignInByPhoneVC(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:22000] nonexistent phone"), err)
}

// TestUserSignInByPhoneFormatIsError 手机号格式错误
func (suite *UserSignInTestSuite) TestUserSignInByPhoneFormatIsError() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getSignInSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getSmsVerificationCode(suite.JinmuIDService, *suite.Account)
	resp := new(proto.UserSignInByPhoneVCResponse)
	req := new(proto.UserSignInByPhoneVCRequest)
	req.Phone = suite.Account.PhoneError
	req.Mvc = mvc
	req.SerialNumber = serialNumber
	err := suite.JinmuIDService.UserSignInByPhoneVC(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:22000] nonexistent phone"), err)
}

// TestUserSignInByPhoneIsNull  登录手机号为空
func (suite *UserSignInTestSuite) TestUserSignInByPhoneIsNull() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getSignInSerialNumber(suite.JinmuIDService, *suite.Account)
	// 获取最新验证码
	mvc := getSmsVerificationCode(suite.JinmuIDService, *suite.Account)
	resp := new(proto.UserSignInByPhoneVCResponse)
	req := new(proto.UserSignInByPhoneVCRequest)
	req.Phone = suite.Account.phoneIsNull
	req.Mvc = mvc
	req.SerialNumber = serialNumber
	err := suite.JinmuIDService.UserSignInByPhoneVC(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:22000] nonexistent phone"), err)
}

// TestUserSignInByPhoneMvcIsNull 验证码为空
func (suite *UserSignInTestSuite) TestUserSignInByPhoneMvcIsNull() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	serialNumber := getSignInSerialNumber(suite.JinmuIDService, *suite.Account)
	resp := new(proto.UserSignInByPhoneVCResponse)
	req := new(proto.UserSignInByPhoneVCRequest)
	req.Phone = suite.Account.phone
	req.Mvc = suite.Account.mvcIsNull
	req.SerialNumber = serialNumber
	err := suite.JinmuIDService.UserSignInByPhoneVC(ctx, req, resp)
	assert.NoError(t, err)
	assert.Error(t, errors.New("[errcode:31000] vc is invalid"), err)
}

// TestUserSignInByPhoneSerialNumberIsNull  SerialNumber为空
func (suite *UserSignInTestSuite) TestUserSignInByPhoneSerialNumberIsNull() {
	t := suite.T()
	ctx := context.Background()
	// 获取最新验证码
	mvc := getSmsVerificationCode(suite.JinmuIDService, *suite.Account)
	resp := new(proto.UserSignInByPhoneVCResponse)
	req := new(proto.UserSignInByPhoneVCRequest)
	req.Phone = suite.Account.phone
	req.Mvc = mvc
	req.SerialNumber = suite.Account.serialNumberIsNull
	err := suite.JinmuIDService.UserSignInByPhoneVC(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:31000] vc is invalid"), err)
}

// TestUserSignInByUsernamePasswordIsNull  密码为空
func (suite *UserSignInTestSuite) TestUserSignInByUsernamePasswordIsNull() {
	t := suite.T()
	ctx := context.Background()
	resp := new(proto.UserSignInByUsernamePasswordResponse)
	err := suite.JinmuIDService.UserSignInByUsernamePassword(ctx, &proto.UserSignInByUsernamePasswordRequest{
		Username:       suite.Account.phone,
		HashedPassword: suite.Account.hashedPasswordIsNull,
		Seed:           suite.Account.seedNull,
	}, resp)
	assert.Error(t, errors.New("[errcode:21000] username does not exist"), err)
}

// TestUserSignInByUsernamePasswordIsError  密码错误
func (suite *UserSignInTestSuite) TestUserSignInByUsernamePasswordIsError() {
	t := suite.T()
	ctx := context.Background()
	resp := new(proto.UserSignInByUsernamePasswordResponse)
	err := suite.JinmuIDService.UserSignInByUsernamePassword(ctx, &proto.UserSignInByUsernamePasswordRequest{
		Username:       suite.Account.username,
		HashedPassword: suite.Account.password,
		Seed:           suite.Account.seedNull,
	}, resp)
	assert.Error(t, errors.New("[errcode:31000] vc is invalid"), err)
}

// TestUserSignInByPhoneUserIsDeactivated 用户被停用
func (suite *UserSignInTestSuite) TestUserSignInByPhoneUserIsDeactivated() {
	t := suite.T()
	ctx := context.Background()
	resp := new(proto.UserSignInByPhonePasswordResponse)
	err := suite.JinmuIDService.UserSignInByPhonePassword(ctx, &proto.UserSignInByPhonePasswordRequest{
		Phone:          suite.Account.phone,
		HashedPassword: suite.Account.hashedPassword,
		Seed:           suite.Account.seed,
		NationCode:     suite.Account.nationCode,
	}, resp)
	assert.Error(t, errors.New("[errcode:2700] user is deactivated"), err)
}

//TestUserSignInByPhonePasswordNotMatch
func (suite *UserSignInTestSuite) TestUserSignInByPhonePasswordNotMatch() {
	t := suite.T()
	ctx := context.Background()
	resp := new(proto.UserSignInByPhonePasswordResponse)
	err := suite.JinmuIDService.UserSignInByPhonePassword(ctx, &proto.UserSignInByPhonePasswordRequest{
		Phone:          suite.Account.phone,
		HashedPassword: suite.Account.hashedPassword,
		Seed:           suite.Account.seed,
		NationCode:     suite.Account.nationCode,
	}, resp)
	assert.Error(t, errors.New("[errcode:1700] password and phone does not match"), err)
}

func (suite *UserSignInTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}

func TestUserSignInTestSuite(t *testing.T) {
	suite.Run(t, new(UserSignInTestSuite))
}
