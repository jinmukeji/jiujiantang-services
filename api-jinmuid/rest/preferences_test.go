package rest_test

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/iris-contrib/httpexpect"
	r "github.com/jinmukeji/jiujiantang-services/api-jinmuid/rest"
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
	app := r.NewApp("", "jinmuhealth", false)
	envFilepath := filepath.Join("testdata", "local.api-jinmuid.env")
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

// TestGetUserPreferences 测试GetUserPreferences
func (suite *PreferencesTestSuite) TestGetUserPreferences() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, userID, errGetAccessToken := getAccessTokenAndUserID(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	e.GET("/user/{user_id}/preferences").
		WithHeaders(headers).
		WithPath("user_id", userID).
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
