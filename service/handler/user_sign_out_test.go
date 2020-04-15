package handler

import (
	"context"
	"path/filepath"
	"testing"

	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// UserSignOutTestSuite 是用户登录的单元测试的 Test Suite
type UserSignOutTestSuite struct {
	suite.Suite
	jinmuHealth *JinmuHealth
}

// SetupSuite 设置测试环境
func (suite *UserSignOutTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-biz-core.env")
	suite.jinmuHealth = newTestingJinmuHealthFromEnvFile(envFilepath)
	suite.jinmuHealth.datastore, _ = newTestingDbClientFromEnvFile(envFilepath)
	suite.jinmuHealth.mailClient, _ = newTestingMailClientFromEnvFile(envFilepath)
}

// TestUserSiginin 测试用户登录
func (suite *UserSignOutTestSuite) TestUserSignOut() {
	// 先进行登录操作
	t := suite.T()
	ctx := context.Background()
	const username = "4"
	const passwordHash = "97951af80347d78d63bf3a7b7962fb42dd21ef56d4e590f2be2b7954475f4089"
	const clientID = "jm-10005"
	const name = "JinmuHealth-Android-app"
	const zone = "CN"
	ctx = mockAuth(ctx, clientID, name, zone)
	const registerType = "username"
	ctx, err := mockSignin(ctx, suite.jinmuHealth, username, passwordHash, registerType, proto.SignInMethod_SIGN_IN_METHOD_GENERAL)
	assert.NoError(t, err)
	// 退出登录操作
	signoutReq, signoutResp := new(proto.UserSignOutRequest), new(proto.UserSignOutResponse)
	assert.NoError(t, suite.jinmuHealth.UserSignOut(ctx, signoutReq, signoutResp))
}

func TestUserSignOutTestSuite(t *testing.T) {
	suite.Run(t, new(UserSignOutTestSuite))
}
