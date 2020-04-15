package mysqldb

import (
	"context"
	"time"

	tx "github.com/jinmukeji/plat-pkg/v2/store"
)

// Datastore 定义数据访问接口
type Datastore interface {
	tx.Tx // 需要支持事务处理
	// SafeCloseDB 安全的关闭数据库连接
	SafeCloseDB(ctx context.Context)
	// FindClientByClientID 查找一条 Client 数据记录
	FindClientByClientID(ctx context.Context, clientID string) (*Client, error)
	// FindUserByPhone 通过电话找到User
	FindUserByPhone(ctx context.Context, phone string, nationCode string) (*User, error)
	// FindUserByUsername 通过用户名找到User
	FindUserByUsername(ctx context.Context, username string) (*User, error)
	// CreateToken 保存token 并返回一条记录
	CreateToken(ctx context.Context, token string, userID int32, tokenTime time.Duration) (*Token, error)
	// FindUserIDByToken 根据 token 返回 userID，如果token失效返回 error
	FindUserIDByToken(ctx context.Context, token string) (int32, error)
	// SetLanguageByUserID 通过userID设置Language
	SetLanguageByUserID(ctx context.Context, userID int32, language string) error
	// FindLanguageByUserID 通过userID找到Language
	FindLanguageByUserID(ctx context.Context, userID int32) (string, error)
	// ModifyUserProfile 修改用户档案
	ModifyUserProfile(ctx context.Context, profile *UserProfile) error
	// FindUserProfile 找到用户档案
	FindUserProfile(ctx context.Context, userID int32) (*UserProfile, error)
	// ExistUserByUserID 查看 user 能否存在
	ExistUserByUserID(ctx context.Context, userID int32) (bool, error)
	// ExistPasswordByUserID 查看 password 能否存在
	ExistPasswordByUserID(ctx context.Context, userID int32) (bool, error)
	// SetPasswordByUserID 通过userID设置密码
	SetPasswordByUserID(ctx context.Context, userID int32, encryptedPassword string, seed string) error
	// CreateAuditUserCredentialUpdate 新增一个审计记录
	CreateAuditUserCredentialUpdate(ctx context.Context, auditUserCredentialUpdate *AuditUserCredentialUpdate) error
	// ExistSignInPhone 登录手机号是否已经存在

	ExistSignInPhone(ctx context.Context, phone string, nationCode string) (bool, error)
	// ExistPhone 手机号是否已经存在
	ExistPhone(ctx context.Context, phone string, nationCode string) (bool, error)
	// CreateUserByPhone 通过电话创建User
	CreateUserByPhone(ctx context.Context, u *User) (int32, error)
	// GetUserPreferencesByUserID 返回数据库中的用户偏好
	GetUserPreferencesByUserID(ctx context.Context, userID int32) (*UserPreferences, error)
	// DeleteToken 删除数据库内指定的 token
	DeleteToken(ctx context.Context, token string) error
	// CreateUserPreferences 创建UserPreferences
	CreateUserPreferences(ctx context.Context, userID int32) error
	// SetPasswordByPhone 根据手机号重置密码
	SetPasswordByPhone(ctx context.Context, phone string, nationCode string, encryptedPassword string, seed string) error
	// SetPasswordByUsername 根据用户名重置密码
	SetPasswordByUsername(ctx context.Context, username string, encryptedPassword string, seed string) error
	// FindSecureQuestionByPhone 通过电话号码找到密保问题和答案
	FindSecureQuestionByPhone(ctx context.Context, phone string, nationCode string) ([]SecureQuestion, error)
	// FindSecureQuestionByUsername 通过用户名找到密保问题和答案
	FindSecureQuestionByUsername(ctx context.Context, username string) ([]SecureQuestion, error)
	// IsPasswordSameByPhone 根据手机号判断密码是否与之前密码相同
	IsPasswordSameByPhone(ctx context.Context, phone string, nationCode string, encryptedPassword string) (bool, error)
	// IsPasswordSameByUsername 根据用户名判断密码是否与之前密码相同
	IsPasswordSameByUsername(ctx context.Context, username string, encryptedPassword string) (bool, error)
	// FindUserIDByPhone 通过电话号码找到userID
	FindUserIDByPhone(ctx context.Context, phoneNumber string, nationCode string) (int32, error)
	// FindUserIDByUsername 通过用户名找到userID
	FindUserIDByUsername(ctx context.Context, username string) (int32, error)
	// ExistsSecureQuestion 用户是否已经设置了密保问题
	ExistsSecureQuestion(ctx context.Context, userID int32) (bool, error)
	// SetSecureQuestion 设置密保问题
	SetSecureQuestion(ctx context.Context, userID int32, secureQuestion []SecureQuestion) error
	// FindSecureQuestionByUserID 通过userID找到密保问题和答案
	FindSecureQuestionByUserID(ctx context.Context, userID int32) ([]SecureQuestion, error)
	// ExistUsername 用户名是否存在
	ExistUsername(ctx context.Context, username string) (bool, error)
	// HasSetSecureEmailByAnyone 安全邮箱是否已经被任何人设置
	HasSetSecureEmailByAnyone(ctx context.Context, email string) (bool, error)
	// SecureEmailExists 当前用户已经设置了安全邮箱
	SecureEmailExists(ctx context.Context, userID int32) (bool, error)
	// MatchSecureEmail 安全邮箱是否与原来邮箱一致
	MatchSecureEmail(ctx context.Context, email string, userID int32) (bool, error)
	//SetSecureEmail 设置安全邮箱
	SetSecureEmail(ctx context.Context, email string, userID int32) error
	// UnsetSecureEmail 解除设置安全邮箱
	UnsetSecureEmail(ctx context.Context, userID int32) error
	// CreateVcRecord 创建验证码记录
	CreateVcRecord(ctx context.Context, record *VcRecord) error
	// SearchVcRecordCountsIn24hours 搜索24小时内的验证码记录个数
	SearchVcRecordCountsIn24hours(ctx context.Context, sendTo string) (int, error)
	// SearchVcRecordCountsIn1Minute 搜索1分钟内的验证码记录
	SearchVcRecordCountsIn1Minute(ctx context.Context, sendTo string) (int, error)
	// SearchVcRecordEarliestTimeIn1Minute 搜索1分钟最早的验证码记录时间
	SearchVcRecordEarliestTimeIn1Minute(ctx context.Context, sendTo string) (*VcRecord, error)
	// FindVcRecord 查找验证码记录
	FindVcRecord(ctx context.Context, sn string, vc string, sendTo string, usage Usage) (*VcRecord, error)
	// ModifyVcRecordStatus 修改验证码的状态
	ModifyVcRecordStatus(ctx context.Context, recordID int32) error
	// SearchSpecificVcRecordCountsIn24hours 搜索24小时内特定模板的验证码记录个数
	SearchSpecificVcRecordCountsIn24hours(ctx context.Context, sendTo string, usage Usage) (int, error)
	// SearchSpecificVcRecordEarliestTimeIn24hours 搜索24小时指定模板最早的验证码记录
	SearchSpecificVcRecordEarliestTimeIn24hours(ctx context.Context, sendTo string, usage Usage) (*VcRecord, error)
	// FindLatestVcRecord 查找最新验证码记录
	FindLatestVcRecord(ctx context.Context, sendTo string, usage Usage) (*VcRecord, error)
	// SetSigninPhoneByUserID 设置登录手机号
	SetSigninPhoneByUserID(ctx context.Context, userID int32, signinPhone string, nationCode string) error
	// SetUserRegion 设置选择区域
	SetUserRegion(ctx context.Context, userID int32, region Region) error
	// VerifyMVC 验证手机号
	VerifyMVC(ctx context.Context, sn string, vc string, sendTo string, nationCode string) (bool, error)
	// SearchVcRecord 搜索VcRecord
	SearchVcRecord(ctx context.Context, sn string, vc string, sendTo string, nationCode string) (*VcRecord, error)
	// VerifyVerificationNumber 验证 VerificationNumber是否有效
	VerifyVerificationNumber(ctx context.Context, verificationType VerificationType, verificationNumber string, UserID int32) (bool, error)
	// CreatePhoneOrEmailVerfication 创建手机邮箱验证记录
	CreatePhoneOrEmailVerfication(ctx context.Context, record *PhoneOrEmailVerfication) error
	// SetVerificationNumberAsUsed 设置VerificationNumber已经使用
	SetVerificationNumberAsUsed(ctx context.Context, verificationType VerificationType, verificationNumber string) error
	// VerifyVerificationNumberByPhone 手机号验证 VerificationNumber是否有效
	VerifyVerificationNumberByPhone(ctx context.Context, verificationNumber string, phone, nationCode string) (bool, error)
	// SearchLatestPhoneVerificationCode 搜索最新的电话验证码
	SearchLatestPhoneVerificationCode(ctx context.Context, sendTo string, nationCode string) (string, error)
	// SearchLatestEmailVerificationCode 搜索最新的邮件验证码
	SearchLatestEmailVerificationCode(ctx context.Context, sendTo string) (string, error)
	// FindUserBySecureEmail 通过安全邮箱找到User
	FindUserBySecureEmail(ctx context.Context, email string) (*User, error)
	// GetSecureQuestionListToModifyByUserID 通过userID找到密保问题Key
	GetSecureQuestionListToModifyByUserID(ctx context.Context, userID int32) ([]string, error)
	// FindUsernameBySecureEmail 通过邮箱查找用户名
	FindUsernameBySecureEmail(ctx context.Context, email string) (string, error)
	// GetSecureQuestionsByPhone 根据手机号获取当前设置的密保问题
	GetSecureQuestionsByPhone(ctx context.Context, nationCode, phone string) ([]string, error)
	// GetSecureQuestionsByUsername 根据用户名获取当前设置的密保问题
	GetSecureQuestionsByUsername(ctx context.Context, username string) ([]string, error)
	// FindUserByEmail 通过邮箱找到User
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	// VerifyVerificationNumberByEmail 邮箱 VerificationNumber是否有效
	VerifyVerificationNumberByEmail(ctx context.Context, verificationNumber string, email string) (bool, error)
	// VerifyMVCBySecureEmail 根据安全邮箱寻找验证码信息
	VerifyMVCBySecureEmail(ctx context.Context, sn string, vc string, email string) (*VcRecord, error)
	// SetSecureEmailByUserID 根据userID重置安全邮箱
	SetSecureEmailByUserID(ctx context.Context, userID int32, email string) error
	// ModifyVcRecordStatusByEmail 根据安全邮箱修改验证码的状态
	ModifyVcRecordStatusByEmail(ctx context.Context, email, verificationCode, serialNumber string) error
	// SetMVCAsUsed 设置vc为使用过的
	SetVcAsUsed(ctx context.Context, sn string, vc string, sendTo string, nationCode string) error
	// ModifyHasSetUserProfileStatus 修改HasSetUserProfile的状态
	ModifyHasSetUserProfileStatus(ctx context.Context, userID int32) error
	// HasSecureEmailSet 当前安全邮箱是否被任何人设置
	HasSecureEmailSet(ctx context.Context, email string) (bool, error)
	// CreateAuditUserSigninSignout 新增登录/登出审计记录
	CreateAuditUserSigninSignout(ctx context.Context, auditUserSigninSignout *AuditUserSigninSignout) error
	// FindUserProfileByRecordID 通过记录ID获取用户档案
	FindUserProfileByRecordID(ctx context.Context, recordID int32) (*UserProfile, error)
	// FindUserByUserID 获取用户和档案信息
	FindUserByUserID(ctx context.Context, userID int32) (*User, error)
	// CreateUserProfile 创建UserProfile
	CreateUserProfile(ctx context.Context, u *UserProfile) (*UserProfile, error)
	// ModifyUser 修改用户信息
	ModifyUser(ctx context.Context, user *User) error
	// FindUsingClients 得到正在使用的客户端
	FindUsingClients(ctx context.Context, userID int32) ([]Client, error)
	// FindUserSigninRecord 查询用户登录记录
	FindUserSigninRecord(ctx context.Context, userID int32) ([]AuditUserSigninSignout, error)
	// HasUserSetNotificationPreferences 查看用户是否设置通知配置首选项
	HasUserSetNotificationPreferences(ctx context.Context, userID int32) (bool, error)
	// CreateNotificationPreferences 创建通知配置首选项
	CreateNotificationPreferences(ctx context.Context, notificationPreferences *NotificationPreferences) error
	// UpdateNotificationPreferences 更新通知配置首选项
	UpdateNotificationPreferences(ctx context.Context, notificationPreferences *NotificationPreferences) error
	// GetNotificationPreferences 获取通知配置首选项
	GetNotificationPreferences(ctx context.Context, userID int32) (*NotificationPreferences, error)
	// DeleteTokenByUserID 通过用户ID删除token
	DeleteTokenByUserID(ctx context.Context, userID int32) error
	// CreateUserAndUserProfile 创建用户和用户资料
	CreateUserAndUserProfile(ctx context.Context, user *User, userProfile *UserProfile) error
	// SearchVcRecordFrom1MinuteTo2Mintue
	SearchVcRecordFrom1MinuteTo2Mintue(ctx context.Context, sendTo string, usage Usage) (*VcRecord, error)
	// DeleteUser 删除用户
	DeleteUser(ctx context.Context, userID int32) error
}
