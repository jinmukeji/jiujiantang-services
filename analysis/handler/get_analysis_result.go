package handler

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"strconv"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/ae/v2/biz"
	"github.com/jinmukeji/ae/v2/engine/core"
	"github.com/jinmukeji/ae/v2/engine/render"
	"github.com/jinmukeji/go-pkg/mathutil"
	"github.com/jinmukeji/plat-pkg/rpc/errors"
	"github.com/jinmukeji/plat-pkg/rpc/errors/codes"
	"github.com/jinmukeji/plat-pkg/rpc/errors/errmsg"

	"fmt"

	"time"

	"github.com/jinmukeji/go-pkg/crypto/rand"
	"github.com/jinmukeji/jiujiantang-services/analysis/aws"
	"github.com/jinmukeji/jiujiantang-services/analysis/mysqldb"
	jinmuidpb "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	analysispb "github.com/jinmukeji/proto/gen/micro/idl/jm/analysis/v1"
	pulsetestinfopb "github.com/jinmukeji/proto/gen/micro/idl/jm/pulsetestinfo/v1"
	subscriptionpb "github.com/jinmukeji/proto/gen/micro/idl/jm/subscription/v1"
	generalpb "github.com/jinmukeji/proto/gen/micro/idl/ptypes/v2"
	ptypesv2 "github.com/jinmukeji/proto/gen/micro/idl/ptypes/v2"
)

// DefaultReportVersion 默认的分析报告版本
const (
	DefaultReportVersion = "2.0"
	customSanshui        = "custom_sanshui"
)

// AnswerKeys 答案的Key
type AnswerKeys []string

// Answer 答案
type Answer map[string]AnswerKeys

