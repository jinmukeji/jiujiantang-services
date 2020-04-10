package mysqldb

import (
	"os"
	"path/filepath"
	"testing"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ClientTestSuite 是 Client 的单元测试的 Test Suite
type ClientTestSuite struct {
	suite.Suite
	db *DbClient
}

// SetupSuite 设置数据库连接
func (suite *ClientTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.db, _ = newTestingDbClientFromEnvFile(envFilepath)
}

// TearDownSuite 测试结束关闭连接
func (suite *ClientTestSuite) TearDownSuite() {
	safeCloseDB(suite.db)
}

// TestFindClientByClientID 测试 FindClientByClientID 成功返回记录
func (suite *ClientTestSuite) TestFindClientByClientID() {
	t := suite.T()
	ctx := context.Background()
	client, err := suite.db.FindClientByClientID(ctx, os.Getenv("X_TEST_CLIENT_ID"))
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "CHSnHWkepLThkmPw8IUX", client.SecretKey)
}

// TestCompanyTestSuite 启动测试
func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}
