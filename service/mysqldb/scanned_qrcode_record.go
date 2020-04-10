package mysqldb

import (
	"context"
	"time"
)

// ScannedQRCodeRecord 二维码扫码记录
type ScannedQRCodeRecord struct {
	RecordID  int32      `gorm:"primary_key"` // 记录ID
	SceneID   int32      `gorm:"scene_id"`    // 场景ID
	CreatedAt time.Time  // 创建时间
	UpdatedAt time.Time  // 更新时间
	DeletedAt *time.Time // 删除时间
}

// TableName 返回 QRCodeRecord 所在的表名
func (q ScannedQRCodeRecord) TableName() string {
	return "scanned_qrcode_record"
}

// CreateScannedQRCodeRecord 创建二维码扫码记录
func (db *DbClient) CreateScannedQRCodeRecord(ctx context.Context, record *ScannedQRCodeRecord) error {
	return db.Create(record).Error
}
