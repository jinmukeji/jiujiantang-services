package handler

import (
	"context"
	"fmt"
	"path/filepath"

	jinmuidpb "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	corepb "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	sempb "github.com/jinmukeji/proto/gen/micro/idl/jm/sem/v1"
	smspb "github.com/jinmukeji/proto/gen/micro/idl/jm/sms/v1"
	subscriptionpb "github.com/jinmukeji/proto/gen/micro/idl/jm/subscription/v1"
	generalpb "github.com/jinmukeji/proto/gen/micro/idl/ptypes/v2"
	"github.com/micro/go-micro/client"
)

const (
	rpcSmsServiceName       = "com.himalife.srv.svc-sms-gw"
	rpcSemServiceName       = "com.himalife.srv.svc-sem-gw"
	rpcServiceName          = "com.himalife.srv.svc-jinmuid"
	rpcBizServiceName       = "com.himalife.srv.svc-biz-core"
	subscriptionServiceName = "com.himalife.srv.svc-subscription"
)

// 初始化
func newJinmuIDServiceForTest() *JinmuIDService {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	db, err := newTestingDbClientFromEnvFile(envFilepath)
	if err != nil {
		panic(fmt.Sprintln("failed to init db:", err))
	}
	encryptKey := newTestingEncryptKeyFromEnvFile(envFilepath)
	smsSvc := smspb.NewSmsAPIService(rpcSmsServiceName, client.DefaultClient)
	semSvc := sempb.NewSemAPIService(rpcSemServiceName, client.DefaultClient)
	rpcUserManagerSvc := jinmuidpb.NewUserManagerAPIService(rpcServiceName, client.DefaultClient)
	bizSvc := corepb.NewJinmuhealthAPIService(rpcBizServiceName, client.DefaultClient)
	subscriptionSvc := subscriptionpb.NewSubscriptionManagerAPIService(subscriptionServiceName, client.DefaultClient)
	return NewJinmuIDService(db, smsSvc, semSvc, rpcUserManagerSvc, bizSvc, subscriptionSvc, encryptKey)
}

// 短信注册
func getSignUpSerialNumber(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	resp := new(jinmuidpb.SmsNotificationResponse)
	req := new(jinmuidpb.SmsNotificationRequest)
	req.Phone = account.phone
	req.Action = jinmuidpb.TemplateAction_TEMPLATE_ACTION_SIGN_UP
	req.Language = generalpb.Language_LANGUAGE_ENGLISH
	req.NationCode = account.nationCode
	_ = jinmuIDService.SmsNotification(ctx, req, resp)
	return resp.SerialNumber
}

// 短信登录
func getSignInSerialNumber(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	resp := new(jinmuidpb.SmsNotificationResponse)
	req := new(jinmuidpb.SmsNotificationRequest)
	req.Phone = account.phone
	req.Action = jinmuidpb.TemplateAction_TEMPLATE_ACTION_SIGN_IN
	req.Language = generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE
	req.NationCode = account.nationCode
	_ = jinmuIDService.SmsNotification(ctx, req, resp)
	return resp.SerialNumber
}

// 获取短信验证码
func getSmsVerificationCode(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	respGetLatestVerificationCodes := new(jinmuidpb.GetLatestVerificationCodesResponse)
	reqGetLatestVerificationCodes := new(jinmuidpb.GetLatestVerificationCodesRequest)
	reqGetLatestVerificationCodes.SendTo = []*jinmuidpb.SingleGetLatestVerificationCode{
		&jinmuidpb.SingleGetLatestVerificationCode{
			SendVia:    jinmuidpb.SendVia_SEND_VIA_PHONE_SEND_VIA,
			Phone:      account.phone,
			NationCode: account.nationCode,
		},
	}
	_ = jinmuIDService.GetLatestVerificationCodes(ctx, reqGetLatestVerificationCodes, respGetLatestVerificationCodes)
	mvc := respGetLatestVerificationCodes.LatestVerificationCodes[0].VerificationCode
	return mvc
}

// 获取邮件验证码
func getEmailVerificationCode(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	respGetLatestVerificationCodes := new(jinmuidpb.GetLatestVerificationCodesResponse)
	reqGetLatestVerificationCodes := new(jinmuidpb.GetLatestVerificationCodesRequest)
	reqGetLatestVerificationCodes.SendTo = []*jinmuidpb.SingleGetLatestVerificationCode{
		&jinmuidpb.SingleGetLatestVerificationCode{
			SendVia: jinmuidpb.SendVia_SEND_VIA_USERNAME_SEND_VIA,
			Email:   account.email,
		},
	}
	_ = jinmuIDService.GetLatestVerificationCodes(ctx, reqGetLatestVerificationCodes, respGetLatestVerificationCodes)
	mvc := respGetLatestVerificationCodes.LatestVerificationCodes[0].VerificationCode
	return mvc
}

