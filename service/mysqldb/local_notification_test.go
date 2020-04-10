package mysqldb

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	context "golang.org/x/net/context"
)

// generateLocalNotification 生成随机测试用例
func generateLocalNotification() *LocalNotification {
	now := time.Now()
	return &LocalNotification{
		Title:         "",
		Content:       "喜马宝宝提醒您，又到检查时间啦！",
		EventHappenAt: "2018-01-01T20:00:00",
		Timezone:      "local",
		Frequency:     "FrequencyDaily",
		Interval:      1,
		CreatedAt:     now.UTC(),
		UpdatedAt:     now.UTC(),
	}
}

// LocalNotificationTestSuite 是 LocalNotification 的单元测试 Test Suite
type LocalNotificationTestSuite struct {
	suite.Suite
	db *DbClient
}

// SetupSuite 设置测试数据库
func (suite *LocalNotificationTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-biz-core.env")
	suite.db, _ = newTestingDbClientFromEnvFile(envFilepath)
}

// TestCreateLocalNotification 测试创建一个本地消息
func (suite *LocalNotificationTestSuite) TestCreateLocalNotification() {
	t := suite.T()
	localnotification := generateLocalNotification()
	ctx := context.Background()
	assert.NoError(t, suite.db.CreateLocalNotification(ctx, localnotification))
	assert.NotZero(t, localnotification.Timezone)
}

// TestGetUserPreferencesByUserID 测试 得到LocalNotification
func (suite *LocalNotificationTestSuite) TestGetLocalNotification() {
	t := suite.T()
	ctx := context.Background()
	_, err := suite.db.GetLocalNotifications(ctx)
	assert.NoError(t, err)
}

// TestDeleteLocalNotification 测试 TestDeleteLocalNotification
func (suite *LocalNotificationTestSuite) TestDeleteLocalNotification() {
	t := suite.T()
	ctx := context.Background()
	lnID := 1
	err := suite.db.DeleteLocalNotification(ctx, lnID)
	assert.NoError(t, err)
}

func TestLocalNotificationTestSuite(t *testing.T) {
	suite.Run(t, new(LocalNotificationTestSuite))
}
