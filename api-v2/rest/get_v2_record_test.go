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

// GetV2RecordTestSuite 是GetV2Record的单元测试的 Test Suite
type GetV2RecordTestSuite struct {
	suite.Suite
	Expect  *httpexpect.Expect
	Account *Account
}

// SetupSuite 设置测试环境
func (suite *GetV2RecordTestSuite) SetupSuite() {
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

// TestGetVersion 版本
func (suite *GetV2RecordTestSuite) TestGetV2AnalyzeReportByRecordID() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, errGetAccessToken := getAccessToken(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	e.GET("/v2-api/owner/measurements/{record_id}/v2/analyze").
		WithHeaders(headers).
		WithPath("record_id", suite.Account.RecordID).Expect().Body().
		Contains("c0").Contains("c1").Contains("c2").Contains("c3").Contains("c4").Contains("c5").Contains("c6").Contains("c7")
}

func TestGetV2RecordTestSuite(t *testing.T) {
	suite.Run(t, new(GetV2RecordTestSuite))
}
