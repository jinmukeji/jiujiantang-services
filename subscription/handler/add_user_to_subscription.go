package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jinmukeji/jiujiantang-services/subscription/mysqldb"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/subscription/v1"
)

// AddUsersIntoSubscription 将用户添加到订阅中
func (j *SubscriptionService) AddUsersIntoSubscription(ctx context.Context, req *proto.AddUsersIntoSubscriptionRequest, resp *proto.AddUsersIntoSubscriptionResponse) error {
	token, ok := TokenFromContext(ctx)
	if !ok {
		return NewError(ErrInvalidUser, errors.New("failed to get token from context"))
	}
	userID, err := j.datastore.FindUserIDByToken(ctx, token)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find userID from token: %s", err.Error()))
	}
	if userID != req.OwnerId {
		return NewError(ErrInvalidUser, fmt.Errorf("current user %d from token is not the owner %d", userID, req.OwnerId))
	}
	users := make([]*mysqldb.UserSubscriptionSharing, len(req.UserIdList))
	now := time.Now()
	for idx, uid := range req.UserIdList {
		users[idx] = &mysqldb.UserSubscriptionSharing{
			SubscriptionID: req.SubscriptionId,
			UserID:         uid,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
	}

	errCreateUserSubscriptionSharing := j.datastore.CreateMultiUserSubscriptionSharing(ctx, users)
	if errCreateUserSubscriptionSharing != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to create multi user subscription sharing: %s", errCreateUserSubscriptionSharing.Error()))
	}

	return nil
}
