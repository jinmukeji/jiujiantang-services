package mysqldb

import (
	"context"
	"time"
)

// Feedback 用户反馈
type Feedback struct {
	FeedbackID int        `gorm:"primary_key"`    // 用户反馈唯一标识
	UserID     int32      `gorm:"column:user_id"` // 用户ID
	ContactWay string     // 联系方式
	Content    string     // 意见内容
	IsValid    int        // 是否有效
	CreatedAt  time.Time  // 创建时间
	UpdatedAt  time.Time  // 更新时间
	DeletedAt  *time.Time // 删除时间
}

// TableName 返回 FeedBack 所在表名
func (feedback *Feedback) TableName() string {
	return "feedback"
}

// CreateFeedback 新增一个用户反馈
func (db *DbClient) CreateFeedback(ctx context.Context, feedback *Feedback) error {
	return db.Create(feedback).Error
}

// FindFeedbackByFeedBackID 查找一个用户反馈
func (db *DbClient) FindFeedbackByFeedBackID(ctx context.Context, feedbackID int) (*Feedback, error) {
	var feedback Feedback
	if err := db.First(&feedback, "feedback_id = ?", feedbackID).Error; err != nil {
		return nil, err
	}
	return &feedback, nil
}
