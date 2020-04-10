package handler

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/golang/protobuf/ptypes"

	"github.com/jinmukeji/jiujiantang-services/service/mysqldb"

	"fmt"

	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
)

// SessionDuration Session的期限
const SessionDuration = time.Minute * 20

// CreateSession 创建Session
func (j *JinmuHealth) CreateSession(ctx context.Context, req *proto.CreateSessionRequest, resp *proto.CreateSessionResponse) error {
	now := time.Now().UTC()
	sid := uuid.New()
	session := &mysqldb.Session{
		SessionID:  sid.String(),
		State:      req.Session.State,
		OpenID:     req.Session.OpenId,
		UnionID:    req.Session.UnionId,
		UserID:     req.Session.UserId,
		Authorized: req.Session.Authorized,
		ExpiredAt:  now.Add(SessionDuration).UTC(),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	session, err := j.datastore.CreateSession(ctx, session)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to create session : %s", err.Error()))
	}
	resp.Sid = session.SessionID
	expiredAt, _ := ptypes.TimestampProto(session.ExpiredAt)
	resp.ExpiredTime = expiredAt
	return nil
}

// GetSession 得到Session
func (j *JinmuHealth) GetSession(ctx context.Context, req *proto.GetSessionRequest, resp *proto.GetSessionResponse) error {
	session, err := j.datastore.FindSession(ctx, req.Sid)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find session by sid: %s", err.Error()))
	}
	expiredAt, _ := ptypes.TimestampProto(session.ExpiredAt)
	resp.Session = &proto.SessionInfo{
		State:       session.State,
		OpenId:      session.OpenID,
		UnionId:     session.UnionID,
		UserId:      session.UserID,
		Authorized:  session.Authorized,
		ExpiredTime: expiredAt,
	}
	return nil
}

// UpdateSession 更新Session
func (j *JinmuHealth) UpdateSession(ctx context.Context, req *proto.UpdateSessionRequest, resp *proto.UpdateSessionResponse) error {
	now := time.Now().UTC()
	session := &mysqldb.Session{
		State:      req.Session.State,
		OpenID:     req.Session.OpenId,
		UnionID:    req.Session.UnionId,
		UserID:     req.Session.UserId,
		Authorized: req.Session.Authorized,
		UpdatedAt:  now,
	}
	return j.datastore.UpdateSession(ctx, req.Sid, session)
}
