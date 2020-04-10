package mail

import (
	"testing"

	"path/filepath"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// MailTestSuite 是邮件发送的测试suite
type MailTestSuite struct {
	suite.Suite
	client *Client
}

func (suite *MailTestSuite) SetupSuite() {
	options, _ := newMailOptionfromEnvFile(filepath.Join("testdata", "local.svc-biz-core.env"))
	suite.client = &Client{options}
}

// TestSendMailRetrievePassword 测试找回密码邮件的发送
func (suite *MailTestSuite) TestSendMailRetrievePassword() {
	t := suite.T()
	err := suite.client.SendMail(Contact{"zhuyingjie@jinmuhealth.com", "朱英杰"}, "主题", "内容", textHTML)
	assert.NoError(t, err)
}

// TestMailTestSuite 启动邮件测试
func TestMailTestSuite(t *testing.T) {
	suite.Run(t, new(MailTestSuite))
}
