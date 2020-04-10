package rest_test

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/iris-contrib/httpexpect"
	r "github.com/jinmukeji/gf-api2/api-v2/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// PreferencesTestSuite 是Preferences的单元测试的 Test Suite
type PreferencesTestSuite struct {
	suite.Suite
	Expect  *httpexpect.Expect
	Account *Account
}

// SetupSuite 设置测试环境
func (suite *PreferencesTestSuite) SetupSuite() {
	t := suite.T()
	app := r.NewApp("v2-api", "jinmuhealth")
	envFilepath := filepath.Join("testdata", "local.api-v2.env")
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
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

// TestOwnerGetUserPreferences 测试OwnerGetUserPreferences
func (suite *PreferencesTestSuite) TestOwnerGetUserPreferences() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, errGetAccessToken := getAccessToken(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	e.GET("/v2-api/owner/users/{user_id}/preferences").
		WithHeaders(headers).
		WithPath("user_id", suite.Account.UserID).
		Expect().Body().
		Contains("enable_heart_rate_chart").
		Contains("enable_pulse_wave_chart").
		Contains("enable_warm_prompt").
		Contains("enable_choose_status").
		Contains("enable_constitution_differentiation").
		Contains("enable_syndrome_differentiation").
		Contains("enable_western_medicine_analysis").
		Contains("enable_meridian_bar_graph").
		Contains("enable_comment").
		Contains("enable_health_trending")
}

func TestPreferencesTestSuite(t *testing.T) {
	suite.Run(t, new(PreferencesTestSuite))
}
