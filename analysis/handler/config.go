package handler

import (
	"errors"
	"fmt"

	ptypesv2 "github.com/jinmukeji/proto/gen/micro/idl/ptypes/v2"
)

const (
	// dengyunClient 登云的client
	dengyunClient = "dengyun-10001"
	// moaiClient 摩爱的client
	moaiClient = "moai-10001"
	// hongjingtangClient 弘经堂的client
	hongjingtangClient = "hongjingtang-10001"
	// kangmeiClient 康美的client
	kangmeiClient = "kangmei-10001"
	// AssistantName 提示助手的名字
	AssistantName = "金姆宝宝"
)

// getAnalysisOptionsByClientID 根据客户端 ID 获得对应的分析引擎的开关
func getAnalysisOptionsByClientID(clientID string) map[string]bool {
	switch clientID {
	case dengyunClient:
		return map[string]bool{
			"EnabledCd":                  true,
			"EnabledStressStateJudgment": true,
		}
	case moaiClient:
		return map[string]bool{
			"EnabledCd":                  true,
			"EnabledCc":                  true,
			"EnabledSd":                  true,
			"EnabledStressStateJudgment": true,
		}
	case hongjingtangClient:
		return map[string]bool{
			"EnabledCd":                  true,
			"EnabledCc":                  true,
			"EnabledSd":                  true,
			"EnabledStressStateJudgment": true,
		}
	case kangmeiClient:
		return map[string]bool{
			"EnabledCd":                  true,
			"EnabledCc":                  true,
			"EnabledSd":                  true,
			"EnabledStressStateJudgment": true,
		}
	}
	return map[string]bool{
		"EnabledCd":                    true,
		"EnabledCc":                    true,
		"EnabledSd":                    true,
		"EnabledFactorInterpretation":  true,
		"EnabledBreastHealth":          true,
		"EnabledMenstrualHealth":       true,
		"EnabledFatigueAndPressure":    true,
		"EnabledCerebralInsufficiency": true,
		"EnabledHormoneLevelMale":      true,
		"EnabledUterineHealth":         true,
		"EnabledGastritis":             true,
		"EnabledSpine":                 true,
		"EnabledImmunity":              true,
		"EnabledChd":                   true,
		"EnabledFacialSkin":            true,
		"EnabledSleepProblems":         true,
		"EnabledHormoneLevel":          true,
		"EnabledDepression":            true,
		"EnabledReproductiveAge":       true,
		"EnabledAcutePharyngitis":      true,
		"EnabledChronicCough":          true,
		"EnabledBreastCancer":          true,
		"EnabledLymphaticHealth":       true,
		"EnabledBloodSugar":            true,
		"EnabledBloodPressure":         true,
		"EnabledHyperlipidemia":        true,
		"EnabledSpinalDisease":         true,
		"EnabledRenalDysfunction":      true,
		"EnabledHypomotilityOfStomach": true,
		"EnabledInflammation":          true,
		"EnabledFacialSkinMale":        true,
		"EnabledChineseMedicineAdvice": true,
		"EnabledDietaryAdvice":         true,
		"EnabledPhysicalTherapyAdvice": true,
		"EnabledSportsAdvice":          true,
		"EnabledStressStateJudgment":   true,
		"EnabledDirtyDialectic":        true,
		"EnabledAnxiety":               true,
		"EnabledInflammationRisk":      true,
		"EnabledProstateDisease":       true,
	}
}

// getAssistantNameByClientId 根据客户端 ID 获得对应的提示助手名称
func getAssistantNameByClientId(clientID string) string {
	return AssistantName
}

const (
	// AEGenderMale 男性
	AEGenderMale string = "male"
	// AEGenderFemale 女性
	AEGenderFemale string = "female"
)

