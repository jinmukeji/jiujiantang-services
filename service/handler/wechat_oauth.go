package handler

import (
	"context"

	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
)

// WechatBuildOAuthURL 生成网页授权地址
func (j *JinmuHealth) WechatBuildOAuthURL(ctx context.Context, req *proto.WechatBuildOAuthURLRequest, resp *proto.WechatBuildOAuthURLResponse) error {
	w := j.wechat
	resp.AuthCodeUrl = w.BuildOAuthURL(req.AuthRedirectUrl, req.State)
	return nil
}
