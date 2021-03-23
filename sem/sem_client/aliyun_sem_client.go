package sem

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	// DysemAPIEndpoint 阿里云邮件接口URL
	DysemAPIEndpoint = "https://dm.aliyuncs.com/?"
	// SimplifiedChineseContent 简体邮件html基本内容
	SimplifiedChineseContent = `<html><head><meta http-equiv="Content-Type" content="text/html; charset=UTF-8"><meta http-equiv="Content-Style-Type" content="text/css"><meta name="generator" content="Aspose.Words for .NET 15.1.0.0"><title></title><div></div></head><body><div><p style="margin:0pt"><span style="font-family:'Times New Roman'; font-size:12pt">&nbsp; 亲爱的%s，您好：</span></p><p style="margin:0pt"><span style="font-family:'Times New Roman'; font-size:12pt">&nbsp;&nbsp;&nbsp;&nbsp; 您正在%s</span></p><p style="margin:0pt"><span style="font-family:'Times New Roman'; font-size:12pt">&nbsp;&nbsp;&nbsp;&nbsp; 邮箱验证码为：%s</span></p><p style="margin:0pt"><span style="font-family:'Times New Roman'; font-size:12pt">&nbsp;&nbsp;&nbsp;&nbsp; 验证码30分钟内有效</span></p><p style="margin:0pt"><span style="font-family:'Times New Roman'; font-size:12pt">&nbsp;&nbsp;&nbsp;&nbsp; 如非本人操作，可能密码已泄露，请尽快登录并更改密码。</span></p><p style="margin:0pt"><span style="font-family:'Times New Roman'; font-size:12pt">&nbsp;&nbsp;&nbsp;&nbsp; 本邮件由系统自动发送，请勿直接回复。</span></p><p style="margin:0pt"><span style="font-family:'Times New Roman'; font-size:12pt">&nbsp;</span></p></div></body></html>`
	// TraditionalChineseContent 繁体邮件html基本内容
	TraditionalChineseContent = `<html><head><meta http-equiv="Content-Type" content="text/html; charset=UTF-8"><meta http-equiv="Content-Style-Type" content="text/css"><meta name="generator" content="Aspose.Words for .NET 15.1.0.0"><title></title></head><body><div><p style="margin:0pt"><span style="font-family:'Times New Roman'; font-size:12pt">&nbsp; 親愛的%s，您好：</span></p><p style="margin:0pt"><span style="font-family:'Times New Roman'; font-size:12pt">&nbsp;&nbsp;&nbsp;&nbsp; 您正在%s</span></p><p style="margin:0pt"><span style="font-family:'Times New Roman'; font-size:12pt">&nbsp;&nbsp;&nbsp;&nbsp; 郵箱驗證碼為：%s</span></p><p style="margin:0pt"><span style="font-family:'Times New Roman'; font-size:12pt">&nbsp;&nbsp;&nbsp;&nbsp; 驗證碼30分鐘內有效</span></p><p style="margin:0pt"><span style="font-family:'Times New Roman'; font-size:12pt">&nbsp;&nbsp;&nbsp;&nbsp; 如非本人操作，可能密碼已洩露，請盡快登錄并更改密碼。</span></p><p style="margin:0pt"><span style="font-family:'Times New Roman'; font-size:12pt">&nbsp;&nbsp;&nbsp;&nbsp; 本郵件由系統自動發送，請勿直接回復。</span></p><p style="margin:0pt"><span style="font-family:'Times New Roman'; font-size:12pt">&nbsp;</span></p></div></body></html>`
	// EnglishContent 英文邮件html基本内容
	EnglishContent = `<html><head><meta http-equiv="Content-Type" content="text/html; charset=UTF-8"><meta http-equiv="Content-Style-Type" content="text/css"><meta name="generator" content="Aspose.Words for .NET 15.1.0.0"><title></title></head><body><div><p style="margin:0pt"><span style="font-family:'Times New Roman'; font-size:12pt">&nbsp; Hi %s，</span></p><p style="margin:0pt"><span style="font-family:'Times New Roman'; font-size:12pt">&nbsp; We received a request to %s.</span></p><p style="margin:0pt"><span style="font-family:'Times New Roman'; font-size:12pt">&nbsp; Use verification code: %s within 30mins.</span></p><p style="margin:0pt"><span style="font-family:'Times New Roman'; font-size:12pt">&nbsp; This email can't receive replies.</span></p><p style="margin:0pt"><span style="font-family:'Times New Roman'; font-size:12pt">&nbsp; For more information,visit http://www.jinmuhealth.com</span></p><p style="margin:0pt"><span style="font-family:'Times New Roman'; font-size:12pt">&nbsp;</span></p></div></body></html>`
	// Method 请求方法
	Method = "GET"
	// AccountName 管理控制台中配置的发信地址
	AccountName = "testsend@noreply.jinmuhealth.com"
)

