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

// AnalysisReportRemarkTestSuite 是AnalysisReportRemark的单元测试的 Test Suite
type AnalysisReportRemarkTestSuite struct {
	suite.Suite
	Expect  *httpexpect.Expect
	Account *Account
}

// SetupSuite 设置测试环境
func (suite *AnalysisReportRemarkTestSuite) SetupSuite() {
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

// TestSubmitRemark 测试SubmitRemark
func (suite *AnalysisReportRemarkTestSuite) TestSubmitRemark() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, errGetAccessToken := getAccessToken(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	e.PUT("/v2-api/owner/measurements/{record_id}/remark").
		WithHeaders(headers).
		WithPath("record_id", suite.Account.RecordID).
		WithJSON(&r.SubmitRemarkReq{
			UserID: suite.Account.UserID,
			Remark: suite.Account.Remark,
		}).
		Expect().Body().Contains("ok").Contains("true")
}

func TestAnalysisReportRemarkTestSuite(t *testing.T) {
	suite.Run(t, new(AnalysisReportRemarkTestSuite))
}
