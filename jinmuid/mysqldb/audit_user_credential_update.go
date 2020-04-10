package mysqldb

import (
	"context"
	"time"
)

// UpdatedRecordType 更新记录类型
type UpdatedRecordType string

const (
	// UsernameUpdated 用户名更新
	UsernameUpdated UpdatedRecordType = "username"
	// PhoneUpdated 手机号更新
	PhoneUpdated UpdatedRecordType = "phone"
	// EmailUpdated 邮件更新
	EmailUpdated UpdatedRecordType = "email"
	// PasswordUpdated 密码更新
	PasswordUpdated UpdatedRecordType = "password"
)

// AuditUserCredentialUpdate 审核用户凭证书更新
type AuditUserCredentialUpdate struct {
	RecordID          string            `gorm:"primary_key"`         // record_id 表主键
	UserID            int32             `gorm:"user_id"`             // 用户ID
	ClientID          string            `gorm:"client_id"`           // 客户端ID
	UpdatedRecordType UpdatedRecordType `gorm:"updated_record_type"` // 更新记录类型
	OldValue          string            `grom:"old_value"`           // 修改前的值
	NewValue          string            `grom:"new_value"`           // 修改后的值
	CreatedAt         time.Time         // 创建时间
	UpdatedAt         time.Time         // 更新时间
	DeletedAt         *time.Time
}

// TableName 返回表名
func (a AuditUserCredentialUpdate) TableName() string {
	return "audit_user_credential_update"
}

// CreateAuditUserCredentialUpdate 新增一个审计记录
func (db *DbClient) CreateAuditUserCredentialUpdate(ctx context.Context, auditUserCredentialUpdate *AuditUserCredentialUpdate) error {
	return db.DB(ctx).Create(auditUserCredentialUpdate).Error
}
