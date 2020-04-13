package rest

import (
	"encoding/json"
	"time"

	"github.com/jinmukeji/go-pkg/v2/crypto/rand"

	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/ae-v1/core"
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	corepb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	ptypespb "github.com/jinmukeji/proto/v3/gen/micro/idl/ptypes/v1"
	generalpb "github.com/jinmukeji/proto/v3/gen/micro/idl/ptypes/v2"
	"github.com/kataras/iris/v12"
)

const seamlessClient = "seamlessiot-10001"

// GetAnalyzeResult 得到智能分析结果
func (h *v2Handler) GetAnalyzeResult(ctx iris.Context) {
	if ctx.Values().GetString(ClientIDKey) == seamlessClient {
		writeError(
			ctx,
			wrapError(ErrDeniedToAccessAPI, "", fmt.Errorf("%s is denied to access this API", seamlessClient)),
			false,
		)
		return
	}
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
		writeError(
			ctx,
			wrapError(ErrInvalidValue, fmt.Sprintf("cid %d is invalid", analysisReportBody.Cid), nil),
			false,
		)
		return
	}
	tags, errGetAnalysisSystemTags := h.getAnalysisSystemTags(ctx)
	if errGetAnalysisSystemTags != nil {
		writeRPCInternalError(ctx, errGetAnalysisSystemTags, false)
		return
	}
	transactionNumber, errGetAnalysisReportTransactionNumber := h.getAnalysisReportTransactionNumber(ctx)
	if errGetAnalysisReportTransactionNumber != nil {
		writeRPCInternalError(ctx, errGetAnalysisReportTransactionNumber, false)
		return
	}
	req := new(corepb.GetAnalyzeResultRequest)
	req.RecordId = int32(recordID)
	req.SystemTags = tags
	answerTags := make([]string, 0)
	for _, answer := range analysisReportBody.Answers {
		for _, value := range answer.Values {
			answerTags = append(answerTags, core.BuildFullQuestionChoiceTag(answer.QuestionKey, value))
		}
	}
	req.AnswerTags = answerTags
	req.TransactionNumber = transactionNumber
	resp, errGetAnalyzeResult := h.rpcSvc.GetAnalyzeResult(
		newRPCContext(ctx), req,
	)
	if errGetAnalyzeResult != nil {
		writeRPCInternalError(ctx, errGetAnalyzeResult, false)
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
			writeRPCInternalError(ctx, errSaveAnswers, false)
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

// MapGeneralExplains 把 proto中的GeneralExplain 转成GeneralExplain
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
	if ctx.Values().GetString(ClientIDKey) == seamlessClient {
		writeError(
			ctx,
			wrapError(ErrDeniedToAccessAPI, "", fmt.Errorf("%s is denied to access this API", seamlessClient)),
			false,
		)
		return
	}
	recordID, err := ctx.Params().GetInt("record_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	tags, errGetAnalysisSystemTags := h.getAnalysisSystemTags(ctx)
	if errGetAnalysisSystemTags != nil {
		writeRPCInternalError(ctx, errGetAnalysisSystemTags, false)
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
		writeRPCInternalError(ctx, errGetAnalysisReportTransactionNumber, false)
		return
	}
	req.TransactionNumber = transactionNumber
	resp, errGetAnalyzeResult := h.rpcSvc.GetAnalyzeResult(
		newRPCContext(ctx), req,
	)
	if errGetAnalyzeResult != nil {
		writeRPCInternalError(ctx, errGetAnalyzeResult, false)
		return
	}
	analysisReport := AnalysisReport{}
	analysisReport.ReportID = resp.TransactionNo
	analysisReport.ReportVersion = resp.AnalysisReport.ReportVersion
	analysisReport.Content, _ = h.getAnalysisReportContent(ctx, resp.AnalysisReport)
	rest.WriteOkJSON(ctx, analysisReport)
}

// GetAnalyzeReportByToken 通过token拿到分析报告
func (h *v2Handler) GetAnalyzeReportByToken(ctx iris.Context) {
	if ctx.Values().GetString(ClientIDKey) == seamlessClient {
		writeError(
			ctx,
			wrapError(ErrDeniedToAccessAPI, "", fmt.Errorf("%s is denied to access this API", seamlessClient)),
			false,
		)
		return
	}
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
	int64Gender, errmapProtoGenderToRest := mapProtoGenderToRest(analysisReport.Content.UserProfile.Gender)
	if errmapProtoGenderToRest != nil {
		return Content{}, errmapProtoGenderToRest
	}
	createdAt, _ := ptypes.Timestamp(analysisReport.Content.CreatedTime)
	finger, errmapProtoFingerToRest := mapProtoFingerToRest(analysisReport.Content.MeasurementResult.Finger)
	if errmapProtoFingerToRest != nil {
		return Content{}, errmapProtoFingerToRest
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
		UserProfile: ReportUserProfile{
			UserID:    int64(analysisReport.Content.UserProfile.UserId),
			Nickname:  analysisReport.Content.UserProfile.Nickname,
			Birthday:  birthday,
			Age:       int64(analysisReport.Content.UserProfile.Age),
			Gender:    int64Gender,
			Height:    int64(analysisReport.Content.UserProfile.Height),
			Weight:    int64(analysisReport.Content.UserProfile.Weight),
			AvatarURL: analysisReport.Content.UserProfile.AvatarUrl,
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
		MeasurementResult: MeasurementResult{
			Finger:              int64(finger),
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
		Remark:    analysisReport.Content.Remark,
		CreatedAt: createdAt,
	}
	content.M0 = ParseNullableInt32Value(analysisReport.Content.M0)
	content.M1 = ParseNullableInt32Value(analysisReport.Content.M1)
	content.M2 = ParseNullableInt32Value(analysisReport.Content.M2)
	content.M3 = ParseNullableInt32Value(analysisReport.Content.M3)
	riskEstimates := make([]RiskEstimate, 0)
	riskEstimates = append(riskEstimates,
		RiskEstimate{
			DiseaseName: "hypertension",
			Status:      0,
		},
		RiskEstimate{
			DiseaseName: "hyperglycemia",
			Status:      1,
		},
		RiskEstimate{
			DiseaseName: "hyperlipidemia",
			Status:      2,
		},
		RiskEstimate{
			DiseaseName: "coronary_heart_disease",
			Status:      0,
		},
		RiskEstimate{
			DiseaseName: "anxiety",
			Status:      1,
		},

		RiskEstimate{
			DiseaseName: "immunity",
			Status:      0,
		},
		RiskEstimate{
			DiseaseName: "fatigue_and_pressure",
			Status:      2,
		},
		RiskEstimate{
			DiseaseName: "sleep_problems",
			Status:      1,
		},
		RiskEstimate{
			DiseaseName: "depression",
			Status:      2,
		},
		RiskEstimate{
			DiseaseName: "gastritis",
			Status:      1,
		},
		RiskEstimate{
			DiseaseName: "insufficient_stomach",
			Status:      0,
		},
		RiskEstimate{
			DiseaseName: "acute_pharyngitis",
			Status:      3,
		},
		RiskEstimate{
			DiseaseName: "chronic_cough",
			Status:      1,
		},
		RiskEstimate{
			DiseaseName: "scoliosis",
			Status:      2,
		},
		RiskEstimate{
			DiseaseName: "cervical_spondylosis",
			Status:      1,
		},
		RiskEstimate{
			DiseaseName: "insufficiency_blood_supply",
			Status:      0,
		},
	)
	if analysisReport.Content.UserProfile.Gender == generalpb.Gender_GENDER_MALE {
		riskEstimates = append(riskEstimates, RiskEstimate{
			DiseaseName: "male_prostate",
			Status:      0,
		})
	} else {
		riskEstimates = append(riskEstimates, RiskEstimate{
			DiseaseName: "gynecological_inflammation",
			Status:      0,
		})
	}
	riskEstimates = append(riskEstimates,
		RiskEstimate{
			DiseaseName: "renal_dysfunction",
			Status:      0,
		},
	)

	content.RiskEstimate = riskEstimates
	content.RiskEstimateTips = []string{
		"有高血压、高血糖、心梗、冠心病这类病，如果已经采取了治疗方法，建议在未服用药物或未治疗前重新测量。并且实际风险可能要比测量的更高。",
	}
	treatmentAdvices := make([]TreatmentAdvice, 0)
	treatmentAdvices = append(treatmentAdvices,
		TreatmentAdvice{
			TreatmentName: "dietary_advice",
			Weigth:        1,
			Advice:        "1.脂肪是较好的能量来源，减少碳水化合物的摄入（脂肪燃烧产生的能量比碳水化合物高，所以较少的脂肪就可以满足较多碳水化合物的同样的能量; 同时碳水化合物产生的废热是脂肪的三倍以上）。<br/>2.碳水化合物、饱和油和不饱和油不要超过自己需要的总量。<br/>3.多吃纤维素，也就是蔬菜瓜果，它们是食物中最好的填充料。叶类植物也是碳水化合物，但是我们不消化纤维素，只摄取了碳水化合物，一举两得有了营养又满足饱感。<br/>4.含碳酸糖水饮料是最差的的碳水化合物，我们测脉诊仪就发现常年喝饮料的人几乎都脾虚厉害，运化失常。<br/>5.蛋白质是最不好的热量来源，蛋白质分解成氨基酸产生的能量最少，产生的废热却最多。<br/>6.多喝水，会促进身体的基本功能自行平衡酸碱度。",
		},
		TreatmentAdvice{
			TreatmentName: "sports_advice",
			Weigth:        0,
			Advice:        "现代社会中大多数人都没有固定的时间去做运动，而一有空的时候反而会呆在家里睡觉或者做大量的运动，真是两个奇怪的现象。这两种方法都很伤害身体的 ，下面来说说正确的运动方式吧。<br/>1.一有时间就休息是不可取的，生命在于运动，总是躺在不仅身体不健康连人也会跟着懒堕的。不防留些时间做一些针对自己的运动吧。像羽毛球 高尔夫 壁球 跳绳 游泳等，不怎么剧烈的，这样的运动会使身体在运动中放松，舒展从而达到健康。<br/>2.一有时间就大量做运动也是不可取的，有些朋友听说运动健康就不管不顾的去运动，不休息，殊不知这样反而会伤害身体，因为长久的不运动，身体忽然做大量的超出范围的运动是很累也很伤身体的。 什么都有有个度，过犹不及，适当是最好的。<br/>3.健康的运动，代表是让身体自己呼吸 舒展 放松，所以不要做剧烈的运动，也不要做大量的运动。慢跑，可以促进肺部的呼吸，对于工作空气环境很恶劣的朋友有一定的作用。同时要多食木耳等清肺。<br/>4.对于长期坐在办公室工作的人们，应该做些让全身归位的运动，由于长时间坐着所以会有一些不良的坐姿 长时间的不活动都会使骨头扭曲。这时应该避免瑜伽，那样会加重扭曲的，可以慢慢散步，或者做一些温和的有氧运动都可以的。<br/>最后是运动后的步骤了。不要再风中吹风图凉快，因为这样是很容易中风的。要最快的时间找到休息的地方，擦干身上的汗液，还有就是不要马上吃喝，胃和肺是需要缓解的。等一会是最好的。",
		},
	)
	content.TreatmentAdvices = treatmentAdvices
	content.SyndromeDifferentiationConstitution = "此次测量的体质是您的瞬间体质，如果想要更加了解您的体质，请查看报周告、月报告。"
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
