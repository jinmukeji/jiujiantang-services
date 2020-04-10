package rest

import (
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	"github.com/kataras/iris/v12"
)

// Device mac地址
type Device struct {
	Mac string `json:"mac"`
}

// DeviceList Mac的集合
type DeviceList []Device

// OrganizationID 组织的ID
type OrganizationID struct {
	OrganizationID int32 `json:"organization_id"`
}

// OwnerBulkBindDevices 关联多个Device
type OwnerBulkBindDevices struct {
	SuccessfulDevices []*proto.Device `json:"successful_devices"`
	FailedDevices     []*proto.Device `json:"failed_devices"`
}

// OwnerBulkUnbindDevices 解除关联多个Device
type OwnerBulkUnbindDevices struct {
	SuccessfulDevices []*proto.Device `json:"successful_devices"`
	FailedDevices     []*proto.Device `json:"failed_devices"`
}

// OwnerCheckDeviceBindable 确定Device是否关联
type OwnerCheckDeviceBindable struct {
	Bindable bool `json:"bindable"`
}

// JinmuLOwnerBulkBindDevices 关联多个Device
func (h *v2Handler) JinmuLOwnerBulkBindDevices(ctx iris.Context) {
	var deviceList DeviceList
	err := ctx.ReadJSON(&deviceList)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}
	organizationID, err := ctx.Params().GetInt("organization_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	devices := make([]*proto.Device, len(deviceList))
	for idx, device := range deviceList {
		devices[idx] = &proto.Device{
			Mac: device.Mac,
		}
	}
	req := new(proto.JinmuLBulkBindDevicesRequest)
	req.Devices = devices
	req.OrganizationId = int32(organizationID)
	resp, errBulkBindDevices := h.rpcSvc.JinmuLBulkBindDevices(
		newRPCContext(ctx), req,
	)
	if errBulkBindDevices != nil {
		writeRpcInternalError(ctx, errBulkBindDevices, false)
		return
	}
	rest.WriteOkJSON(ctx, OwnerBulkBindDevices{
		SuccessfulDevices: resp.SuccessfulDevices,
		FailedDevices:     resp.FailedDevices,
	})
}

// JinmuLOwnerBulkUnbindDevices 解除关联多个Device
func (h *v2Handler) JinmuLOwnerBulkUnbindDevices(ctx iris.Context) {
	var deviceList DeviceList
	err := ctx.ReadJSON(&deviceList)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}
	devices := make([]*proto.Device, len(deviceList))
	for idx, device := range deviceList {
		devices[idx] = &proto.Device{
			Mac: device.Mac,
		}
	}
	req := new(proto.JinmuLBulkUnbindDevicesRequest)
	req.Devices = devices
	organizationID, err := ctx.Params().GetInt("organization_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	req.OrganizationId = int32(organizationID)
	resp, errBulkUnbindDevices := h.rpcSvc.JinmuLBulkUnbindDevices(
		newRPCContext(ctx), req,
	)
	if errBulkUnbindDevices != nil {
		writeError(ctx, wrapError(ErrRPCInternal, "", errBulkUnbindDevices), false)
		return
	}
	rest.WriteOkJSON(ctx, OwnerBulkUnbindDevices{
		SuccessfulDevices: resp.SuccessfulDevices,
		FailedDevices:     resp.FailedDevices,
	})
}

// JinmuLOwnerCheckDeviceBindable 确定Device是否关联
func (h *v2Handler) JinmuLOwnerCheckDeviceBindable(ctx iris.Context) {
	organizationID, err := ctx.Params().GetInt("organization_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	var device Device
	err = ctx.ReadJSON(&device)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}
	req := new(proto.JinmuLCheckDeviceBindableRequest)
	req.Device = &proto.Device{
		Mac: device.Mac,
	}
	req.OrganizationId = int32(organizationID)
	resp, errCheckDeviceBindable := h.rpcSvc.JinmuLCheckDeviceBindable(
		newRPCContext(ctx), req,
	)
	if errCheckDeviceBindable != nil {
		writeRpcInternalError(ctx, errCheckDeviceBindable, false)
		return
	}
	rest.WriteOkJSON(ctx, OwnerCheckDeviceBindable{
		Bindable: resp.Bindable,
	})
}

// JinmuLOwnerGetOrganizationDeviceList 查找关联到的Device
func (h *v2Handler) JinmuLOwnerGetOrganizationDeviceList(ctx iris.Context) {
	organizationID, err := ctx.Params().GetInt("organization_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	req := new(proto.JinmuLGetDeviceListRequest)
	req.OrganizationId = int32(organizationID)
	resp, err := h.rpcSvc.JinmuLGetDeviceList(
		newRPCContext(ctx), req,
	)
	if err != nil {
		writeRpcInternalError(ctx, err, false)
		return
	}
	devices := make([]Device, len(resp.Devices))
	for idx, device := range resp.Devices {
		devices[idx] = Device{
			Mac: device.Mac,
		}
	}
	rest.WriteOkJSON(ctx, devices)
}
