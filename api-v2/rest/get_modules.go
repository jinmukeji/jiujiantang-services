package rest

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	age "github.com/jinmukeji/go-pkg/age"
	analysispb "github.com/jinmukeji/proto/gen/micro/idl/jm/analysis/v1"
)

func getAnalysisAskQuestions(resp *analysispb.GetAnalyzeResultResponse) map[string]Questions {
	// 如果提问的问题非空
	if resp.GetQuestions() != nil {
		askQuestions := mapProtoQuestionsToAskQuestions(resp.GetQuestions())
		return askQuestions
	}
	return nil
}

// 个人信息模块
func getUserProfileModule(profile *analysispb.UserProfile) (UserProfileModule, error) {
	protoBirthday := profile.GetBirthday()
	birthday := time.Date(int(protoBirthday.GetYear()), time.Month(protoBirthday.GetMonth()), int(protoBirthday.GetDay()), 0, 0, 0, 0, time.UTC)
	gender, _ := mapProtoGenderToRest(profile.GetGender())
	age := int64(age.Age(birthday))
	userProfile := UserProfileModule{
		UserID:          int64(profile.GetRecordId()),
		Nickname:        profile.GetNickname(),
		NicknameInitial: profile.GetNicknameInitial(),
		Birthday:        birthday,
		Age:             age,
		Gender:          gender,
		Height:          int64(profile.GetHeight()),
		Weight:          int64(profile.GetWeight()),
	}
	return userProfile, nil
}

// 测量上下文模块
func getPulseTestModule(pulseTest *analysispb.PulseTest) (PulseTestModule, error) {
	fingers := make([]int, len(pulseTest.GetFingers()))
	for idx, value := range pulseTest.GetFingers() {
		mapProtoFingerToRest, _ := mapAnalysisProtoFingerToRest(value)
		fingers[idx] = mapProtoFingerToRest
	}
	pulseTestModule := PulseTestModule{
		Fingers: fingers,
	}
	return pulseTestModule, nil

}

// 备注模块
func getRemarkModule(remark *analysispb.Remark) (RemarkModule, error) {
	remarkModule := RemarkModule{
		Enabled: remark.GetEnabled(),
		Content: remark.GetContent(),
	}
	return remarkModule, nil

}

