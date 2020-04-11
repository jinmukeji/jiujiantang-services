package handler

import (
	"net/http"
	"testing"

	"github.com/golang/protobuf/ptypes"
	analysispb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/analysis/v1"
	ptypesv2 "github.com/jinmukeji/proto/v3/gen/micro/idl/ptypes/v2"
	"github.com/micro/go-micro/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// MonthlyReportTestSuite 是 monthly_report 的单元测试的 Test Suite
type MonthlyReportTestSuite struct {
	suite.Suite
	analysisManagerService *AnalysisManagerService
}

// SetupSuite 设置测试环境
func (suite *MonthlyReportTestSuite) SetupSuite() {
	suite.analysisManagerService = newAnalysisManagerServiceForTest()
}

// TestGetMonthlyAnalyzeResult 测试TestGetMonthlyAnalyzeResult
func (suite *MonthlyReportTestSuite) TestGetMonthlyAnalyzeResult() {
	t := suite.T()
	reqSignIn, ctx := getSignInReq()
	userId := 1
	respSignIn, err := suite.analysisManagerService.jinmuidSvc.UserSignInByUsernamePassword(ctx, reqSignIn)
	assert.Nil(t, err)
	ctx = metadata.NewContext(ctx, map[string]string{
		http.CanonicalHeaderKey(ClientIDKey):    ClientIDsForTest,
		http.CanonicalHeaderKey(AccessTokenKey): respSignIn.AccessToken,
	})
	reqGetMonthlyAnalyzeResult := &analysispb.GetMonthlyAnalyzeResultRequest{
		UserId:   int32(userId),
		Language: ptypesv2.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		Cid:      "Cid",
		CInfo: &analysispb.CInfo{
			C0: 2,
			C1: 4,
			C2: 2,
			C3: 1,
			C4: 0,
			C5: 6,
			C6: 4,
			C7: 2,
		},
		PhysicalDialectics: []string{"T0017", "TZN0001", "T0017", "TZN0001", "T0017", "TZN0001", "T0017", "TZN0001", "T0017", "TZN0001", "T0017", "TZN0001", "T0017", "TZN000"},
	}
	respGetMonthlyAnalyzeResult := &analysispb.GetMonthlyAnalyzeResultResponse{}
	err = suite.analysisManagerService.GetMonthlyAnalyzeResult(ctx, reqGetMonthlyAnalyzeResult, respGetMonthlyAnalyzeResult)
	assert.Nil(t, err)
	physicalTherapyIndex := &analysispb.PhysicalTherapyIndexModule{}
	err = ptypes.UnmarshalAny(respGetMonthlyAnalyzeResult.Report.Modules["physical_therapy_index"], physicalTherapyIndex)
	assert.Nil(t, err)
	assert.Equal(t, int32(25), physicalTherapyIndex.GetF0().GetValue())
	assert.Equal(t, int32(100), physicalTherapyIndex.GetF1().GetValue())
	assert.Equal(t, int32(0), physicalTherapyIndex.GetF2().GetValue())
	assert.Equal(t, int32(0), physicalTherapyIndex.GetF3().GetValue())
}

// TestMonthlyReportTestSuite 启动 TestSuite
func TestMonthlyReportTestSuite(t *testing.T) {
	suite.Run(t, new(WeeklyReportTestSuite))
}
