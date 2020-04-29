package handler

import (
	"context"

	"fmt"

	"github.com/jinmukeji/jiujiantang-services/service/auth"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
)

// SubmitRemark 用户修改备注
func (j *JinmuHealth) SubmitRemark(ctx context.Context, req *proto.SubmitRemarkRequest, resp *proto.SubmitRemarkResponse) error {
	accessTokenType, _ := auth.AccessTokenTypeFromContext(ctx)
	if accessTokenType != AccessTokenTypeWeChatValue {
		var organizationID int
		// 比较传入的userID与根据recordID得到的UserID是否是同一个组织
		userID, _ := j.datastore.GetUserIDByRecordID(ctx, req.RecordId)
		userId := int(req.UserId)
		if req.UserId == -1 {
			u, _ := j.datastore.FindUserByUsername(ctx, req.Username)
			userId = u.UserID
		}
		o, _ := j.datastore.FindOrganizationByUserID(ctx, userId)
		organizationID = o.OrganizationID
		organization, _ := j.datastore.FindOrganizationByUserID(ctx, int(userID))
		if organization.OrganizationID != organizationID {
			return NewError(ErrNoPermissionSubmitRemark, fmt.Errorf("user %d has no permission to submit remark to record %d", req.UserId, req.RecordId))
		}
	}
	if err := j.datastore.UpdateRemarkByRecordID(ctx, int(req.RecordId), req.Remark); err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to update remark by record %d: %s", req.RecordId, err.Error()))
	}
	return nil
}