// getAnalysisModules 获得引擎分析结果中的模块
func getAnalysisModules(reportModules map[string]*any.Any) (*AnalysisReportContent, error) {
	content := &AnalysisReportContent{}
	// 构建各个报告中的模块
	// 心率模块
	if reportModules["heart_rate"] != nil {
		heartRate := &analysispb.HeartRateModule{}
		err := ptypes.UnmarshalAny(reportModules["heart_rate"], heartRate)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal module [%s] from response: %s", "heart_rate", err.Error())
		}
		heartRateModule := HeartRateModule{
			Enabled:          heartRate.GetEnabled(),
			AverageHeartRate: heartRate.GetAverageHeartRate(),
			HighestHeartRate: heartRate.GetHighestHeartRate(),
			LowestHeartRate:  heartRate.GetLowestHeartRate(),
		}
		content.HeartRate = heartRateModule
	}
	// 风险预估模块
	if reportModules["risk_estimate"] != nil {
		riskEstimate := &analysispb.RiskEstimateModule{}
		err := ptypes.UnmarshalAny(reportModules["risk_estimate"], riskEstimate)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal module [%s] from response: %s", "risk_estimate", err.Error())
		}
		riskEstimateModule := RiskEstimateModule{
			Enabled:         riskEstimate.GetEnabled(),
			DiseaseEstimate: getRestLookupFromProto(riskEstimate.GetDiseaseEstimate()),
			PromptMessage:   getRestLookupFromProto(riskEstimate.GetPromptMessage()),
		}
		content.RiskEstimate = riskEstimateModule
	}

	// 体质辩证模块
	if reportModules["physical_dialectics"] != nil {
		physicalDialectics := &analysispb.PhysicalDialecticsModule{}
		err := ptypes.UnmarshalAny(reportModules["physical_dialectics"], physicalDialectics)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal module [%s] from response: %s", "physical_dialectics", err.Error())
		}
		physicalDialecticsModule := PhysicalDialecticsModule{
			Enabled: physicalDialectics.GetEnabled(),
			Lookups: getRestLookupFromProto(physicalDialectics.GetLookups()),
		}
		content.PhysicalDialectics = physicalDialecticsModule
	}

	// 脏腑辩证模块模块
	if reportModules["dirty_dialectic"] != nil {
		dirtyDialectic := &analysispb.DirtyDialecticModule{}
		err := ptypes.UnmarshalAny(reportModules["dirty_dialectic"], dirtyDialectic)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal module [%s] from response: %s", "dirty_dialectic", err.Error())
		}
		dirtyDialecticModule := DirtyDialecticModule{
			Enabled: dirtyDialectic.GetEnabled(),
			Lookups: getRestLookupFromProto(dirtyDialectic.GetLookups()),
		}
		content.DirtyDialectic = dirtyDialecticModule
	}

	// 理疗指数模块
	if reportModules["physical_therapy_index"] != nil {
		physicalTherapyIndex := &analysispb.PhysicalTherapyIndexModule{}
		err := ptypes.UnmarshalAny(reportModules["physical_therapy_index"], physicalTherapyIndex)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal module [%s] from response: %s", "physical_therapy_index", err.Error())
		}

		physicalTherapyIndexModule := PhysicalTherapyIndexModule{
			Enabled: physicalTherapyIndex.GetEnabled(),
			Lookups: getRestLookupFromProto(physicalTherapyIndex.GetLookups()),
			F0:      physicalTherapyIndex.GetF0().GetValue(),
			F1:      physicalTherapyIndex.GetF1().GetValue(),
			F2:      physicalTherapyIndex.GetF2().GetValue(),
			F3:      physicalTherapyIndex.GetF3().GetValue(),
		}
		content.PhysicalTherapyIndex = physicalTherapyIndexModule
	}

	// 调理建议模块
	if reportModules["conditioning_advice"] != nil {
		conditioningAdvice := &analysispb.ConditioningAdviceModule{}
		err := ptypes.UnmarshalAny(reportModules["conditioning_advice"], conditioningAdvice)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal module [%s] from response: %s", "conditioning_advice", err.Error())
		}
		conditioningAdviceModule := ConditioningAdviceModule{
			Enabled: conditioningAdvice.GetEnabled(),
			DietaryAdvice: &DietaryAdviceModule{
				Enabled: conditioningAdvice.GetDietaryAdvice().GetEnabled(),
				Lookups: getRestLookupFromProto(conditioningAdvice.GetDietaryAdvice().GetLookups()),
			},
			SportsAdvice: &SportsAdviceModule{
				Enabled: conditioningAdvice.GetSportsAdvice().GetEnabled(),
				Lookups: getRestLookupFromProto(conditioningAdvice.GetSportsAdvice().GetLookups()),
			},
			ChineseMedicineAdvice: &ChineseMedicineAdviceModule{
				Enabled: conditioningAdvice.GetChineseMedicineAdvice().GetEnabled(),
				Lookups: getRestLookupFromProto(conditioningAdvice.GetChineseMedicineAdvice().GetLookups()),
			},
			PhysicalTherapyAdvice: &PhysicalTherapyAdviceModule{
				Enabled: conditioningAdvice.GetPhysicalTherapyAdvice().GetEnabled(),
				Lookups: getRestLookupFromProto(conditioningAdvice.GetPhysicalTherapyAdvice().GetLookups()),
			},
		}
		content.ConditioningAdvice = conditioningAdviceModule
	}

	// 局部脉搏波模块
	if reportModules["partial_pulse_wave"] != nil {
		partialPulseWave := &analysispb.PartialPulseWaveModule{}
		err := ptypes.UnmarshalAny(reportModules["partial_pulse_wave"], partialPulseWave)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal module [%s] from response: %s", "partial_pulse_wave", err.Error())
		}
		partialPulseWaveModule := PartialPulseWaveModule{
			Enabled: partialPulseWave.GetEnabled(),
			Points:  partialPulseWave.GetPoints(),
		}
		content.PartialPulseWave = partialPulseWaveModule
	}

	//  经络柱状图模块
	if reportModules["meridian_bar_chart"] != nil {
		meridianBarChart := &analysispb.MeridianBarChartModule{}
		err := ptypes.UnmarshalAny(reportModules["meridian_bar_chart"], meridianBarChart)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal module [%s] from response: %s", "meridian_bar_chart", err.Error())
		}
		testingTime, _ := ptypes.Timestamp(meridianBarChart.GetMeridianValue().GetTestTime())
		meridianBarChartModule := MeridianBarChartModule{
			Enabled: meridianBarChart.GetEnabled(),
			MeridianValue: &CInfo{
				C0:       meridianBarChart.GetMeridianValue().GetC0(),
				C1:       meridianBarChart.GetMeridianValue().GetC1(),
				C2:       meridianBarChart.GetMeridianValue().GetC2(),
				C3:       meridianBarChart.GetMeridianValue().GetC3(),
				C4:       meridianBarChart.GetMeridianValue().GetC4(),
				C5:       meridianBarChart.GetMeridianValue().GetC5(),
				C6:       meridianBarChart.GetMeridianValue().GetC6(),
				C7:       meridianBarChart.GetMeridianValue().GetC7(),
				TestTime: testingTime,
			},
		}
		content.MeridianBarChart = meridianBarChartModule
	}

	// 经络解读模块
	if reportModules["meridian_explain"] != nil {
		meridianExplain := &analysispb.MeridianExplainModule{}
		err := ptypes.UnmarshalAny(reportModules["meridian_explain"], meridianExplain)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal module [%s] from response: %s", "meridian_explain", err.Error())
		}
		meridianExplainModule := MeridianExplainModule{
			Enabled: meridianExplain.GetEnabled(),
			Lookups: getRestLookupFromProto(meridianExplain.GetLookups()),
		}
		content.MeridianExplain = meridianExplainModule
	}

	// 测量异常判断模块
	if reportModules["measurement_judgment"] != nil {
		measurementJudgments := &analysispb.MeasurementJudgment{}
		err := ptypes.UnmarshalAny(reportModules["measurement_judgment"], measurementJudgments)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal module [%s] from response: %s", "measurement_judgment", err.Error())
		}
		measurementJudgmentModule := MeasurementJudgmentModule{
			Enabled: measurementJudgments.GetEnabled(),
			Lookups: getRestLookupFromProto(measurementJudgments.GetLookups()),
		}
		content.MeasurementJudgment = measurementJudgmentModule
	}

	// 温馨提示模块
	if reportModules["tips"] != nil {
		tips := &analysispb.Tips{}
		err := ptypes.UnmarshalAny(reportModules["tips"], tips)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal module [%s] from response: %s", "tips", err.Error())
		}
		tipsModule := TipsModule{
			Enabled: tips.GetEnabled(),
			Lookups: getRestLookupFromProto(tips.GetLookups()),
		}
		content.Tips = tipsModule
	}

	// 应激态模块
	if reportModules["stress_state_judgment"] != nil {
		stressStateJudgment := &analysispb.StressStateJudgment{}
		err := ptypes.UnmarshalAny(reportModules["stress_state_judgment"], stressStateJudgment)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal module [%s] from response: %s", "stress_state_judgment", err.Error())
		}
		stressStateJudgmentModule := StressStateJudgmentModule{
			Enabled:                stressStateJudgment.GetEnabled(),
			HasStressState:         stressStateJudgment.GetHasStressState(),
			HasDoneSports:          stressStateJudgment.GetHasDoneSports(),
			HasDrinkedWine:         stressStateJudgment.GetHasDrinkedWine(),
			HasHadCold:             stressStateJudgment.GetHasHadCold(),
			HasRhinitisEpisode:     stressStateJudgment.GetHasRhinitisEpisode(),
			HasAbdominalPain:       stressStateJudgment.GetHasAbdominalPain(),
			HasViralInfection:      stressStateJudgment.GetHasViralInfection(),
			HasPhysiologicalPeriod: stressStateJudgment.GetHasPhysiologicalPeriod(),
			HasOvulation:           stressStateJudgment.GetHasOvulation(),
			HasPregnant:            stressStateJudgment.GetHasPregnant(),
			HasLactation:           stressStateJudgment.GetHasLactation(),
		}
		content.StressStateJudgment = stressStateJudgmentModule
	}

	return content, nil
}

