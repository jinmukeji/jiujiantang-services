package handler

import (
	"context"
	"path/filepath"
	"testing"

	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SetNotificationPreferencesTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

// SetupSuite 设置测试环境
func (suite *SetNotificationPreferencesTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestSetNotificationPrefences  设置个人资料
func (suite *SetNotificationPreferencesTestSuite) TestSetNotificationPrefences() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)

	req := new(proto.ModifyNotificationPreferencesRequest)
	req.UserId = userID
	req.PhoneEnabled = true
	req.WechatEnabled = true
	req.WeiboEnabled = true
	resp := new(proto.ModifyNotificationPreferencesResponse)
	err = suite.JinmuIDService.ModifyNotificationPreferences(ctx, req, resp)

	assert.NoError(t, err)
}

func (suite *SetNotificationPreferencesTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}

func TestSetNotificationPreferencesTestSuite(t *testing.T) {
	suite.Run(t, new(SetNotificationPreferencesTestSuite))
}
