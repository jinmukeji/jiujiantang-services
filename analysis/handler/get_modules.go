package handler

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	any "github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/jinmukeji/ae/v2/engine/core"
	"github.com/jinmukeji/jiujiantang-services/analysis/mysqldb"
	"github.com/jinmukeji/plat-pkg/v2/micro/errors"
	"github.com/jinmukeji/plat-pkg/v2/micro/errors/codes"
	analysispb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/analysis/v1"
	ptypesv2 "github.com/jinmukeji/proto/v3/gen/micro/idl/ptypes/v2"
)

// getModulesFromAEOutput 解析引擎输出的Output，获得与引擎相关的输出模块，用于单次测量
func getModulesFromAEOutput(ctxData core.ContextData, options map[string]bool, gender ptypesv2.Gender) (map[string]*any.Any, error) {
	out := ctxData.Output

	// 构建各个模块
	modules := make(map[string]*any.Any)
	// 构建体质辩证模块
	physicalDialecticsModule := analysispb.PhysicalDialecticsModule{
		Enabled: options["DisplayCdExplain"],
		Lookups: parseLookupsIncludeKeyContent(out["physical_dialectics"]),
	}
	err := setModule(modules, "physical_dialectics", &physicalDialecticsModule)
	if err != nil {
		return nil, fmt.Errorf("failed to set value of module physical_dialectics :[%s]", err.Error())
	}
	// 构建其他模块
	sameModules, err := buildSameModulesOfWeeklyMonthlyAndMeasurement(ctxData, options, gender)
	if err != nil {
		return nil, fmt.Errorf("failed to build modules :[%s]", err.Error())
	}
	for key, value := range sameModules {
		modules[key] = value
	}
	return modules, nil

}

func (j *AnalysisManagerService) buildModulesNotAboutAE(record *mysqldb.Record, displayOptions map[string]bool, measurementJudgment bool) (map[string]*any.Any, error) {

	respModules := make(map[string]*any.Any)
	// 构建心率模块
	heartRateModule := analysispb.HeartRateModule{
		Enabled: displayOptions["DisplayHeartRate"],

		AverageHeartRate: int32(record.HeartRate),
		HighestHeartRate: record.AlgorithmHighestHeartRate,
		LowestHeartRate:  record.AlgorithmLowestHeartRate,
	}
	any, err := ptypes.MarshalAny(&heartRateModule)
	if err != nil {
		return nil, errors.ErrorWithCause(codes.InvalidArgument, err, "failed to marshal heart rate module")
	}
	respModules["heart_rate"] = any
	// 构建经络柱状图模块
	createdAt, _ := ptypes.TimestampProto(record.CreatedAt)
	meridianBarChartModule := analysispb.MeridianBarChartModule{
		Enabled: displayOptions["DisplayCcBarChart"],
		MeridianValue: &analysispb.CInfo{
			C0:       int32(record.C0),
			C1:       int32(record.C1),
			C2:       int32(record.C2),
			C3:       int32(record.C3),
			C4:       int32(record.C4),
			C5:       int32(record.C5),
			C6:       int32(record.C6),
			C7:       int32(record.C7),
			TestTime: createdAt,
		},
	}
	any, err = ptypes.MarshalAny(&meridianBarChartModule)
	if err != nil {
		return nil, errors.ErrorWithCause(codes.InvalidArgument, err, "failed to marshal meridian bar chart module")
	}
	respModules["meridian_bar_chart"] = any
	// 如果为单次测量，则需要构建局部脉搏波模块，否则返回默认值
	if measurementJudgment {
		// 构建局部脉搏波模块
		var partialData []int32
		dataArray, errgetPulseTestDataIntArray := getPulseTestDataIntArray(record.S3Key, j.awsClient)
		if errgetPulseTestDataIntArray != nil {
			log.WithError(errgetPulseTestDataIntArray).Warn("failed to get pulse test data int array")
		}
		if len(dataArray) >= 4000 {
			partialData = dataArray[3000:4000]
		}

		partialPulseWaveModule := analysispb.PartialPulseWaveModule{
			Enabled: displayOptions["DisplayPartialData"],
			Points:  partialData,
		}
		any, err = ptypes.MarshalAny(&partialPulseWaveModule)
		if err != nil {
			return nil, errors.ErrorWithCause(codes.InvalidArgument, err, "failed to marshal partial pulse wave module")
		}
		respModules["partial_pulse_wave"] = any
	}

	return respModules, nil
}

