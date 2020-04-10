package mysqldb

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	enableLog = true
	maxConns  = 1
)

// FeedbackTestSuite 是 Feedback 的单元测试 Test Suite
type SmsRecordTestSuite struct {
	suite.Suite
	db *DbClient
}

// SetupSuite 设置测试数据库
func (suite *SmsRecordTestSuite) SetupSuite() {
	suite.db, _ = newTestingDbClientFromEnvFile("../../build/local.svc-sms-gw.env")
}

// TestCreateSmsRecord 测试 CreateSmsRecord
func (suite *SmsRecordTestSuite) TestCreateSmsRecord() {
	t := suite.T()
	now := time.Now()
	record := &SmsRecord{
		Phone:          "18805177594",
		SmsStatus:      SendSucceed,
		TemplateAction: "SignUp",
		PlatformType:   "Aliyun",
		TemplateParam:  "1234",
		NationCode:     "+86",
		Language:       SimpleChinese,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	assert.NoError(t, suite.db.CreateSmsRecord(record))
}

// TestSearchSmsRecordCountsIn24hours 测试 SearchSmsRecordCountsIn24hours
func (suite *SmsRecordTestSuite) TestSearchSmsRecordCountsIn24hours() {
	t := suite.T()
	count, err := suite.db.SearchSmsRecordCountsIn24hours("13700007474", "+86")
	assert.NoError(t, err)
	assert.Equal(t, 8, count)
}

// TestSmsRecordTestSuite 启动测试
func TestSmsRecordTestSuite(t *testing.T) {
	suite.Run(t, new(SmsRecordTestSuite))
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
