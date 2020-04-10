package rest_test

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
	r "github.com/jinmukeji/gf-api2/api-v2/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// LocalNamePrefixTestSuitet 是LocalNamePrefix的单元测试的 Test Suite
type LocalNamePrefixTestSuite struct {
	suite.Suite
	Expect *httpexpect.Expect
}

// SetupSuite 设置测试环境
func (suite *LocalNamePrefixTestSuite) SetupSuite() {
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

// TestGetBluetoothNamePrefixes 测试得到蓝牙名前缀
func (suite *LocalNamePrefixTestSuite) TestGetBluetoothNamePrefixes() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, errGetAccessToken := getAccessToken(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	e.GET("/v2-api/bluetooth_name_prefix").
		WithHeaders(headers).
		Expect().Body().Contains("JinMu").Contains("HJT").Contains("KM")
}

func TestLocalNamePrefixTestSuite(t *testing.T) {
	suite.Run(t, new(LocalNamePrefixTestSuite))
}
