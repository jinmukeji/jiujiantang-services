package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"

	"github.com/jinmukeji/ae-v1/core"
	"github.com/jinmukeji/gf-api2/pkg/rpc"
	"github.com/jinmukeji/gf-api2/service/auth"
	mysql "github.com/jinmukeji/gf-api2/service/mysqldb"
	jinmuidpb "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	corepb "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	generalpb "github.com/jinmukeji/proto/gen/micro/idl/ptypes/v2"
)

const (
	// FemaleFactionFlag 女性健康标记
	FemaleFactionFlag = "HasFemaleAnalysisExplains"
	// RecordType1_8 1.8的record_type是7
	RecordType1_8 = 7
)

// CreateReportShareToken 创建分享报告的token
func (j *JinmuHealth) CreateReportShareToken(ctx context.Context, req *corepb.CreateReportShareTokenRequest, resp *corepb.CreateReportShareTokenResponse) error {
	if exist, err := j.datastore.ExistRecordByRecordID(ctx, req.RecordId); !exist || err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to check if record %d exists: %s", req.RecordId, err.Error()))
	}
	record, errFindRecordByID := j.datastore.FindRecordByID(ctx, int(req.RecordId))
	if errFindRecordByID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find record by recordID %d: %s", req.RecordId, errFindRecordByID.Error()))
	}
	accessTokenType, _ := auth.AccessTokenTypeFromContext(ctx)
	if accessTokenType != AccessTokenTypeLValue && accessTokenType != AccessTokenTypeWeChatValue {
		userID, _ := auth.UserIDFromContext(ctx)
		o, err := j.datastore.FindOrganizationByUserID(ctx, int(userID))
		if err != nil {
			return fmt.Errorf("failed to find organization by userID %d: %s", int(userID), err.Error())
		}
		organization, err := j.datastore.FindOrganizationByUserID(ctx, record.UserID)
		if err != nil {
			return fmt.Errorf("failed to find organization by userID %d in record: %s", int(userID), err.Error())
		}
		if o.OrganizationID != organization.OrganizationID {
			return NewError(ErrNoPermissionGetShareToken, errors.New("No permission to get share token"))
		}
	}

	resp.Token = record.RecordToken
	if record.RecordToken == "" {
		token := uuid.New().String()
		err := j.datastore.UpdateRecordToken(ctx, req.RecordId, token)
		if err != nil {
			return fmt.Errorf("failed to update record token of recordID %d: %s", req.RecordId, err.Error())
		}
		resp.Token = token
	}
	resp.Link = fmt.Sprintf("%s#/analysisreport?token=%s", j.wechat.Options.JinmuH5ServerbaseV2_1, resp.Token)
	return nil
}

// GetAnalyzeResultByToken 通过分享的token拿到分析数据
func (j *JinmuHealth) GetAnalyzeResultByToken(ctx context.Context, req *corepb.GetAnalyzeResultByTokenRequest, resp *corepb.GetAnalyzeResultByTokenResponse) error {
	record, err := j.datastore.FindRecordByToken(ctx, req.Token)
	if record == nil || err != nil {
		return fmt.Errorf("failed to get record by token: %s", err.Error())
	}
	measurementResult, errbuildMeasurementResult := j.buildAnalysisResult(ctx, record, int(record.UserID))
	if errbuildMeasurementResult != nil {
		return NewError(ErrGetMeasureResultFailure, errbuildMeasurementResult)
	}
	reportUserProfile, errbuildUserProfile := j.buildAnalysisUserProfile(ctx, int32(record.RecordID))
	if errbuildUserProfile != nil {
		return NewError(ErrGetUserFailure, errbuildUserProfile)
	}
	//  处理站姿C0-C7
	if corepb.MeasurementPosture(record.MeasurementPosture) == corepb.MeasurementPosture_MEASUREMENT_POSTURE_STANDING {
		DealWithMeasurementPosture(measurementResult.ReportMeasurementResult, reportUserProfile.Gender)
	}
	tags := make([]string, 0)
	client, _ := clientFromContext(ctx)
	if client.ClientID == dengyunClient {
		tags = append(tags, core.EnabledCCSystemTag)
	} else {
		tags = append(tags, core.EnabledCDSystemTag)
		tags = append(tags, core.EnabledCCSystemTag)
		tags = append(tags, core.EnabledSDSystemTag)
		tags = append(tags, core.EnabledFactorSystemTag)
	}
	accessTokenType, _ := auth.AccessTokenTypeFromContext(ctx)
	// 女性部分内容是根据一体机客户进行开关的
	if accessTokenType == AccessTokenTypeLValue || accessTokenType == AccessTokenTypeWeChatValue {
		tags = append(tags, core.EnabledBreastSystemTag)
		tags = append(tags, core.EnabledEmotionalSystemTag)
		tags = append(tags, core.EnabledFacialSystemTag)
		tags = append(tags, core.EnabledGynecologicalSystemTag)
		tags = append(tags, core.EnabledHormoneSystemTag)
		tags = append(tags, core.EnabledLymphaticSystemTag)
		tags = append(tags, core.EnabledMammarySystemTag)
		tags = append(tags, core.EnabledMenstruationSystemTag)
		tags = append(tags, core.EnabledReproductiveSystemTag)
		tags = append(tags, core.EnabledUterineSystemTag)
	}
	recordAnswers := record.Answers
	answerTags := make([]string, 0)
	if recordAnswers != "" {
		var answers []Answer
		errUnmarshal := json.Unmarshal([]byte(recordAnswers), &answers)
		if errUnmarshal != nil {
			return NewError(ErrJSONUnmarshalFailure, errUnmarshal)
		}
		for _, answer := range answers {
			for _, value := range answer.Values {
				answerTags = append(answerTags, core.BuildFullQuestionChoiceTag(answer.QuestionKey, value))
			}
		}
	}
	out, errRunAnalysisEngine := j.RunAnalysisEngine(ctx, measurementResult, reportUserProfile, tags, answerTags)
	if errRunAnalysisEngine != nil {
		return NewError(ErrRunAnalysisEngineFailure, errRunAnalysisEngine)
	}
	questions := BuildOutQuestions(out)
	createdAt, _ := ptypes.TimestampProto(measurementResult.CreatedAt)
	resp.Cid = int32(record.RecordID)
	resp.Questionnaire = &corepb.Questionnaire{
		Title:       "",
		CreatedTime: createdAt,
		Questions:   questions,
	}
	content := BuildProtoContent(ctx, out, measurementResult, reportUserProfile)
	transactionNo := record.TransactionNumber
	resp.TransactionNo = transactionNo
	resp.AnalysisReport = j.BuildAnalysisReportOptions(ctx, transactionNo, content, reportUserProfile, out, record.RecordType)
	resp.AnalysisDone = true
	return nil
}

