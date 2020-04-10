package rest

import (
	"encoding/json"
	"time"

	"github.com/jinmukeji/go-pkg/crypto/rand"

	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/ae-v1/core"
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	corepb "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	ptypespb "github.com/jinmukeji/proto/gen/micro/idl/ptypes/v1"
	"github.com/kataras/iris/v12"
)

// GetAnalyzeResult 得到智能分析结果
func (h *v2Handler) GetAnalyzeResult(ctx iris.Context) {
	var analysisReportBody AnalysisReportBody
	errReadJSON := ctx.ReadJSON(&analysisReportBody)
	if errReadJSON != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", errReadJSON), false)
		return
	}

	recordID, err := ctx.Params().GetInt("record_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	if analysisReportBody.Cid != int64(recordID) {
		writeError(ctx, wrapError(ErrInvalidValue, fmt.Sprintf("cid %d is invalid", analysisReportBody.Cid), err), false)
		return
	}
	tags, errGetAnalysisSystemTags := h.getAnalysisSystemTags(ctx)
	if errGetAnalysisSystemTags != nil {
		writeRpcInternalError(ctx, errGetAnalysisSystemTags, false)
		return
	}
	transactionNumber, errGetAnalysisReportTransactionNumber := h.getAnalysisReportTransactionNumber(ctx)
	if errGetAnalysisReportTransactionNumber != nil {
		writeRpcInternalError(ctx, errGetAnalysisReportTransactionNumber, false)
		return
	}
	req := new(corepb.GetAnalyzeResultRequest)
	req.RecordId = int32(recordID)
	req.UserId = int32(analysisReportBody.UserID)
	req.SystemTags = tags
	answerTags := make([]string, 0)
	for _, answer := range analysisReportBody.Answers {
		for _, value := range answer.Values {
			answerTags = append(answerTags, core.BuildFullQuestionChoiceTag(answer.QuestionKey, value))
		}
	}
	req.AnswerTags = answerTags
	req.TransactionNumber = transactionNumber
	req.Answers = mapAnswer(analysisReportBody.Answers)
	resp, errGetAnalyzeResult := h.rpcSvc.GetAnalyzeResult(
		newRPCContext(ctx), req,
	)
	if errGetAnalyzeResult != nil {
		writeRpcInternalError(ctx, errGetAnalyzeResult, false)
		return
	}
	retData := buildAnalysisReportData()
	retData.Cid = (int64)(resp.Cid)
	retData.AnalysisReport.ReportID = resp.AnalysisReport.ReportId
	retData.AnalysisReport.ReportVersion = resp.AnalysisReport.ReportVersion
	retData.TransactionNo = resp.TransactionNo
	questionnaire := resp.Questionnaire
	retData.Questionnaire.Title = questionnaire.Title
	retData.Questionnaire.Questions = make([]Question, len(questionnaire.Questions))
	retData.Questionnaire.CreateAt = time.Now().UTC()
	retData.Questionnaire.Answers = analysisReportBody.Answers
	retData.AnalysisSession, _ = rand.RandomStringWithMask(rand.MaskLetterDigits, 60)
	for idx, question := range questionnaire.Questions {
		choices := make([]Choice, len(question.Choices))
		for idx, questionChoice := range question.Choices {
			choices[idx] = Choice{
				Key:          questionChoice.Key,
				Name:         questionChoice.Name,
				Value:        questionChoice.Key,
				ConflictKeys: questionChoice.ConflictKeys,
			}
		}
		retData.Questionnaire.Questions[idx] = Question{
			Key:         question.Key,
			Title:       question.Title,
			Description: question.Description,
			Tip:         question.Tip,
			Type:        question.Type,
			Choices:     choices,
			DefaultKeys: question.DefaultKeys,
		}
	}
	retData.AnalysisDone = resp.AnalysisDone
	retData.AnalysisReport.Content, _ = h.getAnalysisReportContent(ctx, resp.AnalysisReport)
	retData.AnalysisResultError = resp.AnalysisResultError
	// 分析完成，保存Answers
	if retData.AnalysisDone {
		errSaveAnswers := h.saveAnswers(ctx, analysisReportBody.Answers, int32(recordID))
		if errSaveAnswers != nil {
			writeRpcInternalError(ctx, errSaveAnswers, false)
			return
		}
	}
	rest.WriteOkJSON(ctx, retData)
}

