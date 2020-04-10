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

// ResetPasswordSuite 是ResetPassword的单元测试的 Test Suite
type ResetPasswordSuite struct {
	suite.Suite
	Expect  *httpexpect.Expect
	Account *Account
}

// SetupSuite 设置测试环境
func (suite *ResetPasswordSuite) SetupSuite() {
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

// TestResetPassword 测试重置密码
func (suite *ResetPasswordSuite) TestResetPassword() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, userID, errGetAccessToken := getAccessTokenAndUserID(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	body := e.POST("/notification/sms").WithHeader("Authorization", auth).WithJSON(
		&r.SmsNotificationBody{
			Phone:            suite.Account.SignInPhone,
			NotificationType: r.ResetPasswordSmsNotification,
			Language:         suite.Account.Language,
			NationCode:       suite.Account.NationCode,
			UserID:           userID,
		}).Expect().Body()
	var respSms SmsNotification
	errUnmarshalSms := json.Unmarshal([]byte(body.Raw()), &respSms)
	assert.NoError(t, errUnmarshalSms)
	// 获取最新的mvc
	mvc, errGetLatestMVC := getLatestMVC(e)
	assert.NoError(t, errGetLatestMVC)
	bodyVN := e.POST("/validate_signin_phone").
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
	var respVN VerifyUserSigninPhoneResp
	errUnmarshalVN := json.Unmarshal([]byte(bodyVN.Raw()), &respVN)
	assert.NoError(t, errUnmarshalVN)
	// 重置密码
	e.POST("/user/{user_id}/reset_password").
		WithHeaders(headers).
		WithPath("user_id", userID).WithJSON(
		&r.ResetPasswordBody{
			PlainPassword:      suite.Account.PlainPassword,
			VerificationNumber: respVN.Data.VerificationNumber,
			VerificationType:   r.PhoneVerificationType,
		}).
		Expect().Body().Contains("新密码不能与原密码相同")
}

func TestResetPasswordSuite(t *testing.T) {
	suite.Run(t, new(ResetPasswordSuite))
}
