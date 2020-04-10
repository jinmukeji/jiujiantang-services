package handler

import (
	"os"
	"strconv"

	db "github.com/jinmukeji/jiujiantang-services/service/mysqldb"
	"github.com/joho/godotenv"
)

const (
	enableLog = false
	maxConns  = 1
)

type Account struct {
	secretKeyHash       string
	seed                string
	userAccount         string
	password            string
	machineUuid         string
	errorAccount        string
	errorPassword       string
	errorMachineUuid    string
	userAccountNotExist string
	errorFormatPassword string

	organizationID int32
	userID         int32
	userName       string
	passwordHash   string
	clientID       string
	name           string
	zone           string
}

func newTestingAccountFromEnvFile(filepath string) *Account {
	_ = godotenv.Load(filepath)
	userID, _ := strconv.Atoi(os.Getenv("X_TEST_USER_ID"))
	organizationID, _ := strconv.Atoi(os.Getenv("X_TEST_ORGANIZATION_ID"))
	return &Account{
		os.Getenv("X_TEST_SECRET_KEY_HASH"),
		os.Getenv("X_TEST_SEED"),
		os.Getenv("X_TEST_USER_ACCOUNT"),
		os.Getenv("X_TEST_PASSWORD"),
		os.Getenv("X_TEST_MACHINE_UUID"),
		os.Getenv("X_TEST_ERROR_ACCOUNT"),
		os.Getenv("X_TEST_ERROR_PASSWORD"),
		os.Getenv("X_TEST_ERROR_MACHINE_UUID"),
		os.Getenv("X_TEST_USER_ACCOUNT_NOT_EXIST"),
		os.Getenv("X_TEST_ERROR_FORMAT_PASSWORD"),
		int32(organizationID),
		int32(userID),
		os.Getenv("X_TEST_USER_NAME"),
		os.Getenv("X_TEST_PASSWORD_HASH"),
		os.Getenv("X_TEST_CLIENT_ID"),
		os.Getenv("X_TEST_NAME"),
		os.Getenv("X_TEST_ZONE"),
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

func newTestingJinmuHealthFromEnvFile(path string) *JinmuHealth {
	err := godotenv.Load(path)
	if err != nil {
		panic(err)
	}
	j := new(JinmuHealth)
	return j
}
