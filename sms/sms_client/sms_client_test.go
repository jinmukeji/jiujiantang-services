package sms

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SMSClientSuite struct {
	suite.Suite
	*AliyunSMSClient
}

// Setup 初始化测试
func (suite *SMSClientSuite) SetupSuite() {
	suite.AliyunSMSClient = newClientOptionsFromEnvFile("../../build/local.svc-sms-gw.env")
}

// TestSendAliyunSMSMessage 测试发送阿里云的发送短信
func (suite *SMSClientSuite) TestSendAliyunSMSMessage() {
	t := suite.T()
	var client SMSClient
	client, errNewAliyunSMSClient := NewAliyunSMSClient(suite.AccessKeyID, suite.AccessKeySecret)
	isSucceed, errSendSms := client.SendSms("18805177594", "+86", SignUp, SimpleChinese, map[string]string{
		"code": "1234",
	})
	assert.Equal(t, true, isSucceed)
	assert.NoError(t, errNewAliyunSMSClient)
	assert.NoError(t, errSendSms)
}

func TestSMSClientSuite(t *testing.T) {
	suite.Run(t, new(SMSClientSuite))
}

// newClientOptionsFromEnvFile 读取环境变脸配置文件，返回算法服务器连接配置
func newClientOptionsFromEnvFile(filepath string) *AliyunSMSClient {
	err := godotenv.Load(filepath)
	if err != nil {
		panic(err)
	}
	return &AliyunSMSClient{
		AccessKeyID:     os.Getenv("X_ALIYUN_SMS_ACCESS_KEY_ID"),
		AccessKeySecret: os.Getenv("X_ALIYUN_SMS_ACCESS_KEY_Secret"),
	}
}
