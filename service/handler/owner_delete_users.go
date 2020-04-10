package handler

import (
	"context"
	"errors"

	"fmt"

	"github.com/jinmukeji/jiujiantang-services/service/auth"
	corepb "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	subscriptionpb "github.com/jinmukeji/proto/gen/micro/idl/jm/subscription/v1"
)

// OwnerDeleteUsers Owner删除用户
func (j *JinmuHealth) OwnerDeleteUsers(ctx context.Context, req *corepb.OwnerDeleteUsersRequest, resp *corepb.OwnerDeleteUsersResponse) error {
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok || userID != req.OwnerId {
		return NewError(ErrInvalidUser, fmt.Errorf("ownerID %d from request and userID %d from context are inconsistent", req.OwnerId, userID))
	}
	o, err := j.datastore.FindOrganizationByUserID(ctx, int(userID))
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("ownerID %d from request and userID %d from context are inconsistent", req.OwnerId, userID))
	}
	organizationID := o.OrganizationID
	isOwner, err := j.datastore.CheckOrganizationOwner(ctx, int(userID), organizationID)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to check organization %d and user %d: %s", organizationID, userID, err.Error()))
	}
	if !isOwner {
		return NewError(ErrInvalidUser, fmt.Errorf("user %d is not owner of organization", userID))
	}
	// 获取当前拥有者的订阅
	reqGetUserSubscriptions := new(subscriptionpb.GetUserSubscriptionsRequest)
	reqGetUserSubscriptions.UserId = userID
	respGetUserSubscriptions, errGetUserSubscriptions := j.subscriptionSvc.GetUserSubscriptions(ctx, reqGetUserSubscriptions)
	if errGetUserSubscriptions != nil {
		return errGetUserSubscriptions
	}
	var subscriptionID int32
	// 获取当前正在使用的订阅
	for _, item := range respGetUserSubscriptions.Subscriptions {
		if item.IsSelected {
			subscriptionID = item.SubscriptionId
		}
	}

	for _, item := range req.UserIdList {
		if item == userID {
			return NewError(ErrForbidToRemoveSubscriptionOwner, errors.New("unable to delete owner"))
		}
		ok, err := j.datastore.CheckOrganizationUser(ctx, int(item), organizationID)
		if !ok || err != nil {
			return NewError(ErrUserNotInOrganization, fmt.Errorf("user %d is not user of organization %d", userID, organizationID))
		}

		// 删除与组织的绑定关系
		err = j.datastore.DeleteOrganizationUser(ctx, int(item), organizationID)
		if err != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to delete user %d from organization %d", item, organizationID))
		}
		// 删除订阅绑定关系
		reqDeleteUsersFromSubscription := new(subscriptionpb.DeleteUsersFromSubscriptionRequest)
		reqDeleteUsersFromSubscription.OwnerId = userID
		reqDeleteUsersFromSubscription.SubscriptionId = subscriptionID
		reqDeleteUsersFromSubscription.UserIdList = []int32{int32(item)}
		_, errDeleteUsersFromSubscription := j.subscriptionSvc.DeleteUsersFromSubscription(ctx, reqDeleteUsersFromSubscription)
		if errDeleteUsersFromSubscription != nil {
			return errDeleteUsersFromSubscription
		}
	}
	return nil
}
