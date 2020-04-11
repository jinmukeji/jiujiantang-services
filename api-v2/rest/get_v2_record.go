package rest

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	analysispb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/analysis/v1"
	"github.com/kataras/iris/v12"
)

func (h *v2Handler) GetV2AnalyzeReportByRecordID(ctx iris.Context) {
	recordID, _ := ctx.Params().GetInt("record_id")
	req := new(analysispb.GetAnalyzeResultByRecordIDRequest)
	req.RecordId = int32(recordID)
	// req.Cid = rest.GetCidFromContext(ctx)
	resp, err := h.rpcAnalysisSvc.GetAnalyzeResultByRecordID(
		newRPCContext(ctx), req,
	)
	if err != nil {
		writeError(ctx, wrapError(ErrRPCInternal, "", err), false)
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
