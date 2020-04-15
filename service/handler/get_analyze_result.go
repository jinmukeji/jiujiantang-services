package handler

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/ae-v1/core"
	"github.com/jinmukeji/ae-v1/model"
	"github.com/jinmukeji/go-pkg/v2/age"
	"github.com/jinmukeji/go-pkg/v2/crypto/rand"
	"github.com/jinmukeji/jiujiantang-services/pkg/rpc"
	"github.com/jinmukeji/jiujiantang-services/service/auth"
	"github.com/jinmukeji/jiujiantang-services/service/mysqldb"
	corepb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	jinmuidpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	ptypespb "github.com/jinmukeji/proto/v3/gen/micro/idl/ptypes/v1"
	generalpb "github.com/jinmukeji/proto/v3/gen/micro/idl/ptypes/v2"
)

const (
	// DefaultAnalysisReportVersion 默认分析报告的版本
	DefaultAnalysisReportVersion = "1.0"
	// ReportErrorKey 报告错误的key
	ReportErrorKey = "C0001.0"
	// zoneDengyun 登云的Zone
	zoneDengyun = "CN-X"
)

var regContentLink *regexp.Regexp

func init() {
	/*
		aaa[xxx](#key#)bbb
		xxx 任意中文字符串
		key 任意 数字 字母 下划线 .
	*/
	regContentLink = regexp.MustCompile(`\[([\w\p{Han}]+)\]\(#[\w\.]+#\)`)
}

// GetAnalyzeResult 得到分析结果
func (j *JinmuHealth) GetAnalyzeResult(ctx context.Context, req *corepb.GetAnalyzeResultRequest, resp *corepb.GetAnalyzeResultResponse) error {
	l := rpc.ContextLogger(ctx)
	accessTokenType, _ := auth.AccessTokenTypeFromContext(ctx)

	record, errGetRecordByRecordID := j.datastore.FindRecordByID(ctx, int(req.RecordId))
	if record == nil || errGetRecordByRecordID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find record by recordID %d: %s", req.RecordId, errGetRecordByRecordID.Error()))
	}
	measurementResult, errbuildMeasurementResult := j.buildMeasurementResult(ctx, record, int(req.UserId))
	if errbuildMeasurementResult != nil {
		return NewError(ErrGetMeasureResultFailure, errbuildMeasurementResult)
	}
	reportUserProfile, errbuildUserProfile := j.buildUserProfile(ctx, int32(record.RecordID))
	if errbuildUserProfile != nil {
		return NewError(ErrGetUserFailure, errbuildUserProfile)
	}
	result := measurementResult.ReportMeasurementResult
	//  处理站姿C0-C7
	if measurementResult.MeasurementPosture == corepb.MeasurementPosture_MEASUREMENT_POSTURE_STANDING {
		DealWithMeasurementPosture(result, reportUserProfile.Gender)
	}
	out, errRunAnalysisEngine := j.RunAnalysisEngine(ctx, measurementResult, reportUserProfile, req.SystemTags, req.AnswerTags)
	if errRunAnalysisEngine != nil {
		return NewError(ErrRunAnalysisEngineFailure, errRunAnalysisEngine)
	}
	questions := BuildOutQuestions(out)
	createdAt, _ := ptypes.TimestampProto(measurementResult.CreatedAt)

	resp.Cid = req.RecordId
	resp.Questionnaire = &corepb.Questionnaire{
		Title:       "",
		CreatedTime: createdAt,
		Questions:   questions,
	}

	content := BuildProtoContent(ctx, out, measurementResult, reportUserProfile)
	if record.TransactionNumber == "" {
		// 记录流水号
		transactionNo, err := genReportID(req.TransactionNumber)
		if err != nil {
			return NewError(ErrGenRandomString, fmt.Errorf("failed to generate transaction number: %s", err.Error()))
		}
		errUpdateRecordTransactionNumber := j.datastore.UpdateRecordTransactionNumber(ctx, int32(record.RecordID), transactionNo)
		if errUpdateRecordTransactionNumber != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to update record transaction number of record %d: %s", record.RecordID, errUpdateRecordTransactionNumber.Error()))
		}
		resp.TransactionNo = transactionNo
	} else {
		resp.TransactionNo = record.TransactionNumber
	}
	resp.AnalysisReport = j.BuildAnalysisReportOptions(ctx, resp.TransactionNo, content, reportUserProfile, out, record.RecordType)
	resp.AnalysisDone = hasAnswerAll(questions)
	analysisResultError := hasReportError(out.MeasurementTips)
	resp.AnalysisResultError = analysisResultError
	if analysisResultError {
		resp.AnalysisDone = false
		err := j.datastore.UpdateRecordHasAEError(ctx, int32(record.RecordID), mysqldb.DbValidValue)
		if err != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to update recordHasAEError of record %d: %s", record.RecordID, err.Error()))
		}
		return nil
	}
	if resp.AnalysisDone && (accessTokenType == AccessTokenTypeLValue || accessTokenType == AccessTokenTypeWeChatValue) {
		wxUser, errWxUser := j.datastore.FindWXUserByUserID(ctx, reportUserProfile.UserId)
		if errWxUser != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to find WXuser by userID %d: %s", reportUserProfile.UserId, errWxUser.Error()))
		}

		// 测量完成状态下且当前记录还没有发送微信模板消息的情况下，发送一次微信模板消息
		if resp.AnalysisDone && !record.HasSentWxViewReportNotification && !analysisResultError {
			// machineUUID, _ := auth.MachineUUIDFromContext(ctx) 暂时不使用
			err := j.wechat.SendViewReportTemplateMessage(wxUser.OpenID, record.RecordID, wxUser.Nickname, record.CreatedAt)
			if err != nil {
				l.Warnf("failed to send wx template message. error: %v", err)
				return NewError(ErrSendTemplateMessageFailure, err)
			}

			err = j.datastore.SetRecordSentWxViewReportNotification(ctx, record.RecordID, true)
			if err != nil {
				return NewError(ErrDatabase, fmt.Errorf("failed to set record sent WxViewReport notification of record %d: %s", record.RecordID, err.Error()))
			}
		}
	}
	return nil
}

