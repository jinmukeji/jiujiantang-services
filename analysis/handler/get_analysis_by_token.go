package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/ae/v2/biz"
	"github.com/jinmukeji/ae/v2/engine/core"
	"github.com/jinmukeji/go-pkg/v2/age"
	"github.com/jinmukeji/jiujiantang-services/analysis/mysqldb"
	analysispb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/analysis/v1"
	jinmuidpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	ptypesv2 "github.com/jinmukeji/proto/v3/gen/micro/idl/ptypes/v2"
)

// GetAnalyzeResultByToken 根据分享 token 得到分析结果
func (j *AnalysisManagerService) GetAnalyzeResultByToken(ctx context.Context, req *analysispb.GetAnalyzeResultByTokenRequest, resp *analysispb.GetAnalyzeResultByTokenResponse) error {

	recordToken, err := j.database.FindAnalysisBodyByToken(req.Token)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find analysis body by token %s: %s", req.Token, err.Error()))
	}
	if recordToken == nil {
		return NewError(ErrInvalidAnalysisStatus, fmt.Errorf("analysis status of token %s is invalid", req.Token))
	}
	record, _ := j.database.FindRecordByRecordID(recordToken.RecordID)
	// 如果分析状态不是完成，则无法获得分析结果
	if record.AnalyzeStatus != mysqldb.AnalysisStatusCompeleted {
		return NewError(ErrDatabase, fmt.Errorf("analyze status [%d] of token [%s] is wrong", record.AnalyzeStatus, req.Token))
	}
	// 获得用户个人信息
	reqGetUserProfile := new(jinmuidpb.GetUserProfileRequest)
	reqGetUserProfile.IsSkipVerifyToken = true
	reqGetUserProfile.UserId = record.UserID
	respGetUserProfile, err := j.jinmuidSvc.GetUserProfile(ctx, reqGetUserProfile)
	if err != nil {
		return NewError(ErrInvalidUser, fmt.Errorf("failed to get user profile of user [%d]: %s", record.UserID, err.Error()))
	}
	userProfile := respGetUserProfile.GetProfile()
	clientID := record.ClientID

	// 加载 preset
	err = j.biz.LoadPresets(j.presetsFilePath, biz.MeasurementJudgmentModule)
	if err != nil {
		return NewError(ErrLoadConfig, fmt.Errorf("failed to load presets: %s", err.Error()))
	}
	preset, err := j.biz.SelectPresetByClientID(clientID)
	if err != nil {
		return NewError(ErrClientID, fmt.Errorf("failed to get preset by client_id %s: %s", clientID, err.Error()))
	}

	// 构建 ae 的输入
	analyzeBody := &AnalyzeBody{}
	err = json.Unmarshal([]byte(record.AnalyzeBody), analyzeBody)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to unmarshal analyze body of record [%d]: %s", record.RecordID, err.Error()))
	}
	ctxData, err := buildAEInput(userProfile, record, analyzeBody.QuestionAnswers, analyzeBody.TransactionID, true)
	if err != nil {
		return NewError(ErrLoadConfig, fmt.Errorf("failed to build ae input of record [%d]: %s", record.RecordID, err.Error()))
	}

	language := analyzeBody.Language

	// 根据客户端 ID 获得对应的提示助手名称
	assistantName := getAssistantNameByClientId(clientID)

	// 构建分析需要的上下文信息
	presetAgent := biz.PresetAgent{
		// 助手名称
		AssistantName: assistantName,
		// 用户昵称
		ReporterName: userProfile.GetNickname(),
	}

	errRunEngine := preset.RunEngine(ctxData, language, presetAgent)
	if errRunEngine != nil {
		return NewError(ErrAEError, fmt.Errorf("error occurs when run engine: %s", errRunEngine.Error()))
	}
	// 如果还有需要提问的问题，则分析错误
	if !ctxData.Output["has_answered_all_questions"].(bool) {
		return NewError(ErrAEError, fmt.Errorf("error occurs when run engine: %s", errRunEngine.Error()))
	}

	// 构建返回的内容
	hasAbnormalMeasurementFromAeOutput := checkIfHasAbnormalMeasurementFromAeOutput(*ctxData)
	displayOptions := getDisplayOption(record.ClientID, userProfile.GetGender(), checkIfHasStressStateFromAeOutput(*ctxData), hasAbnormalMeasurementFromAeOutput, true)

	respModules, err := getModulesFromAEOutput(*ctxData, displayOptions, userProfile.GetGender())
	if err != nil {
		return NewError(ErrBuildReturnModule, fmt.Errorf("failed to get modules from content output"))
	}
	// 下面构建报告内容里与引擎输出无关的模块
	restModules, err := j.buildModulesNotAboutAE(record, displayOptions, true)
	if err != nil {
		return NewError(ErrBuildReturnModule, fmt.Errorf("failed to get modules"))
	}
	for key, value := range restModules {
		respModules[key] = value
	}

	userProfileModule, err := buildUserProfileModule(userProfile, record, displayOptions)
	if err != nil {
		return NewError(ErrBuildReturnModule, fmt.Errorf("failed to build user profile module"))
	}

	// 构建测量上下文模块
	pulseTestModule, err := buildPulseTestModule(record, displayOptions)
	if err != nil {
		return NewError(ErrBuildReturnModule, fmt.Errorf("failed to build pulse test module"))
	}

	remarkModule := buildRemarkModule(record, displayOptions)

	protoAnalysisFinishTime, err := ptypes.TimestampProto(record.CreatedAt)
	if err != nil {
		return NewError(ErrBuildReturnModule, fmt.Errorf("failed to get timestamp of created_at of record [%d]: %s", record.RecordID, err.Error()))
	}

	resp.RecordId = int32(record.RecordID)
	resp.ReportVersion = DefaultReportVersion
	resp.Report = &analysispb.ReportContent{
		RecordId:    int32(record.RecordID),
		UserProfile: userProfileModule,
		PulseTest:   pulseTestModule,
		Remark:      remarkModule,
		Modules:     respModules,
		CreatedTime: protoAnalysisFinishTime,
	}
	resp.TransactionId = record.TransactionNumber

	return nil

}

