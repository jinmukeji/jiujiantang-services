package mysqldb

import (
	"context"
	"time"
)

// JinmuLAccessToken 是一体机账户登录会话信息数据
type JinmuLAccessToken struct {
	Account     string    `gorm:"account"`             // 账户
	Token       string    `gorm:"primary_key"`         // token 是登录凭证
	MachineUUID string    `gorm:"column:machine_uuid"` // machine_uuid
	ExpiredAt   time.Time // token 过期的时间
	CreatedAt   time.Time // 创建时间
	UpdatedAt   time.Time // 更新时间
}

// TableName 返回 token 所在的表名
func (t JinmuLAccessToken) TableName() string {
	return "jinmu_l_access_token"
}

// CreateAccessToken 保存token 并立即返回创建的token
func (db *DbClient) CreateAccessToken(ctx context.Context, token string, account string, machineUUID string, availableDuration time.Duration) (*JinmuLAccessToken, error) {
	now := time.Now()
	tk := JinmuLAccessToken{
		Token:       token,
		Account:     account,
		MachineUUID: machineUUID,
		CreatedAt:   now,
		ExpiredAt:   now.Add(availableDuration),
		UpdatedAt:   now,
	}
	if err := db.GetDB(ctx).Create(&tk).Error; err != nil {
		return nil, err
	}
	return &tk, nil
}

// FindJinmuLAccountByToken 通过token找 JinmuL account
func (db *DbClient) FindJinmuLAccountByToken(ctx context.Context, token string) (string, error) {
	var t JinmuLAccessToken
	now := time.Now()

	if err := db.GetDB(ctx).Where("expired_at > ? ", now).First(&t, "token = ?", token).Error; err != nil {
		return "", err
	}
	return t.Account, nil
}

// FindMachineUUIDByToken 通过token找 Machine UUID
func (db *DbClient) FindMachineUUIDByToken(ctx context.Context, token string) (string, error) {
	var t JinmuLAccessToken
	now := time.Now()

	if err := db.GetDB(ctx).Where("expired_at > ? ", now).First(&t, "token = ?", token).Error; err != nil {
		return "", err
	}
	return t.MachineUUID, nil
}

// DeleteJinmuLAccessToken 删除AccessToken
func (db *DbClient) DeleteJinmuLAccessToken(ctx context.Context, token string) error {
	return db.GetDB(ctx).Delete(JinmuLAccessToken{}, "token = ?", token).Error
}
