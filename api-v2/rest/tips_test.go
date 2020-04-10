package rest_test

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
	r "github.com/jinmukeji/jiujiantang-services/api-v2/rest"
	"github.com/stretchr/testify/suite"
)

// TipsTestSuite 是tips的单元测试的 Test Suite
type TipsTestSuite struct {
	suite.Suite
	Expect *httpexpect.Expect
}

// SetupSuite 设置测试环境
func (suite *TipsTestSuite) SetupSuite() {
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

// TestGetTips tips
func (suite *TipsTestSuite) TestGetTips() {
	suite.Expect.GET("/v2-api/tips").Expect().Body().Contains("content").Contains("duration")
}

func TestTipsTestSuite(t *testing.T) {
	suite.Run(t, new(TipsTestSuite))
}
