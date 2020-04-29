package mysqldb

import (
	"context"
	"time"
)

// RecordType 记录类型
type RecordType string

const (
	// SigninRecordType 登录记录类型
	SigninRecordType RecordType = "signin"
	// SignoutRecordType 登出记录类型
	SignoutRecordType RecordType = "signout"
)

// AuditUserSigninSignout 登录/登出审计表
type AuditUserSigninSignout struct {
	RecordID      string     `gorm:"primary_key;column:record_id"` // record_id 表主键
	UserID        int32      `gorm:"column:user_id"`               // 用户ID
	ClientID      string     `gorm:"column:client_id"`             // 客户端ID
	RecordType    RecordType `gorm:"column:record_type"`           // 更新记录类型
	IP            string     `gorm:"column:ip"`                    // ip
	ExtraParams   string     `gorm:"column:extra_params"`          // 参数
	SignInMachine string     `gorm:"column:sign_in_machine"`       // 登录设备
	CreatedAt     time.Time  // 创建时间
	UpdatedAt     time.Time  // 更新时间
	DeletedAt     *time.Time // 删除时间
}

// TableName 返回表名
func (a AuditUserSigninSignout) TableName() string {
	return "audit_user_signin_signout"
}

// CreateAuditUserSigninSignout 新增登录/登出审计记录
func (db *DbClient) CreateAuditUserSigninSignout(ctx context.Context, auditUserSigninSignout *AuditUserSigninSignout) error {
	return db.GetDB(ctx).Create(auditUserSigninSignout).Error
}

// FindUsingClients 寻找正在使用的客户端
func (db *DbClient) FindUsingClients(ctx context.Context, userID int32) ([]Client, error) {
	var clients []Client
	err := db.GetDB(ctx).Raw("SELECT distinct(AU.client_id),C.remark,C.`usage` FROM audit_user_signin_signout as AU inner join `client` as C on C.client_id = AU.client_id where AU.user_id = ? and AU.client_id <> '' and AU.deleted_at IS NULL", userID).Scan(&clients).Error
	if err != nil {
		return nil, err
	}
	return clients, nil
}

// FindUserSigninRecord 查询登录记录
func (db *DbClient) FindUserSigninRecord(ctx context.Context, userID int32) ([]AuditUserSigninSignout, error) {
	var auditUserSigninSignout []AuditUserSigninSignout
	err := db.GetDB(ctx).Raw(`
    SELECT sign_in_machine,created_at from (
        SELECT AU.sign_in_machine as sign_in_machine, min(AU.created_at) as created_at
            FROM audit_user_signin_signout as AU
            where AU.user_id = ? and AU.record_type = 'signin' and AU.deleted_at IS NULL and AU.sign_in_machine <> ''
            group by AU.sign_in_machine
            ) as a order by created_at
    `, userID).Scan(&auditUserSigninSignout).Error
	if err != nil {
		return nil, err
	}
	return auditUserSigninSignout, nil
}
