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

// SearchWeekAndMonthHistoryTestSuite 是SearchWeekAndMonthHistory的单元测试的 Test Suite
type SearchWeekAndMonthHistoryTestSuite struct {
	suite.Suite
	Account *Account
	Expect  *httpexpect.Expect
}

// SetupSuite 设置测试环境
func (suite *SearchWeekAndMonthHistoryTestSuite) SetupSuite() {
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

// TestSearchWeekHistory 周的趋势
func (suite *SearchWeekAndMonthHistoryTestSuite) TestSearchWeekHistory() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, errGetAccessToken := getAccessToken(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	e.GET("/v2-api/owner/week_measurements").
		WithHeaders(headers).
		WithQuery("user_id", suite.Account.UserID).
		Expect().Body().Contains("c0").Contains("c1").Contains("c2").Contains("c3").Contains("c4").Contains("c5").Contains("c6").Contains("c7")
}

// SearchMonthHistory 月的趋势
func (suite *SearchWeekAndMonthHistoryTestSuite) TestSearchMonthHistory() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, errGetAccessToken := getAccessToken(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	e.GET("/v2-api/owner/month_measurements").
		WithHeaders(headers).
		WithQuery("user_id", suite.Account.UserID).
		Expect().Body().Contains("c0").Contains("c1").Contains("c2").Contains("c3").Contains("c4").Contains("c5").Contains("c6").Contains("c7")
}

func TestSearchWeekAndMonthHistoryTestSuite(t *testing.T) {
	suite.Run(t, new(SearchWeekAndMonthHistoryTestSuite))
}
