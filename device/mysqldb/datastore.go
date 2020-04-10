package mysqldb

import "context"

// Datastore 定义数据访问接口
type Datastore interface {
	// FindUserIDByToken 根据 token 返回 userID，如果token失效返回 error
	// TODO 以后重写，删除这个方法
	FindUserIDByToken(ctx context.Context, token string) (int32, error)
	// GetDeviceByDeviceID 通过device_id 获取device
	GetDeviceByDeviceID(ctx context.Context, deviceID int32) (*Device, error)
	// ExistUserDevice UserDevice是否已经存在
	ExistUserDevice(ctx context.Context, userDevice *UserDevice) (bool, error)
	// CreateUserDevice 创建UserDevice
	CreateUserDevice(ctx context.Context, userDevice *UserDevice) error
	// UserGetUsedDevice 用户得到使用过的设备
	UserGetUsedDevices(ctx context.Context, userID int32) ([]*Device, error)
	// GetDeviceByMac 通关MAC获取device
	GetDeviceByMac(ctx context.Context, mac uint64) (*Device, error)
}
