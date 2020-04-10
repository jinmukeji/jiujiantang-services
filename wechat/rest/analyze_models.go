package rest

import "time"

// AnalysisReportBody 分析报告body
type AnalysisReportBody struct {
	Cid             int64    `json:"cid"`
	AnalysisSession string   `json:"analysis_session"`
	Answers         []Answer `json:"answers"`
}

// AnalysisReportData 分析报告数据
type AnalysisReportData struct {
	Cid             int64           `json:"cid"`
	TransactionNo   string          `json:"transaction_no"`
	AnalysisSession string          `json:"analysis_session"`
	AnalysisDone    bool            `json:"analysis_done"`
	Questionnaire   Questionnaire   `json:"questionnaire"`
	AnalysisReport  *AnalysisReport `json:"analysis_report"`
}

// AnalysisReport 分析报告
type AnalysisReport struct {
	ReportVersion string  `json:"report_version"`
	ReportID      string  `json:"report_id"`
	Content       Content `json:"content"`
}

// Content 内容
type Content struct {
	Lead                                      []GeneralExplain       `json:"lead"`
	UserProfile                               ReportUserProfile      `json:"user_profile"`
	MeasurementResult                         MeasurementResult      `json:"measurement_result"`
	Tags                                      []GeneralExplain       `json:"tags"`
	TipsForWoman                              []GeneralExplain       `json:"tips_for_woman"`
	ChannelsAndCollateralsExplains            []GeneralExplain       `json:"channels_and_collaterals_explains"`
	ConstitutionDifferentiationExplains       []GeneralExplain       `json:"constitution_differentiation_explains"`
	SyndromeDifferentiationExplains           []GeneralExplain       `json:"syndrome_differentiation_explains"`
	FactorExplains                            []GeneralExplain       `json:"factor_explains"`
	DictionaryEntries                         []GeneralExplain       `json:"dictionary_entries"`
	ConstitutionDifferentiationExplainNotices []GeneralExplain       `json:"constitution_differentiation_explain_notices"`
	MeasurementTips                           []GeneralExplain       `json:"measurement_tips"`
	HealthDescriptions                        []GeneralExplain       `json:"health_descriptions"`
	ChannelsAndCollateralsStrength            []CCStrengthItem       `json:"channels_and_collaterals_strength"`
	BabyTips                                  []GeneralExplain       `json:"baby_tips"`
	CCExplainNotices                          []GeneralExplain       `json:"channels_and_collaterals_explain_notices"`
	PTExplain                                 PhysicalTherapyExplain `json:"physical_therapy_explain"`

	UterineHealthIndexes                  []GeneralExplain `json:"uterine_health_indexes"`
	UterusAttentionPrompts                []GeneralExplain `json:"uterus_attention_prompts"`
	UterineHealthDescriptions             []GeneralExplain `json:"uterine_health_descriptions"`
	MenstrualHealthValues                 []GeneralExplain `json:"menstrual_health_values"`
	MenstrualHealthDescriptions           []GeneralExplain `json:"menstrual_health_descriptions"`
	GynecologicalInflammations            []GeneralExplain `json:"gynecological_inflammations"`
	GynecologicalInflammationDescriptions []GeneralExplain `json:"gynecological_inflammation_descriptions"`
	BreastHealth                          []GeneralExplain `json:"breast_health"`
	BreastHealthDescriptions              []GeneralExplain `json:"breast_health_descriptions"`
	EmotionalHealthIndexes                []GeneralExplain `json:"emotional_health_indexes"`
	EmotionalHealthDescriptions           []GeneralExplain `json:"emotional_health_descriptions"`
	FacialSkins                           []GeneralExplain `json:"facial_skins"`
	FacialSkinDescriptions                []GeneralExplain `json:"facial_skin_descriptions"`
	ReproductiveAgeConsiderations         []GeneralExplain `json:"reproductive_age_considerations"`
	BreastCancerOvarianCancers            []GeneralExplain `json:"breast_cancer_ovarian_cancers"`
	BreastCancerOvarianCancerDescriptions []GeneralExplain `json:"breast_cancer_ovarian_cancer_descriptions"`
	HormoneLevels                         []GeneralExplain `json:"hormone_levels"`
	LymphaticHealth                       []GeneralExplain `json:"lymphatic_health"`
	LymphaticHealthDescriptions           []GeneralExplain `json:"lymphatic_health_descriptions"`
	M0                                    *int32           `json:"m0"`
	M1                                    *int32           `json:"m1"`
	M2                                    *int32           `json:"m2"`
	M3                                    *int32           `json:"m3"`
	F100                                  int32            `json:"f100"`
	F101                                  int32            `json:"f101"`
	F102                                  int32            `json:"f102"`
	F103                                  int32            `json:"f103"`
	F104                                  int32            `json:"f104"`
	F105                                  int32            `json:"f105"`
	F106                                  int32            `json:"f106"`
	F107                                  int32            `json:"f107"`
	HasPaid                               bool             `json:"has_paid"`
	Options                               DisplayOptions   `json:"display_options"`
	Remark                                string           `json:"remark"`
	CreatedAt                             time.Time        `json:"created_at"`
}

