package mysqldb

import (
	"context"
	"time"
)

// Client 是客户端信息
type Client struct {
	ClientID       string `gorm:"primary_key"`            // client 表主键
	SecretKey      string `gorm:"column:secret_key"`      // 客户端密钥
	Name           string `gorm:"column:name"`            // 客户端名称
	Zone           string `gorm:"column:zone"`            // 客户端区域
	CustomizedCode string `gorm:"column:customized_code"` // 定制化代码
	Remark         string `gorm:"column:remark"`          // 备注
	Usage          string `gorm:"column:usage"`           // 用途
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

// TableName 返回表名
func (c Client) TableName() string {
	return "client"
}

// FindClientByClientID 查找一条 Client 数据记录
func (db *DbClient) FindClientByClientID(ctx context.Context, clientID string) (*Client, error) {
	var client Client
	if err := db.DB(ctx).First(&client, "( client_id = ? ) ", clientID).Error; err != nil {
		return nil, err
	}
	return &client, nil
}
