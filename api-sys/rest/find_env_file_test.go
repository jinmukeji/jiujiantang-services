package rest_test

import (
	"os"

	"github.com/joho/godotenv"
)

// Client 客户端
type Client struct {
	ClientID      string
	SecretKeyHash string
	Seed          string
	ClientVersion string
	Environment   string
}

func newTestingClientFromEnvFile(filepath string) *Client {
	_ = godotenv.Load(filepath)
	return &Client{
		os.Getenv("X_TEST_CLIENT_ID"),
		os.Getenv("X_TEST_SECRET_KEY_HASH"),
		os.Getenv("X_TEST_SEED"),
		os.Getenv("X_TEST_CLIENT_VERSION"),
		os.Getenv("X_TEST_ENVIRONMENT"),
	}
}
