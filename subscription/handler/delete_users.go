package handler

import (
	"context"
	"errors"
	"fmt"

	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/subscription/v1"
)

// DeleteUsersFromSubscription 将用户从订阅中删除
func (j *SubscriptionService) DeleteUsersFromSubscription(ctx context.Context, req *proto.DeleteUsersFromSubscriptionRequest, resp *proto.DeleteUsersFromSubscriptionResponse) error {
	token, ok := TokenFromContext(ctx)
	if !ok {
		return NewError(ErrInvalidUser, errors.New("failed to get token from context"))
	}
	ownerID, err := j.datastore.FindUserIDByToken(ctx, token)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find userID from token %s: %s", token, err.Error()))
	}
	if ownerID != req.OwnerId {
		return NewError(ErrInvalidUser, fmt.Errorf("current user %d from token is not the owner %d", ownerID, req.OwnerId))
	}

	ok, errCheckSubscriptionOwner := j.datastore.CheckSubscriptionOwner(ctx, int(ownerID), int(req.SubscriptionId))
	if errCheckSubscriptionOwner != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to check if %d is owner of subscription %d: %s", ownerID, req.SubscriptionId, errCheckSubscriptionOwner.Error()))
	}
	if !ok {
		return NewError(ErrDatabase, errors.New("failed to check owner of subscription"))
	}

	for _, userID := range req.UserIdList {
		if userID == ownerID {
			return NewError(ErrForbidToRemoveSubscriptionOwner, fmt.Errorf("cannot delete subscription owner %d", userID))
		}
	}
	if err := j.datastore.DeleteSubscriptionUsers(ctx, req.UserIdList, req.SubscriptionId); err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to delete users %d in subscription %d: %s", req.UserIdList, req.SubscriptionId, err.Error()))
	}
	return nil
}
