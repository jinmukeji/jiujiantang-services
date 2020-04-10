package sms

// TemplateLanguage 模块语言
type TemplateLanguage int32

const (
	// SimpleChinese 简体中文
	SimpleChinese TemplateLanguage = 0
	// TraditionalChinese 繁体中文
	TraditionalChinese TemplateLanguage = 1
	// English 英文
	English TemplateLanguage = 2
)

// TemplateAction 模块行为
type TemplateAction int32

const (
	// SignUp 手机号注册
	SignUp TemplateAction = 0
	// SignIn 手机号登录
	SignIn TemplateAction = 1
	// ResetPassword 找回/重置密码
	ResetPassword TemplateAction = 2
	// SetPhoneNumber 设置手机号
	SetPhoneNumber TemplateAction = 3
	// ModifyPhoneNumber 修改手机号
	ModifyPhoneNumber TemplateAction = 4
)

// SMSClient 短信Client
type SMSClient interface {
	// SendSms 发送短信  phoneNumber 手机号 nationCode 国际代码 templateAction 模版行为 language 语言 templateParam 模版参数
	SendSms(phoneNumber, nationCode string, templateAction TemplateAction, language TemplateLanguage, templateParam map[string]string) (bool, error)
}
