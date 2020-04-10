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

// generateRandomFeedback 生成内容随机的用户反馈
func generateRandomFeedback(userID int32) *Feedback {
	now := time.Now()
	return &Feedback{
		UserID:     userID,
		ContactWay: uuid.New().String(),
		Content:    uuid.New().String(),
		IsValid:    1,
		CreatedAt:  now.UTC(),
		UpdatedAt:  now.UTC(),
	}
}

// FeedbackTestSuite 是 Feedback 的单元测试 Test Suite
type FeedbackTestSuite struct {
	suite.Suite
	db *DbClient
}

// SetupSuite 设置测试数据库
func (suite *FeedbackTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-biz-core.env")
	suite.db, _ = newTestingDbClientFromEnvFile(envFilepath)
}

// TestCreateFeedback 测试创建一个用户反馈
func (suite *FeedbackTestSuite) TestCreateFeedback() {
	const userID = 1
	t := suite.T()
	ctx := context.Background()
	feedback := generateRandomFeedback(userID)
	assert.NoError(t, suite.db.CreateFeedback(ctx, feedback))
	assert.NotZero(t, feedback.FeedbackID)
	feedbackFound, _ := suite.db.FindFeedbackByFeedBackID(ctx, feedback.FeedbackID)
	assert.Equal(t, feedback.Content, feedbackFound.Content)
}

func TestFeedBackTestSuite(t *testing.T) {
	suite.Run(t, new(FeedbackTestSuite))
}