// SendSemReply 发送邮件返回
type SendSemReply struct {
	RequestID string `json:"RequestId,omitempty"` // 请求ID
	HostID    string `json:"HostId,omitempty"`    // 主机ID
	Code      string `json:"Code,omitempty"`      // 状态码
	Message   string `json:"Message,omitempty"`   // 状态码的描述
}

// SingleTemplateParam 触发邮件模板变量
type SingleTemplateParam struct {
	Code string `json:"code,omitempty"` // 对应模版中的 ${code}
}

// 简体版邮件模板参数
var simplifiedChineseVariable = map[TemplateAction][]string{
	FindResetPassword: {"重置喜马把脉ID密码", "喜马把脉平台", "重置喜马把脉ID密码"},   // 找回/重置密码
	FindUsername:      {"找回喜马把脉ID用户名", "喜马把脉平台", "找回喜马把脉ID用户名"}, // 找回用户名
	SetSecureEmail:    {"修改您的安全邮箱", "喜马把脉平台", "修改喜马把脉ID安全邮箱"},   // 设置安全邮箱
	ModifySecureEmail: {"修改您的安全邮箱", "喜马把脉平台", "修改喜马把脉ID安全邮箱"},   // 修改安全邮箱
	UnsetSecureEmail:  {"修改您的安全邮箱", "喜马把脉平台", "修改喜马把脉ID安全邮箱"},   // 解绑安全邮箱
}

// 繁体版邮件模板参数
var traditionalChineseVariable = map[TemplateAction][]string{
	FindResetPassword: {"重置喜马把脉ID密碼", "喜马把脉平台", "重置喜马把脉ID密碼"},   // 找回/重置密码
	FindUsername:      {"找回喜马把脉ID用户名", "喜马把脉平台", "找回喜马把脉ID用戶名"}, // 找回用户名
	SetSecureEmail:    {"修改您的安全郵箱", "喜马把脉平台", "修改喜马把脉ID安全郵箱"},   // 设置安全邮箱
	ModifySecureEmail: {"修改您的安全郵箱", "喜马把脉平台", "修改喜马把脉ID安全郵箱"},   // 修改安全邮箱
	UnsetSecureEmail:  {"修改您的安全郵箱", "喜马把脉平台", "修改喜马把脉ID安全郵箱"},   // 解绑安全邮箱
}

// 英文版邮件模板参数
var englishVariable = map[TemplateAction][]string{
	FindResetPassword: {"reset your password", "Jinmu Platform", "Reset Password"},               // 找回/重置密码
	FindUsername:      {"retrieve your username", "Jinmu Platform", "Retrieve Username"},         // 找回用户名
	SetSecureEmail:    {"change your security email", "Jinmu Platform", "Change Security Email"}, // 设置安全邮箱
	ModifySecureEmail: {"change your security email", "Jinmu Platform", "Change Security Email"}, // 修改安全邮箱
	UnsetSecureEmail:  {"change your security email", "Jinmu Platform", "Change Security Email"}, // 解绑安全邮箱
}

// AliyunSEMClient 阿里云邮箱Client
type AliyunSEMClient struct {
	AccessKeyID     string
	AccessKeySecret string
}

