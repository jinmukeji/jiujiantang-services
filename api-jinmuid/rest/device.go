package rest

import (
	"github.com/jinmukeji/gf-api2/pkg/rest"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/device/v1"
	"github.com/kataras/iris/v12"
)

// Device 设备
type Device struct {
	ClientID string `json:"client_id"`
	SN       string `json:"sn"`
	Model    string `json:"model"`
	Mac      string `json:"mac"`
}

// UserUsedDevice 用户使用的设备
func (h *webHandler) UserUsedDevice(ctx iris.Context) {
	userID, err := ctx.Params().GetInt("user_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	req := new(proto.UserGetUsedDevicesRequest)
	req.UserId = int32(userID)
	resp, errUserGetUsedDevices := h.rpcDeviceSvc.UserGetUsedDevices(newRPCContext(ctx), req)
	if errUserGetUsedDevices != nil {
		writeRpcInternalError(ctx, errUserGetUsedDevices, false)
		return
	}
	devices := make([]Device, len(resp.Devices))
	for idx, device := range resp.Devices {
		devices[idx] = Device{
			ClientID: device.ClientId,
			SN:       device.Sn,
			Model:    device.Model,
			Mac:      device.Mac,
		}
	}
	rest.WriteOkJSON(ctx, devices)
}
