package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jinmukeji/gf-api2/service/auth"
	"github.com/jinmukeji/gf-api2/service/mysqldb"

	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
)

// DefaultMaxSize 默认最大的通知数量100
const DefaultMaxSize = 100

// ReadPushNotification 阅读通知
func (j *JinmuHealth) ReadPushNotification(ctx context.Context, req *proto.ReadPushNotificationRequest, resp *proto.ReadPushNotificationResponse) error {
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return NewError(ErrInvalidUser, errors.New("failed to get userID from context"))
	}
	exist, err := j.datastore.ExistPnRecord(ctx, req.PnId, userID)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("nonexistent pn record %d of user %d: %s", req.PnId, userID, err.Error()))
	}
	if !exist {
		now := time.Now().UTC()
		pnRecord := &mysqldb.PnRecord{
			PnID:      req.PnId,
			UserID:    userID,
			CreatedAt: now,
			UpdatedAt: now,
		}
		err := j.datastore.CreatePnRecord(ctx, pnRecord)
		if err != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to create pn record %d: %s", req.PnId, err.Error()))
		}
	}
	return nil
}

// GetPushNotifications 得到通知
func (j *JinmuHealth) GetPushNotifications(ctx context.Context, req *proto.GetPushNotificationsRequest, resp *proto.GetPushNotificationsResponse) error {
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return NewError(ErrInvalidUser, errors.New("failed to get userID from context"))
	}
	pns, err := j.datastore.GetPnsByUserID(ctx, userID, DefaultMaxSize)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to get pns by user %d: %s", userID, err.Error()))
	}
	notifications := make([]*proto.Notification, 0)
	for _, pn := range pns {
		notifications = append(notifications, &proto.Notification{
			PnId:          pn.PnID,
			PnTitle:       pn.PnTitle,
			PnContentUrl:  pn.PnContentURL,
			PnDisplayTime: pn.PnDisplayTime,
			PnImageUrl:    pn.PnImageURL,
		})
	}
	resp.Notifications = notifications
	return nil
}