// GetAnalyzeResult 得到分析结果
func (j *AnalysisManagerService) GetAnalyzeResult(ctx context.Context, req *analysispb.GetAnalyzeResultRequest, resp *analysispb.GetAnalyzeResultResponse) error {

	record, err := j.database.FindRecordByRecordID(req.RecordId)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find analysis params of record %d: %s", req.RecordId, err.Error()))
	}

	// 如果分析完毕，不可以再次走问答
	if record.AnalyzeStatus == mysqldb.AnalysisStatusCompeleted {
		return NewError(ErrDatabase, fmt.Errorf("status of current record %d has been completed", req.RecordId))
	}
	// 如果分析出错，不可以再次走问答
	if record.AnalyzeStatus == mysqldb.AnalysisStatusError {
		return NewError(ErrDatabase, fmt.Errorf("status of current record %d is error", req.RecordId))
	}

	var ownerID int32
	// 不需言忽略 token 检查
	if !req.IsSkipVerifyToken {
		token, _ := TokenFromContext(ctx)
		ownerID, _ = j.database.FindUserIDByToken(token)
		reqCheckUserHaveSameSubscription := new(subscriptionpb.CheckUserHaveSameSubscriptionRequest)
		reqCheckUserHaveSameSubscription.UserId = int32(record.UserID)
		reqCheckUserHaveSameSubscription.OwnerId = ownerID
		respGetSelectedUserSubscription, errGetSelectedUserSubscription := j.subscriptionSvc.CheckUserHaveSameSubscription(ctx, reqCheckUserHaveSameSubscription)
		if errGetSelectedUserSubscription != nil || !respGetSelectedUserSubscription.IsSameSubscription {
			return fmt.Errorf("user %d from request does not have same subscription of owner %d from token", record.UserID, ownerID)
		}
	}
	// 获得用户个人信息
	reqGetUserProfile := new(jinmuidpb.GetUserProfileRequest)
	if ownerID != record.UserID {
		reqGetUserProfile.IsSkipVerifyToken = true
	}
	reqGetUserProfile.UserId = record.UserID
	respGetUserProfile, err := j.jinmuidSvc.GetUserProfile(ctx, reqGetUserProfile)
	if err != nil {
		return NewError(ErrInvalidUser, fmt.Errorf("failed to get user profile of user [%d]: %s", record.UserID, err.Error()))
	}
	userProfile := respGetUserProfile.GetProfile()

	// 站姿走滤镜
	if record.MeasurementPosture == mysqldb.MeasurementPostureStanging {
		record.C0, record.C1, record.C2, record.C3, record.C4, record.C5, record.C6, record.C7 = ConvertToStandingCFloat64Values(userProfile.GetGender(), int(record.C0), int(record.C1), int(record.C2), int(record.C3), int(record.C4), int(record.C5), int(record.C6), int(record.C7))
	}

	clientID := record.ClientID
	// 根据客户端 ID 获得对应的提示助手名称
	assistantName := getAssistantNameByClientId(clientID)

	// 构建 ae 需要的问答
	var answers map[string]AnalysisReportAnswers
	answers = getAnswers(req.GetQuestionAnswers())

	ctxData, err := buildAEInput(userProfile, record, answers, req.GetCid(), true)
	if err != nil {
		return NewError(ErrLoadConfig, fmt.Errorf("failed to build ae input of record [%d]: %s", record.RecordID, err.Error()))
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
		// 更新 ae 分析状态为有错
		err = j.database.UpdateRecordHasAEError(req.GetRecordId(), mysqldb.HasAeError)
		if err != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to update record has ae error of record %d: %s", req.GetRecordId(), err.Error()))
		}
		return NewError(ErrAEError, fmt.Errorf("error occurs when run engine: %s", errRunEngine.Error()))
	}

	// 如果还有需要提问的问题，则更新分析状态为等待用户输入，则返回问题
	if !ctxData.Output["has_answered_all_questions"].(bool) {
		err = j.database.UpdateAnalysisStatusInProgress(req.RecordId)
		if err != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to update analysis status in progress of record %d: %s", req.RecordId, err.Error()))
		}

		askQuestions, _ := ctxData.Output["ask_questions"].(map[interface{}]interface{})
		resp.Questions = parseOutputQuestions(askQuestions)
		return nil

	}
	// 如果所有问题回答完毕
	if ctxData.Output["has_answered_all_questions"].(bool) {
		// 存储体制辩证信息，便于周报月报分析
		physicalDialecticLookups := parseLookupsIncludeKeyContent(ctxData.Output["physical_dialectics"])
		physicalDialectics := make([]AnalysisReportRequestBodyInputKey, len(physicalDialecticLookups))
		for idx, value := range physicalDialecticLookups {
			physicalDialectics[idx] = AnalysisReportRequestBodyInputKey{
				Key:   value.GetKey(),
				Score: int32(value.GetScore()),
			}
		}
		// 存储分析 body
		analyzeBody := AnalyzeBody{
			TransactionID:      strconv.Itoa(int(req.GetRecordId())),
			QuestionAnswers:    answers,
			Language:           language,
			PhysicalDialectics: physicalDialectics,
		}
		body, _ := json.Marshal(&analyzeBody)

		// 构建脉搏波测量分析表的 tags
		hasStressState := false
		stressStatus := map[string]bool{}
		if ctxData.Output["stress_state_judgment"] != nil {
			stressStateJudgment, ok := ctxData.Output["stress_state_judgment"].(map[string]interface{})
			if ok {
				if parseOutputBool(stressStateJudgment, "Has_stress_state") {
					hasStressState = true
				}
				stressStatus["has_stress_state"] = parseOutputBool(stressStateJudgment, "Has_stress_state")
				stressStatus["has_done_sports"] = parseOutputBool(stressStateJudgment, "has_done_sports")
				stressStatus["has_drinked_wine"] = parseOutputBool(stressStateJudgment, "has_drinked_wine")
				stressStatus["has_had_cold"] = parseOutputBool(stressStateJudgment, "has_had_cold")
				stressStatus["has_rhinitis_episode"] = parseOutputBool(stressStateJudgment, "has_rhinitis_episode")
				stressStatus["has_abdominal_pain"] = parseOutputBool(stressStateJudgment, "has_abdominal_pain")
				stressStatus["has_viral_infection"] = parseOutputBool(stressStateJudgment, "has_viral_infection")
				stressStatus["has_physiological_period"] = parseOutputBool(stressStateJudgment, "has_physiological_period")
				stressStatus["has_ovulation"] = parseOutputBool(stressStateJudgment, "has_ovulation")
				stressStatus["has_pregnant"] = parseOutputBool(stressStateJudgment, "has_pregnant")
				stressStatus["has_lactation"] = parseOutputBool(stressStateJudgment, "has_lactation")
			}

		}

		stringStressStatus, _ := json.Marshal(&stressStatus)
		// 记录流水号
		transactionNo, err := genTransactionNumber(req.RecordId)
		if err != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to generate transaction number: %s", err.Error()))
		}

		errUpdateRecordTransactionNumber := j.database.UpdateRecordTransactionNumber(ctx, req.RecordId, transactionNo)
		if errUpdateRecordTransactionNumber != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to update record transaction number of record %d: %s", record.RecordID, errUpdateRecordTransactionNumber.Error()))
		}
		// 更新分析状态为完成
		r := &mysqldb.Record{
			RecordID:       req.RecordId,
			HasStressState: hasStressState,
			StressState:    string(stringStressStatus),
			AnalyzeBody:    string(body),
			AnalyzeStatus:  mysqldb.AnalysisStatusCompeleted,
		}

		err = j.database.UpdateAnalysisRecord(r)
		if err != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to update analysis record of recordID %d: %s", req.RecordId, err.Error()))
		}
		// 构建返回的内容
		hasAbnormalMeasurementFromAeOutput := checkIfHasAbnormalMeasurementFromAeOutput(*ctxData)
		displayOptions := getDisplayOption(record.ClientID, userProfile.GetGender(), checkIfHasStressStateFromAeOutput(*ctxData), hasAbnormalMeasurementFromAeOutput, true)

		respModules, err := getModulesFromAEOutput(*ctxData, displayOptions, userProfile.GetGender())
		if err != nil {
			return errors.ErrorWithCause(codes.InvalidOperation, err, "failed to get modules from ae output")
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

		resp.RecordId = req.GetRecordId()
		resp.ReportVersion = DefaultReportVersion
		resp.Report = &analysispb.ReportContent{
			RecordId:    req.GetRecordId(),
			UserProfile: userProfileModule,
			PulseTest:   pulseTestModule,
			Remark:      remarkModule,
			Modules:     respModules,
			CreatedTime: protoAnalysisFinishTime,
		}
		resp.TransactionId = record.TransactionNumber
	}

	return nil
}

