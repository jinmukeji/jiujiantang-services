package mail

import (
	"bytes"
	"html/template"
)

// SendMailRetrievePassword 发送找回密码邮件
func (c *Client) SendMailRetrievePassword(to Contact, password string) error {
	// 生成内容
	var content string
	var err error
	if content, err = c.MailContentRetrievePassword(InfoRetrievePassword{to.Nickname, password}); err != nil {
		return err
	}
	return c.SendMail(to, "找回密码", content, textHTML)
}

// InfoRetrievePassword 是用于生成模版的结构体
type InfoRetrievePassword struct {
	AccountName string
	Password    string
}

// MailContentRetrievePassword 生成邮件主题内容
func (c *Client) MailContentRetrievePassword(info InfoRetrievePassword) (string, error) {
	buf := bytes.NewBufferString("")
	tpl, err := template.New("mailContent").Parse(`<h3>尊敬的 {{.AccountName}}</h3><p>您好，系统已收到您的找回密码申请，以下为你的脉诊仪账户所需要的重要信息，请注意保密。</p><p>您的脉诊仪账号密码为：{{.Password}}</p>
		<p align='center'><img width= '200pt' height ='200pt' src= 'http://www.jinmuhealth.com/assets/img/QRcode/wechatQRcode.jpg',  alt='喜马把脉科技'/></p>
		<p align='center'>扫一扫关注“喜马把脉科技”</p>
		<p>请注意，该邮件地址不接收回复邮件，请联系喜马把脉客服进行咨询</p>
		<p>客服电话：0519-81180075</p>
		<p>E-mail:information@jinmuhealth.com</p>`)
	if err != nil {
		return "", err
	}
	if err := tpl.Execute(buf, info); err != nil {
		return "", err
	}
	return buf.String(), nil
}
