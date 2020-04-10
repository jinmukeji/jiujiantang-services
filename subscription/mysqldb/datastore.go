package mysqldb

import "context"

// Datastore 定义数据访问接口
type Datastore interface {
	// FindUserIDByToken 根据 token 返回 userID，如果token失效返回 error
	FindUserIDByToken(ctx context.Context, token string) (int32, error)
	// CreateSubscription 创建订阅
	CreateSubscription(ctx context.Context, subscription *Subscription) (*Subscription, error)
	// GetSubscriptionByUserID 通过userID获取订阅
	GetSubscriptionsByUserID(ctx context.Context, userID int32) ([]*Subscription, error)
	// GetSubscriptionActivationCodeInfo 得到订阅信息
	GetSubscriptionActivationCodeInfo(ctx context.Context, code string) (*SubscriptionActivationCode, error)
	// ActivateSubscriptionActivationCode 设置激活码为已经激活的
	ActivateSubscriptionActivationCode(ctx context.Context, code string, subscriptionID, userID int32) error
	// CreateUserSubscriptionSharing 创建单个订阅拥有或者被分享记录
	CreateUserSubscriptionSharing(ctx context.Context, UserSubscriptionSharing *UserSubscriptionSharing) error
	// ActivateSubscription 激活订阅
	ActivateSubscription(ctx context.Context, subscription *Subscription) error
	// FindSelectedSubscriptionByUserID 通过用户ID找订阅
	FindSelectedSubscriptionByUserID(ctx context.Context, userID int32) (*Subscription, error)
	// 获取订阅下的用户数量
	GetUserCountOfSubscription(ctx context.Context, subscriptionID int32) (int32, error)
	// CreateMultiUserSubscriptionSharing 创建订阅多个拥有或者被分享记录
	CreateMultiUserSubscriptionSharing(ctx context.Context, userSubscriptionSharing []*UserSubscriptionSharing) error
	// CheckSubscriptionOwner 检查用户是否是订阅的拥有者
	CheckSubscriptionOwner(ctx context.Context, ownerID, subscriptionID int) (bool, error)
	// DeleteSubscriptionUsers 删除订阅下的用户
	DeleteSubscriptionUsers(ctx context.Context, userIDList []int32, subscriptionID int32) error
	// UpdateSubscriptionIsSelectedStatus 更新订阅的默认状态
	UpdateSubscriptionIsSelectedStatus(ctx context.Context, ownerID int32) error
	// UseSubscriptionActivationCode 使用激活码
	UseSubscriptionActivationCode(ctx context.Context, userID int32, code *SubscriptionActivationCode, uuid string) error
	// GetSelectedSubscriptionByUserID 获取选中的订阅
	GetSelectedSubscriptionByUserID(ctx context.Context, userID int32) (*Subscription, error)
	// CheckExistUserSubscriptionSharing 检查user是否拥有或被分享Subscription
	CheckExistUserSubscriptionSharing(ctx context.Context, userID int32, subscriptionID int32) (bool, error)
}
