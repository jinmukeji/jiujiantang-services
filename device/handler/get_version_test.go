package handler

import (
	"context"
	"path/filepath"
	"testing"

	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/device/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// GetVersionTestSuite 是 GetVersion rpc 的单元测试的 Test Suite
type GetVersionTestSuite struct {
	suite.Suite
	deviceManagerService *DeviceManagerService
}

// SetupSuite 设置测试环境
func (suite *GetVersionTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-device.env")
	db, _ := newTestingDbClientFromEnvFile(envFilepath)
	suite.deviceManagerService = NewDeviceManagerService(db)
}

// TestGetVersion 测试密码找回功能
func (suite *GetVersionTestSuite) TestGetVersion() {
	t := suite.T()
	ctx := context.Background()
	req := new(proto.GetVersionRequest)
	resp := new(proto.GetVersionResponse)
	err := suite.deviceManagerService.GetVersion(ctx, req, resp)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
func TestGetVersionTestSuite(t *testing.T) {
	suite.Run(t, new(GetVersionTestSuite))
}
