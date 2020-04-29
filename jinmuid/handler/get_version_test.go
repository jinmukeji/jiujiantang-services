package handler

import (
	"context"
	"path/filepath"
	"testing"

	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GetVersionTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

// SetupSuite 设置测试环境
func (suite *GetVersionTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestGetVersion 获取版本号
func (suite *GetVersionTestSuite) TestGetVersion() {
	t := suite.T()
	ctx := context.Background()
	resp := new(proto.GetVersionResponse)
	req := new(proto.GetVersionRequest)
	err := suite.JinmuIDService.GetVersion(ctx, req, resp)
	assert.NoError(t, err)
}

func (suite *GetVersionTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}
func TestGetVersionTestSuite(t *testing.T) {
	suite.Run(t, new(GetVersionTestSuite))
}
