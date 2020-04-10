package mysqldb

import (
	"context"
	"time"
)

// Datastore 是数据访问的接口定义
type Datastore interface {
	// FindValidRecordByID 查找指定 ID 的一条 Record 数据记录
	FindValidRecordByID(ctx context.Context, recordID int) (*Record, error)

	// SetRecordSentWxViewReportNotification 设置是否已经发送过微信查看报告通知的状态标记 0: 未发送，1: 已发送
	SetRecordSentWxViewReportNotification(ctx context.Context, recordID int, sent bool) error

	// FindRecordByID 查找指定 ID 的一条 Record 数据记录
	FindRecordByID(ctx context.Context, recordID int) (*Record, error)

	// CreateToken 保存token 并返回一条记录
	CreateToken(ctx context.Context, token string, userID int32, tokenTime time.Duration) (*Token, error)

	// FindUserIDByToken 根据 token 返回 userID，如果token失效返回 error
	FindUserIDByToken(ctx context.Context, token string) (int32, error)

	// DeleteToken 删除数据库里面的指定的token 失效
	DeleteToken(ctx context.Context, token string) error

	// UpdateUserProfile 更新用户个人信息
	UpdateUserProfile(ctx context.Context, profile ProtoUserProfile, userID int32) error

	// UpdateSimpleUserProfile 简易更新用户个人信息
	UpdateSimpleUserProfile(ctx context.Context, user *User) error

	// CreateFeedback 新增一个用户反馈
	CreateFeedback(ctx context.Context, feedback *Feedback) error

	// FindFeedbackByFeedBackID 查找一个用户反馈
	FindFeedbackByFeedBackID(ctx context.Context, feedbackID int) (*Feedback, error)

	// UpdateRemarkByRecordID 更新备注
	UpdateRemarkByRecordID(ctx context.Context, recordID int, comment string) error

	// FindValidRecordsByDateRange 返回给定时间范围内的 record
	FindValidRecordsByDateRange(ctx context.Context, subjectID int, start time.Time, end time.Time) ([]Record, error)

	// CreateRecord 增加测量记录
	CreateRecord(ctx context.Context, record *Record) error

	// CheckUserHasRecord 验证用户和测量记录是否有关联
	CheckUserHasRecord(ctx context.Context, accountID int32, recordID int32) (bool, error)

	// FindValidPaginatedRecordsByDateRange 返回给定时间和给定数量范围内的有效 record
	FindValidPaginatedRecordsByDateRange(ctx context.Context, subjectID, offset, size int, start, end time.Time) ([]Record, error)

	// FindClientByClientID 查找一条 Client 数据记录
	FindClientByClientID(ctx context.Context, clientID string) (*Client, error)

	// FindUserByUsername 返回数据库中的用户
	FindUserByUsername(ctx context.Context, username string) (*User, error)

	// FindFirstOrganizationByOwner 查找指定 owner 拥有的第一个组织
	FindFirstOrganizationByOwner(ctx context.Context, userID int) (*Organization, error)

	// CreateUser 创建一个新用户
	CreateUser(ctx context.Context, u *User) (*User, error)

	// FindUserByUseID 返回数据库中的用户
	FindUserByUserID(ctx context.Context, userID int) (*User, error)

	// BindDeviceToOrganization 关联设备到组织
	BindDeviceToOrganization(ctx context.Context, d *DeviceOrganizationBinding) error

	// CreateOrganization 创建组织
	CreateOrganization(ctx context.Context, o *Organization) error

	// CreateOrganizationOwner 在 organization_owner 新增一条记录
	CreateOrganizationOwner(ctx context.Context, o *OrganizationOwner) error

	// CheckOrganizationOwner 检查用户是否为组织的拥有者
	CheckOrganizationOwner(ctx context.Context, userID int, organizationID int) (bool, error)

	// FindOrganizationsByOwner 查找指定 owner 拥有的所有
	FindOrganizationsByOwner(ctx context.Context, ownerID int) ([]*Organization, error)

	// UpdateOrganizationProfile 更新组织信息
	UpdateOrganizationProfile(ctx context.Context, o *Organization) error

	// FindOrganizationByID 从组织 ID 查找组织
	FindOrganizationByID(ctx context.Context, organizationID int) (*Organization, error)

	// DeleteOrganizationByID 删除组织
	DeleteOrganizationByID(ctx context.Context, organizationID int) error

	// CreateOrganizationUsers 在 organization_user 新增多条记录
	CreateOrganizationUsers(ctx context.Context, users []*OrganizationUser) error

	// DeleteOrganizationUser 在 organization_user 删除一条记录
	DeleteOrganizationUser(ctx context.Context, userID, organizationID int) error

	// CheckOrganizationUser 检查 user 是否为组织下用户
	CheckOrganizationUser(ctx context.Context, userID, organizationID int) (bool, error)

	// DeleteOrganizationUsers 在 organization_user 删除多条记录
	DeleteOrganizationUsers(ctx context.Context, userIDList []int32, organizationID int32) error

	// FindOrganizationUsers 查看组织下用户
	FindOrganizationUsers(ctx context.Context, organizationID int) ([]*User, error)

	// 创建订阅
	CreateSubscription(ctx context.Context, s *Subscription) error

	// GetExistingUserCountByOrganizationID 查找指定 user的数量
	GetExistingUserCountByOrganizationID(ctx context.Context, organizationID int) (int, error)

	// FindSubscriptionsByOrganizationID 查找 Subscription 通过 OrganizationID
	FindSubscriptionsByOrganizationID(ctx context.Context, organizationID int) ([]*Subscription, error)

	// GetOrganizationCountByOwnerID 查找组织的数量
	GetOrganizationCountByOwnerID(ctx context.Context, ownerID int) (int, error)

	// ExistOrganizationDeviceByID 检查组织和deviceID关联是否存在
	ExistOrganizationDeviceByID(ctx context.Context, organizationID int32, deviceID int) (bool, error)

	// UnbindOrganizationDevice 解除组织和Device的关联关系
	UnbindOrganizationDevice(ctx context.Context, organizationID int32, deviceID int) error

	// GetOrganizationDeviceList 通过 organizationID 查询与 Device 的关联关系
	GetOrganizationDeviceList(ctx context.Context, organizationID int32) ([]*DeviceOrganizationBinding, error)

	// GetUserIDByRecordID 通过 RecordID 得到 UserID
	GetUserIDByRecordID(ctx context.Context, recordID int32) (int32, error)

	// FindOrganizationByUserID 通过UserID找到组织
	FindOrganizationByUserID(ctx context.Context, userID int) (*Organization, error)

	// GetDeviceByDeviceID 通过 DeviceID 查询 Device
	GetDeviceByDeviceID(ctx context.Context, deviceID int) (*Device, error)

	// GetDeviceByMac 通过 mac 查询 Device
	GetDeviceByMac(ctx context.Context, mac uint64) (*Device, error)

	// ExistDeviceByMac 查看 mac 能否存在
	ExistDeviceByMac(ctx context.Context, mac uint64) (bool, error)

	// GetUserPreferencesByUserID 返回数据库中的用户偏好
	GetUserPreferencesByUserID(ctx context.Context, userID int32) (*UserPreferences, error)

	// ExistOrganizationDeviceByDeviceID 检查组织和deviceID关联是否存在
	ExistOrganizationDeviceByDeviceID(ctx context.Context, deviceID int) (bool, error)

	// FindOrganizationUsersByKeyword 通过keyword和分页搜索用户
	FindOrganizationUsersByKeyword(ctx context.Context, organizationID int32, keyword string, size int32, offset int32) ([]*User, error)

	// FindOrganizationUsersByOffset 通过分页搜索用户
	FindOrganizationUsersByOffset(ctx context.Context, organizationID int32, size int32, offset int32) ([]*User, error)

	// CheckOrganizationIsValid 检查组织是否有效
	CheckOrganizationIsValid(ctx context.Context, organizationID int) bool

	// UpdateRecordStatus 更新状态
	UpdateRecordStatus(ctx context.Context, record *Record) error

	// FindTransactionNumberByCurrentDate 查询流水号
	FindTransactionNumberByCurrentDate(ctx context.Context) (*TransactionNumber, error)

	// IsExistTransactionNumberByCurrentDate 流水号是否存在
	IsExistTransactionNumberByCurrentDate(ctx context.Context) (bool, error)

	// CreateTransactionNumber 创建流水号
	CreateTransactionNumber(ctx context.Context) (*TransactionNumber, error)

	// GetUserByRecordID 通过RecordID 获取 User
	GetUserByRecordID(ctx context.Context, recordID int32) (*User, error)

	// ExistRecordByRecordID 检查Record是否存在
	ExistRecordByRecordID(ctx context.Context, recordID int32) (bool, error)

	// UpdateRecordAnswers 更新状态
	UpdateRecordAnswers(ctx context.Context, recordID int32, answers string) error

	// CountUserByCustomizedCode 通过customizedCode算出User的数量
	CountUserByCustomizedCode(ctx context.Context, customizedCode string) (int, error)

	// ActivateSubscription 更新订阅
	ActivateSubscription(ctx context.Context, s *Subscription) error

	// CreateQRCode 创建二维码
	CreateQRCode(ctx context.Context, qrcode *QRCode) (*QRCode, error)

	// UpdateQRCode 更新二维码信息
	UpdateQRCode(ctx context.Context, qrcode *QRCode) error

	// CreateScannedQRCodeRecord 创建二维码扫码记录
	CreateScannedQRCodeRecord(ctx context.Context, record *ScannedQRCodeRecord) error

	// FindJinmuLAccountByToken 根据 token 返回 account，如果token失效返回 error
	FindJinmuLAccountByToken(ctx context.Context, token string) (string, error)

	// FindJinmuLAccount 查找account信息
	FindJinmuLAccount(ctx context.Context, account string) (*JinmuLAccount, error)

	// CreateAccessToken 保存token 并立即返回创建的token
	CreateAccessToken(ctx context.Context, token string, account string, machineUUID string, availableDuration time.Duration) (*JinmuLAccessToken, error)

	// FindMachineUUIDByToken 根据 token 返回 Machine UUID，如果token失效返回 error
	FindMachineUUIDByToken(ctx context.Context, token string) (string, error)

	// ExistWXUser 存在该微信用户
	ExistWXUser(ctx context.Context, unionID string) (bool, error)

	// CreateWXUser 创建微信用户
	CreateWXUser(ctx context.Context, wXUser *WXUser) error

	// FindWXUserByUnionID通过UnionId找WXUser
	FindWXUserByUnionID(ctx context.Context, UnionID string) (*WXUser, error)

	// DeleteJinmuLAccessToken 删除数据库里面的指定的token 失效
	DeleteJinmuLAccessToken(ctx context.Context, token string) error

	// FindWXUserByUserID 通过userID找WXUser
	FindWXUserByUserID(ctx context.Context, userID int32) (*WXUser, error)

	// FindWXUserByOpenID 通过openID找WXUser
	FindWXUserByOpenID(ctx context.Context, openID string) (*WXUser, error)

	// FindValidPaginatedRecordsByUserID 通过UserID拿到records
	FindValidPaginatedRecordsByUserID(ctx context.Context, userID int32) ([]Record, error)

	// UpdateRecordHasPaid 更新record的hasPaid
	UpdateRecordHasPaid(ctx context.Context, recordID int32) error

	// UpdateUserProfileContainGender 更新含有Gender的用户信息
	UpdateUserProfileContainGender(ctx context.Context, profile ProtoUserProfile, userID int32, gender string) error

	// FindQRCodeBySceneID 通过SceneID拿到QRCode
	FindQRCodeBySceneID(ctx context.Context, SceneID int32) (*QRCode, error)

	// CreateSession 创建Session
	CreateSession(ctx context.Context, s *Session) (*Session, error)

	// FindSession 寻找Session
	FindSession(ctx context.Context, sessionID string) (*Session, error)

	// UpdateSession 更新Session
	UpdateSession(ctx context.Context, sessionID string, session *Session) error

	// UpdateRecordToken 更新record_token
	UpdateRecordToken(ctx context.Context, recordID int32, token string) error

	// FindRecordByToken 查找指定 token 的一条 Record 数据记录
	FindRecordByToken(ctx context.Context, token string) (*Record, error)

	// ExistPnRecord 是否已经存在PnRecord
	ExistPnRecord(ctx context.Context, pnID int32, UserID int32) (bool, error)

	// CreatePnRecord 创建通知记录
	CreatePnRecord(ctx context.Context, pr *PnRecord) error

	// GetPnsByUserID 通过userID拿到未读通知记录，按时间倒序
	GetPnsByUserID(ctx context.Context, UserID int32, size int32) ([]PushNotification, error)

	// UpdateRecordHasAEError 更新记录AE是否存在异常
	UpdateRecordHasAEError(ctx context.Context, recordID int32, hasAEError int32) error

	// CreateLocalNotification 创建本地消息推送
	CreateLocalNotification(ctx context.Context, record *LocalNotification) error

	//GetLocalNotification 得到本地消息推送
	GetLocalNotifications(ctx context.Context) ([]LocalNotification, error)

	// DeleteLocalNotification 删除本地消息推送
	DeleteLocalNotification(ctx context.Context, lnID int) error

	// CreateAccountLRecord 创建一体机账户与记录关联表
	CreateAccountLRecord(ctx context.Context, account string, recordID int32) error

	// FindOrganizationUsersIDList 查看组织下用户ID的List
	FindOrganizationUsersIDList(ctx context.Context, organizationID int) ([]int, error)

	// UpdateRecordTransactionNumber 更新记录的流水号
	UpdateRecordTransactionNumber(ctx context.Context, recordID int32, transactionNumber string) error

	// GetIsRemovableStatus 获取用户是否可移除的状态
	GetIsRemovableStatus(ctx context.Context, username string) (bool, error)

	// DeleteRecord 删除记录
	DeleteRecord(ctx context.Context, recordID int32) error

	// CheckUserOwnerBelongToSameOrganization 检查用户和拥有这是否处于相同的组织下
	CheckUserOwnerBelongToSameOrganization(ctx context.Context, userID int32, ownerID int32) (bool, error)

	// GetOwnerIDByOrganizationID 根据组织的ID获取组织拥有者的ID
	GetOwnerIDByOrganizationID(ctx context.Context, organizationID int) ([]int, error)
}
