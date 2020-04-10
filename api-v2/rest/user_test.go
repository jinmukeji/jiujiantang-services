package rest_test

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"testing"

	"time"

	"github.com/iris-contrib/httpexpect"
	r "github.com/jinmukeji/jiujiantang-services/api-v2/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// UserTestSuite 是User的单元测试的 Test Suite
type UserTestSuite struct {
	suite.Suite
	Account *Account
	Expect  *httpexpect.Expect
}

// SetupSuite 设置测试环境
func (suite *UserTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.api-v2.env")
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
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

// UserSignInResponse 登录的返回
type UserSignInResponse struct {
	Data r.UserSignInResponse `json:"data"`
}

// TestUserSignIn 测试登录
func (suite *UserTestSuite) TestUserSignIn() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	body := e.POST("/v2-api/users/signin").WithHeader("Authorization", auth).WithJSON(&r.UserSignIn{
		SignInKey:    suite.Account.SignInKey,
		PasswordHash: suite.Account.PasswordHash,
	},
	).Expect().Body()
	var resp UserSignInResponse
	errUnmarshalSignIn := json.Unmarshal([]byte(body.Raw()), &resp)
	assert.NoError(t, errUnmarshalSignIn)
	assert.NotNil(t, resp.Data.AccessToken)
	assert.NotEqual(t, "", resp.Data.AccessToken)
}

// TestGetUserProfile 测试得到user_profile
func (suite *UserTestSuite) TestGetUserProfile() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, errGetAccessToken := getAccessToken(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	e.GET("/v2-api/owner/users/{user_id}/profile").WithPath("user_id", suite.Account.UserID).WithHeaders(headers).Expect().Body().Contains(suite.Account.Nickname)
}

// TestModifyUserProfile 测试修改UserProfile
func (suite *UserTestSuite) TestModifyUserProfile() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, errGetAccessToken := getAccessToken(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	year1993 := time.Date(1993, 1, 1, 0, 0, 0, 0, time.UTC)
	var gender = int32(r.GenderMale)
	var height = int32(170)
	var weight = int32(70)
	e.PUT("/v2-api/owner/users/{user_id}/profile").WithPath("user_id", suite.Account.UserID).WithHeaders(headers).WithJSON(&r.ModifyUserProfile{
		UserProfile: r.UserProfile{
			Nickname: suite.Account.Nickname,
			Birthday: year1993,
			Gender:   &gender,
			Height:   &height,
			Weight:   &weight,
		},
	},
	).Expect().Body().Contains(suite.Account.Nickname)
}

// TestUserSignUp 注册
func (suite *UserTestSuite) TestUserSignUp() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, errGetAccessToken := getAccessToken(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	var nickname = "liu"
	year1993 := time.Date(1993, 1, 1, 0, 0, 0, 0, time.UTC)
	var gender = int32(r.GenderMale)
	var height = int32(170)
	var weight = int32(70)
	envFilepath := filepath.Join("testdata", "local.api-v2.env")
	client := newTestingClientFromEnvFile(envFilepath)
	e.POST("/v2-api/owner/users/signup").WithHeaders(headers).WithJSON(&r.UserSignUp{
		UserProfile: r.UserProfile{
			Nickname: nickname,
			Birthday: year1993,
			Gender:   &gender,
			Height:   &height,
			Weight:   &weight,
		},
		RegisterType: suite.Account.RegisterType,
		Username:     suite.Account.SignInKey,
		ClientID:     client.ClientID,
	},
	).Expect().Body().Contains(nickname)
}

// TestOwnerUserSignUp 测试OwnerUserSignUp
func (suite *UserTestSuite) TestOwnerUserSignUp() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, errGetAccessToken := getAccessToken(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	var nickname = "liu"
	year1993 := time.Date(1993, 1, 1, 0, 0, 0, 0, time.UTC)
	var gender = int32(r.GenderMale)
	var height = int32(170)
	var weight = int32(70)
	envFilepath := filepath.Join("testdata", "local.api-v2.env")
	client := newTestingClientFromEnvFile(envFilepath)
	e.POST("/v2-api/owner/{owner_id}/users/sign_up").WithHeaders(headers).WithJSON(&r.OwnerUserSignUpBody{
		Nickname:     nickname,
		Birthday:     year1993,
		Gender:       &gender,
		Height:       height,
		Weight:       weight,
		RegisterType: suite.Account.RegisterType,
		ClientID:     client.ClientID,
	}).
		WithPath("owner_id", suite.Account.UserID).
		Expect().Body().Contains(nickname)
}

// TestSignOut 测试登出
func (suite *UserTestSuite) TestSignOut() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, errGetAccessToken := getAccessToken(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	e.POST("/v2-api/users/signout").WithHeaders(headers).Expect().Body().Contains("ok").Contains("true")
}

// getAccessToken 得到AccessToken
func getAccessToken(e *httpexpect.Expect) (string, error) {
	envFilepath := filepath.Join("testdata", "local.api-v2.env")
	account := newTestingAccountFromEnvFile(envFilepath)
	auth, errGetAuthorization := getAuthorization(e)
	if errGetAuthorization != nil {
		return "", errGetAuthorization
	}
	body := e.POST("/v2-api/users/signin").WithHeader("Authorization", auth).WithJSON(&r.UserSignIn{
		SignInKey:    account.SignInKey,
		PasswordHash: account.PasswordHash,
	},
	).Expect().Body()
	var resp UserSignInResponse
	errUnmarshalSignIn := json.Unmarshal([]byte(body.Raw()), &resp)
	return resp.Data.AccessToken, errUnmarshalSignIn
}

// getAuthorization 得到Authorization
func getAuthorization(e *httpexpect.Expect) (string, error) {
	envFilepath := filepath.Join("testdata", "local.api-v2.env")
	client := newTestingClientFromEnvFile(envFilepath)
	body := e.POST("/v2-api/client/auth").WithJSON(&r.ClientAuthReq{
		ClientID:      client.ClientID,
		SecretKeyHash: client.SecretKeyHash,
		Seed:          client.Seed,
	},
	).Expect().Body()
	var auth ClientAuth
	errUnmarshalAuth := json.Unmarshal([]byte(body.Raw()), &auth)
	return auth.Data.Authorization, errUnmarshalAuth
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