// 获取新邮件验证码
func getNewEmailVerificationCode(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	respGetLatestVerificationCodes := new(jinmuidpb.GetLatestVerificationCodesResponse)
	reqGetLatestVerificationCodes := new(jinmuidpb.GetLatestVerificationCodesRequest)
	reqGetLatestVerificationCodes.SendTo = []*jinmuidpb.SingleGetLatestVerificationCode{
		&jinmuidpb.SingleGetLatestVerificationCode{
			SendVia: jinmuidpb.SendVia_SEND_VIA_USERNAME_SEND_VIA,
			Email:   account.emailNew,
		},
	}
	_ = jinmuIDService.GetLatestVerificationCodes(ctx, reqGetLatestVerificationCodes, respGetLatestVerificationCodes)
	mvc := respGetLatestVerificationCodes.LatestVerificationCodes[0].VerificationCode
	return mvc
}

// GetSetEmailSerialNumber 设置邮箱获取serialNumber
func getSetEmailSerialNumber(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	ctx, userID, _ := mockSigninByPhonePassword(ctx, jinmuIDService, account.phone, account.phonePassword, account.seed, account.nationCode)
	resp := new(jinmuidpb.LoggedInEmailNotificationResponse)
	req := new(jinmuidpb.LoggedInEmailNotificationRequest)
	req.Email = account.email
	req.Action = jinmuidpb.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_SET_SECURE_EMAIL
	req.Language = generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE
	req.UserId = userID
	req.SendToNewIfModify = false
	_ = jinmuIDService.LoggedInEmailNotification(ctx, req, resp)
	return resp.SerialNumber
}

// GetUnSetEmailSerialNumber  解绑邮箱获取serialNumber
func getUnSetEmailSerialNumber(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	ctx, userID, _ := mockSigninByPhonePassword(ctx, jinmuIDService, account.phone, account.phonePassword, account.seed, account.nationCode)
	resp := new(jinmuidpb.LoggedInEmailNotificationResponse)
	req := new(jinmuidpb.LoggedInEmailNotificationRequest)
	req.Email = account.email
	req.Action = jinmuidpb.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_UNSET_SECURE_EMAIL
	req.Language = generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE
	req.UserId = userID
	req.SendToNewIfModify = false
	_ = jinmuIDService.LoggedInEmailNotification(ctx, req, resp)
	return resp.SerialNumber
}

// GetModifyEmailSerialNumber 修改邮箱获取serialNumber
func getModifyEmailSerialNumber(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	ctx, userID, _ := mockSigninByPhonePassword(ctx, jinmuIDService, account.phone, account.phonePassword, account.seed, account.nationCode)
	resp := new(jinmuidpb.LoggedInEmailNotificationResponse)
	req := new(jinmuidpb.LoggedInEmailNotificationRequest)
	req.Email = account.email
	req.Action = jinmuidpb.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_MODIFY_SECURE_EMAIL
	req.Language = generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE
	req.UserId = userID
	req.SendToNewIfModify = false
	_ = jinmuIDService.LoggedInEmailNotification(ctx, req, resp)
	return resp.SerialNumber
}

// GetModifyEmailSerialNumberNewEmail 修改邮箱获取serialNumber
func getModifyEmailSerialNumberNewEmail(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	ctx, userID, _ := mockSigninByPhonePassword(ctx, jinmuIDService, account.phone, account.phonePassword, account.seed, account.nationCode)
	resp := new(jinmuidpb.LoggedInEmailNotificationResponse)
	req := new(jinmuidpb.LoggedInEmailNotificationRequest)
	req.Email = account.emailNew
	req.Action = jinmuidpb.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_MODIFY_SECURE_EMAIL
	req.Language = generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE
	req.UserId = userID
	req.SendToNewIfModify = true
	_ = jinmuIDService.LoggedInEmailNotification(ctx, req, resp)
	return resp.SerialNumber
}

