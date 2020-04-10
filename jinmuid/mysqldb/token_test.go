package mysqldb

import (
	"path/filepath"
	"testing"
	"time"
	"context"
	"github.com/jinmukeji/gf-api2/service/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// TokenTestSuite 是 Token 的 testSuite
type TokenTestSuite struct {
	suite.Suite
	db *DbClient
}

// SetupSuite 准备设置 Test Suite 执行
func (suite *TokenTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	db, _ := newTestingDbClientFromEnvFile(envFilepath)
	suite.db = db
}

// TestCreateTokenSuccess 测试 token 成功创建
func (suite *TokenTestSuite) TestCreateTokenSuccess() {
	t := suite.T()
	const userID = int32(1)
	token := auth.GenerateToken()
	ctx := context.Background()
	// 创建 token
	tk, err := suite.db.CreateToken(ctx, token, userID, time.Hour*12)
	assert.NoError(t, err)
	assert.Equal(t, userID, tk.UserID)
}

// TestFindUserIDByTokenSuccess 测试成功从 token 获取 userID
func (suite *TokenTestSuite) TestFindUserIDByTokenSuccess() {
	t := suite.T()
	const userID = int32(1)
	token := auth.GenerateToken()
	ctx := context.Background()
	// 创建 token
	tk, err := suite.db.CreateToken(ctx, token, userID, time.Hour*12)
	assert.NoError(t, err)

	// 从 token 取出 accoount
	user, _ := suite.db.FindUserIDByToken(ctx, tk.Token)
	assert.Equal(t, user, userID)
}

// TestFindUserIDByTokenFail 测试 token 超时 返回空的account
func (suite *TokenTestSuite) TestFindUserIDByTokenFail() {
	t := suite.T()
	token := auth.GenerateToken()
    ctx := context.Background()
	// 插入一条 1 秒后失效的记录
    _, errCreateToken := suite.db.CreateToken(ctx, token, 1, time.Second)
    assert.Error(t, errCreateToken)

	// 休眠2秒
	time.Sleep(time.Second * 2)
	userID, err := suite.db.FindUserIDByToken(ctx, token)

	// 失效返回error
	assert.True(t, userID == int32(0))
	assert.Error(t, err)
}

// TestDeleteToken 测试删除 token
func (suite *TokenTestSuite) TestDeleteToken() {
	t := suite.T()
	const userID = int32(1)
	token := auth.GenerateToken()
	ctx := context.Background()
	// 生成 token
	tk, err := suite.db.CreateToken(ctx, token, userID, time.Hour*12)
	assert.NoError(t, err)
	assert.Equal(t, userID, tk.UserID)

	// 删除 token
	err = suite.db.DeleteToken(ctx, tk.Token)
	assert.NoError(t, err)

	// 删除后找不到 token
	user, err := suite.db.FindUserIDByToken(ctx, tk.Token)
	assert.Equal(t, int32(0), user)
	assert.Error(t, err)
}

// TestTokenTestSuite 启动测试
func TestTokenTestSuite(t *testing.T) {
	suite.Run(t, new(TokenTestSuite))
}
