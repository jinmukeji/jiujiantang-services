package wechat

import (
	"github.com/micro/go-micro/broker"
	"gopkg.in/chanxuehong/wechat.v2/mp/core"
	"gopkg.in/chanxuehong/wechat.v2/mp/jssdk"
	mpoauth2 "gopkg.in/chanxuehong/wechat.v2/mp/oauth2"
	"gopkg.in/chanxuehong/wechat.v2/oauth2"
)

const (
	oauth2Scope = "snsapi_userinfo" // 微信 OAuth 的 Scope
)

// WxmpOptions 微信公众平台选项参数
type WxmpOptions struct {
	WxAppID               string // 微信公众平台 开发者 ID
	WxAppSecret           string // 微信公众平台 开发者密码
	WxOriID               string // 微信公众平台 公众号原始 ID
	WxToken               string // 微信公众平台 令牌 (Token)
	WxEncodedAESKey       string // 微信公众平台 消息加解密密钥 (EncodingAESKey)
	WxCallbackServerBase  string // 微信公众平台 回调服务器的地址
	WxTemplateID          string // 微信模版ID
	WxH5ServerBase        string // 微信H5地址
	JinmuH5Serverbase     string // jinmu H5地址
	JinmuH5ServerbaseV2_0 string // jinmu V2_0 H5地址
	JinmuH5ServerbaseV2_1 string // jinmu V2_1 H5地址
}

// Wxmp 微信公众平台
type Wxmp struct {
	Options *WxmpOptions

	// FIXME: 中控服务需要定时刷新 Token

	Oauth2Endpoint    oauth2.Endpoint
	AccessTokenServer core.AccessTokenServer
	WechatClient      *core.Client
	JsTicketServer    jssdk.TicketServer

	Broker broker.Broker
}

// NewWxmp 创建Wxmp
func NewWxmp(bk broker.Broker, ops *WxmpOptions) *Wxmp {
	u := &Wxmp{
		Broker:  bk,
		Options: ops,
	}
	// u.prepare()

	return u
}

func (u *Wxmp) prepare() {
	u.Oauth2Endpoint = mpoauth2.NewEndpoint(u.Options.WxAppID, u.Options.WxAppSecret)
	u.AccessTokenServer = NewDefaultClusterAccessTokenServer(u.Options.WxAppID, u.Options.WxAppSecret, nil, DefaultClusterRefreshTokenTopic, u.Broker)
	u.WechatClient = core.NewClient(u.AccessTokenServer, nil)
	u.JsTicketServer = jssdk.NewDefaultTicketServer(u.WechatClient)
}

// GetOriginID 返回当前微信公众号的原始 ID
func (u *Wxmp) GetOriginID() string {
	return u.Options.WxOriID
}
