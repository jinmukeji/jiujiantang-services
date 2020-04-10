package main

import (
	"os"
	"path"

	calcpb "github.com/jinmukeji/proto/gen/micro/idl/platform/calc/v2"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/config/cmd"
	"github.com/micro/go-micro/transport"

	ae "github.com/jinmukeji/ae-v1/core"
	"github.com/jinmukeji/jiujiantang-services/analysis/aws"
	logger "github.com/jinmukeji/jiujiantang-services/pkg/rpc"
	"github.com/jinmukeji/jiujiantang-services/service/config"
	handler "github.com/jinmukeji/jiujiantang-services/service/handler"
	"github.com/jinmukeji/jiujiantang-services/service/mail"
	"github.com/jinmukeji/jiujiantang-services/service/mysqldb"
	"github.com/jinmukeji/jiujiantang-services/service/wechat"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
)

func main() {
	authenticationWrapper := new(handler.AuthenticationWrapper)

	versionMeta := config.NewVersionMetadata()

	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		// Service Basic Info
		micro.Name(config.FullServiceName()),
		micro.Version(config.ProductVersion),

		// Fault Tolerance - Heartbeating
		// 	 See also: https://micro.mu/docs/fault-tolerance.html#heartbeating
		micro.RegisterTTL(config.DefaultRegisterTTL),
		micro.RegisterInterval(config.DefaultRegisterInterval),

		// Setup wrappers
		micro.WrapHandler(logger.LogWrapper, authenticationWrapper.AuthWrapper(), authenticationWrapper.HandleWrapper()),

		// Setup runtime flags
		dbClientOptions(), mailClientOptions(), algorithmClientOptions(), awsClientOptions(), jinmuHealthOptions(), wechatOptions(),

		// Setup --version flag
		micro.Flags(
			cli.BoolFlag{
				Name:  "version",
				Usage: "Show version information",
			},
		),

		// Setup metadata
		micro.Metadata(versionMeta),
	)

	// optionally setup command line usage
	service.Init(
		micro.Action(func(c *cli.Context) {
			if c.Bool("version") {
				config.PrintFullVersionInfo()
				os.Exit(0)
			}
		}),
	)

	log.Infof("Starting service: %s", config.FullServiceName())
	log.Infof("Product Version: %s", config.ProductVersion)
	log.Infof("Git SHA: %s", config.GitSHA)
	log.Infof("Git Branch: %s", config.GitBranch)
	log.Infof("Go Version: %s", config.GoVersion)
	log.Infof("Build Version: %s", config.BuildVersion)
	log.Infof("Build Time: %s", config.BuildTime)

	// Register handler
	server := service.Server()

	// Connect to MySQL instance
	db, err := newDbClient()
	mailClient := newMailClient()
	algorithmClient := newAlgorithmClient()
	if err != nil {
		log.Panicf("Failed to connect to MySQL instance at %s. Error: %v", dbAddress, err)
	}
	log.Infoln("Connected to MySQL instance at", dbAddress)
	awsClient, err := newAWSClient()
	if err != nil {
		log.Panicf("Failed to connect to aws s3 server at %s. Error: %v", awsBucketName, err)
	}

	if err := broker.Init(); err != nil {
		log.Fatalf("Broker Init error: %v", err)
	}

	if err := broker.Connect(); err != nil {
		log.Fatalf("Broker Connect error: %v", err)
	}
	ae := newAE()
	wx := newWechat()

	jinmuHealth := handler.NewJinmuHealth(db, mailClient, algorithmClient, awsClient, ae, wx, algorithmServerAddress)
	if err := proto.RegisterJinmuhealthAPIHandler(server, jinmuHealth); err != nil {
		log.Fatalln(err)
	}
	// 设置认证 handlerWrapper 的数据库
	authenticationWrapper.SetDataStore(db)
	if errInit := cmd.Init(); errInit != nil {
		log.Fatalln(errInit)
	}

	// Run the server
	if err := service.Run(); err != nil {
		log.Fatalln(err)
	}
}

// 数据库和邮件服务器连接信息
var (
	dbAddress          string
	dbUsername         string
	dbPassword         string
	dbDatabase         string
	dbEnableLog        = false
	dbMaxConns         = 1
	mailAddress        string
	mailUsername       string
	mailPassword       string
	mailPort           int
	mailCharset        string
	mailSenderNickname string
	mailReply          string
)

