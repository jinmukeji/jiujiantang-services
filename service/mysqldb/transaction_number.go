package mysqldb

import (
	"context"
	"time"
)

// TransactionNumber 流水号
type TransactionNumber struct {
	TransactionDate   time.Time `gorm:"primary_key"` // 日期
	TransactionNumber int       // 流水号
	CreatedAt         time.Time // 创建时间
	UpdatedAt         time.Time // 更新时间
}

// TableName 返回 TransactionNumber 所在的表名
func (t TransactionNumber) TableName() string {
	return "transaction_number"
}

// FindTransactionNumberByCurrentDate 查询流水号
func (d *DbClient) FindTransactionNumberByCurrentDate(ctx context.Context) (*TransactionNumber, error) {
	var t TransactionNumber
	now := time.Now().UTC()
	format := "2006-01-02"
	d.Model(&TransactionNumber{}).Where("transaction_date = ?", now.Format(format)).Scan(&t)

	d.Model(&TransactionNumber{}).Where("transaction_date = ?", now.Format(format)).Update(map[string]interface{}{
		"transaction_number": t.TransactionNumber + 1,
		"create_at":          now,
		"updated_at":         now,
	})
	return &t, nil
}

// IsExistTransactionNumberByCurrentDate 流水号是否存在
func (d *DbClient) IsExistTransactionNumberByCurrentDate(ctx context.Context) (bool, error) {
	var count int
	now := time.Now().UTC()
	format := "2006-01-02"
	d.Model(&TransactionNumber{}).Where("transaction_date = ?", now.Format(format)).Count(&count)
	return count != 0, nil
}

// CreateTransactionNumber 创建流水号
func (d *DbClient) CreateTransactionNumber(ctx context.Context) (*TransactionNumber, error) {
	now := time.Now().UTC()
	tn := TransactionNumber{
		TransactionDate:   now,
		TransactionNumber: 1,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	if err := d.Create(&tn).Error; err != nil {
		return nil, err
	}
	return &tn, nil
}