// getModulesOfWeeklyReport 解析周报分析的Output，获得与引擎相关的输出模块
func getModulesOfWeeklyReport(ctxData core.ContextData, options map[string]bool, gender ptypesv2.Gender) (map[string]*any.Any, error) {
	out := ctxData.Output

	// 构建各个模块
	modules := make(map[string]*any.Any)
	weeklyReport := out["weekly_report"]
	mapWeeklyReport, ok := weeklyReport.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to marshal weekly report module")
	}
	dialecticsStatSummary := mapWeeklyReport["physical_dialectics_stat_summary"]
	dialecticsDescription := mapWeeklyReport["physical_dialectics_description"]
	// 构建体质辩证模块
	statSummary := parseIncludeKeyContent(dialecticsStatSummary)
	description := parseIncludeKeyContent(dialecticsDescription)
	physicalDialecticsLookups := make([]*analysispb.Lookup, len(statSummary)+len(description))
	copy(physicalDialecticsLookups, statSummary)
	for idx2, value2 := range description {
		physicalDialecticsLookups[idx2+len(statSummary)] = value2
	}
	physicalDialecticsModule := analysispb.PhysicalDialecticsModule{
		Enabled: options["DisplayCdExplain"],
		Lookups: physicalDialecticsLookups,
	}
	err := setModule(modules, "physical_dialectics", &physicalDialecticsModule)
	if err != nil {
		return nil, fmt.Errorf("failed to set value of module physical_dialectics :[%s]", err.Error())
	}

	// 构建其他模块
	sameModules, err := buildSameModulesOfWeeklyMonthlyAndMeasurement(ctxData, options, gender)
	if err != nil {
		return nil, fmt.Errorf("failed to build modules :[%s]", err.Error())
	}
	for key, value := range sameModules {
		modules[key] = value
	}
	return modules, nil
}

// getModulesOfMonthlyReport 解析月报分析的Output，获得与引擎相关的输出模块
func getModulesOfMonthlyReport(ctxData core.ContextData, options map[string]bool, gender ptypesv2.Gender) (map[string]*any.Any, error) {
	out := ctxData.Output

	// 构建各个模块
	modules := make(map[string]*any.Any)
	monthlyReport := out["monthly_report"]
	mapMonthlyReport, ok := monthlyReport.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to marshal weekly report module")
	}
	dialecticsStatSummary := mapMonthlyReport["physical_dialectics_stat_summary"]
	dialecticsDescription := mapMonthlyReport["physical_dialectics_description"]
	// 构建体质辩证模块
	statSummary := parseIncludeKeyContent(dialecticsStatSummary)
	description := parseIncludeKeyContent(dialecticsDescription)
	physicalDialecticsLookups := make([]*analysispb.Lookup, len(statSummary)+len(description))
	copy(physicalDialecticsLookups, statSummary)
	for idx2, value2 := range description {
		physicalDialecticsLookups[idx2+len(statSummary)] = value2
	}
	physicalDialecticsModule := analysispb.PhysicalDialecticsModule{
		Enabled: options["DisplayCdExplain"],
		Lookups: physicalDialecticsLookups,
	}
	err := setModule(modules, "physical_dialectics", &physicalDialecticsModule)
	if err != nil {
		return nil, fmt.Errorf("failed to set value of module physical_dialectics :[%s]", err.Error())
	}

	// 构建其他模块
	sameModules, err := buildSameModulesOfWeeklyMonthlyAndMeasurement(ctxData, options, gender)
	if err != nil {
		return nil, fmt.Errorf("failed to build modules :[%s]", err.Error())
	}
	for key, value := range sameModules {
		modules[key] = value
	}
	return modules, nil
}

