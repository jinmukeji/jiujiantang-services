package mail

// TODO: https://github.com/jinmukeji/issues-engineer/issues/617

import (
	"github.com/go-gomail/gomail"
)

// Client 包含了邮件服务有关的逻辑
type Client struct {
	options *Options
}

// Contact 邮件收件人 Username 是用户名 Nickname 是昵称
type Contact struct {
	Username string
	Nickname string
}

// NewMailClient 生成Client
func NewMailClient(opts ...Option) *Client {
	options := newOptions(opts...)
	return &Client{options}
}

// SendMail 发送任意内容和主题的邮件
func (c *Client) SendMail(to Contact, subject, body string, contentType string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(c.options.Username, c.options.SenderNickname))
	m.SetHeader("To", m.FormatAddress(to.Username, to.Nickname))
	m.SetHeader("Subject", subject)
	m.SetHeader("Reply-To", c.options.ReplyToAddress)
	m.SetBody(contentType, body)
	d := gomail.NewDialer(c.options.Address, c.options.Port, c.options.Username, c.options.Password)
	return d.DialAndSend(m)
}

// SetOptions 是 options 字段的 setter
func (c *Client) SetOptions(options *Options) {
	c.options = options
}
