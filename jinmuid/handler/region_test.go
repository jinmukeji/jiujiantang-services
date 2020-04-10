package handler

import (
	"context"
	"path/filepath"
	"testing"

	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RegionTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

// SetupSuite 设置测试环境
func (suite *RegionTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestRegionTaiwan  测试选择区域中国台湾
func (suite *RegionTestSuite) TestRegionTaiwan() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(proto.UserSelectRegionRequest)
	req.UserId = userID
	req.Region = proto.Region_REGION_TAIWAN
	resp := new(proto.UserSelectRegionResponse)
	err = suite.JinmuIDService.UserSelectRegion(ctx, req, resp)
	assert.NoError(t, err)
}

// TestRegionMainlandChina  测试选择区域中国大陆
func (suite *RegionTestSuite) TestRegionMainlandChina() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(proto.UserSelectRegionRequest)
	req.UserId = userID
	req.Region = proto.Region_REGION_MAINLAND_CHINA
	resp := new(proto.UserSelectRegionResponse)
	err = suite.JinmuIDService.UserSelectRegion(ctx, req, resp)
	assert.NoError(t, err)
}

// TestRegionAbroad 测试选择区域中国国外
func (suite *RegionTestSuite) TestRegionAbroad() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(proto.UserSelectRegionRequest)
	req.UserId = userID
	req.Region = proto.Region_REGION_ABROAD
	resp := new(proto.UserSelectRegionResponse)
	err = suite.JinmuIDService.UserSelectRegion(ctx, req, resp)
	assert.NoError(t, err)
}

func (suite *RegionTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}

func TestRegionTestSuite(t *testing.T) {
	suite.Run(t, new(RegionTestSuite))
}
