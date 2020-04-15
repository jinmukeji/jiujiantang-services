package mysqldb

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// UserPreferencesTestSuite 是 UserPreferences 的 testSuite
type UserPreferencesTestSuite struct {
	suite.Suite
	db *DbClient
}

// SetupSuite 准备设置 Test Suite 执行
func (suite *UserPreferencesTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	db, _ := newTestingDbClientFromEnvFile(envFilepath)
	suite.db = db
}

// TestGetUserPreferencesByUserID 测试 得到UserPreferences
func (suite *UserPreferencesTestSuite) TestGetUserPreferencesByUserID() {
	t := suite.T()
	userID, _ := strconv.Atoi(os.Getenv("X_TEST_USER_ID"))
	ctx := context.Background()
	u, err := suite.db.GetDB(ctx).GetUserPreferencesByUserID(ctx, int32(userID))
	fmt.Println(u)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), u.EnableSyndromeDifferentiation) // 是否开启中医脏腑判读
}

func TestUserPreferencesTestSuite(t *testing.T) {
	suite.Run(t, new(UserPreferencesTestSuite))
}