func buildUserProfileModule(userProfile *jinmuidpb.UserProfile, record *mysqldb.Record, displayOptions map[string]bool) (*analysispb.UserProfile, error) {
	// 构建返回的 UserProfile
	birthday, err := ptypes.Timestamp(userProfile.GetBirthdayTime())
	if err != nil {
		return nil, NewError(ErrInvalidAge, fmt.Errorf("failed to get timestamp birthday of user [%d]: %s", record.UserID, err.Error()))
	}
	respUserProfile := &analysispb.UserProfile{
		Enabled:         displayOptions["DisplayUserProfile"],
		RecordId:        record.UserID,
		Nickname:        userProfile.GetNickname(),
		NicknameInitial: userProfile.GetNicknameInitial(),
		Birthday: &ptypesv2.Date{
			Year:  int32(birthday.Year()),
			Month: int32(birthday.Month()),
			Day:   int32(birthday.Day()),
		},
		Gender: userProfile.GetGender(),
		Height: userProfile.GetHeight(),
		Weight: userProfile.GetWeight(),
	}
	return respUserProfile, nil
}

func buildPulseTestModule(record *mysqldb.Record, displayOptions map[string]bool) (*analysispb.PulseTest, error) {
	// 构建返回的测量上下文，也就是手指信息
	protoFinger, err := mapDBFingerToProto(record.Finger)
	if err != nil {
		return nil, NewError(ErrInvalidFinger, fmt.Errorf("invalid finger %d", record.Finger))
	}
	return &analysispb.PulseTest{
		Fingers: []analysispb.Finger{protoFinger},
	}, nil

}

// 构建备注模块
func buildRemarkModule(record *mysqldb.Record, displayOptions map[string]bool) *analysispb.Remark {
	return &analysispb.Remark{
		Enabled: displayOptions["DisplayRemark"],
		Content: record.Remark,
	}

}

func buildAEInputAnswers(moduleAnswers map[string]AnalysisReportAnswers) map[string]Answer {
	returnAnswers := make(map[string]Answer)
	for module, answers := range moduleAnswers {
		answerKeys := make(map[string]AnswerKeys)
		for _, a := range answers {
			answerKeys[a.QuestionKey] = a.AnswerKeys
		}
		returnAnswers[module] = answerKeys
	}
	return returnAnswers
}

