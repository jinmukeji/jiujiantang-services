package mysqldb

import (
	"context"
	"errors"
	"time"
)

const (
	// QuestionCount 密保问题的数量
	QuestionCount = 3
)

// Language 语言
type Language string

const (
	// SimpleChinese 简体中文
	LanguageSimpleChinese Language = "zh-Hans"
	// TraditionalChinese 繁体中文
	LanguageTraditionalChinese Language = "zh-Hant"
	// English 英文
	LanguageEnglish Language = "en"
	// LanguageInvalid 非法的语言
	LanguageInvalid Language = ""
)

// ValidationType 安全问题验证类型
type ValidationType string

const (
	// PhoneValidationType 安全问题验证类型为电话
	PhoneValidationType ValidationType = "phone"
	// UsernameValidationType 安全问题验证类型为用户名
	UsernameValidationType ValidationType = "username"
)

// Region 区域
type Region string

const (
	// MainlandChina 中国大陆
	MainlandChina Region = "mainland_china"
	// Taiwan 台湾
	Taiwan Region = "taiwan_and_abroad"
	// Abroad 国外
	Abroad Region = "abroad"
)

// User 用户
type User struct {
	UserID                         int32      `gorm:"primary_key"`                               // 用户 id
	RegisterType                   string     `gorm:"column:register_type"`                      // 用户注册方式
	RegisterTime                   time.Time  `gorm:"column:register_time"`                      // 注册时间
	Zone                           string     `gorm:"column:zone"`                               // 用户选择的区域
	CustomizedCode                 string     `gorm:"column:customized_code"`                    // 定制化代码
	UserDefinedCode                string     `gorm:"column:user_defined_code"`                  // 用户自定义代码
	Remark                         string     `gorm:"column:remark"`                             // 备注
	EncryptedPassword              string     `gorm:"column:encrypted_password"`                 // 加密后的密码
	Seed                           string     `gorm:"column:seed"`                               // 种子
	SecureEmail                    string     `gorm:"column:secure_email"`                       // 安全邮箱
	SigninPhone                    string     `gorm:"column:signin_phone"`                       // 登录电话
	SigninUsername                 string     `gorm:"column:signin_username"`                    // 登录用户名
	NationCode                     string     `gorm:"column:nation_code"`                        // 国家代码
	HasSetEmail                    bool       `gorm:"column:has_set_email"`                      // 是否设置邮箱
	HasSetPhone                    bool       `gorm:"column:has_set_phone"`                      // 是否设置电话
	HasSetUsername                 bool       `gorm:"column:has_set_username"`                   // 是否设置用户名
	HasSetPassword                 bool       `gorm:"column:has_set_password"`                   // 是否设置密码
	HasSetSecureQuestions          bool       `gorm:"column:has_set_secure_questions"`           // 是否设置密保问题
	HasSetUserProfile              bool       `gorm:"column:has_set_user_profile"`               // 是否设置用户详情
	HasSetLanguage                 bool       `gorm:"column:has_set_language"`                   // 是否设置语言
	Language                       Language   `gorm:"column:language"`                           // 语言
	RegisterSource                 string     `gorm:"column:register_source"`                    // 注册来源
	LatestLoginTime                *time.Time `gorm:"column:latest_login_time"`                  // 最新登录时间
	SecureQuestion1                string     `gorm:"column:secure_question_1"`                  // 保密问题1
	SecureQuestion2                string     `gorm:"column:secure_question_2"`                  // 保密问题2
	SecureQuestion3                string     `gorm:"column:secure_question_3"`                  // 保密问题3
	SecureAnswer1                  string     `gorm:"column:secure_answer_1"`                    // 保密答案1
	SecureAnswer2                  string     `gorm:"column:secure_answer_2"`                    // 保密答案2
	SecureAnswer3                  string     `gorm:"column:secure_answer_3"`                    // 保密答案3
	LatestUpdatedEmailAt           *time.Time `gorm:"column:latest_updated_email_at"`            // 最新修改邮箱时间
	LatestUpdatedPhoneAt           *time.Time `gorm:"column:latest_updated_phone_at"`            // 最新修改手机时间
	LatestUpdatedPasswordAt        *time.Time `gorm:"column:latest_updated_password_at"`         // 最新修改密码时间
	LatestUpdatedUsernameAt        *time.Time `gorm:"column:latest_updated_username_at"`         // 最新修改用户名时间
	LatestUpdatedSecureQuestionsAt *time.Time `gorm:"column:latest_updated_secure_questions_at"` // 最新修改密保问题时间
	IsProfileCompleted             bool       `gorm:"column:is_profile_completed"`               // profile是否完整/完成初始化 1是true 0是false
	IsRemovable                    bool       `gorm:"-"`                                         // 是否能够删除
	HasSetRegion                   bool       `gorm:"column:has_set_region"`                     // 是否设置区域
	Region                         Region     `gorm:"column:region"`                             // 区域
	Nickname                       string     `gorm:"-"`
	NicknameInitial                string     `gorm:"-"`              // 昵称首字母
	Gender                         Gender     `gorm:"-"`              // 用户性别 M 男 F 女
	Birthday                       time.Time  `gorm:"-"`              // 用户生日
	Height                         int32      `gorm:"-"`              // 用户身高 单位厘米
	Weight                         int32      `gorm:"-"`              // 用户体重 单位千克
	IsActivated                    bool       `gorm:"is_activated"`   // 是否激活 1激活 0没有激活
	ActivatedAt                    *time.Time `gorm:"activated_at"`   // 激活时间
	DeactivatedAt                  *time.Time `gorm:"deactivated_at"` // 禁用时间
	CreatedAt                      time.Time  // 创建时间
	UpdatedAt                      time.Time  // 更新时间
	DeletedAt                      *time.Time // 删除时间
}

