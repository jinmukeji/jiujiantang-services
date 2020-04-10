package wechat

import (
	"fmt"
	"net/url"

	mpoauth2 "gopkg.in/chanxuehong/wechat.v2/mp/oauth2"
)

// BuildOAuthURL 构建一个发起微信 OAuth 请求 URL
func (u *Wxmp) BuildOAuthURL(redirectURL, state string) string {
	fullURL := u.buildOAuthRedirectUrl(redirectURL)
	return mpoauth2.AuthCodeURL(u.Options.WxAppID, fullURL, oauth2Scope, state)
}

// buildOAuthRedirectUrl 构建微信OAuth转发URL
func (u *Wxmp) buildOAuthRedirectUrl(redirectPath string) string {
	return fmt.Sprintf("%s/wx/oauth/callback?redirect=%s", u.Options.WxCallbackServerBase, url.QueryEscape(redirectPath))
}
