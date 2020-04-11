package rest

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	analysispb "github.com/jinmukeji/proto/gen/micro/idl/jm/analysis/v1"
	"github.com/kataras/iris/v12"
)

func (h *v2Handler) GetV2AnalyzeReportByToken(ctx iris.Context) {

	token := ctx.Params().Get("token")
	req := new(analysispb.GetAnalyzeResultByTokenRequest)
	req.Token = token
	resp, errGetAnalyzeResultByToken := h.rpcAnalysisSvc.GetAnalyzeResultByToken(
		newRPCContext(ctx), req,
	)
	if errGetAnalyzeResultByToken != nil {
		writeError(ctx, wrapError(ErrRPCInternal, "", errGetAnalyzeResultByToken), false)
		return
	}
	analysisReportResponse := AnalysisReportResponse{
		ReportVersion: resp.GetReportVersion(),
		ReportID:      resp.GetReport().GetRecordId(),
		TransactionID: resp.GetTransactionId(),
	}

	// 与引擎分析结果相关模块
	analysisReportContent, err := getAnalysisModules(resp.GetReport().GetModules())
	if err != nil {
		writeRPCInternalError(ctx, wrapError(ErrRPCInternal, "", err), false)
		return
	}
	// 个人信息模块
	userProfileModule, err := getUserProfileModule(resp.GetReport().GetUserProfile())
	if err != nil {
		writeRPCInternalError(ctx, wrapError(ErrRPCInternal, "", err), false)
		return
	}
	analysisReportContent.UserProfile = userProfileModule
	// 测量上下文模块
	pulseTestModule, err := getPulseTestModule(resp.GetReport().GetPulseTest())
	if err != nil {
		writeRPCInternalError(ctx, wrapError(ErrRPCInternal, "", err), false)
		return
	}
	// 备注模块
	analysisReportContent.PulseTest = pulseTestModule
	remarkModule, err := getRemarkModule(resp.GetReport().GetRemark())
	if err != nil {
		writeRPCInternalError(ctx, wrapError(ErrRPCInternal, "", err), false)
		return
	}
	analysisReportContent.Remark = remarkModule
	// 测量时间
	startTime, _ := ptypes.Timestamp(resp.GetReport().GetCreatedTime())
	analysisReportContent.CreatedTime = startTime

	analysisReportResponse.ReportContent = *analysisReportContent

	rest.WriteOkJSON(ctx, analysisReportResponse)
}
