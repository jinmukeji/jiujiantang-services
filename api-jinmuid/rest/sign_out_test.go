package rest_test

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
	r "github.com/jinmukeji/jiujiantang-services/api-jinmuid/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// SignOutSuite 是SignOut的单元测试的 Test Suite
type SignOutSuite struct {
	suite.Suite
	Expect *httpexpect.Expect
}

// SetupSuite 设置测试环境
func (suite *SignOutSuite) SetupSuite() {
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

// TestSignOut 测试登出
func (suite *SignOutSuite) TestSignOut() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, _, errGetAccessToken := getAccessTokenAndUserID(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	e.POST("/user/signout").
		WithHeaders(headers).Expect().Body().Contains("ok").Contains("true")
}

func TestSignOutSuite(t *testing.T) {
	suite.Run(t, new(SignOutSuite))
}
