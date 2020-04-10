package mysqldb

import (
	"context"
	"time"

	"github.com/jinmukeji/jiujiantang-services/analysis/aws"
)

const (
	// UnknownStatus 表示未知的状态
	UnknownStatus int32 = -1
	// PulseTestRawDataS3KeyPrefix S3原始数据地址后缀
	PulseTestRawDataS3KeyPrefix string = "spec-v2"
)

// TODO: record表以后迁移到分析服务

// AnalysisStatus 分析状态
type AnalysisStatus int32

const (
	// AnalysisStatusPending 待定
	AnalysisStatusPending AnalysisStatus = 0
	// AnalysisStatusInProgress 进行中
	AnalysisStatusInProgress AnalysisStatus = 1
	// AnalysisStatusCompeleted 完成
	AnalysisStatusCompeleted AnalysisStatus = 2
	// AnalysisStatusError 错误
	AnalysisStatusError AnalysisStatus = 3
)

// 手指
type Finger int32

const (
	// FingerLeft1 左小拇指
	FingerLeft1 Finger = 1
	// FingerLeft2 左无名指x
	FingerLeft2 Finger = 2
	// FingerLeft3 左中指
	FingerLeft3 Finger = 3
	// FingerLeft4 左食指
	FingerLeft4 Finger = 4
	// FingerLeft5 左大拇指
	FingerLeft5 Finger = 5
	// FingerRight5 右大拇指
	FingerRight5 Finger = 6
	// FingerRight4 右食指
	FingerRight4 Finger = 7
	// FingerRight3 右中指
	FingerRight3 Finger = 8
	// FingerRight2 右无名指
	FingerRight2 Finger = 9
	// FingerRight1 右小拇指
	FingerRight1 Finger = 10
	// FingerInvalid  非法的手指
	FingerInvalid Finger = -1
)

// 测量姿势
type MeasurementPosture int32

const (
	// MeasurementPostureSetting 坐姿
	MeasurementPostureSetting MeasurementPosture = 0
	// MeasurementPostureStanging 站姿
	MeasurementPostureStanging MeasurementPosture = 1
	// MeasurementPostureLying 躺姿
	MeasurementPostureLying MeasurementPosture = 2
	// MeasurementPostureInvalid 错误的姿势
	MeasurementPostureInvalid MeasurementPosture = -1
)

