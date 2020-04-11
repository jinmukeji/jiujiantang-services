package main

import (
	"os"

	logger "github.com/jinmukeji/jiujiantang-services/pkg/rpc"
	"github.com/jinmukeji/jiujiantang-services/subscription/config"

	handler "github.com/jinmukeji/jiujiantang-services/subscription/handler"
	"github.com/jinmukeji/jiujiantang-services/subscription/mysqldb"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/subscription/v1"
	"github.com/micro/cli"
	micro "github.com/micro/go-micro"
)

func main() {
	versionMeta := config.NewVersionMetadata()
	// Create a new service. Optionally include some options here.
	authenticationWrapper := new(handler.AuthenticationWrapper)
	service := micro.NewService(
		// Service Basic Info
		micro.Name(config.FullServiceName()),
		micro.Version(config.ProductVersion),

		// Fault Tolerance - Heartbeating
		// 	 See also: https://micro.mu/docs/fault-tolerance.html#heartbeating
		micro.RegisterTTL(config.DefaultRegisterTTL),
		micro.RegisterInterval(config.DefaultRegisterInterval),

		// Setup wrappers
		micro.WrapHandler(logger.LogWrapper, authenticationWrapper.HandleWrapper()),

		// // Setup runtime flags
		dbClientOptions(), activationCodeEntryptKeyOptions(),

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
	authenticationWrapper.SetDataStore(db)
	subscriptionService := handler.NewSubscriptionService(db, activationCodeEntryptKey)
	if err := proto.RegisterSubscriptionManagerAPIHandler(server, subscriptionService); err != nil {
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

var activationCodeEntryptKey string

func activationCodeEntryptKeyOptions() micro.Option {
	return micro.Flags(
		cli.StringFlag{
			Name:        "x_activation_code_encrypt_key",
			Value:       "",
			Usage:       "激活码加密的key",
			EnvVar:      "X_ACTIVATION_CODE_ENCRYPT_KEY",
			Destination: &activationCodeEntryptKey,
		},
	)
}
