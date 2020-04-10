package mysqldb

import (
	"context"
	"time"
)

// Subscription 订阅
type Subscription struct {
	SubscriptionID   int        `gorm:"primary_key"` // 订阅ID
	OrganizationID   int        // 组织ID
	SubscriptionType int        // 0 定制化 1 试用版 2 黄金姆 3 白金姆 4 钻石姆 5 礼品版
	MaxUserLimits    int        // 组织下最大用户数量
	Active           int        // 是否激活
	CustomizedCode   string     // 自定义代码
	ActivatedAt      time.Time  // 合同开始日期
	ExpiredAt        time.Time  // 合同结束日期
	ContractYear     int        // 合同期限
	CreatedAt        time.Time  // 创建时间
	UpdatedAt        time.Time  // 更新时间
	DeletedAt        *time.Time // 删除时间
}

// TableName 返回 Subscription 所在的表名
func (s Subscription) TableName() string {
	return "subscription"
}

// CreateSubscription 创建订阅
func (db *DbClient) CreateSubscription(ctx context.Context, s *Subscription) error {
	return db.Create(s).Error
}

// FindSubscriptionsByOrganizationID 查找 Subscription 通过 OrganizationID
func (db *DbClient) FindSubscriptionsByOrganizationID(ctx context.Context, organizationID int) ([]*Subscription, error) {
	var subscriptions []*Subscription
	db.Raw(`SELECT 
	S.subscription_id,
	S.subscription_type,
	S.max_user_limits,
	S.activated as active,
	S.customized_code,
	S.activated_at,
	S.expired_at,
    S.contract_year,
	S.created_at,
	S.updated_at,
	S.deleted_at 
	FROM subscription as S 
	inner join organization_owner as OO on OO.owner_id = S.owner_id
	where OO.organization_id = ?`, organizationID).Scan(&subscriptions)
	return subscriptions, nil
}

// ActivateSubscription 更新订阅
func (db *DbClient) ActivateSubscription(ctx context.Context, s *Subscription) error {
	return db.Model(&Subscription{}).Where("subscription_id = ?", s.SubscriptionID).Update(map[string]interface{}{
		"active":       s.Active,
		"activated_at": s.ActivatedAt,
		"expired_at":   s.ExpiredAt,
		"updated_at":   s.UpdatedAt,
	}).Error
}
