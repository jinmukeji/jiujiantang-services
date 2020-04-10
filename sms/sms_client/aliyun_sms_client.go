package sms

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"crypto/hmac"
	"crypto/sha1"

	"github.com/google/uuid"
)

const (
	// DysmsAPIEndpoint 阿里云短信接口URL
	DysmsAPIEndpoint = "http://dysmsapi.aliyuncs.com"
)

// simpleChineseTemplateAction 国内简体模版
var simpleChineseTemplateAction = map[TemplateAction]string{
	SignUp:            "SMS_157065159", // 手机号注册
	SignIn:            "SMS_157070152", // 手机号登录
	ResetPassword:     "SMS_157065163", // 找回/重置密码
	SetPhoneNumber:    "SMS_157070159", // 设置手机号
	ModifyPhoneNumber: "SMS_157070166", // 修改手机号
}

// simpleChineseInternationalTemplateAction 国际简体模版
var simpleChineseInternationalTemplateAction = map[TemplateAction]string{
	SignUp:            "SMS_157065508", // 手机号注册
	SignIn:            "SMS_157065512", // 手机号登录
	ResetPassword:     "SMS_157070484", // 找回/重置密码
	SetPhoneNumber:    "SMS_157065515", // 设置手机号
	ModifyPhoneNumber: "SMS_157070490", // 修改手机号
}

// 阿里的非国际短信,仅支持简体和中文
// traditionalChineseInternationalTemplateAction 国际繁体模版
var traditionalChineseInternationalTemplateAction = map[TemplateAction]string{
	SignUp:            "SMS_157065521", // 手机号注册
	SignIn:            "SMS_157070497", // 手机号登录
	ResetPassword:     "SMS_157070500", // 找回/重置密码
	SetPhoneNumber:    "SMS_157065530", // 设置手机号
	ModifyPhoneNumber: "SMS_157070508", // 修改手机号
}

// englishTemplateAction 国内英文模版
var englishTemplateAction = map[TemplateAction]string{
	SignUp:            "SMS_157065661", // 手机号注册
	SignIn:            "SMS_157065666", // 手机号登录
	ResetPassword:     "SMS_157065672", // 找回/重置密码
	SetPhoneNumber:    "SMS_157070840", // 设置手机号
	ModifyPhoneNumber: "SMS_157070634", // 修改手机号
}

// englishInternationalTemplateAction 国际英文模版
var englishInternationalTemplateAction = map[TemplateAction]string{
	SignUp:            "SMS_157070643", // 手机号注册
	SignIn:            "SMS_157065690", // 手机号登录
	ResetPassword:     "SMS_157065703", // 找回/重置密码
	SetPhoneNumber:    "SMS_157070663", // 设置手机号
	ModifyPhoneNumber: "SMS_157070666", // 修改手机号
}

// SendSmsReply 发送短信返回
type SendSmsReply struct {
	Code      string `json:"Code,omitempty"`      // 状态码
	Message   string `json:"Message,omitempty"`   // 状态码的描述
	RequestID string `json:"RequestId,omitempty"` // 请求ID
	BizID     string `json:"BizId,omitempty"`     // 发送回执ID
}

// TemplateParam 短信模板变量
type TemplateParam struct {
	Code string `json:"code,omitempty"` // 对应模版中的 ${code}
}

// AliyunSMSClient 阿里云SMS客户端
type AliyunSMSClient struct {
	AccessKeyID     string
	AccessKeySecret string
}

// NewAliyunSMSClient 生成SMSClient
func NewAliyunSMSClient(accessKeyID, accessKeySecret string) (*AliyunSMSClient, error) {
	if accessKeyID == "" {
		return nil, errors.New("AccessKeyId should be not empty")
	}
	if accessKeySecret == "" {
		return nil, errors.New("SecretAccessKey should be not empty")
	}
	return &AliyunSMSClient{
		AccessKeyID:     accessKeyID,
		AccessKeySecret: accessKeySecret,
	}, nil
}

// SendSms 发送短信
func (client *AliyunSMSClient) SendSms(phoneNumber, nationCode string, templateAction TemplateAction, language TemplateLanguage, templateParam map[string]string) (bool, error) {
	templateParamJSON := convertTemplateParam(templateParam)
	// 转化成阿里云短信模版
	templateCode := convertTemplateAction(nationCode, templateAction, language)
	// 转化成阿里云所需的手机号
	phoneNumber = joinPhoneNumber(phoneNumber, nationCode)

	sortQueryStringTmp, errGetSortQueryString := client.getSortQueryStringTmp(
		phoneNumber, templateCode, templateParamJSON, language)

	if errGetSortQueryString != nil {
		return false, errGetSortQueryString
	}
	// 去除第一个多余的&符号
	sortedQueryString := sortQueryStringTmp[1:]
	// * HTTPMethod + “&” + specialUrlEncode(“/”) + ”&” + specialUrlEncode(sortedQueryString)
	stringToSign := fmt.Sprintf("GET&%s&%s", specialURLEncode("/"), specialURLEncode(sortedQueryString))

	sign, errGetSign := getSign(client.AccessKeySecret, stringToSign)
	if errGetSign != nil {
		return false, errGetSign
	}
	// 签名最后也要做特殊URL编码
	signature := specialURLEncode(sign)

	return sendDySmsAPI(signature, sortQueryStringTmp)
}

