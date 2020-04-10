package rest_test

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/iris-contrib/httpexpect"
	r "github.com/jinmukeji/gf-api2/api-v2/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// AnalysisV2TestSuite 是AnalysisV2的单元测试的 Test Suite
type AnalysisV2TestSuite struct {
	suite.Suite
	Expect  *httpexpect.Expect
	Account *Account
}

// SetupSuite 设置测试环境
func (suite *AnalysisV2TestSuite) SetupSuite() {
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

// TestGetV2AnalyzeResult 测试GetV2AnalyzeResult
func (suite *AnalysisV2TestSuite) TestGetV2AnalyzeResult() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, errGetAccessToken := getAccessToken(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	envFilepath := filepath.Join("testdata", "analysis_report_request_body.json")
	jsondata, _ := ioutil.ReadFile(envFilepath)
	e.POST("/v2-api/owner/measurements/{record_id}/v2/analyze").
		WithHeaders(headers).
		WithPath("record_id", suite.Account.RecordID).
		WithBytes(jsondata).
		Expect().Body().Contains("ok").Contains("true").
		Contains("chinese_medicine_advice").Contains("dietary_advice").Contains("physical_therapy_advice").Contains("sports_advice")
}

func TestAnalysisV2TestSuite(t *testing.T) {
	suite.Run(t, new(AnalysisV2TestSuite))
}
