package mail

var (
	// textHtml 是 html 文档
	textHTML = "text/html"
)

// Options 是  邮件服务器 的配置参数
type Options struct {
	// 服务器地址 - host
	Address string

	// 用户名
	Username string

	// 密码
	Password string

	// 字符集
	Charset string

	// 端口号
	Port int

	// 发件人的昵称
	SenderNickname string

	// 客户回复的邮件服务器地址
	ReplyToAddress string
}

// Option 是 设置 Options 的函数
type Option func(options *Options)

// Address 设置服务器地址
func Address(address string) Option {
	return func(options *Options) {
		options.Address = address
	}
}

// Username 设置邮件服务器用户名
func Username(username string) Option {
	return func(options *Options) {
		options.Username = username
	}
}

// Password 设置邮件服务器密码
func Password(password string) Option {
	return func(options *Options) {
		options.Password = password
	}
}

// Charset 设置邮件编码
func Charset(charset string) Option {
	return func(options *Options) {
		options.Charset = charset
	}
}

// Port 设置邮件编码
func Port(port int) Option {
	return func(options *Options) {
		options.Port = port
	}
}

// SenderNickname 设置发件人的昵称
func SenderNickname(nickname string) Option {
	return func(options *Options) {
		options.SenderNickname = nickname
	}
}

// ReplyToAddress 设置回复地址
func ReplyToAddress(reply string) Option {
	return func(options *Options) {
		options.ReplyToAddress = reply
	}
}

// defaultOptions 返回默认参数
func defaultOptions() *Options {
	return &Options{
		Address:  "localhost",
		Username: "root",
		Password: "",
		Charset:  "UTF-8",
		Port:     25,
	}
}

// newOptions 设置 Options
func newOptions(opt ...Option) *Options {
	opts := defaultOptions()

	for _, o := range opt {
		o(opts)
	}

	return opts
}
