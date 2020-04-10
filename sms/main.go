package main

import (
	"os"

	logger "github.com/jinmukeji/gf-api2/pkg/rpc"
	"github.com/jinmukeji/gf-api2/sms/config"
	"github.com/jinmukeji/gf-api2/sms/mysqldb"
	sms "github.com/jinmukeji/gf-api2/sms/sms_client"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/sms/v1"

	handler "github.com/jinmukeji/gf-api2/sms/handler"
	"github.com/micro/cli"
	micro "github.com/micro/go-micro"
)

func main() {
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
		micro.WrapHandler(logger.LogWrapper),

		// // Setup runtime flags
		dbClientOptions(), aliyunSmsOptions(), tencentYunSmsOptions(),

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
	db, err := newDbClient()
	if err != nil {
		log.Panicf("Failed to connect to MySQL instance at %s. Error: %v", dbAddress, err)
	}
	log.Infoln("Connected to MySQL instance at", dbAddress)

	aliyunSmsClient, errNewAliyunSmsClient := newAliyunSmsClient()
	if errNewAliyunSmsClient != nil {
		log.Fatalln(errNewAliyunSmsClient)
	}
	tencentYunSmsClient, errTencentYunSmsClient := newTencentYunSmsClient()
	if errTencentYunSmsClient != nil {
		log.Fatalln(errTencentYunSmsClient)
	}
	smsGateway := handler.NewSMSGateway(db, aliyunSmsClient, tencentYunSmsClient)
	if err := proto.RegisterSmsAPIHandler(server, smsGateway); err != nil {
		log.Fatalln(err)
	}
	// Run the server
	if err := service.Run(); err != nil {
		log.Fatalln(err)
	}
}

// 数据库和邮件服务器连接信息
var (
	dbAddress   string
	dbUsername  string
	dbPassword  string
	dbDatabase  string
	dbEnableLog = false
	dbMaxConns  = 1
)

// 阿里云短信连接信息
var (
	aliyunSmsAccessKeyID     string
	aliyunSmsAccessKeySecret string
)

var (
	tencentYunSmsAccessKeyID     string
	tencentYunSmsAccessKeySecret string
)

// aliyunSmsOptions 构建命令行启动参数
func aliyunSmsOptions() micro.Option {
	return micro.Flags(
		cli.StringFlag{
			Name:        "x_aliyun_sms_access_key_id",
			Usage:       "Aliyun SMS Access Key ID",
			EnvVar:      "X_ALIYUN_SMS_ACCESS_KEY_ID",
			Destination: &aliyunSmsAccessKeyID,
		},
		cli.StringFlag{
			Name:        "x_aliyun_sms_access_key_secret",
			Usage:       "Aliyun SMS Access Key Secret",
			EnvVar:      "X_ALIYUN_SMS_ACCESS_KEY_Secret",
			Destination: &aliyunSmsAccessKeySecret,
		},
	)
}

// tencentYunSmsOptions 构建命令行启动参数
func tencentYunSmsOptions() micro.Option {
	return micro.Flags(
		cli.StringFlag{
			Name:        "x_tencent_yun_sms_access_key_id",
			Usage:       "Tencent Yun SMS Access Key ID",
			EnvVar:      "X_TENCENT_YUN_SMS_ACCESS_APP_ID",
			Destination: &tencentYunSmsAccessKeyID,
		},
		cli.StringFlag{
			Name:        "x_tencent_yun_sms_access_key_secret",
			Usage:       "Tencent Yun Access Key Secret",
			EnvVar:      "X_TENCENT_YUN_SMS_ACCESS_KEY_Secret",
			Destination: &tencentYunSmsAccessKeySecret,
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

// newAliyunSmsClient 创建阿里云短信网关客户端
func newAliyunSmsClient() (*sms.AliyunSMSClient, error) {
	return sms.NewAliyunSMSClient(aliyunSmsAccessKeyID, aliyunSmsAccessKeySecret)
}

// newTencentYunSmsClient 创建腾讯云短信网关客户端
func newTencentYunSmsClient() (*sms.TencentYunSMSClient, error) {
	return sms.NewTencentYunSMSClient(tencentYunSmsAccessKeyID, tencentYunSmsAccessKeySecret)
}
