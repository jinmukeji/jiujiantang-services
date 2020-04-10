package rest_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/iris-contrib/httpexpect"
	r "github.com/jinmukeji/gf-api2/api-v2/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// SearchHistoryTestSuite 是SearchHistory的单元测试的 Test Suite
type SearchHistoryTestSuite struct {
	suite.Suite
	Account *Account
	Expect  *httpexpect.Expect
}

// SetupSuite 设置测试环境
func (suite *SearchHistoryTestSuite) SetupSuite() {
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

// TestSearchHistory 测试历史记录
func (suite *SearchHistoryTestSuite) TestSearchHistory() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, errGetAccessToken := getAccessToken(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	e.GET("/v2-api/owner/measurements").
		WithHeaders(headers).
		WithQuery("offset", 0).
		WithQuery("size", 20).
		WithQuery("user_id", suite.Account.UserID).
		WithQuery("start", "2017-07-19T07:31:20Z").
		WithQuery("end", "2019-07-19T07:31:20Z").
		Expect().Body().Contains("c0").Contains("c1").Contains("c2").Contains("c3").Contains("c4").Contains("c5").Contains("c6").Contains("c7")
}

// SubmitMeasurementData 测量的返回
type SubmitMeasurementData struct {
	Data r.SubmitMeasurementData `json:"data"`
}

// TestDeleteRecords 测试删除记录
func (suite *SearchHistoryTestSuite) TestDeleteRecords() {
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
	// 添加一条记录
	body := e.POST("/v2-api/owner/measurements").
		WithHeaders(headers).
		WithBytes(jsondata).
		Expect().Body()
	var resp SubmitMeasurementData
	errUnmarshalSignIn := json.Unmarshal([]byte(body.Raw()), &resp)
	assert.NoError(t, errUnmarshalSignIn)
	// 删除记录
	e.POST("/v2-api/user/measurements/{user_id}/delete").
		WithHeaders(headers).
		WithPath("user_id", suite.Account.UserID).
		WithJSON(&r.DeleteRecordsBody{
			RecordIDList: []int{int(resp.Data.RecordID)},
		}).
		Expect().Body().Contains("ok").Contains("true")
}

func TestSearchHistoryTestSuite(t *testing.T) {
	suite.Run(t, new(SearchHistoryTestSuite))
}