// GetFindUserNameSerialNumber 通过邮箱找用户名serialNumber
func getFindUserNameSerialNumber(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	resp := new(jinmuidpb.NotLoggedInEmailNotificationResponse)
	req := new(jinmuidpb.NotLoggedInEmailNotificationRequest)
	req.Email = account.email
	req.Action = jinmuidpb.NotLoggedInSemTemplateAction_NOT_LOGGED_IN_SEM_TEMPLATE_ACTION_FIND_USERNAME
	req.Language = generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE
	_ = jinmuIDService.NotLoggedInEmailNotification(ctx, req, resp)
	return resp.SerialNumber
}

// GetVerificationNumber 获取VerificationNumber
func getVerificationNumber(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	// 发送通知
	serialNumber := getModifyEmailSerialNumber(jinmuIDService, account)
	// 获取最新验证码
	mvc := getEmailVerificationCode(jinmuIDService, account)
	resp := new(jinmuidpb.ValidateEmailVerificationCodeResponse)
	req := new(jinmuidpb.ValidateEmailVerificationCodeRequest)
	req.VerificationCode = mvc
	req.SerialNumber = serialNumber
	req.Email = account.email
	req.VerificationType = account.verificationType
	_ = jinmuIDService.ValidateEmailVerificationCode(ctx, req, resp)
	return resp.VerificationNumber
}

func getSignUpVerificationNumber(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	mvc := getSmsVerificationCode(jinmuIDService, account)
	serialNumber := getSignUpSerialNumber(jinmuIDService, account)
	resp := new(jinmuidpb.ValidatePhoneVerificationCodeResponse)
	req := new(jinmuidpb.ValidatePhoneVerificationCodeRequest)
	req.Phone = account.phone
	req.Mvc = mvc
	req.SerialNumber = serialNumber
	req.NationCode = account.nationCode
	_ = jinmuIDService.ValidatePhoneVerificationCode(ctx, req, resp)
	return resp.VerificationNumber
}

// 手机号获取短信验证码-香港手机号码
func getSignUpVerificationNumberHK(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	mvc := getSmsVerificationCodeHK(jinmuIDService, account)
	serialNumber := getSignUpSerialNumberHK(jinmuIDService, account)
	resp := new(jinmuidpb.ValidatePhoneVerificationCodeResponse)
	req := new(jinmuidpb.ValidatePhoneVerificationCodeRequest)
	req.Phone = account.phoneHK
	req.Mvc = mvc
	req.SerialNumber = serialNumber
	req.NationCode = account.nationCodeHK
	_ = jinmuIDService.ValidatePhoneVerificationCode(ctx, req, resp)

	return resp.VerificationNumber
}

// 获取短信验证码-香港手机号码
func getSmsVerificationCodeHK(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	respGetLatestVerificationCodes := new(jinmuidpb.GetLatestVerificationCodesResponse)
	reqGetLatestVerificationCodes := new(jinmuidpb.GetLatestVerificationCodesRequest)
	reqGetLatestVerificationCodes.SendTo = []*jinmuidpb.SingleGetLatestVerificationCode{
		&jinmuidpb.SingleGetLatestVerificationCode{
			SendVia:    jinmuidpb.SendVia_SEND_VIA_PHONE_SEND_VIA,
			Phone:      account.phoneHK,
			NationCode: account.nationCodeHK,
		},
	}
	_ = jinmuIDService.GetLatestVerificationCodes(ctx, reqGetLatestVerificationCodes, respGetLatestVerificationCodes)
	mvc := respGetLatestVerificationCodes.LatestVerificationCodes[0].VerificationCode
	return mvc
}

// 手机验证短信注册验证码是否正确-香港手机号码
func getSignUpSerialNumberHK(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	resp := new(jinmuidpb.SmsNotificationResponse)
	req := new(jinmuidpb.SmsNotificationRequest)
	req.Phone = account.phoneHK
	req.Action = jinmuidpb.TemplateAction_TEMPLATE_ACTION_SIGN_UP
	req.Language = generalpb.Language_LANGUAGE_ENGLISH
	req.NationCode = account.nationCodeHK
	_ = jinmuIDService.SmsNotification(ctx, req, resp)
	return resp.SerialNumber
}

// 手机号获取短信验证码-美国手机号码
func getSignUpVerificationNumberUS(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	mvc := getSmsVerificationCodeUS(jinmuIDService, account)
	serialNumber := getSignUpSerialNumberUS(jinmuIDService, account)
	resp := new(jinmuidpb.ValidatePhoneVerificationCodeResponse)
	req := new(jinmuidpb.ValidatePhoneVerificationCodeRequest)
	req.Phone = account.phoneUS
	req.Mvc = mvc
	req.SerialNumber = serialNumber
	req.NationCode = account.nationCodeUSA
	_ = jinmuIDService.ValidatePhoneVerificationCode(ctx, req, resp)
	return resp.VerificationNumber
}