// buildSameModulesOfWeeklyMonthlyAndMeasurement 解析单次测量和周报月报分析的Output，获得与共同输出模块
func buildSameModulesOfWeeklyMonthlyAndMeasurement(ctxData core.ContextData, options map[string]bool, gender ptypesv2.Gender) (map[string]*any.Any, error) {
	out := ctxData.Output

	// 构建各个模块
	modules := make(map[string]*any.Any)

	// 构建测量异常判断模块
	measurementJudgmentModule := analysispb.MeasurementJudgment{
		Enabled: options["DisplayMeasurementJudgment"],
		Lookups: parseLookupsIncludeKeyContent(out["measurement_judgment"]),
	}
	err := setModule(modules, "measurement_judgment", &measurementJudgmentModule)
	if err != nil {
		return nil, fmt.Errorf("failed to set value of module measurement_judgment :[%s]", err.Error())
	}
	// 构建温馨提示模块
	tips := analysispb.Tips{
		Enabled: options["DisplayTips"],
		Lookups: parseLookupsIncludeKeyContent(out["stress_state_judgment"]),
	}
	err = setModule(modules, "tips", &tips)
	if err != nil {
		return nil, fmt.Errorf("failed to set value of module tips :[%s]", err.Error())
	}
	// 构建身体状况模块
	stressStateJudgment := getPhysicalConditions(ctxData, options)
	err = setModule(modules, "stress_state_judgment", &stressStateJudgment)
	if err != nil {
		return nil, fmt.Errorf("failed to set value of module stress_state_judgment :[%s]", err.Error())
	}
	// 构建脏腑辩证模块
	dirtyDialecticModule := analysispb.DirtyDialecticModule{
		Enabled: options["DisplaySdExplain"],
		Lookups: parseLookupsIncludeKeyContent(out["dirty_dialectic"]),
	}
	err = setModule(modules, "dirty_dialectic", &dirtyDialecticModule)
	if err != nil {
		return nil, fmt.Errorf("failed to set value of module dirty_dialectic :[%s]", err.Error())
	}
	// 构建理疗指数模块
	physicalTherapyIndexModule := analysispb.PhysicalTherapyIndexModule{
		Enabled: options["DisplayPhysicalTherapyIndex"],
		Lookups: parseLookupsIncludeKeyContent(out["factor_interpretation"]),
	}
	if out["factor_interpretation"] != nil {
		if _, ok := out["factor_interpretation"].(map[string]interface{}); ok {
			physicalTherapyIndexModule.F0 = &wrappers.Int32Value{Value: int32(out["factor_interpretation"].(map[string]interface{})["f0"].(float64))}
			physicalTherapyIndexModule.F1 = &wrappers.Int32Value{Value: int32(out["factor_interpretation"].(map[string]interface{})["f1"].(float64))}
			physicalTherapyIndexModule.F2 = &wrappers.Int32Value{Value: int32(out["factor_interpretation"].(map[string]interface{})["f2"].(float64))}
			physicalTherapyIndexModule.F3 = &wrappers.Int32Value{Value: int32(out["factor_interpretation"].(map[string]interface{})["f3"].(float64))}
		}
	}
	err = setModule(modules, "physical_therapy_index", &physicalTherapyIndexModule)
	if err != nil {
		return nil, fmt.Errorf("failed to set value of module physical_therapy_index :[%s]", err.Error())
	}

	// 构建风险预估模块
	disease, promptMessage := getDiseasesMessageFromInput(ctxData.Output)
	riskEstimateModule := analysispb.RiskEstimateModule{
		Enabled:         options["DisplayRiskEstimate"],
		DiseaseEstimate: disease,
		PromptMessage:   promptMessage,
	}
	err = setModule(modules, "risk_estimate", &riskEstimateModule)
	if err != nil {
		return nil, fmt.Errorf("failed to set value of module risk_estimate :[%s]", err.Error())
	}
	// 构建调理建议模块
	conditioningAdviceModule := analysispb.ConditioningAdviceModule{
		Enabled: options["DisplayConditioningAdvice"],
		DietaryAdvice: &analysispb.DietaryAdviceModule{
			Enabled: options["DisplayDietaryAdvice"],
			Lookups: parseLookupsIncludeKeyContent(out["dietary_advice"]),
		},
		SportsAdvice: &analysispb.SportsAdviceModule{
			Enabled: options["DisplaySportsAdvice"],
			Lookups: parseLookupsIncludeKeyContent(out["sports_advice"]),
		},
		ChineseMedicineAdvice: &analysispb.ChineseMedicineAdviceModule{
			Enabled: options["DisplayChineseMedicineAdvice"],
			Lookups: parseLookupsIncludeKeyContent(out["chinese_medicine_advice"]),
		},
		PhysicalTherapyAdvice: &analysispb.PhysicalTherapyAdviceModule{
			Enabled: options["DisplayPhysicalTherapyAdvice"],
			Lookups: parseLookupsIncludeKeyContent(out["physical_therapy_advice"]),
		},
	}
	err = setModule(modules, "conditioning_advice", &conditioningAdviceModule)
	if err != nil {
		return nil, fmt.Errorf("failed to set value of module conditioning_advice :[%s]", err.Error())
	}
	// 构建经络解读模块 proto 合并后取消下面一段的注释
	meridianExplainModule := analysispb.MeridianExplainModule{
		Enabled: options["DisplayCcExplain"],
	}
	meridianExplainModule.Lookups = parseLookupsIncludeKeyContent(out["meridian_explain"])
	err = setModule(modules, "meridian_explain", &meridianExplainModule)
	if err != nil {
		return nil, fmt.Errorf("failed to set value of module meridian_explain :[%s]", err.Error())
	}
	// 构建情绪健康模块
	emotionalHealthModule := analysispb.EmotionalHealthModule{
		Enabled: options["DisplayF3"],
	}
	emotionalHealthModule.Lookups = parseLookupsIncludeKeyContent(out["emotional_health"])
	if out["emotional_health"] != nil {
		if _, ok := out["factor_interpretation"].(map[string]interface{}); ok {
			emotionalHealthModule.F103 = &wrappers.Int32Value{Value: int32(out["emotional_health"].(map[string]interface{})["f103"].(float64))}
		}
	}
	err = setModule(modules, "emotional_health", &emotionalHealthModule)
	if err != nil {
		return nil, fmt.Errorf("failed to set value of module emotional_health :[%s]", err.Error())
	}
	// 构建面部美肤模块
	// 如果为男性则输出男性面部美肤，如果为女性输出女性面部美肤
	if gender == ptypesv2.Gender_GENDER_FEMALE {
		facialSkinModule := analysispb.FacialSkinModule{
			Enabled: options["DisplayFacialSkin"],
			Lookups: parseLookupsIncludeKeyContent(out["facial_skin"]),
		}
		if out["facial_skin"] != nil {
			if _, ok := out["facial_skin"].(map[string]interface{}); ok {
				facialSkinModule.F104 = &wrappers.Int32Value{Value: int32(out["facial_skin"].(map[string]interface{})["f104"].(float64))}
			}
		}
		err = setModule(modules, "facial_skin", &facialSkinModule)
		if err != nil {
			return nil, fmt.Errorf("failed to set value of module facial_skin :[%s]", err.Error())
		}
	}
	if gender == ptypesv2.Gender_GENDER_MALE {
		facialSkinMaleModule := analysispb.FacialSkinMaleModule{
			Enabled: options["DisplayFacialSkin"],
			Lookups: parseLookupsIncludeKeyContent(out["facial_skin_male"]),
		}
		if out["facial_skin_male"] != nil {
			facialSkinMaleModule.F109 = &wrappers.Int32Value{Value: int32(out["facial_skin_male"].(map[string]interface{})["f109"].(float64))}
		}
		err = setModule(modules, "facial_skin_male", &facialSkinMaleModule)
		if err != nil {
			return nil, fmt.Errorf("failed to set value of module facial_skin_male :[%s]", err.Error())
		}
	}
	// 构建妇科疾病风险预估模块
	gynecologicalDiseaseRiskEstimateModule := analysispb.GynecologicalDiseaseRiskEstimateModule{
		Enabled: options["DisplayGynecologicalDiseaseRisk"],
		Lookups: parseLookupsIncludeKeyContent(out["gynecological_disease_risk"]),
	}
	if out["gynecological_disease_risk"] != nil {
		if _, ok := out["gynecological_disease_risk"].(map[string]interface{}); ok {
			gynecologicalDiseaseRiskEstimateModule.F101 = &wrappers.Int32Value{Value: int32(out["gynecological_disease_risk"].(map[string]interface{})["f101"].(float64))}
		}

	}
	err = setModule(modules, "gynecological_disease_risk_estimate", &gynecologicalDiseaseRiskEstimateModule)
	if err != nil {
		return nil, fmt.Errorf("failed to set value of module gynecological_disease_risk_estimate :[%s]", err.Error())
	}
	// 构建妇科炎症指数模块
	gynecologicalInflammationModule := analysispb.GynecologicalInflammationModule{
		Enabled: options["DisplayGynecologicalInflammation"],
	}
	gynecologicalInflammation := getGynecologicalInflammation(out["inflammation_risk"])
	if gynecologicalInflammation != nil {
		gynecologicalInflammationModule.F102 = &wrappers.Int32Value{Value: int32(gynecologicalInflammation.GetScore())}
	}

	err = setModule(modules, "gynecological_inflammations", &gynecologicalInflammationModule)
	if err != nil {
		return nil, fmt.Errorf("failed to set value of module gynecological_inflammations :[%s]", err.Error())
	}
	// 构建激素水平模块
	hormoneLevelModule := analysispb.HormoneLevelModule{
		Enabled: options["DisplayHormoneLevel"],
		Lookups: parseLookupsIncludeKeyContent(out["hormone_level"]),
	}
	if out["hormone_level"] != nil {
		if _, ok := out["hormone_level"].(map[string]interface{}); ok {
			hormoneLevelModule.F106 = &wrappers.Int32Value{Value: int32(out["hormone_level"].(map[string]interface{})["f106"].(float64))}
		}
	}
	err = setModule(modules, "hormone_level", &hormoneLevelModule)
	if err != nil {
		return nil, fmt.Errorf("failed to set value of module hormone_level :[%s]", err.Error())
	}
	// 构建子宫健康模块
	uterineHealthModule := analysispb.UterineHealthModule{
		Enabled: options["DisplayUterineHealth"],
		Lookups: parseLookupsIncludeKeyContent(out["uterine_health_index"]),
	}
	if out["uterine_health_index"] != nil {
		if _, ok := out["uterine_health_index"].(map[string]interface{}); ok {
			uterineHealthModule.F100 = &wrappers.Int32Value{Value: int32(out["uterine_health_index"].(map[string]interface{})["f100"].(float64))}
		}
	}
	err = setModule(modules, "uterine_health", &uterineHealthModule)
	if err != nil {
		return nil, fmt.Errorf("failed to set value of module uterine_health :[%s]", err.Error())
	}
	// 构建月经（天葵）健康模块
	menstrualHealthModule := analysispb.MenstrualSunflowerModule{
		Enabled: options["DisplayMenstrualHealth"],
		Lookups: parseLookupsIncludeKeyContent(out["menstrual_sunflower"]),
	}
	if out["menstrual_sunflower"] != nil {
		if _, ok := out["menstrual_sunflower"].(map[string]interface{}); ok {
			menstrualHealthModule.M0 = &wrappers.Int32Value{Value: int32(out["menstrual_sunflower"].(map[string]interface{})["m0"].(float64))}
		}
	}
	err = setModule(modules, "menstrual_sunflower", &menstrualHealthModule)
	if err != nil {
		return nil, fmt.Errorf("failed to set value of module menstrual_sunflower :[%s]", err.Error())
	}
	// 构建月经不调模块
	irregularMenstruationModule := analysispb.IrregularMenstruationModule{
		Enabled: options["DisplayMenstruationIrregular"],
		Lookups: parseLookupsIncludeKeyContent(out["irregular_menstruation"]),
	}
	if out["irregular_menstruation"] != nil {
		if _, ok := out["irregular_menstruation"].(map[string]interface{}); ok {
			irregularMenstruationModule.M1 = &wrappers.Int32Value{Value: int32(out["irregular_menstruation"].(map[string]interface{})["m1"].(float64))}
		}
	}
	err = setModule(modules, "irregular_menstruation", &irregularMenstruationModule)
	if err != nil {
		return nil, fmt.Errorf("failed to set value of module irregular_menstruation :[%s]", err.Error())
	}
	// 构建痛经模块
	dysmenorrheaIndexModule := analysispb.DysmenorrheaIndexModule{
		Enabled: options["DisplayDysmenorrhea"],
		Lookups: parseLookupsIncludeKeyContent(out["dysmenorrhea_index"]),
	}
	if out["dysmenorrhea_index"] != nil {
		if _, ok := out["dysmenorrhea_index"].(map[string]interface{}); ok {
			dysmenorrheaIndexModule.M2 = &wrappers.Int32Value{Value: int32(out["dysmenorrhea_index"].(map[string]interface{})["m2"].(float64))}
		}
	}
	err = setModule(modules, "dysmenorrhea_index", &dysmenorrheaIndexModule)
	if err != nil {
		return nil, fmt.Errorf("failed to set value of module dysmenorrhea_index :[%s]", err.Error())
	}
	// 构建生殖年龄模块
	reproductiveAgeModule := analysispb.ReproductiveAgeModule{
		Enabled: options["DisplayReproductiveAge"],
		Lookups: parseLookupsIncludeKeyContent(out["reproductive_age"]),
	}
	if out["reproductive_age"] != nil {
		if _, ok := out["reproductive_age"].(map[string]interface{}); ok {
			reproductiveAgeModule.F105 = &wrappers.Int32Value{Value: int32(out["reproductive_age"].(map[string]interface{})["f105"].(float64))}
		}
	}
	err = setModule(modules, "reproductive_age", &reproductiveAgeModule)
	if err != nil {
		return nil, fmt.Errorf("failed to set value of module reproductive_age :[%s]", err.Error())
	}
	// 构建淋巴健康模块
	lymphaticHealthModule := analysispb.LymphaticHealthModule{
		Enabled: options["DisplayLymphaticHealth"],
		Lookups: parseLookupsIncludeKeyContent(out["lymphatic_health"]),
	}
	if out["lymphatic_health"] != nil {
		if _, ok := out["lymphatic_health"].(map[string]interface{}); ok {
			lymphaticHealthModule.F107 = &wrappers.Int32Value{Value: int32(out["lymphatic_health"].(map[string]interface{})["f107"].(float64))}
		}
	}
	err = setModule(modules, "lymphatic_health", &lymphaticHealthModule)
	if err != nil {
		return nil, fmt.Errorf("failed to set value of module lymphatic_health :[%s]", err.Error())
	}
	// 构建乳腺健康模块
	breastHealthModule := analysispb.BreastHealthModule{
		Enabled: options["DisplayBreastHealth"],
		Lookups: parseLookupsIncludeKeyContent(out["breast_health"]),
	}
	if out["breast_health"] != nil {
		if _, ok := out["breast_health"].(map[string]interface{}); ok {
			breastHealthModule.M3 = &wrappers.Int32Value{Value: int32(out["breast_health"].(map[string]interface{})["m3"].(float64))}
		}
	}
	err = setModule(modules, "breast_health", &breastHealthModule)
	if err != nil {
		return nil, fmt.Errorf("failed to set value of module breast_cancers :[%s]", err.Error())
	}
	// 构建乳腺癌卵巢癌风险模块
	breastCancerModule := analysispb.BreastCancerModule{
		Enabled: options["DisplayBreastCancer"],
		Lookups: parseLookupsIncludeKeyContent(out["breast_cancer"]),
	}
	err = setModule(modules, "breast_cancer", &breastCancerModule)
	if err != nil {
		return nil, fmt.Errorf("failed to set value of module breast_cancer :[%s]", err.Error())
	}
	return modules, nil
}

