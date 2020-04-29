package handler

import (
	"context"
	"encoding/json"
	"time"

	"fmt"

	mysqldb "github.com/jinmukeji/jiujiantang-services/sms/mysqldb"
	sms "github.com/jinmukeji/jiujiantang-services/sms/sms_client"
	smspb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/sms/v1"
	generalpb "github.com/jinmukeji/proto/v3/gen/micro/idl/ptypes/v2"
)

const (
	// limitsIn24Hours 单人发送短信24小时内上限（10条）
	limitsIn24Hours = 10
	// limitsIn1Minutes 单人 一分钟内上限(1条)
	limitsIn1Minutes = 60
	// AliyunPlatform 阿里云平台
	AliyunPlatform = "Aliyun"
)

// SMSTemplateParam 用于解析proto传来的短信模版参数
type SMSTemplateParam struct {
	Code string `json:"code"`
}

// SendMessage 发送短信
func (j *SMSGateway) SendMessage(ctx context.Context, req *smspb.SendMessageRequest, resp *smspb.SendMessageResponse) error {
	counts, err := j.datastore.SearchSmsRecordCountsIn1Minute(req.Phone, req.NationCode)
	if err != nil {
		return fmt.Errorf("failed to search sms record counts of %s%s in 1 minute: %s", req.NationCode, req.Phone, err.Error())
	}
	// 1分钟限制
	if counts > limitsIn1Minutes && !req.IsForced {
		return fmt.Errorf("send message of %s%s reach 1 minute limits", req.NationCode, req.Phone)
	}
	countsHours, err := j.datastore.SearchSmsRecordCountsIn24hours(req.Phone, req.NationCode)
	if err != nil {
		return fmt.Errorf("failed to search sms record counts of %s%s in 24 hours: %s", req.NationCode, req.Phone, err.Error())
	}
	// 24小时限制
	if countsHours > limitsIn24Hours && !req.IsForced {
		return fmt.Errorf("send message of %s%s reach 24 hours limits", req.NationCode, req.Phone)
	}
	return j.SendSMSMessage(req, resp)
}

// SendSMSMessage 发送短信
func (j *SMSGateway) SendSMSMessage(req *smspb.SendMessageRequest, resp *smspb.SendMessageResponse) error {
	_, err := j.SendAliyunMessage(req, resp)
	return err
}

// SendAliyunMessage 阿里云处理短信逻辑
func (j *SMSGateway) SendAliyunMessage(req *smspb.SendMessageRequest, resp *smspb.SendMessageResponse) (bool, error) {
	var client sms.SMSClient = j.aliyunSMSClient
	// TODO: SendSms,CreateSmsRecord 要在同一个事务
	// 这里先用阿里云发送短信
	smsLanguage, errMapProtoLanguageToSms := mapProtoLanguageToSms(req.Language)
	if errMapProtoLanguageToSms != nil {
		return false, errMapProtoLanguageToSms
	}
	templateAction, errMapProtoTemplateActionToSms := mapProtoTemplateActionToSms(req.TemplateAction)
	if errMapProtoTemplateActionToSms != nil {
		return false, errMapProtoTemplateActionToSms
	}
	isSucceed, errSendSms := client.SendSms(req.Phone, req.NationCode, templateAction, smsLanguage, req.TemplateParam)

	// 生成记录
	now := time.Now().UTC()
	templateParam, errTemplateParam := json.Marshal(req.TemplateParam)
	if errTemplateParam != nil {
		return false, fmt.Errorf("failed to get marshal format %v: %s", req.TemplateParam, errTemplateParam.Error())
	}
	language, errmapProtoLanguageToDB := mapProtoLanguageToDB(req.Language)
	if errmapProtoLanguageToDB != nil {
		return false, errmapProtoLanguageToDB
	}
	record := &mysqldb.SmsRecord{
		Phone:          req.Phone,
		TemplateAction: req.TemplateAction.String(),
		PlatformType:   AliyunPlatform,
		TemplateParam:  string(templateParam),
		NationCode:     req.NationCode,
		Language:       language,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if isSucceed {
		record.SmsStatus = mysqldb.SendSucceed
	} else {
		record.SmsStatus = mysqldb.SendFailed
	}
	if errSendSms != nil {
		record.SmsErrorLog = errSendSms.Error()
	}
	errCreateSmsRecord := j.datastore.CreateSmsRecord(record)
	if errSendSms != nil {
		return false, fmt.Errorf("failed to create sem record of phone %s%s: %s", req.NationCode, req.Phone, errSendSms.Error())
	}
	return isSucceed, errCreateSmsRecord
}

func mapProtoLanguageToSms(language generalpb.Language) (sms.TemplateLanguage, error) {
	switch language {
	case generalpb.Language_LANGUAGE_INVALID:
		return sms.SimpleChinese, fmt.Errorf("invalid proto language %d", language)
	case generalpb.Language_LANGUAGE_UNSET:
		return sms.SimpleChinese, fmt.Errorf("invalid proto language %d", language)
	case generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE:
		return sms.SimpleChinese, nil
	case generalpb.Language_LANGUAGE_TRADITIONAL_CHINESE:
		return sms.TraditionalChinese, nil
	case generalpb.Language_LANGUAGE_ENGLISH:
		return sms.English, nil

	}
	return sms.SimpleChinese, fmt.Errorf("invalid proto language %d", language)
}

func mapProtoTemplateActionToSms(templateAction smspb.TemplateAction) (sms.TemplateAction, error) {
	switch templateAction {
	case smspb.TemplateAction_TEMPLATE_ACTION_INVALID:
		return sms.SignUp, fmt.Errorf("invalid proto template action %d", templateAction)
	case smspb.TemplateAction_TEMPLATE_ACTION_UNSET:
		return sms.SignIn, fmt.Errorf("invalid proto template action %d", templateAction)
	case smspb.TemplateAction_TEMPLATE_ACTION_SIGN_UP:
		return sms.SignUp, nil
	case smspb.TemplateAction_TEMPLATE_ACTION_SIGN_IN:
		return sms.SignIn, nil
	case smspb.TemplateAction_TEMPLATE_ACTION_RESET_PASSWORD:
		return sms.ResetPassword, nil
	case smspb.TemplateAction_TEMPLATE_ACTION_SET_PHONE_NUMBER:
		return sms.SetPhoneNumber, nil
	case smspb.TemplateAction_TEMPLATE_ACTION_MODIFY_PHONE_NUMBER:
		return sms.ModifyPhoneNumber, nil
	}
	return sms.SignUp, fmt.Errorf("invalid proto template action %d", templateAction)
}

func mapProtoLanguageToDB(language generalpb.Language) (mysqldb.Language, error) {
	switch language {
	case generalpb.Language_LANGUAGE_INVALID:
		return mysqldb.SimpleChinese, fmt.Errorf("invalid proto language %d", generalpb.Language_LANGUAGE_INVALID)
	case generalpb.Language_LANGUAGE_UNSET:
		return mysqldb.SimpleChinese, fmt.Errorf("invalid proto language %d", generalpb.Language_LANGUAGE_UNSET)
	case generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE:
		return mysqldb.SimpleChinese, nil
	case generalpb.Language_LANGUAGE_TRADITIONAL_CHINESE:
		return mysqldb.TraditionalChinese, nil
	case generalpb.Language_LANGUAGE_ENGLISH:
		return mysqldb.English, nil
	}
	return mysqldb.SimpleChinese, fmt.Errorf("invalid proto language %d", generalpb.Language_LANGUAGE_INVALID)
}
