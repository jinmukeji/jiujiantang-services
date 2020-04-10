package handler

import (
	"context"
	"path/filepath"
	"testing"

	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// UserGetInformation 获取用户信息
type ModifyUserInformationTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

// ModifyUserInformationTestSuite 设置测试环境
func (suite *ModifyUserInformationTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestModifyUserInformation 修改用户信息
func (suite *ModifyUserInformationTestSuite) TestModifyUserInformation() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	req := new(proto.ModifyUserInformationRequest)
	req.UserId = userID
	req.SigninUsername = suite.Account.username
	req.SigninPhone = suite.Account.phone
	req.SecureEmail = suite.Account.email
	req.Remark = suite.Account.remark
	req.CustomizedCode = suite.Account.customizedCode
	req.HasSetUserProfile = false

	resp := new(proto.ModifyUserInformationResponse)
	err = suite.JinmuIDService.ModifyUserInformation(ctx, req, resp)
	assert.NoError(t, err)
}

func (suite *ModifyUserInformationTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}

func TestModifyUserInformationTestSuite(t *testing.T) {
	suite.Run(t, new(ModifyUserInformationTestSuite))
}
