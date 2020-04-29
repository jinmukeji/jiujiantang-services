package rest

import (
	"errors"
	"fmt"

	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	"github.com/kataras/iris/v12"
)

// RegionBody 区域的body
type RegionBody struct {
	UserID int32  `json:"user_id"`
	Region *int32 `json:"region"`
}

const (
	RegionMainlandChina = 0
	RegionTaiwan        = 1
	RegionAboard        = 2
)

// SelectRegion 选择区域
func (h *webHandler) SelectRegion(ctx iris.Context) {
	var body RegionBody
	errReadJSON := ctx.ReadJSON(&body)
	if errReadJSON != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", errReadJSON), false)
		return
	}
	if body.Region == nil {
		writeError(ctx, wrapError(ErrEmptyRegion, "", errors.New("region is empty")), false)
		return
	}
	req := new(proto.UserSelectRegionRequest)
	req.UserId = body.UserID
	protoRegion, errMapRestRegionToProto := mapRestRegionToProto(*body.Region)
	if errMapRestRegionToProto != nil {
        writeError(ctx, wrapError(ErrInvalidValue, "", errMapRestRegionToProto), false)
        return
	}
	req.Region = protoRegion
	_, errUserSelectRegion := h.rpcSvc.UserSelectRegion(newRPCContext(ctx), req)
	if errUserSelectRegion != nil {
		writeRpcInternalError(ctx, errUserSelectRegion, false)
		return
	}
	rest.WriteOkJSON(ctx, nil)
}

// mapRestRegionToProto 将 Rest 的 Region 转化成proto中格式
func mapRestRegionToProto(region int32) (proto.Region, error) {
	switch region {
	case RegionMainlandChina:
		return proto.Region_REGION_MAINLAND_CHINA, nil
	case RegionTaiwan:
		return proto.Region_REGION_TAIWAN, nil
	case RegionAboard:
		return proto.Region_REGION_ABROAD, nil
	}
	return proto.Region_REGION_MAINLAND_CHINA, fmt.Errorf("invalid int32 region %d", region)
}
