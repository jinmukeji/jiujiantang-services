package mysqldb

import (
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"math/rand"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	context "golang.org/x/net/context"
)

const UserRegisterTypeLegacy = "LEGACY"

// UserTestSuite 是 User 的 testSuite
type UserTestSuite struct {
	suite.Suite
	db *DbClient
}

// SetupSuite 准备设置 Test Suite 执行
func (suite *UserTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-biz-core.env")
	db, _ := newTestingDbClientFromEnvFile(envFilepath)
	suite.db = db
}

// TestCreateUserSuccess 测试 token 成功创建
func (suite *UserTestSuite) TestCreateUserSuccess() {
	t := suite.T()
	ctx := context.Background()
	now := time.Now()
	u := &User{
		Username:     "abcabc",
		Nickname:     "abcabc",
		Gender:       "M",
		Birthday:     now.AddDate(-20, 0, 0).UTC(),
		CreatedAt:    now.UTC(),
		RegisterType: UserRegisterTypeLegacy,
		RegisterTime: now.UTC(),
		UpdatedAt:    now.UTC(),
	}
	_, err := suite.db.CreateUser(ctx, u)
	assert.NoError(t, err)
}

// TestFindUserByUserIDSuccess 测试查找用户成功
func (suite *UserTestSuite) TestFindUserByUserIDSuccess() {
	t := suite.T()
	ctx := context.Background()
	u, err := suite.db.FindUserByUserID(ctx, 1)
	assert.NotNil(t, u)
	assert.NoError(t, err)
}

// TestUpdateUserProfile 测试修改用户个人信息成功
func (suite *UserTestSuite) TestUpdateUserProfile() {
	const userID = 1
	t := suite.T()
	ctx := context.Background()
	randName := strconv.Itoa(rand.Int())
	var p ProtoUserProfile
	err := suite.db.UpdateUserProfile(ctx, p, userID)
	assert.NoError(t, err)
	u, _ := suite.db.FindUserByUserID(ctx, userID)
	assert.Equal(t, randName, u.Nickname)
}

// TestTokenTestSuite 启动测试
func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