// 算法服务器连接信息
var (
	algorithmServerAddress string
)

// aws 存储桶连接信息
var (
	awsBucketName                          string
	awsAccessKey                           string
	awsSecretKey                           string
	awsRegion                              string
	PulseTestRawDataEnvironmentS3KeyPrefix string
	pulseTestRawDataS3KeyPrefix            string
)

// AE 信息
var (
	aeConfigDir string
)

// 微信信息
var (
	wxAppID               string
	wxAppSecret           string // 微信公众平台 开发者密码
	wxOriID               string // 微信公众平台 公众号原始 ID
	wxToken               string // 微信公众平台 令牌 (Token)
	wxEncodedAESKey       string // 微信公众平台 消息加解密密钥 (EncodingAESKey)
	wxCallbackServerBase  string // 微信公众平台 回调服务器的地址
	wxTemplateID          string // 微信模版ID
	wxH5ServerBase        string // 微信h5地址
	jinmuH5ServerBase     string // jinmu的h5地址
	jinmuH5ServerbaseV2_0 string // v2 jinmu的h5地址
	jinmuH5ServerbaseV2_1 string // v2.1 jinmu的h5地址
)

// ip和mac过滤的信息
var (
	blockerDBFile     string // ip配置的数据库文件
	blockerConfigFile string // ip和mac过滤的配置文件
)

// jinmuHealthOptions 构建命令行启动参数
func jinmuHealthOptions() micro.Option {
	return micro.Flags(
		cli.StringFlag{
			Name:        "x_ae_config_dir",
			Value:       "",
			Usage:       "Analysis Engine configuration directory",
			EnvVar:      "X_AE_CONFIG_DIR",
			Destination: &aeConfigDir,
		},
	)
}

// dbClientOptions 构建命令行启动参数
func dbClientOptions() micro.Option {
	return micro.Flags(
		cli.StringFlag{
			Name:        "x_db_address",
			Value:       "localhost:3306",
			Usage:       "MySQL instance `ADDRESS` - [host]:[port]",
			EnvVar:      "X_DB_ADDRESS",
			Destination: &dbAddress,
		},
		cli.StringFlag{
			Name:        "x_db_username",
			Usage:       "MySQL login `USERNAME`",
			EnvVar:      "X_DB_USERNAME",
			Destination: &dbUsername,
		},
		cli.StringFlag{
			Name:        "x_db_password",
			Usage:       "MySQL login `PASSWORD`",
			EnvVar:      "X_DB_PASSWORD",
			Destination: &dbPassword,
		},
		cli.StringFlag{
			Name:        "x_db_database",
			Usage:       "MySQL database name",
			EnvVar:      "X_DB_DATABASE",
			Destination: &dbDatabase,
		},
		cli.BoolFlag{
			Name:        "x_db_enable_log",
			Usage:       "Enable MySQL client log",
			EnvVar:      "X_DB_ENABLE_LOG",
			Destination: &dbEnableLog,
		},
		cli.IntFlag{
			Name:        "x_db_max_connections",
			Usage:       "Max connections of MySQL client",
			EnvVar:      "X_DB_MAX_CONNECTIONS",
			Value:       1,
			Destination: &dbMaxConns,
		},
	)
}

// mailClientOptions 构建命令行启动参数
func mailClientOptions() micro.Option {
	return micro.Flags(
		cli.StringFlag{
			Name:        "x_mail_address",
			Usage:       "mail server address without port",
			EnvVar:      "X_MAIL_ADDRESS",
			Destination: &mailAddress,
		},
		cli.StringFlag{
			Name:        "x_mail_password",
			Usage:       "password of mail account",
			EnvVar:      "X_MAIL_PASSWORD",
			Destination: &mailPassword,
		},
		cli.IntFlag{
			Name:        "x_mail_port",
			Usage:       "listening port of mail server",
			EnvVar:      "X_MAIL_PORT",
			Destination: &mailPort,
		},
		cli.StringFlag{
			Name:        "x_mail_charset",
			Usage:       "charset of mail content",
			EnvVar:      "X_MAIL_CHARSET",
			Destination: &mailCharset,
		},
		cli.StringFlag{
			Name:        "x_mail_username",
			Usage:       "username of mail account in mail server",
			EnvVar:      "X_MAIL_USERNAME",
			Destination: &mailUsername,
		},
		cli.StringFlag{
			Name:        "x_mail_sender_nickname",
			Usage:       "nickname of mail account in mail server",
			EnvVar:      "X_MAIL_SENDER_NICKNAME",
			Destination: &mailSenderNickname,
		},
		cli.StringFlag{
			Name:        "x_mail_reply",
			Usage:       "mail address for reply",
			EnvVar:      "X_MAIL_REPLY",
			Destination: &mailReply,
		},
	)
}

