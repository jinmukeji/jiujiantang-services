package mysqldb

import (
	"context"
	"time"
)

// Frequency 性别
type Frequency string

const (
	// FrequencyDaily 推送基本单位为天
	FrequencyDaily Frequency = "FrequencyDaily"
	// FrequencyWeekly 推送基本单位为周
	FrequencyWeekly Frequency = "FrequencyWeekly"
	// FrequencyMonthly 推送基本单位为月
	FrequencyMonthly Frequency = "FrequencyMonthly"
)

// LocalNotification 本地消息推送
type LocalNotification struct {
	LnID    int32  `gorm:"primary_key;column:ln_id"`
	Title   string `gorm:"column:title"` // 推送标题
	Content string `gorm:"content"`      // 推送内容

	EventHappenAt string `gorm:"column:event_happen_at"` // 推送时间
	Timezone      string `gorm:"column:timezone"`        // 时区信息

	Frequency Frequency `gorm:"column:frequency"` // 推送基本单位
	Interval  int32     `gorm:"column:interval"`  // 推送时间间隔

	HasWeekdays          bool    `gorm:"column:has_weekdays"`           // 是否需要以周为基本单位
	Weekdays             string  `gorm:"column:weekdays"`               // 一周内有哪些天推送
	HasMonthDays         bool    `gorm:"column:has_month_days"`         // 是否需要以月为基本单位
	MonthDays            string  `gorm:"column:month_days"`             // 一个月内有哪些天推送
	MaxNotificationTimes *int32  `gorm:"column:max_notification_times"` // 最大推送次数
	EndAt                *string `gorm:"column:end_at"`                 // 推送结束时间

	CreatedAt time.Time  // 创建时间
	UpdatedAt time.Time  // 更新时间
	DeletedAt *time.Time // 删除时间
}

// TableName 返回 LocalNotification 所在的表名
func (d LocalNotification) TableName() string {
	return "local_notifications"
}

// CreateLocalNotification 在 local_notifications新增一条记录
func (db *DbClient) CreateLocalNotification(ctx context.Context, record *LocalNotification) error {
	return db.GetDB(ctx).Create(record).Error
}

// GetLocalNotifications 拿到本地通知
func (db *DbClient) GetLocalNotifications(ctx context.Context) ([]LocalNotification, error) {

	var pns []LocalNotification
	err := db.GetDB(ctx).Raw(`SELECT
	L.title,
	L.content,
	L.event_happen_at,
	L.timezone,
	L.frequency,
    L.interval,
    L.has_weekdays,
    L.weekdays,
    L.has_month_days,
    L.month_days,
    L.max_notification_times,
    L.end_at
    FROM local_notifications as L WHERE L.deleted_at is null`).Scan(&pns).Error
	if err != nil {
		return nil, err
	}
	return pns, nil
}

// DeleteLocalNotification 根据ID删除数据库指定的本地推送内容
func (db *DbClient) DeleteLocalNotification(ctx context.Context, lnID int) error {
	return db.GetDB(ctx).Delete(LocalNotification{}, "ln_id = ?", lnID).Error
}