// NewAliyunSEMClient 生成SEMClient
func NewAliyunSEMClient(accessKeyID, accessKeySecret string) (*AliyunSEMClient, error) {

	if accessKeyID == "" {
		return nil, errors.New("Aliyun AccessKeyId should be not empty")
	}
	if accessKeySecret == "" {
		return nil, errors.New("Aliyun AccessKeySecret should be not empty")
	}
	return &AliyunSEMClient{
		AccessKeyID:     accessKeyID,
		AccessKeySecret: accessKeySecret,
	}, nil
}

// convertTemplateParam 转化成阿里的所需的模版参数
func convertTemplateParam(templateParam map[string]string) string {
	params := &SingleTemplateParam{
		Code: templateParam["code"],
	}
	return params.Code
}

func getParams(templateAction TemplateAction, language TemplateLanguage) []string {
	switch language {
	case SimplifiedChinese:
		return simplifiedChineseVariable[templateAction]
	case TraditionalChinese:
		return traditionalChineseVariable[templateAction]
	case English:
		return englishVariable[templateAction]
	}
	return simplifiedChineseVariable[templateAction]
}

func getLanguageHTML(language TemplateLanguage) string {
	switch language {
	case SimplifiedChinese:
		return SimplifiedChineseContent
	case TraditionalChinese:
		return TraditionalChineseContent
	case English:
		return EnglishContent
	}
	return SimplifiedChineseContent
}

// SendEmail 发送邮件
func (client *AliyunSEMClient) SendEmail(toAddress string, templateAction TemplateAction, language TemplateLanguage, templateParam map[string]string) (bool, error) {
	params := map[string]string{
		"Format":           "JSON",                                          // 返回值的类型
		"Version":          "2015-11-23",                                    // API 版本号
		"AccessKeyId":      "",                                              // 阿里云颁发给用户的访问服务所用的密钥ID
		"Signature":        "",                                              // 签名结果串
		"SignatureMethod":  "HMAC-SHA1",                                     // 签名方式
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"), // 请求的时间戳
		"SignatureVersion": "1.0",                                           // 签名算法版本
		"SignatureNonce":   uuid.New().String(),                             // 唯一随机数
		"RegionID":         "",                                              // 机房信息.可选
		"Action":           "SingleSendMail",                                // 操作接口名
		"AccountName":      AccountName,                                     // 管理控制台中配置的发信地址
		"ReplyToAddress":   "false",                                         // 使用管理控制台中配置的回信地址 bool类型
		"AddressType":      "0",                                             // 取值范围 0~1: 0 为随机账号；1 为发信地址
		"ToAddress":        "",                                              // 目标地址,可以有多个
		"FromAlias":        "",                                              // 发信人昵称,可选
		"Subject":          "",                                              // 邮件主题,可选
		"HtmlBody":         "",                                              // 邮件 html 正文,可选
		"TxtBody":          "",                                              // 邮件 text 正文,可选
		"ClickTrace":       "0",                                             // 是否打开数字跟踪功能,可选
	}

	params["ToAddress"] = toAddress
	params["AccessKeyId"] = client.AccessKeyID
	replaces := getParams(templateAction, language)
	params["FromAlias"] = replaces[1]
	params["Subject"] = replaces[2]
	if !strings.Contains(toAddress, ",") {
		params["ToAddress"] = toAddress
		params["HtmlBody"] = fmt.Sprintf(getLanguageHTML(language), toAddress, replaces[0], convertTemplateParam(templateParam))
		return SendSem(client.AccessKeySecret, params)
	}
	for _, v := range strings.Split(toAddress, ",") {
		params["ToAddress"] = v
		params["HtmlBody"] = fmt.Sprintf(getLanguageHTML(language), v, replaces[0], convertTemplateParam(templateParam))
		params["SignatureNonce"] = uuid.New().String()
		ok, err := SendSem(client.AccessKeySecret, params)
		if !ok {
			return false, err
		}
	}
	return true, nil
}

