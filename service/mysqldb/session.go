package mysqldb

import (
	"context"
	"time"
)

// TableName 返回 Organization 对应的数据库数据表名
func (s Session) TableName() string {
	return "wechat_h5_session"
}

// Session 订阅
type Session struct {
	SessionID  string     `gorm:"primary_key"` // Session ID
	State      string     `grom:"state"`       // 微信 OAuth 验证的 state'
	OpenID     string     `grom:"open_id"`     // 微信OpenID
	UnionID    string     `grom:"union_id"`    // 微信UnionID
	UserID     int64      `grom:"user_id"`     // 用户ID
	Authorized bool       `grom:"authorized"`  // 是否已经验证通过
	ExpiredAt  time.Time  // 到期时间
	CreatedAt  time.Time  // 创建时间
	UpdatedAt  time.Time  // 更新时间
	DeletedAt  *time.Time // 删除时间
}

// CreateSession 创建Session
func (db *DbClient) CreateSession(ctx context.Context, s *Session) (*Session, error) {
	if err := db.GetDB(ctx).Create(&s).Error; err != nil {
		return nil, err
	}
	return s, nil
}

// FindSession 寻找Session
func (db *DbClient) FindSession(ctx context.Context, sessionID string) (*Session, error) {
	var session Session
	err := db.GetDB(ctx).Model(&Session{}).Where("session_id = ?", sessionID).Scan(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// UpdateSession 更新Session
func (db *DbClient) UpdateSession(ctx context.Context, sessionID string, session *Session) error {
	return db.GetDB(ctx).Model(&Session{}).Where("session_id = ?", sessionID).Updates(map[string]interface{}{
		"state":      session.State,
		"open_id":    session.OpenID,
		"union_id":   session.UnionID,
		"user_id":    session.UserID,
		"authorized": session.Authorized,
		"updated_at": session.UpdatedAt,
	}).Error
}