// convertTemplateParam 转化成阿里的所需的模版参数
func convertTemplateParam(templateParam map[string]string) string {
	params := &TemplateParam{
		Code: templateParam["code"],
	}
	jsonTemplateParam, _ := json.Marshal(params)
	return string(jsonTemplateParam)
}

// convertTemplateCode 转化成阿里的模版ID 如02转---->SMS_151997087
func convertTemplateAction(nationCode string, templateAction TemplateAction, language TemplateLanguage) string {
	if nationCode == "" || nationCode == "+86" {
		switch language {
		case SimpleChinese:
			return simpleChineseTemplateAction[templateAction]
		case English:
			return englishTemplateAction[templateAction]
		}
		return simpleChineseTemplateAction[templateAction]
	}
	switch language {
	case SimpleChinese:
		return simpleChineseInternationalTemplateAction[templateAction]
	case TraditionalChinese:
		return traditionalChineseInternationalTemplateAction[templateAction]
	case English:
		return englishInternationalTemplateAction[templateAction]
	}
	return simpleChineseInternationalTemplateAction[templateAction]
}

// joinPhoneNumber 拼接国际号码
func joinPhoneNumber(phoneNumber, nationCode string) string {
	// phoneNumber 的0是需要去掉再上传到发送平台
	if nationCode != "" && strings.HasPrefix(phoneNumber, "0") {
		phoneNumber = phoneNumber[1:]
	}
	if nationCode == "+86" || nationCode == "86" {
		return phoneNumber
	}
	return dealNationCode(nationCode) + phoneNumber
}

// getSortQueryStringTmp 得到sortQueryStringTmp
func (client *AliyunSMSClient) getSortQueryStringTmp(phoneNumber, templateCode, templateParamJSON string, language TemplateLanguage) (string, error) {
	params := map[string]string{
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureNonce":   uuid.New().String(),
		"AccessKeyId":      client.AccessKeyID,
		"SignatureVersion": "1.0",
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"Format":           "JSON",
		"Action":           "SendSms",
		"Version":          "2017-05-25",
		"RegionId":         "cn-hangzhou",
		"PhoneNumbers":     phoneNumber,
		"SignName":         "金姆平台",
		"TemplateParam":    templateParamJSON,
		"TemplateCode":     templateCode,
	}
	if language == English {
		params["SignName"] = "Jinmu"
	}
	var keys []string

	for k := range params {
		keys = append(keys, k)
	}
	// 根据参数Key排序
	sort.Strings(keys)

	var sortQueryStringTmp string

	for _, value := range keys {
		// specialUrlEncode(参数Key) + "=" + specialUrlEncode(参数值)
		sortQueryStringTmp = fmt.Sprintf("%s&%s=%s", sortQueryStringTmp, specialURLEncode(value), specialURLEncode(params[value]))
	}
	return sortQueryStringTmp, nil
}

// sendDySmsAPI 发送阿里云短信api
func sendDySmsAPI(signature, sortQueryStringTmp string) (bool, error) {
	// http://dysmsapi.aliyuncs.com/?Signature=" + signature + sortQueryStringTmp
	str := fmt.Sprintf("%s/?Signature=%s%s", DysmsAPIEndpoint, signature, sortQueryStringTmp)

	resp, err := http.Get(str)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	ssr := &SendSmsReply{}

	if err := json.Unmarshal(body, ssr); err != nil {
		return false, err
	}
	if ssr.Code != "OK" {
		return false, errors.New(ssr.Code)
	}
	return true, nil
}

// specialURLEncode 处理特殊URL编码
func specialURLEncode(in string) string {
	// URLEncode后 加号（+）替换成 %20、星号（*）替换成 %2A、%7E 替换回波浪号（~）
	return strings.NewReplacer("+", "%20", "*", "%2A", "%7E", "~").Replace(url.QueryEscape(in))
}

// getSign 签名采用HmacSHA1算法 + Base64
func getSign(accessSecret, stringToSign string) (string, error) {
	mac := hmac.New(sha1.New, []byte(fmt.Sprintf("%s&", accessSecret)))
	_, err := mac.Write([]byte(stringToSign))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}
