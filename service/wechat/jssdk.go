package wechat

import (
	"fmt"
	"time"

	"github.com/jinmukeji/go-pkg/v2/crypto/rand"
	"gopkg.in/chanxuehong/wechat.v2/mp/jssdk"
)

// JsSdkSignConfig JS SDK Sign配置
type JsSdkSignConfig struct {
	AppID     string `json:"app_id"`
	Timestamp string `json:"timestamp"`
	NonceStr  string `json:"noncestr"`
	Signature string `json:"signature"`
}

// GetWxJsSdkConfig 返回微信 JS SDK 配置信息
func (u *Wxmp) GetWxJsSdkConfig(url string) (*JsSdkSignConfig, error) {

	// https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421141115

	// 获取 ticket
	ticket, err := u.JsTicketServer.Ticket()
	if err != nil {
		return nil, err
	}

	nonceStr := generateNonceStr()
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	signature := jssdk.WXConfigSign(ticket, nonceStr, timestamp, url)

	cfg := JsSdkSignConfig{
		AppID:     u.Options.WxAppID,
		Timestamp: timestamp,
		NonceStr:  nonceStr,
		Signature: signature,
	}

	return &cfg, nil
}

func generateNonceStr() string {
	v, _ := rand.RandomStringWithMask(rand.MaskLetterDigits, 16)
	return v
}
