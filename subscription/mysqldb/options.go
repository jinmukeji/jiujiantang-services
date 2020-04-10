package mysqldb

// Options 是 DbClient 的配置参数
type Options struct {
	// 是否启用日志
	EnableLog bool

	// 是否转换时间
	ParseTime bool

	// 最大连接数
	MaxConnections int

	// 服务器地址 - host:port
	Address string

	// 用户名
	Username string

	// 密码
	Password string

	// 数据库名
	Database string

	// 字符集
	Charset string

	// 区域设置
	Locale string
}

// Option 是设置 Options 的函数
type Option func(*Options)

// Address 设置服务器地址 - host:port
func Address(addr string) Option {
	return func(o *Options) {
		o.Address = addr
	}
}

// Username 设置用户名
func Username(username string) Option {
	return func(o *Options) {
		o.Username = username
	}
}

// Password 设置密码
func Password(pwd string) Option {
	return func(o *Options) {
		o.Password = pwd
	}
}

// Database 设置数据库名
func Database(db string) Option {
	return func(o *Options) {
		o.Database = db
	}
}

// EnableLog 设置是否启用日志
func EnableLog(enable bool) Option {
	return func(o *Options) {
		o.EnableLog = enable
	}
}

// MaxConnections 设置最大连接数
func MaxConnections(maxConns int) Option {
	return func(o *Options) {
		o.MaxConnections = maxConns
	}
}

// Charset 设置字符集
func Charset(charset string) Option {
	return func(o *Options) {
		o.Charset = charset
	}
}

// ParseTime 设置转换时间
func ParseTime(parseTime bool) Option {
	return func(o *Options) {
		o.ParseTime = parseTime
	}
}

// Locale 设置区域设置
func Locale(locale string) Option {
	return func(o *Options) {
		o.Locale = locale
	}
}

// defaultOptions 返回默认配置的 Options
func defaultOptions() Options {
	return Options{
		Address:        "localhost:3306",
		EnableLog:      false,
		MaxConnections: 1,
		Charset:        "utf8mb4",
		ParseTime:      true,
		Locale:         "UTC", // 注意: 这里字母必须大写，否则找不到 Timezone 文件
	}
}

// newOptions 设置 Options
func newOptions(opt ...Option) Options {
	opts := defaultOptions()

	for _, o := range opt {
		o(&opts)
	}

	return opts
}
