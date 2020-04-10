package handler

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/jiujiantang-services/device/mysqldb"
	"github.com/jinmukeji/go-pkg/mac"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/device/v1"
)

// UserGetUsedDevices 用户使用过的设备
func (j *DeviceManagerService) UserGetUsedDevices(ctx context.Context, req *proto.UserGetUsedDevicesRequest, resp *proto.UserGetUsedDevicesResponse) error {
	token, ok := TokenFromContext(ctx)
	if !ok {
		return NewError(ErrInvalidUser, errors.New("failed to get token from context"))
	}
	userID, err := j.database.FindUserIDByToken(ctx, token)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find userID from token: %s", err.Error()))
	}
	if userID != req.UserId {
		return NewError(ErrInvalidUser, fmt.Errorf("current user %d from token and user %d from request are inconsistent", userID, req.UserId))
	}
	usedDevices, errUserGetUsedDevices := j.database.UserGetUsedDevices(ctx, req.UserId)
	if errUserGetUsedDevices != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to get used devices of user %d: %s", req.UserId, errUserGetUsedDevices.Error()))
	}
	devices := make([]*proto.Device, len(usedDevices))

	for idx, device := range usedDevices {
		mac := strconv.FormatInt(device.MAC, 16)

		if len(mac)%2 == 1 {
			mac = fmt.Sprintf("%s%s", "0", mac)
		}
		devices[idx] = &proto.Device{
			Mac:            mac,
			Sn:             device.Sn,
			Pin:            device.Pin,
			Zone:           device.Zone,
			Model:          device.Model,
			CustomizedCode: device.CustomizedCode,
			Tags:           device.Tags,
			Remarks:        device.Remarks,
			ClientId:       device.ClientID,
		}
	}
	resp.Devices = devices
	return nil
}

// UserUseDevice 用户使用设备
func (j *DeviceManagerService) UserUseDevice(ctx context.Context, req *proto.UserUseDeviceRequest, resp *proto.UserUseDeviceResponse) error {
	now := time.Now()
	userDevice := &mysqldb.UserDevice{
		UserID:    req.UserId,
		DeviceID:  req.DeviceId,
		ClientID:  req.ClientId,
		CreatedAt: now.UTC(),
		UpdatedAt: now.UTC(),
	}
	exist, errExistUserDevice := j.database.ExistUserDevice(ctx, userDevice)
	if errExistUserDevice != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to check existence of user %d, device %d and client %s: %s", req.UserId, req.DeviceId, req.ClientId, errExistUserDevice.Error()))
	}
	if exist {
		return nil
	}
	errCreateUserDevice := j.database.CreateUserDevice(ctx, userDevice)
	if errCreateUserDevice != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to create record of user %d, device %d and client %s: %s", req.UserId, req.DeviceId, req.ClientId, errCreateUserDevice.Error()))
	}
	return nil
}

// GetDeviceByMac 通过MAC获取设备
func (j *DeviceManagerService) GetDeviceByMac(ctx context.Context, req *proto.GetDeviceByMacRequest, resp *proto.GetDeviceByMacResponse) error {
	hexMac, _ := strconv.ParseUint(mac.NormalizeMac(req.Mac), 16, 64)
	device, errGetDeviceByMac := j.database.GetDeviceByMac(ctx, hexMac)
	if errGetDeviceByMac != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to get device by mac %d: %s", hexMac, errGetDeviceByMac.Error()))
	}
	createdAt, _ := ptypes.TimestampProto(device.CreatedAt)
	resp.Device = &proto.Device{
		DeviceId:       device.DeviceID,
		Mac:            req.Mac,
		Sn:             device.Sn,
		Pin:            device.Pin,
		Zone:           device.Zone,
		Model:          device.Model,
		CustomizedCode: device.CustomizedCode,
		Tags:           device.Tags,
		Remarks:        device.Remarks,
		CreatedTime:    createdAt,
	}
	return nil
}
