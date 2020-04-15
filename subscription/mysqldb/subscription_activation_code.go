package mysqldb

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
)

// SubscriptionActivationCode 激活码
type SubscriptionActivationCode struct {
	Code             string           `gorm:"primary_key"`              // 激活码
	UserID           int32            `gorm:"colusmn:user_id"`          // 用户ID
	SubscriptionID   int32            `gorm:"column:subscription_id"`   // 订阅ID
	MaxUserLimits    int32            `gorm:"column:max_user_limits"`   // 组织下最大用户数量
	ContractYear     int32            `gorm:"column:contract_year"`     // 年限
	SubscriptionType SubscriptionType `gorm:"column:subscription_type"` // 0 定制化 1 试用版 2 黄喜马把脉 3 白喜马把脉 4 钻石姆 5 礼品版
	Checksum         string           `gorm:"column:checksum"`          // 校验位
	Activated        bool             `gorm:"column:activated"`         // 是否激活
	ActivatedAt      *time.Time       `gorm:"column:activated_at"`      // 激活时间
	Sold             bool             `gorm:"column:sold"`              // 是否售出
	SoldAt           *time.Time       `gorm:"column:sold_at"`           // 售出时间
	ExpiredAt        *time.Time       `gorm:"column:expired_at"`        // 到期时间
	CreatedAt        time.Time        // 创建时间
	UpdatedAt        time.Time        // 更新时间
	DeletedAt        *time.Time       // 删除时间
}

// TableName 返回 SubscriptionActivationCode 所在的表名
func (s SubscriptionActivationCode) TableName() string {
	return "subscription_activation_code"
}

// GetSubscriptionActivationCodeInfo 得到激活码的信息
func (db *DbClient) GetSubscriptionActivationCodeInfo(ctx context.Context, code string) (*SubscriptionActivationCode, error) {
	var subscriptionActivationCode SubscriptionActivationCode
	err := db.GetDB(ctx).Model(&SubscriptionActivationCode{}).Where("code = ?", code).Scan(&subscriptionActivationCode).Error
	if err != nil {
		return nil, err
	}
	return &subscriptionActivationCode, nil
}

// ActivateSubscriptionActivationCode 设置激活码为已经激活的
func (db *DbClient) ActivateSubscriptionActivationCode(ctx context.Context, code string, subscriptionID, userID int32) error {
	return db.GetDB(ctx).Model(&SubscriptionActivationCode{}).Where("code = ?", code).Updates(map[string]interface{}{
		"subscription_id": subscriptionID,
		"user_id":         userID,
		"activated":       true,
		"activated_at":    time.Now().UTC(),
		"updated_at":      time.Now().UTC(),
	}).Error
}