// BatchSendEmail 批量邮件邮件
func (client *AliyunSEMClient) BatchSendEmail(TemplateName, ReceiversName, TagName string) (bool, error) {
	params := map[string]string{
		"Format":           "JSON",                                          // 返回值的类型
		"Version":          "2015-11-23",                                    // API 版本号
		"AccessKeyId":      "",                                              // 阿里云颁发给用户的访问服务所用的密钥ID
		"Signature":        "",                                              // 签名结果串
		"SignatureMethod":  "HMAC-SHA1",                                     // 签名方式
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"), // 请求的时间戳
		"SignatureVersion": "1.0",                                           // 签名算法版本
		"SignatureNonce":   uuid.New().String(),                             // 唯一随机数
		"RegionID":         "",                                              // 机房信息.可选
		"Action":           "BatchSendMail",                                 // 操作接口名
		"AccountName":      AccountName,                                     // 管理控制台中配置的发信地址
		"AddressType":      "0",                                             // 取值范围 0~1: 0 为随机账号；1 为发信地址
		"TemplateName":     "",                                              // 预先创建且通过审核的模板名称
		"ReceiversName":    "",                                              // 预先创建且上传了收件人的收件人列表名称
		"TagName":          "",                                              // 邮件标签名称,可选
		"ClickTrace":       "0",                                             // 是否打开数字跟踪功能,可选
	}
	params["AccessKeyId"] = client.AccessKeyID
	params["TemplateName"] = TemplateName
	params["ReceiversName"] = ReceiversName
	params["TagName"] = TagName
	return SendSem(client.AccessKeySecret, params)
}

// SendSem 发送邮件
func SendSem(accessKeySecret string, params map[string]string) (bool, error) {

	sign, err := getSHA1Signature(accessKeySecret, params)
	if err != nil {
		return false, err
	}
	params["Signature"] = sign
	parsedParams := getParamsStr(params)
	singleurl := DysemAPIEndpoint + parsedParams[1:]
	resp, err := http.Get(singleurl)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	ssr := &SendSemReply{}
	if err := json.Unmarshal(body, ssr); err != nil {
		return false, err
	}
	if ssr.Code != "" {
		return false, errors.New(ssr.Code)
	}
	return true, nil
}

// getParamsStr 得到请求参数
func getParamsStr(params map[string]string) string {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	// 根据参数Key排序
	sort.Strings(keys)
	var sortQueryStringTmp string
	for _, value := range keys {
		if params[value] == "" {
			continue
		}
		// specialUrlEncode(参数Key) + "=" + specialUrlEncode(参数值)
		sortQueryStringTmp = fmt.Sprintf("%s&%s=%s", sortQueryStringTmp, specialURLEncode(value), specialURLEncode(params[value]))
	}
	return sortQueryStringTmp
}

// getSHA1Signature 计算签名
func getSHA1Signature(accessKeySecret string, params map[string]string) (string, error) {
	stringToSign := Method + "&" + specialURLEncode("/") + "&" + specialURLEncode(getSHA1ParamsStr(params))
	sign, err := getSign(accessKeySecret, stringToSign)
	if err != nil {
		return "", err
	}
	return sign, nil
}

// getParamsStr 得到用于计算签名时的请求参数
func getSHA1ParamsStr(params map[string]string) string {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	// 根据参数Key排序
	sort.Strings(keys)
	var sortQueryStringTmp string
	for _, value := range keys {
		if params[value] == "" || value == "Signature" {
			continue
		}
		// specialUrlEncode(参数Key) + "=" + specialUrlEncode(参数值)
		sortQueryStringTmp = fmt.Sprintf("%s&%s=%s", sortQueryStringTmp, specialURLEncode(value), specialURLEncode(params[value]))
	}
	return sortQueryStringTmp[1:]
}

// specialURLEncode 处理URL编码
func specialURLEncode(in string) string {
	// URLEncode后 加号（+）替换成 %20、星号（*）替换成 %2A、%7E 替换回波浪号（~）
	return strings.NewReplacer("+", "%20", "*", "%2A", "%7E", "~").Replace(url.QueryEscape(in))
}

// getSign 计算签名
func getSign(accessSecret, stringToSign string) (string, error) {
	mac := hmac.New(sha1.New, []byte(fmt.Sprintf("%s&", accessSecret)))
	_, err := mac.Write([]byte(stringToSign))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}
