package rest_test

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Client 客户端
type Client struct {
	ClientID      string
	SecretKeyHash string
	Seed          string
}

// newTestingClientFromEnvFile 从环境文件中获取Client
func newTestingClientFromEnvFile(filepath string) *Client {
	_ = godotenv.Load(filepath)
	return &Client{
		os.Getenv("X_TEST_CLIENT_ID"),
		os.Getenv("X_TEST_SECRET_KEY_HASH"),
		os.Getenv("X_TEST_SEED"),
	}
}

// Account 账户
type Account struct {
	Username       string
	HashedPassword string
	Seed           string
	SignInMachine  string
	SignInPhone    string
	NationCode     string
	Language       string
	Email          string
	PlainPassword  string
	Nickname       string
	Region         int32
}

// newTestingAccountFromEnvFile 从环境文件中获取Account
func newTestingAccountFromEnvFile(filepath string) *Account {
	_ = godotenv.Load(filepath)
	region, _ := strconv.Atoi(os.Getenv("X_TEST_REGION"))
	return &Account{
		os.Getenv("X_TEST_USERNAME"),
		os.Getenv("X_TEST_PASSWORD_HASH"),
		os.Getenv("X_TEST_SIGN_IN_SEED"),
		os.Getenv("X_TEST_SIGN_IN_MACHINE"),
		os.Getenv("X_TEST_SIGN_IN_PHONE"),
		os.Getenv("X_TEST_NATION_CODE"),
		os.Getenv("X_TEST_LANGUAGE"),
		os.Getenv("X_TEST_EMAIL"),
		os.Getenv("X_TEST_PLAIN_PASSWORD"),
		os.Getenv("X_TEST_NICKNAME"),
		int32(region),
	}
}
