package rest

import (
	"fmt"
	"time"

	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	analysispb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/analysis/v1"
	"github.com/kataras/iris/v12"
)

// GetWeeklyReportBody 请求周报的body
type GetWeeklyReportBody struct {
	Language Language `json:"language"`
}

// GetWeeklyReport 周报
func (h *v2Handler) GetWeeklyReport(ctx iris.Context) {
	userID, err := ctx.Params().GetInt("user_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	if ctx.Values().GetString(ClientIDKey) == seamlessClient {
		writeError(
			ctx,
			wrapError(ErrDeniedToAccessAPI, "", fmt.Errorf("%s is denied to access this API", seamlessClient)),
			false,
		)
		return
	}

	var body GetWeeklyReportBody
	errReadJSON := ctx.ReadJSON(&body)
	if errReadJSON != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", errReadJSON), false)
		return
	}
	language, err := mapRestLanguageToProto(body.Language)
	if err != nil {
		writeError(ctx, wrapError(ErrValueRequired, "", err), false)
		return
	}
	timeZone := getTimeZone(ctx)

	// 周历史记录信息
	location, _ := time.LoadLocation(timeZone)
	now := time.Now().In(location)
	// 从当天的23点59分59秒，往前推7天
	weekStart := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, location).AddDate(0, 0, -7)
	// 从当天的23点59分59秒算结束时间
	weekEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, location).AddDate(0, 0, 0)
	statData, err := h.getWeekOrMonthStatData(ctx, userID, weekStart.UTC(), weekEnd.UTC(), WeekStatData)
	if err != nil {
		writeError(ctx, wrapError(ErrRPCInternal, "", err), false)
		return
	}
	// 错误提示不为空，表示记录信息不全，返回错误提示
	if statData.ErrorMessage != "" {
		rest.WriteOkJSON(ctx, WeeklyOrMonthlyReportResponse{
			ErrorMessage: statData.ErrorMessage,
		})
		return
	}
	req := new(analysispb.GetWeeklyAnalyzeResultRequest)
	req.UserId = int32(userID)
	req.Language = language
	req.Cid = rest.GetCidFromContext(ctx)
	req.CInfo = &analysispb.CInfo{
		C0: statData.AverageMeridian.C0,
		C1: statData.AverageMeridian.C1,
		C2: statData.AverageMeridian.C2,
		C3: statData.AverageMeridian.C3,
		C4: statData.AverageMeridian.C4,
		C5: statData.AverageMeridian.C5,
		C6: statData.AverageMeridian.C6,
		C7: statData.AverageMeridian.C7,
	}
	req.PhysicalDialectics = statData.PhysicalDialectics
	resp, err := h.rpcAnalysisSvc.GetWeeklyAnalyzeResult(
		newRPCContext(ctx), req,
	)
	if err != nil {
		writeRPCInternalError(ctx, err, false)
		return
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

	rest.WriteOkJSON(ctx, WeeklyOrMonthlyReportResponse{
		ReportVersion: resp.ReportVersion,
		ReportContent: analysisReportContent,
		StartTime:     weekStart,
		EndTime:       weekEnd,
	})

}
