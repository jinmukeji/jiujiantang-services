package handler

import (
	"context"

	jinmuidpb "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	corepb "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
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
