package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jinmukeji/jiujiantang-services/jinmuid/mysqldb"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
)

// ModifyNotificationPreferences 修改通知配置首选项
func (j *JinmuIDService) ModifyNotificationPreferences(ctx context.Context, req *proto.ModifyNotificationPreferencesRequest, resp *proto.ModifyNotificationPreferencesResponse) error {
	token, ok := TokenFromContext(ctx)
	if !ok {
		return NewError(ErrInvalidUser, errors.New("failed to get token from context"))
	}
	userID, err := j.datastore.FindUserIDByToken(ctx, token)
	if err != nil {
		return NewError(ErrUserUnauthorized, fmt.Errorf("failed to get userID by token: %s", err.Error()))
	}
	if userID != req.UserId {
		return NewError(ErrInvalidUser, fmt.Errorf("user %d from request and user %d from token are inconsistent", req.UserId, userID))
	}
	// 判断UserID是否存在
	exist, errExistUserByUserID := j.datastore.ExistUserByUserID(ctx, req.UserId)
	if !exist || errExistUserByUserID != nil {
		return NewError(ErrInvalidUser, fmt.Errorf("failed to check existence of user %d", req.UserId))
	}
	now := time.Now()
	notificationPreferences := &mysqldb.NotificationPreferences{
		UserID:                 userID,
		PhoneEnabled:           req.PhoneEnabled,
		PhoneEnabledUpdatedAt:  now,
		WechatEnabled:          req.WechatEnabled,
		WechatEnabledUpdatedAt: now,
		WeiboEnabled:           req.WeiboEnabled,
		WeiboEnabledUpdatedAt:  now,
		UpdatedAt:              now,
	}
	errUpdateNotificationPreferences := j.datastore.UpdateNotificationPreferences(ctx, notificationPreferences)
	if errUpdateNotificationPreferences != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to update notification preferences: %s", errUpdateNotificationPreferences.Error()))
	}

	return nil
}
