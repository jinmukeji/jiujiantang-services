package aws

// Options 是 aws 连接配置
type Options struct {
	// aws 存储桶名称
	BucketName string
	// aws api 访问身份认证
	AccessKeyID string
	// aws api 访问密钥
	SecretKey string
	// aws 所在区域
	Region string
	// 波形数据存储 key 的前缀
	PulseTestRawDataEnvironmentS3KeyPrefix string
	// 当前的数据存储 key 的前缀
	PulseTestRawDataS3KeyPrefix string
}

// Option 是 aws 连接配置设置方法
type Option func(opt *Options)

// BucketName 设置 aws 存储桶地址
func BucketName(name string) Option {
	return func(options *Options) {
		options.BucketName = name
	}
}

// AccessKeyID 设置 aws api 访问身份认证信息
func AccessKeyID(id string) Option {
	return func(options *Options) {
		options.AccessKeyID = id
	}
}

// SecretKey 设置 aws api 访问密钥
func SecretKey(key string) Option {
	return func(options *Options) {
		options.SecretKey = key
	}
}

// Region 设置 aws 区域
func Region(region string) Option {
	return func(options *Options) {
		options.Region = region
	}
}

// PulseTestRawDataEnvironmentS3KeyPrefix 设置存储前缀
func PulseTestRawDataEnvironmentS3KeyPrefix(prefix string) Option {
	return func(options *Options) {
		options.PulseTestRawDataEnvironmentS3KeyPrefix = prefix
	}
}

// PulseTestRawDataS3KeyPrefix 设置当前存储前缀
func PulseTestRawDataS3KeyPrefix(prefix string) Option {
	return func(options *Options) {
		options.PulseTestRawDataS3KeyPrefix = prefix
	}
}

// defaultOptions 返回 aws 默认连接配置
func defaultOptions() *Options {
	return &Options{
		BucketName:  "bucket-name",
		AccessKeyID: "0xff",
		SecretKey:   "0xff",
	}
}

// newOptions 返回新的 aws 连接配置
func newOptions(opts ...Option) *Options {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return options
}
