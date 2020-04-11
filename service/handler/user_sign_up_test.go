package handler

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/golang/protobuf/ptypes"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	"github.com/stretchr/testify/suite"
)

// UserSignUpTestSuite 是用户登录的单元测试的 Test Suite
type UserSignUpTestSuite struct {
	suite.Suite
	jinmuHealth *JinmuHealth
}

// SetupSuite 设置测试环境
func (suite *UserSignUpTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-biz-core.env")
	suite.jinmuHealth = newTestingJinmuHealthFromEnvFile(envFilepath)
	suite.jinmuHealth.datastore, _ = newTestingDbClientFromEnvFile(envFilepath)
	suite.jinmuHealth.mailClient, _ = newTestingMailClientFromEnvFile(envFilepath)
}

// TestUserSiginup 测试注册用户
func (suite *UserSignUpTestSuite) TestUserSignUp() {
	t := suite.T()
	ctx := context.Background()
	now := ptypes.TimestampNow()
	randUsername := uuid.New().String()
	profile := &proto.UserProfile{
		Nickname:     "abcabc",
		Gender:       0,
		BirthdayTime: now,
	}
	req, resp := new(proto.UserSignUpRequest), new(proto.UserSignUpResponse)
	req.Username = randUsername
	req.Password = "password"
	req.UserProfile = profile
	assert.NoError(t, suite.jinmuHealth.UserSignUp(ctx, req, resp))
	u, err := suite.jinmuHealth.datastore.FindUserByUserID(ctx, int(resp.User.UserId))
	assert.NoError(t, err)
	assert.Equal(t, proto.SignInMethod_name[int32(proto.SignInMethod_SIGN_IN_METHOD_LEGACY)], u.RegisterType)
}

func TestUserSignUpTestSuite(t *testing.T) {
	suite.Run(t, new(UserSignUpTestSuite))
}
