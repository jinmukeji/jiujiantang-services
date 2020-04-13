package mysqldb

import (
	"context"
	"time"
)

// PushNotification 通知
type PushNotification struct {
	PnID          int32      `gorm:"pn_id"`
	PnDisplayTime string     `gorm:"pn_display_time"`
	PnTitle       string     `gorm:"pn_title"`
	PnImageURL    string     `gorm:"pn_image_url"`
	PnContentURL  string     `gorm:"pn_content_url"`
	PnType        int32      `gorm:"pn_type"`
	CreatedAt     time.Time  // 创建时间
	UpdatedAt     time.Time  // 更新时间
	DeletedAt     *time.Time // 删除时间
}

// TableName 返回 QRCode 所在的表名
func (pn PushNotification) TableName() string {
	return "push_notification"
}

// GetPnsByUserID 通过userID拿到未读通知记录，按时间倒序
func (db *DbClient) GetPnsByUserID(ctx context.Context, UserID int32, size int32) ([]PushNotification, error) {
	var pns []PushNotification
	err := db.GetDB(ctx).Raw(`Select 
		PN.pn_id,
		PN.pn_title,
		PN.pn_display_time,
		PN.pn_image_url,
		PN.pn_content_url
	    from push_notification as PN where PN.pn_id not in 
		(select PR.pn_id from pn_record AS PR where PR.user_id = ? ) order by PN.created_at desc`, UserID).Scan(&pns).Error
	if err != nil {
		return nil, err
	}
	return pns, nil
}
