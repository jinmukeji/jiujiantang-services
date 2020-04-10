package rest_test

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/iris-contrib/httpexpect"
	r "github.com/jinmukeji/gf-api2/api-jinmuid/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// DeviceSuite 是Device的单元测试的 Test Suite
type DeviceSuite struct {
	suite.Suite
	Expect  *httpexpect.Expect
	Account *Account
}

// SetupSuite 设置测试环境
func (suite *DeviceSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.api-jinmuid.env")
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
	t := suite.T()
	app := r.NewApp("", "jinmuhealth", false)
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

// TestUserUsedDevice 测试UserUsedDevice
func (suite *DeviceSuite) TestUserUsedDevice() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, userID, errGetAccessToken := getAccessTokenAndUserID(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	e.GET("/user/{user_id}/devices").
		WithHeaders(headers).
		WithPath("user_id", userID).Expect().Body().Contains("ok").Contains("true")
}

func TestDeviceSuite(t *testing.T) {
	suite.Run(t, new(DeviceSuite))
}
