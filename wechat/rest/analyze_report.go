package rest

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/ae-v1/core"
	"github.com/jinmukeji/gf-api2/pkg/rest"
	corepb "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	ptypespb "github.com/jinmukeji/proto/gen/micro/idl/ptypes/v1"
	"github.com/kataras/iris/v12"
)

// GetAnalyzeReport 拿到分析报告
func (h *handler) GetAnalyzeReport(ctx iris.Context) {
	recordID, err := ctx.Params().GetInt("record_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}

	// Session
	session, err := h.getSession(ctx)
	if err != nil || !session.Authorized || session.IsExpired {
		path := fmt.Sprintf("%s/app.html#/analysisreport?record_id=%d", h.WxH5ServerBase, recordID)
		redirectURL := fmt.Sprintf("%s/wx/oauth?redirect=%s", h.WxCallbackServerBase, url.QueryEscape(path))
		writeSessionErrorJSON(ctx, redirectURL, err)
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
	req.UserId = int32(session.UserID)
	recordAnswers, _ := h.getAnswers(ctx, recordID)
	answerTags := make([]string, 0)
	if recordAnswers != "" {
		var answers []Answer
		errUnmarshal := json.Unmarshal([]byte(recordAnswers), &answers)
		if errUnmarshal != nil {
			writeError(ctx, wrapError(ErrRPCInternal, "", errUnmarshal), false)
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

// GetAnalyzeReportByToken 通过token拿到分析报告
func (h *handler) GetAnalyzeReportByToken(ctx iris.Context) {
	token := ctx.Params().Get("token")
	req := new(corepb.GetAnalyzeResultByTokenRequest)
	req.Token = token
	resp, errGetAnalyzeResultByToken := h.rpcSvc.GetAnalyzeResultByToken(
		newRPCContext(ctx), req,
	)
	if errGetAnalyzeResultByToken != nil {
		writeError(ctx, wrapError(ErrRPCInternal, "", errGetAnalyzeResultByToken), false)
		return
	}
	analysisReport := AnalysisReport{}
	analysisReport.ReportID = resp.TransactionNo
	analysisReport.ReportVersion = resp.AnalysisReport.ReportVersion
	analysisReport.Content, _ = h.getAnalysisReportContent(ctx, resp.AnalysisReport)
	rest.WriteOkJSON(ctx, analysisReport)
}

// getAnalysisSystemTags 得到分析的tags
func (h *handler) getAnalysisSystemTags(ctx iris.Context) ([]string, error) {
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
func (h *handler) getAnalysisReportContent(ctx iris.Context, analysisReport *corepb.AnalysisReport) (Content, error) {
	birthday, _ := ptypes.Timestamp(analysisReport.Content.UserProfile.BirthdayTime)
	createdAt, _ := ptypes.Timestamp(analysisReport.Content.CreatedTime)
	content := Content{
		Lead:                                      MapGeneralExplains(analysisReport.Content.Lead),
		TipsForWoman:                              MapGeneralExplains(analysisReport.Content.TipsForWoman),
		ChannelsAndCollateralsExplains:            MapGeneralExplains(analysisReport.Content.ChannelsAndCollateralsExplains),
		ConstitutionDifferentiationExplains:       MapGeneralExplains(analysisReport.Content.ConstitutionDifferentiationExplains),
		SyndromeDifferentiationExplains:           MapGeneralExplains(analysisReport.Content.SyndromeDifferentiationExplains),
		ChannelsAndCollateralsStrength:            mapCCStrengthItems(analysisReport.Content.ChannelsAndCollateralsStrength),
		BabyTips:                                  MapGeneralExplains(analysisReport.Content.BabyTips),
		ConstitutionDifferentiationExplainNotices: MapGeneralExplains(analysisReport.Content.ConstitutionDifferentiationExplainNotices),
		Tags:               MapGeneralExplains(analysisReport.Content.Tags),
		DictionaryEntries:  MapGeneralExplains(analysisReport.Content.DictionaryEntries),
		FactorExplains:     MapGeneralExplains(analysisReport.Content.FactorExplains),
		HealthDescriptions: MapGeneralExplains(analysisReport.Content.HealthDescriptions),
		MeasurementTips:    MapGeneralExplains(analysisReport.Content.MeasurementTips),
		CCExplainNotices:   MapGeneralExplains(analysisReport.Content.ChannelsAndCollateralsExplainNotices),

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
			Gender:    int64(analysisReport.Content.UserProfile.Gender),
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
		PTExplain: PhysicalTherapyExplain{
			F0: analysisReport.Content.F0,
			F1: analysisReport.Content.F1,
			F2: analysisReport.Content.F2,
			F3: analysisReport.Content.F3,
		},
		Options: DisplayOptions{
			DisplayNavbar:                 analysisReport.Options.DisplayNavbar,
			DisplayTags:                   analysisReport.Options.DisplayTags,
			DisplayPartialData:            analysisReport.Options.DisplayPartialInfo,
			DisplayUserProfile:            analysisReport.Options.DisplayUserProfile,
			DisplayHeartRate:              analysisReport.Options.DisplayHeartRate,
			DisplayCcBarChart:             analysisReport.Options.DisplayCcBarChart,
			DisplayCcExplain:              analysisReport.Options.DisplayCcExplain,
			DisplayCdExplain:              analysisReport.Options.DisplayCdExplain,
			DisplaySdExplain:              analysisReport.Options.DisplaySdExplain,
			DisplayF0:                     analysisReport.Options.DisplayF0,
			DisplayF1:                     analysisReport.Options.DisplayF1,
			DisplayF2:                     analysisReport.Options.DisplayF2,
			DisplayF3:                     analysisReport.Options.DisplayF3,
			DisplayPhysicalTherapyExplain: analysisReport.Options.DisplayPhysicalTherapyExplain,
			DisplayRemark:                 analysisReport.Options.DisplayRemark,
			DisplayMeasurementResult:      analysisReport.Options.DisplayMeasurementResult,
			DisplayBabyTips:               analysisReport.Options.DisplayBabyTips,
			DisplayWh:                     analysisReport.Options.DisplayWh,
		},
		Remark:    analysisReport.Content.Remark,
		HasPaid:   analysisReport.Content.HasPaid,
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

// getAnalysisReportTransactionNumber 得到流水号
func (h *handler) getAnalysisReportTransactionNumber(ctx iris.Context) (int32, error) {
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
func (h *handler) getAnswers(ctx iris.Context, recordID int) (string, error) {
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

// mapCCStrengthItems corepb.CCStrengthItem数组转成CCStrengthItem数组
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

// MapGeneralExplains corepb.GeneralExplain数组 转成 GeneralExplain数组
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
