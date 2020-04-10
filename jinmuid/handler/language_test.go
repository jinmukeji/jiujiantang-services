package handler

import (
	"context"
	"path/filepath"
	"testing"

	generalpb "github.com/jinmukeji/proto/gen/micro/idl/ptypes/v2"
	jinmuidpb "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// LanguageTestSuite 语言测试
type LanguageTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

// LanguageTestSuite 设置测试环境
func (suite *LanguageTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestSetJinmuIDWebLanguage 测试设置金姆ID语言
func (suite *UserTestSuite) TestSetJinmuIDWebLanguage() {
	t := suite.T()
	ctx := context.Background()
	resp := new(jinmuidpb.SetJinmuIDWebLanguageResponse)
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	err = suite.JinmuIDService.SetJinmuIDWebLanguage(ctx, &jinmuidpb.SetJinmuIDWebLanguageRequest{
		UserId:   userID,
		Language: generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE,
	}, resp)
	assert.NoError(t, err)
}

// TestSetJinmuIDWebLanguageUserIdIsNull 测试设置金姆ID语言
func (suite *UserTestSuite) TestSetJinmuIDWebLanguageUserIdIsNull() {
	t := suite.T()
	ctx := context.Background()
	resp := new(jinmuidpb.SetJinmuIDWebLanguageResponse)
	err := suite.JinmuIDService.SetJinmuIDWebLanguage(ctx, &jinmuidpb.SetJinmuIDWebLanguageRequest{
		UserId:   suite.Account.userID,
		Language: generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE,
	}, resp)
	assert.NoError(t, err)
}

// TestGetJinmuIDWebLanguage 测试得到金姆ID语言
func (suite *UserTestSuite) TestGetJinmuIDWebLanguage() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	resp := new(jinmuidpb.GetJinmuIDWebLanguageResponse)
	err = suite.JinmuIDService.GetJinmuIDWebLanguage(ctx, &jinmuidpb.GetJinmuIDWebLanguageRequest{
		UserId: userID,
	}, resp)
	assert.NoError(t, err)
	assert.Equal(t, generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE, resp.Language)
}

// TestGetJinmuIDWebLanguageUserIdIsError  测试得到金姆ID语言
func (suite *UserTestSuite) TestGetJinmuIDWebLanguageUserIdIsError() {
	t := suite.T()
	ctx := context.Background()
	resp := new(jinmuidpb.GetJinmuIDWebLanguageResponse)
	err := suite.JinmuIDService.GetJinmuIDWebLanguage(ctx, &jinmuidpb.GetJinmuIDWebLanguageRequest{
		UserId: suite.Account.userID,
	}, resp)
	assert.NoError(t, err)
	assert.Equal(t, generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE, resp.Language)
}

func (suite *LanguageTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}

func TestLanguageTestSuite(t *testing.T) {
	suite.Run(t, new(LanguageTestSuite))
}
