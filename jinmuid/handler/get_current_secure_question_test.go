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

// UserGetSecureQuestionsByPhoneOrUsernameTestSuite 根据用户名或者手机号获取当前设置的密保问题测试
type UserGetSecureQuestionsByPhoneOrUsernameTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

// UserGetSecureQuestionsByPhoneOrUsernameTestSuite 设置测试环境
func (suite *UserGetSecureQuestionsByPhoneOrUsernameTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestGetSecureQuestionsByUsername 测试根据用户名获取当前设置的密保问题
func (suite *UserGetSecureQuestionsByPhoneOrUsernameTestSuite) TestGetSecureQuestionsByUsername() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.GetSecureQuestionsByPhoneOrUsernameRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_USERNAME,
		Username:       suite.Account.usernameExist,
	}
	resp := new(proto.GetSecureQuestionsByPhoneOrUsernameResponse)
	err := suite.JinmuIDService.GetSecureQuestionsByPhoneOrUsername(ctx, req, resp)
	assert.NoError(t, err)
}

// TestGetSecureQuestionsByPhoneOrUsernameSecureQustionsNotSet   用户未设置密保问题
func (suite *UserGetSecureQuestionsByPhoneOrUsernameTestSuite) TestGetSecureQuestionsByPhoneOrUsernameSecureQustionsNotSet() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.GetSecureQuestionsByPhoneOrUsernameRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_USERNAME,
		Username:       suite.Account.username,
	}
	resp := new(proto.GetSecureQuestionsByPhoneOrUsernameResponse)
	err := suite.JinmuIDService.GetSecureQuestionsByPhoneOrUsername(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:49000] nonexistent secure questions"), err)
}

// TestGetSecureQuestionsByPhone  根据手机号获取当前设置的密保问题
func (suite *UserGetSecureQuestionsByPhoneOrUsernameTestSuite) TestGetSecureQuestionsByPhone() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.GetSecureQuestionsByPhoneOrUsernameRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_PHONE,
		Phone:          suite.Account.phone,
		NationCode:     suite.Account.nationCode,
	}
	resp := new(proto.GetSecureQuestionsByPhoneOrUsernameResponse)
	err := suite.JinmuIDService.GetSecureQuestionsByPhoneOrUsername(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:49000] nonexistent secure questions"), err)
}

// TestGetSecureQuestionsByPhoneIsNull     手机号为空
func (suite *UserGetSecureQuestionsByPhoneOrUsernameTestSuite) TestGetSecureQuestionsByPhoneIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.GetSecureQuestionsByPhoneOrUsernameRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_PHONE,
		Phone:          suite.Account.phoneIsNull,
		NationCode:     suite.Account.nationCode,
	}
	resp := new(proto.GetSecureQuestionsByPhoneOrUsernameResponse)
	err := suite.JinmuIDService.GetSecureQuestionsByPhoneOrUsername(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:49000] nonexistent secure questions"), err)
}

// TestGetSecureQuestionsByPhoneFormatError  手机号格式为空
func (suite *UserGetSecureQuestionsByPhoneOrUsernameTestSuite) TestGetSecureQuestionsByPhoneFormatError() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.GetSecureQuestionsByPhoneOrUsernameRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_PHONE,
		Phone:          suite.Account.PhoneError,
		NationCode:     suite.Account.nationCode,
	}
	resp := new(proto.GetSecureQuestionsByPhoneOrUsernameResponse)
	err := suite.JinmuIDService.GetSecureQuestionsByPhoneOrUsername(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:49000] nonexistent secure questions"), err)
}

// TestGetSecureQuestionsByPhoneNationCodeIsNull   NationCode为空
func (suite *UserGetSecureQuestionsByPhoneOrUsernameTestSuite) TestGetSecureQuestionsByPhoneNationCodeIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.GetSecureQuestionsByPhoneOrUsernameRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_PHONE,
		Phone:          suite.Account.phone,
		NationCode:     suite.Account.nationCode,
	}
	resp := new(proto.GetSecureQuestionsByPhoneOrUsernameResponse)
	err := suite.JinmuIDService.GetSecureQuestionsByPhoneOrUsername(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:49000] nonexistent secure questions"), err)
}

// TestGetSecureQuestionsByUsernameIsNull     用户名为空
func (suite *UserGetSecureQuestionsByPhoneOrUsernameTestSuite) TestGetSecureQuestionsByUsernameIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.GetSecureQuestionsByPhoneOrUsernameRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_PHONE,
		Username:       suite.Account.usernameNull,
	}
	resp := new(proto.GetSecureQuestionsByPhoneOrUsernameResponse)
	err := suite.JinmuIDService.GetSecureQuestionsByPhoneOrUsername(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:49000] nonexistent secure questions"), err)
}

func (suite *UserGetSecureQuestionsByPhoneOrUsernameTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}
func TestGetSecureQuestionsByPhoneOrUsernameTestSuite(t *testing.T) {
	suite.Run(t, new(UserGetSecureQuestionsByPhoneOrUsernameTestSuite))
}
