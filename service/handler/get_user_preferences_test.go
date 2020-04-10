package handler

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// UserPreferencesTestSuite 是 UserPreferences 的 testSuite
type UserPreferencesTestSuite struct {
	suite.Suite
	jinmuHealth *JinmuHealth
}

// SetupSuite 设置测试环境
func (suite *UserPreferencesTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-biz-core.env")
	suite.jinmuHealth = newTestingJinmuHealthFromEnvFile(envFilepath)
	suite.jinmuHealth.datastore, _ = newTestingDbClientFromEnvFile(envFilepath)
	suite.jinmuHealth.mailClient, _ = newTestingMailClientFromEnvFile(envFilepath)
}

// TestGetUserPreferencesByUserID 测试 得到UserPreferences
func (suite *UserPreferencesTestSuite) TestGetUserPreferencesByUserID() {
	t := suite.T()
	const userID = int32(1)
	ctx := context.Background()
	u, err := suite.jinmuHealth.datastore.GetUserPreferencesByUserID(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, 1, u.EnableSyndromeDifferentiation)
}

func TestUserPreferencesTestSuite(t *testing.T) {
	suite.Run(t, new(UserPreferencesTestSuite))
}
