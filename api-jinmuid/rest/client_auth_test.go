package rest_test

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/iris-contrib/httpexpect"
	r "github.com/jinmukeji/jiujiantang-services/api-jinmuid/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ClientAuthTestSuite 是ClientAuth的单元测试的 Test Suite
type ClientAuthTestSuite struct {
	suite.Suite
	Client *Client
}

// SetupSuite 设置测试环境
func (suite *ClientAuthTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.api-jinmuid.env")
	suite.Client = newTestingClientFromEnvFile(envFilepath)
}

// ClientAuth 授权的返回
type ClientAuth struct {
	Data r.ClientAuth `json:"data"`
}

// TestClientAuth 测试授权
func (suite *ClientAuthTestSuite) TestClientAuth() {
	t := suite.T()
	app := r.NewApp("", "jinmuhealth", false)
	e := httpexpect.WithConfig(httpexpect.Config{
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
	body := e.POST("/client/auth").WithJSON(&r.ClientAuthReq{
		ClientID:      suite.Client.ClientID,
		SecretKeyHash: suite.Client.SecretKeyHash,
		Seed:          suite.Client.Seed,
	},
	).Expect().Body()
	var auth ClientAuth
	errUnmarshal := json.Unmarshal([]byte(body.Raw()), &auth)
	assert.NoError(t, errUnmarshal)
	assert.NotNil(t, auth.Data.Authorization)
}

func TestClientAuthTestSuite(t *testing.T) {
	suite.Run(t, new(ClientAuthTestSuite))
}
