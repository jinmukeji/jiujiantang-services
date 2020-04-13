package mysqldb

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
	context "golang.org/x/net/context"
)

// AccountLRecordSuite 是 AccountLRecord 的单元测试的 Test Suite
type AccountLRecordSuite struct {
	suite.Suite
	db *DbClient
}

// SetupSuite 设置数据库连接
func (suite *AccountLRecordSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.env")
	suite.db, _ = newTestingDbClientFromEnvFile(envFilepath)
}

// TestCreateAccountLRecord 测试创建一体机账户与记录关联表
func (suite *AccountLRecordSuite) TestCreateAccountLRecord() {
	t := suite.T()
	ctx := context.Background()
	account := "JML003"
	recordID := 10000000
	err := suite.db.GetDB(ctx).CreateAccountLRecord(ctx, account, int32(recordID))
	assert.NoError(t, err)
}
func TestAccountLRecordSuite(t *testing.T) {
	suite.Run(t, new(AccountLRecordSuite))
}
