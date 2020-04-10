package handler

import (
	"context"
	"path/filepath"
	"testing"

	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// EchoTestSuite 是 Echo rpc 的单元测试的 Test Suite
type EchoTestSuite struct {
	suite.Suite
	jinmuHealth *JinmuHealth
}

// SetupSuite 设置测试环境
func (suite *EchoTestSuite) SetupSuite() {
	suite.jinmuHealth = new(JinmuHealth)
	envFilepath := filepath.Join("testdata", "local.svc-biz-core.env")
	suite.jinmuHealth.datastore, _ = newTestingDbClientFromEnvFile(envFilepath)
	suite.jinmuHealth.mailClient, _ = newTestingMailClientFromEnvFile(envFilepath)
}

// TestEcho 测试Echo回显
func (suite *EchoTestSuite) TestEcho() {
	t := suite.T()
	ctx := context.Background()
	req := new(proto.EchoRequest)
	req.Content = "This is Echo test."
	resp := new(proto.EchoResponse)
	err := suite.jinmuHealth.Echo(ctx, req, resp)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Content, resp.Content)
}
func TestEchoTestSuite(t *testing.T) {
	suite.Run(t, new(EchoTestSuite))
}
