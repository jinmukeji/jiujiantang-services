package rest

// AnalysisReportRequestBody 分析报告的body
type AnalysisReportRequestBody struct {
	TransactionID   string                           `json:"transaction_id"`
	QuestionAnswers map[string]AnalysisReportAnswers `json:"question_answers"`
	Language        Language                         `json:"language"`
}

// AnalysisReportAnswers 分析报告的答案
type AnalysisReportAnswers []AnalysisReportAnswer

// AnalysisReportAnswer 答案
type AnalysisReportAnswer struct {
	QuestionKey string   `json:"question_key"`
	AnswerKeys  []string `json:"answer_keys"`
}
