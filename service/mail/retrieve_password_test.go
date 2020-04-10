package mail

import (
	"testing"

	"path/filepath"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// MailTestSuite 是邮件发送的测试suite
type RetrievePasswordTestSuite struct {
	suite.Suite
	client *Client
}

func (suite *RetrievePasswordTestSuite) SetupSuite() {
	options, _ := newMailOptionfromEnvFile(filepath.Join("testdata", "local.svc-biz-core.env"))
	suite.client = &Client{options}
}

// TestMailContentRetrievePassword 测试找回密码内容的生成
func (suite *RetrievePasswordTestSuite) TestMailContentRetrievePassword() {
	t := suite.T()
	option, err := newMailOptionfromEnvFile(filepath.Join("testdata", "local.svc-biz-core.env"))
	assert.NoError(t, err)
	assert.NotNil(t, option)
	content, _ := suite.client.MailContentRetrievePassword(InfoRetrievePassword{
		"jinmu", "123",
	})
	assert.NotNil(t, content)
	assert.Equal(t, content, `<h3>尊敬的 jinmu</h3><p>您好，系统已收到您的找回密码申请，以下为你的脉诊仪账户所需要的重要信息，请注意保密。</p><p>您的脉诊仪账号密码为：123</p>
		<p align='center'><img width= '200pt' height ='200pt' src= 'http://www.jinmuhealth.com/assets/img/QRcode/wechatQRcode.jpg',  alt='金姆健康科技'/></p>
		<p align='center'>扫一扫关注“金姆健康科技”</p>
		<p>请注意，该邮件地址不接收回复邮件，请联系金姆客服进行咨询</p>
		<p>客服电话：0519-81180075</p>
		<p>E-mail:information@jinmuhealth.com</p>`)
	assert.NoError(t, err)
}

// TestSendMailRetrievePassword 测试找回密码邮件的发送
func (suite *RetrievePasswordTestSuite) TestSendMailRetrievePassword() {
	t := suite.T()
	err := suite.client.SendMailRetrievePassword(Contact{"zhuyingjie@jinmuhealth.com", "朱英杰"}, "123")
	assert.NoError(t, err)
}

// TestMailTestSuite 启动邮件测试
func TestRetrievePasswordTestSuite(t *testing.T) {
	suite.Run(t, new(RetrievePasswordTestSuite))
}
