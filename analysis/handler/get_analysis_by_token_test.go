package handler

import (
	"net/http"
	"testing"

	"github.com/golang/protobuf/ptypes"
	analysispb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/analysis/v1"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/analysis/v1"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// GetGetAnalyzeResultByTokenTestSuite 是 GetGetAnalyzeResultByTokenTestSuite 的单元测试的 Test Suite
type GetGetAnalyzeResultByTokenTestSuite struct {
	suite.Suite
	analysisManagerService *AnalysisManagerService
}

// SetupSuite 设置测试环境
func (suite *GetGetAnalyzeResultByTokenTestSuite) SetupSuite() {
	suite.analysisManagerService = newAnalysisManagerServiceForTest()
}

// TestGetAnalyzeResultContains 测试有问答的GetAnalyzeResult
func (suite *GetGetAnalyzeResultByTokenTestSuite) TestGetGetAnalyzeResultByToken() {
	t := suite.T()
	reqSignIn, ctx := getSignInReq()
	respSignIn, err := suite.analysisManagerService.jinmuidSvc.UserSignInByUsernamePassword(ctx, reqSignIn)
	assert.Nil(t, err)
	ctx = metadata.NewContext(ctx, map[string]string{
		http.CanonicalHeaderKey(ClientIDKey):    ClientIDsForTest,
		http.CanonicalHeaderKey(AccessTokenKey): respSignIn.AccessToken,
	})

	req, resp := new(proto.GetAnalyzeResultByTokenRequest), new(proto.GetAnalyzeResultByTokenResponse)
	req.Token = "b7527995-5ea5-400e-91a5-72c8d7e90de8"
	req.Cid = "Cid"
	assert.NoError(t, suite.analysisManagerService.GetAnalyzeResultByToken(ctx, req, resp))

	physicalTherapyIndex := &analysispb.PhysicalTherapyIndexModule{}
	err = ptypes.UnmarshalAny(resp.Report.Modules["physical_therapy_index"], physicalTherapyIndex)
	assert.Nil(t, err)
	assert.Equal(t, int32(65), physicalTherapyIndex.GetF0().GetValue())
	assert.Equal(t, int32(50), physicalTherapyIndex.GetF1().GetValue())
	assert.Equal(t, int32(5), physicalTherapyIndex.GetF2().GetValue())
	assert.Equal(t, int32(0), physicalTherapyIndex.GetF3().GetValue())
}

// TestGetGetAnalyzeResultByTokenTestSuite 启动 TestSuite
func TestGetGetAnalyzeResultByTokenTestSuite(t *testing.T) {
	suite.Run(t, new(GetGetAnalyzeResultByTokenTestSuite))
}
