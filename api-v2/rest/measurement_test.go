package rest_test

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/iris-contrib/httpexpect"
	r "github.com/jinmukeji/jiujiantang-services/api-v2/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// MeasurementTestSuite 是measurement的单元测试的 Test Suite
type MeasurementTestSuite struct {
	suite.Suite
	Expect *httpexpect.Expect
}

// SetupSuite 设置测试环境
func (suite *MeasurementTestSuite) SetupSuite() {
	t := suite.T()
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

// TestSubmitMeasurementData 测试SubmitMeasurementData
func (suite *MeasurementTestSuite) TestSubmitMeasurementData() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, errGetAccessToken := getAccessToken(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	headers["Content-Type"] = "application/json"
	envFilepath := filepath.Join("testdata", "algorithmServerTest.json")
	jsondata, _ := ioutil.ReadFile(envFilepath)
	e.POST("/v2-api/owner/measurements").
		WithHeaders(headers).
		WithBytes(jsondata).
		Expect().Body().
		Contains("c0").Contains("c1").Contains("c2").Contains("c3").Contains("c4").Contains("c5").Contains("c6").Contains("c7").Contains("record_id")
}

func TestMeasurementTestSuite(t *testing.T) {
	suite.Run(t, new(MeasurementTestSuite))
}
