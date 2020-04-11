package rest_test

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
	r "github.com/jinmukeji/jiujiantang-services/api-jinmuid/rest"
	"github.com/stretchr/testify/suite"
)

// VersionTestSuite 是version的单元测试的 Test Suite
type VersionTestSuite struct {
	suite.Suite
	Expect *httpexpect.Expect
}

// SetupSuite 设置测试环境
func (suite *VersionTestSuite) SetupSuite() {
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

// TestGetVersion 版本
func (suite *VersionTestSuite) TestGetVersion() {
	suite.Expect.GET("/version").Expect().Body().Contains("com.himalife.srv.svc-jinmuid")
}

func TestVersionTestSuite(t *testing.T) {
	suite.Run(t, new(VersionTestSuite))
}