func mapCCStrengthItems(protoCCStrengthItems []*corepb.CCStrengthItem) []CCStrengthItem {
	ccStrengthItems := make([]CCStrengthItem, len(protoCCStrengthItems))
	for idx, protoCCStrengthItem := range protoCCStrengthItems {
		labels := protoCCStrengthItem.Labels
		ccStrengthLabel := make([]CCStrengthLabel, len(labels))
		for idx, protoLabels := range labels {
			ccStrengthLabel[idx] = CCStrengthLabel{
				Label: protoLabels.Label,
				CC:    protoLabels.Cc,
			}
		}
		ccStrengthItems[idx] = CCStrengthItem{
			Key:      protoCCStrengthItem.Key,
			Disabled: protoCCStrengthItem.Disabled,
			Remark:   protoCCStrengthItem.Remark,
			Labels:   ccStrengthLabel,
		}
	}
	return ccStrengthItems
}

// MapGeneralExplains corepb.GeneralExplain数组转GeneralExplain数组
func MapGeneralExplains(protoGeneralExplains []*corepb.GeneralExplain) []GeneralExplain {
	generalExplains := make([]GeneralExplain, len(protoGeneralExplains))
	for idx, protoGeneralExplain := range protoGeneralExplains {
		generalExplains[idx] = GeneralExplain{
			Key:     protoGeneralExplain.Key,
			Label:   protoGeneralExplain.Label,
			Content: protoGeneralExplain.Content,
		}
	}
	return generalExplains
}

