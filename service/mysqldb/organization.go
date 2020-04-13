package mysqldb

import (
	"context"
	"time"
)

// Organization 组织
type Organization struct {
	OrganizationID int        `gorm:"primary_key"` // 组织ID
	Name           string     // 组织名称
	Phone          string     // 固定电话
	Contact        string     // 联系人
	Type           string     // 组织机构类型
	State          string     // 组织所在省份
	City           string     // 组织所在城市
	Street         string     // 组织所在街道
	Address        string     // 组织地址
	IsValid        int        // 是否有效
	Email          string     // 邮箱
	District       string     // 区域
	Country        string     // 国家
	PostalCode     string     // 邮编
	CustomizedCode string     `gorm:"column:customized_code"` // 自定义代码
	CreatedAt      time.Time  // 创建时间
	UpdatedAt      time.Time  // 更新时间
	DeletedAt      *time.Time // 删除时间
}

// TableName 返回 Organization 对应的数据库数据表名
func (o Organization) TableName() string {
	return "organization"
}

// CreateOrganization 创建组织
func (db *DbClient) CreateOrganization(ctx context.Context, o *Organization) error {
	return db.GetDB(ctx).Create(o).Error
}

// FindFirstOrganizationByOwner 查找指定 owner 拥有的第一个组织
func (db *DbClient) FindFirstOrganizationByOwner(ctx context.Context, ownerID int) (*Organization, error) {
	var o Organization
	if err := db.GetDB(ctx).Raw(`SELECT 
		organization.organization_id, 
		organization.name, 
		organization.phone,
		organization.contact, 
		organization.type, 
		organization.state,
		organization.city, 
		organization.street, 
		organization.address,
		organization.email, 
		organization.district, 
		organization.country, 
		organization.postal_code,
		organization.customized_code,
		organization.created_at,
		organization.updated_at,
		organization.deleted_at
		FROM organization INNER JOIN organization_owner 
		ON organization_owner.organization_id = organization.organization_id 
		WHERE organization_owner.owner_id = ? AND organization.deleted_at IS NULL`, ownerID).Scan(&o).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

// FindOrganizationByID 从组织 ID 查找组织
func (db *DbClient) FindOrganizationByID(ctx context.Context, organizationID int) (*Organization, error) {
	var o Organization
	if err := db.GetDB(ctx).First(&o, "organization_id = ?", organizationID).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

// FindOrganizationsByOwner 查找指定 owner 拥有的组织
func (db *DbClient) FindOrganizationsByOwner(ctx context.Context, ownerID int) ([]*Organization, error) {
	var os []*Organization
	if err := db.GetDB(ctx).Raw(`SELECT 
		organization.organization_id, 
		organization.name, 
		organization.phone,
		organization.contact, 
		organization.type, 
		organization.state,
		organization.city, 
		organization.street, 
		organization.address, 
		organization.email, 
		organization.district, 
		organization.country, 
		organization.postal_code,
		organization.customized_code,
		organization.created_at,
		organization.updated_at,
		organization.deleted_at
		FROM organization INNER JOIN organization_owner
		ON organization_owner.organization_id = organization.organization_id 
		WHERE organization_owner.owner_id = ? AND organization_owner.deleted_at is NULL order by organization_owner.organization_id`, ownerID).Scan(&os).Error; err != nil {
		return nil, err
	}
	return os, nil
}

// UpdateOrganizationProfile 更新组织信息
func (db *DbClient) UpdateOrganizationProfile(ctx context.Context, o *Organization) error {
	return db.GetDB(ctx).Model(&Organization{}).Where("organization_id = ?", o.OrganizationID).Update(map[string]interface{}{
		"name":        o.Name,
		"phone":       o.Phone,
		"contact":     o.Contact,
		"type":        o.Type,
		"state":       o.State,
		"city":        o.City,
		"street":      o.Street,
		"email":       o.Email,
		"country":     o.Country,
		"district":    o.District,
		"postal_code": o.PostalCode,
	}).Error
}

// DeleteOrganizationByID 删除组织
func (db *DbClient) DeleteOrganizationByID(ctx context.Context, organizationID int) error {
	tx := db.GetDB(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	// 删除组织和owner的关联表
	if err := tx.Where("organization_id = ?", organizationID).Delete(&OrganizationOwner{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	now := time.Now()
	// 组织is_valid更改为0
	err := tx.Model(&Organization{}).Where("organization_id = ?", organizationID).Update(map[string]interface{}{
		"is_valid":   0,
		"updated_at": now.UTC(),
		"deleted_at": now.UTC(),
	}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// GetExistingUserCountByOrganizationID 查找指定 user的数量
func (db *DbClient) GetExistingUserCountByOrganizationID(ctx context.Context, organizationID int) (int, error) {
	var count int
	db.GetDB(ctx).Table("organization_user").Where("organization_id=? AND organization_user.deleted_at is NULL", organizationID).Count(&count)
	return count, nil
}

// GetOrganizationCountByOwnerID 查找组织的数量
func (db *DbClient) GetOrganizationCountByOwnerID(ctx context.Context, ownerID int) (int, error) {
	var count int
	db.GetDB(ctx).Table("organization_owner").Where("owner_id=? AND organization_owner.deleted_at is NULL", ownerID).Count(&count)
	return count, nil
}

// FindOrganizationByUserID 从UserID 查找组织
func (db *DbClient) FindOrganizationByUserID(ctx context.Context, userID int) (*Organization, error) {
	// TODO: 以后多组织需要重构这个方法
	var organizationUser OrganizationUser
	if err := db.GetDB(ctx).Table("organization_user").Where("user_id = ?", userID).Scan(&organizationUser).Error; err != nil {
		return nil, err
	}
	var o Organization
	if err := db.GetDB(ctx).First(&o, "organization_id = ?", organizationUser.OrganizationID).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

// CheckOrganizationIsValid 检查组织是否有效
func (db *DbClient) CheckOrganizationIsValid(ctx context.Context, organizationID int) bool {
	var organization Organization
	db.GetDB(ctx).First(&organization, "organization_id = ?", organizationID)
	return organization.IsValid == 1
}
