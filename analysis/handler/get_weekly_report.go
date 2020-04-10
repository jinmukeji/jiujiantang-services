package handler

import (
	"context"
	"fmt"

	"github.com/jinmukeji/ae/v2/biz"
	"github.com/jinmukeji/jiujiantang-services/analysis/mysqldb"
	jinmuidpb "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	analysispb "github.com/jinmukeji/proto/gen/micro/idl/jm/analysis/v1"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/analysis/v1"
)

// GetWeeklyReportContent 得到分析报告的内容
func (j *AnalysisManagerService) GetWeeklyAnalyzeResult(ctx context.Context, req *proto.GetWeeklyAnalyzeResultRequest, resp *proto.GetWeeklyAnalyzeResultResponse) error {
	// 获得用户个人信息
	reqGetUserProfile := new(jinmuidpb.GetUserProfileRequest)
	// 如果请求的不是当前用户，则不需要 token 验证
	token, _ := TokenFromContext(ctx)
	ownerID, _ := j.database.FindUserIDByToken(token)
	if ownerID != req.GetUserId() {
		reqGetUserProfile.IsSkipVerifyToken = true
	}

	reqGetUserProfile.UserId = req.GetUserId()
	respGetUserProfile, err := j.jinmuidSvc.GetUserProfile(ctx, reqGetUserProfile)
	if err != nil {
		return NewError(ErrInvalidUser, fmt.Errorf("failed to get user profile of user [%d]: %s", req.GetUserId(), err.Error()))
	}
	userProfile := respGetUserProfile.GetProfile()
	clientID, ok := ClientIDFromContext(ctx)
	if !ok {
		return NewError(ErrUserUnauthorized, fmt.Errorf("unauthorized client [%s]", clientID))
	}
	// 根据客户端 ID 获得对应的提示助手名称
	assistantName := getAssistantNameByClientId(clientID)
	record := &mysqldb.Record{
		ClientID: clientID,
		C0:       float64(req.GetCInfo().GetC0()),
		C1:       float64(req.GetCInfo().GetC1()),
		C2:       float64(req.GetCInfo().GetC2()),
		C3:       float64(req.GetCInfo().GetC3()),
		C4:       float64(req.GetCInfo().GetC4()),
		C5:       float64(req.GetCInfo().GetC5()),
		C6:       float64(req.GetCInfo().GetC6()),
		C7:       float64(req.GetCInfo().GetC7()),
		// 这里没有姿势矫正的概念，所以传入一个非站姿的姿势
		MeasurementPosture: mysqldb.MeasurementPostureSetting,
		HeartRate:          0,
	}

	// 问答传空
	var answers map[string]AnalysisReportAnswers
	ctxData, err := buildAEInput(userProfile, record, answers, req.GetCid(), true)
	if err != nil {
		return NewError(ErrLoadConfig, fmt.Errorf("failed to build ae input of record [%ds]: %s", record.RecordID, err.Error()))
	}
	language, err := mapProtoLanguageToAE(req.GetLanguage())
	if err != nil {
		return NewError(ErrInvalidLanguage, fmt.Errorf("invalid language %s: %s", req.GetLanguage(), err.Error()))
	}
	err = j.biz.LoadPresets(j.presetsFilePath, biz.MeasurementJudgmentModule)
	if err != nil {
		return NewError(ErrLoadConfig, fmt.Errorf("failed to load presets: %s", err.Error()))
	}
	preset, err := j.biz.SelectPresetByClientID(clientID)
	if err != nil {
		return NewError(ErrClientID, fmt.Errorf("failed to get preset by client_id %s: %s", clientID, err.Error()))
	}
	// 构建分析需要的上下文信息
	presetAgent := biz.PresetAgent{
		// 助手名称
		AssistantName: assistantName,
		// 用户名称
		ReporterName: userProfile.GetNickname(),
	}
	errRunEngine := preset.RunEngine(ctxData, language, presetAgent)
	if errRunEngine != nil {
		// 分析有错
		return NewError(ErrAEError, fmt.Errorf("error occurs when run engine: %s", errRunEngine.Error()))
	}

	newCtxData, _ := buildAEInput(userProfile, record, answers, req.GetCid(), false)
	dirtyDialecticLookups := parseLookupsIncludeKeyContent(ctxData.Output["dirty_dialectic"])
	disease, _ := getDiseasesMessageFromInput(ctxData.Output)
	physicalDialectics := req.GetPhysicalDialectics()
	// 补全上下文信息的体质辩证，脏腑辩证和疾病
	newCtxData.Input["physical_dialectics"] = convertPhysicalDialecticsToAEInput(physicalDialectics)
	newCtxData.Input["dirty_dialectic"] = convertAEInputKey(dirtyDialecticLookups)
	newCtxData.Input["disease"] = convertAEInputKey(disease)

	// 加载周报的 preset
	err = j.biz.LoadPresets(j.presetsFilePath, biz.WeeklyReportModule)
	if err != nil {
		return NewError(ErrClientID, fmt.Errorf("failed to get preset by client_id %s: %s", clientID, err.Error()))
	}
	newPreset, err := j.biz.SelectPresetByClientID(clientID)
	if err != nil {
		return NewError(ErrClientID, fmt.Errorf("failed to get preset by client_id %s: %s", clientID, err.Error()))
	}
	errRunEngine = newPreset.RunEngine(newCtxData, language, presetAgent)
	if errRunEngine != nil {
		// 分析有错
		return NewError(ErrAEError, fmt.Errorf("error occurs when run weekly engine: %s", errRunEngine.Error()))
	}
	// 构建返回的内容
	hasAbnormalMeasurementFromAeOutput := checkIfHasAbnormalMeasurementFromAeOutput(*newCtxData)
	displayOptions := getDisplayOption(record.ClientID, userProfile.GetGender(), checkIfHasStressStateFromAeOutput(*newCtxData), hasAbnormalMeasurementFromAeOutput, false)

	respModules, err := getModulesOfWeeklyReport(*newCtxData, displayOptions, userProfile.GetGender())
	if err != nil {
		return NewError(ErrBuildReturnModule, fmt.Errorf("failed to get modules from content output"))
	}
	// 下面构建报告内容里与引擎输出无关的模块
	restModules, err := j.buildModulesNotAboutAE(record, displayOptions, false)
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

	resp.ReportVersion = DefaultReportVersion
	resp.Report = &analysispb.ReportContent{
		UserProfile: userProfileModule,
		Modules:     respModules,
	}
	return nil
}
