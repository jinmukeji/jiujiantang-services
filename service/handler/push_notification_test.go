package handler

import (
	"context"
	"path/filepath"

	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// PushNotification是 PushNotification rpc 的单元测试的 Test Suite
type PushNotificationTestSuite struct {
	suite.Suite
	jinmuHealth *JinmuHealth
}

// SetupSuite 设置测试环境
func (suite *PushNotificationTestSuite) SetupSuite() {
	suite.jinmuHealth = new(JinmuHealth)
	envFilepath := filepath.Join("testdata", "local.svc-biz-core.env")
	suite.jinmuHealth.datastore, _ = newTestingDbClientFromEnvFile(envFilepath)
	suite.jinmuHealth.mailClient, _ = newTestingMailClientFromEnvFile(envFilepath)
}

// TestGetPushNotifications 测试查找通知
func (suite *PushNotificationTestSuite) TestGetPushNotifications() {
	t := suite.T()
	ctx := context.Background()
	req := new(proto.GetPushNotificationsRequest)
	req.UserId = 96
	resp := new(proto.GetPushNotificationsResponse)
	assert.NoError(t, suite.jinmuHealth.GetPushNotifications(ctx, req, resp))
}

// TestReadPushNotification 测试阅读通知
func (suite *PushNotificationTestSuite) TestReadPushNotification() {
	t := suite.T()
	ctx := context.Background()
	req := new(proto.ReadPushNotificationRequest)
	req.PnId = 5
	req.UserId = 96
	resp := new(proto.ReadPushNotificationResponse)
	assert.NoError(t, suite.jinmuHealth.ReadPushNotification(ctx, req, resp))
}
