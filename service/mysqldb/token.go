package mysqldb

import (
	"context"
	"time"
)

// Token 是账户登录会话信息数据
type Token struct {
	UserID    int32     `gorm:"user_id"`     // 用户 id
	Token     string    `gorm:"primary_key"` // token 是登录凭证
	ExpiredAt time.Time // token 过期的时间
	CreatedAt time.Time // 创建时间
	UpdatedAt time.Time // 更新时间
}

// TableName 返回 token 所在的表名
func (t Token) TableName() string {
	return "user_access_token"
}

// CreateToken 保存token 并立即返回创建的token
func (db *DbClient) CreateToken(ctx context.Context, token string, userID int32, availableDuration time.Duration) (*Token, error) {
	now := time.Now()
	tk := Token{
		Token:     token,
		UserID:    userID,
		CreatedAt: now,
		ExpiredAt: now.Add(availableDuration),
		UpdatedAt: now,
	}
	if err := db.GetDB(ctx).Create(&tk).Error; err != nil {
		return nil, err
	}
	return &tk, nil
}

// FindUserIDByToken 根据 token 返回 UserID，如果 token 失效返回 error
func (db *DbClient) FindUserIDByToken(ctx context.Context, token string) (int32, error) {
	var t Token
	now := time.Now()

	if err := db.GetDB(ctx).Where("expired_at > ? ", now).First(&t, "token = ?", token).Error; err != nil {
		return 0, err
	}

	return t.UserID, nil
}

// DeleteToken 删除数据库内指定的 token
func (db *DbClient) DeleteToken(ctx context.Context, token string) error {
	return db.GetDB(ctx).Delete(Token{}, "token = ?", token).Error
}