// algorithmClientOptions 构建命令行启动参数
func algorithmClientOptions() micro.Option {
	return micro.Flags(
		cli.StringFlag{
			Name:        "x_algorithm_server_address",
			Usage:       "algorithm server address",
			EnvVar:      "X_ALGORITHM_SERVER_ADDRESS",
			Destination: &algorithmServerAddress,
		},
	)
}

// awsClientOptions 构建命令行启动参数
func awsClientOptions() micro.Option {
	return micro.Flags(
		cli.StringFlag{
			Name:        "x_aws_bucket_name",
			Usage:       "aws bucket name",
			EnvVar:      "X_AWS_BUCKET_NAME",
			Destination: &awsBucketName,
		},
		cli.StringFlag{
			Name:        "x_aws_access_key",
			Usage:       "aws access key",
			EnvVar:      "X_AWS_ACCESS_KEY",
			Destination: &awsAccessKey,
		},
		cli.StringFlag{
			Name:        "x_aws_secret_key",
			Usage:       "aws secret key",
			EnvVar:      "X_AWS_SECRET_KEY",
			Destination: &awsSecretKey,
		},
		cli.StringFlag{
			Name:        "x_aws_region",
			Usage:       "aws region",
			EnvVar:      "X_AWS_REGION",
			Destination: &awsRegion,
		},
		cli.StringFlag{
			Name:        "x_wave_data_key_prefix",
			Usage:       "S3 prefix key for wave data",
			EnvVar:      "X_WAVE_DATA_KEY_PREFIX",
			Destination: &PulseTestRawDataEnvironmentS3KeyPrefix,
		},
		cli.StringFlag{
			Name:        "x_pulse_test_raw_data_s3_key_prefix",
			Usage:       "pulse test raw data s3 key prefix",
			EnvVar:      "X_PULSE_TEST_RAW_DATA_S3_KEY_PREFIX",
			Destination: &pulseTestRawDataS3KeyPrefix,
		},
	)
}

func wechatOptions() micro.Option {
	return micro.Flags(

		cli.StringFlag{
			Name:        "x_wx_app_id",
			Value:       "",
			Usage:       "微信公众平台 开发者 ID",
			EnvVar:      "X_WX_APP_ID",
			Destination: &wxAppID,
		},

		cli.StringFlag{
			Name:        "x_wx_app_secret",
			Value:       "",
			Usage:       "微信公众平台 开发者密码",
			EnvVar:      "X_WX_APP_SECRET",
			Destination: &wxAppSecret,
		},

		cli.StringFlag{
			Name:        "x_wx_ori_id",
			Value:       "",
			Usage:       "微信公众平台 公众号原始 ID",
			EnvVar:      "X_WX_ORI_ID",
			Destination: &wxOriID,
		},

		cli.StringFlag{
			Name:        "x_wx_token",
			Value:       "",
			Usage:       "微信公众平台 令牌",
			EnvVar:      "X_WX_TOKEN",
			Destination: &wxToken,
		},

		cli.StringFlag{
			Name:        "x_wx_encoded_aes_key",
			Value:       "",
			Usage:       "微信公众平台 消息加解密密钥",
			EnvVar:      "X_WX_ENCODED_AES_KEY",
			Destination: &wxEncodedAESKey,
		},

		cli.StringFlag{
			Name:        "x_wx_callback_server_base",
			Value:       "",
			Usage:       "微信公众平台 回调服务器的地址",
			EnvVar:      "X_WX_CALLBACK_SERVER_BASE",
			Destination: &wxCallbackServerBase,
		},
		cli.StringFlag{
			Name:        "x_wx_template_id",
			Value:       "",
			Usage:       "微信公众平台 模版ID",
			EnvVar:      "X_WX_Template_ID",
			Destination: &wxTemplateID,
		},
		cli.StringFlag{
			Name:        "x_wx_h5_server_base",
			Value:       "",
			Usage:       "微信公众平台 H5地址",
			EnvVar:      "X_WX_H5_SERVER_BASE",
			Destination: &wxH5ServerBase,
		},
		cli.StringFlag{
			Name:        "x_jinmu_h5_server_base",
			Value:       "",
			Usage:       "金姆 H5地址",
			EnvVar:      "X_JINMU_H5_SERVER_BASE",
			Destination: &jinmuH5ServerBase,
		},
		cli.StringFlag{
			Name:        "x_jinmu_h5_server_base_v2_0",
			Value:       "",
			Usage:       "金姆 v2 H5地址",
			EnvVar:      "X_JINMU_H5_SERVER_BASE_V2_0",
			Destination: &jinmuH5ServerbaseV2_0,
		},
		cli.StringFlag{
			Name:        "x_jinmu_h5_server_base_v2_1",
			Value:       "",
			Usage:       "金姆 v2 H5地址",
			EnvVar:      "X_JINMU_H5_SERVER_BASE_V2_1",
			Destination: &jinmuH5ServerbaseV2_1,
		},
	)
}

