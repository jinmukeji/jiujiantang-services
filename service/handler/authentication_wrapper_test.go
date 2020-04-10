package handler

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

// AuthenticationWrapperTestSuite 是客户端认证中间件单元测试的 Test Suite
type AuthenticationWrapperTestSuite struct {
	suite.Suite
	jinmuHealth *JinmuHealth
}

// SetupSuite 设置测试环境
func (suite *AuthenticationWrapperTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-biz-core.env")
	suite.jinmuHealth = newTestingJinmuHealthFromEnvFile(envFilepath)
	suite.jinmuHealth.datastore, _ = newTestingDbClientFromEnvFile(envFilepath)
	suite.jinmuHealth.mailClient, _ = newTestingMailClientFromEnvFile(envFilepath)
}

func TestAuthenticationWrapperTestSuite(t *testing.T) {
	suite.Run(t, new(AuthenticationWrapperTestSuite))
}