func getGynecologicalInflammation(outModule interface{}) *analysispb.Lookup {
	if outModule == nil {
		return nil
	}
	nextModules, ok := outModule.(map[string]interface{})
	if ok {
		if nextModules["disease_estimate"] == nil {
			return nil
		}
		for _, value := range nextModules["disease_estimate"].([]interface{}) {

			diseaseScore, ok := value.(map[interface{}]interface{})
			if ok {
				diseaseLookup := &analysispb.Lookup{
					Key:     diseaseScore["key"].(string),
					Content: diseaseScore["content"].(string),
					Score:   diseaseScore["score"].(float64),
				}
				return diseaseLookup
			}

		}
		return nil
	}
	return nil
}

// getDiseasesMessageFromInput 用于构建风险预估模块
func getDiseasesMessageFromInput(output map[string]interface{}) ([]*analysispb.Lookup, []*analysispb.Lookup) {

	var disease []*analysispb.Lookup
	var promptMessage []*analysispb.Lookup

	for _, v := range output {
		nextModules, ok := v.(map[string]interface{})
		if ok {
			if nextModules["disease_estimate"] == nil {
				continue
			}
			for _, value := range nextModules["disease_estimate"].([]interface{}) {

				diseaseScore, ok := value.(map[interface{}]interface{})
				if ok {
					diseaseLookup := analysispb.Lookup{
						Key:     diseaseScore["key"].(string),
						Content: diseaseScore["content"].(string),
						Score:   diseaseScore["score"].(float64),
					}
					disease = append(disease, &diseaseLookup)
				}

			}
			promptMessage = append(promptMessage, parseLookupsIncludeKeyContent(nextModules)...)
		}
	}
	return disease, promptMessage
}

