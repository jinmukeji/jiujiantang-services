package mysqldb

import (
	"context"
	"time"
)

// Gender 性别
type Gender string

const (
	// GenderMale 性别为男
	GenderMale Gender = "M"
	// GenderFemale 性别为女
	GenderFemale Gender = "F"
)

// UserProfile 用户档案
type UserProfile struct {
	UserID          int32      `gorm:"primary_key"`             // 用户 id
	Nickname        string     `gorm:"column:nickname"`         // 用户名
	NicknameInitial string     `gorm:"column:nickname_initial"` // 昵称首字母
	Gender          Gender     `gorm:"column:gender"`           // 用户性别 M 男 F 女
	Birthday        time.Time  `gorm:"column:birthday"`         // 用户生日
	Height          int32      `gorm:"column:height"`           // 用户身高 单位厘米
	Weight          int32      `gorm:"column:weight"`           // 用户体重 单位千克
	CreatedAt       time.Time  // 创建时间
	UpdatedAt       time.Time  // 更新时间
	DeletedAt       *time.Time // 删除时间
}

// TableName 返回 UserProfile 所在的表名
func (u UserProfile) TableName() string {
	return "user_profile"
}

// CreateUserProfile 创建UserProfile
func (db *DbClient) CreateUserProfile(ctx context.Context, u *UserProfile) (*UserProfile, error) {
	if err := db.DB(ctx).Create(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

// ModifyUserProfile 修改用户档案
func (db *DbClient) ModifyUserProfile(ctx context.Context, profile *UserProfile) error {
	return db.DB(ctx).Model(&UserProfile{}).Where("user_id = ?", profile.UserID).Updates(map[string]interface{}{
		"nickname":         profile.Nickname,
		"nickname_initial": profile.NicknameInitial,
		"gender":           profile.Gender,
		"birthday":         profile.Birthday,
		"height":           profile.Height,
		"weight":           profile.Weight,
	}).Error
}

// FindUserProfile 找到用户档案
func (db *DbClient) FindUserProfile(ctx context.Context, userID int32) (*UserProfile, error) {
	var profile UserProfile
	if err := db.DB(ctx).First(&profile, "( user_id = ? ) ", userID).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

// FindUserProfileByRecordID 通过记录ID获取用户档案
func (db *DbClient) FindUserProfileByRecordID(ctx context.Context, recordID int32) (*UserProfile, error) {
	var userProfile UserProfile
	if err := db.DB(ctx).Raw(`SELECT 
	UP.user_id, 
    UP.nickname,
    case when UP.nickname_initial = '~' THEN '#' ELSE UP.nickname_initial END as nickname_initial,
	UP.gender, 
	UP.birthday,
	UP.weight, 
	UP.height, 
	UP.gender
	FROM user_profile as UP INNER JOIN record as R
	ON UP.user_id = R.user_id 
	WHERE R.record_id = ? AND R.deleted_at IS NULL`, recordID).Scan(&userProfile).Error; err != nil {
		return nil, err
	}
	return &userProfile, nil
}

// CreateUserAndUserProfile 创建用户和用户资料
func (db *DbClient) CreateUserAndUserProfile(ctx context.Context, user *User, userProfile *UserProfile) error {
	tx := db.DB(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return err
	}
	userProfile.UserID = user.UserID
	if err := tx.Create(userProfile).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
