package handler

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// SubmitCommentTestSuite 是 SubmitComment rpc 的单元测试的 Test Suite
type SubmitCommentTestSuite struct {
	suite.Suite
	jinmuHealth *JinmuHealth
}

// SetupSuite 设置测试环境
func (suite *SubmitCommentTestSuite) SetupSuite() {
	suite.jinmuHealth = new(JinmuHealth)
	envFilepath := filepath.Join("testdata", "local.svc-biz-core.env")
	suite.jinmuHealth.datastore, _ = newTestingDbClientFromEnvFile(envFilepath)
	suite.jinmuHealth.mailClient, _ = newTestingMailClientFromEnvFile(envFilepath)
}

// TestSubmitFeedbackSuccess 测试提交数据成功
func (suite *SubmitCommentTestSuite) TestSubmitFeedbackSuccess() {
	const username = "4"
	const passwordHash = "97951af80347d78d63bf3a7b7962fb42dd21ef56d4e590f2be2b7954475f4089"
	const clientID = "jm-10005"
	const name = "JinmuHealth-Android-app"
	const zone = "CN"
	t := suite.T()
	ctx := context.Background()
	ctx = mockAuth(ctx, clientID, name, zone)
	const registerType = "username"
	ctx, err := mockSignin(ctx, suite.jinmuHealth, username, passwordHash, registerType, proto.SignInMethod_SIGN_IN_METHOD_GENERAL)
	assert.NoError(t, err)
	req, resp := new(proto.SubmitFeedbackRequest), new(proto.SubmitFeedbackResponse)
	req.Content = uuid.New().String()
	req.ContactWay = uuid.New().String()
	assert.NoError(t, suite.jinmuHealth.SubmitFeedback(ctx, req, resp))
}

func TestSubmitCommentTestSuite(t *testing.T) {
	suite.Run(t, new(SubmitCommentTestSuite))
}