func getRestLookupFromProto(lookups []*analysispb.Lookup) []*Lookup {
	restLookups := make([]*Lookup, len(lookups))
	for idx, value := range lookups {
		restLookups[idx] = &Lookup{
			Key:     value.GetKey(),
			Content: value.GetContent(),
			Score:   value.GetScore(),
			LinkKey: value.GetLinkKey(),
		}
	}
	return restLookups
}

// getAnswers 建立模块名到回答到对应关系
func mapProtoQuestionsToAskQuestions(answers map[string]*analysispb.Questions) map[string]Questions {
	questionAnswers := make(map[string]Questions)
	for key, value := range answers {
		webQuestion := make([]AnalysisReportQuestion, len(value.GetQuestions()))
		for idxQuestion, valueQuestion := range value.GetQuestions() {
			webQuestionChoice := make([]AnalysisReportChoice, len(valueQuestion.GetChoices()))
			for idxChoice, valueChoice := range valueQuestion.GetChoices() {
				webQuestionChoice[idxChoice] = AnalysisReportChoice{
					Key:          valueChoice.GetChoiceKey(),
					Content:      valueChoice.GetContent(),
					ConflictKeys: valueChoice.GetConflictKeys(),
					Selected:     valueChoice.GetSelected(),
				}
			}
			webQuestion[idxQuestion] = AnalysisReportQuestion{
				Key:     valueQuestion.GetQuestionKey(),
				Content: valueQuestion.GetContent(),
				Type:    valueQuestion.GetType(),
				Choices: webQuestionChoice,
			}
		}
		questionAnswers[key] = webQuestion
	}
	return questionAnswers
}

