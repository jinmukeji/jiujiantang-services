package mysqldb

import (
	"context"
	"time"
)

// Client 是客户端信息
type Client struct {
	ClientID       string `gorm:"primary_key"` // client 表主键
	SecretKey      string // 客户端密钥
	Name           string // 客户端名称
	Zone           string // 客户端区域
	CustomizedCode string // 定制化代码
	Remark         string // 备注
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
	if err := db.First(&client, "( client_id = ? ) ", clientID).Error; err != nil {
		return nil, err
	}
	return &client, nil
}
