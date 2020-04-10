package mysqldb

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// UserProfileTestSuite 是 UserProfile 的 testSuite
type UserProfileTestSuite struct {
	suite.Suite
	db *DbClient
}

// SetupSuite 准备设置 Test Suite 执行
func (suite *UserProfileTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	db, _ := newTestingDbClientFromEnvFile(envFilepath)
	suite.db = db
}

// TestModifyUserProfile 修改用户档案
func (suite *UserProfileTestSuite) TestModifyUserProfile() {
	t := suite.T()
	userID, _ := strconv.Atoi(os.Getenv("X_TEST_USER_ID"))
	height, _ := strconv.Atoi(os.Getenv("X_TEST_HEIGHT"))
	weight, _ := strconv.Atoi(os.Getenv("X_TEST_WEIGHT"))
	profile := &UserProfile{
		UserID:   int32(userID),
		Nickname: os.Getenv("X_TEST_NICKNAME"),
		Gender:   GenderMale,
		Birthday: time.Now().UTC(),
		Height:   int32(height),
		Weight:   int32(weight),
	}
	ctx := context.Background()
	err := suite.db.ModifyUserProfile(ctx, profile)
	assert.NoError(t, err)
}

// TestFindUserProfile 测试找到用户档案
func (suite *UserProfileTestSuite) TestFindUserProfile() {
	t := suite.T()
	userID, _ := strconv.Atoi(os.Getenv("X_TEST_USER_ID"))
	nickname := os.Getenv("X_TEST_NICKNAME")
	ctx := context.Background()
	u, err := suite.db.FindUserProfile(ctx, int32(userID))
	assert.NoError(t, err)
	assert.Equal(t, nickname, u.Nickname)
}

func TestUserProfileTestSuite(t *testing.T) {
	suite.Run(t, new(UserProfileTestSuite))
}