func replaceGeneralItem(generalItems []model.GeneralItem) []model.GeneralItem {
	items := make([]model.GeneralItem, len(generalItems))
	copy(items, generalItems)
	for idx, item := range generalItems {
		items[idx].Content = replaceMatchString(item.Content)
	}
	return items
}

func replaceMatchString(content string) string {
	return regContentLink.ReplaceAllString(content, `$1`)
}

// BuildProtoContent 构造proto的content
func BuildProtoContent(ctx context.Context, out *core.Output, measurementResult *MeasurementResult, reportUserProfile *corepb.ReportUserProfile) *corepb.Content {
	createdAt, _ := ptypes.TimestampProto(measurementResult.CreatedAt)
	content := &corepb.Content{
		Lead:                                      mapOutput(out.Leads),
		TipsForWoman:                              mapOutput(out.TipsForWoman),
		ChannelsAndCollateralsExplains:            mapOutput(out.CCExplains),
		ConstitutionDifferentiationExplains:       mapOutput(out.CDExplains),
		SyndromeDifferentiationExplains:           mapOutput(out.SDExplains),
		ChannelsAndCollateralsStrength:            mapOutputCCStrengthItem(out.CCStrength),
		BabyTips:                                  mapOutput(out.BabyTips),
		ConstitutionDifferentiationExplainNotices: mapOutput(out.CDExplainNotices),
		Tags:                                  mapOutput(out.AnswerTags),
		DictionaryEntries:                     mapOutput(out.DictionaryEntries),
		FactorExplains:                        mapOutput(out.FactorExplains),
		HealthDescriptions:                    mapOutput(out.HealthDescriptions),
		MeasurementTips:                       mapOutput(out.MeasurementTips),
		ChannelsAndCollateralsExplainNotices:  mapOutput(out.CCExplainNotices),
		UserProfile:                           reportUserProfile,
		MeasurementResult:                     measurementResult.ReportMeasurementResult,
		HasPaid:                               measurementResult.HasPaid,
		CreatedTime:                           createdAt,
		F0:                                    int32(out.F0),
		F1:                                    int32(out.F1),
		F2:                                    int32(out.F2),
		F3:                                    int32(out.F3),
		Remark:                                measurementResult.Remark,
		UterineHealthIndexes:                  mapOutput(out.UterineHealthIndexes),
		UterusAttentionPrompts:                mapOutput(out.UterusAttentionPrompts),
		UterineHealthDescriptions:             mapOutput(out.UterineHealthDescriptions),
		MenstrualHealthValues:                 mapOutput(out.MenstrualHealthValues),
		MenstrualHealthDescriptions:           mapOutput(out.MenstrualHealthDescriptions),
		GynecologicalInflammations:            mapOutput(out.GynecologicalInflammations),
		GynecologicalInflammationDescriptions: mapOutput(out.GynecologicalInflammationDescriptions),
		BreastHealth:                          mapOutput(out.BreastHealth),
		BreastHealthDescriptions:              mapOutput(out.BreastHealthDescriptions),
		EmotionalHealthIndexes:                mapOutput(out.EmotionalHealthIndexes),
		EmotionalHealthDescriptions:           mapOutput(out.EmotionalHealthDescriptions),
		FacialSkins:                           mapOutput(out.FacialSkins),
		FacialSkinDescriptions:                mapOutput(out.FacialSkinDescriptions),
		ReproductiveAgeConsiderations:         mapOutput(out.ReproductiveAgeConsiderations),
		BreastCancerOvarianCancers:            mapOutput(out.BreastCancerOvarianCancers),
		BreastCancerOvarianCancerDescriptions: mapOutput(out.BreastCancerOvarianCancerDescriptions),
		HormoneLevels:                         mapOutput(out.HormoneLevels),
		LymphaticHealth:                       mapOutput(out.LymphaticHealth),
		LymphaticHealthDescriptions:           mapOutput(out.LymphaticHealthDescriptions),
		F100:                                  int32(out.F100()),
		F101:                                  int32(out.F101()),
		F102:                                  int32(out.F102()),
		F103:                                  int32(out.F103()),
		F104:                                  int32(out.F104()),
		F105:                                  int32(out.F105()),
		F106:                                  int32(out.F106()),
		F107:                                  int32(out.F107()),
	}
	client, _ := clientFromContext(ctx)
	if client.ClientID == dengyunClient {
		content.MeasurementResult.C0 = content.MeasurementResult.C0 * 10
		content.MeasurementResult.C1 = content.MeasurementResult.C1 * 10
		content.MeasurementResult.C2 = content.MeasurementResult.C2 * 10
		content.MeasurementResult.C3 = content.MeasurementResult.C3 * 10
		content.MeasurementResult.C4 = content.MeasurementResult.C4 * 10
		content.MeasurementResult.C5 = content.MeasurementResult.C5 * 10
		content.MeasurementResult.C6 = content.MeasurementResult.C6 * 10
		content.MeasurementResult.C7 = content.MeasurementResult.C7 * 10
	}
	content.M0 = NewNullableInt32Value(out.M0)
	content.M1 = NewNullableInt32Value(out.M1)
	content.M2 = NewNullableInt32Value(out.M2)
	content.M3 = NewNullableInt32Value(out.M3)
	return content
}

