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

// GetCurrentSecureQuestionSuite 是GetCurrentSecureQuestion的单元测试的 Test Suite
type GetCurrentSecureQuestionSuite struct {
	suite.Suite
	Expect  *httpexpect.Expect
	Account *Account
}

// SetupSuite 设置测试环境
func (suite *GetCurrentSecureQuestionSuite) SetupSuite() {
	t := suite.T()
	app := r.NewApp("", "jinmuhealth", false)
	envFilepath := filepath.Join("testdata", "local.api-jinmuid.env")
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

// TestGetSecureQuestionsByPhoneOrUsername 测试GetSecureQuestionsByPhoneOrUsername
func (suite *GetCurrentSecureQuestionSuite) TestGetSecureQuestionsByPhoneOrUsername() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, _, errGetAccessToken := getAccessTokenAndUserID(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	e.POST("/user/secure_question").
		WithHeaders(headers).
		WithJSON(
			&r.GetSecureQuestionsByPhoneOrUsernameRequest{
				ValidationType: r.ValidationTypePhone,
				Phone:          suite.Account.SignInPhone,
				NationCode:     suite.Account.NationCode,
			},
		).
		Expect().Body().Contains("未设置密保问题")
}

func TestGetCurrentSecureQuestionSuite(t *testing.T) {
	suite.Run(t, new(GetCurrentSecureQuestionSuite))
}
