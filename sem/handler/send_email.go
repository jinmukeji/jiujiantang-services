package handler

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"fmt"

	mysqldb "github.com/jinmukeji/jiujiantang-services/sem/mysqldb"
	sem "github.com/jinmukeji/jiujiantang-services/sem/sem_client"
	generalpb "github.com/jinmukeji/proto/gen/micro/idl/ptypes/v2"
	sempb "github.com/jinmukeji/proto/gen/micro/idl/jm/sem/v1"
)

const (
	// AliyunPlatform 阿里云平台
	AliyunPlatform = "Aliyun"
	//NetEasePlatform 网易平台
	NetEasePlatform = "NetEase"
	// AliyunSeparator 阿里云分隔邮件的分隔符
	AliyunSeparator = ","
	// NetEaseSeparator 网易分隔邮件的分隔符
	NetEaseSeparator = ";"
	// StoreDBSeparator 存储到数据库里邮箱之间的分隔符
	StoreDBSeparator = " "
)

// SEMTemplateParam 用于解析proto传来的邮件模版参数
type SEMTemplateParam struct {
	Code string `json:"code"`
}

// EmailNotification 邮件通知
func (j *SEMGateway) EmailNotification(ctx context.Context, req *sempb.EmailNotificationRequest, resp *sempb.EmailNotificationResponse) error {
	return j.SendSEMEmail(req)
}

// SendSEMEmail 发送邮件
func (j *SEMGateway) SendSEMEmail(req *sempb.EmailNotificationRequest) error {
	// 优先使用阿里云发送邮件，如果不成功使用网易163
	isSuccess, _ := j.SendAliyunEmail(req)
	if !isSuccess {
		_, errSendSem := j.SendNetEaseEmail(req)
		if errSendSem != nil {
			return errSendSem
		}
	}
	return nil
}

