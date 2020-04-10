package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/jinmukeji/jiujiantang-services/service/auth"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
)

// JinmuLAccountSignOut 账户登出
func (j *JinmuHealth) JinmuLAccountSignOut(ctx context.Context, req *proto.JinmuLAccountSignOutRequest, resp *proto.JinmuLAccountSignOutResponse) error {
	token, ok := auth.TokenFromContext(ctx)
	if !ok {
		return NewError(ErrInvalidAccessToken, errors.New("failed to get access token from context"))
	}
	if err := j.datastore.DeleteJinmuLAccessToken(ctx, token); err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to delete token %s: %s", token, err.Error()))
	}
	return nil
}
