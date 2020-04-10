package rest_test

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/iris-contrib/httpexpect"
	r "github.com/jinmukeji/jiujiantang-services/api-jinmuid/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// UserSuite 是User的单元测试的 Test Suite
type UserSuite struct {
	suite.Suite
	Expect  *httpexpect.Expect
	Account *Account
}

// SetupSuite 设置测试环境
func (suite *UserSuite) SetupSuite() {
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

// TestSetUserPassword 测试设置密码
func (suite *UserSuite) TestSetUserPassword() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, userID, errGetAccessToken := getAccessTokenAndUserID(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	e.POST("/user/{user_id}/password").
		WithHeaders(headers).
		WithPath("user_id", userID).WithJSON(
		&r.SetUserPasswordBody{
			PlainPassword: suite.Account.PlainPassword,
		},
	).Expect().Body().Contains("密码已经设置")
}

// TestModifyUserPassword 测试修改用户名密码
func (suite *UserSuite) TestModifyUserPassword() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, userID, errGetAccessToken := getAccessTokenAndUserID(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	e.PUT("/user/{user_id}/password").
		WithHeaders(headers).
		WithPath("user_id", userID).WithJSON(
		&r.ModifyUserPasswordBody{
			OldHashedPassword: suite.Account.HashedPassword,
			Seed:              suite.Account.Seed,
			NewPlainPassword:  suite.Account.PlainPassword,
		},
	).Expect().Body().Contains("新密码不能与原密码相同")
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}
