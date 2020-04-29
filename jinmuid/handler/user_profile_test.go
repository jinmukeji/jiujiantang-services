package handler

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"

	jinmuidpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// UserProfileTestSuite 用户档案测试
type UserProfileTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

// SetupSuite 设置测试环境
func (suite *UserProfileTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestModifyUserProfile 测试修改用户档案
func (suite *UserProfileTestSuite) TestModifyUserProfile() {
	t := suite.T()
	// 登录
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	resp := new(jinmuidpb.ModifyUserProfileResponse)
	nickname := suite.Account.nickname
	birthday, _ := ptypes.TimestampProto(time.Now())
	protoGender := suite.Account.gender
	err = suite.JinmuIDService.ModifyUserProfile(ctx, &jinmuidpb.ModifyUserProfileRequest{
		UserId: userID,
		UserProfile: &jinmuidpb.UserProfile{
			Nickname:     nickname,
			Gender:       protoGender,
			BirthdayTime: birthday,
			Height:       suite.Account.height,
			Weight:       suite.Account.weight,
		},
	}, resp)
	assert.NoError(t, err)
	assert.Equal(t, nickname, resp.UserProfile.Nickname)
}

// TestGetUserProfile 测试得到用户档案
func (suite *UserProfileTestSuite) TestGetUserProfile() {
	t := suite.T()
	ctx := context.Background()
	// 登录
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	resp := new(jinmuidpb.GetUserProfileResponse)
	nickname := suite.Account.nickname
	err = suite.JinmuIDService.GetUserProfile(ctx, &jinmuidpb.GetUserProfileRequest{
		UserId: userID,
	}, resp)
	assert.NoError(t, err)
	assert.Equal(t, nickname, resp.Profile.Nickname)
}

// TestGetUserProfileByRecordID  通过记录id 获取用户信息
func (suite *UserProfileTestSuite) TestGetUserProfileByRecordID() {
	t := suite.T()
	ctx := context.Background()
	resp := new(jinmuidpb.GetUserProfileByRecordIDResponse)
	err := suite.JinmuIDService.GetUserProfileByRecordID(ctx, &jinmuidpb.GetUserProfileByRecordIDRequest{
		RecordId:          4,
		IsSkipVerifyToken: true,
	}, resp)
	assert.NoError(t, err)
}

func (suite *UserProfileTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}

func TestUserProfileTestSuite(t *testing.T) {
	suite.Run(t, new(UserProfileTestSuite))
}
