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

// GetLatestVerificationCodeSuite 是GetLatestVerificationCode的单元测试的 Test Suite
type GetLatestVerificationCodeSuite struct {
	suite.Suite
	Expect  *httpexpect.Expect
	Account *Account
}

// SetupSuite 设置测试环境
func (suite *GetLatestVerificationCodeSuite) SetupSuite() {
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

// LatestVerificationCodes 验证码的返回
type LatestVerificationCodes struct {
	Data []r.LatestVerificationCode `json:"data"`
}

// TestGetLatestVerificationCodes 获取最新的验证码
func (suite *GetLatestVerificationCodeSuite) TestGetLatestVerificationCodes() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	sendInformation := make([]r.GetLatestVerificationCode, 1)
	sendInformation[0] = r.GetLatestVerificationCode{
		SendVia:    r.SendViaPhone,
		Phone:      suite.Account.SignInPhone,
		NationCode: suite.Account.NationCode,
	}
	body := e.POST("/_debug/user/latest_verification_code").WithHeader("Authorization", auth).WithJSON(&r.LatestVerificationCodeBody{
		SendInformation: sendInformation,
	},
	).Expect().Body()
	var resp LatestVerificationCodes
	errUnmarshalSignIn := json.Unmarshal([]byte(body.Raw()), &resp)
	assert.NoError(t, errUnmarshalSignIn)
	assert.NotNil(t, resp.Data[0].VerificationCode)
}

func TestGetLatestVerificationCodeSuite(t *testing.T) {
	suite.Run(t, new(GetLatestVerificationCodeSuite))
}
