package rest

import (
	"fmt"

	"github.com/jinmukeji/go-pkg/crypto/rand"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	"github.com/kataras/iris/v12"

	"net/url"
)

// WeChatAuth 微信公众号服务器接入配置验证回调
func (h *handler) WeChatOAuth(ctx iris.Context) {
	state := generateState()
	session := Session{
		State: state,
	}
	sid, err := h.addSession(ctx, &session)
	if err != nil {
		writeError(ctx, wrapError(ErrRPCInternal, "", err), false)
		log.Println(err)
		return
	}
	ctx.SetCookieKV("sid", sid, CookieDomain(domain))

	// 302 跳转微信 OAuth 授权页面
	redirectOnSuccessfulAuth := ctx.URLParam("redirect")

	req := new(proto.WechatBuildOAuthURLRequest)
	req.AuthRedirectUrl = redirectOnSuccessfulAuth
	req.State = state

	resp, err := h.rpcSvc.WechatBuildOAuthURL(
		newRPCContext(ctx), req,
	)
	if err != nil {
		writeError(ctx, wrapError(ErrRPCInternal, "", err), false)
		log.Println(err)
		return
	}

	ctx.Redirect(resp.AuthCodeUrl, iris.StatusFound)
}

// generateState 生成微信 OAuth 登录用的 state
func generateState() string {
	v, _ := rand.RandomStringWithMask(rand.MaskLetterDigits, 16)
	return v
}

// WeChatOAuthCallback 微信 OAuth 登录的回调响应
func (h *handler) WeChatOAuthCallback(ctx iris.Context) {
	session, err := h.getSession(ctx)
	if err != nil {
		_, errText := ctx.Text("非法的会话")
		if errText != nil {
			return
		}
		log.Warnln("invalid session:", err)
		return
	}
	// Session 过期
	if session.IsExpired {
		_, errText := ctx.Text("非法的会话,session过期")
		if errText != nil {
			return
		}
		log.Warnln("非法的会话,session过期")
		return
	}
	// 检查 code
	code := ctx.URLParam("code")
	if code == "" {
		_, errText := ctx.Text("用户禁止授权")
		if errText != nil {
			return
		}
		log.Warnln("用户禁止授权")
		return
	}

	// 检查 state
	savedState := session.State
	queryState := ctx.URLParam("state")
	if queryState == "" {
		_, errText := ctx.Text("非法的会话，state错误")
		if errText != nil {
			return
		}
		log.Warnln("state 参数为空")
		return
	}
	if savedState != queryState {
		_, errText := ctx.Text("非法的会话，state 错误")
		if errText != nil {
			return
		}
		str := fmt.Sprintf("state 不匹配, session 中为 %q, url 传递过来的是 %q", savedState, queryState)
		log.Warnln(str)
		return
	}
	req := new(proto.WechatGetUserInfoRequest)
	req.Code = code

	wxUser, err := h.rpcSvc.WechatGetUserInfo(
		newRPCContext(ctx), req,
	)
	if err != nil {
		writeError(ctx, wrapError(ErrWxOAuth, "", err), false)
		log.Println(err)
		return
	}
	session.Authorized = true
	session.WxOpenID = wxUser.OpenId
	session.WxUnionID = wxUser.UnionId

	reqWxUser := new(proto.GetWechatUserRequest)
	reqWxUser.OpenId = wxUser.OpenId
	resp, err := h.rpcSvc.GetWechatUser(
		newRPCContext(ctx), reqWxUser,
	)
	if err != nil {
		writeError(ctx, wrapError(ErrWxOAuth, "", err), false)
		log.Println(err)
		return
	}
	session.UserID = int64(resp.UserId)

	err = h.updateSession(ctx, session.SID, session)
	if err != nil {
		writeError(ctx, wrapError(ErrWxOAuth, "", err), false)
		log.Println(err)
		return
	}
	redirectOnSuccessfulAuth := ctx.URLParam("redirect")
	url, errParse := url.Parse(redirectOnSuccessfulAuth)
	if errParse != nil {
		writeError(ctx, wrapError(ErrWxOAuth, "", errParse), false)
		log.Println(errParse)
		return
	}
	q := url.Query()
	rs, _ := rand.RandomStringWithMask(rand.MaskLetterDigits, 4)
	q.Set("random", rs)
	url.RawQuery = q.Encode()
	// 302 跳转
	ctx.Redirect(url.String(), iris.StatusFound)
}
