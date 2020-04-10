package mail

import (
	"testing"

	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// MailTestSuite 是邮件发送的测试suite
type OptionsTestSuite struct {
	suite.Suite
}

// newTestingDbClientFromEnvFile 从配置文件里面读取环境变量并创建 mail options
func newMailOptionfromEnvFile(filepath string) (*Options, error) {
	err := godotenv.Load(filepath)
	if err != nil {
		log.Fatalf("Error loading %s file", filepath)
		return nil, err
	}
	port, _ := strconv.ParseInt(os.Getenv("X_MAIL_PORT"), 0, 32)

	options := newOptions(
		Address(os.Getenv("X_MAIL_ADDRESS")),
		Username(os.Getenv("X_MAIL_USERNAME")),
		Password(os.Getenv("X_MAIL_PASSWORD")),
		Port(int(port)),
		Charset(os.Getenv("X_MAIL_CHARSET")),
		SenderNickname(os.Getenv("X_MAIL_NICKNAME")),
		ReplyToAddress(os.Getenv("X_([a-zA-Z]*?).([a-zA-Z]*?)Response")),
	)
	return options, nil

}

// TestDefaultOptions 测试 defaultOptions 方法成功返回 Options 记录
func (suite *OptionsTestSuite) TestDefaultOptions() {
	t := suite.T()
	opts := defaultOptions()

	assert.EqualValues(t, Options{
		Address:  "localhost",
		Username: "root",
		Password: "",
		Charset:  "UTF-8",
		Port:     25,
	}, *opts)
}

// TestNewOptions 测试 newOptions 方法成功返回 Options 记录
func (suite *OptionsTestSuite) TestNewOptions() {
	t := suite.T()
	opts := newOptions()

	assert.EqualValues(t, Options{
		Address:  "localhost",
		Username: "root",
		Password: "",
		Charset:  "UTF-8",
		Port:     25,
	}, *opts)
}

// TestNewOptions 测试带参的 newOptions 方法成功返回 Options 记录
func (suite *OptionsTestSuite) TestNewOptionsWithParameters() {
	t := suite.T()
	opts := newOptions(
		Address("0.0.0.0:6606"),
		Charset("utf8"),
		Username("xx"),
		Password("yy"),
		Port(11),
	)

	assert.EqualValues(t, Options{
		Address:  "0.0.0.0:6606",
		Charset:  "utf8",
		Username: "xx",
		Password: "yy",
		Port:     11,
	}, *opts)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestOptionsTestSuite(t *testing.T) {
	suite.Run(t, new(OptionsTestSuite))
}
