package handler

import (
	"context"

	"github.com/jinmukeji/gf-api2/sms/config"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/sms/v1"
)

// GetVersion 获取服务版本信息
func (j *SMSGateway) GetVersion(ctx context.Context, req *proto.GetVersionRequest, resp *proto.GetVersionResponse) error {
	resp.ServiceName = config.FullServiceName()
	resp.Version = config.ProductVersion
	return nil
}