// UseSubscriptionActivationCode 使用订阅激活码
func (db *DbClient) UseSubscriptionActivationCode(ctx context.Context, userID int32, activationCode *SubscriptionActivationCode, uuid string) error {
	tx := db.GetDB(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	// 加锁
	errSetLock := tx.Exec(`update subscription_activation_code set activation_lock = ? where code = ? and activated = 0 and activation_lock is null`, uuid, activationCode.Code).Error
	if errSetLock != nil {
		tx.Rollback()
		return errSetLock
	}
	count, errGetLastAffectedRowCount := getLastAffectedRowCount(tx)
	// 有错或者影响不是1行就回滚
	if errGetLastAffectedRowCount != nil || count != 1 {
		tx.Rollback()
		return errGetLastAffectedRowCount
	}
	// 更新激活码状态
	errUpdateSubscriptionActivationCodeActivated := tx.Exec(`update subscription_activation_code set activated = 1 , activated_at = ? where code = ? and activation_lock = ?`, time.Now().UTC(), activationCode.Code, uuid).Error
	if errUpdateSubscriptionActivationCodeActivated != nil {
		tx.Rollback()
		return errUpdateSubscriptionActivationCodeActivated
	}
	count, errGetLastAffectedRowCount = getLastAffectedRowCount(tx)
	// 有错或者影响不是1行就回滚
	if errGetLastAffectedRowCount != nil || count != 1 {
		tx.Rollback()
		return errGetLastAffectedRowCount
	}
	// 正在使用的老订阅的数量
	var oidSubscriptionCount int32
	errOidSubscriptionCount := tx.Model(&Subscription{}).Where("owner_id = ? and is_selected = 1 and deleted_at is null", userID).Count(&oidSubscriptionCount).Error
	if errOidSubscriptionCount != nil {
		tx.Rollback()
		return errOidSubscriptionCount
	}
	var oldSubscription Subscription
	if oidSubscriptionCount > 0 {
		if errFindOidSubscription := tx.Where("is_selected = 1 and deleted_at is null").First(&oldSubscription, "owner_id = ?", userID).Error; errFindOidSubscription != nil {
			tx.Rollback()
			return errFindOidSubscription
		}
	}

	// 更新已经存在订阅的状态
	errUpdateSubscriptionStatus := tx.Exec(`update subscription set is_selected = 0 where owner_id = ?`, userID).Error
	if errUpdateSubscriptionStatus != nil {
		tx.Rollback()
		return errUpdateSubscriptionStatus
	}
	now := time.Now()
	// 生成订阅
	subscription := &Subscription{
		SubscriptionType: activationCode.SubscriptionType,
		Activated:        true,
		ActivatedAt:      now,
		MaxUserLimits:    activationCode.MaxUserLimits,
		ContractYear:     activationCode.ContractYear,
		ExpiredAt:        now.AddDate(int(activationCode.ContractYear), 0, 0),
		OwnerID:          userID,
		CreatedAt:        now,
		UpdatedAt:        now,
		IsSelected:       true,
	}
	errCreateSubscription := tx.Create(subscription).Error
	if errCreateSubscription != nil {
		tx.Rollback()
		return errCreateSubscription
	}
	// 没有旧的订阅就创建UserSubscriptionSharing表,有就的订阅迁移UserSubscriptionSharing表
	if oidSubscriptionCount == 0 {
		// 创建分享记录
		now := time.Now()
		errCreateUserSubscriptionSharing := tx.Create(&UserSubscriptionSharing{
			SubscriptionID: subscription.SubscriptionID,
			UserID:         userID,
			CreatedAt:      now.UTC(),
			UpdatedAt:      now.UTC(),
		}).Error
		if errCreateUserSubscriptionSharing != nil {
			tx.Rollback()
			return errCreateUserSubscriptionSharing
		}
	} else {
		// 迁移分享表的数据
		errBulkInsertSubscriptionSharing := tx.Exec(`insert  INTO  user_subscription_sharing
			(`+"`subscription_id`,`user_id`,`created_at`,`updated_at`,`deleted_at`"+`) 
			SELECT ?, US.user_id, ?, ?, NULL AS deleted_at 
			FROM user_subscription_sharing AS US 
			INNER JOIN subscription AS S on US.subscription_id = S.subscription_id AND S.subscription_id = ? AND S.deleted_at IS NULL
			WHERE US.deleted_at IS NULL`, subscription.SubscriptionID, now, now, oldSubscription.SubscriptionID).Error
		if errBulkInsertSubscriptionSharing != nil {
			tx.Rollback()
			return errBulkInsertSubscriptionSharing
		}
	}
	// 更新激活码的拥有者和生成的订阅信息
	errUpdateSubscriptionActivationCode := tx.Exec(`update subscription_activation_code set user_id = ? , subscription_id = ? where code = ?`, userID, subscription.SubscriptionID, activationCode.Code).Error
	if errUpdateSubscriptionActivationCode != nil {
		tx.Rollback()
		return errUpdateSubscriptionActivationCode
	}
	return tx.Commit().Error
}

func getLastAffectedRowCount(db *gorm.DB) (int, error) {
	var count int
	err := db.Raw(`SELECT ROW_COUNT() as count`).Row().Scan(&count)
	return count, err
}
