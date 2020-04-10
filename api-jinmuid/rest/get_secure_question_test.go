package rest_test

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/iris-contrib/httpexpect"
	r "github.com/jinmukeji/jiujiantang-services/api-jinmuid/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// GetSecureQuestionSuite 是GetSecureQuestion的单元测试的 Test Suite
type GetSecureQuestionSuite struct {
	suite.Suite
	Expect  *httpexpect.Expect
	Account *Account
}

// SetupSuite 设置测试环境
func (suite *GetSecureQuestionSuite) SetupSuite() {
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

// TestGetSecureQuestionList 测试GetSecureQuestionList
func (suite *GetSecureQuestionSuite) TestGetSecureQuestionList() {
	t := suite.T()
	e := suite.Expect
	auth, errGetAuthorization := getAuthorization(e)
	assert.NoError(t, errGetAuthorization)
	token, userID, errGetAccessToken := getAccessTokenAndUserID(e)
	assert.NoError(t, errGetAccessToken)
	headers := make(map[string]string)
	headers["Authorization"] = auth
	headers["X-Access-Token"] = token
	e.GET("/user/{user_id}/secure_question_list").
		WithHeaders(headers).
		WithPath("user_id", userID).
		Expect().Body().
		Contains("你少年时代最好的朋友叫什么名字？").
		Contains("你学会做的第一道菜是什么？").
		Contains("你第一次去电影院看的是哪一部电影？").
		Contains("你第一次坐飞机是去哪里？").
		Contains("你上小学时最喜欢的老师姓什么？").
		Contains("你的父母是在哪里认识的？").
		Contains("你的第一个上司叫什么名字？").
		Contains("你从小长大的那条街叫什么？").
		Contains("你去过的第一个游乐场是哪一个？").
		Contains("你购买的第一张专辑是什么？").
		Contains("你最喜欢哪个球队？").
		Contains("你的理想工作是什么？").
		Contains("你小时候最喜欢哪一本书？").
		Contains("你童年时的绰号是什么？").
		Contains("你拥有的第一辆车是什么型号？").
		Contains("你在学生时代最喜欢的电影明星是谁？").
		Contains("你最喜欢哪个乐队或歌手？")
}

func TestGetSecureQuestionSuite(t *testing.T) {
	suite.Run(t, new(GetSecureQuestionSuite))
}
