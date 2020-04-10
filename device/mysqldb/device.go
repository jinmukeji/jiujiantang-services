package mysqldb

import (
	"context"
	"time"
)

// Device 用于表示设备信息
type Device struct {
	DeviceID       int32      `gorm:"primary_key;column:device_id"` // 表记录标识
	MAC            int64      `gorm:"column:mac"`                   // 测量设备的MAC地址
	Sn             string     `gorm:"column:sn"`                    // SN号
	Pin            string     `gorm:"column:pin"`                   // 验证码
	Zone           string     `gorm:"column:zone"`                  // 地区
	Model          string     `gorm:"column:model"`                 // 设备型号
	CustomizedCode string     `gorm:"column:customized_code"`       // 自定义代码
	Tags           string     `gorm:"column:tags"`                  // 标签列表
	Remarks        string     `gorm:"column:remarks"`               // 备注
	ClientID       string     `gorm:"-"`                            // 客户端ID
	CreatedAt      time.Time  // 创建时间
	UpdatedAt      time.Time  // 更新时间
	DeletedAt      *time.Time // 删除时间
}

// TableName 返回 Device 所在的表名
func (d Device) TableName() string {
	return "device"
}

// GetDeviceByDeviceID 通过DeviceID查询Device
func (db *DbClient) GetDeviceByDeviceID(ctx context.Context, deviceID int32) (*Device, error) {
	var device Device
	err := db.Model(&Device{}).Where("device_id = ?", deviceID).Scan(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

// UserGetUsedDevices 用户得到使用过的设备
func (db *DbClient) UserGetUsedDevices(ctx context.Context, userID int32) ([]*Device, error) {
	var devices []*Device
	err := db.Raw(`SELECT 
	D.device_id,
	D.mac,
	D.sn,
	D.pin,
	D.model,
	D.customized_code,
	D.tags,
	D.remarks,
	D.zone,
	UD.client_id as client_id
	FROM user_used_device as UD
	inner join device as D on D.device_id = UD.device_id  and UD.user_id = ? and D.deleted_at IS NULL
	where  UD.deleted_at IS NULL`, userID).Scan(&devices).Error
	if err != nil {
		return nil, err
	}
	return devices, nil
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
