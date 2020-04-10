package mysqldb

import (
	"context"
	"time"
)

// OrganizationOwner 组织所有者
type OrganizationOwner struct {
	OrganizationID int        `gorm:"primary_key;column:organization_id"` // 组织ID
	OwnerID        int        `gorm:"primary_key;column:owner_id"`        // 所有者ID
	CreatedAt      time.Time  // 创建时间
	UpdatedAt      time.Time  // 更新时间
	DeletedAt      *time.Time // 删除时间
}

// TableName 返回 Organization 对应的数据库数据表名
func (o OrganizationOwner) TableName() string {
	return "organization_owner"
}

// CreateOrganizationOwner 在 organization_owner 新增一条记录
func (db *DbClient) CreateOrganizationOwner(ctx context.Context, o *OrganizationOwner) error {
	return db.Create(o).Error
}

// CheckOrganizationOwner 检查用户是否为组织的拥有者
func (db *DbClient) CheckOrganizationOwner(ctx context.Context, userID int, organizationID int) (bool, error) {
	var count int
	if err := db.Raw("select count(*) from organization_owner where owner_id = ? AND organization_id = ? AND deleted_at IS NULL", userID, organizationID).Count(&count).Error; err != nil {
		return false, err
	}
	return count != 0, nil
}

// CheckUserOwnerBelongToSameOrganization 检查用户和拥有这是否处于相同的组织下
func (db *DbClient) CheckUserOwnerBelongToSameOrganization(ctx context.Context, userID int32, ownerID int32) (bool, error) {
	var count int
	if err := db.Raw(`SELECT COUNT(OU1.organization_id) FROM organization_user AS OU1 
	INNER JOIN organization_user AS OU2 ON OU1.organization_id = OU2.organization_id 
    AND OU2.user_id = ? AND OU2.deleted_at IS NULL WHERE OU1.user_id = ? AND OU1.deleted_at IS NULL`, userID, ownerID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetOwnerIDByOrganizationID 根据组织的ID获取组织拥有者的ID
func (db *DbClient) GetOwnerIDByOrganizationID(ctx context.Context, organizationID int) ([]int, error) {
	var organizationOwner []OrganizationOwner
	db.Raw(`SELECT
        owner_id
    FROM 
        organization_owner
    WHERE
        organization_id = ? AND deleted_at is NULL`, organizationID).Scan(&organizationOwner)

	ownersIDList := make([]int, len(organizationOwner))
	for idx, item := range organizationOwner {
		ownersIDList[idx] = item.OwnerID
	}
	return ownersIDList, nil
}
