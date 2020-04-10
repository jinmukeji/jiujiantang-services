package mysqldb

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    context "golang.org/x/net/context"
)

// UserPreferencesTestSuite 是 UserPreferences 的 testSuite
type UserPreferencesTestSuite struct {
	suite.Suite
	db *DbClient
}

// SetupSuite 准备设置 Test Suite 执行
func (suite *UserPreferencesTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-biz-core.env")
	db, _ := newTestingDbClientFromEnvFile(envFilepath)
	suite.db = db
}

// TestGetUserPreferencesByUserID 测试 得到UserPreferences
func (suite *UserPreferencesTestSuite) TestGetUserPreferencesByUserID() {
	t := suite.T()
	ctx := context.Background()
	const userID = int32(1)
	u, err := suite.db.GetUserPreferencesByUserID(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, 1, u.EnableSyndromeDifferentiation)
}

func TestUserPreferencesTestSuite(t *testing.T) {
	suite.Run(t, new(UserPreferencesTestSuite))
}
