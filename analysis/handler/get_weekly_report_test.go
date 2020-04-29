package handler

import (
	"net/http"
	"testing"

	"github.com/golang/protobuf/ptypes"
	analysispb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/analysis/v1"
	ptypesv2 "github.com/jinmukeji/proto/v3/gen/micro/idl/ptypes/v2"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// WeeklyReportTestSuite 是 weekly_report 的单元测试的 Test Suite
type WeeklyReportTestSuite struct {
	suite.Suite
	analysisManagerService *AnalysisManagerService
}

// SetupSuite 设置测试环境
func (suite *WeeklyReportTestSuite) SetupSuite() {
	suite.analysisManagerService = newAnalysisManagerServiceForTest()
}

// TestGetWeeklyReportContent 测试GetWeeklyReportContent
func (suite *WeeklyReportTestSuite) TestGetWeeklyReportContent() {
	t := suite.T()
	reqSignIn, ctx := getSignInReq()
	userId := 1
	respSignIn, err := suite.analysisManagerService.jinmuidSvc.UserSignInByUsernamePassword(ctx, reqSignIn)
	ctx = metadata.NewContext(ctx, map[string]string{
		http.CanonicalHeaderKey(AccessTokenKey): respSignIn.AccessToken,
	})
	assert.Nil(t, err)
	reqGetWeeklyReportContent := &analysispb.GetWeeklyAnalyzeResultRequest{
		UserId:   int32(userId),
		Language: ptypesv2.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		Cid:      "cid",
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
	respGetWeeklyReportContent := &analysispb.GetWeeklyAnalyzeResultResponse{}
	ctx = metadata.NewContext(ctx, map[string]string{
		http.CanonicalHeaderKey(ClientIDKey):    "jm-10005",
		http.CanonicalHeaderKey(AccessTokenKey): respSignIn.AccessToken,
	})
	err = suite.analysisManagerService.GetWeeklyAnalyzeResult(ctx, reqGetWeeklyReportContent, respGetWeeklyReportContent)
	assert.Nil(t, err)
	physicalTherapyIndex := &analysispb.PhysicalTherapyIndexModule{}
	err = ptypes.UnmarshalAny(respGetWeeklyReportContent.Report.Modules["physical_therapy_index"], physicalTherapyIndex)
	assert.Nil(t, err)
	assert.Equal(t, int32(25), physicalTherapyIndex.GetF0().GetValue())
	assert.Equal(t, int32(100), physicalTherapyIndex.GetF1().GetValue())
	assert.Equal(t, int32(0), physicalTherapyIndex.GetF2().GetValue())
	assert.Equal(t, int32(0), physicalTherapyIndex.GetF3().GetValue())
}

// TestWeeklyReportTestSuite 启动 TestSuite
func TestWeeklyReportTestSuite(t *testing.T) {
	suite.Run(t, new(WeeklyReportTestSuite))
}
