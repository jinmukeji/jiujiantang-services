package handler

import (
	"context"

	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
)

var (
	// 语言 languages
	languages = []string{"zh-Hans", "zh-Hant", "en"}
	// 区域 regions
	regions = []string{"mainland_china", "taiwan", "abroad"}
)

// GerResourceList 获取资源列表
func (j *JinmuIDService) GerResourceList(ctx context.Context, req *proto.GerResourceListRequest, resp *proto.GerResourceListResponse) error {
	resp.Languages = languages
	resp.Regions = regions
	return nil
}