// getAnswers 建立模块名到回答到对应关系
func getAnswers(answers map[string]*analysispb.Answers) map[string]AnalysisReportAnswers {
	questionAnswers := make(map[string]AnalysisReportAnswers)
	for module, qa := range answers {
		a := []AnalysisReportAnswer{}
		for _, answer := range qa.Answers {
			keys := make([]string, len(answer.AnswerKeys))
			copy(keys, answer.AnswerKeys)
			a = append(a, AnalysisReportAnswer{
				QuestionKey: answer.GetQuestionKey(),
				AnswerKeys:  keys,
			})
		}
		questionAnswers[module] = a
	}
	return questionAnswers
}

const (
	// AELanguageSimplifiedChinese 简体中文
	AELanguageSimplifiedChinese = "zh-Hans"
	// AELanguageTraditionalChinese 繁体中文
	AELanguageTraditionalChinese = "zh-Hant"
	// AELanguageEnglish 英文
	AELanguageEnglish = "en"
)

func mapProtoLanguageToAE(language ptypesv2.Language) (string, error) {
	switch language {
	case ptypesv2.Language_LANGUAGE_SIMPLIFIED_CHINESE:
		return AELanguageSimplifiedChinese, nil
	case ptypesv2.Language_LANGUAGE_TRADITIONAL_CHINESE:
		return AELanguageTraditionalChinese, nil
	case ptypesv2.Language_LANGUAGE_ENGLISH:
		return AELanguageEnglish, nil
	}
	return "", fmt.Errorf(errmsg.InvalidValue("language", language))
}

