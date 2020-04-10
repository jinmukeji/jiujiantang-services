package handler

import (
	"os"
	"strconv"

	db "github.com/jinmukeji/gf-api2/subscription/mysqldb"
	"github.com/joho/godotenv"
)

const (
	enableLog = true
	maxConns  = 1
)

func newTestingAccountFromEnvFile(filepath string) *Account {
	_ = godotenv.Load(filepath)
	userID, _ := strconv.Atoi(os.Getenv("X_TEST_USER_ID"))
	contractYear, _ := strconv.Atoi(os.Getenv("X_CONTRACT_YEAR"))
	maxUserLimits, _ := strconv.Atoi(os.Getenv("X_MAX_USER_LIMITS"))
	return &Account{
		os.Getenv("X_TEST_ACCOUNT"),
		os.Getenv("X_TEST_PASSWORD"),
		int32(userID),
		os.Getenv("X_TEST_SEED"),
		os.Getenv("X_TEST_HASHED_PASSWORD"),
		os.Getenv("X_TEST_CODE"),
		os.Getenv("X_ACTIVATION_CODE_ENCRYPT_KEY"),
		os.Getenv("X_ACTIVATION_CODE"),
		int32(contractYear),
		int32(maxUserLimits),
	}
}

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