// SecureQuestion 密保问题
type SecureQuestion struct {
	SecureQuestionKey string // 密保问题
	SecureAnswer      string // 密保答案
}

// TableName 返回 User对应的数据库数据表名
func (user User) TableName() string {
	return "user"
}

// FindUserByPhone 通过电话找到用户
func (db *DbClient) FindUserByPhone(ctx context.Context, phone string, nationCode string) (*User, error) {
	var user User
	if err := db.GetDB(ctx).First(&user, "( signin_phone = ? and nation_code = ?) ", phone, nationCode).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindUserByUsername 通过用户名找到base64密码
func (db *DbClient) FindUserByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	if err := db.GetDB(ctx).First(&user, "( signin_username = ? ) ", username).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// SetLanguageByUserID 通过userID设置Language
func (db *DbClient) SetLanguageByUserID(ctx context.Context, userID int32, language string) error {
	return db.GetDB(ctx).Model(&User{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"language":         language,
		"has_set_language": true,
		"updated":          time.Now().UTC(),
	}).Error
}

// FindLanguageByUserID 通过userID找到Language
func (db *DbClient) FindLanguageByUserID(ctx context.Context, userID int32) (string, error) {
	var user User
	if err := db.GetDB(ctx).First(&user, "( user_id = ? ) ", userID).Error; err != nil {
		return "", err
	}
	return string(user.Language), nil
}

// ExistUserByUserID 查看 user 能否存在
func (db *DbClient) ExistUserByUserID(ctx context.Context, userID int32) (bool, error) {
	var count int
	err := db.GetDB(ctx).Model(&User{}).Where("user_id = ?", userID).Count(&count).Error
	return count == 1, err
}

// ExistPasswordByUserID 查看 password 能否存在
func (db *DbClient) ExistPasswordByUserID(ctx context.Context, userID int32) (bool, error) {
	var user User
	if err := db.GetDB(ctx).First(&user, "( user_id = ? ) ", userID).Error; err != nil {
		return false, err
	}
	return user.HasSetPassword, nil
}

// SetPasswordByUserID 通过userID设置密码
func (db *DbClient) SetPasswordByUserID(ctx context.Context, userID int32, encryptedPassword string, seed string) error {
	return db.GetDB(ctx).Model(&User{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"encrypted_password":         encryptedPassword,
		"seed":                       seed,
		"latest_updated_password_at": time.Now().UTC(),
		"has_set_password":           true,
	}).Error
}

// FindSecureQuestionByUserID 通过userID找到密保问题和答案
func (db *DbClient) FindSecureQuestionByUserID(ctx context.Context, userID int32) ([]SecureQuestion, error) {
	var user User
	if err := db.GetDB(ctx).First(&user, "( user_id = ? AND has_set_secure_questions = ?) ", userID, true).Error; err != nil {
		return nil, err
	}
	secureQuestions := []SecureQuestion{
		{user.SecureQuestion1, user.SecureAnswer1},
		{user.SecureQuestion2, user.SecureAnswer2},
		{user.SecureQuestion3, user.SecureAnswer3},
	}
	return secureQuestions, nil
}

// FindSecureQuestionByPhone 通过电话号码找到找到密保问题和答案
func (db *DbClient) FindSecureQuestionByPhone(ctx context.Context, phone string, nationCode string) ([]SecureQuestion, error) {
	var user User

	if err := db.GetDB(ctx).First(&user, "( signin_phone = ?  AND nation_code = ?  AND has_set_phone = ? AND has_set_secure_questions = ?) ", phone, nationCode, true, true).Error; err != nil {
		return nil, err
	}
	secureQuestions := []SecureQuestion{
		{user.SecureQuestion1, user.SecureAnswer1},
		{user.SecureQuestion2, user.SecureAnswer2},
		{user.SecureQuestion3, user.SecureAnswer3},
	}

	return secureQuestions, nil
}

// FindSecureQuestionByUsername 通过用户名找到密保问题和答案
func (db *DbClient) FindSecureQuestionByUsername(ctx context.Context, username string) ([]SecureQuestion, error) {
	var user User

	if err := db.GetDB(ctx).First(&user, "( signin_username = ? AND has_set_username = ? AND has_set_secure_questions = ?) ", username, true, true).Error; err != nil {
		return nil, err
	}
	secureQuestions := []SecureQuestion{
		{user.SecureQuestion1, user.SecureAnswer1},
		{user.SecureQuestion2, user.SecureAnswer2},
		{user.SecureQuestion3, user.SecureAnswer3},
	}
	return secureQuestions, nil
}

// SetPasswordByPhone 根据手机号重置密码
func (db *DbClient) SetPasswordByPhone(ctx context.Context, phone string, nationCode string, encryptedPassword string, seed string) error {
	return db.GetDB(ctx).Model(&User{}).Where("signin_phone = ? AND nation_code = ?  AND has_set_phone = ?", phone, nationCode, true).Updates(map[string]interface{}{
		"encrypted_password":         encryptedPassword,
		"seed":                       seed,
		"latest_updated_password_at": time.Now().UTC(),
		"has_set_password":           true,
	}).Error
}

// SetPasswordByUsername 根据用户名重置密码
func (db *DbClient) SetPasswordByUsername(ctx context.Context, username string, encryptedPassword string, seed string) error {
	return db.GetDB(ctx).Model(&User{}).Where("signin_username = ? AND has_set_username = ?", username, true).Updates(map[string]interface{}{
		"encrypted_password":         encryptedPassword,
		"seed":                       seed,
		"latest_updated_password_at": time.Now().UTC(),
		"has_set_password":           true,
	}).Error
}

// IsPasswordSameByPhone 根据手机号判断密码是否与之前密码相同
func (db *DbClient) IsPasswordSameByPhone(ctx context.Context, phone string, nationCode string, encryptedPassword string) (bool, error) {
	var count int
	err := db.GetDB(ctx).Model(&User{}).Where("signin_phone = ? AND nation_code = ? AND encrypted_password = ? AND has_set_phone = ? ", phone, nationCode, encryptedPassword, true).Count(&count).Error
	return count != 0, err
}

// IsPasswordSameByUsername 根据用户名判断密码是否与之前密码相同
func (db *DbClient) IsPasswordSameByUsername(ctx context.Context, username string, encryptedPassword string) (bool, error) {
	var count int
	err := db.GetDB(ctx).Model(&User{}).Where("signin_username = ? AND encrypted_password = ? AND has_set_username = ?", username, encryptedPassword, true).Count(&count).Error
	return count != 0, err
}

// FindUserIDByPhone 通过电话号码找到userID
func (db *DbClient) FindUserIDByPhone(ctx context.Context, phoneNumber string, nationCode string) (int32, error) {
	var user User

	if err := db.GetDB(ctx).First(&user, "( signin_phone = ? AND nation_code = ? AND has_set_phone = ? ) ", phoneNumber, nationCode, true).Error; err != nil {
		return 0, err
	}
	return user.UserID, nil
}

// FindUserIDByUsername 通过用户名找到userID
func (db *DbClient) FindUserIDByUsername(ctx context.Context, username string) (int32, error) {
	var user User

	if err := db.GetDB(ctx).First(&user, "( signin_username = ? AND has_set_username = ? ) ", username, true).Error; err != nil {
		return 0, err
	}
	return user.UserID, nil
}

// ExistUsername 用户名是否存在
func (db *DbClient) ExistUsername(ctx context.Context, username string) (bool, error) {
	var count int
	err := db.GetDB(ctx).Model(&User{}).Where("signin_username = ? AND has_set_username = ? ", username, true).Count(&count).Error
	return count != 0, err
}

// HasSetSecureEmailByAnyone 安全邮箱是否已经被任何人设置
func (db *DbClient) HasSetSecureEmailByAnyone(ctx context.Context, email string) (bool, error) {
	var count int
	err := db.GetDB(ctx).Model(&User{}).Where("secure_email = ? AND has_set_email = ?", email, true).Count(&count).Error
	return count != 0, err
}

// ExistPhone 手机号是否已经存在
func (db *DbClient) ExistPhone(ctx context.Context, phone string, nationCode string) (bool, error) {
	var count int
	err := db.GetDB(ctx).Model(&User{}).Where("signin_phone = ? AND nation_code = ? AND has_set_phone = ? ", phone, nationCode, true).Count(&count).Error
	return count != 0, err
}

// ExistSignInPhone 登录手机号是否已经存在
func (db *DbClient) ExistSignInPhone(ctx context.Context, phone string, nationCode string) (bool, error) {
	var count int
	err := db.GetDB(ctx).Model(&User{}).Where("signin_phone = ? AND nation_code = ? AND has_set_phone = ? ", phone, nationCode, true).Count(&count).Error
	return count != 0, err
}

// SecureEmailExists 当前用户是否已经设置了安全邮箱
func (db *DbClient) SecureEmailExists(ctx context.Context, userID int32) (bool, error) {
	var count int
	err := db.GetDB(ctx).Model(&User{}).Where("user_id = ? AND has_set_email = ?", userID, true).Count(&count).Error
	return count != 0, err

}

// MatchSecureEmail 安全邮箱是否与原来邮箱一致
func (db *DbClient) MatchSecureEmail(ctx context.Context, email string, userID int32) (bool, error) {
	var count int
	err := db.GetDB(ctx).Model(&User{}).Where("user_id = ? AND has_set_email = ? AND secure_email = ?", userID, true, email).Count(&count).Error
	return count == 1, err
}

// CreateUserByPhone 创建user通过电话
func (db *DbClient) CreateUserByPhone(ctx context.Context, u *User) (int32, error) {
	if err := db.GetDB(ctx).Create(u).Error; err != nil {
		return 0, err
	}
	return u.UserID, nil
}

// SetSecureEmail 设置安全邮箱
func (db *DbClient) SetSecureEmail(ctx context.Context, email string, userID int32) error {
	return db.GetDB(ctx).Model(&User{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"secure_email":            email,
		"has_set_email":           true,
		"latest_updated_email_at": time.Now().UTC(),
	}).Error
}

// UnsetSecureEmail 解除设置安全邮箱
func (db *DbClient) UnsetSecureEmail(ctx context.Context, userID int32) error {
	return db.GetDB(ctx).Model(&User{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"secure_email":            "",
		"has_set_email":           false,
		"latest_updated_email_at": time.Now().UTC(),
	}).Error
}

// ExistsSecureQuestion 用户是否已经设置了密保问题
func (db *DbClient) ExistsSecureQuestion(ctx context.Context, userID int32) (bool, error) {
	var user User
	if err := db.GetDB(ctx).First(&user, "( user_id = ? AND deleted_at is null ) ", userID).Error; err != nil {
		return false, err
	}
	if !user.HasSetSecureQuestions {
		return false, nil
	}
	return true, nil
}

// SetSecureQuestion 设置密保问题
func (db *DbClient) SetSecureQuestion(ctx context.Context, userID int32, secureQuestion []SecureQuestion) error {
	if len(secureQuestion) != QuestionCount {
		return errors.New("wrong count of question")
	}
	return db.GetDB(ctx).Model(&User{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"secure_question_1":                  secureQuestion[0].SecureQuestionKey,
		"secure_answer_1":                    secureQuestion[0].SecureAnswer,
		"secure_question_2":                  secureQuestion[1].SecureQuestionKey,
		"secure_answer_2":                    secureQuestion[1].SecureAnswer,
		"secure_question_3":                  secureQuestion[2].SecureQuestionKey,
		"secure_answer_3":                    secureQuestion[2].SecureAnswer,
		"has_set_secure_questions":           true,
		"latest_updated_secure_questions_at": time.Now().UTC(),
	}).Error
}

// SetSigninPhoneByUserID 通过userID设置登录手机号
func (db *DbClient) SetSigninPhoneByUserID(ctx context.Context, userID int32, signinPhone string, nationCode string) error {
	now := time.Now()
	return db.GetDB(ctx).Model(&User{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"signin_phone":            signinPhone,
		"has_set_phone":           true,
		"latest_updated_phone_at": now.UTC(),
		"nation_code":             nationCode,
		"updated_at":              now.UTC(),
	}).Error
}

// SetUserRegion 设置用户区域
func (db *DbClient) SetUserRegion(ctx context.Context, userID int32, region Region) error {
	return db.GetDB(ctx).Model(&User{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"has_set_region": true,
		"region":         region,
	}).Error
}

// GetSecureQuestionListToModifyByUserID 通过userID找到密保问题
func (db *DbClient) GetSecureQuestionListToModifyByUserID(ctx context.Context, userID int32) ([]string, error) {
	var user User
	if err := db.GetDB(ctx).First(&user, "( user_id = ? and has_set_secure_questions = ?) ", userID, true).Error; err != nil {
		return nil, err
	}
	return []string{user.SecureQuestion1, user.SecureQuestion2, user.SecureQuestion3}, nil
}

// FindUserBySecureEmail 通过安全邮箱找到User
func (db *DbClient) FindUserBySecureEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := db.GetDB(ctx).First(&user, "( secure_email = ? and has_set_email = ?) ", email, true).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindUsernameBySecureEmail 通过邮箱查找用户名
func (db *DbClient) FindUsernameBySecureEmail(ctx context.Context, email string) (string, error) {
	var user User

	if err := db.GetDB(ctx).First(&user, "( secure_email = ? AND has_set_email = ? ) ", email, true).Error; err != nil {
		return "", err
	}
	if user.HasSetUsername {
		return user.SigninUsername, nil
	}
	return "", errors.New("you have not set the username before")
}

// GetSecureQuestionsByPhone 根据手机号获取当前设置的密保问题
func (db *DbClient) GetSecureQuestionsByPhone(ctx context.Context, nationCode, phone string) ([]string, error) {
	var user User

	if err := db.GetDB(ctx).First(&user, "( signin_phone = ?  AND nation_code = ?  AND has_set_phone = ? AND has_set_secure_questions = ?) ", phone, nationCode, true, true).Error; err != nil {
		return nil, err
	}
	secureQuestions := []string{user.SecureQuestion1, user.SecureQuestion2, user.SecureQuestion3}

	return secureQuestions, nil
}

// GetSecureQuestionsByUsername 根据用户名获取当前设置的密保问题
func (db *DbClient) GetSecureQuestionsByUsername(ctx context.Context, username string) ([]string, error) {
	var user User

	if err := db.GetDB(ctx).First(&user, "( signin_username = ?  AND has_set_username = ? AND has_set_secure_questions = ?) ", username, true, true).Error; err != nil {
		return nil, err
	}
	secureQuestions := []string{user.SecureQuestion1, user.SecureQuestion2, user.SecureQuestion3}

	return secureQuestions, nil
}

// FindUserByEmail 通过邮箱找到User
func (db *DbClient) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := db.GetDB(ctx).First(&user, "( secure_email = ? and has_set_email = ?) ", email, true).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// SetSecureEmailByUserID 根据userID重置安全邮箱
func (db *DbClient) SetSecureEmailByUserID(ctx context.Context, userID int32, email string) error {
	now := time.Now()
	return db.GetDB(ctx).Model(&User{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"secure_email":            email,
		"has_set_email":           true,
		"latest_updated_email_at": now.UTC(),
		"updated_at":              now.UTC(),
	}).Error
}

// ModifyHasSetUserProfileStatus 修改HasSetUserProfile状态
func (db *DbClient) ModifyHasSetUserProfileStatus(ctx context.Context, userID int32) error {
	return db.GetDB(ctx).Model(&User{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"has_set_user_profile": true,
		"is_profile_completed": true,
		"updated_at":           time.Now().UTC(),
	}).Error
}

// HasSecureEmailSet 当前安全邮箱是否被任何人设置
func (db *DbClient) HasSecureEmailSet(ctx context.Context, email string) (bool, error) {
	var count int
	err := db.GetDB(ctx).Model(&User{}).Where("secure_email = ? AND has_set_email = ?", email, true).Count(&count).Error
	return count == 1, err
}

// FindUserByUserID 获取用户和用户档案信息
func (db *DbClient) FindUserByUserID(ctx context.Context, userID int32) (*User, error) {
	var user User
	if err := db.GetDB(ctx).Raw(`SELECT 
    U.signin_username, 
    UP.nickname as nickname,
    case when UP.nickname_initial = '~' THEN '#' ELSE UP.nickname_initial END as nickname_initial,
    UP.gender as gender, 
    UP.birthday as birthday, 
    UP.height as height,
    UP.weight as weight, 
    U.signin_phone,
    U.nation_code,
    U.secure_email,
    U.encrypted_password,
    U.remark,
    U.register_type,
    U.customized_code,
    U.register_time,
    U.user_defined_code, 
    U.has_set_user_profile,
    U.has_set_email,
    U.has_set_phone,
    U.has_set_username,
    U.has_set_password,
	U.has_set_secure_questions,
	U.is_profile_completed,
	(OO.owner_id IS NULL) AS is_removable,
	U.language,
	U.region,
	U.seed
    FROM user AS U
	LEFT JOIN user_profile AS UP ON U.user_id = UP.user_id
	LEFT JOIN organization_owner AS OO ON U.user_id = OO.owner_id
    WHERE U.user_id = ? AND UP.deleted_at IS NULL`, userID).Scan(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// ModifyUser 修改用户信息
func (db *DbClient) ModifyUser(ctx context.Context, user *User) error {
	return db.GetDB(ctx).Model(&User{}).Where("user_id = ?", user.UserID).Updates(map[string]interface{}{
		"customized_code":   user.CustomizedCode,
		"remark":            user.Remark,
		"user_defined_code": user.UserDefinedCode,
		"updated_at":        time.Now().UTC(),
	}).Error
}

// DeleteUser 删除用户
func (db *DbClient) DeleteUser(ctx context.Context, userID int32) error {
	now := time.Now()
	return db.GetDB(ctx).Model(&User{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"updated_at": now.UTC(),
		"deleted_at": now.UTC(),
	}).Error
}