// parseLookupsIncludeKeyContent 解析需要 key 和 content Lookup
func parseLookupsIncludeKeyContent(outModule interface{}) []*analysispb.Lookup {
	if outModule == nil {
		return nil
	}
	module := outModule.(map[string]interface{})
	lookups, ok := module["lookups"].([]interface{})
	var lenModuleOutputLookups int
	if ok {
		lenModuleOutputLookups = len(lookups)
	} else {
		lenModuleOutputLookups = 0
	}
	moduleOutputLookups := make([]*analysispb.Lookup, lenModuleOutputLookups)
	if ok && lookups != nil {
		for idx, lookup := range lookups {
			v, ok := lookup.(map[interface{}]interface{})
			if ok {
				var content string
				if v["content"] != nil {
					content = v["content"].(string)
				}
				moduleOutputLookups[idx] = &analysispb.Lookup{
					Key:     v["key"].(string),
					Content: content,
				}
			}
		}
	}
	return moduleOutputLookups
}

// parseIncludeKeyContent 解析 content
func parseIncludeKeyContent(outModule interface{}) []*analysispb.Lookup {
	if outModule == nil {
		return nil
	}
	lookups, ok := outModule.([]interface{})
	var lenModuleOutputLookups int
	if ok {
		lenModuleOutputLookups = len(lookups)
	} else {
		lenModuleOutputLookups = 0
	}
	moduleOutputLookups := make([]*analysispb.Lookup, lenModuleOutputLookups)
	if ok && lookups != nil {
		for idx, lookup := range lookups {
			v, ok := lookup.(map[interface{}]interface{})
			if ok {
				var content string
				var score float64
				if v["content"] != nil {
					content = v["content"].(string)
				}
				if v["score_percent"] != nil {
					score = v["score_percent"].(float64)
				}
				moduleOutputLookups[idx] = &analysispb.Lookup{
					Key:     v["key"].(string),
					Content: content,
					Score:   score,
				}
			}
		}
	}
	return moduleOutputLookups
}

