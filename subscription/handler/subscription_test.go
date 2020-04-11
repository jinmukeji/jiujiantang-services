package handler

import (
	"context"
	"path/filepath"
	"testing"

	jinmuidpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/subscription/v1"
	"github.com/micro/go-micro/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// SubscriptionTestSuite 测试帐号和 mac 的关联
type SubscriptionTestSuite struct {
	suite.Suite
	subscriptionService *SubscriptionService
	account             *Account
	jinmuidSrv          jinmuidpb.UserManagerAPIService
}

type Account struct {
	account                  string
	password                 string
	userID                   int32
	seed                     string
	hashPassword             string
	code                     string
	activationCodeEncryptKey string
	activationCode           string
	contractYear             int32
	maxUserLimits            int32
}

// SetupSuite 初始化测试
func (suite *SubscriptionTestSuite) SetupSuite() {
	suite.subscriptionService = new(SubscriptionService)
	envFilepath := filepath.Join("testdata", "local.svc-subscription.env")
	suite.subscriptionService.datastore, _ = newTestingDbClientFromEnvFile(envFilepath)
	suite.account = newTestingAccountFromEnvFile(envFilepath)
	suite.jinmuidSrv = jinmuidpb.NewUserManagerAPIService(rpcJinmuidServiceName, client.DefaultClient)
}

// TestGetUserSubscriptions 测试得到使用中的订阅
func (suite *SubscriptionTestSuite) TestGetUserSubscriptions() {
	t := suite.T()
	ctx := context.Background()
	ctx, err := mockSignin(ctx, suite.jinmuidSrv, suite.account.account, suite.account.hashPassword, suite.account.seed)
	assert.NoError(t, err)
	resp := new(proto.GetUserSubscriptionsResponse)
	err = suite.subscriptionService.GetUserSubscriptions(ctx, &proto.GetUserSubscriptionsRequest{
		UserId: suite.account.userID,
	}, resp)
	assert.NoError(t, err)
	assert.Equal(t, int32(300), resp.Subscriptions[0].MaxUserLimits)
}

func TestSubscriptionTestSuite(t *testing.T) {
	suite.Run(t, new(SubscriptionTestSuite))
}

// mockSignin 模拟登录
func mockSignin(ctx context.Context, rpcUserManagerSrv jinmuidpb.UserManagerAPIService, username string, passwordHash, seed string) (context.Context, error) {
	resp, err := rpcUserManagerSrv.UserSignInByUsernamePassword(ctx, &jinmuidpb.UserSignInByUsernamePasswordRequest{
		Username:       username,
		HashedPassword: passwordHash,
		Seed:           seed,
	})
	if err != nil {
		return nil, err
	}
	return AddContextToken(ctx, resp.AccessToken), nil
}
