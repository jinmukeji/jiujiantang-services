package handler

import (
	"net/http"
	"testing"

	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/analysis/v1"
	ptypesv2 "github.com/jinmukeji/proto/v3/gen/micro/idl/ptypes/v2"
	"github.com/micro/go-micro/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// AnalysisResultTestSuite 是 AnalysisResultTestSuite 的单元测试的 Test Suite
type AnalysisResultTestSuite struct {
	suite.Suite
	analysisManagerService *AnalysisManagerService
}

// SetupSuite 设置测试环境
func (suite *AnalysisResultTestSuite) SetupSuite() {
	suite.analysisManagerService = newAnalysisManagerServiceForTest()
}

// TestGetAnalyzeResultContains 测试有问答的GetAnalyzeResult
func (suite *AnalysisResultTestSuite) TestGetAnalyzeResultContains() {
	t := suite.T()
	reqSignIn, ctx := getSignInReq()
	respSignIn, err := suite.analysisManagerService.jinmuidSvc.UserSignInByUsernamePassword(ctx, reqSignIn)
	assert.Nil(t, err)
	ctx = metadata.NewContext(ctx, map[string]string{
		http.CanonicalHeaderKey(ClientIDKey):    ClientIDsForTest,
		http.CanonicalHeaderKey(AccessTokenKey): respSignIn.AccessToken,
	})

	req, resp := new(proto.GetAnalyzeResultRequest), new(proto.GetAnalyzeResultResponse)
	req.Language = ptypesv2.Language_LANGUAGE_SIMPLIFIED_CHINESE
	req.RecordId = 493038
	req.TransactionId = "TransactionId"
	req.QuestionAnswers = nil
	req.IsSkipVerifyToken = true
	req.Cid = "Cid"
	assert.NoError(t, suite.analysisManagerService.GetAnalyzeResult(ctx, req, resp))
	assert.Nil(t, resp.Report)
	assert.NotEqual(t, nil, resp.Questions)
	assert.Equal(t, "Q0001", resp.Questions["stress_state_judgment"].Questions[0].GetQuestionKey())
	assert.Equal(t, "single_choice", resp.Questions["stress_state_judgment"].Questions[0].GetType())
	assert.Equal(t, "QC0001", resp.Questions["stress_state_judgment"].Questions[0].GetChoices()[0].GetChoiceKey())
	assert.Equal(t, "QC0002", resp.Questions["stress_state_judgment"].Questions[0].GetChoices()[1].GetChoiceKey())
}

// TestAnalysisResultTestSuite 启动 TestSuite
func TestAnalysisResultTestSuite(t *testing.T) {
	suite.Run(t, new(AnalysisResultTestSuite))
}
