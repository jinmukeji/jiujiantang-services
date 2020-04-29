package handler

import (
	"context"
	"path/filepath"
	"testing"

	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
)

// LocalNotificationTestSuite 是本地通知单元测试的 Test Suite
type LocalNotificationTestSuite struct {
	suite.Suite
	jinmuHealth *JinmuHealth
}

// SetupSuite 设置测试环境
func (suite *LocalNotificationTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-biz-core.env")
	suite.jinmuHealth = newTestingJinmuHealthFromEnvFile(envFilepath)
	suite.jinmuHealth.datastore, _ = newTestingDbClientFromEnvFile(envFilepath)
	suite.jinmuHealth.mailClient, _ = newTestingMailClientFromEnvFile(envFilepath)
}

// TestClientAuth 测试客户端授权
func (suite *LocalNotificationTestSuite) TestCreateLocalNotification() {
	t := suite.T()
	ctx := context.Background()
	req := new(proto.CreateLocalNotificationRequest)
	resp := new(proto.CreateLocalNotificationResponse)
	req.LocalNotification = &proto.LocalNotification{
		Content: "喜马宝宝提醒您，又到检查时间啦！",
		Schedule: &proto.Schedule{
			EventHappenAt: "2018-10-10T20:00:00",
			Timezone:      "local",
			Repeat: &proto.RepeatSchedule{
				Frequency: proto.Frequency_FREQUENCY_DAILY,
				Interval:  1,
			},
		},
	}
	assert.NoError(t, suite.jinmuHealth.CreateLocalNotification(ctx, req, resp))
}

func (suite *LocalNotificationTestSuite) GetLocalNotifications() {
	t := suite.T()
	ctx := context.Background()
	req := new(proto.GetLocalNotificationsRequest)
	resp := new(proto.GetLocalNotificationsResponse)
	assert.NoError(t, suite.jinmuHealth.GetLocalNotifications(ctx, req, resp))
}

func TestCreateLocalNotificationTestSuite(t *testing.T) {
	suite.Run(t, new(LocalNotificationTestSuite))
}
