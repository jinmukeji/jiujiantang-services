package mysqldb

import (
	"context"
	"time"
)

// VerificationType 验证VerificationType
type VerificationType string

const (
	// VerificationPhone 验证手机
	VerificationPhone VerificationType = "phone"
	// VerificationEmail 验证邮箱
	VerificationEmail VerificationType = "email"
)

// PhoneOrEmailVerfication 手机邮箱验证
type PhoneOrEmailVerfication struct {
	RecordID           int32            `gorm:"primary_key"`                // 记录id
	VerificationType   VerificationType `gorm:"column:verification_type"`   // 验证类型
	VerificationNumber string           `gorm:"column:verification_number"` // 验证号
	SendTo             string           `gorm:"column:send_to"`             // 接收人
	NationCode         string           `gorm:"column:nation_code"`         // 国家代码
	UserID             int32            `grom:"column:user_id"`             // 用户ID
	ExpiredAt          *time.Time       `gorm:"column:expired_at"`          // 到期时间
	HasUsed            bool             `gorm:"column:has_used"`            // 是否使用过
	CreatedAt          time.Time        // 创建时间
	UpdatedAt          time.Time        // 更新时间
	DeletedAt          *time.Time       // 删除时间
}

// TableName 返回 VcRecord 所在的表名
func (v PhoneOrEmailVerfication) TableName() string {
	return "phone_or_email_verfication"
}

// CreatePhoneOrEmailVerfication 创建手机邮箱验证记录
func (db *DbClient) CreatePhoneOrEmailVerfication(ctx context.Context, record *PhoneOrEmailVerfication) error {
	return db.GetDB(ctx).Create(record).Error
}

// VerifyVerificationNumber 验证 VerificationNumber是否有效
func (db *DbClient) VerifyVerificationNumber(ctx context.Context, verificationType VerificationType, verificationNumber string, UserID int32) (bool, error) {
	var count int
	err := db.GetDB(ctx).Raw(`SELECT count(*) FROM phone_or_email_verfication 
	    where verification_type = ? and verification_number = ? and user_id = ? 
		AND expired_at > NOW() and has_used = false;`, verificationType, verificationNumber, UserID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 1, nil
}

// SetVerificationNumberAsUsed 设置VerificationNumber已经使用
func (db *DbClient) SetVerificationNumberAsUsed(ctx context.Context, verificationType VerificationType, verificationNumber string) error {
	return db.GetDB(ctx).Model(&PhoneOrEmailVerfication{}).Where("verification_type = ? and verification_number = ? ", verificationType, verificationNumber).Updates(map[string]interface{}{
		"has_used": true,
	}).Error
}

// VerifyVerificationNumberByPhone 手机号验证 VerificationNumber是否有效
func (db *DbClient) VerifyVerificationNumberByPhone(ctx context.Context, verificationNumber string, phone, nationCode string) (bool, error) {
	var count int
	err := db.GetDB(ctx).Raw(`SELECT count(*) FROM phone_or_email_verfication 
	    where send_to = ? and verification_number = ? and nation_code = ? 
		AND expired_at > NOW() and has_used = false;`, phone, verificationNumber, nationCode).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 1, nil
}

// VerifyVerificationNumberByEmail 邮箱 VerificationNumber是否有效
func (db *DbClient) VerifyVerificationNumberByEmail(ctx context.Context, verificationNumber string, email string) (bool, error) {
	var count int
	err := db.GetDB(ctx).Raw(`SELECT count(*) FROM phone_or_email_verfication 
	    where send_to = ? AND verification_number = ? 
		AND expired_at > NOW() and has_used = false;`, email, verificationNumber).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 1, nil
}
