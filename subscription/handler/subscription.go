package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/jiujiantang-services/subscription/mysqldb"

	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/subscription/v1"
)

// GetUserSubscriptions 得到用户的订阅
func (j *SubscriptionService) GetUserSubscriptions(ctx context.Context, req *proto.GetUserSubscriptionsRequest, resp *proto.GetUserSubscriptionsResponse) error {
	token, ok := TokenFromContext(ctx)
	if !ok {
		return NewError(ErrInvalidUser, errors.New("failed get token from context"))
	}
	userID, err := j.datastore.FindUserIDByToken(ctx, token)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find userID from token: %s", err.Error()))
	}
	if userID != req.UserId {
		return NewError(ErrInvalidUser, fmt.Errorf("userID %d from token and userID %d from request are inconsistent", userID, req.UserId))
	}
	subscriptions, errGetSubscriptionByUserID := j.datastore.GetSubscriptionsByUserID(ctx, req.UserId)
	if errGetSubscriptionByUserID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to get subscription of userID %d: %s", req.UserId, errGetSubscriptionByUserID.Error()))
	}
	if len(subscriptions) == 0 {
		return NewError(ErrNoneExistSubscription, fmt.Errorf("user %d has no subscription", req.UserId))
	}
	resp.Subscriptions = make([]*proto.Subscription, len(subscriptions))
	for idx, subscription := range subscriptions {
		activatedAt, _ := ptypes.TimestampProto(subscription.ActivatedAt)
		expiredAt, _ := ptypes.TimestampProto(subscription.ExpiredAt)
		resp.Subscriptions[idx] = &proto.Subscription{
			SubscriptionId:      subscription.SubscriptionID,
			SubscriptionType:    mapDBSubscriptionTypeToProto(subscription.SubscriptionType),
			Activated:           subscription.Activated,
			ActivatedTime:       activatedAt,
			ExpiredTime:         expiredAt,
			CustomizedCode:      subscription.CustomizedCode,
			MaxUserLimits:       subscription.MaxUserLimits,
			ContractYear:        subscription.ContractYear,
			IsSelected:          subscription.IsSelected,
			IsMigratedActivated: subscription.IsMigratedActivated,
		}
	}

	return nil
}

// ActivateSubscription 激活订阅
func (j *SubscriptionService) ActivateSubscription(ctx context.Context, req *proto.ActivateSubscriptionRequest, resp *proto.ActivateSubscriptionResponse) error {
	token, ok := TokenFromContext(ctx)
	if !ok {
		return NewError(ErrInvalidUser, errors.New("failed get token from context"))
	}
	userID, err := j.datastore.FindUserIDByToken(ctx, token)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find userID from token: %s", err.Error()))
	}
	subscription, errFindSelectedSubscriptionByUserID := j.datastore.FindSelectedSubscriptionByUserID(ctx, userID)
	if errFindSelectedSubscriptionByUserID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find subscription by useID %d: %s", userID, errFindSelectedSubscriptionByUserID.Error()))
	}
	if subscription.SubscriptionID != req.SubscriptionId {
		return NewError(ErrSubscriptionNotBelongToUser, fmt.Errorf("subscription %d not belong to user %d", subscription.SubscriptionID, userID))
	}
	now := time.Now()
	s := &mysqldb.Subscription{
		SubscriptionID: subscription.SubscriptionID,
		Activated:      true,
		ActivatedAt:    now,
		ExpiredAt:      now.AddDate(int(subscription.ContractYear), 0, 0),
		UpdatedAt:      now,
	}
	errActivateSubscription := j.datastore.ActivateSubscription(ctx, s)
	if errActivateSubscription != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to activate subscription %d: %s", req.SubscriptionId, errActivateSubscription.Error()))
	}
	resp.ExpiredTime, err = ptypes.TimestampProto(now.AddDate(int(subscription.ContractYear), 0, 0))
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to parse time at %s: %s", now.AddDate(int(subscription.ContractYear), 0, 0), err.Error()))
	}
	return nil
}

