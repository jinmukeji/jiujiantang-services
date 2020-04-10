package handler

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

// SignUpViaPassword 获取用户信息
type SignUpViaPasswordTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

// SignUpViaPassword 设置测试环境
func (suite *SignUpViaPasswordTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

func (suite *SignUpViaPasswordTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}

func TestSignUpViaPasswordTestSuite(t *testing.T) {
	suite.Run(t, new(SignUpViaPasswordTestSuite))
}
