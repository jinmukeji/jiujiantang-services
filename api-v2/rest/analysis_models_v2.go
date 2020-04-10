package rest

import "time"

// WeeklyOrMonthlyReportResponse 周报或者月报分析的返回
type WeeklyOrMonthlyReportResponse struct {
	ReportVersion string                 `json:"report_version"`
	ReportContent *AnalysisReportContent `json:"report_content"`
	StartTime     time.Time              `json:"start_time"`
	EndTime       time.Time              `json:"end_time"`
	ErrorMessage  string                 `json:"error_message"`
}

// AnalysisReportResponse 分析报告的返回
type AnalysisReportResponse struct {
	ReportVersion string `json:"report_version"`
	ReportID      int32  `json:"report_id"`
	TransactionID string `json:"transaction_id"`
	// 提问的问题
	AskQuestions  map[string]Questions  `json:"ask_questions"`
	ReportContent AnalysisReportContent `json:"report_content"`
}

// AnalysisMonthlyReportRequestBody 月报的body
type AnalysisMonthlyReportRequestBody struct {
	C0                 int32    `json:"c0"`      // 心包经测量指标
	C1                 int32    `json:"c1"`      // 肝经测量指标
	C2                 int32    `json:"c2"`      // 肾经测量指标
	C3                 int32    `json:"c3"`      // 脾经测量指标
	C4                 int32    `json:"c4"`      // 肺经测量指标
	C5                 int32    `json:"c5"`      // 胃经测量指标
	C6                 int32    `json:"c6"`      // 胆经测量指标
	C7                 int32    `json:"c7"`      // 膀胱经测量指标
	UserID             int32    `json:"user_id"` // useriD
	Language           Language `json:"language"`
	PhysicalDialectics []string `json:"physical_dialectics"`
}

// AnalysisWeeklyReportRequestBody 周报的body
type AnalysisWeeklyReportRequestBody struct {
	C0                 int32    `json:"c0"`      // 心包经测量指标
	C1                 int32    `json:"c1"`      // 肝经测量指标
	C2                 int32    `json:"c2"`      // 肾经测量指标
	C3                 int32    `json:"c3"`      // 脾经测量指标
	C4                 int32    `json:"c4"`      // 肺经测量指标
	C5                 int32    `json:"c5"`      // 胃经测量指标
	C6                 int32    `json:"c6"`      // 胆经测量指标
	C7                 int32    `json:"c7"`      // 膀胱经测量指标
	UserID             int32    `json:"user_id"` // useriD
	Language           Language `json:"language"`
	PhysicalDialectics []string `json:"physical_dialectics"`
}

// InputKey 疾病，脏腑，体质的结构体
type InputKey struct {
	Key   string `json:"key"`
	Score int32  `json:"score"`
}

// PresetTemplateData preset模板数据
type PresetTemplateData struct {
	AssistantName string `json:"assistant_name"`
	ReporterName  string `json:"reporter_name"`
}

// Questions 问题
type Questions []AnalysisReportQuestion

// AnalysisReportQuestion 分析报告问题
type AnalysisReportQuestion struct {
	Key     string                 `json:"key"`
	Type    string                 `json:"type"`
	Content string                 `json:"content"`
	Choices []AnalysisReportChoice `json:"choices"`
}

// AnalysisReportChoice 分析报告问题的选项
type AnalysisReportChoice struct {
	Key          string   `json:"key"`
	Content      string   `json:"content"`
	ConflictKeys []string `json:"conflict_keys"`
	Selected     bool     `json:"selected"`
}

// AnalysisMonthlyReportResponse 月报的返回
type AnalysisMonthlyReportResponse struct {
	ReportVersion string               `json:"report_version"`
	Content       MonthlyReportContent `json:"content"`
}

// AnalysisWeeklyReportResponse 周报的返回
type AnalysisWeeklyReportResponse struct {
	ReportVersion string              `json:"report_version"`
	Content       WeeklyReportContent `json:"content"`
}

type RemarkModule struct {
	// 是否显示
	Enabled bool `json:"enabled"`
	// 备注内容
	Content string `json:"content"`
}

type HeartRateModule struct {
	// 是否显示
	Enabled bool `json:"enabled"`
	// 平均心率
	AverageHeartRate int32 `json:"average_heart_rate"`
	// 最高心率
	HighestHeartRate int32 `json:"highest_heart_rate"`
	// 最低心率
	LowestHeartRate int32 `json:"lowest_heart_rate"`
}