// PhysicalTherapyExplain 理疗指数
type PhysicalTherapyExplain struct {
	F0 int32 `json:"f0"`
	F1 int32 `json:"f1"`
	F2 int32 `json:"f2"`
	F3 int32 `json:"f3"`
}

// CCStrengthLabel 经络强度分析的项目的标签
type CCStrengthLabel struct {
	Label string `json:"label"` // 文本标签，可以是HTML
	CC    string `json:"cc"`    // 经络标识
}

// CCStrengthItem 经络强度分析的项目的Item
type CCStrengthItem struct {
	Key      string            `json:"key"`                // 编号，即唯一标识
	Labels   []CCStrengthLabel `json:"labels,flow"`        // 标签清单
	Disabled bool              `json:"disabled,omitempty"` // 是否已经弃用
	Remark   string            `json:"remark,omitempty"`   // 备注
}

// GeneralExplain 一般解释
type GeneralExplain struct {
	Key     string `json:"key"`
	Label   string `json:"label"`
	Content string `json:"content"`
}

// MeasurementResult 测量结果
type MeasurementResult struct {
	Finger              int64   `json:"finger"`
	C0                  int64   `json:"c0"`
	C1                  int64   `json:"c1"`
	C2                  int64   `json:"c2"`
	C3                  int64   `json:"c3"`
	C4                  int64   `json:"c4"`
	C5                  int64   `json:"c5"`
	C6                  int64   `json:"c6"`
	C7                  int64   `json:"c7"`
	HeartRate           int64   `json:"heart_rate"`
	AppHeartRate        int64   `json:"app_heart_rate"`
	PartialPulseWave    []int32 `json:"partial_pulse_wave"`
	AppHighestHeartRate int32   `json:"app_highest_heart_rate"` // app最高心率
	AppLowestHeartRate  int32   `json:"app_lowest_heart_rate"`  // app最低心率
}

// ReportUserProfile 分析的用户档案
type ReportUserProfile struct {
	UserID    int64     `json:"user_id"`
	Nickname  string    `json:"nickname"`
	Birthday  time.Time `json:"birthday"`
	Age       int64     `json:"age"`
	Gender    int64     `json:"gender"`
	Height    int64     `json:"height"`
	Weight    int64     `json:"weight"`
	AvatarURL string    `json:"avatar_url"`
}

// Questionnaire 问卷
type Questionnaire struct {
	Title     string     `json:"title"`
	Questions []Question `json:"questions"`
	Answers   []Answer   `json:"answers"`
	CreateAt  time.Time  `json:"create_at"`
}

// Answer 答案
type Answer struct {
	QuestionKey string   `json:"question_key"`
	Values      []string `json:"values"`
}

// Question 问题
type Question struct {
	Key         string   `json:"key"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tip         string   `json:"tip"`
	Type        string   `json:"type"`
	Choices     []Choice `json:"choices"`
	DefaultKeys []string `json:"default_keys"`
}

// Choice 选项
type Choice struct {
	Key          string   `json:"key"`
	Name         string   `json:"name"`
	Value        string   `json:"value"`
	ConflictKeys []string `json:"conflict_keys"`
}

// DisplayOptions 显示配置
type DisplayOptions struct {
	DisplayNavbar                 bool `json:"display_navbar"`
	DisplayTags                   bool `json:"display_tags"`
	DisplayPartialData            bool `json:"display_partial_data"`
	DisplayUserProfile            bool `json:"display_user_profile"`
	DisplayHeartRate              bool `json:"display_heart_rate"`
	DisplayCcBarChart             bool `json:"display_cc_bar_chart"`
	DisplayCcExplain              bool `json:"display_cc_explain"`
	DisplayCdExplain              bool `json:"display_cd_explain"`
	DisplaySdExplain              bool `json:"display_sd_explain"`
	DisplayF0                     bool `json:"display_f0"`
	DisplayF1                     bool `json:"display_f1"`
	DisplayF2                     bool `json:"display_f2"`
	DisplayF3                     bool `json:"display_f3"`
	DisplayPhysicalTherapyExplain bool `json:"display_physical_therapy_explain"`
	DisplayRemark                 bool `json:"display_remark"`
	DisplayMeasurementResult      bool `json:"display_measurement_result"`
	DisplayBabyTips               bool `json:"display_baby_tips"`
	DisplayWh                     bool `json:"display_wh"`
}
