package rest_test

import (
	"net/http"
	"path/filepath"

	"encoding/json"
	"testing"

	"github.com/iris-contrib/httpexpect"
	r "github.com/jinmukeji/gf-api2/api-jinmuid/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// FindUsernameByEmailTestSuite 是FindUsernameByEmail的单元测试的 Test Suite
type FindUsernameByEmailTestSuite struct {
	suite.Suite
	Expect  *httpexpect.Expect
	Account *Account
}

// SetupSuite 设置测试环境
func (suite *FindUsernameByEmailTestSuite) SetupSuite() {
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

// TestFindUsernameBySecureEmail 测试FindUsernameBySecureEmail
func (suite *FindUsernameByEmailTestSuite) TestFindUsernameBySecureEmail() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	// 发送邮件
	bodyEmailNotification := e.POST("/notification/email").WithHeader("Authorization", auth).WithJSON(
		&r.EmailNotificationRequest{
			Email:    suite.Account.Email,
			Type:     r.FindUsernameSemNotification,
			Language: suite.Account.Language,
		},
	).Expect().Body()
	var resp EmailNotificationResponse
	errUnmarshalEmailNotification := json.Unmarshal([]byte(bodyEmailNotification.Raw()), &resp)
	assert.NoError(t, errUnmarshalEmailNotification)
	// 获取最新的evc
	evc, errGetLatestEVC := getLatestEVC(e)
	assert.NoError(t, errGetLatestEVC)
	// 找回用户名
	e.POST("/user/find_username_by_email").WithHeader("Authorization", auth).WithJSON(&r.FindUsernameBySecureEmailBody{
		Email:            suite.Account.Email,
		SerialNumber:     resp.Data.SerialNumber,
		VerificationCode: evc,
	},
	).Expect().Body().Contains(suite.Account.Username)
}

// getLatestEVC 获取最新的email vc
func getLatestEVC(e *httpexpect.Expect) (string, error) {
	envFilepath := filepath.Join("testdata", "local.api-jinmuid.env")
	account := newTestingAccountFromEnvFile(envFilepath)
	auth, _ := getAuthorization(e)
	sendInformation := make([]r.GetLatestVerificationCode, 1)
	sendInformation[0] = r.GetLatestVerificationCode{
		SendVia: r.SendViaEmail,
		Email:   account.Email,
	}
	body := e.POST("/_debug/user/latest_verification_code").WithHeader("Authorization", auth).WithJSON(&r.LatestVerificationCodeBody{
		SendInformation: sendInformation,
	},
	).Expect().Body()
	var resp LatestVerificationCodes
	errUnmarshal := json.Unmarshal([]byte(body.Raw()), &resp)
	return resp.Data[0].VerificationCode, errUnmarshal
}

func TestFindUsernameByEmailTestSuite(t *testing.T) {
	suite.Run(t, new(FindUsernameByEmailTestSuite))
}
