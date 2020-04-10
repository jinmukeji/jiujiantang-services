package wechat

import (
	mpoauth2 "gopkg.in/chanxuehong/wechat.v2/mp/oauth2"
	wxuser "gopkg.in/chanxuehong/wechat.v2/mp/user"
	"gopkg.in/chanxuehong/wechat.v2/oauth2"
)

// GetUserInfo 根据微信 OpenID 获得用户信息
func (u *Wxmp) GetUserInfo(openID string) (*wxuser.UserInfo, error) {
	const lang = "zh_CN"
	return wxuser.Get(u.WechatClient, openID, lang)
}

// GetUserInfoByOauthCode 根据微信 OAuth 回调返回的 Code 获取对应微信用户的信息
func (u *Wxmp) GetUserInfoByOauthCode(code string) (*mpoauth2.UserInfo, error) {
	oauth2Client := oauth2.Client{
		Endpoint: u.Oauth2Endpoint,
	}
	token, err := oauth2Client.ExchangeToken(code)
	if err != nil {
		return nil, err
	}

	// 微信第三方登录验证成功，获取微信用户信息
	return mpoauth2.GetUserInfo(token.AccessToken, token.OpenId, "", nil)
}