// 获取短信验证码-美国手机号码
func getSmsVerificationCodeUS(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	respGetLatestVerificationCodes := new(jinmuidpb.GetLatestVerificationCodesResponse)
	reqGetLatestVerificationCodes := new(jinmuidpb.GetLatestVerificationCodesRequest)
	reqGetLatestVerificationCodes.SendTo = []*jinmuidpb.SingleGetLatestVerificationCode{
		&jinmuidpb.SingleGetLatestVerificationCode{
			SendVia:    jinmuidpb.SendVia_SEND_VIA_PHONE_SEND_VIA,
			Phone:      account.phoneUS,
			NationCode: account.nationCodeUSA,
		},
	}
	_ = jinmuIDService.GetLatestVerificationCodes(ctx, reqGetLatestVerificationCodes, respGetLatestVerificationCodes)
	mvc := respGetLatestVerificationCodes.LatestVerificationCodes[0].VerificationCode
	return mvc
}

// 手机验证短信注册验证码是否正确-美国手机号码
func getSignUpSerialNumberUS(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	resp := new(jinmuidpb.SmsNotificationResponse)
	req := new(jinmuidpb.SmsNotificationRequest)
	req.Phone = account.phoneUS
	req.Action = jinmuidpb.TemplateAction_TEMPLATE_ACTION_SIGN_UP
	req.Language = generalpb.Language_LANGUAGE_ENGLISH
	req.NationCode = account.nationCodeUSA
	_ = jinmuIDService.SmsNotification(ctx, req, resp)
	return resp.SerialNumber
}

// 手机号获取短信验证码-台湾手机号码
func getSignUpVerificationNumberTW(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	mvc := getSmsVerificationCodeTW(jinmuIDService, account)
	serialNumber := getSignUpSerialNumberTW(jinmuIDService, account)
	resp := new(jinmuidpb.ValidatePhoneVerificationCodeResponse)
	req := new(jinmuidpb.ValidatePhoneVerificationCodeRequest)
	req.Phone = account.phoneTW
	req.Mvc = mvc
	req.SerialNumber = serialNumber
	req.NationCode = account.nationCodeTW
	_ = jinmuIDService.ValidatePhoneVerificationCode(ctx, req, resp)
	return resp.VerificationNumber
}

// 获取短信验证码-台湾手机号码
func getSmsVerificationCodeTW(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	respGetLatestVerificationCodes := new(jinmuidpb.GetLatestVerificationCodesResponse)
	reqGetLatestVerificationCodes := new(jinmuidpb.GetLatestVerificationCodesRequest)
	reqGetLatestVerificationCodes.SendTo = []*jinmuidpb.SingleGetLatestVerificationCode{
		&jinmuidpb.SingleGetLatestVerificationCode{
			SendVia:    jinmuidpb.SendVia_SEND_VIA_PHONE_SEND_VIA,
			Phone:      account.phoneTW,
			NationCode: account.nationCodeTW,
		},
	}
	_ = jinmuIDService.GetLatestVerificationCodes(ctx, reqGetLatestVerificationCodes, respGetLatestVerificationCodes)
	mvc := respGetLatestVerificationCodes.LatestVerificationCodes[0].VerificationCode
	return mvc
}

// 手机验证短信注册验证码是否正确-台湾手机号码
func getSignUpSerialNumberTW(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	resp := new(jinmuidpb.SmsNotificationResponse)
	req := new(jinmuidpb.SmsNotificationRequest)
	req.Phone = account.phoneTW
	req.Action = jinmuidpb.TemplateAction_TEMPLATE_ACTION_SIGN_UP
	req.Language = generalpb.Language_LANGUAGE_ENGLISH
	req.NationCode = account.nationCodeTW
	_ = jinmuIDService.SmsNotification(ctx, req, resp)
	return resp.SerialNumber
}

