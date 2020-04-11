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

// UserGetInformation 获取用户信息
type UserGetInformationTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

// UserGetInformationTestSuite 设置测试环境
func (suite *UserGetInformationTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestUserGetInformation 测试获取用户信息
func (suite *UserGetInformationTestSuite) TestUserGetInformation() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	req := &proto.GetUserAndProfileInformationRequest{
		UserId: userID,
	}

	resp := new(proto.GetUserAndProfileInformationResponse)
	err = suite.JinmuIDService.GetUserAndProfileInformation(ctx, req, resp)
	assert.NoError(t, err)
}

// TestUserGetInformationInvalidError 用户名无效
func (suite *UserGetInformationTestSuite) TestUserGetInformationInvalidError() {
	t := suite.T()
	ctx := context.Background()
	req := &proto.GetUserAndProfileInformationRequest{
		UserId: 0,
	}
	resp := new(proto.GetUserAndProfileInformationResponse)
	err := suite.JinmuIDService.GetUserAndProfileInformation(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:10001] failed to find user profile"), err)
}

// TestUserGetInformationUsername 测试获取用户信息用户名登录
func (suite *UserGetInformationTestSuite) TestUserGetInformationUsername() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockUsernameSignin(ctx, suite.JinmuIDService, suite.Account.username, suite.Account.hashedPassword, suite.Account.seedNull)
	assert.NoError(t, err)
	req := &proto.GetUserAndProfileInformationRequest{
		UserId: userID,
	}

	resp := new(proto.GetUserAndProfileInformationResponse)
	err = suite.JinmuIDService.GetUserAndProfileInformation(ctx, req, resp)
	assert.NoError(t, err)
}

func (suite *UserGetInformationTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}

func TestUserGetInformationTestSuite(t *testing.T) {
	suite.Run(t, new(UserGetInformationTestSuite))
}
