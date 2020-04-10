package handler

import (
	"path/filepath"
	"testing"

	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// SubmitMeasurementStatusTestSuite 是 SubmitMeasurementStatus rpc 的单元测试的 Test Suite
type SubmitMeasurementStatusTestSuite struct {
	suite.Suite
	jinmuHealth *JinmuHealth
}

// SetupSuite 设置测试环境
func (suite *SubmitMeasurementStatusTestSuite) SetupSuite() {
	suite.jinmuHealth = new(JinmuHealth)
	envFilepath := filepath.Join("testdata", "local.svc-biz-core.env")
	suite.jinmuHealth.datastore, _ = newTestingDbClientFromEnvFile(envFilepath)
	suite.jinmuHealth.mailClient, _ = newTestingMailClientFromEnvFile(envFilepath)
}

// TestMeasurementStatus测试提交测量时状态
func (suite *SubmitMeasurementStatusTestSuite) TestMeasurementStatus() {
	t := suite.T()
	// mock 一次登录
	var password = "release1"
	var testRecordID int32 = 1
	ctx, _ := mockLogin(suite.jinmuHealth, testUsername, password)
	req, resp := new(proto.SubmitMeasurementStatusRequest), new(proto.SubmitMeasurementStatusResponse)
	req.RecordId = testRecordID
	req.Lactation = proto.Status_STATUS_UNSELECTED_STATUS
	req.Pregnancy = proto.Status_STATUS_SELECTED_STATUS
	assert.NoError(t, suite.jinmuHealth.SubmitMeasurementStatus(ctx, req, resp))
	// 查找这个 record
}

func TestSubmitMeasurementStatusTestSuite(t *testing.T) {
	suite.Run(t, new(SubmitMeasurementStatusTestSuite))
}
