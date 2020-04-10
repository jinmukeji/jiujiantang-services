package handler

// AnalysisReportRequestBody 分析报告请求的body
type AnalysisReportRequestBody struct {
	TransactionID      string                              `json:"transaction_id"`
	QuestionAnswers    map[string]AnalysisReportAnswers    `json:"question_answers"`
	Language           Language                            `json:"language"`
	PhysicalDialectics []AnalysisReportRequestBodyInputKey `json:"physical_dialectics"`
	Disease            []AnalysisReportRequestBodyInputKey `json:"disease"`
	DirtyDialectic     []AnalysisReportRequestBodyInputKey `json:"dirty_dialectic"`
}

// AnalysisReportRequestBodyInputKey 分析报告请求的body中的InputKey
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
