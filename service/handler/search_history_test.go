package handler

import (
	"context"
	"path/filepath"
	"testing"

	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// SearchHistory是 SearchHistory rpc 的单元测试的 Test Suite
type SearchHistoryTestSuite struct {
	suite.Suite
	jinmuHealth *JinmuHealth
}

// SetupSuite 设置测试环境
func (suite *SearchHistoryTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-biz-core.env")
	db, _ := newTestingDbClientFromEnvFile(envFilepath)
	suite.jinmuHealth = NewJinmuHealth(db, nil, nil, nil, nil, nil, nil, "")
}

// TestDeleteRecord 测试删除记录
func (suite *SearchHistoryTestSuite) TestDeleteRecord() {
	t := suite.T()
	ctx := context.Background()
	const username = "4"
	const passwordHash = "97951af80347d78d63bf3a7b7962fb42dd21ef56d4e590f2be2b7954475f4089"
	const clientID = "jm-10005"
	const name = "JinmuHealth-Android-app"
	const zone = "CN"
	ctx = mockAuth(ctx, clientID, name, zone)
	const registerType = "username"
	ctx, err := mockSignin(ctx, suite.jinmuHealth, username, passwordHash, registerType, proto.SignInMethod_SIGN_IN_METHOD_GENERAL)
	assert.NoError(t, err)
	req, resp := new(proto.DeleteRecordRequest), new(proto.DeleteRecordResponse)
	req.UserId = int32(96)
	req.RecordId = int32(204452)
	assert.NoError(t, suite.jinmuHealth.DeleteRecord(ctx, req, resp))
}

func TestSearchHistoryTestSuite(t *testing.T) {
	suite.Run(t, new(SearchHistoryTestSuite))
}
