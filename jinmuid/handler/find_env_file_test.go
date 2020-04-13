package handler

import (
	"os"

	"github.com/jinmukeji/go-pkg/v2/mysqldb"
	db "github.com/jinmukeji/jiujiantang-services/jinmuid/mysqldb"
	"github.com/joho/godotenv"
)

const (
	enableLog = false
	maxConns  = 1
)

func newTestingDbClientFromEnvFile(filepath string) (*db.DbClient, error) {
	_ = godotenv.Load(filepath)
	db, err := db.NewDbClient(
		mysqldb.Address(os.Getenv("X_DB_ADDRESS")),
		mysqldb.Username(os.Getenv("X_DB_USERNAME")),
		mysqldb.Password(os.Getenv("X_DB_PASSWORD")),
		mysqldb.Database(os.Getenv("X_DB_DATABASE")),
		mysqldb.EnableLog(enableLog),
		mysqldb.MaxConnections(maxConns),
	)
	return db, err
}

func newTestingEncryptKeyFromEnvFile(filepath string) string {
	_ = godotenv.Load(filepath)
	return os.Getenv("X_ENCRYPT_KEY")
}
