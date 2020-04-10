package mysqldb

import (
	"context"
	"time"
)

// WXUser 微信用户
type WXUser struct {
	UnionID        string     `gorm:"primary_key"`      // union_id 微信UnionID
	OpenID         string     `gorm:"open_id"`          // open_id 微信OpenID
	UserID         int32      `gorm:"user_id"`          // 用户ID
	OriginID       string     `gorm:"origin_id"`        // 公众号原始ID
	Nickname       string     `gorm:"nickname"`         // 用户昵称
	AvatarImageURL string     `gorm:"avatar_image_url"` // 头像
	CreatedAt      time.Time  // 创建时间
	UpdatedAt      time.Time  // 更新时间
	DeletedAt      *time.Time // 删除时间
}

// TableName 返回 WXUser 所在的表名
func (w WXUser) TableName() string {
	return "wechat_user"
}

// CreateWXUser 创建微信用户
func (db *DbClient) CreateWXUser(ctx context.Context, wXUser *WXUser) error {
	if err := db.Create(wXUser).Error; err != nil {
		return err
	}
	return nil
}

// ExistWXUser 存在微信用户
func (db *DbClient) ExistWXUser(ctx context.Context, unionID string) (bool, error) {
	var count int
	db.Table("wechat_user").Where("union_id=?", unionID).Count(&count)
	return count > 0, nil
}

// FindWXUserByUnionID 通过UnionId找WXUser
func (db *DbClient) FindWXUserByUnionID(ctx context.Context, UnionID string) (*WXUser, error) {
	var wxUser WXUser
	err := db.First(&wxUser, "union_id = ? ", UnionID).Error
	if err != nil {
		return nil, err
	}
	return &wxUser, nil
}

// FindWXUserByUserID 通过userID找WXUser
func (db *DbClient) FindWXUserByUserID(ctx context.Context, userID int32) (*WXUser, error) {
	var wxUser WXUser
	err := db.First(&wxUser, "user_id = ? ", userID).Error
	if err != nil {
		return nil, err
	}
	return &wxUser, nil
}

// FindWXUserByOpenID 通过openID找WXUser
func (db *DbClient) FindWXUserByOpenID(ctx context.Context, openID string) (*WXUser, error) {
	var wxUser WXUser
	err := db.First(&wxUser, "open_id = ? ", openID).Error
	if err != nil {
		return nil, err
	}
	return &wxUser, nil
}
