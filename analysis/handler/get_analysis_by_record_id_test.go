package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/ae/v2/biz"
	"github.com/jinmukeji/jiujiantang-services/analysis/aws"
	jinmuidpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	analysispb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/analysis/v1"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/analysis/v1"
	"github.com/joho/godotenv"
	"github.com/micro/go-micro/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
)

// GetAnalyzeResultByRecordIDTestSuite 是 GetAnalyzeResultByRecordIDTestSuite 的单元测试的 Test Suite
type GetAnalyzeResultByRecordIDTestSuite struct {
	suite.Suite
	analysisManagerService *AnalysisManagerService
}

func newAnalysisManagerServiceForTest() *AnalysisManagerService {
	envFilepath := filepath.Join("testdata", "local.svc-analysis.env")
	datastore, _ := newTestingDbClientFromEnvFile(envFilepath)
	biz := biz.NewBizEngineManager(
		biz.LuaSrcPath("../../build/assets/ae_data_v2/v2.8.19/lua_src-v2.8.19"),
		biz.TemplatesDir("../../build/assets/ae_data_v2/v2.8.19/lookups-v2.8.19"),
		biz.QuestionDir("../../build/assets/ae_data_v2/v2.8.19/question-v2.8.19"),
		biz.PoolSize(2),
	)
	presetsFilePath := "../../build/assets/ae_data_v2/v2.8.19/biz_conf-v2.8.19/presets.yaml"
	_ = godotenv.Load(envFilepath)
	awsClient, _ := aws.NewClient(
		aws.BucketName(os.Getenv("X_AWS_BUCKET_NAME")),
		aws.AccessKeyID(os.Getenv("X_AWS_ACCESS_KEY")),
		aws.SecretKey(os.Getenv("X_AWS_SECRET_KEY")),
		aws.Region(os.Getenv("X_AWS_REGION")),
		aws.PulseTestRawDataEnvironmentS3KeyPrefix(os.Getenv("X_WAVE_DATA_KEY_PREFIX")),
		aws.PulseTestRawDataS3KeyPrefix(os.Getenv("X_PULSE_TEST_RAW_DATA_S3_KEY_PREFIX")),
	)
	return NewAnalysisManagerService(datastore, biz, presetsFilePath, awsClient)
}

func getSignInReq() (*jinmuidpb.UserSignInByUsernamePasswordRequest, context.Context) {
	signInUsername := "1"
	hashedPassword := "9626c7444717aab7a3bbdd509bcafa35a7491e9478d421b38e539a621f695edd"
	ctx := context.Background()
	ctx = metadata.NewContext(ctx, map[string]string{
		http.CanonicalHeaderKey(ClientIDKey): "jm-10005",
	})
	reqSignIn := &jinmuidpb.UserSignInByUsernamePasswordRequest{
		Username:       signInUsername,
		HashedPassword: hashedPassword,
		Seed:           "",
	}
	return reqSignIn, ctx
}

const (
	ClientIDsForTest = "jm-10005"
)

// SetupSuite 设置测试环境
func (suite *GetAnalyzeResultByRecordIDTestSuite) SetupSuite() {
	suite.analysisManagerService = newAnalysisManagerServiceForTest()
}

// TestGetAnalyzeResultContains 测试有问答的GetAnalyzeResult
func (suite *GetAnalyzeResultByRecordIDTestSuite) TestGetAnalyzeResultByRecordID() {
	t := suite.T()
	reqSignIn, ctx := getSignInReq()
	respSignIn, err := suite.analysisManagerService.jinmuidSvc.UserSignInByUsernamePassword(ctx, reqSignIn)
	assert.Nil(t, err)
	ctx = metadata.NewContext(ctx, map[string]string{
		http.CanonicalHeaderKey(ClientIDKey):    ClientIDsForTest,
		http.CanonicalHeaderKey(AccessTokenKey): respSignIn.AccessToken,
	})
	req, resp := new(proto.GetAnalyzeResultByRecordIDRequest), new(proto.GetAnalyzeResultByRecordIDResponse)
	req.RecordId = 493050
	req.Cid = "Cid"
	assert.NoError(t, suite.analysisManagerService.GetAnalyzeResultByRecordID(ctx, req, resp))

	physicalTherapyIndex := &analysispb.PhysicalTherapyIndexModule{}
	err = ptypes.UnmarshalAny(resp.Report.Modules["physical_therapy_index"], physicalTherapyIndex)
	assert.Nil(t, err)
	assert.Equal(t, int32(65), physicalTherapyIndex.GetF0().GetValue())
	assert.Equal(t, int32(50), physicalTherapyIndex.GetF1().GetValue())
	assert.Equal(t, int32(5), physicalTherapyIndex.GetF2().GetValue())
	assert.Equal(t, int32(0), physicalTherapyIndex.GetF3().GetValue())
}

// TestGetAnalyzeResultByRecordIDTestSuite 启动 TestSuite
func TestGetAnalyzeResultByRecordIDTestSuite(t *testing.T) {
	suite.Run(t, new(GetAnalyzeResultByRecordIDTestSuite))
}
