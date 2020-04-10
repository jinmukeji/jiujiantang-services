package sem

import (
	"log"
	"os"
	"testing"

	encry "github.com/jinmukeji/go-pkg/crypto/rand"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SEMClientSuite struct {
	suite.Suite
	*AliyunSEMClient
	*NetEaseSEMClient
}

const (
	digits = "0123456789"
	length = 6
)

// Setup 初始化测试
func (suite *SEMClientSuite) SetupSuite() {
	suite.AliyunSEMClient, suite.NetEaseSEMClient = newClientOptionsFromEnvFile("../../build/local.svc-sem-gw.env")
}

// TestSendNetEaseSEMEmail 测试使用网易单发邮件
func (suite *SEMClientSuite) TestSendSingleNetEaseSEMEmail() {
	t := suite.T()
	var client SEMClient
	code := encry.RandomStringWithMask(digits, length)
	log.Println("this is code", code)
	client, errNewNetEaseSEMClient := NewNetEaseSEMClient(suite.USER, suite.PASSWD)
	isSucceed, errSendSem := client.SendEmail("tech@jinmuhealth.com", FindResetPassword, SimplifiedChinese, map[string]string{
		"code": code,
	})
	assert.Equal(t, true, isSucceed)
	assert.NoError(t, errNewNetEaseSEMClient)
	assert.NoError(t, errSendSem)
}

// TestSendNetEaseSEMEmail 测试使用网易群发邮件
func (suite *SEMClientSuite) TestSendMassNetEaseSEMEmail() {
	t := suite.T()
	var client SEMClient
	code := encry.RandomStringWithMask(digits, length)
	log.Println("this is code", code)
	client, errNewNetEaseSEMClient := NewNetEaseSEMClient(suite.USER, suite.PASSWD)
	// 需将testemail替换为具体的测试email地址，否则会报错
	isSucceed, errSendSem := client.SendEmail("tech@jinmuhealth.com;testemail", FindResetPassword, SimplifiedChinese, map[string]string{
		"code": code,
	})
	assert.Equal(t, true, isSucceed)
	assert.NoError(t, errNewNetEaseSEMClient)
	assert.NoError(t, errSendSem)
}

// TestSendAliyunSEMEmail 测试使用阿里云单发邮件
func (suite *SEMClientSuite) TestSendSingleAliyunSEMEmail() {
	t := suite.T()
	var client SEMClient
	code := encry.RandomStringWithMask(digits, length)
	log.Println("this is code", code)
	client, errNewAliyunSEMClient := NewAliyunSEMClient(suite.AccessKeyID, suite.AccessKeySecret)
	isSucceed, errSendSem := client.SendEmail("tech@jinmuhealth.com", FindResetPassword, SimplifiedChinese, map[string]string{
		"code": code,
	})
	assert.Equal(t, true, isSucceed)
	assert.NoError(t, errNewAliyunSEMClient)
	assert.NoError(t, errSendSem)
}

// TestSendAliyunSEMEmail 测试使用阿里云群发邮件
func (suite *SEMClientSuite) TestSendMassAliyunSEMEmail() {
	t := suite.T()
	var client SEMClient
	code := encry.RandomStringWithMask(digits, length)
	log.Println("this is code", code)
	client, errNewAliyunSEMClient := NewAliyunSEMClient(suite.AccessKeyID, suite.AccessKeySecret)
	// 需将testemail替换为具体的测试email地址，否则会报错
	isSucceed, errSendSem := client.SendEmail("tech@jinmuhealth.com,testemail", FindResetPassword, SimplifiedChinese, map[string]string{
		"code": code,
	})
	assert.Equal(t, true, isSucceed)
	assert.NoError(t, errNewAliyunSEMClient)
	assert.NoError(t, errSendSem)
	// 测试阿里云的批量发送邮件
	// isSucceed, errSendSem = client.BatchSendEmail("test", "ninatest", "")
	// assert.Equal(t, true, isSucceed)
	// assert.NoError(t, errNewAliyunSEMClient)
	// assert.NoError(t, errSendSem)
}

func TestSEMClientSuite(t *testing.T) {
	suite.Run(t, new(SEMClientSuite))
}

// newClientOptionsFromEnvFile 读取环境变脸配置文件，返回算法服务器连接配置
func newClientOptionsFromEnvFile(filepath string) (*AliyunSEMClient, *NetEaseSEMClient) {
	err := godotenv.Load(filepath)
	if err != nil {
		panic(err)
	}
	return &AliyunSEMClient{
			AccessKeyID:     os.Getenv("X_ALIYUN_SEM_ACCESS_KEY_ID"),
			AccessKeySecret: os.Getenv("X_ALIYUN_SEM_ACCESS_KEY_Secret"),
		},
		&NetEaseSEMClient{
			USER:   os.Getenv("X_NETEASE_SEM_USER"),
			PASSWD: os.Getenv("X_NETEASE_SEM_PASSWD"),
		}
}
