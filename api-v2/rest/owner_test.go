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

// OwnerTestSuite 是Owner的单元测试的 Test Suite
type OwnerTestSuite struct {
	suite.Suite
	Account *Account
	Expect  *httpexpect.Expect
}

// SetupSuite 设置测试环境
func (suite *OwnerTestSuite) SetupSuite() {
	t := suite.T()
	envFilepath := filepath.Join("testdata", "local.api-v2.env")
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
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

// OwnerUserSignUpReply 注册的返回
type OwnerUserSignUpReply struct {
	Data r.OwnerUserSignUpReply `json:"data"`
}

// TestOwnerDeleteUsers 测试删除用户
func (suite *OwnerTestSuite) TestOwnerDeleteUsers() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, errGetAccessToken := getAccessToken(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	var nickname = "火娃"
	year1993 := time.Date(1993, 1, 1, 0, 0, 0, 0, time.UTC)
	var gender = int32(r.GenderMale)
	var height = int32(170)
	var weight = int32(70)
	envFilepath := filepath.Join("testdata", "local.api-v2.env")
	client := newTestingClientFromEnvFile(envFilepath)
	// 注册一个用户
	body := e.POST("/v2-api/owner/{owner_id}/users/sign_up").WithHeaders(headers).WithJSON(&r.OwnerUserSignUpBody{
		Nickname:     nickname,
		Birthday:     year1993,
		Gender:       &gender,
		Height:       height,
		Weight:       weight,
		RegisterType: suite.Account.RegisterType,
		ClientID:     client.ClientID,
	}).
		WithPath("owner_id", suite.Account.UserID).
		Expect().Body()
	var resp OwnerUserSignUpReply
	errUnmarshalSignIn := json.Unmarshal([]byte(body.Raw()), &resp)
	assert.NoError(t, errUnmarshalSignIn)
	// 删除用户 UserIDList
	e.POST("/v2-api/owner/{owner_id}/users/delete").WithHeaders(headers).
		WithPath("owner_id", suite.Account.UserID).
		WithJSON(&r.UserIDList{
			UserIDList: []int32{int32(resp.Data.UserID)},
		}).
		Expect().Body().Contains("ok").Contains("true")
}

func TestOwnerTestSuite(t *testing.T) {
	suite.Run(t, new(OwnerTestSuite))
}
