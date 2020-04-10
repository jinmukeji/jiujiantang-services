package handler

import (
	"context"
	"path/filepath"
	"testing"

	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ClientAuth 是客户端认证的单元测试的 Test Suite
type ClientAuthTestSuite struct {
	suite.Suite
	jinmuHealth *JinmuHealth
	Account     *Account
}

// SetupSuite 设置测试环境
func (suite *ClientAuthTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-biz-core.env")
	suite.jinmuHealth = newTestingJinmuHealthFromEnvFile(envFilepath)
	suite.jinmuHealth.datastore, _ = newTestingDbClientFromEnvFile(envFilepath)
	suite.jinmuHealth.mailClient, _ = newTestingMailClientFromEnvFile(envFilepath)
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestClientAuth 测试客户端授权
func (suite *ClientAuthTestSuite) TestClientAuth() {
	t := suite.T()
	ctx := context.Background()
	resp := new(proto.ClientAuthResponse)
	err := suite.jinmuHealth.ClientAuth(ctx, &proto.ClientAuthRequest{
		ClientId:      suite.Account.clientID,
		SecretKeyHash: suite.Account.secretKeyHash,
		Seed:          suite.Account.seed,
	}, resp)
	assert.NoError(t, err)
	assert.Equal(t, "jm-10005", resp.ClientId)
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientAuthTestSuite))
}

func mockAuth(ctx context.Context, clientID, name, zone string) context.Context {
	client := metaClient{
		ClientID:       clientID,
		Name:           name,
		Zone:           zone,
		CustomizedCode: "",
	}
	return addContextClient(ctx, client)
}

func (suite *ClientAuthTestSuite) TearDownSuite() {
	//To Do
}
