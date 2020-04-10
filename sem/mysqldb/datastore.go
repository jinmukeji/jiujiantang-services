package mysqldb

// Datastore 定义数据访问接口
type Datastore interface {
	// CreateSemRecord 创建邮件记录
	CreateSemRecord(record *SemRecord) error
	// SearchSemRecordCountsIn1Minute 搜索1分钟内的邮件记录数目
	SearchSemRecordCountsIn1Minute(email string) (int, error)
}