// AnalyzeBody 新分析接口的body
type AnalyzeBody struct {
	TransactionID      string                              `json:"transaction_id"`
	QuestionAnswers    map[string]AnalysisReportAnswers    `json:"question_answers"`
	Language           string                              `json:"language"`
	PhysicalDialectics []AnalysisReportRequestBodyInputKey `json:"physical_dialectics"`
	Disease            []AnalysisReportRequestBodyInputKey `json:"disease"`
	DirtyDialectic     []AnalysisReportRequestBodyInputKey `json:"dirty_dialectic"`
}

// parseOutputQuestions 解析问题
func parseOutputQuestions(askQuestions map[interface{}]interface{}) map[string]*analysispb.Questions {
	questions := make(map[string]*analysispb.Questions)
	for module, v := range askQuestions {
		array, ok := v.([]interface{})
		qs := make([]*analysispb.Question, len(array))
		if ok {
			for idx, question := range array {
				if question, ok := question.(map[string]interface{}); ok {
					choices, ok := question["choices"].([]render.Choice)
					cs := make([]*analysispb.QuestionChoice, len(choices))
					if ok {
						for idx, choice := range choices {
							cs[idx] = &analysispb.QuestionChoice{
								ChoiceKey:    choice.Key,
								Content:      choice.Content,
								ConflictKeys: choice.ConflictKeys,
							}
						}
					}

					qs[idx] = &analysispb.Question{
						QuestionKey: question["key"].(string),
						Content:     question["content"].(string),
						Type:        question["type"].(string),
						Choices:     cs,
					}
					questions[module.(string)] = &analysispb.Questions{
						Questions: qs,
					}
				}
			}
		}
	}
	return questions
}

// parseOutputF 解析Bool
func parseOutputBool(module map[string]interface{}, f string) bool {
	if module[f] == nil {
		return false
	}
	b, ok := module[f].(bool)
	if !ok {
		return false
	}
	return b
}

// checkIfHasAbnormalMeasurementFromAeOutput 根据 ae 的输出判断测量异常
func checkIfHasAbnormalMeasurementFromAeOutput(ctxData core.ContextData) bool {
	out := ctxData.Output
	mapMeasurementJudgment, ok := out["measurement_judgment"].(map[string]interface{})
	if ok {
		if mapMeasurementJudgment["lookups"] != nil {
			return true
		}
	}
	// 默认测量正常
	return false
}

// checkIfHasStressStateFromAeOutput 根据 ae 的输出判断是否处于应激态
func checkIfHasStressStateFromAeOutput(ctxData core.ContextData) bool {
	out := ctxData.Output
	mapStressState, ok := out["stress_state_judgment"].(map[string]interface{})
	if ok {
		return parseOutputBool(mapStressState, "has_stress_state")
	}
	// 默认处于非应激态
	return false
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
	return fmt.Sprintf("%s%d", str, getCheckSum(str))
}

// getCheckSum 得到CheckSum
func getCheckSum(str string) int {
	bytes := []byte(str)
	sum := 0
	for i := 0; i < len(bytes); i++ {
		sum = sum + int(bytes[i])
	}
	return sum % 10
}

// getPulseTestDataIntArray 从aws上得到波形数据
func getPulseTestDataIntArray(s3Key string, client aws.PulseTestRawDataS3Client) ([]int32, error) {
	pulseTestRawData, err := client.Download(s3Key)
	if err != nil {
		return []int32{}, err
	}
	return ParsePayload(pulseTestRawData)
}

// 蓝牙数据每三个为一组
const blutoothDataSegmentLength = 3