// RiskEstimateModule 风险预估模块
type RiskEstimateModule struct {
	// 是否显示该模块
	Enabled bool `json:"enabled"`
	// 疾病分值
	DiseaseEstimate []*Lookup `json:"disease_estimate"`
	// 提示信息
	PromptMessage []*Lookup `json:"prompt_message"`
}

// Lookup 模块对应的 key 和 value.
type Lookup struct {
	// 序号
	Key string `json:"key"`
	// 显示内容
	Content string `json:"content"`
	// 得分或排序值
	Score float64 `json:"score"`
	// 词条链接 Key
	LinkKey string `json:"link_key"`
}

// PhysicalDialecticsModule 体质辩证模块.
type PhysicalDialecticsModule struct {
	// 是否显示该模块
	Enabled bool      `json:"enabled"`
	Lookups []*Lookup `json:"lookups"`
}

// DirtyDialecticModule 脏腑辩证模块.
type DirtyDialecticModule struct {
	// 是否显示该模块
	Enabled bool      `json:"enabled"`
	Lookups []*Lookup `json:"lookups"`
}

// PhysicalTherapyIndexModule 理疗指数模块.
type PhysicalTherapyIndexModule struct {
	// 是否显示该模块
	Enabled bool `json:"enabled"`
	// 阴虚指数
	F0 int32 `json:"f0"`
	// 阳虚指数
	F1 int32 `json:"f1"`
	// 湿气指数
	F2 int32 `json:"f2"`
	// 血瘀指数
	F3      int32     `json:"f3"`
	Lookups []*Lookup `json:"lookups"`
}

// ConditioningAdviceModule 调理建议模块.
type ConditioningAdviceModule struct {
	// 是否显示该模块
	Enabled bool `json:"enabled"`
	// 食疗建议
	DietaryAdvice *DietaryAdviceModule `json:"dietary_advice"`
	// 运动方案
	SportsAdvice *SportsAdviceModule `json:"sports_advice"`
	// 中药调理建议
	ChineseMedicineAdvice *ChineseMedicineAdviceModule `json:"chinese_medicine_advice"`
	// 理疗建议
	PhysicalTherapyAdvice *PhysicalTherapyAdviceModule `json:"physical_therapy_advice"`
}

// DietaryAdviceModule 食疗建议模块.
type DietaryAdviceModule struct {
	// 是否显示该模块
	Enabled bool      `json:"enabled"`
	Lookups []*Lookup `json:"lookups"`
}

// SportsAdviceModule 运动方案模块.
type SportsAdviceModule struct {
	// 是否显示该模块
	Enabled bool      `json:"enabled"`
	Lookups []*Lookup `json:"lookups"`
}

// ChineseMedicineAdviceModule 中药调理建议模块.
type ChineseMedicineAdviceModule struct {
	// 是否显示该模块
	Enabled bool      `json:"enabled"`
	Lookups []*Lookup `json:"lookups"`
}

// PhysicalTherapyAdviceModule 理疗建议模块.
type PhysicalTherapyAdviceModule struct {
	// 是否显示该模块
	Enabled bool      `json:"enabled"`
	Lookups []*Lookup `json:"lookups"`
}

// PartialPulseWaveModule 局部脉搏波模块.
type PartialPulseWaveModule struct {
	// 是否显示该模块
	Enabled bool    `json:"enabled"`
	Points  []int32 `json:"points"`
}

// MeridianExplainModule 经络解读模块.
type MeridianExplainModule struct {
	// 是否显示该模块
	Enabled bool      `json:"enabled"`
	Lookups []*Lookup `json:"lookups"`
}

// Measurement judgment.
// 测量判断.
type MeasurementJudgmentModule struct {
	Enabled bool      `json:"enabled"`
	Lookups []*Lookup `json:"lookups"`
}

// TipsModule 温馨提示模块.
type TipsModule struct {
	Enabled bool      `json:"enabled"`
	Lookups []*Lookup `json:"lookups"`
}

