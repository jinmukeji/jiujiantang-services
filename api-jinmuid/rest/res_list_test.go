package rest_test

import (
	"net/http"
	"testing"

	"github.com/iris-contrib/httpexpect"
	r "github.com/jinmukeji/gf-api2/api-jinmuid/rest"
	"github.com/stretchr/testify/suite"
)

// ResListTestSuite 是ResList的单元测试的 Test Suite
type ResListTestSuite struct {
	suite.Suite
	Expect  *httpexpect.Expect
	Account *Account
}

// SetupSuite 设置测试环境
func (suite *ResListTestSuite) SetupSuite() {
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

// TestGetResourceList 测试GetResourceList
func (suite *ResListTestSuite) TestGetResourceList() {
	e := suite.Expect
	e.GET("/resource").
		Expect().Body().
		Contains("zh-Hans").Contains("zh-Hant").Contains("en").Contains("mainland_china").Contains("taiwan").Contains("abroad")
}

func TestResListTestSuite(t *testing.T) {
	suite.Run(t, new(ResListTestSuite))
}
