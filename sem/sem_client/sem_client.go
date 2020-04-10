package sem

// TemplateLanguage 语言
type TemplateLanguage int32

// 邮件语言
const (
	// UndefinedLanguage 没有定义的语言
	UndefinedLanguage TemplateLanguage = 0
	// SimplifiedChinese 简体中文
	SimplifiedChinese TemplateLanguage = 1
	// TraditionalChinese 繁体中文
	TraditionalChinese TemplateLanguage = 2
	//  English 英语
	English TemplateLanguage = 3
)

// TemplateAction 模块行为
type TemplateAction int32

// 发送邮件的目的
const (
	// UndefinedAction 未定义行为
	UndefinedAction TemplateAction = 0
	// FindResetPassword 找回/重置密码
	FindResetPassword TemplateAction = 1
	// FindUsername 找回用户名
	FindUsername TemplateAction = 2
	// SetSecureEmail 设置安全邮箱
	SetSecureEmail TemplateAction = 3
	// ModifySecureEmail 修改安全邮箱
	ModifySecureEmail TemplateAction = 4
	// UnsetSecureEmail 解绑安全邮箱
	UnsetSecureEmail TemplateAction = 5
)

// SEMClient 邮件Client
type SEMClient interface {
	// SendEmail 发送触发邮件  toAddress 收件邮箱地址 templateAction 模版行为 language 语言 templateParam 模版参数
	SendEmail(toAddress string, templateAction TemplateAction, language TemplateLanguage, templateParam map[string]string) (bool, error)
}