// GetAnalyzeReport 拿到分析报告
func (h *v2Handler) GetAnalyzeReport(ctx iris.Context) {
	recordID, err := ctx.Params().GetInt("record_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	tags, errGetAnalysisSystemTags := h.getAnalysisSystemTags(ctx)
	if errGetAnalysisSystemTags != nil {
		writeError(ctx, wrapError(ErrRPCInternal, "", errGetAnalysisSystemTags), false)
		return
	}
	req := new(corepb.GetAnalyzeResultRequest)
	req.RecordId = int32(recordID)
	req.SystemTags = tags
	recordAnswers, _ := h.getAnswers(ctx, recordID)
	answerTags := make([]string, 0)
	if recordAnswers != "" {
		var answers []Answer
		errUnmarshal := json.Unmarshal([]byte(recordAnswers), &answers)
		if errUnmarshal != nil {
			writeError(ctx, wrapError(ErrRPCInternal, "", fmt.Errorf("failed to unmarshal %s: %s", recordAnswers, errUnmarshal.Error())), false)
			return
		}
		for _, answer := range answers {
			for _, value := range answer.Values {
				answerTags = append(answerTags, core.BuildFullQuestionChoiceTag(answer.QuestionKey, value))
			}
		}
	}
	req.AnswerTags = answerTags
	transactionNumber, errGetAnalysisReportTransactionNumber := h.getAnalysisReportTransactionNumber(ctx)
	if errGetAnalysisReportTransactionNumber != nil {
		writeError(ctx, wrapError(ErrRPCInternal, "", errGetAnalysisReportTransactionNumber), false)
		return
	}
	req.TransactionNumber = transactionNumber
	resp, errGetAnalyzeResult := h.rpcSvc.GetAnalyzeResult(
		newRPCContext(ctx), req,
	)
	if errGetAnalyzeResult != nil {
		writeError(ctx, wrapError(ErrRPCInternal, "", errGetAnalyzeResult), false)
		return
	}
	analysisReport := AnalysisReport{}
	analysisReport.ReportID = resp.TransactionNo
	analysisReport.ReportVersion = resp.AnalysisReport.ReportVersion
	analysisReport.Content, _ = h.getAnalysisReportContent(ctx, resp.AnalysisReport)
	rest.WriteOkJSON(ctx, analysisReport)
}

// getAnalysisSystemTags 得到分析的tags
func (h *v2Handler) getAnalysisSystemTags(ctx iris.Context) ([]string, error) {
	req := new(corepb.GetAnalysisSystemTagsRequest)
	resp, errGetAnalysisSystemTags := h.rpcSvc.GetAnalysisSystemTags(
		newRPCContext(ctx), req,
	)
	if errGetAnalysisSystemTags != nil {
		return nil, errGetAnalysisSystemTags
	}
	return resp.Tags, nil
}

// getAnalysisReportContent 得到分析报告中的content
func (h *v2Handler) getAnalysisReportContent(ctx iris.Context, analysisReport *corepb.AnalysisReport) (Content, error) {
	birthday, _ := ptypes.Timestamp(analysisReport.Content.UserProfile.BirthdayTime)
	createdAt, _ := ptypes.Timestamp(analysisReport.Content.CreatedTime)
	gender, errMapProtoGenderToRest := mapProtoGenderToRest(analysisReport.Content.UserProfile.Gender)
	if errMapProtoGenderToRest != nil {
		return Content{}, errMapProtoGenderToRest
	}
	content := Content{
		Lead:                                      MapGeneralExplains(analysisReport.Content.Lead),
		TipsForWoman:                              MapGeneralExplains(analysisReport.Content.TipsForWoman),
		ChannelsAndCollateralsExplains:            MapGeneralExplains(analysisReport.Content.ChannelsAndCollateralsExplains),
		ConstitutionDifferentiationExplains:       MapGeneralExplains(analysisReport.Content.ConstitutionDifferentiationExplains),
		SyndromeDifferentiationExplains:           MapGeneralExplains(analysisReport.Content.SyndromeDifferentiationExplains),
		ChannelsAndCollateralsStrength:            mapCCStrengthItems(analysisReport.Content.ChannelsAndCollateralsStrength),
		BabyTips:                                  MapGeneralExplains(analysisReport.Content.BabyTips),
		ConstitutionDifferentiationExplainNotices: MapGeneralExplains(analysisReport.Content.ConstitutionDifferentiationExplainNotices),
		Tags:                                  MapGeneralExplains(analysisReport.Content.Tags),
		DictionaryEntries:                     MapGeneralExplains(analysisReport.Content.DictionaryEntries),
		FactorExplains:                        MapGeneralExplains(analysisReport.Content.FactorExplains),
		HealthDescriptions:                    MapGeneralExplains(analysisReport.Content.HealthDescriptions),
		MeasurementTips:                       MapGeneralExplains(analysisReport.Content.MeasurementTips),
		CCExplainNotices:                      MapGeneralExplains(analysisReport.Content.ChannelsAndCollateralsExplainNotices),
		UterineHealthIndexes:                  MapGeneralExplains(analysisReport.Content.UterineHealthIndexes),
		UterusAttentionPrompts:                MapGeneralExplains(analysisReport.Content.UterusAttentionPrompts),
		UterineHealthDescriptions:             MapGeneralExplains(analysisReport.Content.UterineHealthDescriptions),
		MenstrualHealthValues:                 MapGeneralExplains(analysisReport.Content.MenstrualHealthValues),
		MenstrualHealthDescriptions:           MapGeneralExplains(analysisReport.Content.MenstrualHealthDescriptions),
		GynecologicalInflammations:            MapGeneralExplains(analysisReport.Content.GynecologicalInflammations),
		GynecologicalInflammationDescriptions: MapGeneralExplains(analysisReport.Content.GynecologicalInflammationDescriptions),
		BreastHealth:                          MapGeneralExplains(analysisReport.Content.BreastHealth),
		BreastHealthDescriptions:              MapGeneralExplains(analysisReport.Content.BreastHealthDescriptions),
		EmotionalHealthIndexes:                MapGeneralExplains(analysisReport.Content.EmotionalHealthIndexes),
		EmotionalHealthDescriptions:           MapGeneralExplains(analysisReport.Content.EmotionalHealthDescriptions),
		FacialSkins:                           MapGeneralExplains(analysisReport.Content.FacialSkins),
		FacialSkinDescriptions:                MapGeneralExplains(analysisReport.Content.FacialSkinDescriptions),
		ReproductiveAgeConsiderations:         MapGeneralExplains(analysisReport.Content.ReproductiveAgeConsiderations),
		BreastCancerOvarianCancers:            MapGeneralExplains(analysisReport.Content.BreastCancerOvarianCancers),
		BreastCancerOvarianCancerDescriptions: MapGeneralExplains(analysisReport.Content.BreastCancerOvarianCancerDescriptions),
		HormoneLevels:                         MapGeneralExplains(analysisReport.Content.HormoneLevels),
		LymphaticHealth:                       MapGeneralExplains(analysisReport.Content.LymphaticHealth),
		LymphaticHealthDescriptions:           MapGeneralExplains(analysisReport.Content.LymphaticHealthDescriptions),
		F100:                                  analysisReport.Content.F100,
		F101:                                  analysisReport.Content.F101,
		F102:                                  analysisReport.Content.F102,
		F103:                                  analysisReport.Content.F103,
		F104:                                  analysisReport.Content.F104,
		F105:                                  analysisReport.Content.F105,
		F106:                                  analysisReport.Content.F106,
		F107:                                  analysisReport.Content.F107,
		UserProfile: ReportUserProfile{
			UserID:    int64(analysisReport.Content.UserProfile.UserId),
			Nickname:  analysisReport.Content.UserProfile.Nickname,
			Birthday:  birthday,
			Age:       int64(analysisReport.Content.UserProfile.Age),
			Gender:    int64(gender),
			Height:    int64(analysisReport.Content.UserProfile.Height),
			Weight:    int64(analysisReport.Content.UserProfile.Weight),
			AvatarURL: analysisReport.Content.UserProfile.AvatarUrl,
		},
		MeasurementResult: MeasurementResult{
			Finger:              int64(analysisReport.Content.MeasurementResult.Finger),
			C0:                  int64(analysisReport.Content.MeasurementResult.C0),
			C1:                  int64(analysisReport.Content.MeasurementResult.C1),
			C2:                  int64(analysisReport.Content.MeasurementResult.C2),
			C3:                  int64(analysisReport.Content.MeasurementResult.C3),
			C4:                  int64(analysisReport.Content.MeasurementResult.C4),
			C5:                  int64(analysisReport.Content.MeasurementResult.C5),
			C6:                  int64(analysisReport.Content.MeasurementResult.C6),
			C7:                  int64(analysisReport.Content.MeasurementResult.C7),
			HeartRate:           int64(analysisReport.Content.MeasurementResult.Hr),
			AppHeartRate:        int64(analysisReport.Content.MeasurementResult.AppHr),
			PartialPulseWave:    analysisReport.Content.MeasurementResult.PartialInfo,
			AppHighestHeartRate: analysisReport.Content.MeasurementResult.AppHighestHr,
			AppLowestHeartRate:  analysisReport.Content.MeasurementResult.AppLowestHr,
		},
		CreatedAt: createdAt,
	}
	content.M0 = ParseNullableInt32Value(analysisReport.Content.M0)
	content.M1 = ParseNullableInt32Value(analysisReport.Content.M1)
	content.M2 = ParseNullableInt32Value(analysisReport.Content.M2)
	content.M3 = ParseNullableInt32Value(analysisReport.Content.M3)
	return content, nil
}

// ParseNullableInt32Value 解析可空的int32值
func ParseNullableInt32Value(value *ptypespb.NullableInt32Value) *int32 {
	switch value.GetKind().(type) {
	case *ptypespb.NullableInt32Value_NullValue:
		return nil
	case *ptypespb.NullableInt32Value_Int32Value:
		i := value.GetInt32Value()
		return &i
	default:
		return nil
	}
}

// saveAnswers 保存答案
func (h *v2Handler) saveAnswers(ctx iris.Context, answers []Answer, recordID int32) error {
	reqSubmitRecordAnswersRequest := new(corepb.SubmitRecordAnswersRequest)
	json, errMarshal := json.Marshal(answers)
	if errMarshal != nil {
		return errMarshal
	}
	reqSubmitRecordAnswersRequest.RecordId = recordID
	reqSubmitRecordAnswersRequest.Answers = string(json)
	_, errSubmitRecordAnswers := h.rpcSvc.SubmitRecordAnswers(
		newRPCContext(ctx), reqSubmitRecordAnswersRequest,
	)
	if errSubmitRecordAnswers != nil {
		return errSubmitRecordAnswers
	}
	return nil
}

// buildAnalysisReportData 构造 AnalysisReportData
func buildAnalysisReportData() *AnalysisReportData {
	return &AnalysisReportData{
		Cid:             0,
		AnalysisSession: "",
		AnalysisDone:    false,
		Questionnaire: Questionnaire{
			Title:     "",
			Questions: []Question{},
			Answers:   []Answer{},
			CreateAt:  time.Now().UTC(),
		},
		AnalysisReport: &AnalysisReport{},
	}
}

// getAnalysisReportTransactionNumber 得到流水号
func (h *v2Handler) getAnalysisReportTransactionNumber(ctx iris.Context) (int32, error) {
	reqGetAnalysisReportTransactionNumberRequest := new(corepb.GetAnalysisReportTransactionNumberRequest)
	respGetAnalysisReportTransactionNumber, errGetAnalysisReportTransactionNumber := h.rpcSvc.GetAnalysisReportTransactionNumber(
		newRPCContext(ctx), reqGetAnalysisReportTransactionNumberRequest,
	)
	if errGetAnalysisReportTransactionNumber != nil {
		return 0, errGetAnalysisReportTransactionNumber
	}
	return respGetAnalysisReportTransactionNumber.TransactionNumber, nil
}

// getAnswers 得到Answers
func (h *v2Handler) getAnswers(ctx iris.Context, recordID int) (string, error) {
	req := new(corepb.GetMeasurementRecordRequest)
	req.RecordId = int32(recordID)
	respGetMeasurementRecord, errGetMeasurementRecord := h.rpcSvc.GetMeasurementRecord(
		newRPCContext(ctx), req,
	)
	if errGetMeasurementRecord != nil {
		return "", errGetMeasurementRecord
	}
	return respGetMeasurementRecord.Answers, nil
}

func mapAnswer(inputAnswers []Answer) []*corepb.Answer {
	answers := make([]*corepb.Answer, 0)
	for _, answer := range inputAnswers {
		values := make([]string, 0)
		values = append(values, answer.Values...)
		answers = append(answers, &corepb.Answer{
			QuestionKey: answer.QuestionKey,
			Values:      values,
		})
	}
	return answers
}