// GetSelectedUserSubscription 得到选中用户的订阅
func (j *SubscriptionService) GetSelectedUserSubscription(ctx context.Context, req *proto.GetSelectedUserSubscriptionRequest, resp *proto.GetSelectedUserSubscriptionResponse) error {
	if !req.IsSkipVerifyToken {
		token, ok := TokenFromContext(ctx)
		if !ok {
			return NewError(ErrInvalidUser, errors.New("failed get token from context"))
		}
		userID, err := j.datastore.FindUserIDByToken(ctx, token)
		if err != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to find userID from token: %s", err.Error()))
		}
		if userID != req.OwnerId {
			return NewError(ErrInvalidUser, fmt.Errorf("userID %d from token and userID %d from request are inconsistent", userID, req.OwnerId))
		}
	}
	subscription, errGetSubscriptionByUserID := j.datastore.GetSelectedSubscriptionByUserID(ctx, req.UserId)
	if errGetSubscriptionByUserID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to get selected subscription by ownerID %d: %s", req.OwnerId, errGetSubscriptionByUserID.Error()))
	}
	if subscription == nil {
		return NewError(ErrNoneExistSubscription, fmt.Errorf("user %d has no subscription", req.UserId))
	}
	activatedAt, _ := ptypes.TimestampProto(subscription.ActivatedAt)
	expiredAt, _ := ptypes.TimestampProto(subscription.ExpiredAt)
	resp.Subscription = &proto.Subscription{
		SubscriptionId:      subscription.SubscriptionID,
		SubscriptionType:    mapDBSubscriptionTypeToProto(subscription.SubscriptionType),
		Activated:           subscription.Activated,
		ActivatedTime:       activatedAt,
		ExpiredTime:         expiredAt,
		CustomizedCode:      subscription.CustomizedCode,
		MaxUserLimits:       subscription.MaxUserLimits,
		ContractYear:        subscription.ContractYear,
		IsSelected:          subscription.IsSelected,
		IsMigratedActivated: subscription.IsMigratedActivated,
	}
	return nil
}

// CheckUserHaveSameSubscription 检查用户有相同的订阅
func (j *SubscriptionService) CheckUserHaveSameSubscription(ctx context.Context, req *proto.CheckUserHaveSameSubscriptionRequest, resp *proto.CheckUserHaveSameSubscriptionResponse) error {
	subscription, err := j.datastore.GetSelectedSubscriptionByUserID(ctx, req.OwnerId)
	if err != nil {
		resp.IsSameSubscription = false
		return NewError(ErrDatabase, fmt.Errorf("failed to get selected subscription by ownerID %d: %s", req.OwnerId, err.Error()))
	}
	exist, errCheckExistUserSubscriptionSharing := j.datastore.CheckExistUserSubscriptionSharing(ctx, req.UserId, subscription.SubscriptionID)
	if errCheckExistUserSubscriptionSharing != nil {
		resp.IsSameSubscription = false
		return NewError(ErrDatabase, fmt.Errorf("failed to get selected subscription by ownerID %d: %s", req.OwnerId, errCheckExistUserSubscriptionSharing.Error()))
	}
	if !exist {
		resp.IsSameSubscription = false
		return NewError(ErrUserNotShareSubscription, fmt.Errorf("user %d does not have subscription %d", req.UserId, subscription.SubscriptionID))
	}
	resp.IsSameSubscription = true
	return nil
}

func mapDBSubscriptionTypeToProto(dbType mysqldb.SubscriptionType) proto.SubscriptionType {
	switch dbType {
	case mysqldb.SubscriptionTypeCustomizedVersion:
		return proto.SubscriptionType_SUBSCRIPTION_TYPE_CUSTOMIZED_VERSION
	case mysqldb.SubscriptionTypeTrialVersion:
		return proto.SubscriptionType_SUBSCRIPTION_TYPE_TRIAL_VERSION
	case mysqldb.SubscriptionTypeGoldenVersion:
		return proto.SubscriptionType_SUBSCRIPTION_TYPE_GOLDEN_VERSION
	case mysqldb.SubscriptionTypePlatinumVersion:
		return proto.SubscriptionType_SUBSCRIPTION_TYPE_PLATINUM_VERSION
	case mysqldb.SubscriptionTypeDiamondVersion:
		return proto.SubscriptionType_SUBSCRIPTION_TYPE_DIAMOND_VERSION
	case mysqldb.SubscriptionTypeGiftVersion:
		return proto.SubscriptionType_SUBSCRIPTION_TYPE_GIFT_VERSION
	}
	return proto.SubscriptionType_SUBSCRIPTION_TYPE_INVALID
}
