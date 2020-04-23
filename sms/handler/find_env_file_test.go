package handler

import (
	"os"

	db "github.com/jinmukeji/jiujiantang-services/sms/mysqldb"
	sms "github.com/jinmukeji/jiujiantang-services/sms/sms_client"
	"github.com/joho/godotenv"
)

const (
	enableLog = true
	maxConns  = 1
)

// NewClientOptionsFromEnvFile 读取环境变脸配置文件，返回算法服务器连接配置
func newClientOptionsFromEnvFile(filepath string) *sms.AliyunSMSClient {
	err := godotenv.Load(filepath)
	if err != nil {
		panic(err)
	}
	return &sms.AliyunSMSClient{
		AccessKeyID:     os.Getenv("X_ALIYUN_SMS_ACCESS_KEY_ID"),
		AccessKeySecret: os.Getenv("X_ALIYUN_SMS_ACCESS_KEY_Secret"),
	}
}

// NewTestingDbClientFromEnvFile 从配置文件里面读取环境变量并创建 DbClient
func newTestingDbClientFromEnvFile(filepath string) (*db.DbClient, error) {
	err := godotenv.Load(filepath)
	if err != nil {
		panic(err)
	}
	db, err := db.NewDbClient(
		db.Address(os.Getenv("X_DB_ADDRESS")),
		db.Username(os.Getenv("X_DB_USERNAME")),
		db.Password(os.Getenv("X_DB_PASSWORD")),
		db.Database(os.Getenv("X_DB_DATABASE")),
		db.EnableLog(enableLog),
		db.MaxConnections(maxConns),
	)
	return db, err
}