// StressStateJudgmentModule 应激态模块.
type StressStateJudgmentModule struct {
	Enabled                bool `json:"enabled"`
	HasStressState         bool `json:"has_stress_state"`
	HasDoneSports          bool `json:"has_done_sports"`
	HasDrinkedWine         bool `json:"has_drinked_wine"`
	HasHadCold             bool `json:"has_had_cold"`
	HasRhinitisEpisode     bool `json:"has_rhinitis_episode"`
	HasAbdominalPain       bool `json:"has_abdominal_pain"`
	HasViralInfection      bool `json:"has_viral_infection"`
	HasPhysiologicalPeriod bool `son:"has_physiological_period"`
	HasOvulation           bool `json:"has_ovulation"`
	HasPregnant            bool `json:"has_pregnant"`
	HasLactation           bool `json:"has_lactations"`
}

// MeridianBarChartModule 经络柱状图模块.
type MeridianBarChartModule struct {
	// 是否显示该模块
	Enabled bool `json:"enabled"`
	// 经络值
	MeridianValue *CInfo `json:"meridian_value"`
}

// CInfo 经络的值和测量时间.
type CInfo struct {
	C0       int32     `json:"c0"`
	C1       int32     `json:"c1"`
	C2       int32     `json:"c2"`
	C3       int32     `json:"c3"`
	C4       int32     `json:"c4"`
	C5       int32     `json:"c5"`
	C6       int32     `json:"c6"`
	C7       int32     `json:"c7"`
	TestTime time.Time `json:"test_time"`
}

// PulseTestModule 测量上下文.
type PulseTestModule struct {
	// 采样使用的手指
	Fingers []int `json:"fingers"`
}

// EmotionalHealthModule 情绪健康模块.
type EmotionalHealthModule struct {
	Enabled bool      `json:"enabled"`
	Lookups []*Lookup `json:"lookups"`
	// 情绪健康指数
	F103 int `json:"f103"`
}

// FacialSkinModule 面部美肤模块.
type FacialSkinModule struct {
	Enabled bool      `json:"enabled"`
	Lookups []*Lookup `json:"lookups"`
	// 面部美肤指数
	F104 int `json:"f104"`
}

// FacialSkinMaleModule 男性面部美肤模块.
type FacialSkinMaleModule struct {
	Enabled bool      `json:"enabled"`
	Lookups []*Lookup `json:"lookups"`
	// 面部美肤指数
	F109 int `json:"f109"`
}

// GynecologicalDiseaseRiskEstimateModule 妇科疾病风险预估模块.
type GynecologicalDiseaseRiskEstimateModule struct {
	Enabled bool      `json:"enabled"`
	Lookups []*Lookup `json:"lookups"`
	// 妇科疾病风险指数
	F101 int `json:"f101"`
}

// GynecologicalInflammationModule 妇科炎症模块.
type GynecologicalInflammationModule struct {
	Enabled bool      `json:"enabled"`
	Lookups []*Lookup `json:"lookups"`
	// 妇科炎症指数
	F102 int `json:"f102"`
}

// HormoneLevelModule 激素水平模块.
type HormoneLevelModule struct {
	Enabled bool      `json:"enabled"`
	Lookups []*Lookup `json:"lookups"`
	// 激素水平指数
	F106 int `json:"f106"`
}

// UterineHealthModule 子宫健康模块.
type UterineHealthModule struct {
	Enabled bool      `json:"enabled"`
	Lookups []*Lookup `json:"lookups"`
	// 子宫健康指数
	F100 int `json:"f100"`
}

// MenstrualSunflowerModule 月经（天葵）健康模块.
type MenstrualSunflowerModule struct {
	Enabled bool      `json:"enabled"`
	Lookups []*Lookup `json:"lookups"`
	M0      int       `json:"m0"`
}

// IrregularMenstruationModule 月经不调模块.
type IrregularMenstruationModule struct {
	Enabled bool      `json:"enabled"`
	Lookups []*Lookup `json:"lookups"`
	M1      int       `json:"m1"`
}

// DysmenorrheaIndexModule 月经不调模块.
type DysmenorrheaIndexModule struct {
	Enabled bool      `json:"enabled"`
	Lookups []*Lookup `json:"lookups"`
	M2      int       `json:"m2"`
}

// ReproductiveAgeModule 生殖年龄模块.
type ReproductiveAgeModule struct {
	Enabled bool      `json:"enabled"`
	Lookups []*Lookup `json:"lookups"`
	F05     int       `json:"f05"`
}

