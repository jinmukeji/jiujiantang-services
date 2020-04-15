package handler

import (
	"context"

	analysispb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/analysis/v1"
	generalpb "github.com/jinmukeji/proto/v3/gen/micro/idl/ptypes/v2"
)

// UpdateAnalyzeStatus 更新分析的状态
func (j *AnalysisManagerService) UpdateAnalyzeStatus(ctx context.Context, req *analysispb.UpdateAnalyzeStatusRequest, resp *analysispb.UpdateAnalyzeStatusResponse) error {
	// TODO: 无用的rpc，以后删除
	return nil
}

// ConvertToStandingCValues 转化站姿中的C0-C7
func ConvertToStandingCValues(gender generalpb.Gender, c0, c1, c2, c3, c4, c5, c6, c7 int) (sc0, sc1, sc2, sc3, sc4, sc5, sc6, sc7 int) {
	if c2 >= 7 && c2 <= 10 {
		c2 = c2 - 3
	} else if c2 >= 4 && c2 <= 6 {
		c2 = c2 - 2
	} else if c2 >= 2 && c2 <= 3 {
		c2 = c2 - 1
	}
	c4 = c4 + 1
	c0 = c0 + 1
	if gender == generalpb.Gender_GENDER_FEMALE {
		c3 = c3 - 1
	}
	return IntValBoundedBy10(c0), IntValBoundedBy10(c1), IntValBoundedBy10(c2), IntValBoundedBy10(c3), IntValBoundedBy10(c4), IntValBoundedBy10(c5), IntValBoundedBy10(c6), IntValBoundedBy10(c7)
}

// IntValBoundedBy10 返回 -10 到 10 之间的整数
func IntValBoundedBy10(val int) int {
	if val < -10 {
		return -10
	}
	if val > 10 {
		return 10
	}
	return val
}

// ConvertToStandingCFloat64Values 转化成站姿的C0-C7值
func ConvertToStandingCFloat64Values(gender generalpb.Gender, c0, c1, c2, c3, c4, c5, c6, c7 int) (sc0, sc1, sc2, sc3, sc4, sc5, sc6, sc7 float64) {
	rsc0, rsc1, rsc2, rsc3, rsc4, rsc5, rsc6, rsc7 := ConvertToStandingCValues(gender, c0, c1, c2, c3, c4, c5, c6, c7)
	return float64(rsc0), float64(rsc1), float64(rsc2), float64(rsc3), float64(rsc4), float64(rsc5), float64(rsc6), float64(rsc7)
}

// AnalysisReportRequestBody 分析报告请求的body
type AnalysisReportRequestBody struct {
	TransactionID      string                              `json:"transaction_id"`
	QuestionAnswers    map[string]AnalysisReportAnswers    `json:"question_answers"`
	Language           Language                            `json:"language"`
	PhysicalDialectics []AnalysisReportRequestBodyInputKey `json:"physical_dialectics"`
	Disease            []AnalysisReportRequestBodyInputKey `json:"disease"`
	DirtyDialectic     []AnalysisReportRequestBodyInputKey `json:"dirty_dialectic"`
}

// AnalysisReportRequestBodyInputKey 分析报告请求body中的InputKey
type AnalysisReportRequestBodyInputKey struct {
	Key   string `json:"key"`
	Score int32  `json:"score"`
}

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

// AnalysisReportAnswers 分析报告的回答
type AnalysisReportAnswers []AnalysisReportAnswer

// AnalysisReportAnswer 回答
type AnalysisReportAnswer struct {
	QuestionKey string   `json:"question_key"`
	AnswerKeys  []string `json:"answer_keys"`
}
