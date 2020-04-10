package handler

import (
	"os"

	db "github.com/jinmukeji/gf-api2/analysis/mysqldb"
	"github.com/joho/godotenv"
)

const (
	enableLog = true
	maxConns  = 1
)

// newTestingDbClientFromEnvFile 从配置文件里面读取环境变量并创建 DbClient
func newTestingDbClientFromEnvFile(filepath string) (*db.DbClient, error) {
	errLoad := godotenv.Load(filepath)
	if errLoad != nil {
		return nil, errLoad
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
