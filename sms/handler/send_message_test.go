package handler

import (
	"context"
	"path/filepath"
	"testing"

	smspb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/sms/v1"
	generalpb "github.com/jinmukeji/proto/v3/gen/micro/idl/ptypes/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// SendMessageSuite 测试发送短信
type SendMessageSuite struct {
	suite.Suite
	smsGateway *SMSGateway
}

// SetupSuite 设置测试环境
func (suite *SendMessageSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-sms-gw.env")
	db, _ := newTestingDbClientFromEnvFile(envFilepath)
	aliyunSMSClient := newClientOptionsFromEnvFile(envFilepath)
	suite.smsGateway = NewSMSGateway(db, aliyunSMSClient)
}

// TestSendMessage 测试发送短信
func (suite *SendMessageSuite) TestSendMessage() {
	t := suite.T()
	req, resp := new(smspb.SendMessageRequest), new(smspb.SendMessageResponse)
	req.IsForced = false
	req.Phone = "18805177594"
	req.TemplateAction = smspb.TemplateAction_TEMPLATE_ACTION_SIGN_UP
	req.NationCode = "+86"
	req.TemplateParam = map[string]string{
		"code": "1000",
	}
	req.Language = generalpb.Language_LANGUAGE_ENGLISH
	assert.NoError(t, suite.smsGateway.SendMessage(context.Background(), req, resp))
}

func TestSendMessageTestSuite(t *testing.T) {
	suite.Run(t, new(SendMessageSuite))
}
