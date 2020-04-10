package rest_test

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/iris-contrib/httpexpect"
	r "github.com/jinmukeji/jiujiantang-services/api-v2/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// LocalNotificationTestSuite 是LocalNotification的单元测试的 Test Suite
type LocalNotificationTestSuite struct {
	suite.Suite
	Account *Account
	Expect  *httpexpect.Expect
}

// SetupSuite 设置测试环境
func (suite *LocalNotificationTestSuite) SetupSuite() {
	t := suite.T()
	envFilepath := filepath.Join("testdata", "local.api-v2.env")
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
	app := r.NewApp("v2-api", "jinmuhealth")
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

// TestGetLocalNotifications 得到LocalNotifications
func (suite *LocalNotificationTestSuite) TestGetLocalNotifications() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, errGetAccessToken := getAccessToken(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	e.GET("/v2-api/configuration/local_notification").
		WithHeaders(headers).
		Expect().Body().Contains("content").Contains("喜马宝宝提醒您，又到检查时间啦！")
}

func TestLocalNotificationTestSuite(t *testing.T) {
	suite.Run(t, new(LocalNotificationTestSuite))
}