// newDbClient 创建一个 DbClient
func newDbClient() (*mysqldb.DbClient, error) {
	return mysqldb.NewDbClient(
		mysqldb.Address(dbAddress),
		mysqldb.Username(dbUsername),
		mysqldb.Password(dbPassword),
		mysqldb.Database(dbDatabase),
		mysqldb.EnableLog(dbEnableLog),
		mysqldb.MaxConnections(dbMaxConns),
	)
}

// newMailClient 创建一个 mailCLient
func newMailClient() *mail.Client {
	return mail.NewMailClient(
		mail.Address(mailAddress),
		mail.Username(mailUsername),
		mail.Password(mailPassword),
		mail.Port(mailPort),
		mail.Charset(mailCharset),
		mail.SenderNickname(mailSenderNickname),
	)
}

const (
	AlgorithmFQDN = ""
)

// newAlgorithmClient 创建一个算法服务器连接
func newAlgorithmClient() calcpb.CalcAPIService {

	client := client.NewClient(
		client.Transport(transport.NewTransport(transport.Secure(true))),
	)
	cl := calcpb.NewCalcAPIService(AlgorithmFQDN, client)

	return cl
}

// newAWSClient 创建一个 aws 连接
func newAWSClient() (*aws.Client, error) {
	return aws.NewClient(
		aws.BucketName(awsBucketName),
		aws.Region(awsRegion),
		aws.AccessKeyID(awsAccessKey),
		aws.SecretKey(awsSecretKey),
		aws.PulseTestRawDataEnvironmentS3KeyPrefix(PulseTestRawDataEnvironmentS3KeyPrefix),
		aws.PulseTestRawDataS3KeyPrefix(pulseTestRawDataS3KeyPrefix),
	)
}

func newAE() *ae.Engine {
	eg := ae.NewEngine()
	eg.LoadTemplatesFromDir(path.Join(aeConfigDir, "lookups"))
	eg.LoadRuleSetDoc(path.Join(aeConfigDir, "rules/data-rules-latest.yaml"))

	return eg
}

func newWechat() *wechat.Wxmp {
	return wechat.NewWxmp(
		broker.DefaultBroker,
		&wechat.WxmpOptions{
			WxAppID:               wxAppID,
			WxAppSecret:           wxAppSecret,
			WxOriID:               wxOriID,
			WxToken:               wxToken,
			WxEncodedAESKey:       wxEncodedAESKey,
			WxCallbackServerBase:  wxCallbackServerBase,
			WxTemplateID:          wxTemplateID,
			WxH5ServerBase:        wxH5ServerBase,
			JinmuH5Serverbase:     jinmuH5ServerBase,
			JinmuH5ServerbaseV2_0: jinmuH5ServerbaseV2_0,
			JinmuH5ServerbaseV2_1: jinmuH5ServerbaseV2_1,
		})
}