// LymphaticHealthModule 淋巴健康模块.
type LymphaticHealthModule struct {
	Enabled bool      `json:"enabled"`
	Lookups []*Lookup `json:"lookups"`
	F07     int       `json:"f07"`
}

// BreastHealthModule 乳腺健康模块.
type BreastHealthModule struct {
	Enabled bool      `json:"enabled"`
	Lookups []*Lookup `json:"lookups"`
	F108    int       `json:"f108"`
	M3      int       `json:"m3"`
}

// BreastCancerModule 乳腺癌卵巢癌风险模块.
type BreastCancerModule struct {
	Enabled bool      `json:"enabled"`
	Lookups []*Lookup `json:"lookups"`
}

// AnalysisReportContent 分析报告内容
type AnalysisReportContent struct {
	// 测量时间
	CreatedTime time.Time `json:"created_time"`
	// 用户的个人信息
	UserProfile UserProfileModule `json:"user_profile"`
	// 测量上下文
	PulseTest PulseTestModule `json:"pulse_test"`
	// 备注模块
	Remark RemarkModule `json:"remark"`
	// 下面是各个测量结果模块
	// heartrate 心率模块
	HeartRate HeartRateModule `json:"heartrate"`
	// risk_estimate 风险预估模块
	RiskEstimate RiskEstimateModule `json:"risk_estimate"`
	// physical_dialectics 体质辩证模块
	PhysicalDialectics PhysicalDialecticsModule `json:"physical_dialectics"`
	// dirty_dialectic 脏腑辩证模块
	DirtyDialectic DirtyDialecticModule `json:"dirty_dialectic"`
	// physical_therapy_index 理疗指数模块
	PhysicalTherapyIndex PhysicalTherapyIndexModule `json:"physical_therapy_index"`
	// conditioning_advice 调理建议模块
	ConditioningAdvice ConditioningAdviceModule `json:"conditioning_advice"`
	// partial_pulse_wave 局部脉搏波模块
	PartialPulseWave PartialPulseWaveModule `json:"partial_pulse_wave"`
	// meridian_bar_chart 经络柱状图模块
	MeridianBarChart MeridianBarChartModule `json:"meridian_bar_chart"`
	// meridian_explain 经络解读模块
	MeridianExplain MeridianExplainModule `json:"meridian_explain"`
	// measurement_judgment 测量异常判断模块
	MeasurementJudgment MeasurementJudgmentModule `json:"measurement_judgment"`

	// tips 温馨提示模块
	Tips TipsModule `json:"tips"`
	// stress_state_judgment 应激态模块
	StressStateJudgment StressStateJudgmentModule `json:"stress_state_judgment"`

	// 一体机模块

	// emotional_health 情绪健康模块
	EmotionalHealth EmotionalHealthModule `json:"emotional_health"`
	// facial_skin 面部美肤模块
	FacialSkin FacialSkinModule `json:"facial_skin"`
	// facial_skin_male 男性面部美肤模块
	FacialSkinMale FacialSkinMaleModule `json:"facial_skin_male"`
	// gynecological_disease_risk_estimate 妇科疾病风险预估模块
	GynecologicalDiseaseRiskEstimate GynecologicalDiseaseRiskEstimateModule `json:"gynecological_disease_risk_estimate"`
	// gynecological_inflammation 妇科炎症模块
	GynecologicalInflammation GynecologicalInflammationModule `json:"gynecological_inflammations"`
	// hormone_level 激素水平模块
	HormoneLevel HormoneLevelModule `json:"hormone_level"`
	// uterine_health 子宫健康模块
	UterineHealth UterineHealthModule `json:"uterine_health"`
	// menstrual_sunflower 月经（天葵）健康模块
	MenstrualSunflower MenstrualSunflowerModule `json:"menstrual_sunflower"`
	// irregular_menstruation 月经不调模块
	IrregularMenstruation IrregularMenstruationModule `json:"irregular_menstruation"`
	// dysmenorrhea_index 痛经模块
	DysmenorrheaIndexModule DysmenorrheaIndexModule `json:"dysmenorrhea_index"`
	// reproductive_age 生殖年龄模块
	ReproductiveAge ReproductiveAgeModule `json:"reproductive_age"`
	// lymphatic_health 淋巴健康模块
	LymphaticHealth LymphaticHealthModule `json:"lymphatic_health"`
	// breast_health 乳腺健康模块
	BreastHealthModule BreastHealthModule `json:"breast_health"`
	// breast_cancer 乳腺癌卵巢癌风险模块
	BreastCancerModule BreastCancerModule `json:"breast_cancer"`
}