func setModule(modules map[string]*any.Any, name string, value proto.Message) error {
	any, err := ptypes.MarshalAny(value)
	if err != nil {
		return fmt.Errorf("failed to marshal module [%s]", name)
	}
	modules[name] = any
	return nil
}

// getPhysicalConditions 获得身体状况模块
func getPhysicalConditions(ctxData core.ContextData, options map[string]bool) analysispb.StressStateJudgment {
	stressStateJudgment := analysispb.StressStateJudgment{
		Enabled: options["DisplayPhysicalConditions"],
	}
	outStressStateJudgment, ok := ctxData.Output["stress_state_judgment"].(map[string]interface{})
	if !ok {
		return stressStateJudgment
	}
	hasStressState, ok := outStressStateJudgment["has_stress_state"].(bool)
	if ok {
		stressStateJudgment.HasStressState = hasStressState
	}
	hasDoneSports, ok := outStressStateJudgment["has_done_sports"].(bool)
	if ok {
		stressStateJudgment.HasDoneSports = hasDoneSports
	}
	hasDrinkedWine, ok := outStressStateJudgment["has_drinked_wine"].(bool)
	if ok {
		stressStateJudgment.HasDrinkedWine = hasDrinkedWine
	}
	hasHadCold, ok := outStressStateJudgment["has_had_cold"].(bool)
	if ok {
		stressStateJudgment.HasHadCold = hasHadCold
	}
	hasRhinitisEpisode, ok := outStressStateJudgment["has_rhinitis_episode"].(bool)
	if ok {
		stressStateJudgment.HasRhinitisEpisode = hasRhinitisEpisode
	}
	hasAbdominalPain, ok := outStressStateJudgment["has_abdominal_pain"].(bool)
	if ok {
		stressStateJudgment.HasAbdominalPain = hasAbdominalPain
	}
	hasViralInfection, ok := outStressStateJudgment["has_viral_infection"].(bool)
	if ok {
		stressStateJudgment.HasViralInfection = hasViralInfection
	}
	HasPhysiologicalPeriod, ok := outStressStateJudgment["has_physiological_period"].(bool)
	if ok {
		stressStateJudgment.HasPhysiologicalPeriod = HasPhysiologicalPeriod
	}
	hasOvulation, ok := outStressStateJudgment["has_ovulation"].(bool)
	if ok {
		stressStateJudgment.HasOvulation = hasOvulation
	}
	hasPregnant, ok := outStressStateJudgment["has_pregnant"].(bool)
	if ok {
		stressStateJudgment.HasPregnant = hasPregnant
	}
	hasLactation, ok := outStressStateJudgment["has_lactation"].(bool)
	if ok {
		stressStateJudgment.HasLactation = hasLactation
	}

	return stressStateJudgment
}
