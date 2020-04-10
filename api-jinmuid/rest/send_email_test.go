package rest_test

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/iris-contrib/httpexpect"
	r "github.com/jinmukeji/jiujiantang-services/api-jinmuid/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// SendEmailSuite 是SendEmail的单元测试的 Test Suite
type SendEmailSuite struct {
	suite.Suite
	Expect  *httpexpect.Expect
	Account *Account
}

// SetupSuite 设置测试环境
func (suite *SendEmailSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.api-jinmuid.env")
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
	t := suite.T()
	app := r.NewApp("", "jinmuhealth", false)
	suite.Expect = httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(app),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewCurlPrinter(t),
			httpexpect.NewDebugPrinter(t, true),
		},
	})
}

// EmailNotificationResponse 邮件的返回
type EmailNotificationResponse struct {
	Data r.EmailNotificationResponse `json:"data"`
}

// TestSendEmailNotificationToFindUsername 测试发送邮件找回用户名
func (suite *SendEmailSuite) TestSendEmailNotificationToFindUsername() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	body := e.POST("/notification/email").WithHeader("Authorization", auth).WithJSON(
		&r.EmailNotificationRequest{
			Email:    suite.Account.Email,
			Type:     r.FindUsernameSemNotification,
			Language: suite.Account.Language,
		},
	).Expect().Body()
	var resp EmailNotificationResponse
	errUnmarshalSignIn := json.Unmarshal([]byte(body.Raw()), &resp)
	assert.NoError(t, errUnmarshalSignIn)
	assert.NotNil(t, resp.Data.SerialNumber)
	assert.NotEqual(t, "", resp.Data.SerialNumber)
}

func TestSendEmailSuite(t *testing.T) {
	suite.Run(t, new(SendEmailSuite))
}
