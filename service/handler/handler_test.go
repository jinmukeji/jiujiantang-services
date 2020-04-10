package handler

import (
	"context"
	"os"
	"strconv"

	"github.com/jinmukeji/jiujiantang-services/service/auth"
	"github.com/jinmukeji/jiujiantang-services/service/mail"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	"github.com/joho/godotenv"
)

const testUsername = "xx"

// newTestingDbClientFromEnvFile 从配置文件里面读取环境变量并创建 mail options
func newTestingMailClientFromEnvFile(filepath string) (*mail.Client, error) {
	_ = godotenv.Load(filepath)
	port, _ := strconv.ParseInt(os.Getenv("X_MAIL_PORT"), 0, 32)

	options := &mail.Options{}
	mail.Address(os.Getenv("X_MAIL_ADDRESS"))(options)
	mail.Username(os.Getenv("X_MAIL_USERNAME"))(options)
	mail.Password(os.Getenv("X_MAIL_PASSWORD"))(options)
	mail.Port(int(port))(options)
	mail.Charset(os.Getenv("X_MAIL_CHARSET"))(options)
	mail.SenderNickname(os.Getenv("X_MAIL_NICKNAME"))(options)
	client := &mail.Client{}
	client.SetOptions(options)

	return client, nil
}

// mockLogin mock 一次登录请求 用于其他 api 的单元测试
func mockLogin(j *JinmuHealth, username string, passwordHash string) (context.Context, error) {
	req, resp := new(proto.UserSignInRequest), new(proto.UserSignInResponse)
	ctx := context.Background()
	req.SignInKey = username
	req.PasswordHash = passwordHash
	if err := j.UserSignIn(ctx, req, resp); err != nil {
		return nil, err
	}
	return auth.AddContextUserID(ctx, 1), nil
}
