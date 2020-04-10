package mysqldb

import "context"

// Datastore 定义数据访问接口
type Datastore interface {
	// FindUserIDByToken 根据 token 返回 userID，如果token失效返回 error
	// TODO 以后重写，删除这个方法
	FindUserIDByToken(token string) (int32, error)
	// FindAnalysisParams 查找分析的参数,c0-c7,gender,age,weight,height,heart_rate
	FindAnalysisParams(recordID int32) (*Record, error)
	// UpdateAnalysisRecord 更新分析报告
	UpdateAnalysisRecord(record *Record) error
	// FindRecordByRecordID 通过 recordID 找到 record
	FindRecordByRecordID(recordID int32) (*Record, error)
	// FindAnalysisBodyByToken 查找AnalysisBody
	FindAnalysisBodyByToken(token string) (*Record, error)
	// UpdateAnalysisStatusError 更新分析状态错误
	UpdateAnalysisStatusError(recordID int32) error
	// UpdateAnalysisStatusInProgress 更新分析进行中
	UpdateAnalysisStatusInProgress(recordID int32) error
	// UpdateRecordHasAEError 更新记录有效性
	UpdateRecordHasAEError(recordID int32, hasAEError AEStatus) error
	// UpdateRecordTransactionNumber 更新记录的流水号
	UpdateRecordTransactionNumber(ctx context.Context, recordID int32, transactionNumber string) error
}
