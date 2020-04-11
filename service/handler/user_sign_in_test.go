package handler

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/jinmukeji/jiujiantang-services/service/auth"
	jinmuidpb "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	corepb "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// UserSigninTestSuite 是用户登录的单元测试的 Test Suite
type UserSigninTestSuite struct {
	suite.Suite
	jinmuHealth *JinmuHealth
	Account     *Account
}

// SetupSuite 设置测试环境
func (suite *UserSigninTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-biz-core.env")
	db, _ := newTestingDbClientFromEnvFile(envFilepath)
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
	suite.jinmuHealth = NewJinmuHealth(db, nil, nil, nil, nil, nil, nil, "")
}

// TestUserSiginin 测试用户登录
func (suite *UserSigninTestSuite) TestUserSignin() {
	t := suite.T()
	ctx := context.Background()
	ctx = mockAuth(ctx, suite.Account.clientID, suite.Account.name, suite.Account.zone)
	const registerType = "username"
	_, errMockSignin := mockSignin(ctx, suite.jinmuHealth, suite.Account.userName, suite.Account.passwordHash, registerType, corepb.SignInMethod_SIGN_IN_METHOD_GENERAL)
	assert.NoError(t, errMockSignin)
}

// TestUserSigninTestSuite 运行测试
func TestUserSigninTestSuite(t *testing.T) {
	suite.Run(t, new(UserSigninTestSuite))
}

// mockSignin 模拟登录
func mockSignin(ctx context.Context, j *JinmuHealth, username string, passwordHash string, registerType string, signInMethod corepb.SignInMethod) (context.Context, error) {
	resp, err := j.jinmuidSvc.UserSignInByUsernamePassword(ctx, &jinmuidpb.UserSignInByUsernamePasswordRequest{
		Username:       username,
		HashedPassword: passwordHash,
		Seed:           "",
		SignInMachine:  "",
	})
	if err != nil {
		return nil, err
	}
	return auth.AddContextToken(ctx, resp.AccessToken), nil
}
