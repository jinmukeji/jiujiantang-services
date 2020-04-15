package mysqldb

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// AuditUserCredentialUpdateTestSuite 是 审计记录 的单元测试的 Test Suite
type AuditUserCredentialUpdateTestSuite struct {
	suite.Suite
	db      *DbClient
	account *Account
}

type Account struct {
	account           string
	userID            int32
	seed              string
	oldHashedPassword string
	clientID          string
}

// SetupSuite 设置数据库连接
func (suite *AuditUserCredentialUpdateTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.db, _ = newTestingDbClientFromEnvFile(envFilepath)
	suite.account = newTestingAccountFromEnvFile(envFilepath)
}

// TearDownSuite 测试结束关闭连接
func (suite *AuditUserCredentialUpdateTestSuite) TearDownSuite() {
	safeCloseDB(suite.db)
}

// TestFindClientByClientID 测试 FindClientByClientID 成功返回记录
func (suite *AuditUserCredentialUpdateTestSuite) TestCreateAuditUserCredentialUpdate() {
	t := suite.T()
	record := &AuditUserCredentialUpdate{
		UserID:            suite.account.userID,
		ClientID:          suite.account.clientID,
		UpdatedRecordType: PasswordUpdated,
	}
	ctx := context.Background()
	err := suite.db.GetDB(ctx).CreateAuditUserCredentialUpdate(ctx, record)
	assert.NoError(t, err)
}

func newTestingAccountFromEnvFile(filepath string) *Account {
	_ = godotenv.Load(filepath)
	userID, _ := strconv.Atoi(os.Getenv("X_TEST_USER_ID"))
	return &Account{
		os.Getenv("X_TEST_ACCOUNT"),
		int32(userID),
		os.Getenv("X_TEST_SEED"),
		os.Getenv("X_TEST_HASHED_PASSWORD"),
		os.Getenv("X_TEST_CLIENT_ID"),
	}
}

// TestCompanyTestSuite 启动测试
func TestAuditUserCredentialUpdateTestSuite(t *testing.T) {
	suite.Run(t, new(AuditUserCredentialUpdateTestSuite))
}
