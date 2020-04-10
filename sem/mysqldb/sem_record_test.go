package mysqldb

import (
	"log"
	"os"
	"testing"
	"time"

	encry "github.com/jinmukeji/go-pkg/crypto/rand"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	enableLog = true
	maxConns  = 1
	digits    = "0123456789"
	length    = 6
)

// FeedbackTestSuite 是 Feedback 的单元测试 Test Suite
type SemRecordTestSuite struct {
	suite.Suite
	db *DbClient
}

// SetupSuite 设置测试数据库
func (suite *SemRecordTestSuite) SetupSuite() {
	suite.db, _ = newTestingDbClientFromEnvFile("../../build/local.svc-sem-gw.env")
}

// TestCreateSemRecord 测试 CreateSemRecord
func (suite *SemRecordTestSuite) TestCreateSemRecord() {
	t := suite.T()
	now := time.Now()
	code := encry.RandomStringWithMask(digits, length)
	log.Println("this is code", code)
	record := &SemRecord{
		ToAddress:      "tech@jinmuhealth.com",
		SemStatus:      SendSucceed,
		TemplateAction: "",
		PlatformType:   "Aliyun",
		TemplateParam:  `{"code": code }`,
		Language:       SimpleChinese,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	assert.NoError(t, suite.db.CreateSemRecord(record))
}

// TestSemRecordTestSuite 启动测试
func TestSemRecordTestSuite(t *testing.T) {
	suite.Run(t, new(SemRecordTestSuite))
}

// newTestingDbClientFromEnvFile 从配置文件里面读取环境变量并创建 DbClient
func newTestingDbClientFromEnvFile(filepath string) (*DbClient, error) {
	err := godotenv.Load(filepath)
	if err != nil {
		log.Fatalf("Error loading %s file", filepath)
	}

	db, err := NewDbClient(
		Address(os.Getenv("X_DB_ADDRESS")),
		Username(os.Getenv("X_DB_USERNAME")),
		Password(os.Getenv("X_DB_PASSWORD")),
		Database(os.Getenv("X_DB_DATABASE")),
		EnableLog(enableLog),
		MaxConnections(maxConns),
	)
	return db, err
}
