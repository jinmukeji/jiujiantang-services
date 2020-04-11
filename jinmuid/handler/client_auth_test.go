package handler

import (
	"context"
	"errors"
	"testing"

	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ClientAuthTestSuite 授权测试
type ClientAuthTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
}

// SetupSuite 设置测试环境
func (suite *ClientAuthTestSuite) SetupSuite() {
	suite.JinmuIDService = newJinmuIDServiceForTest()
}

// TestClientAuth 测试客户端授权
func (suite *ClientAuthTestSuite) TestClientAuth() {
	const testSecretKeyHash = "7915550835f93aab534a245f3f498b65397fff77bdda7dab23f1beb1be56bac7"
	t := suite.T()
	ctx := context.Background()
	resp := new(proto.ClientAuthResponse)
	err := suite.JinmuIDService.ClientAuth(ctx, &proto.ClientAuthRequest{
		ClientId:      "jm-10006",
		SecretKeyHash: testSecretKeyHash,
		Seed:          "",
	}, resp)
	assert.NoError(t, err)
	assert.Equal(t, "jm-10006", resp.ClientId)
}

// TestClientAuthIsNull  clientId为空
func (suite *ClientAuthTestSuite) TestClientAuthIsNull() {
	const testSecretKeyHash = "7915550835f93aab534a245f3f498b65397fff77bdda7dab23f1beb1be56bac7"
	t := suite.T()
	ctx := context.Background()
	resp := new(proto.ClientAuthResponse)
	err := suite.JinmuIDService.ClientAuth(ctx, &proto.ClientAuthRequest{
		ClientId:      "",
		SecretKeyHash: testSecretKeyHash,
		Seed:          "",
	}, resp)
	assert.Error(t, errors.New("[errcode:10001] database error"), err)
}

// TestClientAuthNotExist   clientId不存在
func (suite *ClientAuthTestSuite) TestClientAuthNotExist() {
	const testSecretKeyHash = "7915550835f93aab534a245f3f498b65397fff77bdda7dab23f1beb1be56bac7"
	t := suite.T()
	ctx := context.Background()
	resp := new(proto.ClientAuthResponse)
	err := suite.JinmuIDService.ClientAuth(ctx, &proto.ClientAuthRequest{
		ClientId:      "jm-10000000000006",
		SecretKeyHash: testSecretKeyHash,
		Seed:          "",
	}, resp)
	assert.Error(t, errors.New("[errcode:10001] database error"), err)
}

func (suite *ClientAuthTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientAuthTestSuite))
}