// DietaryAdvice 饮食建议
type DietaryAdvice struct {
	AnalysisReportLookup
	AnalysisReportHealthTip
}

// SportsAdvice 运动建议
type SportsAdvice struct {
	AnalysisReportLookup
	AnalysisReportHealthTip
}

// MonthlyReportContent 月报的内容
type MonthlyReportContent struct {
	UserProfile                 ReportUserProfile         `json:"user_profile"`
	DisplayOptions              DisplayOptions            `json:"display_options"`
	MeasurementJudgment         AnalysisReportLookup      `json:"measurement_judgment"`
	SwitchPretreatment          AnalysisReportLookup      `json:"switch_pretreatment"`
	PhysiologicalPeriodJudgment AnalysisReportLookup      `json:"physiological_periodJudgment"`
	LactationJudgment           AnalysisReportLookup      `json:"lactation_judgment"`
	OvulationJudgment           AnalysisReportLookup      `json:"ovulation_judgment"`
	PregnancyJudgment           AnalysisReportLookup      `json:"pregnancy_judgment"`
	StressStateJudgment         StressStateJudgment       `json:"stress_state_judgment"`
	MeridianInterpretation      AnalysisReportLookup      `json:"meridian_interpretation"`
	PhysicalDialectics          AnalysisReportLookup      `json:"physical_dialectics"`
	DirtyDialectic              AnalysisReportLookup      `json:"dirty_dialectic"`
	FactorInterpretation        FactorInterpretation      `json:"factor_interpretation"`
	UterineHealthIndex          UterineHealthIndex        `json:"uterine_health_index"`
	MenstrualHealthIndex        MenstrualHealthIndex      `json:"menstrual_health_index"`
	GynecologicalInflammation   GynecologicalInflammation `json:"gynecological_inflammation"`
	BreastHealth                BreastHealth              `json:"breast_health"`
	EmotionalHealth             EmotionalHealth           `json:"emotional_health"`
	FacialSkin                  FacialSkin                `json:"facial_skin"`
	ReproductiveAge             ReproductiveAge           `json:"reproductive_age"`
	HormoneLevel                HormoneLevel              `json:"hormone_level"`
	BreastCancer                AnalysisReportLookup      `json:"breast_cancer"`
	LymphaticHealth             LymphaticHealth           `json:"lymphatic_health"`
	FacialSkinMale              FacialSkinMale            `json:"facial_skin_male"`
	HormoneLevelMale            HormoneLevelMale          `json:"hormone_level_male"`
	BloodSugar                  AnalysisReportLookup      `json:"blood_sugar"`
	BloodPressure               AnalysisReportLookup      `json:"blood_pressure"`
	Hyperlipidemia              AnalysisReportLookup      `json:"hyperlipidemia"`
	AcutePharyngitis            AnalysisReportLookup      `json:"acute_pharyngitis"`
	Spine                       AnalysisReportLookup      `json:"spine"`
	Anxiety                     AnalysisReportLookup      `json:"anxiety"`
	SpinalDisease               AnalysisReportLookup      `json:"spinal_disease"`
	ChronicCough                AnalysisReportLookup      `json:"chronic_cough"`
	ProstateDisease             AnalysisReportLookup      `json:"prostate_disease"`
	InflammationRisk            AnalysisReportLookup      `json:"inflammation_risk"`
	CerebralInsufficiency       AnalysisReportLookup      `json:"cerebral_insufficiency"`
	Immunity                    AnalysisReportLookup      `json:"immunity"`
	FatigueAndPressure          AnalysisReportLookup      `json:"fatigue_and_pressure"`
	RenalDysfunction            AnalysisReportLookup      `json:"renal_dysfunction"`
	CoronaryHeartDisease        AnalysisReportLookup      `json:"chd"`
	Depression                  AnalysisReportLookup      `json:"depression"`
	SleepProblems               AnalysisReportLookup      `json:"sleep_problems"`
	HypomotilityOfStomach       AnalysisReportLookup      `json:"hypomotility_of_stomach"`
	Gastritis                   AnalysisReportLookup      `json:"gastritis"`
	HasAnsweredAllQuestions     bool                      `json:"has_answered_all_questions"`
	ChineseMedicineAdvice       AnalysisReportLookup      `json:"chinese_medicine_advice"`
	DietaryAdvice               DietaryAdvice             `json:"dietary_advice"`
	PhysicalTherapyAdvice       AnalysisReportLookup      `json:"physical_therapy_advice"`
	SportsAdvice                SportsAdvice              `json:"sports_advice"`
	MonthlyReport               MonthlyReport             `json:"monthly_report"`
}

