package handler

import (
	"context"
	"path/filepath"
	"testing"

	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// GetUserProfileTestSuite 是 GetUserProfile rpc 的单元测试的 Test Suite
type GetUserProfileTestSuite struct {
	suite.Suite
	jinmuHealth *JinmuHealth
}

// SetupSuite 设置测试环境
func (suite *GetUserProfileTestSuite) SetupSuite() {
	suite.jinmuHealth = new(JinmuHealth)
	envFilepath := filepath.Join("testdata", "local.svc-biz-core.env")
	suite.jinmuHealth.datastore, _ = newTestingDbClientFromEnvFile(envFilepath)
	suite.jinmuHealth.mailClient, _ = newTestingMailClientFromEnvFile(envFilepath)
}

// TestGetUserProfile 测试查看某个用户个人档案
func (suite *GetUserProfileTestSuite) TestGetUserProfile() {
	t := suite.T()
	ctx := context.Background()
	req := new(proto.GetUserProfileRequest)
	resp := new(proto.GetUserProfileResponse)
	req.UserId = 1
	err := suite.jinmuHealth.GetUserProfile(ctx, req, resp)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.GetUser())
}
func TestGetUserProfileTestSuite(t *testing.T) {
	suite.Run(t, new(GetUserProfileTestSuite))
}
