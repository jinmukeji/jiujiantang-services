package handler

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Organization 是组织单元测试的 Test Suite
type OrganizationTestSuite struct {
	suite.Suite
	jinmuHealth *JinmuHealth
}

// SetupSuite 设置测试环境
func (suite *OrganizationTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-biz-core.env")
	suite.jinmuHealth = newTestingJinmuHealthFromEnvFile(envFilepath)
	suite.jinmuHealth.datastore, _ = newTestingDbClientFromEnvFile(envFilepath)
	suite.jinmuHealth.mailClient, _ = newTestingMailClientFromEnvFile(envFilepath)
}

// TestCreateOrganization 测试创建组织
func (suite *OrganizationTestSuite) TestCreateOrganization() {
	const username = "4"
	const passwordHash = "97951af80347d78d63bf3a7b7962fb42dd21ef56d4e590f2be2b7954475f4089"

	t := suite.T()
	const registerType = "username"
	ctx, err := mockSignin(context.Background(), suite.jinmuHealth, username, passwordHash, registerType, proto.SignInMethod_SIGN_IN_METHOD_GENERAL)
	assert.NoError(t, err)
	req, resp := new(proto.OwnerCreateOrganizationRequest), new(proto.OwnerCreateOrganizationResponse)
	req.Profile = &proto.OrganizationProfile{
		Name:    uuid.New().String(),
		Address: &proto.Address{},
	}
	assert.Error(t, NewError(ErrOrganizationCountExceedsMaxLimits, errors.New("Organization count exceeds the max limits")), suite.jinmuHealth.OwnerCreateOrganization(ctx, req, resp))
}

// TestOwnerGetOrganizations 测试拥有者获取组织
func (suite *OrganizationTestSuite) TestOwnerGetOrganizations() {
	const username = "4"
	const passwordHash = "97951af80347d78d63bf3a7b7962fb42dd21ef56d4e590f2be2b7954475f4089"
	t := suite.T()
	const registerType = "username"
	ctx, err := mockSignin(context.Background(), suite.jinmuHealth, username, passwordHash, registerType, proto.SignInMethod_SIGN_IN_METHOD_GENERAL)
	assert.NoError(t, err)
	req, resp := new(proto.OwnerGetOrganizationsRequest), new(proto.OwnerGetOrganizationsResponse)
	assert.NoError(t, suite.jinmuHealth.OwnerGetOrganizations(ctx, req, resp))
	assert.Equal(t, "", resp.Organizations[0].Profile.Name)
	assert.Equal(t, "", resp.Organizations[0].Profile.Contact)
}

// TestOwnerAddOrganizationUsers 测试组织下新增用户
func (suite *OrganizationTestSuite) TestOwnerAddOrganizationUsers() {
	t := suite.T()
	const organizationID = 96
	const username = "4"
	const passwordHash = "97951af80347d78d63bf3a7b7962fb42dd21ef56d4e590f2be2b7954475f4089"
	const clientID = "jm-10005"
	const name = "JinmuHealth-Android-app"
	const zone = "CN"
	ctx := mockAuth(context.Background(), clientID, name, zone)
	const registerType = "username"
	ctx, err := mockSignin(ctx, suite.jinmuHealth, username, passwordHash, registerType, proto.SignInMethod_SIGN_IN_METHOD_GENERAL)
	assert.NoError(t, err)
	req, resp := new(proto.UserSignUpRequest), new(proto.UserSignUpResponse)
	req.Password = "release4"
	req.ClientId = clientID
	req.UserProfile = &proto.UserProfile{
		Nickname: "liu",
	}
	assert.NoError(t, suite.jinmuHealth.UserSignUp(ctx, req, resp))
	reqOwnerAddOrganizationUsers, respOwnerAddOrganizationUsers := new(proto.OwnerAddOrganizationUsersRequest), new(proto.OwnerAddOrganizationUsersResponse)
	reqOwnerAddOrganizationUsers.OrganizationId = organizationID
	reqOwnerAddOrganizationUsers.UserIdList = []int32{resp.User.UserId}
	assert.NoError(t, suite.jinmuHealth.OwnerAddOrganizationUsers(ctx, reqOwnerAddOrganizationUsers, respOwnerAddOrganizationUsers))
}

// TestOwnerDeleteOrganizationUsers 测试组织下删除用户
func (suite *OrganizationTestSuite) TestOwnerDeleteOrganizationUsers() {
	t := suite.T()
	const organizationID = 96
	const username = "4"
	const passwordHash = "97951af80347d78d63bf3a7b7962fb42dd21ef56d4e590f2be2b7954475f4089"
	const clientID = "jm-10005"
	const name = "JinmuHealth-Android-app"
	const zone = "CN"
	ctx := mockAuth(context.Background(), clientID, name, zone)
	const registerType = "username"
	ctx, err := mockSignin(ctx, suite.jinmuHealth, username, passwordHash, registerType, proto.SignInMethod_SIGN_IN_METHOD_GENERAL)
	assert.NoError(t, err)
	req, resp := new(proto.UserSignUpRequest), new(proto.UserSignUpResponse)
	req.Password = "release4"
	req.ClientId = clientID
	req.UserProfile = &proto.UserProfile{
		Nickname: "liu",
	}
	assert.NoError(t, suite.jinmuHealth.UserSignUp(ctx, req, resp))
	reqOwnerAddOrganizationUsers, respOwnerAddOrganizationUsers := new(proto.OwnerAddOrganizationUsersRequest), new(proto.OwnerAddOrganizationUsersResponse)
	reqOwnerAddOrganizationUsers.OrganizationId = organizationID
	reqOwnerAddOrganizationUsers.UserIdList = []int32{resp.User.UserId}
	assert.NoError(t, suite.jinmuHealth.OwnerAddOrganizationUsers(ctx, reqOwnerAddOrganizationUsers, respOwnerAddOrganizationUsers))

	reqOwnerDeleteOrganizationUser, respOwnerDeleteOrganizationUser := new(proto.OwnerDeleteOrganizationUsersRequest), new(proto.OwnerDeleteOrganizationUsersResponse)
	reqOwnerDeleteOrganizationUser.OrganizationId = organizationID
	reqOwnerDeleteOrganizationUser.UserIdList = []int32{resp.User.UserId}
	assert.NoError(t, suite.jinmuHealth.OwnerDeleteOrganizationUsers(ctx, reqOwnerDeleteOrganizationUser, respOwnerDeleteOrganizationUser))
}

// TestOwnerGetOrganizationUsers测试组织下查看用户
func (suite *OrganizationTestSuite) TestOwnerGetOrganizationUsers() {
	const organizationID = 96
	const username = "4"
	const passwordHash = "97951af80347d78d63bf3a7b7962fb42dd21ef56d4e590f2be2b7954475f4089"
	t := suite.T()
	const registerType = "username"
	ctx, err := mockSignin(context.Background(), suite.jinmuHealth, username, passwordHash, registerType, proto.SignInMethod_SIGN_IN_METHOD_GENERAL)
	assert.NoError(t, err)
	// 从组织查看这些用户
	getUserReq, getUserRepl := new(proto.OwnerGetOrganizationUsersRequest), new(proto.OwnerGetOrganizationUsersResponse)
	getUserReq.OrganizationId = organizationID
	getUserReq.Size = 1
	assert.NoError(t, suite.jinmuHealth.OwnerGetOrganizationUsers(ctx, getUserReq, getUserRepl))
}

// TestOrganizationTestSuite 启动测试
func TestOrganizationtTestSuite(t *testing.T) {
	suite.Run(t, new(OrganizationTestSuite))
}