// MonthlyReport 月报
type MonthlyReport struct {
	AnalysisReportLookup
	PhysicalDialecticsStatSummary []Lookup `json:"physical_dialectics_stat_summary"`
	PhysicalDialecticsDescription []Lookup `json:"physical_dialectics_description"`
}

// WeeklyReportContent 周报
type WeeklyReportContent struct {
	UserProfile                 ReportUserProfile         `json:"user_profile"`
	DisplayOptions              DisplayOptions            `json:"display_options"`
	MeasurementJudgment         AnalysisReportLookup      `json:"measurement_judgment"`
	SwitchPretreatment          AnalysisReportLookup      `json:"switch_pretreatment"`
	PhysiologicalPeriodJudgment AnalysisReportLookup      `json:"physiological_periodJudgment"`
	LactationJudgment           AnalysisReportLookup      `json:"lactation_judgment"`
	OvulationJudgment           AnalysisReportLookup      `json:"ovulation_judgment"`
	PregnancyJudgment           AnalysisReportLookup      `json:"pregnancy_judgment"`
	StressStateJudgment         StressStateJudgment       `json:"stress_state_judgment"`
	MeridianInterpretation      AnalysisReportLookup      `json:"meridian_interpretation"`
	PhysicalDialectics          AnalysisReportLookup      `json:"physical_dialectics"`
	DirtyDialectic              AnalysisReportLookup      `json:"dirty_dialectic"`
	FactorInterpretation        FactorInterpretation      `json:"factor_interpretation"`
	UterineHealthIndex          UterineHealthIndex        `json:"uterine_health_index"`
	MenstrualHealthIndex        MenstrualHealthIndex      `json:"menstrual_health_index"`
	GynecologicalInflammation   GynecologicalInflammation `json:"gynecological_inflammation"`
	BreastHealth                BreastHealth              `json:"breast_health"`
	EmotionalHealth             EmotionalHealth           `json:"emotional_health"`
	FacialSkin                  FacialSkin                `json:"facial_skin"`
	ReproductiveAge             ReproductiveAge           `json:"reproductive_age"`
	HormoneLevel                HormoneLevel              `json:"hormone_level"`
	BreastCancer                AnalysisReportLookup      `json:"breast_cancer"`
	LymphaticHealth             LymphaticHealth           `json:"lymphatic_health"`
	FacialSkinMale              FacialSkinMale            `json:"facial_skin_male"`
	HormoneLevelMale            HormoneLevelMale          `json:"hormone_level_male"`
	BloodSugar                  AnalysisReportLookup      `json:"blood_sugar"`
	BloodPressure               AnalysisReportLookup      `json:"blood_pressure"`
	Hyperlipidemia              AnalysisReportLookup      `json:"hyperlipidemia"`
	AcutePharyngitis            AnalysisReportLookup      `json:"acute_pharyngitis"`
	Spine                       AnalysisReportLookup      `json:"spine"`
	Anxiety                     AnalysisReportLookup      `json:"anxiety"`
	SpinalDisease               AnalysisReportLookup      `json:"spinal_disease"`
	ChronicCough                AnalysisReportLookup      `json:"chronic_cough"`
	ProstateDisease             AnalysisReportLookup      `json:"prostate_disease"`
	InflammationRisk            AnalysisReportLookup      `json:"inflammation_risk"`
	CerebralInsufficiency       AnalysisReportLookup      `json:"cerebral_insufficiency"`
	Immunity                    AnalysisReportLookup      `json:"immunity"`
	FatigueAndPressure          AnalysisReportLookup      `json:"fatigue_and_pressure"`
	RenalDysfunction            AnalysisReportLookup      `json:"renal_dysfunction"`
	CoronaryHeartDisease        AnalysisReportLookup      `json:"chd"`
	Depression                  AnalysisReportLookup      `json:"depression"`
	SleepProblems               AnalysisReportLookup      `json:"sleep_problems"`
	HypomotilityOfStomach       AnalysisReportLookup      `json:"hypomotility_of_stomach"`
	Gastritis                   AnalysisReportLookup      `json:"gastritis"`
	HasAnsweredAllQuestions     bool                      `json:"has_answered_all_questions"`
	ChineseMedicineAdvice       AnalysisReportLookup      `json:"chinese_medicine_advice"`
	DietaryAdvice               DietaryAdvice             `json:"dietary_advice"`
	PhysicalTherapyAdvice       AnalysisReportLookup      `json:"physical_therapy_advice"`
	SportsAdvice                SportsAdvice              `json:"sports_advice"`
	WeeklyReport                WeeklyReport              `json:"weekly_report"`
}

