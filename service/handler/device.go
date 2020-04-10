package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/jinmukeji/jiujiantang-services/service/mysqldb"
)

// macWithColon 返回由冒号分隔的 MAC 地址
func macWithColon(mac string) string {
	length := len(mac) / 2
	sa := make([]string, length)

	for i := 0; i < length; i++ {
		sa[i] = mac[2*i : 2*i+2]
	}

	return strings.Join(sa, ":")
}

// validateDevice 验证设备
func (j *JinmuHealth) validateDevice(ctx context.Context, device *mysqldb.Device) error {
	client, _ := clientFromContext(ctx)
	customizedCode := client.CustomizedCode
	zone := client.Zone
	if device.CustomizedCode != customizedCode || device.Zone != zone {
		return NewError(ErrDeviceNotBelongToClient, fmt.Errorf("device [MAC: %x] is not belong to current client", device.MAC))
	}
	return nil
}
