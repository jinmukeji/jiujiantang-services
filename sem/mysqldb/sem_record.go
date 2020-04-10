package mysqldb

import (
	"time"
)

// SemSendStatus 邮件发送状态
type SemSendStatus int32

const (
	// Pending 待定
	Pending SemSendStatus = 0
	// Sending 发送中
	Sending SemSendStatus = 1
	// SendSucceed 发送成功
	SendSucceed SemSendStatus = 2
	// SendFailed 发送失败
	SendFailed SemSendStatus = 3
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
	// UndefinedLanguage 没有定义的语言
	UndefinedLanguage Language = "undefined_language"
)

// SemRecord 邮件记录
type SemRecord struct {
	SemID          int32         `gorm:"primary_key"`            // 邮件记录ID
	ToAddress      string        `gorm:"column:to_address"`      // 收信邮箱地址
	SemStatus      SemSendStatus `gorm:"column:sem_status"`      // 是否发送成功
	TemplateAction string        `gorm:"column:template_action"` // 模版行为
	PlatformType   string        `gorm:"column:platform_type"`   // 平台
	TemplateParam  string        `gorm:"column:template_param"`  // 模版的填充值
	Language       Language      `gorm:"column:language"`        // 语言
	SemErrorLog    string        `gorm:"column:sem_error_log"`   // 错误日志
	CreatedAt      time.Time     // 创建时间
	UpdatedAt      time.Time     // 更新时间
	DeletedAt      *time.Time    // 删除时间
}

// TemplateParamAndExpiredAt 模版的填充值和其过期时间
type TemplateParamAndExpiredAt struct {
	TemplateParam string    `gorm:"template_param"`    // 模版的填充值
	ExpiredAt     time.Time `gorm:"column:expired_at"` // 过期时间
}

// TableName 返回 SemRecord对应的数据库数据表名
func (record SemRecord) TableName() string {
	return "sem_record"
}

// CreateSemRecord 创建邮件记录
func (db *DbClient) CreateSemRecord(record *SemRecord) error {
	return db.Create(record).Error
}

// SearchSemRecordCountsIn1Minute 搜索1分钟内的短信记录数目
func (db *DbClient) SearchSemRecordCountsIn1Minute(email string) (int, error) {
	var count int
	err := db.Model(&SemRecord{}).Where("SemRecord = ?and created_at > DATE_SUB(NOW(), INTERVAL 1 minute)", email).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
