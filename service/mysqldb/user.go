package mysqldb

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

// Gender 性别
type Gender string

const (
	// GenderMale 男性
	GenderMale Gender = "M"
	// GenderFemale 女性
	GenderFemale Gender = "F"
	// GenderInvalid 非法的性别
	GenderInvalid Gender = ""
)

// User 用户
type User struct {
	UserID                 int        `gorm:"primary_key"`             // 用户 id
	Username               string     `gorm:"column:username"`         // 用户名
	Nickname               string     `gorm:"column:nickname"`         // 用户昵称
	NicknameInitial        string     `gorm:"column:nickname_initial"` // 昵称首字母
	Password               string     // 用户登录密码
	Gender                 Gender     // 用户性别 M 男 F 女
	Birthday               time.Time  // 用户生日
	Height                 int        // 用户身高 单位厘米
	Weight                 int        // 用户体重 单位千克
	Phone                  string     // 用户电话号码
	Email                  string     // 用户联系邮箱
	State                  string     // 用户所在省份
	City                   string     // 用户所在城市
	Street                 string     // 用户所在街道
	Zone                   string     // 用户选择的区域
	Remark                 string     // 用户备注
	RegisterType           string     `gorm:"column:register_type"` // 用户注册方式
	CustomizedCode         string     // 用户自定义代码
	IsActivated            int        `gorm:"is_activated"`              // 是否激活 1激活 0没有激活
	ActivatedAt            *time.Time `gorm:"activated_at"`              // 激活时间
	DeactivatedAt          *time.Time `gorm:"deactivated_at"`            // 禁用时间
	RegisterSourceClientID string     `gorm:"register_source_client_id"` // 注册的Client_id
	UserDefinedCode        string     `gorm:"user_defined_code"`         // 用户自定义代码
	District               string     // 地区
	Country                string     // 国家
	RegisterTime           time.Time  // 注册时间
	IsProfileCompleted     bool       `gorm:"is_profile_completed"` // profile是否完整/完成初始化 1是true 0是false
	IsRemovable            bool       `gorm:"-"`                    // 是否能够删除
	CreatedAt              time.Time  // 创建时间
	UpdatedAt              time.Time  // 更新时间
	DeletedAt              *time.Time // 删除时间
}

// ProtoUserProfile 接口
type ProtoUserProfile interface {
	GetPhone() string
	GetHeight() int32
	GetWeight() int32
	GetEmail() string
	GetBirthday() *timestamp.Timestamp
	GetState() string
	GetCity() string
	GetStreet() string
	GetRemark() string
	GetNickname() string
	GetCountry() string
	GetDistrict() string
	GetUserDefinedCode() string
}

// TableName 返回 User 所在的表名
func (u User) TableName() string {
	return "user"
}

