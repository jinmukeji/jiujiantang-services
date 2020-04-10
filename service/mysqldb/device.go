package mysqldb

import (
	"context"
	"time"
)

// Device 用于表示设备信息
type Device struct {
	DeviceID       int        `gorm:"primary_key"`            // 表记录标识
	MAC            int64      `gorm:"column:mac"`             // 测量设备的MAC地址
	Sn             string     `gorm:"column:sn"`              // SN号
	Pin            string     `gorm:"column:pin"`             // 验证码
	Zone           string     `gorm:"column:zone"`            // 地区
	Model          string     `gorm:"column:model"`           // 设备型号
	CustomizedCode string     `gorm:"column:customized_code"` // 自定义代码
	Tags           string     `gorm:"column:tags"`            // 标签列表
	Remarks        string     `gorm:"column:remarks"`         // 备注
	CreatedAt      time.Time  // 创建时间
	UpdatedAt      time.Time  // 更新时间
	DeletedAt      *time.Time // 删除时间
}

// TableName 返回 Device 所在的表名
func (d Device) TableName() string {
	return "device"
}

// GetDeviceByDeviceID 通过DeviceID查询Device
func (db *DbClient) GetDeviceByDeviceID(ctx context.Context, deviceID int) (*Device, error) {
	var device Device
	err := db.Model(&Device{}).Where("device_id = ?", deviceID).Scan(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

// GetDeviceByMac 通过 mac 查询Device
func (db *DbClient) GetDeviceByMac(ctx context.Context, mac uint64) (*Device, error) {
	var device Device
	err := db.Model(&Device{}).Where("mac = ?", mac).Scan(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

// ExistDeviceByMac 查看 mac 能否存在
func (db *DbClient) ExistDeviceByMac(ctx context.Context, mac uint64) (bool, error) {
	var count int
	err := db.Model(&Device{}).Where("mac = ?", mac).Count(&count).Error
	return count == 1, err
}
