package rest_test

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/iris-contrib/httpexpect"
	r "github.com/jinmukeji/jiujiantang-services/api-sys/rest"
	"github.com/stretchr/testify/suite"
)

// ClientPreferencesTestSuite 是client_preferences的单元测试的 Test Suite
type ClientPreferencesTestSuite struct {
	suite.Suite
	Client *Client
}

// SetupSuite 设置测试环境
func (suite *ClientPreferencesTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.api-sys.env")
	suite.Client = newTestingClientFromEnvFile(envFilepath)
}

// TestClientPreferences 测试ClientPreferences
func (suite *ClientPreferencesTestSuite) TestClientPreferences() {
	t := suite.T()
	app := r.NewApp("", "../data/config.yml")
	e := httpexpect.WithConfig(httpexpect.Config{
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
	e.POST("/_sys/client_preferences").WithJSON(&r.ClientPreferencesBody{
		ClientID:      suite.Client.ClientID,
		SecretKeyHash: suite.Client.SecretKeyHash,
		Seed:          suite.Client.Seed,
		ClientVersion: suite.Client.ClientVersion,
		Environment:   suite.Client.Environment,
	},
	).Expect().Body().Contains("https://testing-api.jinmuhealth.com:37633/v2-api").
		Contains("https://res-cdn.jinmuhealth.com/v2-testing/app-login/2-1/index.html").
		Contains("https://res-cdn.jinmuhealth.com/v2-testing/app-entry/2-1").
		Contains("https://res-cdn.jinmuhealth.com/v2-testing/app-faq/2-1").
		Contains("https://res-cdn.jinmuhealth.com/v2-testing/app-report/2-1/index.html")
}

func TestClientPreferencesTestSuite(t *testing.T) {
	suite.Run(t, new(ClientPreferencesTestSuite))
}
