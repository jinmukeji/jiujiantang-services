package handler

import (
	"os"

	db "github.com/jinmukeji/jiujiantang-services/sem/mysqldb"
	sem "github.com/jinmukeji/jiujiantang-services/sem/sem_client"
	"github.com/joho/godotenv"
)

const (
	enableLog = true
	maxConns  = 1
)

// newTestingDbClientFromEnvFile 从配置文件里面读取环境变量并创建 DbClient
func newTestingDbClientFromEnvFile(filepath string) (*db.DbClient, error) {
	_ = godotenv.Load(filepath)
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

// newClientOptionsFromEnvFile 读取环境变脸配置文件，返回算法服务器连接配置
func newClientOptionsFromEnvFile(filepath string) (*sem.AliyunSEMClient, *sem.NetEaseSEMClient) {
	err := godotenv.Load(filepath)
	if err != nil {
		panic(err)
	}
	return &sem.AliyunSEMClient{
			AccessKeyID:     os.Getenv("X_ALIYUN_SEM_ACCESS_KEY_ID"),
			AccessKeySecret: os.Getenv("X_ALIYUN_SEM_ACCESS_KEY_Secret"),
		},
		&sem.NetEaseSEMClient{
			USER:   os.Getenv("X_NETEASE_SEM_USER"),
			PASSWD: os.Getenv("X_NETEASE_SEM_PASSWD"),
		}
}