// SendAliyunEmail 阿里云处理邮件
func (j *SEMGateway) SendAliyunEmail(req *sempb.EmailNotificationRequest) (bool, error) {
	var client sem.SEMClient = j.aliyunSEMClient
	action, errMapProtoTemplateActionToSem := mapProtoTemplateActionToSem(req.TemplateAction)
	if errMapProtoTemplateActionToSem != nil {
		return false, errMapProtoTemplateActionToSem
	}
	if action == sem.UndefinedAction {
		return false, errors.New("invalid template action")
	}
	clientLanguage, errMapProtoLanguageToSem := mapProtoLanguageToSem(req.Language)
	if errMapProtoLanguageToSem != nil {
		return false, errMapProtoLanguageToSem
	}
	if clientLanguage == sem.UndefinedLanguage {
		return false, errors.New("invalid language parameter")
	}
	// TODO: SendEmail,CreateSemRecord要在同一个事务
	// 这里先用阿里云发送邮件
	isSucceed, errSendSem := client.SendEmail(strings.Join(req.ToAddress, AliyunSeparator), action, clientLanguage, req.TemplateParam)
	// 生成记录
	now := time.Now().UTC()
	templateParam, errTemplateParam := json.Marshal(req.TemplateParam)
	if errTemplateParam != nil {
		return false, fmt.Errorf("failed to get marshal format %v: %s", req.TemplateParam, errTemplateParam.Error())
	}
	dbLanguage, errmapProtoLanguageToDB := mapProtoLanguageToDB(req.Language)
	if errmapProtoLanguageToDB != nil {
		return false, errmapProtoLanguageToDB
	}
	if dbLanguage == mysqldb.UndefinedLanguage {
		return false, errors.New("undefined language parameter")
	}
	record := &mysqldb.SemRecord{
		ToAddress:      strings.Join(req.ToAddress, StoreDBSeparator),
		TemplateAction: req.TemplateAction.String(),
		PlatformType:   AliyunPlatform,
		TemplateParam:  string(templateParam),
		Language:       dbLanguage,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if isSucceed {
		record.SemStatus = mysqldb.SendSucceed
	} else {
		record.SemStatus = mysqldb.SendFailed
	}
	if errSendSem != nil {
		record.SemErrorLog = errSendSem.Error()
	}
	errCreateSemRecord := j.datastore.CreateSemRecord(record)
	if errSendSem != nil {
		return false, fmt.Errorf("failed to create sem record of address %s: %s", strings.Join(req.ToAddress, StoreDBSeparator), errSendSem.Error())
	}
	return isSucceed, errCreateSemRecord
}

// SendNetEaseEmail 网易处理邮件逻辑
func (j *SEMGateway) SendNetEaseEmail(req *sempb.EmailNotificationRequest) (bool, error) {
	var client sem.SEMClient = j.neteaseSEMClient
	action, errMapProtoTemplateActionToSem := mapProtoTemplateActionToSem(req.TemplateAction)
	if errMapProtoTemplateActionToSem != nil {
		return false, errMapProtoTemplateActionToSem
	}
	if action == sem.UndefinedAction {
		return false, errors.New("undefined template action")
	}
	clientLanguage, errMapProtoLanguageToSem := mapProtoLanguageToSem(req.Language)
	if errMapProtoLanguageToSem != nil {
		return false, errMapProtoLanguageToSem
	}
	if clientLanguage == sem.UndefinedLanguage {
		return false, errors.New("undefined language parameter")
	}
	isSucceed, errSendSem := client.SendEmail(strings.Join(req.ToAddress, NetEaseSeparator), action, clientLanguage, req.TemplateParam)
	// 生成记录
	now := time.Now().UTC()
	templateParam, errTemplateParam := json.Marshal(req.TemplateParam)
	if errTemplateParam != nil {
		return false, fmt.Errorf("failed to get marshal format %v: %s", req.TemplateParam, errTemplateParam.Error())
	}
	dbLanguage, errmapProtoLanguageToDB := mapProtoLanguageToDB(req.Language)
	if errmapProtoLanguageToDB != nil {
		return false, errmapProtoLanguageToDB
	}
	if dbLanguage == mysqldb.UndefinedLanguage {
		return false, errors.New("undefined language parameter")
	}

	record := &mysqldb.SemRecord{
		ToAddress:      strings.Join(req.ToAddress, StoreDBSeparator),
		TemplateAction: req.TemplateAction.String(),
		PlatformType:   NetEasePlatform,
		TemplateParam:  string(templateParam),
		Language:       dbLanguage,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if isSucceed {
		record.SemStatus = mysqldb.SendSucceed
	} else {
		record.SemStatus = mysqldb.SendFailed
	}
	if errSendSem != nil {
		record.SemErrorLog = errSendSem.Error()
	}
	errCreateSemRecord := j.datastore.CreateSemRecord(record)
	if errCreateSemRecord != nil {
		return false, fmt.Errorf("failed to create sem record of address %s: %s", strings.Join(req.ToAddress, StoreDBSeparator), errCreateSemRecord.Error())
	}
	return isSucceed, errCreateSemRecord
}

func mapProtoLanguageToSem(language generalpb.Language) (sem.TemplateLanguage, error) {
	switch language {
	case generalpb.Language_LANGUAGE_INVALID:
		return sem.UndefinedLanguage, fmt.Errorf("invalid proto language %d", generalpb.Language_LANGUAGE_INVALID)
	case generalpb.Language_LANGUAGE_UNSET:
		return sem.UndefinedLanguage, fmt.Errorf("invalid proto language %d", generalpb.Language_LANGUAGE_UNSET)
	case generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE:
		return sem.SimplifiedChinese, nil
	case generalpb.Language_LANGUAGE_TRADITIONAL_CHINESE:
		return sem.TraditionalChinese, nil
	case generalpb.Language_LANGUAGE_ENGLISH:
		return sem.English, nil
	}
	return sem.UndefinedLanguage, fmt.Errorf("invalid proto language %d", generalpb.Language_LANGUAGE_INVALID)
}

// mapProtoTemplateActionToSem sempb.TemplateAction 转 sem.TemplateAction
func mapProtoTemplateActionToSem(templateAction sempb.TemplateAction) (sem.TemplateAction, error) {
	switch templateAction {
	case sempb.TemplateAction_TEMPLATE_ACTION_INVALID:
		return sem.UndefinedAction, fmt.Errorf("invalid sem template action %d", templateAction)
	case sempb.TemplateAction_TEMPLATE_ACTION_UNSET:
		return sem.UndefinedAction, fmt.Errorf("invalid sem template action %d", templateAction)
	case sempb.TemplateAction_TEMPLATE_ACTION_UNKNOWN:
		return sem.UndefinedAction, fmt.Errorf("invalid sem template action %d", templateAction)
	case sempb.TemplateAction_TEMPLATE_ACTION_FIND_RESET_PASSWORD:
		return sem.FindResetPassword, nil
	case sempb.TemplateAction_TEMPLATE_ACTION_FIND_USERNAME:
		return sem.FindUsername, nil
	case sempb.TemplateAction_TEMPLATE_ACTION_SET_SECURE_EMAIL:
		return sem.SetSecureEmail, nil
	case sempb.TemplateAction_TEMPLATE_ACTION_MODIFY_SECURE_EMAIL:
		return sem.ModifySecureEmail, nil
	case sempb.TemplateAction_TEMPLATE_ACTION_UNSET_SECURE_EMAIL:
		return sem.UnsetSecureEmail, nil
	}
	return sem.UndefinedAction, fmt.Errorf("invalid sem template action %d", templateAction)
}

func mapProtoLanguageToDB(language generalpb.Language) (mysqldb.Language, error) {
	switch language {
	case generalpb.Language_LANGUAGE_INVALID:
		return mysqldb.UndefinedLanguage, fmt.Errorf("invalid proto language %d", generalpb.Language_LANGUAGE_INVALID)
	case generalpb.Language_LANGUAGE_UNSET:
		return mysqldb.UndefinedLanguage, fmt.Errorf("invalid proto language %d", generalpb.Language_LANGUAGE_UNSET)
	case generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE:
		return mysqldb.SimpleChinese, nil
	case generalpb.Language_LANGUAGE_TRADITIONAL_CHINESE:
		return mysqldb.TraditionalChinese, nil
	case generalpb.Language_LANGUAGE_ENGLISH:
		return mysqldb.English, nil
	}
	return mysqldb.UndefinedLanguage, fmt.Errorf("invalid proto language %d", generalpb.Language_LANGUAGE_INVALID)
}
