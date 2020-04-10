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

// SigninPhoneSuite 是SigninPhone的单元测试的 Test Suite
type SigninPhoneSuite struct {
	suite.Suite
	Expect  *httpexpect.Expect
	Account *Account
}

// SetupSuite 设置测试环境
func (suite *SigninPhoneSuite) SetupSuite() {
	t := suite.T()
	app := r.NewApp("", "jinmuhealth", true)
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

// TestSetSigninPhone 测试设置登录手机号
func (suite *SigninPhoneSuite) TestSetSigninPhone() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, userID, errGetAccessToken := getAccessTokenAndUserID(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	// 发送设置手机号短信
	bodySms := e.POST("/notification/sms").WithHeader("Authorization", auth).WithJSON(
		&r.SmsNotificationBody{
			Phone:            suite.Account.SignInPhone,
			NotificationType: r.SetPhoneSmsNotification,
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
	// 设置手机号
	e.POST("/user/{user_id}/signin_phone").WithPath("user_id", userID).WithHeaders(headers).WithJSON(
		&r.SetSigninPhoneBody{
			Phone:               suite.Account.SignInPhone,
			NationCode:          suite.Account.NationCode,
			Mvc:                 mvc,
			SerialNumber:        respSms.Data.SerialNumber,
			SmsNotificationType: r.SetPhoneSmsNotification,
		},
	).Expect().Body().Contains("手机号已被注册")
}

// VerifyUserSigninPhoneResp 验证登录手机号的返回
type VerifyUserSigninPhoneResp struct {
	Data r.VerifyUserSigninPhoneResp `json:"data"`
}

// TestVerifySigninPhone 测试VerifySigninPhone
func (suite *SigninPhoneSuite) TestVerifySigninPhone() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, userID, errGetAccessToken := getAccessTokenAndUserID(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	// 发送重置密码的短信
	bodySms := e.POST("/notification/sms").WithHeader("Authorization", auth).WithJSON(
		&r.SmsNotificationBody{
			Phone:            suite.Account.SignInPhone,
			NotificationType: r.ResetPasswordSmsNotification,
			Language:         suite.Account.Language,
			NationCode:       suite.Account.NationCode,
			UserID:           userID,
		}).Expect().Body()
	var respSms SmsNotification
	errUnmarshalSms := json.Unmarshal([]byte(bodySms.Raw()), &respSms)
	assert.NoError(t, errUnmarshalSms)
	// 获取最新的mvc
	mvc, errGetLatestMVC := getLatestMVC(e)
	assert.NoError(t, errGetLatestMVC)
	// 验证登录手机号
	body := e.POST("/validate_signin_phone").
		WithHeaders(headers).
		WithJSON(
			&r.VerifySigninPhoneBody{
				r.SetSigninPhoneBody{
					Phone:               suite.Account.SignInPhone,
					Mvc:                 mvc,
					SerialNumber:        respSms.Data.SerialNumber,
					NationCode:          suite.Account.NationCode,
					SmsNotificationType: r.ResetPasswordSmsNotification,
				},
			}).
		Expect().Body()
	var resp VerifyUserSigninPhoneResp
	errUnmarshal := json.Unmarshal([]byte(body.Raw()), &resp)
	assert.NoError(t, errUnmarshal)
	assert.NotNil(t, resp.Data.VerificationNumber)
	assert.NotEqual(t, "", resp.Data.VerificationNumber)
}

func TestSigninPhoneSuite(t *testing.T) {
	suite.Run(t, new(SigninPhoneSuite))
}
