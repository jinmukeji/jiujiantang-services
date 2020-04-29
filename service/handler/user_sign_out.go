package handler

import (
	"context"

	corepb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	jinmuidpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
)

// UserSignOut 注销用户
func (j *JinmuHealth) UserSignOut(ctx context.Context, req *corepb.UserSignOutRequest, resp *corepb.UserSignOutResponse) error {
	reqUserSignOut := new(jinmuidpb.UserSignOutRequest)
	reqUserSignOut.Ip = req.Ip
	_, errUserSignOut := j.jinmuidSvc.UserSignOut(ctx, reqUserSignOut)
	if errUserSignOut != nil {
		return errUserSignOut
	}
	resp.Tip = "Successfully signed-out"
	return nil
}