// 站姿走滤镜
func buildAEInput(userProfile *jinmuidpb.UserProfile, record *mysqldb.Record, answers map[string]AnalysisReportAnswers, cid string, isMeasurementJudgment bool) (*core.ContextData, error) {
	aeGender, err := mapProtoGenderToAE(userProfile.Gender)
	if err != nil {
		return nil, fmt.Errorf("invalid gender [%s]: %s", userProfile.Gender, err.Error())
	}
	birthday, _ := ptypes.Timestamp(userProfile.GetBirthdayTime())

	if record.MeasurementPosture == mysqldb.MeasurementPostureStanging {
		record.C0, record.C1, record.C2, record.C3, record.C4, record.C5, record.C6, record.C7 = ConvertToStandingCFloat64Values(userProfile.Gender, int(record.C0), int(record.C1), int(record.C2), int(record.C3), int(record.C4), int(record.C5), int(record.C6), int(record.C7))
	}

	clientID := record.ClientID
	// 根据客户端 ID 获得对应的分析引擎的开关
	analysisOptions := getAnalysisOptionsByClientID(clientID)
	age := int(age.Age(birthday))
	if age < 0 || age > 120 {
		return nil, fmt.Errorf("age %d of record is invalid", age)
	}
	if userProfile.Weight < 30 || userProfile.Weight > 500 {
		return nil, fmt.Errorf("weight %d of record is invalid", userProfile.Weight)
	}
	if userProfile.Height < 50 || userProfile.Height > 250 {
		return nil, fmt.Errorf("height %d of record is invalid", userProfile.Height)
	}
	input := core.D{
		"c0":               IntValBoundedBy10FromFloat(record.C0) * 10,
		"c1":               IntValBoundedBy10FromFloat(record.C1) * 10,
		"c2":               IntValBoundedBy10FromFloat(record.C2) * 10,
		"c3":               IntValBoundedBy10FromFloat(record.C3) * 10,
		"c4":               IntValBoundedBy10FromFloat(record.C4) * 10,
		"c5":               IntValBoundedBy10FromFloat(record.C5) * 10,
		"c6":               IntValBoundedBy10FromFloat(record.C6) * 10,
		"c7":               IntValBoundedBy10FromFloat(record.C7) * 10,
		"gender":           aeGender,
		"age":              age,
		"heart_rate":       int(record.HeartRate),
		"weight":           int(userProfile.GetWeight()),
		"height":           int(userProfile.GetHeight()),
		"question_answers": buildAEInputAnswers(answers),
	}
	bag := core.D{
		// app_id 先传 client_id
		"app_id":                          record.ClientID,
		"is_measurement_judgment":         isMeasurementJudgment,
		"enabled_anxiety":                 analysisOptions["EnabledAnxiety"],
		"enabled_cd":                      analysisOptions["EnabledCd"],
		"enabled_cc":                      analysisOptions["EnabledCc"],
		"enabled_sd":                      analysisOptions["EnabledSd"],
		"enabled_factor_interpretation":   analysisOptions["EnabledFactorInterpretation"],
		"enabled_emotional_health":        analysisOptions["EnabledEmotionalHealth"],
		"enabled_lymphatic_health":        analysisOptions["EnabledLymphaticHealth"],
		"enabled_blood_sugar":             analysisOptions["EnabledBloodSugar"],
		"enabled_blood_pressure":          analysisOptions["EnabledBloodPressure"],
		"enabled_hyperlipidemia":          analysisOptions["EnabledHyperlipidemia"],
		"enabled_acute_pharyngitis":       analysisOptions["EnabledAcutePharyngitis"],
		"enabled_spine":                   analysisOptions["EnabledSpine"],
		"enabled_spinal_disease":          analysisOptions["EnabledSpinalDisease"],
		"enabled_chronic_cough":           analysisOptions["EnabledChronicCough"],
		"enabled_immunity":                analysisOptions["EnabledImmunity"],
		"enabled_cerebral_insufficiency":  analysisOptions["EnabledCerebralInsufficiency"],
		"enabled_fatigue_and_pressure":    analysisOptions["EnabledFatigueAndPressure"],
		"enabled_renal_dysfunction":       analysisOptions["EnabledRenalDysfunction"],
		"enabled_chd":                     analysisOptions["EnabledChd"],
		"enabled_depression":              analysisOptions["EnabledDepression"],
		"enabled_sleep_problems":          analysisOptions["EnabledSleepProblems"],
		"enabled_hypomotility_of_stomach": analysisOptions["EnabledHypomotilityOfStomach"],
		"enabled_gastritis":               analysisOptions["EnabledGastritis"],
		"enabled_menstrual_health_index":  analysisOptions["EnabledMenstrualHealthIndex"],
		"enabled_reproductive_age":        analysisOptions["EnabledReproductiveAge"],
		"enabled_uterine_health":          analysisOptions["EnabledUterineHealth"],
		"enabled_inflammation_risk":       analysisOptions["EnabledInflammationRisk"],
		"enabled_inflammation":            analysisOptions["EnabledInflammation"],
		"enabled_breast_health":           analysisOptions["EnabledBreastHealth"],
		"enabled_facial_skin":             analysisOptions["EnabledFacialSkin"],
		"enabled_hormone_level":           analysisOptions["EnabledHormoneLevel"],
		"enabled_breast_cancer":           analysisOptions["EnabledBreastCancer"],
		"enabled_facial_skin_male":        analysisOptions["EnabledFacialSkinMale"],
		"enabled_hormone_level_male":      analysisOptions["EnabledHormoneLevelMale"],
		"enabled_prostate_disease":        analysisOptions["EnabledProstateDisease"],
		"enabled_dietary_advice":          analysisOptions["EnabledDietaryAdvice"],
		"enabled_sports_advice":           analysisOptions["EnabledSportsAdvice"],
		"enabled_chinese_medicine_advice": analysisOptions["EnabledChineseMedicineAdvice"],
		"enabled_physical_therapy_advice": analysisOptions["EnabledPhysicalTherapyAdvice"],
		"enabled_stress_state_judgment":   analysisOptions["EnabledStressStateJudgment"],
		"enabled_dirty_dialectic":         analysisOptions["EnabledDirtyDialectic"],
	}
	ctxData := core.NewContextData(cid, input, bag, nil)
	return ctxData, nil
}
