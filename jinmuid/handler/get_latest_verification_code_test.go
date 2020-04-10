package handler

import (
	"context"
	"path/filepath"
	"testing"

	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GetLatestVerificationCodesTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

// SetupSuite 设置测试环境
func (suite *GetLatestVerificationCodesTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestGetLatestVerificationCodes 测试获取最新验证码
func (suite *GetLatestVerificationCodesTestSuite) TestGetLatestVerificationCodes() {
	t := suite.T()
	ctx := context.Background()
	resp := new(proto.GetLatestVerificationCodesResponse)
	req := new(proto.GetLatestVerificationCodesRequest)
	latestVerificationCodes := make([]*proto.SingleGetLatestVerificationCode, 1)
	latestVerificationCodes[0] = &proto.SingleGetLatestVerificationCode{
		SendVia:    proto.SendVia_SEND_VIA_PHONE_SEND_VIA,
		Phone:      suite.Account.phone,
		NationCode: suite.Account.nationCode,
	}
	req.SendTo = latestVerificationCodes
	err := suite.JinmuIDService.GetLatestVerificationCodes(ctx, req, resp)
	assert.NoError(t, err)
}

func (suite *GetLatestVerificationCodesTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}
func TestGetLatestVeri1ficationCodesTestSuite(t *testing.T) {
	suite.Run(t, new(GetLatestVerificationCodesTestSuite))

}

// TestGetLatestHKVerificationCodes 获取最新的香港号码的验证码
func (suite *GetLatestVerificationCodesTestSuite) TestGetLatestHKVerificationCodes() {
	t := suite.T()
	ctx := context.Background()
	resp := new(proto.GetLatestVerificationCodesResponse)
	req := new(proto.GetLatestVerificationCodesRequest)
	latestVerificationCodes := make([]*proto.SingleGetLatestVerificationCode, 1)
	latestVerificationCodes[0] = &proto.SingleGetLatestVerificationCode{
		SendVia:    proto.SendVia_SEND_VIA_PHONE_SEND_VIA,
		Phone:      suite.Account.phoneHK,
		NationCode: suite.Account.nationCodeHK,
	}
	req.SendTo = latestVerificationCodes
	err := suite.JinmuIDService.GetLatestVerificationCodes(ctx, req, resp)
	assert.NoError(t, err)
}
