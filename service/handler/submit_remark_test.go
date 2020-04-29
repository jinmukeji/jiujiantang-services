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

// SubmitRemarkTestSuite 是 SubmitRemark rpc 的单元测试的 Test Suite
type SubmitRemarkTestSuite struct {
	suite.Suite
	jinmuHealth *JinmuHealth
}

// SetupSuite 设置测试环境
func (suite *SubmitRemarkTestSuite) SetupSuite() {
	suite.jinmuHealth = new(JinmuHealth)
	envFilepath := filepath.Join("testdata", "local.svc-biz-core.env")
	suite.jinmuHealth.datastore, _ = newTestingDbClientFromEnvFile(envFilepath)
	suite.jinmuHealth.mailClient, _ = newTestingMailClientFromEnvFile(envFilepath)
}

// TestSubmitRemarkSuccess 测试提交数据成功
func (suite *SubmitRemarkTestSuite) TestSubmitRemarkSuccess() {
	const username = "4"
	const passwordHash = "97951af80347d78d63bf3a7b7962fb42dd21ef56d4e590f2be2b7954475f4089"
	const clientID = "jm-10005"
	const name = "JinmuHealth-Android-app"
	const zone = "CN"
	const testRecordID = 183289
	const userID = 96
	t := suite.T()
	ctx := context.Background()
	ctx = mockAuth(ctx, clientID, name, zone)
	const registerType = "username"
	ctx, err := mockSignin(ctx, suite.jinmuHealth, username, passwordHash, registerType, proto.SignInMethod_SIGN_IN_METHOD_GENERAL)
	assert.NoError(t, err)
	req, resp := new(proto.SubmitRemarkRequest), new(proto.SubmitRemarkResponse)
	req.Remark, req.RecordId, req.UserId = uuid.New().String(), int32(testRecordID), userID
	assert.NoError(t, suite.jinmuHealth.SubmitRemark(ctx, req, resp))
	// 查找这个 record
}

func TestSubmitRemarkTestSuite(t *testing.T) {
	suite.Run(t, new(SubmitRemarkTestSuite))
}
