package rest_test

import (
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/iris-contrib/httpexpect"
	r "github.com/jinmukeji/gf-api2/api-jinmuid/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// UserProfileSuite 是UserProfile的单元测试的 Test Suite
type UserProfileSuite struct {
	suite.Suite
	Expect  *httpexpect.Expect
	Account *Account
}

// SetupSuite 设置测试环境
func (suite *UserProfileSuite) SetupSuite() {
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

// TestGetUserProfile 测试GetUserProfile
func (suite *UserProfileSuite) TestGetUserProfile() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, userID, errGetAccessToken := getAccessTokenAndUserID(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	e.GET("/user/{user_id}/profile").
		WithHeaders(headers).
		WithPath("user_id", userID).Expect().Body().Contains(suite.Account.Nickname)
}

// TestModifyUserProfile 测试修改UserProfile
func (suite *UserProfileSuite) TestModifyUserProfile() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, userID, errGetAccessToken := getAccessTokenAndUserID(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	year1993 := time.Date(1993, 1, 1, 0, 0, 0, 0, time.UTC)
	var gender = int32(r.GenderMale)
	var height = int32(170)
	var weight = int32(70)
	e.PUT("/user/{user_id}/profile").
		WithHeaders(headers).
		WithPath("user_id", userID).WithJSON(
		&r.UserProfile{
			Nickname: suite.Account.Nickname,
			Gender:   gender,
			Birthday: year1993,
			Height:   height,
			Weight:   weight,
		},
	).Expect().Body().Contains(suite.Account.Nickname)
}

func TestUserProfileSuite(t *testing.T) {
	suite.Run(t, new(UserProfileSuite))
}
