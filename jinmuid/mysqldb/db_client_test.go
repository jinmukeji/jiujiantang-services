package mysqldb

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/jinmukeji/go-pkg/mysqldb"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// DbClientTestSuite 是 DbClient 的单元测试的 Test Suite
type DbClientTestSuite struct {
	suite.Suite
	envFilepath string
}

const (
	enableLog = true
	maxConns  = 1
)

// newTestingDbClientFromEnvFile 从配置文件里面读取环境变量并创建 DbClient
func newTestingDbClientFromEnvFile(filepath string) (*DbClient, error) {
	err := godotenv.Load(filepath)
	if err != nil {
		log.Fatalf("Error loading %s file", filepath)
	}

	db, err := NewDbClient(
		mysqldb.Address(os.Getenv("X_DB_ADDRESS")),
		mysqldb.Username(os.Getenv("X_DB_USERNAME")),
		mysqldb.Password(os.Getenv("X_DB_PASSWORD")),
		mysqldb.Database(os.Getenv("X_DB_DATABASE")),
		mysqldb.EnableLog(enableLog),
		mysqldb.MaxConnections(maxConns),
	)
	return db, err
}

// SetupSuite 准备设置 Test Suite 执行
func (suite *DbClientTestSuite) SetupSuite() {
	// 设置环境变量配置文件地址
	suite.envFilepath = filepath.Join("testdata", "local.svc-jinmuid.env")
}

// TearDownSuite 结束 Test Suite 执行
func (suite *DbClientTestSuite) TearDownSuite() {
}

// safeCloseDB 安全的关闭数据库连接
func safeCloseDB(db *DbClient) {
	if db != nil {
		db.Close()
	}
}

// TestConnectDB 测试 ConnectDB 方法成功返回 *DbClient
func (suite *DbClientTestSuite) TestConnectDB() {
	t := suite.T()
	db, err := newTestingDbClientFromEnvFile(suite.envFilepath)
	defer safeCloseDB(db)

	assert.Nil(t, db)
	assert.Error(t, err)
}

func (suite *DbClientTestSuite) TestConnectDBMore() {
	t := suite.T()
	// db, err := newTestingDbClientFromEnvFile(suite.envFilepath)
	db1, err := newTestingDbClientFromEnvFile(suite.envFilepath)
	length := 50
	dbarray := make([]*DbClient, length)
	for idx := 0; idx < length; idx++ {
		dbarray[idx], err = newTestingDbClientFromEnvFile(suite.envFilepath)
		defer safeCloseDB(dbarray[idx])
	}
	// defer safeCloseDB(db)
	assert.NotNil(t, db1)
	assert.NoError(t, err)
}

// TestConnectDB 测试 ConnectDB 方法失败
func (suite *DbClientTestSuite) TestConnectDBFailed() {
	t := suite.T()
	wrongAddr := "localhost:3306"
	wrongUsername := "wront_username"
	wrongPwd := "this_is_a_wrong_password"

	db, err := NewDbClient(
		mysqldb.Address(wrongAddr),
		mysqldb.Username(wrongUsername),
		mysqldb.Password(wrongPwd),
	)
	defer safeCloseDB(db)

	assert.Nil(t, db)
	assert.Error(t, err)
}

// TestConnectDB 测试 DbClient 的 Options 方法失败
func (suite *DbClientTestSuite) TestDbClientOptions() {
	t := suite.T()
	db, err := newTestingDbClientFromEnvFile(suite.envFilepath)
	defer safeCloseDB(db)

	assert.NotNil(t, db)
	assert.NoError(t, err)

	opts := mysqldb.DefaultOptions()
	assert.Equal(t, os.Getenv("X_DB_ADDRESS"), opts.Address)
	assert.Equal(t, os.Getenv("X_DB_USERNAME"), opts.Username)
	assert.Equal(t, os.Getenv("X_DB_PASSWORD"), opts.Password)
	assert.Equal(t, os.Getenv("X_DB_DATABASE"), opts.Database)
	assert.Equal(t, enableLog, opts.EnableLog)
	assert.Equal(t, maxConns, opts.MaxConnections)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestDBTestSuite(t *testing.T) {
	suite.Run(t, new(DbClientTestSuite))
}