package mysqldb

import (
	"context"
	"time"
)

// DeviceOrganizationBinding Device和组织关联
type DeviceOrganizationBinding struct {
	DeviceID       int       `gorm:"column:device_id"`
	OrganizationID int       `gorm:"column:organization_id"`
	CreatedAt      time.Time // 创建时间
	UpdatedAt      time.Time // 更新时间
}

// TableName 返回 Mac 所在的表名
func (m *DeviceOrganizationBinding) TableName() string {
	return "device_organization_binding"
}

// BindDeviceToOrganization 关联设备到组织
func (db *DbClient) BindDeviceToOrganization(ctx context.Context, d *DeviceOrganizationBinding) error {
	return db.GetDB(ctx).Create(d).Error
}

// ExistOrganizationDeviceByID 检查组织和deviceID关联是否存在
func (db *DbClient) ExistOrganizationDeviceByID(ctx context.Context, organizationID int32, deviceID int) (bool, error) {
	var count int
	err := db.GetDB(ctx).Model(&DeviceOrganizationBinding{}).Where("organization_id = ? AND device_id = ?", organizationID, deviceID).Count(&count).Error
	return count > 0, err
}

// UnbindOrganizationDevice 解除关联设备到组织
func (db *DbClient) UnbindOrganizationDevice(ctx context.Context, organizationID int32, deviceID int) error {
	return db.GetDB(ctx).Where("organization_id = ? AND device_id = ?", organizationID, deviceID).Delete(&DeviceOrganizationBinding{}).Error
}

// GetOrganizationDeviceList 通过organizationID查询与Device的关联关系
func (db *DbClient) GetOrganizationDeviceList(ctx context.Context, organizationID int32) ([]*DeviceOrganizationBinding, error) {
	var deviceOrganizationBindingList []*DeviceOrganizationBinding
	err := db.GetDB(ctx).Model(&DeviceOrganizationBinding{}).Where("organization_id = ?", organizationID).Scan(&deviceOrganizationBindingList).Error
	if err != nil {
		return nil, err
	}
	return deviceOrganizationBindingList, nil
}

// ExistOrganizationDeviceByDeviceID 查看 Device 是否已经被关联
func (db *DbClient) ExistOrganizationDeviceByDeviceID(ctx context.Context, deviceID int) (bool, error) {
	var count int
	err := db.GetDB(ctx).Model(&DeviceOrganizationBinding{}).Where("device_id = ?", deviceID).Count(&count).Error
	return count > 0, err
}
