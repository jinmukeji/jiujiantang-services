package mysqldb

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	context "golang.org/x/net/context"
)

// RecordTestSuite 是 Record 的单元测试的 Test Suite
type RecordTestSuite struct {
	suite.Suite
	db *DbClient
}

// SetupSuite 准备设置 Test Suite 执行
func (suite *RecordTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-biz-core.env")
	db, _ := newTestingDbClientFromEnvFile(envFilepath)
	suite.db = db
}

// TearDownSuite 结束 Test Suite 执行
func (suite *RecordTestSuite) TearDownSuite() {
	safeCloseDB(suite.db)
}

// TestFindValidRecordByID 测试 FindValidRecordByID 方法成功返回 Record 记录
func (suite *RecordTestSuite) TestFindValidRecordByID() {
	t := suite.T()
	ctx := context.Background()
	db := suite.db
	const testRecordID = 1
	record, err := db.GetDB(ctx).FindValidRecordByID(ctx, testRecordID)
	assert.NoError(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, testRecordID, record.RecordID)
}

// TestFindValidRecordByID 测试 FindValidRecordByID 方法返回 Error
func (suite *RecordTestSuite) TestFindValidRecordByIDFailed() {
	t := suite.T()
	db := suite.db
	ctx := context.Background()
	notExistedRecordID := 0
	record, err := db.GetDB(ctx).FindValidRecordByID(ctx, notExistedRecordID)

	assert.Error(t, err)
	assert.Nil(t, record)
}

// TestUpdateCommentByRecordID 测试更新 comment 成功
func (suite *RecordTestSuite) TestUpdateCommentByRecordID() {
	t := suite.T()
	ctx := context.Background()
	testRecordID := 1
	randomComment := uuid.New().String()
	assert.NoError(t, suite.db.GetDB(ctx).UpdateRemarkByRecordID(ctx, testRecordID, randomComment))
	updatedRecord, _ := suite.db.GetDB(ctx).FindValidRecordByID(ctx, testRecordID)
	assert.Equal(t, randomComment, updatedRecord.Remark)
}

// TestCreateRecord 测试插入一条 record 成功
func (suite *RecordTestSuite) TestCreateRecord() {
	t := suite.T()
	ctx := context.Background()
	now := time.Now()
	record := &Record{
		C0:        32,
		CreatedAt: now.UTC(),
		UpdatedAt: now.UTC(),
		IsValid:   1,
	}
	assert.NoError(t, suite.db.GetDB(ctx).CreateRecord(ctx, record))
	assert.NotEqual(t, 0, record.RecordID)
	createdRecord, _ := suite.db.GetDB(ctx).FindValidRecordByID(ctx, record.RecordID)
	assert.Equal(t, record.C0, createdRecord.C0)
}

// TestCheckUserRecordAssociated 测试验证帐号和测量记录是否有关联
func (suite *RecordTestSuite) TestCheckUserRecordAssociated() {
	t := suite.T()
	ctx := context.Background()
	var testUserID int32 = 1
	var testRecordID int32 = 1
	res, err := suite.db.GetDB(ctx).CheckUserHasRecord(ctx, testUserID, testRecordID)
	assert.NoError(t, err)
	assert.True(t, res)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestRecordTestSuite(t *testing.T) {
	suite.Run(t, new(RecordTestSuite))
}
