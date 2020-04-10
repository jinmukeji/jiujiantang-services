package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jinmukeji/jiujiantang-services/jinmuid/mysqldb"

	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
)

// ModifyUserInformation 修改用户信息
func (j *JinmuIDService) ModifyUserInformation(ctx context.Context, req *proto.ModifyUserInformationRequest, resp *proto.ModifyUserInformationResponse) error {
	// 如果需要token验证
	if !req.IsSkipVerifyToken {
		token, ok := TokenFromContext(ctx)
		if !ok {
			return NewError(ErrInvalidUser, errors.New("failed to get token from context"))
		}
		userID, err := j.datastore.FindUserIDByToken(ctx, token)
		if err != nil {
			return NewError(ErrUserUnauthorized, fmt.Errorf("failed to get userID by token: %s", err.Error()))
		}
		if userID != req.UserId {
			return NewError(ErrInvalidUser, fmt.Errorf("user %d from request and user %d from token are inconsistent", req.UserId, userID))
		}
	}
	now := time.Now()
	user := &mysqldb.User{
		UserID:             req.UserId,
		Remark:             req.Remark,
		CustomizedCode:     req.CustomizedCode,
		UserDefinedCode:    req.UserDefinedCode,
		HasSetUserProfile:  true,
		IsProfileCompleted: true,
		UpdatedAt:          now,
	}
	if req.SigninUsername != "" {
		user.SigninUsername = req.SigninUsername
		user.HasSetUsername = true
		user.LatestUpdatedUsernameAt = &now
	}
	errModifyUser := j.datastore.ModifyUser(ctx, user)
	if errModifyUser != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to update user %d: %s", req.UserId, errModifyUser.Error()))
	}
	return nil
}
