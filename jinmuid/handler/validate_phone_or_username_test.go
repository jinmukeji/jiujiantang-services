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

// UserValidateUsernameOrPhoneTestSuite 验证手机号码和用户名是否存在测试
type UserValidateUsernameOrPhoneTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

// UserValidateUsernameOrPhoneTestSuite 设置测试环境
func (suite *UserValidateUsernameOrPhoneTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestUserValidateUsername 测试验证手机号码和用户名是否存在
func (suite *UserValidateUsernameOrPhoneTestSuite) TestUserValidateUsername() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.UserValidateUsernameOrPhoneRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_USERNAME,
		Username:       suite.Account.username,
	}
	resp := new(proto.UserValidateUsernameOrPhoneResponse)
	err := suite.JinmuIDService.UserValidateUsernameOrPhone(ctx, req, resp)
	assert.NoError(t, err)
}

// TestUserValidateUsernameNotExist 测试验证手机号码和用户名是否存在
func (suite *UserValidateUsernameOrPhoneTestSuite) TestUserValidateUsernameNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.UserValidateUsernameOrPhoneRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_USERNAME,
		Username:       suite.Account.usernameExist,
	}
	resp := new(proto.UserValidateUsernameOrPhoneResponse)
	err := suite.JinmuIDService.UserValidateUsernameOrPhone(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:21000] nonexistent username"), err)
}

// TestUserValidateUsernameIsNull 测试验证手机号码和用户名是否存在
func (suite *UserValidateUsernameOrPhoneTestSuite) TestUserValidateUsernameIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.UserValidateUsernameOrPhoneRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_USERNAME,
		Username:       suite.Account.usernameNull,
	}
	resp := new(proto.UserValidateUsernameOrPhoneResponse)
	err := suite.JinmuIDService.UserValidateUsernameOrPhone(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:21000] nonexistent username"), err)
}

// TestUserValidatePhone 测试验证手机号码和用户名是否存在
func (suite *UserValidateUsernameOrPhoneTestSuite) TestUserValidatePhone() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.UserValidateUsernameOrPhoneRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_PHONE,
		Phone:          suite.Account.phone,
		NationCode:     suite.Account.nationCode,
	}
	resp := new(proto.UserValidateUsernameOrPhoneResponse)
	err := suite.JinmuIDService.UserValidateUsernameOrPhone(ctx, req, resp)
	assert.NoError(t, err)
}

// TestUserValidatePhoneNotExist 测试验证手机号码和用户名是否存在
func (suite *UserValidateUsernameOrPhoneTestSuite) TestUserValidatePhoneNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.UserValidateUsernameOrPhoneRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_PHONE,
		Phone:          suite.Account.phoneNotExist,
		NationCode:     suite.Account.nationCode,
	}
	resp := new(proto.UserValidateUsernameOrPhoneResponse)
	err := suite.JinmuIDService.UserValidateUsernameOrPhone(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:22000] nonexistent phone"), err)
}

// TestUserValidatePhoneIsNull 测试验证手机号码和用户名是否存在
func (suite *UserValidateUsernameOrPhoneTestSuite) TestUserValidatePhoneIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.UserValidateUsernameOrPhoneRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_PHONE,
		Phone:          suite.Account.phoneIsNull,
		NationCode:     suite.Account.nationCode,
	}
	resp := new(proto.UserValidateUsernameOrPhoneResponse)
	err := suite.JinmuIDService.UserValidateUsernameOrPhone(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:22000] nonexistent phone"), err)
}

// TestUserValidatePhoneNationCodeIsNull 测试验证手机号码和用户名是否存在
func (suite *UserValidateUsernameOrPhoneTestSuite) TestUserValidatePhoneNationCodeIsNull() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.UserValidateUsernameOrPhoneRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_PHONE,
		Phone:          suite.Account.phone,
		NationCode:     suite.Account.nationCodeIsNull,
	}
	resp := new(proto.UserValidateUsernameOrPhoneResponse)
	err := suite.JinmuIDService.UserValidateUsernameOrPhone(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:22000] nonexistent phone"), err)
}

// TestUserValidatePhoneValidationTypeUnknown  validationtype为ValidationType_VALIDATION_TYPE_UNKNOWN
func (suite *UserValidateUsernameOrPhoneTestSuite) TestUserValidatePhoneValidationTypeUnknown() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.UserValidateUsernameOrPhoneRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_UNKNOWN,
		Phone:          suite.Account.phone,
		NationCode:     suite.Account.nationCode,
	}
	resp := new(proto.UserValidateUsernameOrPhoneResponse)
	err := suite.JinmuIDService.UserValidateUsernameOrPhone(ctx, req, resp)

	assert.Error(t, errors.New("[errcode:11000] invalid secure queston validation type"), err)
}

// TestUserValidatePhoneValidationTypeError validationtype为ValidationType_VALIDATION_TYPE_UNKNOWN
func (suite *UserValidateUsernameOrPhoneTestSuite) TestUserValidatePhoneValidationTypeError() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.UserValidateUsernameOrPhoneRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_USERNAME,
		Phone:          suite.Account.phone,
		NationCode:     suite.Account.nationCode,
	}
	resp := new(proto.UserValidateUsernameOrPhoneResponse)
	err := suite.JinmuIDService.UserValidateUsernameOrPhone(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:10000] invalid validation value"), err)
}

// TestUserValidatePhoneValidationTypeUsernameError validationtype为ValidationType_VALIDATION_TYPE_UNKNOWN
func (suite *UserValidateUsernameOrPhoneTestSuite) TestUserValidatePhoneValidationTypeUsernameError() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.UserValidateUsernameOrPhoneRequest{
		ValidationType: proto.ValidationType_VALIDATION_TYPE_PHONE,
		Username:       suite.Account.username,
	}
	resp := new(proto.UserValidateUsernameOrPhoneResponse)
	err := suite.JinmuIDService.UserValidateUsernameOrPhone(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:10000] invalid validation value"), err)
}

func (suite *UserValidateUsernameOrPhoneTestSuite) TestValidatePhoneFormat() {
	t := suite.T()
	err1 := validatePhoneFormat("900011110", "+886")
	assert.NoError(t, err1)
	err2 := validatePhoneFormat("68100000", "+852")
	assert.NoError(t, err2)
	err3 := validatePhoneFormat("7000000000", "+81")
	assert.NoError(t, err3)
}

func (suite *UserValidateUsernameOrPhoneTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}

func TestUserValidateUsernameOrPhoneTestSuite(t *testing.T) {
	suite.Run(t, new(UserValidateUsernameOrPhoneTestSuite))
}
