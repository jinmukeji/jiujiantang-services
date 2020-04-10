package mysqldb

// Datastore 定义数据访问接口
type Datastore interface {
	// CreateSmsRecord 创建短信记录
	CreateSmsRecord(record *SmsRecord) error
	// SearchSmsRecordCountsIn24hours 搜索24小时内的短信记录数目
	SearchSmsRecordCountsIn24hours(phone string, nationCode string) (int, error)
	// SearchSmsRecordCountsIn1Minute 搜索1分钟内的短信记录数目
	SearchSmsRecordCountsIn1Minute(phone string, nationCode string) (int, error)
}
