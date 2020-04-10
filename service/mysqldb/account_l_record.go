package mysqldb

import (
	"context"
	"time"
)

// AccountLRecord 是一体机账户与记录关联表
type AccountLRecord struct {
	RecordID  int32     `gorm:"primary_key"` // 记录
	Account   string    `gorm:"account"`     // 账户
	CreatedAt time.Time // 创建时间
	UpdatedAt time.Time // 更新时间
	DeletedAt *time.Time
}

// TableName 返回表名
func (a AccountLRecord) TableName() string {
	return "account_l_record"
}

// CreateAccountLRecord 创建一体机账户与记录关联表
func (db *DbClient) CreateAccountLRecord(ctx context.Context, account string, recordID int32) error {
	now := time.Now()
	record := &AccountLRecord{
		RecordID:  recordID,
		Account:   account,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := db.Create(&record).Error; err != nil {
		return err
	}
	return nil
}
