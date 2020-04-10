package rest_test

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/iris-contrib/httpexpect"
	r "github.com/jinmukeji/gf-api2/api-jinmuid/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// SignInSuite 是SignIn的单元测试的 Test Suite
type SignInSuite struct {
	suite.Suite
	Expect  *httpexpect.Expect
	Account *Account
}

// SetupSuite 设置测试环境
func (suite *SignInSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.api-jinmuid.env")
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
	t := suite.T()
	app := r.NewApp("", "jinmuhealth", true)
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

// SignInResp 登录的返回
type SignInResp struct {
	Data r.SignInResp `json:"data"`
}

// TestSignInByUsernamePassword 测试用户名密码登录
func (suite *SignInSuite) TestSignInByUsernamePassword() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	body := e.POST("/signin").WithHeader("Authorization", auth).WithJSON(&r.SignInBody{
		SignInMethod:   "username_password",
		Username:       suite.Account.Username,
		HashedPassword: suite.Account.HashedPassword,
		Seed:           suite.Account.Seed,
		SignInMachine:  suite.Account.SignInMachine,
	},
	).Expect().Body()
	var resp SignInResp
	errUnmarshalSignIn := json.Unmarshal([]byte(body.Raw()), &resp)
	assert.NoError(t, errUnmarshalSignIn)
	assert.NotNil(t, resp.Data.AccessToken)
	assert.NotEqual(t, "", resp.Data.AccessToken)
}

// TestSignInByPhonePassword  测试手机号密码登录
func (suite *SignInSuite) TestSignInByPhonePassword() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	body := e.POST("/signin").WithHeader("Authorization", auth).WithJSON(&r.SignInBody{
		SignInMethod:   "phone_password",
		Phone:          suite.Account.SignInPhone,
		HashedPassword: suite.Account.HashedPassword,
		Seed:           suite.Account.Seed,
		SignInMachine:  suite.Account.SignInMachine,
		NationCode:     suite.Account.NationCode,
	},
	).Expect().Body()
	var resp SignInResp
	errUnmarshalSignIn := json.Unmarshal([]byte(body.Raw()), &resp)
	assert.NoError(t, errUnmarshalSignIn)
	assert.NotNil(t, resp.Data.AccessToken)
	assert.NotEqual(t, "", resp.Data.AccessToken)
}

// TestSignInByPhoneMVC 测试手机号验证码登录
func (suite *SignInSuite) TestSignInByPhoneMVC() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	// 发送短信
	bodySms := e.POST("/notification/sms").WithHeader("Authorization", auth).WithJSON(
		&r.SmsNotificationBody{
			Phone:            suite.Account.SignInPhone,
			NotificationType: r.SigninSmsNotification,
			Language:         suite.Account.Language,
			NationCode:       suite.Account.NationCode,
		}).Expect().Body()
	var respSms SmsNotification
	errUnmarshalSms := json.Unmarshal([]byte(bodySms.Raw()), &respSms)
	assert.NoError(t, errUnmarshalSms)
	assert.NotNil(t, respSms.Data.SerialNumber)
	assert.NotEqual(t, "", respSms.Data.SerialNumber)
	// 获取最新的mvc
	mvc, errGetLatestMVC := getLatestMVC(e)
	assert.NoError(t, errGetLatestMVC)
	body := e.POST("/signin").WithHeader("Authorization", auth).WithJSON(&r.SignInBody{
		SignInMethod:  "phone_mvc",
		Phone:         suite.Account.SignInPhone,
		MVC:           mvc,
		SerialNumber:  respSms.Data.SerialNumber,
		SignInMachine: suite.Account.SignInMachine,
		NationCode:    suite.Account.NationCode,
	},
	).Expect().Body()
	var resp SignInResp
	errUnmarshalSignIn := json.Unmarshal([]byte(body.Raw()), &resp)
	assert.NoError(t, errUnmarshalSignIn)
	assert.NotNil(t, resp.Data.AccessToken)
	assert.NotEqual(t, "", resp.Data.AccessToken)
}

// getLatestMVC 获取最新的mvc
func getLatestMVC(e *httpexpect.Expect) (string, error) {
	envFilepath := filepath.Join("testdata", "local.api-jinmuid.env")
	account := newTestingAccountFromEnvFile(envFilepath)
	auth, _ := getAuthorization(e)
	sendInformation := make([]r.GetLatestVerificationCode, 1)
	sendInformation[0] = r.GetLatestVerificationCode{
		SendVia:    r.SendViaPhone,
		Phone:      account.SignInPhone,
		NationCode: account.NationCode,
	}
	body := e.POST("/_debug/user/latest_verification_code").WithHeader("Authorization", auth).WithJSON(&r.LatestVerificationCodeBody{
		SendInformation: sendInformation,
	},
	).Expect().Body()
	var resp LatestVerificationCodes
	errUnmarshalSignIn := json.Unmarshal([]byte(body.Raw()), &resp)
	return resp.Data[0].VerificationCode, errUnmarshalSignIn
}

// getAccessTokenAndUserID 得到AccessToken和userID
func getAccessTokenAndUserID(e *httpexpect.Expect) (string, int32, error) {
	envFilepath := filepath.Join("testdata", "local.api-jinmuid.env")
	account := newTestingAccountFromEnvFile(envFilepath)
	auth, _ := getAuthorization(e)
	body := e.POST("/signin").WithHeader("Authorization", auth).WithJSON(&r.SignInBody{
		SignInMethod:   "phone_password",
		Phone:          account.SignInPhone,
		HashedPassword: account.HashedPassword,
		Seed:           account.Seed,
		SignInMachine:  account.SignInMachine,
		NationCode:     account.NationCode,
	},
	).Expect().Body()
	var resp SignInResp
	errUnmarshalSignIn := json.Unmarshal([]byte(body.Raw()), &resp)
	return resp.Data.AccessToken, resp.Data.UserID, errUnmarshalSignIn
}

// getAuthorization 得到Authorization
func getAuthorization(e *httpexpect.Expect) (string, error) {
	envFilepath := filepath.Join("testdata", "local.api-jinmuid.env")
	client := newTestingClientFromEnvFile(envFilepath)
	body := e.POST("/client/auth").WithJSON(&r.ClientAuthReq{
		ClientID:      client.ClientID,
		SecretKeyHash: client.SecretKeyHash,
		Seed:          client.Seed,
	},
	).Expect().Body()
	var auth ClientAuth
	errUnmarshalAuth := json.Unmarshal([]byte(body.Raw()), &auth)
	return auth.Data.Authorization, errUnmarshalAuth
}

func TestSignInSuite(t *testing.T) {
	suite.Run(t, new(SignInSuite))
}
