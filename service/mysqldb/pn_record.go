package mysqldb

import (
	"context"
	"time"
)

// PnRecord 通知的阅读记录
type PnRecord struct {
	PnID      int32      `gorm:"pn_id"`
	UserID    int32      `gorm:"user_id"`
	CreatedAt time.Time  // 创建时间
	UpdatedAt time.Time  // 更新时间
	DeletedAt *time.Time // 删除时间
}

// TableName 返回 QRCode 所在的表名
func (pr PnRecord) TableName() string {
	return "pn_record"
}

// CreatePnRecord 创建通知记录
func (db *DbClient) CreatePnRecord(ctx context.Context, pr *PnRecord) error {
	if err := db.Create(pr).Error; err != nil {
		return err
	}
	return nil
}

// ExistPnRecord 是否已经存在PnRecord
func (db *DbClient) ExistPnRecord(ctx context.Context, pnID int32, UserID int32) (bool, error) {
	var count int
	err := db.Model(&PnRecord{}).Where("pn_id = ? and user_id= ?", pnID, UserID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count != 0, nil
}
