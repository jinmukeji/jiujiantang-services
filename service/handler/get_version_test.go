package handler

import (
	"context"
	"path/filepath"
	"testing"

	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// GetVersionTestSuite 是 GetVersion rpc 的单元测试的 Test Suite
type GetVersionTestSuite struct {
	suite.Suite
	jinmuHealth *JinmuHealth
}

// SetupSuite 设置测试环境
func (suite *GetVersionTestSuite) SetupSuite() {
	suite.jinmuHealth = new(JinmuHealth)
	envFilepath := filepath.Join("testdata", "local.svc-biz-core.env")
	suite.jinmuHealth.datastore, _ = newTestingDbClientFromEnvFile(envFilepath)
	suite.jinmuHealth.mailClient, _ = newTestingMailClientFromEnvFile(envFilepath)
}

// TestGetVersion 测试密码找回功能
func (suite *GetVersionTestSuite) TestGetVersion() {
	t := suite.T()
	ctx := context.Background()
	req := new(proto.GetVersionRequest)
	resp := new(proto.GetVersionResponse)
	err := suite.jinmuHealth.GetVersion(ctx, req, resp)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
func TestGetVersionTestSuite(t *testing.T) {
	suite.Run(t, new(GetVersionTestSuite))
}
