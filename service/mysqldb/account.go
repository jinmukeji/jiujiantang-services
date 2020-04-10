package mysqldb

import (
	"context"
	"time"
)

// JinmuLAccount 账户信息
type JinmuLAccount struct {
	Account        string `gorm:"primary_key"`            // 账户
	Password       string `gorm:"column:password"`        // 用户登录密码
	OrganizationID int32  `gorm:"column:organization_id"` // 与设备关联的组织ID
	Remark         string `gorm:"column:remark"`          // 账户备注
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

// TableName 返回表名
func (a JinmuLAccount) TableName() string {
	return "jinmu_l_account"
}

// FindJinmuLAccount 查找account信息
func (db *DbClient) FindJinmuLAccount(ctx context.Context, account string) (*JinmuLAccount, error) {
	var jinmuLAccount JinmuLAccount
	if err := db.First(&jinmuLAccount, "( account = ? ) ", account).Error; err != nil {
		return nil, err
	}
	return &jinmuLAccount, nil
}