// mapProtoGenderToAE 将 proto 类型的 gender 转换为运行分析引擎使用的 gender
func mapProtoGenderToAE(gender ptypesv2.Gender) (string, error) {
	switch gender {
	case ptypesv2.Gender_GENDER_FEMALE:
		return AEGenderFemale, nil
	case ptypesv2.Gender_GENDER_MALE:
		return AEGenderMale, nil
	case ptypesv2.Gender_GENDER_INVALID:
		return "", fmt.Errorf("invalid proto gender %d", gender)
	case ptypesv2.Gender_GENDER_UNSET:
		return "", fmt.Errorf("invalid proto gender %d", gender)
	}
	return "", errors.New("invalid proto gender")
}

// getDisplayOption 根据 appId、性别、是否处于应激态,是否测量异常获得最终显示的模块
func getDisplayOption(appId string, gender ptypesv2.Gender, hasStressState bool, hasAbnormalMeasurement bool, isMeasurement bool) map[string]bool {
	options := map[string]bool{}
	// 如果测量异常，则只显示测量异常模块
	if hasAbnormalMeasurement {
		options["DisplayMeasurementJudgment"] = true
		return options
	}
	// 测量正常会显示的模块
	options["DisplayNavbar"] = true
	options["DisplayTags"] = true

	// 如果是测量，需要显示的模块，这些模块在周报月报的情况下不显示
	if isMeasurement {
		// 显示用户个人信息
		options["DisplayUserProfile"] = true
		// 显示备注
		options["DisplayRemark"] = true
		// 显示心率
		options["DisplayHeartRate"] = true
		// 显示局部脉搏波
		options["DisplayPartialData"] = true
	}
	// 显示理疗指数
	options["DisplayPhysicalTherapyIndex"] = true
	// 显示经络柱状图
	options["DisplayCcBarChart"] = true

	//  处于应激态会显示的模块
	if hasStressState {
		// 显示温馨提示
		options["DisplayTips"] = true
		// 显示身体状况
		options["DisplayPhysicalConditions"] = true
		return options
	}
	// 不处于应激态会显示的模块
	if !hasStressState {
		// 显示风险预估
		options["DisplayRiskEstimate"] = true
		// 显示体质辩证
		options["DisplayCdExplain"] = true
		// 显示脏腑辩证
		options["DisplaySdExplain"] = true
		// 显示调理建议
		options["DisplayConditioningAdvice"] = true
		// 显示食疗建议
		options["DisplayDietaryAdvice"] = true
		// 显示运动建议
		options["DisplaySportsAdvice"] = true
		// 显示中药建议
		// options["DisplayChineseMedicineAdvice"] = true
		// 显示理疗建议
		options["DisplayPhysicalTherapyAdvice"] = true
		// 显示局部脉搏波
		options["DisplayDisplayPartialData"] = true
		// 显示经络柱解读
		options["DisplayCcExplain"] = true
		// 下面注释调的都是一体机的模块
		// 因为App不显示一体机相关模块，所以注释
		// 显示情绪健康
		// options["DisplayF3"] = true
		// // 显示面部美肤
		// options["DisplayFacialSkin"] = true
		// // 显示激素水平
		// options["DisplayHormoneLevel"] = true
		// // 显示淋巴健康
		// options["DisplayLymphaticHealth"] = true

	}

	// 如果性别为女性，会显示的模块
	if gender == ptypesv2.Gender_GENDER_FEMALE {
		// // 显示妇科疾病风险预估
		// options["DisplayGynecologicalDiseaseRisk"] = true
		// // 显示妇科炎症
		// options["DisplayGynecologicalInflammation"] = true
		// // 显示子宫健康
		// options["DisplayUterineHealth"] = true
		// // 显示月经（天葵）健康
		// options["DisplayMenstrualHealth"] = true
		// // 显示月经不调
		// options["DisplayMenstruationIrregular"] = true
		// // 显示痛经
		// options["DisplayDysmenorrhea"] = true
		// // 显示生殖年龄
		// options["DisplayReproductiveAge"] = true
		// // 显示乳腺健康
		// options["DisplayBreastHealth"] = true
		// // 显示乳腺癌卵巢癌风险
		// options["DisplayBreastCancer"] = true
	}

	return options

}