// Record 是测量结果记录
type Record struct {
	RecordID int    `gorm:"primary_key"`      // 测量结果记录ID
	ClientID string `gorm:"column:client_id"` // 客户端 ID
	UserID   int    `gorm:"user_id"`          // 用户档案ID
	DeviceID int    `gorm:"column:device_id"` // 设备ID
	Finger   Finger `gorm:"column:finger"`    // 左右手

	AppHeartRate   float64 `gorm:"column:app_heart_rate"`    // app 计算的心率
	IsSportOrDrunk int     `gorm:"column:is_sport_or_drunk"` // 运动或饮酒
	Cold           int     `gorm:"column:cold"`              // 感冒或病毒感染期
	MenstrualCycle int     `gorm:"column:menstrual_cycle"`   // 生理周期
	OvipositPeriod int     `gorm:"column:oviposit_period"`   // 排卵期
	Lactation      int     `gorm:"column:lactation"`         // 哺乳期
	Pregnancy      int     `gorm:"column:pregnancy"`         // 怀孕
	StatusA        int     `gorm:"column:cm_app_status_a"`   // 口苦口黏，皮肤瘙痒，大便不成形，头重身痛
	StatusB        int     `gorm:"column:cm_app_status_b"`   // 急躁易怒，头晕胀痛
	StatusC        int     `gorm:"column:cm_app_status_c"`   // 口苦听力下降女性带下异味小便黄短
	StatusD        int     `gorm:"column:cm_app_status_d"`   // 口中异味反酸便秘喉咙干痒牙龈出血
	StatusE        int     `gorm:"column:cm_app_status_e"`   // 胃部冷痛，得温缓解
	StatusF        int     `gorm:"column:cm_app_status_f"`   // 失眠多梦健忘眩晕

	C0   float64 `gorm:"column:c0"`   // 心包经测量指标
	C1   float64 `gorm:"column:c1"`   // 肝经测量指标
	C2   float64 `gorm:"column:c2"`   // 肾经测量指标
	C3   float64 `gorm:"column:c3"`   // 脾经测量指标
	C4   float64 `gorm:"column:c4"`   // 肺经测量指标
	C5   float64 `gorm:"column:c5"`   // 胃经测量指标
	C6   float64 `gorm:"column:c6"`   // 胆经测量指标
	C7   float64 `gorm:"column:c7"`   // 膀胱经测量指标
	C0CV float64 `gorm:"column:c0cv"` // 心包经 变异
	C1CV float64 `gorm:"column:c1cv"` // 肝经 变异
	C2CV float64 `gorm:"column:c2cv"` // 肾经 变异
	C3CV float64 `gorm:"column:c3cv"` // 脾经 变异
	C4CV float64 `gorm:"column:c4cv"` // 肺经 变异
	C5CV float64 `gorm:"column:c5cv"` // 胃经 变异
	C6CV float64 `gorm:"column:c6cv"` // 胆经 变异
	C7CV float64 `gorm:"column:c7cv"` // 膀胱经 变异
	G0   int32   `gorm:"column:g0"`   // 风_心包经
	G1   int32   `gorm:"column:g1"`   // 风_肝经
	G2   int32   `gorm:"column:g2"`   // 风_肾经
	G3   int32   `gorm:"column:g3"`   // 风_脾经
	G4   int32   `gorm:"column:g4"`   // 风_肺经
	G5   int32   `gorm:"column:g5"`   // 风_胃经
	G6   int32   `gorm:"column:g6"`   // 风_胆经
	G7   int32   `gorm:"column:g7"`   // 风_膀胱经

	HeartRate                 float64 `gorm:"column:heart_rate"`                   // 算法服务测的心率
	AlgorithmHighestHeartRate int32   `gorm:"column:algorithm_highest_heart_rate"` // 算法服务计算得到的最高心率
	AlgorithmLowestHeartRate  int32   `gorm:"column:algorithm_lowest_heart_rate"`  // 算法服务计算得到的最低心率
	HeartRateCV               float32 `gorm:"column:heart_rate_cv"`                // 心率变异
	SNR                       float32 `gorm:"column:snr"`                          // 信噪比
	DcDrift                   float32 `gorm:"column:dc_drift"`                     // 直流漂移
	Elapsed                   int     `gorm:"column:elapsed"`                      // 计算所耗时间
	Remark                    string  `gorm:"column:remark"`                       // 备注

	RecordType                      int                `gorm:"column:record_type"`                          // 记录类型
	Answers                         string             `gorm:"column:answers"`                              // 智能分析的答案
	IsValid                         int                `gorm:"column:is_valid"`                             // 是否有效
	AppHighestHeartRate             int32              `gorm:"column:app_highest_heart_rate"`               // app最高心率
	AppLowestHeartRate              int32              `gorm:"column:app_lowest_heart_rate"`                // app最低心率
	HasPaid                         bool               `gorm:"column:has_paid"`                             // 是否完成支付
	ShowFullReport                  bool               `gorm:"column:show_full_report"`                     // 是否显示完成测量报告
	HasSentWxViewReportNotification bool               `gorm:"column:has_sent_wx_view_report_notification"` // 是否已经发送过微信查看报告通知
	RecordToken                     string             `grom:"column:record_token"`                         // 分享报告的token
	HasAEError                      int32              `gorm:"column:has_ae_error"`                         // ae得出的结果是否异常
	MeasurementPosture              MeasurementPosture `gorm:"column:measurement_posture"`                  // 测量姿态
	TransactionNumber               string             `gorm:"column:transaction_number"`                   // 流水号
	S3Key                           string             `gorm:"column:s3_key"`                               // S3的key
	StressState                     string             `gorm:"column:stress_state"`                         // 应激态json数组 map[string]bool
	HasStressState                  bool               `gorm:"column:has_stress_state"`                     // 是否是应激态
	AnalyzeStatus                   AnalysisStatus     `gorm:"column:analyze_status"`                       // 分析状态 0 pending,1 in_progress,2 completed,3 error
	AnalyzeBody                     string             `gorm:"column:analyze_body"`                         // 新分析接口的body

	CreatedAt time.Time  // 创建时间
	UpdatedAt time.Time  // 更新时间
	DeletedAt *time.Time // 删除时间
}