// 手机号获取短信验证码-澳门手机号码
func getSignUpVerificationNumberMacao(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	mvc := getSmsVerificationCodeMacao(jinmuIDService, account)
	serialNumber := getSignUpSerialNumberMacao(jinmuIDService, account)
	resp := new(jinmuidpb.ValidatePhoneVerificationCodeResponse)
	req := new(jinmuidpb.ValidatePhoneVerificationCodeRequest)
	req.Phone = account.phoneMacao
	req.Mvc = mvc
	req.SerialNumber = serialNumber
	req.NationCode = account.nationCodeMacao
	_ = jinmuIDService.ValidatePhoneVerificationCode(ctx, req, resp)
	return resp.VerificationNumber
}

// 获取短信验证码-澳门手机号码
func getSmsVerificationCodeMacao(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	respGetLatestVerificationCodes := new(jinmuidpb.GetLatestVerificationCodesResponse)
	reqGetLatestVerificationCodes := new(jinmuidpb.GetLatestVerificationCodesRequest)
	reqGetLatestVerificationCodes.SendTo = []*jinmuidpb.SingleGetLatestVerificationCode{
		&jinmuidpb.SingleGetLatestVerificationCode{
			SendVia:    jinmuidpb.SendVia_SEND_VIA_PHONE_SEND_VIA,
			Phone:      account.phoneMacao,
			NationCode: account.nationCodeMacao,
		},
	}
	_ = jinmuIDService.GetLatestVerificationCodes(ctx, reqGetLatestVerificationCodes, respGetLatestVerificationCodes)
	mvc := respGetLatestVerificationCodes.LatestVerificationCodes[0].VerificationCode
	return mvc
}

// 手机验证短信注册验证码是否正确-澳门手机号码
func getSignUpSerialNumberMacao(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	resp := new(jinmuidpb.SmsNotificationResponse)
	req := new(jinmuidpb.SmsNotificationRequest)
	req.Phone = account.phoneMacao
	req.Action = jinmuidpb.TemplateAction_TEMPLATE_ACTION_SIGN_UP
	req.Language = generalpb.Language_LANGUAGE_ENGLISH
	req.NationCode = account.nationCodeMacao
	_ = jinmuIDService.SmsNotification(ctx, req, resp)
	return resp.SerialNumber
}

// 手机号获取短信验证码-加拿大手机号码
func getSignUpVerificationNumberCanada(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	mvc := getSmsVerificationCodeCanada(jinmuIDService, account)
	serialNumber := getSignUpSerialNumberCanada(jinmuIDService, account)
	resp := new(jinmuidpb.ValidatePhoneVerificationCodeResponse)
	req := new(jinmuidpb.ValidatePhoneVerificationCodeRequest)
	req.Phone = account.phoneCanada
	req.Mvc = mvc
	req.SerialNumber = serialNumber
	req.NationCode = account.nationCodeUSA
	_ = jinmuIDService.ValidatePhoneVerificationCode(ctx, req, resp)
	return resp.VerificationNumber
}

// 获取短信验证码-加拿大手机号码
func getSmsVerificationCodeCanada(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	respGetLatestVerificationCodes := new(jinmuidpb.GetLatestVerificationCodesResponse)
	reqGetLatestVerificationCodes := new(jinmuidpb.GetLatestVerificationCodesRequest)
	reqGetLatestVerificationCodes.SendTo = []*jinmuidpb.SingleGetLatestVerificationCode{
		&jinmuidpb.SingleGetLatestVerificationCode{
			SendVia:    jinmuidpb.SendVia_SEND_VIA_PHONE_SEND_VIA,
			Phone:      account.phoneCanada,
			NationCode: account.nationCodeUSA,
		},
	}
	_ = jinmuIDService.GetLatestVerificationCodes(ctx, reqGetLatestVerificationCodes, respGetLatestVerificationCodes)
	mvc := respGetLatestVerificationCodes.LatestVerificationCodes[0].VerificationCode
	return mvc
}

// 手机验证短信注册验证码是否正确-加拿大手机号码
func getSignUpSerialNumberCanada(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	resp := new(jinmuidpb.SmsNotificationResponse)
	req := new(jinmuidpb.SmsNotificationRequest)
	req.Phone = account.phoneCanada
	req.Action = jinmuidpb.TemplateAction_TEMPLATE_ACTION_SIGN_UP
	req.Language = generalpb.Language_LANGUAGE_ENGLISH
	req.NationCode = account.nationCodeUSA
	_ = jinmuIDService.SmsNotification(ctx, req, resp)
	return resp.SerialNumber
}

