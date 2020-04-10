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

// NotificationPreferencesSuite 是NotificationPreferences的单元测试的 Test Suite
type NotificationPreferencesSuite struct {
	suite.Suite
	Expect  *httpexpect.Expect
	Account *Account
}

// SetupSuite 设置测试环境
func (suite *NotificationPreferencesSuite) SetupSuite() {
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

// TestGetNotificationPreferences 测试GetNotificationPreferences
func (suite *NotificationPreferencesSuite) TestGetNotificationPreferences() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, userID, errGetAccessToken := getAccessTokenAndUserID(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	e.GET("/user/{user_id}/notification_preferences").
		WithHeaders(headers).
		WithPath("user_id", userID).
		Expect().Body().
		Contains("phone_enabled").
		Contains("wechat_enabled").
		Contains("weibo_enabled")
}

// TestModifyNotificationPreferences 测试ModifyNotificationPreferences
func (suite *NotificationPreferencesSuite) TestModifyNotificationPreferences() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, userID, errGetAccessToken := getAccessTokenAndUserID(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	e.POST("/user/{user_id}/notification_preferences").
		WithHeaders(headers).
		WithPath("user_id", userID).
		WithJSON(r.ModifyNotificationPreferencesBody{
			ModifyNotificationPreferences: r.ModifyNotificationPreferences{
				PhoneEnabled:  true,
				WechatEnabled: true,
				WeiboEnabled:  true,
			},
		}).
		Expect().Body().
		Contains("ok").
		Contains("true")
}

func TestNotificationPreferencesSuite(t *testing.T) {
	suite.Run(t, new(NotificationPreferencesSuite))
}
