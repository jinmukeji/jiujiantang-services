package mysqldb

import (
	"context"
	"time"
)

// UserSubscriptionSharing 订阅拥有或者被分享
type UserSubscriptionSharing struct {
	SubscriptionID int32      `gorm:"primary_key"`    // 订阅ID
	UserID         int32      `gorm:"column:user_id"` // 使用者
	CreatedAt      time.Time  // 创建时间
	UpdatedAt      time.Time  // 更新时间
	DeletedAt      *time.Time // 删除时间
}

// TableName 返回 UserSubscriptionSharing 所在的表名
func (s UserSubscriptionSharing) TableName() string {
	return "user_subscription_sharing"
}

// CreateMultiUserSubscriptionSharing 创建订阅多个拥有或者被分享记录
func (db *DbClient) CreateMultiUserSubscriptionSharing(ctx context.Context, userSubscriptionSharing []*UserSubscriptionSharing) error {
	tx := db.GetDB(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	for _, sharing := range userSubscriptionSharing {
		if err := tx.Create(sharing).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

// CreateUserSubscriptionSharing 创建单个订阅拥有或者被分享记录
func (db *DbClient) CreateUserSubscriptionSharing(ctx context.Context, UserSubscriptionSharing *UserSubscriptionSharing) error {
	return db.GetDB(ctx).Create(UserSubscriptionSharing).Error
}

// GetUserCountOfSubscription 获取订阅的使用者数量
func (db *DbClient) GetUserCountOfSubscription(ctx context.Context, subscriptionID int32) (int32, error) {
	var count int
	err := db.GetDB(ctx).Model(&UserSubscriptionSharing{}).Where("subscription_id = ? AND user_subscription_sharing.deleted_at is NULL", subscriptionID).Count(&count).Error
	return int32(count), err

}

// DeleteSubscriptionUsers 删除订阅下的用户
func (db *DbClient) DeleteSubscriptionUsers(ctx context.Context, userIDList []int32, subscriptionID int32) error {
	if err := db.GetDB(ctx).Delete(UserSubscriptionSharing{}, "subscription_id = ? AND user_id IN (?)", subscriptionID, userIDList).Error; err != nil {
		return err
	}
	return nil
}

// CheckExistUserSubscriptionSharing 检查user是否拥有或被分享Subscription
func (db *DbClient) CheckExistUserSubscriptionSharing(ctx context.Context, userID int32, subscriptionID int32) (bool, error) {
	var count int
	err := db.GetDB(ctx).Model(&UserSubscriptionSharing{}).Where("user_id = ? and subscription_id = ?", userID, subscriptionID).Count(&count).Error
	return count >= 1, err
}
