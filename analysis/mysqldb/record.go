package mysqldb

import (
	"context"
	"time"
)

// AnalysisStatus 分析状态
type AnalysisStatus int32

const (
	// AnalysisStatusPending 待决
	AnalysisStatusPending AnalysisStatus = 0
	// AnalysisStatusInProgress 进行中
	AnalysisStatusInProgress AnalysisStatus = 1
	// AnalysisStatusCompeleted 完成
	AnalysisStatusCompeleted AnalysisStatus = 2
	// AnalysisStatusError 错误
	AnalysisStatusError AnalysisStatus = 3
)

// AEStatus AE的有效性
type AEStatus int32

const (
	// HasAeNoError 没有ae错误
	HasAeNoError AEStatus = 0
	// HasAeError 有ae错误
	HasAeError AEStatus = 1
)

// Gender 性别
type Gender string

const (
	// GenderMale 男性
	GenderMale Gender = "M"
	// GenderFemale 女性
	GenderFemale Gender = "F"
	// GenderInvalid 非法的性别
	GenderInvalid Gender = ""
)

// 手指
type Finger int32

const (
	// FingerLeft1 左小拇指
	FingerLeft1 Finger = 1
	// FingerLeft2 左无名指
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

// Record 记录
type Record struct {
	RecordID                  int32              `gorm:"primary_key"`    // 测量结果记录ID
	ClientID                  string             `gorm:"client_id"`      // 客户端ID
	UserID                    int32              `gorm:"column:user_id"` // 用户档案ID
	C0                        float64            `gorm:"column:c0"`      // 心包经测量指标
	C1                        float64            `gorm:"column:c1"`      // 肝经测量指标
	C2                        float64            `gorm:"column:c2"`      // 肾经测量指标
	C3                        float64            `gorm:"column:c3"`      // 脾经测量指标
	C4                        float64            `gorm:"column:c4"`      // 肺经测量指标
	C5                        float64            `gorm:"column:c5"`      // 胃经测量指标
	C6                        float64            `gorm:"column:c6"`      // 胆经测量指标
	C7                        float64            `gorm:"column:c7"`      // 膀胱经测量指标
	HeartRate                 float64            `gorm:"column:heart_rate"`
	AlgorithmHighestHeartRate int32              `gorm:"column:algorithm_highest_heart_rate"` // 算法服务计算得到的最高心率
	AlgorithmLowestHeartRate  int32              `gorm:"column:algorithm_lowest_heart_rate"`  // 算法服务计算得到的最低心率
	Finger                    Finger             `gorm:"column:finger"`                       // 左右手
	Remark                    string             `gorm:"column:remark"`                       // 备注
	HasAEError                int32              `gorm:"column:has_ae_error"`                 // ae得出的结果是否异常
	S3Key                     string             `gorm:"column:s3_key"`                       // S3的key
	CustomizedCode            string             `gorm:"-"`                                   // 用户自定义代码
	HasStressState            bool               `gorm:"column:has_stress_state"`             // 是否是应激态
	StressState               string             `gorm:"column:stress_state"`                 // 应激态json数组 map[string]bool
	AnalyzeBody               string             `gorm:"column:analyze_body"`                 // 新分析接口的body
	AnalyzeStatus             AnalysisStatus     `gorm:"column:analyze_status"`               // 分析状态 0 pending,1 in_progress,2 completed,3 error
	MeasurementPosture        MeasurementPosture `gorm:"column:measurement_posture"`          // 测量姿态
	TransactionNumber         string             `gorm:"column:transaction_number"`           // 流水号
	CreatedAt                 time.Time          // 创建时间
	UpdatedAt                 time.Time          // 更新时间
	DeletedAt                 *time.Time         // 删除时间
}

// TableName 表名
func (r Record) TableName() string {
	return "record"
}

// UpdateAnalysisRecord 更新分析记录
func (db *DbClient) UpdateAnalysisRecord(record *Record) error {
	return db.Model(&Record{}).Where("record_id = ?", record.RecordID).Update(map[string]interface{}{
		"has_stress_state": record.HasStressState,
		"stress_state":     record.StressState,
		"analyze_body":     record.AnalyzeBody,
		"analyze_status":   record.AnalyzeStatus,
		"updated_at":       time.Now().UTC(),
	}).Error
}

// FindAnalysisParams 找到分析的参数
func (db *DbClient) FindAnalysisParams(recordID int32) (*Record, error) {
	var record Record
	err := db.Raw(`SELECT 
    R.c0,
    R.c1,
    R.c2,
    R.c3,
    R.c4,
    R.c5,
    R.c6,
	R.c7,
	R.user_id,
	R.client_id,
	R.finger,
	R.transaction_number,
	UP.nickname,
	UP.nickname_initial,
    UP.gender,
    UP.birthday,
    UP.height,
    UP.weight,
	R.remark,
	R.heart_rate,
	R.s3_key,
	U.customized_code,
	TIMESTAMPDIFF(YEAR,UP.birthday,CURDATE()) age,
	R.created_at
FROM
    record AS R
        INNER JOIN
    user_profile AS UP ON UP.user_id = R.user_id
    INNER JOIN
    user AS U ON U.user_id = R.user_id
WHERE
	R.record_id = ? AND R.deleted_at IS NULL AND UP.deleted_at IS NULL`, recordID).Scan(&record).Error
	return &record, err
}

// FindRecordByRecordID 通过 recordID 找到 record
func (db *DbClient) FindRecordByRecordID(recordID int32) (*Record, error) {

	var record Record
	if err := db.First(&record, "( record_id = ? AND deleted_at IS NULL ) ", recordID).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// FindAnalysisBodyByToken 通过token找AnalysisBody
func (db *DbClient) FindAnalysisBodyByToken(token string) (*Record, error) {
	var record Record
	err := db.Raw(`SELECT
    R.record_id, 
    R.user_id,
    R.analyze_body,
    R.analyze_status,
	R.created_at
	FROM
	record AS R
	WHERE
	R.record_token = ? AND R.deleted_at IS NULL AND R.analyze_status = 2`, token).Scan(&record).Error
	return &record, err
}

// UpdateAnalysisStatusError  更新分析状态错误
func (db *DbClient) UpdateAnalysisStatusError(recordID int32) error {
	return db.Model(&Record{}).Where("record_id = ?", recordID).Update(map[string]interface{}{
		"analyze_status": AnalysisStatusError,
		"updated_at":     time.Now().UTC(),
	}).Error
}

// UpdateAnalysisStatusInProgress  更新分析进行中
func (db *DbClient) UpdateAnalysisStatusInProgress(recordID int32) error {
	return db.Model(&Record{}).Where("record_id = ?", recordID).Update(map[string]interface{}{
		"analyze_status": AnalysisStatusInProgress,
		"updated_at":     time.Now().UTC(),
	}).Error
}

// UpdateRecordHasAEError 更新记录有效性
func (db *DbClient) UpdateRecordHasAEError(recordID int32, hasAEError AEStatus) error {
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
