package handler

import (
	"context"
	"fmt"

	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
)

// GetOwnerIDByOrganizationID 根据organizationID获取ownerID
func (j *JinmuHealth) GetOwnerIDByOrganizationID(ctx context.Context, req *proto.GetOwnerIDByOrganizationIDRequest, resp *proto.GetOwnerIDByOrganizationIDResponse) error {
	ownersIDList, err := j.datastore.GetOwnerIDByOrganizationID(ctx, int(req.OrganizationId))
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find ownersIDList by organization %d, error: %s", req.OrganizationId, err.Error()))
	}
	if len(ownersIDList) > 1 {
		return NewError(ErrMultiOwnersOfOrganization, fmt.Errorf("more than one owner of organization %d", req.OrganizationId))
	}
	if len(ownersIDList) < 1 {
		return NewError(ErrNonexistentOwnerOfOrganization, fmt.Errorf("nonexistent organization %d or organization %d has no owner", req.OrganizationId, req.OrganizationId))
	}
	resp.OwnerId = int32(ownersIDList[0])
	return nil
}
