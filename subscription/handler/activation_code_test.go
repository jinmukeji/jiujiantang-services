package handler

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	ac "github.com/jinmukeji/jiujiantang-services/subscription/activation-code"
	subscriptionpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/subscription/v1"
	jinmuidpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	"github.com/micro/go-micro/v2/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// SubscriptionTestSuite 测试激活码
type SubscriptionActivationCodeTestSuite struct {
	suite.Suite
	subscriptionService *SubscriptionService
	account             *Account
	jinmuidSrv          jinmuidpb.UserManagerAPIService
}

const rpcJinmuidServiceName = "com.himalife.srv.svc-jinmuid"

// SetupSuite 初始化测试
func (suite *SubscriptionActivationCodeTestSuite) SetupSuite() {
	suite.subscriptionService = new(SubscriptionService)
	envFilepath := filepath.Join("testdata", "local.svc-subscription.env")
	db, _ := newTestingDbClientFromEnvFile(envFilepath)
	account := newTestingAccountFromEnvFile(envFilepath)
	suite.subscriptionService = NewSubscriptionService(db, account.activationCodeEncryptKey)
	suite.account = account
	suite.jinmuidSrv = jinmuidpb.NewUserManagerAPIService(rpcJinmuidServiceName, client.DefaultClient)
}

// TestGetSubscriptionActivationCodeInfo 测试获取激活码的内容
func (suite *SubscriptionActivationCodeTestSuite) TestGetSubscriptionActivationCodeInfo() {
	t := suite.T()
	ctx := context.Background()
	resp := new(subscriptionpb.GetSubscriptionActivationCodeInfoResponse)
	err := suite.subscriptionService.GetSubscriptionActivationCodeInfo(ctx, &subscriptionpb.GetSubscriptionActivationCodeInfoRequest{
		Code: suite.account.code,
	}, resp)
	assert.Error(t, errors.New("[errcode:1700] activation code is activated"), err)
}

// TestUseSubscriptionActivationCode 测试使用激活码
func (suite *SubscriptionActivationCodeTestSuite) TestUseSubscriptionActivationCode() {
	t := suite.T()
	ctx := context.Background()
	ctx, err := mockSignin(ctx, suite.jinmuidSrv, suite.account.account, suite.account.hashPassword, suite.account.seed)
	assert.NoError(t, err)
	resp := new(subscriptionpb.UseSubscriptionActivationCodeResponse)
	err = suite.subscriptionService.UseSubscriptionActivationCode(ctx, &subscriptionpb.UseSubscriptionActivationCodeRequest{
		Code:   suite.account.code,
		UserId: suite.account.userID,
	}, resp)
	assert.Error(t, errors.New("[errcode:1700] activation code is activated"), err)
}

// TestCheckActivationCodeAlgorithm 测试激活码算法
func (suite *SubscriptionActivationCodeTestSuite) TestCheckActivationCodeAlgorithm() {
	t := suite.T()
	helper := ac.NewActivationCodeCipherHelper()
	encryptCode := helper.Encrypt(suite.account.activationCode, suite.account.activationCodeEncryptKey, suite.account.contractYear, suite.account.maxUserLimits)
	code := helper.Decrypt(encryptCode, suite.account.activationCodeEncryptKey, suite.account.contractYear, suite.account.maxUserLimits)
	assert.Equal(t, code, suite.account.activationCode)
}

func TestSubscriptionActivationCodeTestSuite(t *testing.T) {
	suite.Run(t, new(SubscriptionActivationCodeTestSuite))
}
