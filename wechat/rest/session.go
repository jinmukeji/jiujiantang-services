package rest

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang/protobuf/ptypes"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

// domain cookie中的domain
const domain = "jinmuhealth.com"

// Session 会话
type Session struct {
	SID        string // Session ID
	State      string // 微信 OAuth 验证的 state
	WxOpenID   string // 微信 OpenID
	WxUnionID  string // 微信 UnionID
	UserID     int64  // 用户 ID
	Authorized bool   // 是否已经验证通过
	IsExpired  bool   // 是否已经到期
}

func (h *handler) getSession(ctx iris.Context) (*Session, error) {
	// 从 Cookie 获取 sid
	sid := ctx.GetCookie("sid", CookieDomain(domain))
	if sid == "" {
		return nil, errors.New("missing sid")
	}
	req := new(proto.GetSessionRequest)
	req.Sid = sid
	resp, err := h.rpcSvc.GetSession(
		newRPCContext(ctx), req,
	)
	if err != nil {
		return nil, err
	}
	expiredAt, err := ptypes.Timestamp(resp.Session.ExpiredTime)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	IsExpired := false
	if expiredAt.Before(now) {
		IsExpired = true
	}
	return &Session{
		SID:        sid,
		State:      resp.Session.State,
		WxOpenID:   resp.Session.OpenId,
		WxUnionID:  resp.Session.UnionId,
		UserID:     resp.Session.UserId,
		Authorized: resp.Session.Authorized,
		IsExpired:  IsExpired,
	}, nil
}

func (h *handler) addSession(ctx iris.Context, session *Session) (string, error) {
	req := new(proto.CreateSessionRequest)
	req.Session = &proto.SessionInfo{
		State:      session.State,
		OpenId:     session.WxOpenID,
		UnionId:    session.WxUnionID,
		UserId:     session.UserID,
		Authorized: session.Authorized,
	}
	resp, err := h.rpcSvc.CreateSession(
		newRPCContext(ctx), req,
	)
	if err != nil {
		return "", err
	}
	return resp.Sid, nil
}

func (h *handler) updateSession(ctx iris.Context, sid string, session *Session) error {
	req := new(proto.UpdateSessionRequest)
	req.Sid = sid
	req.Session = &proto.SessionInfo{
		State:      session.State,
		OpenId:     session.WxOpenID,
		UnionId:    session.WxUnionID,
		UserId:     session.UserID,
		Authorized: session.Authorized,
	}
	_, err := h.rpcSvc.UpdateSession(
		newRPCContext(ctx), req,
	)
	if err != nil {
		return err
	}
	return nil
}

// CookieDomain 设置cookie中的demain
func CookieDomain(domain string) context.CookieOption {
	return func(c *http.Cookie) {
		c.Domain = domain
	}
}
