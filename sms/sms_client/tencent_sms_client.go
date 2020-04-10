package sms

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	// SingleBaseURL 腾讯云单发短信基本URL
	SingleBaseURL = "https://yun.tim.qq.com/v5/tlssmssvr/sendsms?sdkappid=%s&random=%s"
)

// simpleChineseTemplateActionTX 国内简体模版
var simpleChineseTemplateActionTX = map[TemplateAction]int{
	SignUp:            271955, // 手机号注册
	SignIn:            271956, // 手机号登录
	ResetPassword:     271957, // 找回/重置密码
	SetPhoneNumber:    271959, // 设置手机号
	ModifyPhoneNumber: 271960, // 修改手机号
}

// simpleChineseInternationalTemplateActionTX 国际简体模版
var simpleChineseInternationalTemplateActionTX = map[TemplateAction]int{
	SignUp:            271955, // 手机号注册
	SignIn:            271956, // 手机号登录
	ResetPassword:     271957, // 找回/重置密码
	SetPhoneNumber:    271959, // 设置手机号
	ModifyPhoneNumber: 271960, // 修改手机号
}

// traditionalChineseTemplateActionTX 国内繁体模版
var traditionalChineseTemplateActionTX = map[TemplateAction]int{
	SignUp:            271961, // 手机号注册
	SignIn:            271962, // 手机号登录
	ResetPassword:     271963, // 找回/重置密码
	SetPhoneNumber:    271965, // 设置手机号
	ModifyPhoneNumber: 271966, // 修改手机号
}

// traditionalChineseInternationalTemplateActionTX 国际繁体模版
var traditionalChineseInternationalTemplateActionTX = map[TemplateAction]int{
	SignUp:            271961, // 手机号注册
	SignIn:            271962, // 手机号登录
	ResetPassword:     271963, // 找回/重置密码
	SetPhoneNumber:    271965, // 设置手机号
	ModifyPhoneNumber: 271966, // 修改手机号
}

// englishTemplateActionTX 国内英文模版
var englishTemplateActionTX = map[TemplateAction]int{
	SignUp:            271970, // 手机号注册
	SignIn:            271971, // 手机号登录
	ResetPassword:     271973, // 找回/重置密码
	SetPhoneNumber:    271977, // 设置手机号
	ModifyPhoneNumber: 271980, // 修改手机号
}

// englishInternationalTemplateActionTX 国际英文模版
var englishInternationalTemplateActionTX = map[TemplateAction]int{
	SignUp:            271970, // 手机号注册
	SignIn:            271971, // 手机号登录
	ResetPassword:     271973, // 找回/重置密码
	SetPhoneNumber:    271977, // 设置手机号
	ModifyPhoneNumber: 271980, // 修改手机号
}

// TencentYunSMSClient 腾讯SMS客户端
type TencentYunSMSClient struct {
	SDKAppID string
	AppKey   string
}

// SendSingleSmsResponse 单发短信响应
type SendSingleSmsResponse struct {
	Result int32  `json:"result"`
	Errmsg string `json:"errmsg"`
	Ext    string `json:"ext"`
	Fee    int32  `json:"fee"`
	Sid    string `json:"sid"`
}

// Tel 电话及国家信息
type Tel struct {
	Mobile       string `json:"mobile"`
	Nationalcode string `json:"nationcode"`
}

// SendSingleSmsRequestBody 单发短信请求
type SendSingleSmsRequestBody struct {
	Exit   string   `json:"ext"`
	Extend string   `json:"extend"`
	Params []string `json:"params"`
	Sig    string   `json:"sig"`
	Sign   string   `json:"sign"`
	Tel    Tel      `json:"tel"`
	Time   int64    `json:"time"`
	TplID  int      `json:"tpl_id"`
}

// NewTencentYunSMSClient 生成腾讯SMSClient
func NewTencentYunSMSClient(SDKAppID string, AppKey string) (*TencentYunSMSClient, error) {
	if SDKAppID == "" {
		return nil, errors.New("SDKAppID should be not empty")
	}
	if AppKey == "" {
		return nil, errors.New("AppKey should be not empty")
	}
	return &TencentYunSMSClient{
		SDKAppID: SDKAppID,
		AppKey:   AppKey,
	}, nil
}

// getsha256 得到通过sha256加密生成的签名,腾讯官方文档要求
func getsha256(appkey string, random string, now int64, phoneNumbers string) string {
	proclaimed := fmt.Sprintf("appkey=%s&random=%s&time=%d&mobile=%s", appkey, random, now, phoneNumbers)
	h := sha256.New()
	_, _ = h.Write([]byte(proclaimed))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// SendSms 单发短信
// 参考链接 https://cloud.tencent.com/document/product/382/5976
func (client *TencentYunSMSClient) SendSms(phoneNumber, nationCode string, templateAction TemplateAction, language TemplateLanguage, templateParam map[string]string) (bool, error) {
	// 转化成阿里云短信模版
	id := convertTemplateActionTX(nationCode, templateAction, language)
	now := time.Now().UTC().Unix()
	Sign := "金姆平台"
	if language == English {
		Sign = "Jinmu"
	}
	sendSmsReqBody := &SendSingleSmsRequestBody{
		Params: []string{templateParam["code"]},
		Tel: Tel{
			Mobile:       phoneNumber,
			Nationalcode: dealNationCode(nationCode),
		},
		Time:  now,
		TplID: id,
		Sign:  Sign,
	}
	random := uuid.New().String()
	sendSmsReqBody.Sig = getsha256(client.AppKey, random, now, phoneNumber)
	reqBytes, err := json.Marshal(sendSmsReqBody)
	if err != nil {
		return false, fmt.Errorf("failed to marshal tencent request body data. %v", err)
	}
	// URL 示例
	// POST https://yun.tim.qq.com/v5/tlssmssvr/sendsms?sdkappid=xxxxx&random=xxxx
	// dkappid是在腾讯云上申请到的sdkappid，random是生成的随机数
	url := fmt.Sprintf(SingleBaseURL, url.QueryEscape(client.SDKAppID), url.QueryEscape(random))
	return sendTencentYunSmsAPI(url, reqBytes)
}

// sendTencentYunSmsAPI 发送腾讯云短信API
func sendTencentYunSmsAPI(url string, reqBytes []byte) (bool, error) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	ssr := &SendSingleSmsResponse{}
	if err := json.Unmarshal(body, ssr); err != nil {
		return false, err
	}
	if ssr.Result != 0 {
		return false, errors.New(ssr.Errmsg)
	}
	return true, nil
}

