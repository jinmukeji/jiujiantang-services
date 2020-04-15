package handler

import (
	"context"
	"path/filepath"
	"testing"

	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// UserSignOutTestSuite 是用户登录的单元测试的 Test Suite
type UserSignOutTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

// SetupSuite 设置测试环境
func (suite *UserSignOutTestSuite) SetupSuite() {
	envFilepath := filepath.Join("../mysqldb/testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestUserSignOut 测试用户登录
func (suite *UserSignOutTestSuite) TestUserSignOut() {
	t := suite.T()
	ctx := context.Background()
	resp := new(proto.UserSignInByPhonePasswordResponse)
	err := suite.JinmuIDService.UserSignInByPhonePassword(ctx, &proto.UserSignInByPhonePasswordRequest{
		Phone:          suite.Account.phone,
		HashedPassword: suite.Account.phonePassword,
		Seed:           suite.Account.seed,
		NationCode:     suite.Account.nationCode,
	}, resp)
	assert.NoError(t, err)
	assert.Equal(t, 105586, int(resp.UserId))
	// 退出登录操作
	signoutReq, signoutResp := new(proto.UserSignOutRequest), new(proto.UserSignOutResponse)
	ctx = AddContextToken(ctx, resp.AccessToken)
	assert.NoError(t, suite.JinmuIDService.UserSignOut(ctx, signoutReq, signoutResp))
}

func (suite *UserSignOutTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}

func TestUserSignOutTestSuite(t *testing.T) {
	suite.Run(t, new(UserSignOutTestSuite))
}
