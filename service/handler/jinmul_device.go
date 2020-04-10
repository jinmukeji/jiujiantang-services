package handler

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jinmukeji/gf-api2/service/mysqldb"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
)

// JinmuLBulkBindDevices 关联多个 mac
func (j *JinmuHealth) JinmuLBulkBindDevices(ctx context.Context, req *proto.JinmuLBulkBindDevicesRequest, repl *proto.JinmuLBulkBindDevicesResponse) error {
	for _, device := range req.Devices {
		// 16进制字符串转int
		mac, _ := strconv.ParseUint(device.Mac, 16, 64)
		isValid, err := j.datastore.ExistDeviceByMac(ctx, mac)
		if err != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to check if device eof mac %d exists: %s", mac, err.Error()))
		}
		if !isValid {
			return NewError(ErrInvalidDevice, fmt.Errorf("device [MAC: %s] is invalid", device.Mac))
		}
		d, _ := j.datastore.GetDeviceByMac(ctx, mac)
		// Device 与 组织是否已经关联
		existing, err := j.datastore.ExistOrganizationDeviceByID(ctx, req.OrganizationId, d.DeviceID)
		if err != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to check association existence of device %d and organization %d: %s", d.DeviceID, req.OrganizationId, err.Error()))
		}
		if existing {
			continue
		}
		// Device 是否已经被关联
		existingOrganizationDevice, err := j.datastore.ExistOrganizationDeviceByDeviceID(ctx, d.DeviceID)
		if err != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to  check organization existence by device %d: %s", d.DeviceID, err.Error()))
		}
		if existingOrganizationDevice {
			return NewError(ErrBindedDevice, fmt.Errorf("device [MAC: %s] is not bindable", device.Mac))
		}
		if errValidateDevice := j.validateDevice(ctx, d); errValidateDevice != nil {
			return errValidateDevice
		}
	}
	now := time.Now()
	successfulDevices := make([]*proto.Device, 0)
	failedDevices := make([]*proto.Device, 0)
	for _, device := range req.Devices {
		// 16进制字符串转int
		mac, _ := strconv.ParseUint(device.Mac, 16, 64)
		deviceTable, _ := j.datastore.GetDeviceByMac(ctx, mac)
		deviceOrganizationBinding := &mysqldb.DeviceOrganizationBinding{
			DeviceID:       deviceTable.DeviceID,
			OrganizationID: int(req.OrganizationId),
			UpdatedAt:      now,
			CreatedAt:      now,
		}

		// Device 与 组织是否已经关联
		existing, err := j.datastore.ExistOrganizationDeviceByID(ctx, req.OrganizationId, deviceTable.DeviceID)
		if err != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to find bind association between deviceID %d and organizationID %d: %s", deviceTable.DeviceID, req.OrganizationId, err.Error()))
		}
		if existing {
			successfulDevices = append(successfulDevices, &proto.Device{
				Mac: strings.ToUpper(device.Mac),
			})
			continue
		}
		if err := j.datastore.BindDeviceToOrganization(ctx, deviceOrganizationBinding); err != nil {
			failedDevices = append(failedDevices, &proto.Device{
				Mac: strings.ToUpper(device.Mac),
			})
		} else {
			successfulDevices = append(successfulDevices, &proto.Device{
				Mac: strings.ToUpper(device.Mac),
			})
		}
	}
	repl.FailedDevices = failedDevices
	repl.SuccessfulDevices = successfulDevices
	return nil
}

