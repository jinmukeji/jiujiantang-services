package main

import (
	"os"

	"github.com/jinmukeji/ae/v2/biz"
	"github.com/jinmukeji/jiujiantang-services/analysis/aws"
	"github.com/jinmukeji/jiujiantang-services/analysis/config"
	handler "github.com/jinmukeji/jiujiantang-services/analysis/handler"
	"github.com/jinmukeji/jiujiantang-services/analysis/mysqldb"
	logger "github.com/jinmukeji/jiujiantang-services/pkg/rpc"
	dbutilmysql "github.com/jinmukeji/plat-pkg/v2/dbutil/mysql"
	storemysql "github.com/jinmukeji/plat-pkg/v2/store/mysql"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/analysis/v1"
	"github.com/micro/cli/v2"
	micro "github.com/micro/go-micro/v2"
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
		dbClientOptions(), aeOptions(), awsClientOptions(),

		// Setup --version flag
		defaultVersionFlags(),

		// Setup metadata
		micro.Metadata(versionMeta),
	)
	// optionally setup command line usage
	service.Init(
		micro.Action(func(c *cli.Context) error {
			if c.Bool("version") {
				config.PrintFullVersionInfo()
				os.Exit(0)
			}
			return nil
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
	biz := newBizEngineManager()
	awsClient, err := newAWSClient()
	if err != nil {
		log.Panicf("Failed to connect to aws s3 server at %s. Error: %v", awsBucketName, err)
	}
	analysisManagerService := handler.NewAnalysisManagerService(db, biz, presetsFilePath, awsClient)
	if err := proto.RegisterAnalysisManagerAPIHandler(server, analysisManagerService); err != nil {
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

func defaultVersionFlags() micro.Option {
	return micro.Flags(
		&cli.BoolFlag{
			Name:  "version",
			Usage: "Show version information",
		},
	)
}

// dbClientOptions 构建命令行启动参数
func dbClientOptions() micro.Option {
	return micro.Flags(
		&cli.StringFlag{
			Name:        "x_db_address",
			Value:       "localhost:3306",
			Usage:       "MySQL instance `ADDRESS` - [host]:[port]",
			EnvVars:     []string{"X_DB_ADDRESS"},
			Destination: &dbAddress,
		},
		&cli.StringFlag{
			Name:        "x_db_username",
			Usage:       "MySQL login `USERNAME`",
			EnvVars:     []string{"X_DB_USERNAME"},
			Destination: &dbUsername,
		},
		&cli.StringFlag{
			Name:        "x_db_password",
			Usage:       "MySQL login `PASSWORD`",
			EnvVars:     []string{"X_DB_PASSWORD"},
			Destination: &dbPassword,
		},
		&cli.StringFlag{
			Name:        "x_db_database",
			Usage:       "MySQL database name",
			EnvVars:     []string{"X_DB_DATABASE"},
			Destination: &dbDatabase,
		},
		&cli.BoolFlag{
			Name:        "x_db_enable_log",
			Usage:       "Enable MySQL client log",
			EnvVars:     []string{"X_DB_ENABLE_LOG"},
			Destination: &dbEnableLog,
		},
		&cli.IntFlag{
			Name:        "x_db_max_connections",
			Usage:       "Max connections of MySQL client",
			EnvVars:     []string{"X_DB_MAX_CONNECTIONS"},
			Value:       1,
			Destination: &dbMaxConns,
		},
	)
}

// ae 相关配置
var (
	luaSrcPath      string
	templatesDir    string
	questionDir     string
	presetsFilePath string
	productionAELog bool
)

// aeOptions 构建AE命令行启动参数
func aeOptions() micro.Option {
	return micro.Flags(
		&cli.StringFlag{
			Name:        "x_lua_src_path",
			Usage:       "AE LuaSrcPath",
			EnvVars:     []string{"X_LUA_SRC_PATH"},
			Destination: &luaSrcPath,
		},
		&cli.StringFlag{
			Name:        "x_templates_dir",
			Usage:       "AE TemplatesDir",
			EnvVars:     []string{"X_TEMPLATES_DIR"},
			Destination: &templatesDir,
		},
		&cli.StringFlag{
			Name:        "x_question_dir",
			Usage:       "AE QuestionDir",
			EnvVars:     []string{"X_QUESTION_DIR"},
			Destination: &questionDir,
		},
		&cli.StringFlag{
			Name:        "x_presets_file_path",
			Usage:       "AE PresetsFilePath",
			EnvVars:     []string{"X_PRESETS_FILE_PATH"},
			Destination: &presetsFilePath,
		},
		&cli.BoolFlag{
			Name:        "x_production_ae_log",
			Usage:       "Log Mode Of AE",
			EnvVars:     []string{"X_PRODUCTION_AE_LOG"},
			Destination: &productionAELog,
		},
	)
}

// aws 存储桶连接信息
var (
	awsBucketName                          string
	awsAccessKey                           string
	awsSecretKey                           string
	awsRegion                              string
	pulseTestRawDataEnvironmentS3KeyPrefix string
	pulseTestRawDataS3KeyPrefix            string
)

// awsClientOptions 构建命令行启动参数
func awsClientOptions() micro.Option {
	return micro.Flags(
		&cli.StringFlag{
			Name:        "x_aws_bucket_name",
			Usage:       "aws bucket name",
			EnvVars:     []string{"X_AWS_BUCKET_NAME"},
			Destination: &awsBucketName,
		},
		&cli.StringFlag{
			Name:        "x_aws_access_key",
			Usage:       "aws access key",
			EnvVars:     []string{"X_AWS_ACCESS_KEY"},
			Destination: &awsAccessKey,
		},
		&cli.StringFlag{
			Name:        "x_aws_secret_key",
			Usage:       "aws secret key",
			EnvVars:     []string{"X_AWS_SECRET_KEY"},
			Destination: &awsSecretKey,
		},
		&cli.StringFlag{
			Name:        "x_aws_region",
			Usage:       "aws region",
			EnvVars:     []string{"X_AWS_REGION"},
			Destination: &awsRegion,
		},
		&cli.StringFlag{
			Name:        "x_wave_data_key_prefix",
			Usage:       "S3 prefix key for wave data",
			EnvVars:     []string{"X_WAVE_DATA_KEY_PREFIX"},
			Destination: &pulseTestRawDataEnvironmentS3KeyPrefix,
		},
		&cli.StringFlag{
			Name:        "x_pulse_test_raw_data_s3_key_prefix",
			Usage:       "pulse test raw data s3 key prefix",
			EnvVars:     []string{"X_PULSE_TEST_RAW_DATA_S3_KEY_PREFIX"},
			Destination: &pulseTestRawDataS3KeyPrefix,
		},
	)
}

// newAWSClient 创建一个 aws 连接
func newAWSClient() (*aws.Client, error) {
	return aws.NewClient(
		aws.BucketName(awsBucketName),
		aws.Region(awsRegion),
		aws.AccessKeyID(awsAccessKey),
		aws.SecretKey(awsSecretKey),
		aws.PulseTestRawDataEnvironmentS3KeyPrefix(pulseTestRawDataEnvironmentS3KeyPrefix),
		aws.PulseTestRawDataS3KeyPrefix(pulseTestRawDataS3KeyPrefix),
	)
}

// newDbClient 创建一个 DbClient
func newDbClient() (*mysqldb.DbClient, error) {
	cfg := dbutilmysql.NewConfig()
	cfg.User = dbUsername
	cfg.Passwd = dbPassword
	cfg.Net = "tcp"
	cfg.Addr = dbAddress
	cfg.DBName = dbDatabase
	db, err := dbutilmysql.OpenGormDB(
		dbutilmysql.WithMySQLConfig(cfg),
	)
	if err != nil {
		return nil, err
	}

	return mysqldb.NewDbClient(*storemysql.NewStore(db))
}

func newBizEngineManager() *biz.BizEngineManager {
	return biz.NewBizEngineManager(
		biz.LuaSrcPath(luaSrcPath),
		biz.TemplatesDir(templatesDir),
		biz.QuestionDir(questionDir),
		biz.PoolSize(2),
		biz.ProductionMode(productionAELog),
	)
}
