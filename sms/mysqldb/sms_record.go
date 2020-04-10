package mysqldb

import (
	"time"
)

// SmsSendStatus 短信发送状态
type SmsSendStatus int32

const (
	// Pending 待定
	Pending SmsSendStatus = 0
	// Sending 发送中
	Sending SmsSendStatus = 1
	// SendSucceed 发送成功
	SendSucceed SmsSendStatus = 2
	// SendFailed 发送失败
	SendFailed SmsSendStatus = 3
)

// Language 语言
type Language string

const (
	// SimpleChinese 简体中文
	SimpleChinese Language = "zh-Hans"
	// TraditionalChinese 繁体中文
	TraditionalChinese Language = "zh-Hant"
	// English 英文
	English Language = "en"
)

// SmsRecord 短信记录
type SmsRecord struct {
	SmsID          int32         `gorm:"primary_key"`     // 短信记录ID
	Phone          string        `gorm:"phone"`           // 手机号码
	SmsStatus      SmsSendStatus `gorm:"sms_status"`      // 是否发送成功
	TemplateAction string        `gorm:"template_action"` // 模版行为
	NationCode     string        `gorm:"nation_code"`     // 国家代码
	PlatformType   string        `gorm:"platform_type"`   // 平台
	TemplateParam  string        `gorm:"template_param"`  // 模版的填充值
	SmsErrorLog    string        `gorm:"sms_error_log"`   // 错误日志
	Language       Language      `gorm:"language"`        // 语言
	CreatedAt      time.Time     // 创建时间
	UpdatedAt      time.Time     // 更新时间
	DeletedAt      *time.Time    // 删除时间
}

// TableName 返回 SmsRecord对应的数据库数据表名
func (record SmsRecord) TableName() string {
	return "sms_record"
}

// CreateSmsRecord 创建短信记录
func (db *DbClient) CreateSmsRecord(record *SmsRecord) error {
	return db.Create(record).Error
}

// SearchSmsRecordCountsIn24hours 搜索24小时内的短信记录数目
func (db *DbClient) SearchSmsRecordCountsIn24hours(phone, nationCode string) (int, error) {
	var count int
	err := db.Model(&SmsRecord{}).Where("phone = ? and nation_code = ? and created_at > DATE_SUB(CURDATE(), INTERVAL 8 HOUR) and sms_status = 2", phone, nationCode).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SearchSmsRecordCountsIn1Minute 搜索1分钟内的短信记录数目
func (db *DbClient) SearchSmsRecordCountsIn1Minute(phone, nationCode string) (int, error) {
	var count int
	err := db.Model(&SmsRecord{}).Where("phone = ? and  nation_code = ? and created_at > DATE_SUB(NOW(), INTERVAL 1 minute)", phone, nationCode).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