// ParsePayload 解析Payload
func ParsePayload(pulseTestRawData *pulsetestinfopb.PulseTestRawInfo) ([]int32, error) {
	payload := pulseTestRawData.Payloads
	if pulseTestRawData.Spec == 2 {
		waveData := make([]int32, 0)
		for i := 0; i < len(payload); i += blutoothDataSegmentLength {
			val := int32(payload[i])<<16 + int32(payload[i+1])<<8 + int32(payload[i+2])
			waveData = append(waveData, val)
		}
		return waveData, nil
	}
	scanner := bufio.NewScanner(bytes.NewReader(payload))
	waveData := make([]int32, 0)
	for scanner.Scan() {
		line := scanner.Text()
		d, parseErr := strconv.Atoi(line[:len(line)-1])
		if parseErr != nil {
			continue
		}
		waveData = append(waveData, int32(d))
	}
	return waveData, nil
}

// mapDBFingerToProto 将数据库里面存放的 finger 映射为 proto 格式
func mapDBFingerToProto(dbFinger mysqldb.Finger) (analysispb.Finger, error) {
	switch dbFinger {
	case mysqldb.FingerLeft1:
		return analysispb.Finger_FINGER_LEFT_1, nil
	case mysqldb.FingerLeft2:
		return analysispb.Finger_FINGER_LEFT_2, nil
	case mysqldb.FingerLeft3:
		return analysispb.Finger_FINGER_LEFT_3, nil
	case mysqldb.FingerLeft4:
		return analysispb.Finger_FINGER_LEFT_4, nil
	case mysqldb.FingerLeft5:
		return analysispb.Finger_FINGER_LEFT_5, nil
	case mysqldb.FingerRight1:
		return analysispb.Finger_FINGER_RIGHT_1, nil
	case mysqldb.FingerRight2:
		return analysispb.Finger_FINGER_RIGHT_2, nil
	case mysqldb.FingerRight3:
		return analysispb.Finger_FINGER_RIGHT_3, nil
	case mysqldb.FingerRight4:
		return analysispb.Finger_FINGER_RIGHT_4, nil
	case mysqldb.FingerRight5:
		return analysispb.Finger_FINGER_RIGHT_5, nil
	}
	return analysispb.Finger_FINGER_INVALID, fmt.Errorf("invalid database finger %d", dbFinger)
}

// mapDBGenderToProto 将数据库使用的 gender 映射为 proto 格式
func mapDBGenderToProto(gender mysqldb.Gender) (generalpb.Gender, error) {
	switch gender {
	case mysqldb.GenderFemale:
		return generalpb.Gender_GENDER_FEMALE, nil
	case mysqldb.GenderMale:
		return generalpb.Gender_GENDER_MALE, nil
	case mysqldb.GenderInvalid:
		return generalpb.Gender_GENDER_INVALID, fmt.Errorf("invalid database gender %s", gender)
	}
	return generalpb.Gender_GENDER_INVALID, fmt.Errorf("invalid database gender %s", gender)
}

// Int32ValBoundedBy10FromFloat 将 float64 类型的数据转换为 -10 到 10 之间的整数
func Int32ValBoundedBy10FromFloat(val float64) int32 {
	int32Val := mathutil.RoundToInt32(val)
	if int32Val < -10 {
		return -10
	}
	if int32Val > 10 {
		return 10
	}
	return int32Val
}

// IntValBoundedBy10FromFloat 将 float64 类型的数据转换为 -10 到 10 之间的整数
func IntValBoundedBy10FromFloat(val float64) int {
	intVal := mathutil.RoundToInt(val)
	if intVal < -10 {
		return -10
	}
	if intVal > 10 {
		return 10
	}
	return intVal
}

// genTransactionNumber 生成transactionNumber
func genTransactionNumber(recordID int32) (string, error) {
	// 年月日时分（10位）+流水号（前3位）+随机码（2位）+流水号（后3位）+校验码（1位）
	now := time.Now()
	date := fmt.Sprintf("%04d%02d%02d%02d%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute())
	ID := strconv.Itoa(int(recordID))
	if len(ID) < 6 {
		ID = fmt.Sprintf("%06d", recordID)
	}
	rs, _ := rand.RandomStringWithMask(rand.MaskDigits, 2)
	str := fmt.Sprintf("%s%s%s%s", date, Substr(ID, 0, 3), rs, Substr(ID, len(ID)-3, len(ID)))
	return addCheckSumCodeString(str), nil
}
