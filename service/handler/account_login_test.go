package handler

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	"github.com/stretchr/testify/suite"
)

// AccountLoginTestSuite 账户登录的单元测试的 Test Suite
type AccountLoginTestSuite struct {
	suite.Suite
	jinmuHealth *JinmuHealth
	Account     *Account
}

// SetupSuite 设置测试环境
func (suite *AccountLoginTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-biz-core.env")
	suite.jinmuHealth = newTestingJinmuHealthFromEnvFile(envFilepath)
	suite.jinmuHealth.datastore, _ = newTestingDbClientFromEnvFile(envFilepath)
	suite.jinmuHealth.mailClient, _ = newTestingMailClientFromEnvFile(envFilepath)
	suite.Account = newTestingAccountFromEnvFile(envFilepath)

}

// TestAccountLogin  测试账户登录
func (suite *AccountLoginTestSuite) TestAccountLogin() {
	t := suite.T()
	ctx := context.Background()
	req := new(proto.JinmuLAccountLoginRequest)
	resp := new(proto.JinmuLAccountLoginResponse)
	req.Account = suite.Account.userAccount
	req.Password = suite.Account.password
	req.MachineId = suite.Account.machineUuid
	err := suite.jinmuHealth.JinmuLAccountLogin(ctx, req, resp)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
}

// TestErrorAccount  测试账户登录-账户格式错误场景
func (suite *AccountLoginTestSuite) TestErrorAccount() {
	t := suite.T()
	ctx := context.Background()
	req := new(proto.JinmuLAccountLoginRequest)
	resp := new(proto.JinmuLAccountLoginResponse)
	req.Account = suite.Account.errorAccount
	req.Password = suite.Account.password
	req.MachineId = suite.Account.machineUuid
	err := suite.jinmuHealth.JinmuLAccountLogin(ctx, req, resp)
	assert.Error(t, errors.New("error[1600]"), err)
}

// TestAccountNotExist  测试账户登录-账户不存在
func (suite *AccountLoginTestSuite) TestAccountNotExist() {
	t := suite.T()
	ctx := context.Background()
	req := new(proto.JinmuLAccountLoginRequest)
	resp := new(proto.JinmuLAccountLoginResponse)
	req.Account = suite.Account.userAccount
	req.Password = suite.Account.password
	req.MachineId = suite.Account.machineUuid
	err := suite.jinmuHealth.JinmuLAccountLogin(ctx, req, resp)
	assert.NoError(t, err)
	assert.Error(t, errors.New("error[12012]"), err)
}

// TestErrorPassword 测试账户登录-密码错误场景
func (suite *AccountLoginTestSuite) TestErrorPassword() {
	t := suite.T()
	ctx := context.Background()
	req := new(proto.JinmuLAccountLoginRequest)
	resp := new(proto.JinmuLAccountLoginResponse)
	req.Account = suite.Account.userAccount
	req.Password = suite.Account.errorPassword
	req.MachineId = suite.Account.machineUuid
	err := suite.jinmuHealth.JinmuLAccountLogin(ctx, req, resp)
	assert.NoError(t, err)
	assert.Error(t, errors.New("error[12000]"), err)
}

// TestErrorFormatPassword  测试账户登录-密码格式错误
func (suite *AccountLoginTestSuite) TestErrorFormatPassword() {
	t := suite.T()
	ctx := context.Background()
	req := new(proto.JinmuLAccountLoginRequest)
	resp := new(proto.JinmuLAccountLoginResponse)
	req.Account = suite.Account.userAccount
	req.Password = suite.Account.errorFormatPassword
	req.MachineId = suite.Account.machineUuid
	err := suite.jinmuHealth.JinmuLAccountLogin(ctx, req, resp)
	assert.NoError(t, err)
	assert.Error(t, errors.New("error[12015]"), err)
}

// TestErrorMachineId  测试账户登录-machine_uuid 为空场景
func (suite *AccountLoginTestSuite) TestErrorMachineId() {
	t := suite.T()
	ctx := context.Background()
	req := new(proto.JinmuLAccountLoginRequest)
	resp := new(proto.JinmuLAccountLoginResponse)
	req.Account = suite.Account.userAccount
	req.Password = suite.Account.password
	req.MachineId = suite.Account.machineUuid
	err := suite.jinmuHealth.JinmuLAccountLogin(ctx, req, resp)
	assert.NoError(t, err)
	assert.Error(t, errors.New("error[11007]"), err)
}

// 关闭数据库
func TestAccountLoginTestSuite(t *testing.T) {
	suite.Run(t, new(AccountLoginTestSuite))
}
