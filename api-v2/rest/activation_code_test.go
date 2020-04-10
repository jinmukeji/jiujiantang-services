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

// ActivationCodeTestSuite 是ActivationCode的单元测试的 Test Suite
type ActivationCodeTestSuite struct {
	suite.Suite
	Expect  *httpexpect.Expect
	Account *Account
}

// SetupSuite 设置测试环境
func (suite *ActivationCodeTestSuite) SetupSuite() {
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

// TestGetActivationCodeInfo 测试GetActivationCodeInfo
func (suite *ActivationCodeTestSuite) TestGetActivationCodeInfo() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, errGetAccessToken := getAccessToken(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	e.POST("/v2-api/activation_code").
		WithHeaders(headers).
		WithJSON(&r.ActivationCode{
			Code: suite.Account.ActivationCode,
		}).
		Expect().Body().Contains("ok").Contains(suite.Account.ActivationCodeSubscriptionType)
}

// TestUseSubscriptionActivationCode 测试UseSubscriptionActivationCode
func (suite *ActivationCodeTestSuite) TestUseSubscriptionActivationCode() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, errGetAccessToken := getAccessToken(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	e.POST("/v2-api/owner/user/{user_id}/activation_code").
		WithHeaders(headers).
		WithPath("user_id", suite.Account.UserID).
		WithJSON(&r.ActivationCode{
			Code: suite.Account.UsedActivationCode,
		}).
		Expect().Body().Contains("ok").Contains("激活码已失效")
}

func TestActivationCodeTestSuite(t *testing.T) {
	suite.Run(t, new(ActivationCodeTestSuite))
}
