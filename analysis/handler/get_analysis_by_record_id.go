package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/ae/v2/biz"
	"github.com/jinmukeji/jiujiantang-services/analysis/mysqldb"
	"github.com/jinmukeji/plat-pkg/rpc/errors"
	"github.com/jinmukeji/plat-pkg/rpc/errors/codes"
	analysispb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/analysis/v1"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/analysis/v1"
	jinmuidpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
)

// GetAnalyzeResultByRecordID 根据 record_id 得到分析结果
func (j *AnalysisManagerService) GetAnalyzeResultByRecordID(ctx context.Context, req *proto.GetAnalyzeResultByRecordIDRequest, resp *proto.GetAnalyzeResultByRecordIDResponse) error {
	record, err := j.database.FindRecordByRecordID(req.GetRecordId())
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find analysis body by recordID %d: %s", req.GetRecordId(), err.Error()))
	}
	if record == nil {
		return NewError(ErrInvalidAnalysisStatus, fmt.Errorf("analysis status of record %d is invalid", req.GetRecordId()))
	}
	// 如果分析状态不是完成，则无法获得分析结果
	if record.AnalyzeStatus != mysqldb.AnalysisStatusCompeleted {
		return NewError(ErrDatabase, fmt.Errorf("analyze status [%d] of record [%d] is wrong", record.AnalyzeStatus, req.GetRecordId()))
	}

	// 获得用户个人信息
	reqGetUserProfile := new(jinmuidpb.GetUserProfileRequest)
	// 如果请求的不是当前用户，则不需要 token 验证
	token, _ := TokenFromContext(ctx)
	ownerID, _ := j.database.FindUserIDByToken(token)
	if ownerID != record.UserID {
		reqGetUserProfile.IsSkipVerifyToken = true
	}
	reqGetUserProfile.UserId = record.UserID
	respGetUserProfile, err := j.jinmuidSvc.GetUserProfile(ctx, reqGetUserProfile)
	if err != nil {
		return NewError(ErrInvalidUser, fmt.Errorf("failed to get user profile of user [%d]: %s", record.UserID, err.Error()))
	}
	userProfile := respGetUserProfile.GetProfile()
	clientID := record.ClientID

	// 存储分析 body
	analyzeBody := &AnalyzeBody{}
	err = json.Unmarshal([]byte(record.AnalyzeBody), analyzeBody)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to unmarshal analyze body of record [%d]: %s", record.RecordID, err.Error()))
	}
	ctxData, err := buildAEInput(userProfile, record, analyzeBody.QuestionAnswers, req.GetCid(), true)
	if err != nil {
		return NewError(ErrLoadConfig, fmt.Errorf("failed to build ae input of record [%d]: %s", record.RecordID, err.Error()))
	}

	language := analyzeBody.Language
	err = j.biz.LoadPresets(j.presetsFilePath, biz.MeasurementJudgmentModule)
	if err != nil {
		return NewError(ErrLoadConfig, fmt.Errorf("failed to load presets: %s", err.Error()))
	}

	preset, err := j.biz.SelectPresetByClientID(clientID)
	if err != nil {
		return NewError(ErrClientID, fmt.Errorf("failed to get preset by client_id %s: %s", clientID, err.Error()))
	}

	// 根据客户端 ID 获得对应的提示助手名称
	assistantName := getAssistantNameByClientId(clientID)

	// 构建分析需要的上下文信息
	presetAgent := biz.PresetAgent{
		// 助手名称
		AssistantName: assistantName,
		// 用户名称
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
		return errors.ErrorWithCause(codes.InvalidOperation, err, "failed to get modules from content output")
	}
	// 下面构建报告内容里与引擎输出无关的模块
	restModules, err := j.buildModulesNotAboutAE(record, displayOptions, true)
	if err != nil {
		return errors.ErrorWithCause(codes.InvalidOperation, err, "failed to get modules")
	}
	for key, value := range restModules {
		respModules[key] = value
	}
	userProfileModule, err := buildUserProfileModule(userProfile, record, displayOptions)
	if err != nil {
		return errors.ErrorWithCause(codes.InvalidArgument, err, "failed to build user profile module")
	}

	pulseTestModule, err := buildPulseTestModule(record, displayOptions)
	if err != nil {
		return errors.ErrorWithCause(codes.InvalidArgument, err, "failed to build pulse test module")
	}
	remarkModule := buildRemarkModule(record, displayOptions)

	protoAnalysisFinishTime, _ := ptypes.TimestampProto(record.CreatedAt)

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