func mapAnalysisProtoFingerToRest(protoFinger analysispb.Finger) (int, error) {
	switch protoFinger {
	case analysispb.Finger_FINGER_INVALID:
		return FingerInvalid, fmt.Errorf("invalid proto finger %d", protoFinger)
	case analysispb.Finger_FINGER_UNSET:
		return FingerInvalid, fmt.Errorf("invalid proto finger %d", protoFinger)
	case analysispb.Finger_FINGER_LEFT_1:
		return FingerLeft1, nil
	case analysispb.Finger_FINGER_LEFT_2:
		return FingerLeft2, nil
	case analysispb.Finger_FINGER_LEFT_3:
		return FingerLeft3, nil
	case analysispb.Finger_FINGER_LEFT_4:
		return FingerLeft4, nil
	case analysispb.Finger_FINGER_LEFT_5:
		return FingerLeft5, nil
	case analysispb.Finger_FINGER_RIGHT_1:
		return FingerRight1, nil
	case analysispb.Finger_FINGER_RIGHT_2:
		return FingerRight2, nil
	case analysispb.Finger_FINGER_RIGHT_3:
		return FingerRight3, nil
	case analysispb.Finger_FINGER_RIGHT_4:
		return FingerRight4, nil
	case analysispb.Finger_FINGER_RIGHT_5:
		return FingerRight5, nil
	}
	return -1, fmt.Errorf("invalid proto finger %d", protoFinger)
}
