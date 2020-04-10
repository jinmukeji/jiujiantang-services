package mysqldb

import (
	"context"
	"time"
)

// QRCode 二维码
type QRCode struct {
	SceneID     int32      `gorm:"primary_key"`         // 场景ID
	RawURL      string     `gorm:"column:raw_url"`      // 原始URL
	Ticket      string     `gorm:"column:ticket"`       // Ticket
	Account     string     `gorm:"column:account"`      // account
	MachineUUID string     `gorm:"column:machine_uuid"` // machine_uuid
	OriginID    string     `gorm:"column:origin_id"`    // originID
	ExpiredAt   time.Time  // 二维码的过期时间
	CreatedAt   time.Time  // 创建时间
	UpdatedAt   time.Time  // 更新时间
	DeletedAt   *time.Time // 删除时间
}

// TableName 返回 QRCode 所在的表名
func (q QRCode) TableName() string {
	return "wxmp_tmp_qrcode"
}

// CreateQRCode 创建二维码信息
func (db *DbClient) CreateQRCode(ctx context.Context, qrcode *QRCode) (*QRCode, error) {
	if err := db.Create(qrcode).Error; err != nil {
		return nil, err
	}
	return qrcode, nil
}

// UpdateQRCode 更新二维码信息
func (db *DbClient) UpdateQRCode(ctx context.Context, qrcode *QRCode) error {
	return db.Model(&QRCode{}).Where("scene_id = ?", qrcode.SceneID).Updates(map[string]interface{}{
		"raw_url":      qrcode.RawURL,
		"ticket":       qrcode.Ticket,
		"account":      qrcode.Account,
		"machine_uuid": qrcode.MachineUUID,
		"expired_at":   qrcode.ExpiredAt,
		"ori_id":       qrcode.OriginID,
		"updated_at":   qrcode.UpdatedAt,
	}).Error
}

// FindQRCodeBySceneID 通过SceneID拿到QRCode
func (db *DbClient) FindQRCodeBySceneID(ctx context.Context, SceneID int32) (*QRCode, error) {
	var qrcode QRCode
	err := db.First(&qrcode, "scene_id = ? ", SceneID).Error
	if err != nil {
		return nil, err
	}
	return &qrcode, nil
}
