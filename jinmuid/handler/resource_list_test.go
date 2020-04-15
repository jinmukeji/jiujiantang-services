package handler

import (
	"context"
	"testing"

	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ResourceListTestSuite 资源列表单元测试
type ResourceListTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
}

// ResourceListTestSuite 设置测试环境
func (suite *ResourceListTestSuite) SetupSuite() {
	suite.JinmuIDService = newJinmuIDServiceForTest()
}

// TestGetResList 测试资源列表
func (suite *ResourceListTestSuite) TestGerResourceList() {
	t := suite.T()
	ctx := context.Background()

	req := new(proto.GerResourceListRequest)
	resp := new(proto.GerResourceListResponse)
	err := suite.JinmuIDService.GerResourceList(ctx, req, resp)
	assert.NoError(t, err)
}

func (suite *ResourceListTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}

func TestResourceListTestSuite(t *testing.T) {
	suite.Run(t, new(ResourceListTestSuite))
}
