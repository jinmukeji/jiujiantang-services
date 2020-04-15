package mysqldb

import (
	"context"
	"time"
)

type SubscriptionType int32

const (
	// 定制化
	SubscriptionTypeCustomizedVersion = 0
	// 试用版
	SubscriptionTypeTrialVersion = 1
	// 黄喜马把脉
	SubscriptionTypeGoldenVersion = 2
	// 白喜马把脉
	SubscriptionTypePlatinumVersion = 3
	// 钻石姆
	SubscriptionTypeDiamondVersion = 4
	// 礼品版
	SubscriptionTypeGiftVersion = 5
)

// Subscription 订阅
type Subscription struct {
	SubscriptionID      int32            `gorm:"primary_key;column:subscription_id"` // 订阅ID
	SubscriptionType    SubscriptionType `gorm:"column:subscription_type"`           // 0 定制化 1 试用版 2 黄喜马把脉 3 白喜马把脉 4 钻石姆 5 礼品版
	MaxUserLimits       int32            `gorm:"column:max_user_limits"`             // 组织下最大用户数量
	Activated           bool             `gorm:"column:activated"`                   // 是否激活
	CustomizedCode      string           `gorm:"column:customized_code"`             // 自定义代码
	ActivatedAt         time.Time        `gorm:"column:activated_at"`                // 合同开始日期
	ExpiredAt           time.Time        `gorm:"column:expired_at"`                  // 合同结束日期
	ContractYear        int32            `gorm:"column:contract_year"`               // 合同期限
	OwnerID             int32            `gorm:"column:owner_id"`                    // 拥有者
	IsSelected          bool             `gorm:"column:is_selected"`                 // 是否被选择
	IsMigratedActivated bool             `gorm:"column:is_migrated_activated"`       // 迁移前激活状态
	CreatedAt           time.Time        // 创建时间
	UpdatedAt           time.Time        // 更新时间
	DeletedAt           *time.Time       // 删除时间
}

// TableName 返回 Subscription 所在的表名
func (s Subscription) TableName() string {
	return "subscription"
}

// CreateSubscription 创建订阅
func (db *DbClient) CreateSubscription(ctx context.Context, subscription *Subscription) (*Subscription, error) {
	err := db.GetDB(ctx).Create(subscription).Error
	if err != nil {
		return nil, err
	}
	return subscription, err
}

// GetSubscriptionsByUserID 通过userID获取订阅
func (db *DbClient) GetSubscriptionsByUserID(ctx context.Context, userID int32) ([]*Subscription, error) {
	var s []*Subscription
	if err := db.GetDB(ctx).Raw(`SELECT 
	    (S.subscription_id), 
		S.subscription_type, 
		S.max_user_limits,
		S.activated, 
		S.customized_code, 
		S.activated_at,
		S.expired_at, 
		S.contract_year, 
        S.owner_id, 
		S.is_selected,
		S.is_migrated_activated,
		S.created_at,
		S.updated_at,
		S.deleted_at
		FROM user_subscription_sharing as US 
		inner join subscription as S on S.subscription_id = US.subscription_id
		AND US.user_id = ? AND S.deleted_at IS NULL where US.deleted_at IS NULL`, userID).Scan(&s).Error; err != nil {
		return nil, err
	}
	return s, nil
}

// GetSelectedSubscriptionByUserID 通过userID获取订阅
func (db *DbClient) GetSelectedSubscriptionByUserID(ctx context.Context, userID int32) (*Subscription, error) {
	var s Subscription
	err := db.GetDB(ctx).Raw(`SELECT 
    S.subscription_id,
    S.subscription_type,
    S.max_user_limits,
    S.activated,
    S.customized_code,
    S.activated_at,
    S.expired_at,
    S.contract_year,
    S.owner_id,
    S.is_selected,
    S.is_migrated_activated,
    S.created_at,
    S.updated_at,
    S.deleted_at
FROM
    user_subscription_sharing AS US
        INNER JOIN
    subscription AS S ON S.subscription_id = US.subscription_id
        AND S.deleted_at IS NULL AND S.is_selected = '1'
WHERE
	US.user_id = ? AND US.deleted_at IS NULL`, userID).Scan(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// FindSelectedSubscriptionByUserID 通过用户ID找订阅
func (db *DbClient) FindSelectedSubscriptionByUserID(ctx context.Context, userID int32) (*Subscription, error) {
	var s Subscription
	if err := db.GetDB(ctx).First(&s, "owner_id = ? and is_selected = 1", userID).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

// ActivateSubscription 激活订阅
func (db *DbClient) ActivateSubscription(ctx context.Context, subscription *Subscription) error {
	return db.GetDB(ctx).Model(&Subscription{}).Where("subscription_id = ?", subscription.SubscriptionID).Updates(map[string]interface{}{
		"expired_at":   subscription.ExpiredAt,
		"activated":    subscription.Activated,
		"activated_at": subscription.ActivatedAt,
		"updated_at":   subscription.UpdatedAt,
	}).Error
}

// CheckSubscriptionOwner 检查用户是否是订阅的拥有者
func (db *DbClient) CheckSubscriptionOwner(ctx context.Context, ownerID, subscriptionID int) (bool, error) {
	var count int
	if err := db.GetDB(ctx).Raw("select count(*) from subscription where owner_id = ? AND subscription_id = ?", int32(ownerID), int32(subscriptionID)).Count(&count).Error; err != nil {
		return false, err
	}
	return count != 0, nil
}

// UpdateSubscriptionIsSelectedStatus 更新订阅的默认状态
func (db *DbClient) UpdateSubscriptionIsSelectedStatus(ctx context.Context, ownerID int32) error {
	return db.GetDB(ctx).Model(&Subscription{}).Where("owner_id = ?", ownerID).Updates(map[string]interface{}{
		"is_selected": false,
	}).Error
}