func dealNationCode(nationCode string) string {
	if nationCode == "" {
		return "86"
	}
	if nationCode != "" && strings.HasPrefix(nationCode, "+") {
		return nationCode[1:]
	}
	return nationCode
}

func convertTemplateActionTX(nationCode string, templateAction TemplateAction, language TemplateLanguage) int {
	if nationCode == "" || nationCode == "+86" {
		switch language {
		case SimpleChinese:
			return simpleChineseTemplateActionTX[templateAction]
		case TraditionalChinese:
			return traditionalChineseTemplateActionTX[templateAction]
		case English:
			return englishTemplateActionTX[templateAction]
		}
		return simpleChineseTemplateActionTX[templateAction]
	}
	switch language {
	case SimpleChinese:
		return simpleChineseInternationalTemplateActionTX[templateAction]
	case TraditionalChinese:
		return traditionalChineseInternationalTemplateActionTX[templateAction]
	case English:
		return englishInternationalTemplateActionTX[templateAction]
	}
	return simpleChineseInternationalTemplateActionTX[templateAction]
}

// MultiBaseURL 腾讯云群发短信基本URL
const MultiBaseURL = "https://yun.tim.qq.com/v5/tlssmssvr/sendmultisms2?sdkappid=%s&random=%s"

// SendMultiSmsRequest 群发短信的请求
type SendMultiSmsRequest struct {
	Exit   string   `json:"ext"`
	Extend string   `json:"extend"`
	Params []string `json:"params"`
	Sig    string   `json:"sig"`
	Sign   string   `json:"sign"`
	Tel    []Tel    `json:"tel"`
	Time   int64    `json:"time"`
	TplID  int      `json:"tpl_id"`
}

// MultiResponseDetail 群发短信的单个响应
type MultiResponseDetail struct {
	Errmsg     string `json:"errmsg"`
	Fee        int32  `json:"fee"`
	Mobile     string `json:"string"`
	Nationcode string `json:"nationcode"`
	Result     int32  `json:"result"`
	Sid        string `json:"sid"`
}

// SendMultiSmsResponse 群发短信的响应
type SendMultiSmsResponse struct {
	Result int32                 `json:"result"`
	Errmsg string                `json:"errmsg"`
	Ext    string                `json:"ext"`
	Detail []MultiResponseDetail `json:"detail"`
}

// SendMultiSms 群发短信
func (client *TencentYunSMSClient) SendMultiSms(phoneNumbers []string, nationalcodes []string, templateAction TemplateAction, language TemplateLanguage, templateParam map[string]string) (bool, error) {
	now := time.Now().Unix()
	id := convertTemplateActionTX(nationalcodes[0], templateAction, language) // 国内模板和国外模板都一样

	if len(phoneNumbers) != len(nationalcodes) {
		return false, fmt.Errorf("the length of phoneNumbers and nationcodes are not equal")
	}

	detailPhones := make([]Tel, len(phoneNumbers))
	for i := range phoneNumbers {
		detailPhones[i].Mobile = phoneNumbers[i]
		detailPhones[i].Nationalcode = dealNationCode(nationalcodes[i])
	}
	sendSmsReq := &SendMultiSmsRequest{
		Params: []string{templateParam["code"]},
		Tel:    detailPhones,
		Time:   now,
		TplID:  id,
	}

	random := uuid.New().String()
	sendSmsReq.Sig = getsha256(client.AppKey, random, now, strings.Join(phoneNumbers, ","))
	reqBytes, err := json.Marshal(sendSmsReq)
	if err != nil {
		return false, fmt.Errorf("failed to marshal MeasurementResult data. %v", err)
	}

	// URL 示例
	// POST https://yun.tim.qq.com/v5/tlssmssvr/sendmultisms2?sdkappid=xxxxx&random=xxxx
	// dkappid是在腾讯云上申请到的sdkappid，random是生成的随机数
	url := fmt.Sprintf(MultiBaseURL, url.QueryEscape(client.SDKAppID), url.QueryEscape(random))
	return sendTencentYunMultiSmsAPI(url, reqBytes)
}

// sendTencentYunMultiSmsAPI 发送腾讯云短信API
func sendTencentYunMultiSmsAPI(url string, reqBytes []byte) (bool, error) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	ssr := &SendMultiSmsResponse{}
	if err := json.Unmarshal(body, ssr); err != nil {
		return false, err
	}
	if ssr.Result != 0 {
		return false, errors.New(ssr.Errmsg)
	}
	return true, nil
}
