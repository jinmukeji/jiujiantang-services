package handler

import (
	"context"
	"errors"
	"fmt"

	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
)

// GetNotificationPreferences 获取通知配置首选项
func (j *JinmuIDService) GetNotificationPreferences(ctx context.Context, req *proto.GetNotificationPreferencesRequest, resp *proto.GetNotificationPreferencesResponse) error {
	token, ok := TokenFromContext(ctx)
	if !ok {
		return NewError(ErrInvalidUser, errors.New("failed to get token from context"))
	}
	userID, err := j.datastore.FindUserIDByToken(ctx, token)
	if err != nil {
		return NewError(ErrUserUnauthorized, fmt.Errorf("failed to get user_id by token: %s", err.Error()))
	}
	if userID != req.UserId {
		return NewError(ErrInvalidUser, fmt.Errorf("user %d from request and user %d from token are inconsistent", req.UserId, userID))
	}
	_, errFindUserByUserID := j.datastore.FindUserByUserID(ctx, req.UserId)
	if errFindUserByUserID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find user by userID %d: %s", req.UserId, errFindUserByUserID.Error()))
	}

	notificationPreferences, errGetNotificationPreferences := j.datastore.GetNotificationPreferences(ctx, req.UserId)
	if errGetNotificationPreferences != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to get notifications preferences by userID %d: %s", req.UserId, errGetNotificationPreferences.Error()))
	}
	resp.PhoneEnabled = notificationPreferences.PhoneEnabled
	resp.WechatEnabled = notificationPreferences.WechatEnabled
	resp.WeiboEnabled = notificationPreferences.WeiboEnabled
	return nil
}