// TableName 返回 Record 对应的数据库数据表名
func (r Record) TableName() string {
	return "record"
}

// FindRecordByID 查找指定 ID 的一条 Record 数据记录
func (db *DbClient) FindRecordByID(ctx context.Context, recordID int) (*Record, error) {
	var r Record
	if err := db.First(&r, "record_id = ? ", recordID).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

// SetRecordSentWxViewReportNotification 设置是否已经发送过微信查看报告通知的状态标记 0: 未发送，1: 已发送
func (db *DbClient) SetRecordSentWxViewReportNotification(ctx context.Context, recordID int, sent bool) error {
	return db.Model(&Record{}).Where("record_id = ?", recordID).Update(map[string]interface{}{
		"has_sent_wx_view_report_notification": sent,
		"updated_at":                           time.Now().UTC(),
	}).Error
}

// FindValidRecordByID 查找指定 ID 的一条有效的 Record 数据记录
func (db *DbClient) FindValidRecordByID(ctx context.Context, recordID int) (*Record, error) {
	var r Record
	if err := db.First(&r, "record_id = ? and is_valid = ?", recordID, DbValidValue).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

// UpdateRemarkByRecordID 更新备注
func (db *DbClient) UpdateRemarkByRecordID(ctx context.Context, recordID int, remark string) error {
	return db.Model(&Record{}).Where("record_id = ?", recordID).Update(map[string]interface{}{
		"remark":     remark,
		"updated_at": time.Now().UTC(),
	}).Error
}

// FindValidRecordsByDateRange 返回给定时间范围内的有效 record
func (db *DbClient) FindValidRecordsByDateRange(ctx context.Context, subjectID int, start time.Time, end time.Time) ([]Record, error) {
	var records []Record
	if err := db.Model(&Record{}).Order("create_date_utc desc").Where("(subject_id = ?) AND ( create_date_utc between ? AND ? ) AND( is_valid = ?)", subjectID, start, end, DbValidValue).Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

// CreateRecord 增加测量记录
func (db *DbClient) CreateRecord(ctx context.Context, record *Record) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	errCreate := tx.Create(record).Error
	if errCreate != nil {
		tx.Rollback()
		return errCreate
	}
	key := aws.GenerateS3Key(PulseTestRawDataS3KeyPrefix, record.RecordID)
	errUpdateS3Key := tx.Model(&Record{}).Where("record_id = ?", record.RecordID).Update(map[string]interface{}{
		"s3_key":     key,
		"updated_at": time.Now().UTC(),
	}).Error

	if errUpdateS3Key != nil {
		tx.Rollback()
		return errUpdateS3Key
	}
	return tx.Commit().Error
}

// CheckUserHasRecord 验证用户和测量记录是否有关联
func (db *DbClient) CheckUserHasRecord(ctx context.Context, userID int32, recordID int32) (bool, error) {
	var count int
	err := db.Raw("select count(*) from record where user_id = ? and record_id = ?", userID, recordID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindValidPaginatedRecordsByDateRange 返回给定时间和给定数量范围内的有效 record
func (db *DbClient) FindValidPaginatedRecordsByDateRange(ctx context.Context, userID, offset, size int, start, end time.Time) ([]Record, error) {
	var records []Record
	if size == -1 {
		size = maxQuerySize
	}
	if err := db.Model(&Record{}).Order("created_at desc").Where("(user_id = ?) AND ( created_at between ? AND ? ) AND( is_valid = ?) AND ( has_ae_error = 0) AND analyze_status = 2", userID, start, end, DbValidValue).Offset(offset).Limit(size).Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

// FindValidPaginatedRecordsByUserID 通过UserID拿到records
func (db *DbClient) FindValidPaginatedRecordsByUserID(ctx context.Context, userID int32) ([]Record, error) {
	var records []Record
	if err := db.Model(&Record{}).Order("created_at desc").Where("user_id = ? AND is_valid = 1 AND has_ae_error = 0 AND analyze_status = 2", userID).Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

// GetUserIDByRecordID 通过RecordID得到UserID
func (db *DbClient) GetUserIDByRecordID(ctx context.Context, recordID int32) (int32, error) {
	var record Record
	db.Model(&Record{}).Where("record_id = ?", recordID).Scan(&record)
	return int32(record.UserID), nil
}

// UpdateRecordStatus 更新状态
func (db *DbClient) UpdateRecordStatus(ctx context.Context, record *Record) error {
	return db.Model(&Record{}).Where("record_id = ?", record.RecordID).Update(map[string]interface{}{
		"is_sport_or_drunk": record.IsSportOrDrunk,
		"cold":              record.Cold,
		"menstrual_cycle":   record.MenstrualCycle,
		"oviposit_period":   record.OvipositPeriod,
		"lactation":         record.Lactation,
		"pregnancy":         record.Pregnancy,
		"cm_app_status_a":   record.StatusA,
		"cm_app_status_b":   record.StatusB,
		"cm_app_status_c":   record.StatusC,
		"cm_app_status_d":   record.StatusD,
		"cm_app_status_e":   record.StatusE,
		"cm_app_status_f":   record.StatusF,
		"is_valid":          1,
	}).Error
}

// ExistRecordByRecordID 检查Record是否存在
func (db *DbClient) ExistRecordByRecordID(ctx context.Context, recordID int32) (bool, error) {
	var count int
	db.Model(&Record{}).Where("record_id = ? AND is_valid = 1 ", recordID).Count(&count)
	return count != 0, nil
}

// UpdateRecordAnswers 更新Record中的Answers
func (db *DbClient) UpdateRecordAnswers(ctx context.Context, recordID int32, answers string) error {
	return db.Model(&Record{}).Where("record_id = ?", recordID).Update(map[string]interface{}{
		"answers":        answers,
		"analyze_status": AnalysisStatusCompeleted,
		"updated_at":     time.Now().UTC(),
	}).Error
}

// UpdateRecordHasPaid 更新Record中的hasPaid
func (db *DbClient) UpdateRecordHasPaid(ctx context.Context, recordID int32) error {
	return db.Model(&Record{}).Where("record_id = ?", recordID).Update(map[string]interface{}{
		"has_paid":   true,
		"updated_at": time.Now().UTC(),
	}).Error
}

// UpdateRecordToken 更新分享报告的token
func (db *DbClient) UpdateRecordToken(ctx context.Context, recordID int32, token string) error {
	return db.Model(&Record{}).Where("record_id = ?", recordID).Update(map[string]interface{}{
		"record_token": token,
		"updated_at":   time.Now().UTC(),
	}).Error
}

// FindRecordByToken 查找指定 token 的一条 Record 数据记录
func (db *DbClient) FindRecordByToken(ctx context.Context, token string) (*Record, error) {
	var r Record
	if err := db.First(&r, "record_token = ? ", token).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

// UpdateRecordHasAEError 更新记录有效性
func (db *DbClient) UpdateRecordHasAEError(ctx context.Context, recordID int32, hasAEError int32) error {
	return db.Model(&Record{}).Where("record_id = ?", recordID).Update(map[string]interface{}{
		"has_ae_error": hasAEError,
		"updated_at":   time.Now().UTC(),
	}).Error
}

// UpdateRecordTransactionNumber 更新记录的流水号
func (db *DbClient) UpdateRecordTransactionNumber(ctx context.Context, recordID int32, transactionNumber string) error {
	return db.Model(&Record{}).Where("record_id = ?", recordID).Update(map[string]interface{}{
		"transaction_number": transactionNumber,
		"updated_at":         time.Now().UTC(),
	}).Error
}

// DeleteRecord 删除记录
func (db *DbClient) DeleteRecord(ctx context.Context, recordID int32) error {
	now := time.Now()
	return db.Model(&Record{}).Where("record_id = ?", recordID).Update(map[string]interface{}{
		"deleted_at": now.UTC(),
		"updated_at": now.UTC(),
	}).Error
}
