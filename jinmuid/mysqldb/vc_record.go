package mysqldb

import (
	"context"
	"time"
)

// SendMedia 发送方式
type SendMedia string

const (
	// Phone 发送方式为电话号码
	Phone SendMedia = "phone"
	// Email 发送方式为邮件
	Email SendMedia = "email"
)

// Usage 用途
type Usage string

const (
	// Unknown 未知
	Unknown Usage = "Unknown"
	// FindResetPassword 重置密码
	FindResetPassword Usage = "FindResetPassword"
	// FindUsername 找到用户名
	FindUsername Usage = "FindUsername"
	// SetSecureEmail 设置安全邮箱
	SetSecureEmail Usage = "SetEmail"
	// ModifySecureEmail 修改安全邮箱
	ModifySecureEmail Usage = "ModifyEmail"
	// UnsetSecureEmail 解绑安全邮箱
	UnsetSecureEmail Usage = "UnsetEmail"
	// SignUp 注册
	SignUp Usage = "SignUp"
	// SignIn 登陆
	SignIn Usage = "SignIn"
	// ResetPassword 重置密码
	ResetPassword Usage = "ResetPassword"
	// SetPhoneNumber 设置手机号
	SetPhoneNumber Usage = "SetPhoneNumber"
	// ModifyPhoneNumber 修改手机号
	ModifyPhoneNumber Usage = "ModifyPhoneNumber"
)

// VcRecord 验证码记录
type VcRecord struct {
	RecordID   int32      `gorm:"primary_key"`       // 记录id
	Usage      Usage      `gorm:"column:usage"`      // 使用用途
	SN         string     `gorm:"column:sn"`         // 序列号
	Code       string     `gorm:"column:code"`       // 验证码
	SendVia    SendMedia  `gorm:"column:send_via"`   // 发送方式
	SendTo     string     `gorm:"column:send_to"`    // 接收人
	NationCode string     `gorm:"nation_code"`       // 国家代码
	ExpiredAt  *time.Time `gorm:"column:expired_at"` // 到期时间
	HasUsed    bool       `gorm:"column:has_used"`   // 是否使用过
	CreatedAt  time.Time  // 创建时间
	UpdatedAt  time.Time  // 更新时间
	DeletedAt  *time.Time // 删除时间
}

// TableName 返回 VcRecord 所在的表名
func (v VcRecord) TableName() string {
	return "verification_code"
}

// CreateVcRecord 创建验证码记录
func (db *DbClient) CreateVcRecord(ctx context.Context, record *VcRecord) error {
	return db.DB(ctx).Create(record).Error
}

