package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var (
	baseSendTextURL = "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=%s"
	contentType     = "application/json"
)

// TextMessage 文字消息
type TextMessage struct {
	Touser  string      `json:"touser"`
	Msgtype string      `json:"msgtype"`
	Text    TextContent `json:"text"`
}

// TextContent 文字内容
type TextContent struct {
	Content string `json:"content"`
}

// WechatSendTextMessage 微信发送文字消息
func (u *Wxmp) WechatSendTextMessage(openID string, content string) error {
	textMessage := &TextMessage{
		Touser:  openID,
		Msgtype: "text",
		Text: TextContent{
			Content: content,
		},
	}
	msg, errMarshal := json.Marshal(textMessage)
	if errMarshal != nil {
		return errMarshal
	}
	token, errToken := u.AccessTokenServer.Token()
	if errToken != nil {
		return errToken
	}
	serverURL := fmt.Sprintf(baseSendTextURL, url.QueryEscape(token))
	_, err := http.Post(serverURL, contentType, bytes.NewReader(msg))
	if err != nil {
		return err
	}
	return nil
}
