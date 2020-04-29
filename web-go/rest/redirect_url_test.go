package rest

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
	"github.com/stretchr/testify/suite"
)

// RedirectURLTestSuite 是redirect_url的单元测试的 Test Suite
type RedirectURLTestSuite struct {
	suite.Suite
}

// SetupSuite 设置测试环境
func (suite *RedirectURLTestSuite) SetupSuite() {
}

// TestRedirectResURL 测试RedirectResURL
func (suite *RedirectURLTestSuite) TestRedirectResURL() {
	t := suite.T()
	app := NewApp("../data/resource.yml")
	e := httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(app),
			Jar:       httpexpect.NewJar(),
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewCurlPrinter(t),
			httpexpect.NewDebugPrinter(t, true),
		},
	})
	category, key := "e", "c7"
	e.GET(fmt.Sprintf("/%s/%s", category, key)).WithQuery("env", "production").Expect().Body().Contains("Found").Contains("https://res-cdn.jinmuhealth.com/v2/app-entry/2-1/cc/c7.html")
}

func TestRedirectURLTestSuite(t *testing.T) {
	suite.Run(t, new(RedirectURLTestSuite))
}