// SearchVcRecordCountsIn24hours 搜索24小时内的验证码记录个数
func (db *DbClient) SearchVcRecordCountsIn24hours(ctx context.Context, sendTo string) (int, error) {
	var count int
	err := db.DB(ctx).Model(&VcRecord{}).Where("send_to = ? and created_at > DATE_SUB(CURDATE(), INTERVAL 8 HOUR)", sendTo).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SearchVcRecordCountsIn1Minute 搜索1分钟内的验证码记录
func (db *DbClient) SearchVcRecordCountsIn1Minute(ctx context.Context, sendTo string) (int, error) {
	var count int
	err := db.DB(ctx).Model(&VcRecord{}).Where("send_to = ? and created_at > DATE_SUB(NOW(), INTERVAL 1 minute) and has_used = ? ", sendTo, false).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SearchVcRecordEarliestTimeIn1Minute 搜索1分钟最早的验证码记录时间
func (db *DbClient) SearchVcRecordEarliestTimeIn1Minute(ctx context.Context, sendTo string) (*VcRecord, error) {
	var record VcRecord
	err := db.DB(ctx).Raw(`SELECT created_at FROM verification_code WHERE send_to = ?
        AND created_at > DATE_SUB(NOW(), INTERVAL 1 minute) 
		ORDER BY created_at LIMIT 1`, sendTo).Scan(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// FindVcRecord 查找验证码记录
func (db *DbClient) FindVcRecord(ctx context.Context, sn string, vc string, sendTo string, usage Usage) (*VcRecord, error) {
	var record VcRecord
	err := db.DB(ctx).Raw("SELECT record_id,expired_at,has_used FROM verification_code WHERE sn = ? AND code = ? AND send_to = ? AND has_used = ? AND `usage` = ? ORDER BY created_at LIMIT 1", sn, vc, sendTo, false, usage).Scan(&record).Error
	return &record, err
}

// HasSnExpired 判断Sn是否已经过期
func (db *DbClient) HasSnExpired(ctx context.Context, sn string, vc string) (bool, error) {
	var count int
	err := db.DB(ctx).Model(&VcRecord{}).Where("sn = ? and code = ? and expired > NOW() and has_used = ? ", sn, vc, false).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count != 0 {
		return false, nil
	}
	return true, nil
}

// ModifyVcRecordStatus 修改验证码的状态
func (db *DbClient) ModifyVcRecordStatus(ctx context.Context, recordID int32) error {
	return db.DB(ctx).Model(&VcRecord{}).Where("record_id = ?", recordID).Updates(map[string]interface{}{
		"has_used": true,
	}).Error
}

// VerifyMVC 验证MVC
func (db *DbClient) VerifyMVC(ctx context.Context, sn string, vc string, sendTo string, nationCode string) (bool, error) {
	var count int
	err := db.DB(ctx).Raw(`SELECT count(record_id) FROM verification_code WHERE sn = ?
        AND code = ? AND send_to = ? AND expired_at > NOW() AND has_used = ? AND nation_code = ?
		ORDER BY created_at LIMIT 1`, sn, vc, sendTo, false, nationCode).Count(&count).Error
	if err != nil {
		return false, nil
	}
	return count >= 1, nil
}

// SearchVcRecord 查找验证码记录
func (db *DbClient) SearchVcRecord(ctx context.Context, sn string, vc string, sendTo string, nationCode string) (*VcRecord, error) {
	var record VcRecord
	err := db.DB(ctx).Raw(`SELECT record_id,expired_at,has_used FROM verification_code WHERE sn = ?
        AND code = ? AND send_to = ? AND expired_at > NOW() AND has_used = ? AND nation_code = ?
        ORDER BY created_at LIMIT 1`, sn, vc, sendTo, false, nationCode).Scan(&record).Error

	return &record, err
}

// FindLatestVcRecord 查找最新验证码记录
func (db *DbClient) FindLatestVcRecord(ctx context.Context, sendTo string, usage Usage) (*VcRecord, error) {
	var record VcRecord
	err := db.DB(ctx).Raw("SELECT sn, code, record_id, expired_at FROM verification_code WHERE `usage` = ? AND send_to = ? AND has_used = ? AND created_at > DATE_SUB(CURDATE(), INTERVAL 8 HOUR) ORDER BY created_at DESC LIMIT 1", usage, sendTo, false).Scan(&record).Error
	return &record, err
}

// SearchSpecificVcRecordCountsIn24hours 搜索24小时内的指定模板的验证码记录个数
func (db *DbClient) SearchSpecificVcRecordCountsIn24hours(ctx context.Context, sendTo string, usage Usage) (int, error) {
	var count int
	err := db.DB(ctx).Model(&VcRecord{}).Where("send_to = ? and `usage` = ? and created_at > DATE_SUB(CURDATE(), INTERVAL 8 HOUR)", sendTo, usage).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SearchSpecificVcRecordEarliestTimeIn24hours 搜索24小时指定模板最早的验证码记录
func (db *DbClient) SearchSpecificVcRecordEarliestTimeIn24hours(ctx context.Context, sendTo string, usage Usage) (*VcRecord, error) {
	var record VcRecord
	err := db.DB(ctx).Raw("SELECT expired_at FROM verification_code WHERE send_to = ? AND `usage` = ? AND created_at > DATE_SUB(CURDATE(), INTERVAL 8 HOUR) ORDER BY created_at LIMIT 1", sendTo, usage).Scan(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// SearchLatestPhoneVerificationCode 搜索最新的电话验证码
func (db *DbClient) SearchLatestPhoneVerificationCode(ctx context.Context, sendTo string, nationCode string) (string, error) {
	var record VcRecord
	err := db.DB(ctx).Raw(`SELECT code FROM verification_code WHERE send_to = ? AND nation_code = ? 
    ORDER BY created_at  DESC LIMIT 1`, sendTo, nationCode).Scan(&record).Error
	if err != nil {
		return "", err
	}
	return record.Code, nil
}

// SearchLatestEmailVerificationCode 搜索最新的邮件验证码
func (db *DbClient) SearchLatestEmailVerificationCode(ctx context.Context, sendTo string) (string, error) {
	var record VcRecord
	err := db.DB(ctx).Raw(`SELECT code FROM verification_code WHERE send_to = ? ORDER BY created_at DESC LIMIT 1`, sendTo).Scan(&record).Error
	if err != nil {
		return "", err
	}
	return record.Code, nil
}

// VerifyMVCBySecureEmail 根据安全邮箱寻找验证码信息
func (db *DbClient) VerifyMVCBySecureEmail(ctx context.Context, sn string, vc string, email string) (*VcRecord, error) {
	var record VcRecord
	err := db.DB(ctx).Raw(`SELECT expired_at,has_used FROM verification_code 
    WHERE sn = ?  AND code = ? AND send_to = ? 
     ORDER BY created_at DESC LIMIT 1`, sn, vc, email).Scan(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil

}

// ModifyVcRecordStatusByEmail 根据安全邮箱修改验证码的状态
func (db *DbClient) ModifyVcRecordStatusByEmail(ctx context.Context, email, verificationCode, serialNumber string) error {
	return db.DB(ctx).Model(&VcRecord{}).Where("send_to = ? and code = ? and sn = ?", email, verificationCode, serialNumber).Updates(map[string]interface{}{
		"has_used": true,
	}).Error
}

// SetVcAsUsed 设置vc为使用过的
func (db *DbClient) SetVcAsUsed(ctx context.Context, sn string, vc string, sendTo string, nationCode string) error {
	return db.DB(ctx).Model(&VcRecord{}).Where("sn = ? and code = ? and send_to = ? and nation_code = ?", sn, vc, sendTo, nationCode).Updates(map[string]interface{}{
		"has_used": true,
	}).Error
}

// SearchVcRecordFrom1MinuteTo2Mintue 搜索1-2分钟的短信记录
func (db *DbClient) SearchVcRecordFrom1MinuteTo2Mintue(ctx context.Context, sendTo string, usage Usage) (*VcRecord, error) {
	var record VcRecord
	err := db.DB(ctx).Raw("SELECT code FROM verification_code WHERE send_to = ? AND created_at > DATE_SUB(NOW(), INTERVAL 2 minute) and created_at < DATE_SUB(NOW(), INTERVAL 1 minute) and `usage` = ? ORDER BY created_at LIMIT 1", sendTo, usage).Scan(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}