// NewNullableInt32Value 创建一个可空的int32值
func NewNullableInt32Value(v *int) *ptypespb.NullableInt32Value {
	if v != nil {
		return &ptypespb.NullableInt32Value{
			Kind: &ptypespb.NullableInt32Value_Int32Value{
				Int32Value: int32(*v),
			},
		}
	}
	return &ptypespb.NullableInt32Value{
		Kind: &ptypespb.NullableInt32Value_NullValue{},
	}
}

// BuildOutQuestions 构造OutQuestions
func BuildOutQuestions(out *core.Output) []*corepb.Question {
	questions := make([]*corepb.Question, len(out.Questions))
	for idx, question := range out.Questions {
		choices := make([]*corepb.Choice, len(question.Choices))
		for idx, questionChoice := range question.Choices {
			choices[idx] = &corepb.Choice{
				Key:          questionChoice.Key,
				Name:         questionChoice.Label,
				Value:        questionChoice.Key,
				ConflictKeys: questionChoice.ConflictKeys,
			}
		}
		questions[idx] = &corepb.Question{
			Key:         question.Key,
			Title:       question.Label,
			Description: question.Content,
			Tip:         question.Tip,
			Type:        string(question.Type),
			Choices:     choices,
			DefaultKeys: question.DefaultSelectedChoiceKeys,
		}
	}
	return questions
}

