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

// ValidatePhoneVerificationCodeSuite 是ValidatePhoneVerificationCode的单元测试的 Test Suite
type ValidatePhoneVerificationCodeSuite struct {
	suite.Suite
	Expect  *httpexpect.Expect
	Account *Account
}

// SetupSuite 设置测试环境
func (suite *ValidatePhoneVerificationCodeSuite) SetupSuite() {
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

// TestValidatePhoneVerificationCode 测试ValidatePhoneVerificationCode
func (suite *ValidatePhoneVerificationCodeSuite) TestValidatePhoneVerificationCode() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, userID, errGetAccessToken := getAccessTokenAndUserID(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	bodySms := e.POST("/notification/sms").WithHeader("Authorization", auth).WithJSON(
		&r.SmsNotificationBody{
			Phone:            suite.Account.SignInPhone,
			NotificationType: r.SignupSmsNotification,
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
	e.POST("/user/validate_phone_verification_code").
		WithHeaders(headers).
		WithJSON(
			&r.ValidatePhoneVerificationCodeBody{
				Phone:        suite.Account.SignInPhone,
				Mvc:          mvc,
				SerialNumber: respSms.Data.SerialNumber,
				NationCode:   suite.Account.NationCode,
			}).
		Expect().Body().Contains("手机号已被注册")
}

func TestValidatePhoneVerificationCodeSuite(t *testing.T) {
	suite.Run(t, new(ValidatePhoneVerificationCodeSuite))
}
