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
	SignInKey                      string
	PasswordHash                   string
	UserID                         int
	Nickname                       string
	RegisterType                   string
	RecordID                       int
	RecordToken                    string
	Remark                         string
	SubscriptionType               string
	ActivationCode                 string
	ActivationCodeSubscriptionType string
	UsedActivationCode             string
}

// newTestingAccountFromEnvFile 从环境文件中获取Account
func newTestingAccountFromEnvFile(filepath string) *Account {
	_ = godotenv.Load(filepath)
	userID, _ := strconv.Atoi(os.Getenv("X_TEST_USER_ID"))
	recordID, _ := strconv.Atoi(os.Getenv("X_TEST_RECORD_ID"))
	return &Account{
		os.Getenv("X_TEST_SIGN_IN_KEY"),
		os.Getenv("X_TEST_PASSWORD_HASH"),
		userID,
		os.Getenv("X_TEST_NICKNAME"),
		os.Getenv("X_TEST_REGISTER_TYPE"),
		recordID,
		os.Getenv("X_TEST_RECORD_TOKEN"),
		os.Getenv("X_TEST_REMARK"),
		os.Getenv("X_TEST_SUBSCRIPTION_TYPE"),
		os.Getenv("X_TEST_ACTIVATION_CODE"),
		os.Getenv("X_TEST_ACTIVATION_CODE_SUBSCRIPTION_TYPE"),
		os.Getenv("X_TEST_USED_ACTIVATION_CODE"),
	}
}