// RunAnalysisEngine 运行分析引擎
func (j *JinmuHealth) RunAnalysisEngine(ctx context.Context, measurementResult *MeasurementResult, reportUserProfile *corepb.ReportUserProfile, systemTags []string, answerTags []string) (*core.Output, error) {

	rs := j.analysisEngine.RuleSetDoc.Get(j.analysisEngine.RuleSetDoc.DefaultRuleSetKey)
	input := core.NewInput()
	input.SetSystemTags(systemTags...)
	input.HeartRate = int(measurementResult.ReportMeasurementResult.AppHr)
	input.C0 = int(measurementResult.ReportMeasurementResult.C0) * 10
	input.C1 = int(measurementResult.ReportMeasurementResult.C1) * 10
	input.C2 = int(measurementResult.ReportMeasurementResult.C2) * 10
	input.C3 = int(measurementResult.ReportMeasurementResult.C3) * 10
	input.C4 = int(measurementResult.ReportMeasurementResult.C4) * 10
	input.C5 = int(measurementResult.ReportMeasurementResult.C5) * 10
	input.C6 = int(measurementResult.ReportMeasurementResult.C6) * 10
	input.C7 = int(measurementResult.ReportMeasurementResult.C7) * 10
	input.Weight = int(reportUserProfile.Weight)
	input.Age = int(reportUserProfile.Age)
	input.Height = int(reportUserProfile.Height)
	for _, tag := range answerTags {
		input.AnswerTags[tag] = true
	}
	gender, errMapProtoGenderToEngineInput := mapProtoGenderToEngineInput(reportUserProfile.Gender)
	if errMapProtoGenderToEngineInput != nil {
		return nil, errMapProtoGenderToEngineInput
	}
	input.Gender = gender
	aeCtx := j.analysisEngine.NewContext(input)
	rs.Run(aeCtx)
	out := aeCtx.Output()
	client, _ := clientFromContext(ctx)
	customizedCode := client.CustomizedCode
	zone := client.Zone
	if customizedCode == customDengyun && zone == zoneDengyun {
		out.BreastCancerOvarianCancerDescriptions = replaceGeneralItem(out.BreastCancerOvarianCancerDescriptions)
		out.BreastCancerOvarianCancers = replaceGeneralItem(out.BreastCancerOvarianCancers)
		out.BreastHealth = replaceGeneralItem(out.BreastHealth)
		out.BreastHealthDescriptions = replaceGeneralItem(out.BreastHealthDescriptions)
		out.CCExplainNotices = replaceGeneralItem(out.CCExplainNotices)
		out.CCExplains = replaceGeneralItem(out.CCExplains)
		out.CDExplainNotices = replaceGeneralItem(out.CDExplainNotices)
		out.CDExplains = replaceGeneralItem(out.CDExplains)
		out.DictionaryEntries = replaceGeneralItem(out.DictionaryEntries)
		out.EmotionalHealthDescriptions = replaceGeneralItem(out.EmotionalHealthDescriptions)
		out.EmotionalHealthIndexes = replaceGeneralItem(out.EmotionalHealthIndexes)
		out.FacialSkinDescriptions = replaceGeneralItem(out.FacialSkinDescriptions)
		out.FacialSkins = replaceGeneralItem(out.FacialSkins)
		out.FactorExplains = replaceGeneralItem(out.FactorExplains)
		out.GynecologicalInflammationDescriptions = replaceGeneralItem(out.GynecologicalInflammationDescriptions)
		out.GynecologicalInflammations = replaceGeneralItem(out.GynecologicalInflammations)
		out.HealthDescriptions = replaceGeneralItem(out.HealthDescriptions)
		out.HormoneLevels = replaceGeneralItem(out.HormoneLevels)
		out.Leads = replaceGeneralItem(out.Leads)
		out.LymphaticHealth = replaceGeneralItem(out.LymphaticHealth)
		out.LymphaticHealthDescriptions = replaceGeneralItem(out.LymphaticHealthDescriptions)
		out.MeasurementTips = replaceGeneralItem(out.MeasurementTips)
		out.MenstrualHealthDescriptions = replaceGeneralItem(out.MenstrualHealthDescriptions)
		out.MenstrualHealthValues = replaceGeneralItem(out.MenstrualHealthValues)
		out.ReproductiveAgeConsiderations = replaceGeneralItem(out.ReproductiveAgeConsiderations)
		out.SDExplains = replaceGeneralItem(out.SDExplains)
		out.TipsForWoman = replaceGeneralItem(out.TipsForWoman)
		out.UterineHealthDescriptions = replaceGeneralItem(out.UterineHealthDescriptions)
		out.UterineHealthIndexes = replaceGeneralItem(out.UterineHealthIndexes)
		out.UterusAttentionPrompts = replaceGeneralItem(out.UterusAttentionPrompts)
	}
	return out, nil
}