// JinmuLBulkUnbindDevices 解除关联多个Device
func (j *JinmuHealth) JinmuLBulkUnbindDevices(ctx context.Context, req *proto.JinmuLBulkUnbindDevicesRequest, repl *proto.JinmuLBulkUnbindDevicesResponse) error {
	for _, device := range req.Devices {
		// 16进制字符串转int
		mac, _ := strconv.ParseUint(device.Mac, 16, 64)
		isValid, err := j.datastore.ExistDeviceByMac(ctx, mac)
		if err != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to find device by mac %d: %s", mac, err.Error()))
		}
		if !isValid {
			return NewError(ErrInvalidDevice, fmt.Errorf("device [MAC: %s] is invalid", device.Mac))
		}
		d, _ := j.datastore.GetDeviceByMac(ctx, mac)
		existing, err := j.datastore.ExistOrganizationDeviceByID(ctx, req.OrganizationId, d.DeviceID)
		if err != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to find bind association between deviceID %d and organizationID %d: %s", d.DeviceID, req.OrganizationId, err.Error()))
		}
		if !existing {
			return NewError(ErrDeviceNotBelongToOrganization, fmt.Errorf("device [MAC: %s] is not belong to current organization", device.Mac))
		}
		if errValidateDevice := j.validateDevice(ctx, d); errValidateDevice != nil {
			return errValidateDevice
		}
	}
	successfulDevices := make([]*proto.Device, 0)
	failedDevices := make([]*proto.Device, 0)
	for _, device := range req.Devices {
		// 16进制字符串转int
		mac, _ := strconv.ParseUint(device.Mac, 16, 64)
		d, _ := j.datastore.GetDeviceByMac(ctx, mac)
		if err := j.datastore.UnbindOrganizationDevice(ctx, int32(req.OrganizationId), d.DeviceID); err != nil {
			failedDevices = append(failedDevices, &proto.Device{
				Mac: strings.ToUpper(device.Mac),
			})
		} else {
			successfulDevices = append(successfulDevices, &proto.Device{
				Mac: strings.ToUpper(device.Mac),
			})
		}
	}
	repl.FailedDevices = failedDevices
	repl.SuccessfulDevices = successfulDevices
	return nil
}

// JinmuLCheckDeviceBindable 确定能否关联
func (j *JinmuHealth) JinmuLCheckDeviceBindable(ctx context.Context, req *proto.JinmuLCheckDeviceBindableRequest, resp *proto.JinmuLCheckDeviceBindableResponse) error {
	// 16进制字符串转int
	mac, _ := strconv.ParseUint(req.Device.Mac, 16, 64)
	isValid, err := j.datastore.ExistDeviceByMac(ctx, mac)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find device by mac %d: %s", mac, err.Error()))
	}
	if !isValid {
		return NewError(ErrInvalidDevice, fmt.Errorf("device [MAC: %s] is invalid", req.Device.Mac))
	}
	d, _ := j.datastore.GetDeviceByMac(ctx, mac)
	if errValidateDevice := j.validateDevice(ctx, d); errValidateDevice != nil {
		return errValidateDevice
	}
	// Device 与 组织是否已经关联
	existing, err := j.datastore.ExistOrganizationDeviceByID(ctx, req.OrganizationId, d.DeviceID)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find bind association between deviceID %d and organizationID %d: %s", d.DeviceID, req.OrganizationId, err.Error()))
	}
	if existing {
		resp.Bindable = existing
		return nil
	}
	// Device 是否已经被关联
	existingOtherOrganizationDevice, err := j.datastore.ExistOrganizationDeviceByDeviceID(ctx, d.DeviceID)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to  check organization existence  by device %d: %s", d.DeviceID, err.Error()))
	}
	resp.Bindable = !existingOtherOrganizationDevice
	return nil
}

// JinmuLGetDeviceList 查找关联到的 mac
func (j *JinmuHealth) JinmuLGetDeviceList(ctx context.Context, req *proto.JinmuLGetDeviceListRequest, resp *proto.JinmuLGetDeviceListResponse) error {
	deviceOrganizationBindingList, err := j.datastore.GetOrganizationDeviceList(ctx, req.OrganizationId)
	if err != nil {
		return NewError(ErrDeviceNotFoundInOrganization, fmt.Errorf("cannot find deviceOrganizationBindingList of organization %d: %s", req.OrganizationId, err.Error()))
	}
	devices := make([]*proto.Device, len(deviceOrganizationBindingList))
	for idx, deviceOrganizationBinding := range deviceOrganizationBindingList {
		device, _ := j.datastore.GetDeviceByDeviceID(ctx, deviceOrganizationBinding.DeviceID)
		mac := strconv.FormatInt(device.MAC, 16)
		if len(mac)%2 == 1 {
			mac = fmt.Sprintf("%s%s", "0", mac)
		}
		devices[idx] = &proto.Device{
			Mac: strings.ToUpper(mac),
		}
	}
	resp.Devices = devices
	return nil
}
