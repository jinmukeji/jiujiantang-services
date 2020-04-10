package sem

import (
	"errors"
	"fmt"
	"net/smtp"
	"strings"
)

// HOST smtp地址及端口
const HOST = "smtp.163.com:25"

// SendMail 网易暂不支持通过API发送邮件，网易163邮箱发送邮件的逻辑函数
func SendMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var contentType string
	if mailtype == "html" {
		contentType = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		contentType = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	sendTo := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, sendTo, msg)
	return err
}

// NetEaseSEMClient 网易邮件Client
type NetEaseSEMClient struct {
	USER   string
	PASSWD string
}

// NewNetEaseSEMClient 生成SEMClient
func NewNetEaseSEMClient(user, password string) (*NetEaseSEMClient, error) {

	if user == "" {
		return nil, errors.New("NewNetEase User should be not empty")
	}
	if password == "" {
		return nil, errors.New("NewNetEase Password should be not empty")
	}
	return &NetEaseSEMClient{
		USER:   user,
		PASSWD: password,
	}, nil
}

// SendEmail 网易163发送频率不可以太高,否则会被判断为垃圾邮件无法发出
func (client *NetEaseSEMClient) SendEmail(toAddress string, templateAction TemplateAction, language TemplateLanguage, templateParam map[string]string) (bool, error) {
	replaces := getParams(templateAction, language)
	if !strings.Contains(toAddress, ";") {
		body := fmt.Sprintf(getLanguageHTML(language), toAddress, replaces[0], replaces[1], convertTemplateParam(templateParam))
		err := SendMail(client.USER, client.PASSWD, HOST, toAddress, replaces[3], body, "html")
		if err != nil {
			return false, err
		}
		return true, nil
	}
	for _, v := range strings.Split(toAddress, ";") {
		body := fmt.Sprintf(getLanguageHTML(language), v, replaces[0], replaces[1], convertTemplateParam(templateParam))
		err := SendMail(client.USER, client.PASSWD, HOST, v, replaces[3], body, "html")
		if err != nil {
			return false, err
		}
	}
	return true, nil
}
