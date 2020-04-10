package handler

import (
	"context"
	"path/filepath"

	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// RecordTokenSuite 是 RecordToken rpc 的单元测试的 Test Suite
type RecordTokenSuite struct {
	suite.Suite
	jinmuHealth *JinmuHealth
}

// SetupSuite 设置测试环境
func (suite *RecordTokenSuite) SetupSuite() {
	suite.jinmuHealth = new(JinmuHealth)
	envFilepath := filepath.Join("testdata", "local.svc-biz-core.env")
	suite.jinmuHealth.datastore, _ = newTestingDbClientFromEnvFile(envFilepath)
	suite.jinmuHealth.mailClient, _ = newTestingMailClientFromEnvFile(envFilepath)
}

// TestRecordToken 测试生产或者查找的recordToken
func (suite *RecordTokenSuite) TestRecordToken() {
	t := suite.T()
	ctx := context.Background()
	req := new(proto.CreateReportShareTokenRequest)
	req.RecordId = 1
	resp := new(proto.CreateReportShareTokenResponse)
	err := suite.jinmuHealth.CreateReportShareToken(ctx, req, resp)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Token)
}

// TestGetAnalyzeResultByToken 测试token获取分析报告
func (suite *RecordTokenSuite) TestGetAnalyzeResultByToken() {
	t := suite.T()
	ctx := context.Background()
	req := new(proto.GetAnalyzeResultByTokenRequest)
	req.Token = "b92c9636-21a7-46fd-95bf-32d8e34c2d37"
	resp := new(proto.GetAnalyzeResultByTokenResponse)
	err := suite.jinmuHealth.GetAnalyzeResultByToken(ctx, req, resp)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.AnalysisReport)
}

// TestReplaceMatchString 测试replaceMatchString方法
func (suite *RecordTokenSuite) TestReplaceMatchString() {
	t := suite.T()
	test1 := "[高血糖](#comd_15#)肝脏处于极度代偿状态，请立即调整生活及休息，建议结合临床"
	expected1 := "高血糖肝脏处于极度代偿状态，请立即调整生活及休息，建议结合临床"
	assert.Equal(t, test1, expected1)
	test2 := "食肉等是[湿气](#ct0003.0#)的来源"
	expected2 := "食肉等是湿气的来源"
	assert.Equal(t, test2, expected2)
	test3 := "食肉等是[湿气](#ct0003.0#)"
	expected3 := "食肉等是湿气"
	assert.Equal(t, test3, expected3)
	test4 := "[高血糖](#comd_15#)肝脏处于极度代偿[状态]，请立即调整生活及休息，建议结合临床"
	expected4 := "高血糖肝脏处于极度代偿[状态]，请立即调整生活及休息，建议结合临床"
	assert.Equal(t, test4, expected4)
	test5 := "[高血糖](#comd_15#)肝脏处于极度代偿(状态)，请立即调整生活及休息，建议结合临床"
	expected5 := "高血糖肝脏处于极度代偿(状态)，请立即调整生活及休息，建议结合临床"
	assert.Equal(t, test5, expected5)
	test6 := "[高血糖](#comd_15#)肝脏处于极度代偿(状态)[高血糖](#comd_15#)，请立即调整生活及休息，建议结合临床"
	expected6 := "高血糖肝脏处于极度代偿(状态)高血糖，请立即调整生活及休息，建议结合临床"
	assert.Equal(t, test6, expected6)
}
