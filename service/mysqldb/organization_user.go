package mysqldb

import (
	"context"
	"time"
)

// OrganizationUser 组织下用户
type OrganizationUser struct {
	OrganizationID int        `gorm:"primary_key;column:organization_id"` // 组织ID
	UserID         int        `gorm:"primary_key;column:user_id"`         // 所有者ID
	CreatedAt      time.Time  // 创建时间
	UpdatedAt      time.Time  // 更新时间
	DeletedAt      *time.Time // 删除时间
}

// TableName 返回 OrganizationUser对应的数据库数据表名
func (o OrganizationUser) TableName() string {
	return "organization_user"
}

// CreateOrganizationUsers 在 organization_user 新增多条记录
func (db *DbClient) CreateOrganizationUsers(ctx context.Context, users []*OrganizationUser) error {
	tx := db.GetDB(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	for _, u := range users {
		if err := tx.Create(u).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

// DeleteOrganizationUser 在 organization_user 删除一条记录
func (db *DbClient) DeleteOrganizationUser(ctx context.Context, userID, organizationID int) error {
	return db.GetDB(ctx).Where("organization_id = ? AND user_id = ?", organizationID, userID).Delete(&OrganizationUser{}).Error
}

// DeleteOrganizationUsers 在 organization_user 删除多条记录
func (db *DbClient) DeleteOrganizationUsers(ctx context.Context, userIDList []int32, organizationID int32) error {
	tx := db.GetDB(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	for _, uid := range userIDList {
		if err := tx.Where("user_id = ?", uid).Delete(&OrganizationUser{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

// CheckOrganizationUser 检查 user 是否为组织下用户
func (db *DbClient) CheckOrganizationUser(ctx context.Context, userID, organizationID int) (bool, error) {
	var count int
	if err := db.GetDB(ctx).Raw("select count(*) from organization_user where user_id = ? AND organization_id = ?", userID, organizationID).Count(&count).Error; err != nil {
		return false, err
	}
	return count != 0, nil
}

// FindOrganizationUsers 查看组织下用户
func (db *DbClient) FindOrganizationUsers(ctx context.Context, organizationID int) ([]*User, error) {
	var users []*User
	db.GetDB(ctx).Raw(`SELECT
			DISTINCT OU.user_id,
			U.username,
			U.nickname,
			U.password,
			U.gender,
			U.birthday,
			U.height,
			U.weight,
			U.phone,
			U.email,
			U.state,
			U.city,
			U.zone,
			U.remark,
			U.register_type,
			U.customized_code,
			U.register_time,
			U.created_at,
			U.updated_at,
			U.country,
			U.district,
			U.user_defined_code,
			U.is_profile_completed,
			( OO.owner_id IS NULL ) AS is_removable
		FROM
			organization_user AS OU
		INNER JOIN user AS U ON OU.user_id = U.user_id
		LEFT JOIN organization_owner AS OO ON OO.owner_id = OU.user_id
		WHERE
			OU.organization_id = ? AND OU.deleted_at is NULL`, organizationID).Scan(&users)
	return users, nil
}

// FindOrganizationUsersByOffset 通过分页搜索用户
func (db *DbClient) FindOrganizationUsersByOffset(ctx context.Context, organizationID int32, size int32, offset int32) ([]*User, error) {
	var users []*User
	if size == -1 {
		size = maxQuerySize
	}
	db.GetDB(ctx).Limit(size).Offset(offset).Raw(`SELECT
	DISTINCT OU.user_id,
	U.signin_username AS username,
    UP.nickname,
    case when UP.nickname_initial = '~' THEN '#' ELSE UP.nickname_initial END as nickname_initial,
    (UP.nickname_initial = '~') AS at_last,
	UP.gender,
	UP.birthday,
	UP.height,
	UP.weight,
	U.signin_phone AS phone,
	U.secure_email AS email,
	U.zone,
	U.remark,
	U.register_type,
	U.customized_code,
	U.register_time,
	U.created_at,
	U.updated_at,
	U.user_defined_code,
	U.is_profile_completed,
	( OO.owner_id IS NULL ) AS is_removable
FROM
	organization_user AS OU
INNER JOIN user AS U ON OU.user_id = U.user_id AND U.deleted_at IS NULL
INNER JOIN user_profile AS UP ON UP.user_id = U.user_id
LEFT JOIN organization_owner AS OO ON OO.owner_id = OU.user_id
WHERE
	OU.organization_id = ? AND OU.deleted_at is NULL
ORDER BY is_removable, at_last, UP.nickname_initial ASC, convert(UP.nickname using gbk) ASC`, organizationID).Scan(&users)
	return users, nil
}

// FindOrganizationUsersByKeyword 通过keyword和分页搜索用户
func (db *DbClient) FindOrganizationUsersByKeyword(ctx context.Context, organizationID int32, keyword string, size int32, offset int32) ([]*User, error) {
	var users []*User
	db.GetDB(ctx).Raw(`SELECT
	DISTINCT OU.user_id,
	U.signin_username  AS username,
	UP.nickname,
	case when UP.nickname_initial = '~' THEN '#' ELSE UP.nickname_initial END as nickname_initial,
	UP.gender,
	UP.birthday,
	UP.height,
	UP.weight,
	U.signin_phone AS phone, 
	U.secure_email AS email,
	U.zone,
	U.remark,
	U.register_type,
	U.customized_code, 
	U.register_time,
	U.created_at,
	U.updated_at,
	U.user_defined_code,
	U.is_profile_completed,
	( OO.owner_id IS NULL ) AS is_removable
FROM
	organization_user AS OU
INNER JOIN user AS U ON OU.user_id = U.user_id AND U.deleted_at IS NULL
INNER JOIN user_profile AS UP ON UP.user_id = U.user_id
LEFT JOIN organization_owner AS OO ON OO.owner_id = OU.user_id
WHERE
        OU.organization_id = ? AND OU.deleted_at is NULL AND (UP.nickname LIKE ?)
ORDER BY UP.nickname_initial ASC, convert(UP.nickname using gbk) ASC 
limit ? offset ?`, organizationID, sqlEscape(keyword)+"%", size, offset).Scan(&users)
	return users, nil
}

// FindOrganizationUsersIDList 查看组织下用户ID的List
func (db *DbClient) FindOrganizationUsersIDList(ctx context.Context, organizationID int) ([]int, error) {
	var users []*User
	db.GetDB(ctx).Raw(`SELECT
		OU.user_id
    FROM 
        organization_user AS OU
	WHERE
        OU.organization_id = ? AND OU.deleted_at is NULL`, organizationID).Scan(&users)
	usersIDList := make([]int, len(users))
	for idx, item := range users {
		usersIDList[idx] = item.UserID
	}
	return usersIDList, nil
}
