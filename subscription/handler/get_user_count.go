package handler

import (
	"context"

	"fmt"

	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/subscription/v1"
)

// GetUserCountOfSubscription  获取订阅下的用户数量
func (j *SubscriptionService) GetUserCountOfSubscription(ctx context.Context, req *proto.GetUserCountOfSubscriptionRequest, resp *proto.GetUserCountOfSubscriptionResponse) error {
	count, errGetUserCount := j.datastore.GetUserCountOfSubscription(ctx, req.SubscriptionId)
	if errGetUserCount != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to get user count of subscription %d: %s", req.SubscriptionId, errGetUserCount.Error()))
	}
	resp.UserCount = count
	return nil
}
