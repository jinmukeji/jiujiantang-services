package mysqldb

import (
	"context"
	"time"
)

// UserDevice 用户使用设备表
type UserDevice struct {
	UserID    int32      `gorm:"column:user_id"`   // 用户ID
	DeviceID  int32      `gorm:"column:device_id"` // 表记录标识
	ClientID  string     `gorm:"column:client_id"` // 客户端ID
	CreatedAt time.Time  // 创建时间
	UpdatedAt time.Time  // 更新时间
	DeletedAt *time.Time // 删除时间
}

// TableName 返回 UserDevice 所在的表名
func (d UserDevice) TableName() string {
	return "user_used_device"
}

// CreateUserDevice 创建UserDevice
func (db *DbClient) CreateUserDevice(ctx context.Context, userDevice *UserDevice) error {
	return db.Create(userDevice).Error
}

// ExistUserDevice UserDevice是否存在
func (db *DbClient) ExistUserDevice(ctx context.Context, userDevice *UserDevice) (bool, error) {
	var count int
	err := db.Model(&UserDevice{}).Where("user_id = ? and device_id = ? and client_id = ? ", userDevice.UserID, userDevice.DeviceID, userDevice.ClientID).Count(&count).Error
	return count >= 1, err
}
