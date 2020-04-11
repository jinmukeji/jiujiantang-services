package handler

import (
	"context"
	"path/filepath"
	"testing"

	encry "github.com/jinmukeji/go-pkg/crypto/rand"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/sem/v1"
	generalpb "github.com/jinmukeji/proto/gen/micro/idl/ptypes/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	digits = "0123456789"
	length = 6
)

// SendEmailSuite 测试发送邮件
type SendEmailSuite struct {
	suite.Suite
	semGateway *SEMGateway
}

// SetupSuite 设置测试环境
func (suite *SendEmailSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-sem-gw.env")
	datastore, _ := newTestingDbClientFromEnvFile(envFilepath)
	aliyunSEMClient, neteaseSEMClient := newClientOptionsFromEnvFile(envFilepath)
	suite.semGateway = NewSEMGateway(datastore, aliyunSEMClient, neteaseSEMClient)
}

// TestSendEmail 测试单发邮件
func (suite *SendEmailSuite) TestSendSingleEmail() {
	t := suite.T()
	code := encry.RandomStringWithMask(digits, length)
	log.Println("this is code", code)
	req, resp := new(proto.EmailNotificationRequest), new(proto.EmailNotificationResponse)

	req.ToAddress = []string{"tech@jinmuhealth.com"}
	req.TemplateAction = proto.TemplateAction_TEMPLATE_ACTION_FIND_RESET_PASSWORD
	req.TemplateParam = map[string]string{"code": code}
	req.Language = generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE
	assert.NoError(t, suite.semGateway.EmailNotification(context.Background(), req, resp))
}

// TestSendEmail 测试群发邮件
func (suite *SendEmailSuite) TestSendMassEmail() {
	t := suite.T()
	code := encry.RandomStringWithMask(digits, length)
	log.Println("this is code", code)
	req, resp := new(proto.EmailNotificationRequest), new(proto.EmailNotificationResponse)

	// 需将testemail替换为具体的测试email地址，否则会报错
	req.ToAddress = []string{"tech@jinmuhealth.com", "testemail"}
	req.TemplateAction = proto.TemplateAction_TEMPLATE_ACTION_FIND_RESET_PASSWORD
	req.TemplateParam = map[string]string{"code": code}
	req.Language = generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE
	assert.NoError(t, suite.semGateway.EmailNotification(context.Background(), req, resp))
}

func TestSendEmailTestSuite(t *testing.T) {
	suite.Run(t, new(SendEmailSuite))
}