// WeeklyReport 周报
type WeeklyReport struct {
	AnalysisReportLookup
	PhysicalDialecticsStatSummary []Lookup `json:"physical_dialectics_stat_summary"`
	PhysicalDialecticsDescription []Lookup `json:"physical_dialectics_description"`
}

// FactorInterpretation 因子解读
type FactorInterpretation struct {
	AnalysisReportLookup
	F0 *int32 `json:"f0"`
	F1 *int32 `json:"f1"`
	F2 *int32 `json:"f2"`
	F3 *int32 `json:"f3"`
}

// UterineHealthIndex 子宫健康指数
type UterineHealthIndex struct {
	AnalysisReportLookup
	F100 *int32 `json:"f100"`
}

// MenstrualHealthIndex 月经健康指数
type MenstrualHealthIndex struct {
	AnalysisReportLookup
	F101 *int32 `json:"f101"`
	M0   *int32 `json:"m0"`
	M1   *int32 `json:"m1"`
	M2   *int32 `json:"m2"`
}

// GynecologicalInflammation 妇科炎症
type GynecologicalInflammation struct {
	AnalysisReportLookup
	F102 *int32 `json:"f102"`
}

// BreastHealth 乳房健康
type BreastHealth struct {
	AnalysisReportLookup
	F108 *int32 `json:"f108"`
	M3   *int32 `json:"m3"`
}

// EmotionalHealth 情绪健康
type EmotionalHealth struct {
	AnalysisReportLookup
	F103 *int32 `json:"f103"`
}

// FacialSkin 面部美肤指数
type FacialSkin struct {
	AnalysisReportLookup
	F104 *int32 `json:"f104"`
}

// ReproductiveAge 生殖年龄
type ReproductiveAge struct {
	AnalysisReportLookup
	F105 *int32 `json:"f105"`
}

// HormoneLevel 激素水平
type HormoneLevel struct {
	AnalysisReportLookup
	F106 *int32 `json:"f106"`
}

// LymphaticHealth 淋巴健康指数
type LymphaticHealth struct {
	AnalysisReportLookup
	F107 *int32 `json:"f107"`
}

// FacialSkinMale 男性面部美肤指数
type FacialSkinMale struct {
	AnalysisReportLookup
	F109 *int32 `json:"f109"`
}

// HormoneLevelMale 男性激素水平
type HormoneLevelMale struct {
	AnalysisReportLookup
	F110 *int32 `json:"f110"`
}

// StressStateJudgment 应激态判断标识
type StressStateJudgment struct {
	AnalysisReportLookup
	HasStressState         bool `json:"has_stress_state"`
	HasDoneSports          bool `json:"has_done_sports"`
	HasDrinkedWine         bool `json:"has_drinked_wine"`
	HasHadCold             bool `json:"has_had_cold"`
	HasRhinitisEpisode     bool `json:"has_rhinitis_episode"`
	HasAbdominalPain       bool `json:"has_abdominal_pain"`
	HasViralInfection      bool `json:"has_viral_infection"`
	HasPhysiologicalPeriod bool `json:"has_physiological_period"`
	HasOvulation           bool `json:"has_ovulation"`
	HasPregnant            bool `json:"has_pregnant"`
	HasLactation           bool `json:"has_lactation"`
}

// AnalysisReportLookup 分析报告Lookups
type AnalysisReportLookup struct {
	Lookups []Lookup `json:"lookups"`
}

// AnalysisReportHealthTip 分析报告健康提示
type AnalysisReportHealthTip struct {
	HealthTips []HealthTip `json:"health_tips"`
}

// HealthTip 分析报告健康提示
type HealthTip map[string]string
