package wechat

import (
	"fmt"
	"net/url"
	"time"

	"gopkg.in/chanxuehong/wechat.v2/mp/qrcode"
)

// DefaultTempQrCodeExpireSeconds 默认临时二维码有效期 300 秒
const DefaultTempQrCodeExpireSeconds int = 300

// WxmpTempQrCodeURL 微信临时二维码的URL
type WxmpTempQrCodeURL struct {
	ImageURL  string    `json:"image_url"`
	RawURL    string    `json:"raw_url"`
	ExpiredAt time.Time `json:"expired_at"`
	Ticket    string    `json:"-"`
	OriID     string    `json:"-"`
	SceneID   int32     `json:"scene_id"`
}

// GetTempQrCodeUrl 获取微信临时二维码
func (u *Wxmp) GetTempQrCodeUrl(sceneID int32) (*WxmpTempQrCodeURL, error) {
	now := time.Now()

	qrcode, err := qrcode.CreateTempQrcode(u.WechatClient, sceneID, DefaultTempQrCodeExpireSeconds)
	if err != nil {
		return nil, err
	}

	qr := WxmpTempQrCodeURL{
		ImageURL:  fmt.Sprintf("https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=%s", url.QueryEscape(qrcode.Ticket)),
		RawURL:    qrcode.URL,
		ExpiredAt: now.Add(time.Duration(DefaultTempQrCodeExpireSeconds) * time.Second).UTC(),
		Ticket:    qrcode.Ticket,
		OriID:     u.Options.WxOriID,
		SceneID:   sceneID,
	}

	return &qr, nil
}
