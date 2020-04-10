package rest_test

import (
	"net/http"
	"path/filepath"

	"encoding/json"
	"testing"

	"github.com/iris-contrib/httpexpect"
	r "github.com/jinmukeji/gf-api2/api-jinmuid/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// SmsSuite 是Sms的单元测试的 Test Suite
type SmsSuite struct {
	suite.Suite
	Expect  *httpexpect.Expect
	Account *Account
}

// SetupSuite 设置测试环境
func (suite *SmsSuite) SetupSuite() {
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

// SmsNotification 短信的返回
type SmsNotification struct {
	Data r.SmsNotification `json:"data"`
}

// TestSendSigninSmsNotification 发送登录短信
func (suite *SmsSuite) TestSendSigninSmsNotification() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	body := e.POST("/notification/sms").WithHeader("Authorization", auth).WithJSON(
		&r.SmsNotificationBody{
			Phone:            suite.Account.SignInPhone,
			NotificationType: r.SigninSmsNotification,
			Language:         suite.Account.Language,
			NationCode:       suite.Account.NationCode,
		}).Expect().Body()
	var resp SmsNotification
	errUnmarshalSignIn := json.Unmarshal([]byte(body.Raw()), &resp)
	assert.NoError(t, errUnmarshalSignIn)
	assert.NotNil(t, resp.Data.SerialNumber)
	assert.NotEqual(t, "", resp.Data.SerialNumber)
}

func TestSmsSuite(t *testing.T) {
	suite.Run(t, new(SmsSuite))
}