// MeasurementResult 测量记录
type MeasurementResult struct {
	ReportMeasurementResult *corepb.ReportMeasurementResult
	CreatedAt               time.Time
	HasPaid                 bool
	Remark                  string
	MeasurementPosture      corepb.MeasurementPosture
}

// buildMeasurementResult 构建MeasurementResult
func (j *JinmuHealth) buildMeasurementResult(ctx context.Context, record *mysqldb.Record, reqUserID int) (*MeasurementResult, error) {

	accessTokenType, _ := auth.AccessTokenTypeFromContext(ctx)

	// TODO: 以后修改
	if accessTokenType == AccessTokenTypeWeChatValue {
		if record.UserID != reqUserID {
			return nil, errors.New("user_id and record doesn't match")
		}
	}

	if accessTokenType != AccessTokenTypeLValue && accessTokenType != AccessTokenTypeWeChatValue {
		ownerID, ok := auth.UserIDFromContext(ctx)
		if !ok {
			return nil, errors.New("failed to get ownerID from context")
		}
		userOrganization, _ := j.datastore.FindOrganizationByUserID(ctx, record.UserID)
		ownerOrganization, _ := j.datastore.FindOrganizationByUserID(ctx, int(ownerID))
		if ownerOrganization.OrganizationID != userOrganization.OrganizationID {
			return nil, fmt.Errorf("this user_id %d cannot get measurement by record_id", record.UserID)
		}
	}
	return j.buildAnalysisResult(ctx, record, reqUserID)
}

// buildUserProfile 构建proto的用户信息
func (j *JinmuHealth) buildUserProfile(ctx context.Context, recordID int32) (*corepb.ReportUserProfile, error) {
	accessTokenType, _ := auth.AccessTokenTypeFromContext(ctx)
	// TODO: 以后修改
	if accessTokenType != AccessTokenTypeLValue && accessTokenType != AccessTokenTypeWeChatValue {
		userID, errGetUserIDByRecordID := j.datastore.GetUserIDByRecordID(ctx, recordID)
		if errGetUserIDByRecordID != nil {
			return nil, NewError(ErrDatabase, fmt.Errorf("failed to get userID by recordID %d: %s", recordID, errGetUserIDByRecordID.Error()))
		}
		ownerID, ok := auth.UserIDFromContext(ctx)
		if !ok {
			return nil, errors.New("failed to get userID from context")
		}
		userOrganization, _ := j.datastore.FindOrganizationByUserID(ctx, int(userID))
		ownerOrganization, _ := j.datastore.FindOrganizationByUserID(ctx, int(ownerID))
		if ownerOrganization.OrganizationID != userOrganization.OrganizationID {
			return nil, errors.New("failed to get user")
		}
	}

	return j.buildAnalysisUserProfile(ctx, recordID)
}

// buildAnalysisUserProfile 构建proto的用户信息
func (j *JinmuHealth) buildAnalysisUserProfile(ctx context.Context, recordID int32) (*corepb.ReportUserProfile, error) {
	isExsit, err := j.datastore.ExistRecordByRecordID(ctx, recordID)
	if err != nil || !isExsit {
		return nil, errors.New("failed to get user")
	}
	reqGetUserProfileByRecordID := new(jinmuidpb.GetUserProfileByRecordIDRequest)
	reqGetUserProfileByRecordID.RecordId = recordID
	reqGetUserProfileByRecordID.IsSkipVerifyToken = true
	respGetUserProfileByRecordID, errGetUserProfileByRecordID := j.jinmuidSvc.GetUserProfileByRecordID(ctx, reqGetUserProfileByRecordID)
	if errGetUserProfileByRecordID != nil {
		return nil, errGetUserProfileByRecordID
	}
	userProfile := respGetUserProfileByRecordID.UserProfile
	birthday, _ := ptypes.Timestamp(userProfile.BirthdayTime)
	profile := &corepb.ReportUserProfile{
		UserId:       int32(respGetUserProfileByRecordID.UserId),
		Nickname:     userProfile.Nickname,
		BirthdayTime: userProfile.BirthdayTime,
		Gender:       userProfile.Gender,
		Weight:       int32(userProfile.Weight),
		Height:       int32(userProfile.Height),
		Age:          int32(age.Age(birthday)),
	}
	if birthday.IsZero() {
		birthday, _ := ptypes.TimestampProto(time.Now())
		profile.BirthdayTime = birthday
		profile.Age = 0
	}

	wxUser, errFindWXUserByUserID := j.datastore.FindWXUserByUserID(ctx, respGetUserProfileByRecordID.UserId)
	if wxUser == nil || errFindWXUserByUserID != nil {
		return profile, nil
	}
	profile.AvatarUrl = wxUser.AvatarImageURL
	return profile, nil
}