// 手机号获取短信验证码-英国手机号码
func getSignUpVerificationNumberUK(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	mvc := getSmsVerificationCodeUK(jinmuIDService, account)
	serialNumber := getSignUpSerialNumberUK(jinmuIDService, account)
	resp := new(jinmuidpb.ValidatePhoneVerificationCodeResponse)
	req := new(jinmuidpb.ValidatePhoneVerificationCodeRequest)
	req.Phone = account.phoneUK
	req.Mvc = mvc
	req.SerialNumber = serialNumber
	req.NationCode = account.nationCodeUK
	_ = jinmuIDService.ValidatePhoneVerificationCode(ctx, req, resp)
	return resp.VerificationNumber
}

// 获取短信验证码-英国手机号码
func getSmsVerificationCodeUK(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	respGetLatestVerificationCodes := new(jinmuidpb.GetLatestVerificationCodesResponse)
	reqGetLatestVerificationCodes := new(jinmuidpb.GetLatestVerificationCodesRequest)
	reqGetLatestVerificationCodes.SendTo = []*jinmuidpb.SingleGetLatestVerificationCode{
		&jinmuidpb.SingleGetLatestVerificationCode{
			SendVia:    jinmuidpb.SendVia_SEND_VIA_PHONE_SEND_VIA,
			Phone:      account.phoneUK,
			NationCode: account.nationCodeUK,
		},
	}
	_ = jinmuIDService.GetLatestVerificationCodes(ctx, reqGetLatestVerificationCodes, respGetLatestVerificationCodes)
	mvc := respGetLatestVerificationCodes.LatestVerificationCodes[0].VerificationCode
	return mvc
}

// 手机验证短信注册验证码是否正确-英国手机号码
func getSignUpSerialNumberUK(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	resp := new(jinmuidpb.SmsNotificationResponse)
	req := new(jinmuidpb.SmsNotificationRequest)
	req.Phone = account.phoneUK
	req.Action = jinmuidpb.TemplateAction_TEMPLATE_ACTION_SIGN_UP
	req.Language = generalpb.Language_LANGUAGE_ENGLISH
	req.NationCode = account.nationCodeUK
	_ = jinmuIDService.SmsNotification(ctx, req, resp)
	return resp.SerialNumber
}

// 手机号获取短信验证码-日本手机号码
func getSignUpVerificationNumberJP(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	mvc := getSmsVerificationCodeJP(jinmuIDService, account)
	serialNumber := getSignUpSerialNumberJP(jinmuIDService, account)
	resp := new(jinmuidpb.ValidatePhoneVerificationCodeResponse)
	req := new(jinmuidpb.ValidatePhoneVerificationCodeRequest)
	req.Phone = account.phoneJP
	req.Mvc = mvc
	req.SerialNumber = serialNumber
	req.NationCode = account.nationCodeJP
	_ = jinmuIDService.ValidatePhoneVerificationCode(ctx, req, resp)
	return resp.VerificationNumber
}

// 获取短信验证码-日本手机号码
func getSmsVerificationCodeJP(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	respGetLatestVerificationCodes := new(jinmuidpb.GetLatestVerificationCodesResponse)
	reqGetLatestVerificationCodes := new(jinmuidpb.GetLatestVerificationCodesRequest)
	reqGetLatestVerificationCodes.SendTo = []*jinmuidpb.SingleGetLatestVerificationCode{
		&jinmuidpb.SingleGetLatestVerificationCode{
			SendVia:    jinmuidpb.SendVia_SEND_VIA_PHONE_SEND_VIA,
			Phone:      account.phoneJP,
			NationCode: account.nationCodeJP,
		},
	}
	_ = jinmuIDService.GetLatestVerificationCodes(ctx, reqGetLatestVerificationCodes, respGetLatestVerificationCodes)
	mvc := respGetLatestVerificationCodes.LatestVerificationCodes[0].VerificationCode
	return mvc
}

// 手机验证短信注册验证码是否正确-日本手机号码
func getSignUpSerialNumberJP(jinmuIDService *JinmuIDService, account Account) string {
	ctx := context.Background()
	resp := new(jinmuidpb.SmsNotificationResponse)
	req := new(jinmuidpb.SmsNotificationRequest)
	req.Phone = account.phoneJP
	req.Action = jinmuidpb.TemplateAction_TEMPLATE_ACTION_SIGN_UP
	req.Language = generalpb.Language_LANGUAGE_ENGLISH
	req.NationCode = account.nationCodeJP
	_ = jinmuIDService.SmsNotification(ctx, req, resp)
	return resp.SerialNumber
}
