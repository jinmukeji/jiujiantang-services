package rest

import (
	"fmt"
	"time"

	"github.com/jinmukeji/gf-api2/pkg/rest"
	analysispb "github.com/jinmukeji/proto/gen/micro/idl/jm/analysis/v1"
	generalpb "github.com/jinmukeji/proto/gen/micro/idl/ptypes/v2"
	ptypesv2 "github.com/jinmukeji/proto/gen/micro/idl/ptypes/v2"
	"github.com/kataras/iris/v12"
)

const (
	// GenderMale 男性
	GenderMale = 0
	// GenderFemale 女性
	GenderFemale = 1
	// GenderInvalid 性别非法
	GenderInvalid = -1
)

// Language 语言
type Language string

const (
	// LanguageSimpleChinese 简体中文
	LanguageSimpleChinese Language = "zh-Hans"
	// LanguageTraditionalChinese 繁体中文
	LanguageTraditionalChinese Language = "zh-Hant"
	// LanguageEnglish 英文
	LanguageEnglish Language = "en"
)

// GetMonthlyReportBody 请求月报的body
type GetMonthlyReportBody struct {
	Language Language `json:"language"`
}

// GetMonthlyReport 月报
func (h *v2Handler) GetMonthlyReport(ctx iris.Context) {
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

	var body GetMonthlyReportBody
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

	// 月历史记录信息
	location, _ := time.LoadLocation(timeZone)
	now := time.Now().In(location)
	// 从当天的23点59分59秒，往前推30天
	monthStart := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, location).AddDate(0, 0, -30)
	// 从当天的23点59分59秒算结束时间
	monthEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, location).AddDate(0, 0, 0)
	statData, err := h.getWeekOrMonthStatData(ctx, userID, monthStart.UTC(), monthEnd.UTC(), MonthStatData)
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
	req := new(analysispb.GetMonthlyAnalyzeResultRequest)
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
	resp, err := h.rpcAnalysisSvc.GetMonthlyAnalyzeResult(
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
		StartTime:     monthStart,
		EndTime:       monthEnd,
	})
}

func (h *v2Handler) getphysicalDialecticsFromLists(ctx iris.Context, cData []*CData, language ptypesv2.Language) []string {
	var physicalDialectics []string
	for _, value := range cData {
		reqGetAnalyzeResult := &analysispb.GetAnalyzeResultByRecordIDRequest{}
		reqGetAnalyzeResult.RecordId = int32(value.RecordId)
		respGetAnalyzeResult, err := h.rpcAnalysisSvc.GetAnalyzeResultByRecordID(newRPCContext(ctx), reqGetAnalyzeResult)
		if err != nil {
			continue
		}
		analysisReportContent, err := getAnalysisModules(respGetAnalyzeResult.GetReport().GetModules())
		if err != nil {
			continue
		}
		for _, value := range analysisReportContent.PhysicalDialectics.Lookups {
			physicalDialectics = append(physicalDialectics, value.Content)
		}

	}
	return physicalDialectics
}

// mapProtoGenderToRest 将 proto 类型的 gender 转换为 rest 的 int64 类型
func mapProtoGenderToRest(gender generalpb.Gender) (int64, error) {
	switch gender {
	case generalpb.Gender_GENDER_INVALID:
		return GenderInvalid, fmt.Errorf("invalid proto gender %d", gender)
	case generalpb.Gender_GENDER_UNSET:
		return GenderInvalid, fmt.Errorf("invalid proto gender %d", gender)
	case generalpb.Gender_GENDER_MALE:
		return GenderMale, nil
	case generalpb.Gender_GENDER_FEMALE:
		return GenderFemale, nil
	}
	return GenderInvalid, fmt.Errorf("invalid proto gender %d", gender)
}