// DealWithMeasurementPosture 处理站姿C0-C7
func DealWithMeasurementPosture(result *corepb.ReportMeasurementResult, gender generalpb.Gender) {
	// TODO 以后重构
	result.C0, result.C1, result.C2, result.C3, result.C4, result.C5, result.C6, result.C7 = ConvertToStandingCInt32Values(gender, int(result.C0), int(result.C1), int(result.C2), int(result.C3), int(result.C4), int(result.C5), int(result.C6), int(result.C7))
}

// BuildAnalysisReportOptions 构建分析报告Options
func (j *JinmuHealth) BuildAnalysisReportOptions(ctx context.Context, reportID string, content *corepb.Content, reportUserProfile *corepb.ReportUserProfile, out *core.Output, recordType int) *corepb.AnalysisReport {
	analysisReport := &corepb.AnalysisReport{
		ReportVersion: DefaultAnalysisReportVersion,
		ReportId:      reportID,
		Content:       content,
	}
	reqGetUserProfile := new(jinmuidpb.GetUserProfileRequest)
	reqGetUserProfile.IsSkipVerifyToken = true
	reqGetUserProfile.UserId = reportUserProfile.UserId
	respGetUserProfile, _ := j.jinmuidSvc.GetUserProfile(ctx, reqGetUserProfile)
	if respGetUserProfile.CustomizedCode == customSanshui {
		analysisReport.Options = &corepb.DisplayOptions{
			DisplayNavbar:                 false,
			DisplayTags:                   true,
			DisplayPartialInfo:            true,
			DisplayUserProfile:            true,
			DisplayHeartRate:              true,
			DisplayCcBarChart:             true,
			DisplayCcExplain:              true,
			DisplayCdExplain:              false,
			DisplaySdExplain:              false,
			DisplayF0:                     true,
			DisplayF1:                     true,
			DisplayF2:                     true,
			DisplayF3:                     true,
			DisplayPhysicalTherapyExplain: false,
			DisplayRemark:                 true,
			DisplayMeasurementResult:      true,
			DisplayBabyTips:               true,
			DisplayWh:                     false,
		}
	} else {
		analysisReport.Options = &corepb.DisplayOptions{
			DisplayNavbar:                 true,
			DisplayTags:                   true,
			DisplayPartialInfo:            true,
			DisplayUserProfile:            true,
			DisplayHeartRate:              true,
			DisplayCcBarChart:             true,
			DisplayCcExplain:              true,
			DisplayCdExplain:              true,
			DisplaySdExplain:              true,
			DisplayF0:                     true,
			DisplayF1:                     true,
			DisplayF2:                     true,
			DisplayF3:                     true,
			DisplayPhysicalTherapyExplain: true,
			DisplayRemark:                 true,
			DisplayMeasurementResult:      true,
			DisplayBabyTips:               true,
			DisplayWh:                     out.NamedFlags[FemaleFactionFlag],
		}
	}
	// recordType 版本小于1.8的，没有心率不显示
	if recordType < RecordType1_8 {
		analysisReport.Options.DisplayHeartRate = false
	}
	return analysisReport
}

