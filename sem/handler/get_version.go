package handler

import (
	"context"

	"github.com/jinmukeji/gf-api2/sem/config"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/sem/v1"
)

// GetVersion 获取服务版本信息
func (j *SEMGateway) GetVersion(ctx context.Context, req *proto.GetVersionRequest, resp *proto.GetVersionResponse) error {
	resp.ServiceName = config.FullServiceName()
	resp.Version = config.ProductVersion
	return nil
}