const (
	// GenderMale 男性
	GenderMale string = "M"
	// GenderFemale 女性
	GenderFemale string = "F"
	// GenderInvalid 非法的性别
	GenderInvalid string = ""
)

// mapProtoGenderToEngineInput 将proto 类型的 gender 转化为运行引擎需要的 string 类型
func mapProtoGenderToEngineInput(gender generalpb.Gender) (string, error) {
	switch gender {
	case generalpb.Gender_GENDER_FEMALE:
		return GenderFemale, nil
	case generalpb.Gender_GENDER_MALE:
		return GenderMale, nil
	case generalpb.Gender_GENDER_INVALID:
		return GenderInvalid, fmt.Errorf("invalid proto gender %d", generalpb.Gender_GENDER_INVALID)
	case generalpb.Gender_GENDER_UNSET:
		return GenderInvalid, fmt.Errorf("invalid proto gender %d", generalpb.Gender_GENDER_UNSET)
	}
	return GenderInvalid, errors.New("invalid proto gender")
}

func mapOutputCCStrengthItem(source []model.CCStrengthItem) []*corepb.CCStrengthItem {
	ret := make([]*corepb.CCStrengthItem, len(source))
	for idx, o := range source {
		labels := make([]*corepb.CCStrengthLabel, len(o.Labels))
		for idx, ccStrengthLabel := range o.Labels {
			labels[idx] = &corepb.CCStrengthLabel{
				Label: ccStrengthLabel.Label,
				Cc:    ccStrengthLabel.CC,
			}
		}
		ge := &corepb.CCStrengthItem{
			Key:      o.Key,
			Labels:   labels,
			Disabled: o.Disabled,
			Remark:   o.Remark,
		}
		ret[idx] = ge
	}
	return ret
}

func mapOutput(source []model.GeneralItem) []*corepb.GeneralExplain {
	ret := make([]*corepb.GeneralExplain, len(source))
	for idx, o := range source {
		ge := &corepb.GeneralExplain{
			Key:     o.Key,
			Label:   o.Label,
			Content: o.Content,
		}
		ret[idx] = ge
	}

	return ret
}

func genReportID(transactionNumber int32) (string, error) {
	// 年月日时分（10位）+流水号（前3位）+随机码（2位）+流水号（后3位）+校验码（1位）
	now := time.Now()
	format := "2006-01-02 15:04"
	date := Substr(strings.Replace(strings.Replace(strings.Replace(now.Format(format), "-", "", -1), " ", "", -1), ":", "", -1), 2, 12)
	ID := strconv.Itoa(int(transactionNumber))
	if len(ID) >= 6 {
		ID = Substr(ID, 0, 6)
	} else {
		var sl []string
		for i := 0; i < 6-len(ID); i++ {
			sl = append(sl, "0")
		}
		ID = fmt.Sprintf("%s%s", strings.Join(sl, ""), ID)
	}
	rs, _ := rand.RandomStringWithMask(rand.MaskDigits, 2)
	str := fmt.Sprintf("%s%s%s%s", date, Substr(ID, 0, 3), rs, Substr(ID, 3, 6))
	return addCheckSumCodeString(str), nil
}

// Substr 截取字符串 start 起点下标 end 终点下标(不包括)
func Substr(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < 0 || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}

// addCheckSumCodeString 添加验证码
func addCheckSumCodeString(str string) string {
	bytes := []byte(str)
	sum := 0
	for i := 0; i < len(bytes); i++ {
		sum = sum + int(bytes[i])
	}
	return fmt.Sprintf("%s%d", str, sum%10)
}

func hasAnswerAll(questions []*corepb.Question) bool {
	return len(questions) == 0
}

// hasReportError 报告是否异常
func hasReportError(measurementTips []model.GeneralItem) bool {
	for _, tip := range measurementTips {
		if tip.Key == ReportErrorKey {
			return true
		}
	}
	return false
}