// buildAnalysisResult 构建AnalysisResult
func (j *JinmuHealth) buildAnalysisResult(ctx context.Context, record *mysql.Record, reqUserID int) (*MeasurementResult, error) {
	l := rpc.ContextLogger(ctx)
	var partialData []int32
	dataArray, errgetPulseTestDataIntArray := getPulseTestDataIntArray(record.S3Key, j.awsClient)
	if errgetPulseTestDataIntArray != nil {
		l.WithError(errgetPulseTestDataIntArray).Warn("failed to getPulseTestDataIntArray of S3Key")
	}
	if len(dataArray) >= 4000 {
		partialData = dataArray[3000:4000]
	}
	protoFinger, errMapDBFingerToProto := mapDBFingerToProto(record.Finger)
	if errMapDBFingerToProto != nil {
		return nil, errMapDBFingerToProto
	}
	return &MeasurementResult{
		ReportMeasurementResult: &corepb.ReportMeasurementResult{
			C0:           Int32ValBoundedBy10FromFloat(record.C0),
			C1:           Int32ValBoundedBy10FromFloat(record.C1),
			C2:           Int32ValBoundedBy10FromFloat(record.C2),
			C3:           Int32ValBoundedBy10FromFloat(record.C3),
			C4:           Int32ValBoundedBy10FromFloat(record.C4),
			C5:           Int32ValBoundedBy10FromFloat(record.C5),
			C6:           Int32ValBoundedBy10FromFloat(record.C6),
			C7:           Int32ValBoundedBy10FromFloat(record.C7),
			AppHr:        int32(record.AppHeartRate),
			Hr:           int32(record.HeartRate),
			PartialInfo:  partialData,
			Finger:       protoFinger,
			AppHighestHr: record.AppHighestHeartRate,
			AppLowestHr:  record.AppLowestHeartRate,
		},
		CreatedAt:          record.CreatedAt,
		HasPaid:            record.HasPaid,
		Remark:             record.Remark,
		MeasurementPosture: corepb.MeasurementPosture(record.MeasurementPosture),
	}, nil
}

// mapDBFingerToProto 将数据库里面存放的 finger 映射为 proto 格式
func mapDBFingerToProto(dbFinger mysql.Finger) (corepb.Finger, error) {
	switch dbFinger {
	case mysql.FingerLeft1:
		return corepb.Finger_FINGER_LEFT_1, nil
	case mysql.FingerLeft2:
		return corepb.Finger_FINGER_LEFT_2, nil
	case mysql.FingerLeft3:
		return corepb.Finger_FINGER_LEFT_3, nil
	case mysql.FingerLeft4:
		return corepb.Finger_FINGER_LEFT_4, nil
	case mysql.FingerLeft5:
		return corepb.Finger_FINGER_LEFT_5, nil
	case mysql.FingerRight1:
		return corepb.Finger_FINGER_RIGHT_1, nil
	case mysql.FingerRight2:
		return corepb.Finger_FINGER_RIGHT_2, nil
	case mysql.FingerRight3:
		return corepb.Finger_FINGER_RIGHT_3, nil
	case mysql.FingerRight4:
		return corepb.Finger_FINGER_RIGHT_4, nil
	case mysql.FingerRight5:
		return corepb.Finger_FINGER_RIGHT_5, nil
	}
	return corepb.Finger_FINGER_INVALID, fmt.Errorf("invalid database finger %d", dbFinger)
}

// mapProtoFingerToDB 将proto 格式 的 finger 映射为 数据库里面存放的格式
func mapProtoFingerToDB(protoFinger corepb.Finger) (mysql.Finger, error) {
	switch protoFinger {
	case corepb.Finger_FINGER_INVALID:
		return mysql.FingerInvalid, fmt.Errorf("invalid proto finger %d", protoFinger)
	case corepb.Finger_FINGER_UNSET:
		return mysql.FingerInvalid, fmt.Errorf("invalid proto finger %d", protoFinger)
	case corepb.Finger_FINGER_LEFT_1:
		return mysql.FingerLeft1, nil
	case corepb.Finger_FINGER_LEFT_2:
		return mysql.FingerLeft2, nil
	case corepb.Finger_FINGER_LEFT_3:
		return mysql.FingerLeft3, nil
	case corepb.Finger_FINGER_LEFT_4:
		return mysql.FingerLeft4, nil
	case corepb.Finger_FINGER_LEFT_5:
		return mysql.FingerLeft5, nil
	case corepb.Finger_FINGER_RIGHT_1:
		return mysql.FingerRight1, nil
	case corepb.Finger_FINGER_RIGHT_2:
		return mysql.FingerRight2, nil
	case corepb.Finger_FINGER_RIGHT_3:
		return mysql.FingerRight3, nil
	case corepb.Finger_FINGER_RIGHT_4:
		return mysql.FingerRight4, nil
	case corepb.Finger_FINGER_RIGHT_5:
		return mysql.FingerRight5, nil
	}
	return mysql.FingerInvalid, fmt.Errorf("invalid proto finger %d", protoFinger)
}
