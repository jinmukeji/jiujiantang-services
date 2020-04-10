package handler

import (
	"context"
	"math/rand"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/protobuf/ptypes"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	"github.com/stretchr/testify/suite"
)

// ModifyUserProfileTestSuite 是用户登录的单元测试的 Test Suite
type ModifyUserProfileTestSuite struct {
	suite.Suite
	jinmuHealth *JinmuHealth
}

// SetupSuite 设置测试环境
func (suite *ModifyUserProfileTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-biz-core.env")
	suite.jinmuHealth = newTestingJinmuHealthFromEnvFile(envFilepath)
	suite.jinmuHealth.datastore, _ = newTestingDbClientFromEnvFile(envFilepath)
	suite.jinmuHealth.mailClient, _ = newTestingMailClientFromEnvFile(envFilepath)
}

// TestUserSiginin 测试用户登录
func (suite *ModifyUserProfileTestSuite) TestModifyUserProfile() {
	t := suite.T()
	const testUserID = 1
	randName := strconv.Itoa(rand.Int())
	ptypesNow := ptypes.TimestampNow()
	p := &proto.UserProfile{
		BirthdayTime: ptypesNow,
		Nickname:     randName,
	}
	ctx := context.Background()
	req, resp := new(proto.ModifyUserProfileRequest), new(proto.ModifyUserProfileResponse)
	req.UserId = testUserID
	req.UserProfile = p
	assert.NoError(t, suite.jinmuHealth.ModifyUserProfile(ctx, req, resp))
	u, _ := suite.jinmuHealth.datastore.FindUserByUserID(ctx, testUserID)
	assert.Equal(t, randName, u.Nickname)
}

func TestModifyUserProfileTestSuite(t *testing.T) {
	suite.Run(t, new(ModifyUserProfileTestSuite))
}
