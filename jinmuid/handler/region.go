package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/jinmukeji/jiujiantang-services/jinmuid/mysqldb"

	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
)

// UserSelectRegion 用户选择区域
func (j *JinmuIDService) UserSelectRegion(ctx context.Context, req *proto.UserSelectRegionRequest, resp *proto.UserSelectRegionResponse) error {
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
	user, err := j.datastore.FindUserByUserID(ctx, req.UserId)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find user by userID %d: %s", req.UserId, err.Error()))
	}
	// 区域已经被设置就报这个错误
	if user.HasSetRegion {
		return NewError(ErrExsitRegion, fmt.Errorf("region has been set for userId %d", req.UserId))
	}
	region, errmapProtoRegionToDB := mapProtoRegionToDB(req.Region)
	if errmapProtoRegionToDB != nil {
		return NewError(ErrInvalidRegion, errmapProtoRegionToDB)
	}
	errSetUserRegion := j.datastore.SetUserRegion(ctx, req.UserId, region)
	if errSetUserRegion != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to set region for userId %d: %s", req.UserId, errSetUserRegion.Error()))
	}
	return nil
}

func mapProtoRegionToDB(region proto.Region) (mysqldb.Region, error) {
	switch region {
	case proto.Region_REGION_INVALID:
		return mysqldb.MainlandChina, fmt.Errorf("invalid proto region %d", region)
	case proto.Region_REGION_UNSET:
		return mysqldb.MainlandChina, fmt.Errorf("invalid proto region %d", region)
	case proto.Region_REGION_MAINLAND_CHINA:
		return mysqldb.MainlandChina, nil
	case proto.Region_REGION_TAIWAN:
		return mysqldb.Taiwan, nil
	case proto.Region_REGION_ABROAD:
		return mysqldb.Abroad, nil
	}
	return mysqldb.MainlandChina, fmt.Errorf("invalid proto region %d", region)
}
