package mysqldb

import (
	"context"
	"time"
)

// NotificationPreferences 通知配置首选项
type NotificationPreferences struct {
	UserID                 int32      `gorm:"primary_key"`                      // 用户 id
	PhoneEnabled           bool       `gorm:"column:phone_enabled"`             // 是否允许手机通知
	PhoneEnabledUpdatedAt  time.Time  `gorm:"column:phone_enabled_updated_at"`  // 最新更新是否允许手机通知状态的时间
	WechatEnabled          bool       `gorm:"column:wechat_enabled"`            // 是否允许微信通知
	WechatEnabledUpdatedAt time.Time  `gorm:"column:wechat_enabled_updated_at"` // 最新更新是否允许微信通知状态的时间
	WeiboEnabled           bool       `gorm:"column:weibo_enabled"`             // 是否允许微博通知
	WeiboEnabledUpdatedAt  time.Time  `gorm:"column:weibo_enabled_updated_at"`  // 最新更新是否允许微博通知状态的时间
	CreatedAt              time.Time  // 创建时间
	UpdatedAt              time.Time  // 更新时间
	DeletedAt              *time.Time // 删除时间
}

// CreateNotificationPreferences 新增通知配置首选项
func (db *DbClient) CreateNotificationPreferences(ctx context.Context, notificationPreferences *NotificationPreferences) error {
	return db.DB(ctx).Create(notificationPreferences).Error
}

// HasUserSetNotificationPreferences 查看用户是否设置通知配置首选项
func (db *DbClient) HasUserSetNotificationPreferences(ctx context.Context, userID int32) (bool, error) {
	var count int
	err := db.DB(ctx).Model(&NotificationPreferences{}).Where("user_id = ?", userID).Count(&count).Error
	return count == 1, err
}

// GetNotificationPreferences 获取通知配置首选项
func (db *DbClient) GetNotificationPreferences(ctx context.Context, UserID int32) (*NotificationPreferences, error) {
	var notificationPreferences NotificationPreferences
	if err := db.DB(ctx).First(&notificationPreferences, "( user_id = ? AND deleted_at IS NULL) ", UserID).Error; err != nil {
		return nil, err
	}
	return &notificationPreferences, nil
}

// UpdateNotificationPreferences 更新通知配置首选项
func (db *DbClient) UpdateNotificationPreferences(ctx context.Context, notificationPreferences *NotificationPreferences) error {
	return db.DB(ctx).Model(&NotificationPreferences{}).Where("user_id = ?", notificationPreferences.UserID).Updates(map[string]interface{}{
		"phone_enabled":             notificationPreferences.PhoneEnabled,
		"wechat_enabled":            notificationPreferences.WechatEnabled,
		"weibo_enabled":             notificationPreferences.WeiboEnabled,
		"phone_enabled_updated_at":  notificationPreferences.PhoneEnabledUpdatedAt,
		"wechat_enabled_updated_at": notificationPreferences.WechatEnabledUpdatedAt,
		"weibo_enabled_updated_at":  notificationPreferences.WeiboEnabledUpdatedAt,
		"updated_at":                time.Now().UTC(),
	}).Error
}
