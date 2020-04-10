package handler

import (
	"context"

	"github.com/jinmukeji/gf-api2/analysis/config"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/analysis/v1"
)

// GetVersion 获取服务版本信息
func (j *AnalysisManagerService) GetVersion(ctx context.Context, req *proto.GetVersionRequest, resp *proto.GetVersionResponse) error {
	resp.ServiceName = config.FullServiceName()
	resp.Version = config.ProductVersion
	return nil
}