// FindUserByUsername 返回数据库中的用户
func (db *DbClient) FindUserByUsername(ctx context.Context, username string) (*User, error) {
	var u User
	if err := db.Raw(`SELECT 
		U.user_id, 
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
		U.street,
		U.zone,
		U.remark,
		U.register_type,
		U.customized_code,
		U.is_activated,
		U.deactivated_at,
		U.register_source_client_id,
		U.user_defined_code,
		U.district,
		U.register_time,
		(OO.owner_id IS NULL) AS is_removable,
		U.is_profile_completed,
		U.created_at,
		U.updated_at,
		U.deleted_at
		FROM user AS U
		LEFT JOIN organization_owner AS OO ON U.user_id = OO.owner_id
		WHERE U.username = ? AND OO.deleted_at IS NULL`, username).Scan(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// CreateUser 创建一个新用户
func (db *DbClient) CreateUser(ctx context.Context, u *User) (*User, error) {
	if err := db.Create(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

// FindUserByUserID 返回数据库中的用户
func (db *DbClient) FindUserByUserID(ctx context.Context, userID int) (*User, error) {
	var u User
	if err := db.Raw(`SELECT 
	U.user_id, 
	U.signin_username, 
	UP.nickname,
	UP.gender, 
	UP.birthday,
	UP.height, 
	UP.weight, 
	U.zone,
	U.remark,
	U.register_type,
	U.customized_code,
	U.user_defined_code,
	U.register_time,
	(OO.owner_id IS NULL) AS is_removable,
	U.is_profile_completed,
	U.created_at,
	U.updated_at,
	U.deleted_at
	FROM user AS U
	LEFT JOIN organization_owner AS OO ON U.user_id = OO.owner_id
	inner join user_profile as UP on UP.user_id = U.user_id
	WHERE U.user_id = ? AND OO.deleted_at IS NULL`, userID).Scan(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// UpdateUserProfile 更新用户个人信息
func (db *DbClient) UpdateUserProfile(ctx context.Context, profile ProtoUserProfile, userID int32) error {
	birthday, _ := ptypes.Timestamp(profile.GetBirthday())
	return db.Model(&User{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"nickname":             profile.GetNickname(),
		"birthday":             birthday,
		"height":               profile.GetHeight(),
		"weight":               profile.GetWeight(),
		"phone":                profile.GetPhone(),
		"email":                profile.GetEmail(),
		"state":                profile.GetState(),
		"city":                 profile.GetCity(),
		"street":               profile.GetStreet(),
		"remark":               profile.GetRemark(),
		"user_defined_code":    profile.GetUserDefinedCode(),
		"updated_at":           time.Now().UTC(),
		"district":             profile.GetDistrict(),
		"country":              profile.GetCountry(),
		"is_profile_completed": true,
	}).Error
}

// UpdateUserProfileContainGender 更新用户含有gender的个人信息
func (db *DbClient) UpdateUserProfileContainGender(ctx context.Context, profile ProtoUserProfile, userID int32, gender string) error {
	birthday, _ := ptypes.Timestamp(profile.GetBirthday())
	return db.Model(&User{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"nickname":             profile.GetNickname(),
		"birthday":             birthday,
		"height":               profile.GetHeight(),
		"weight":               profile.GetWeight(),
		"gender":               gender,
		"phone":                profile.GetPhone(),
		"email":                profile.GetEmail(),
		"state":                profile.GetState(),
		"city":                 profile.GetCity(),
		"street":               profile.GetStreet(),
		"remark":               profile.GetRemark(),
		"user_defined_code":    profile.GetUserDefinedCode(),
		"updated_at":           time.Now().UTC(),
		"district":             profile.GetDistrict(),
		"country":              profile.GetCountry(),
		"is_profile_completed": true,
	}).Error
}

// UpdateSimpleUserProfile 简易更新用户个人信息
func (db *DbClient) UpdateSimpleUserProfile(ctx context.Context, user *User) error {
	return db.Model(&User{}).Where("user_id = ?", user.UserID).Updates(map[string]interface{}{
		"height":   user.Height,
		"weight":   user.Weight,
		"gender":   user.Gender,
		"birthday": user.Birthday,
	}).Error
}

// GetUserByRecordID 通过RecordID 获取 User
func (db *DbClient) GetUserByRecordID(ctx context.Context, recordID int32) (*User, error) {
	var user User
	if err := db.Raw(`SELECT 
	user.user_id, 
	user.nickname,
	user.gender, 
	user.birthday,
	user.weight, 
	user.height, 
	user.gender
	FROM user INNER JOIN record 
	ON user.user_id = record.user_id 
	WHERE record.record_id = ? AND record.deleted_at IS NULL`, recordID).Scan(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CountUserByCustomizedCode 通过customizedCode计算User的数量
func (db *DbClient) CountUserByCustomizedCode(ctx context.Context, customizedCode string) (int, error) {
	var count int
	count, err := db.countByRaw(ctx, `SELECT count(U.user_id) FROM user AS U
	inner Join organization_user AS OU  ON U.user_id = OU.user_id 
	AND OU.deleted_at IS NULL
	where U.customized_code = ? AND U.deleted_at IS NULL`, customizedCode)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// countByRaw 通过raw拿到数量
func (db *DbClient) countByRaw(ctx context.Context, rawSQL string, params ...interface{}) (int, error) {
	var count int
	row := db.Raw(rawSQL, params...).Row()
    errScan := row.Scan(&count)
    if errScan != nil{
        return -1, errScan
    }
	return count, nil
}

// GetIsRemovableStatus 获取用户是否可移除的状态
func (db *DbClient) GetIsRemovableStatus(ctx context.Context, username string) (bool, error) {
	count, err := db.countByRaw(ctx, `SELECT 
    count(U.user_id)
    FROM user AS U
    INNER JOIN organization_owner AS OO ON U.user_id = OO.owner_id
    WHERE U.signin_username = ? AND U.has_set_username = ? AND OO.deleted_at IS NULL 
    AND U.deleted_at IS NULL`, username, true)
	if count == 0 || err != nil {
		return true, err
	}
	return false, nil
}
