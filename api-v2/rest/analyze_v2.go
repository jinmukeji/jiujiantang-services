package rest

import (
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	analysispb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/analysis/v1"
	corepb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	ptypesv2 "github.com/jinmukeji/proto/v3/gen/micro/idl/ptypes/v2"
	"github.com/kataras/iris/v12"
)

// GetV2AnalyzeResult 得到V2版的分析报告
func (h *v2Handler) GetV2AnalyzeResult(ctx iris.Context) {
	if ctx.Values().GetString(ClientIDKey) == seamlessClient {
		writeError(
			ctx, wrapError(ErrDeniedToAccessAPI, "", fmt.Errorf("%s is denied to access this API", seamlessClient)),
			false,
		)
		return
	}
	var analysisReportRequestBody AnalysisReportRequestBody
	errReadJSON := ctx.ReadJSON(&analysisReportRequestBody)
	if errReadJSON != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", errReadJSON), false)
		return
	}
	recordID, err := ctx.Params().GetInt("record_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	reqGetAnalyzeResult := &analysispb.GetAnalyzeResultRequest{}
	reqGetAnalyzeResult.RecordId = int32(recordID)
	reqGetAnalyzeResult.TransactionId = analysisReportRequestBody.TransactionID
	reqGetAnalyzeResult.Cid = rest.GetCidFromContext(ctx)

	qs := make(map[string]*analysispb.Answers)
	for module, questionAnswer := range analysisReportRequestBody.QuestionAnswers {
		as := make([]*analysispb.Answer, len(questionAnswer))
		for idx, answer := range questionAnswer {
			as[idx] = &analysispb.Answer{
				QuestionKey: answer.QuestionKey,
				AnswerKeys:  answer.AnswerKeys,
			}
		}
		qs[module] = &analysispb.Answers{
			Answers: as,
		}
	}
	langauge, err := mapRestLanguageToProto(analysisReportRequestBody.Language)
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	reqGetAnalyzeResult.Language = langauge
	reqGetAnalyzeResult.QuestionAnswers = qs
	// 忽略 token 检查
	reqGetAnalyzeResult.IsSkipVerifyToken = false
	reqGetAnalyzeResult.Cid = rest.GetCidFromContext(ctx)
	respGetAnalyzeResult, err := h.rpcAnalysisSvc.GetAnalyzeResult(newRPCContext(ctx), reqGetAnalyzeResult)
	if err != nil {
		writeRPCInternalError(ctx, err, false)
		return
	}
	analysisReportResponse := AnalysisReportResponse{
		ReportVersion: respGetAnalyzeResult.GetReportVersion(),
		ReportID:      respGetAnalyzeResult.GetReport().GetRecordId(),
		TransactionID: respGetAnalyzeResult.GetTransactionId(),
	}
	askQuestions := getAnalysisAskQuestions(respGetAnalyzeResult)
	if askQuestions != nil {
		analysisReportResponse.AskQuestions = askQuestions
	} else {
		// 没有提问，则构建分析报告中的模块
		// 与引擎分析结果相关模块
		analysisReportContent, err := getAnalysisModules(respGetAnalyzeResult.GetReport().GetModules())
		if err != nil {
			writeRPCInternalError(ctx, wrapError(ErrRPCInternal, "", err), false)
			return
		}
		// 个人信息模块
		userProfileModule, err := getUserProfileModule(respGetAnalyzeResult.GetReport().GetUserProfile())
		if err != nil {
			writeRPCInternalError(ctx, wrapError(ErrRPCInternal, "", err), false)
			return
		}
		analysisReportContent.UserProfile = userProfileModule
		// 测量上下文模块
		pulseTestModule, err := getPulseTestModule(respGetAnalyzeResult.GetReport().GetPulseTest())
		if err != nil {
			writeRPCInternalError(ctx, wrapError(ErrRPCInternal, "", err), false)
			return
		}
		// 备注模块
		analysisReportContent.PulseTest = pulseTestModule
		remarkModule, err := getRemarkModule(respGetAnalyzeResult.GetReport().GetRemark())
		if err != nil {
			writeRPCInternalError(ctx, wrapError(ErrRPCInternal, "", err), false)
			return
		}
		analysisReportContent.Remark = remarkModule
		// 测量时间
		startTime, _ := ptypes.Timestamp(respGetAnalyzeResult.GetReport().GetCreatedTime())
		analysisReportContent.CreatedTime = startTime

		analysisReportResponse.ReportContent = *analysisReportContent
	}

	rest.WriteOkJSON(ctx, analysisReportResponse)
}

func mapRestLanguageToProto(language Language) (ptypesv2.Language, error) {
	switch language {
	case LanguageSimpleChinese:
		return ptypesv2.Language_LANGUAGE_SIMPLIFIED_CHINESE, nil
	case LanguageTraditionalChinese:
		return ptypesv2.Language_LANGUAGE_TRADITIONAL_CHINESE, nil
	case LanguageEnglish:
		return ptypesv2.Language_LANGUAGE_ENGLISH, nil
	}
	return ptypesv2.Language_LANGUAGE_INVALID, fmt.Errorf("invalid language: %s", language)
}

// mapRestFingerToProto 将传入的格式为 int 的 finger 转化为proto 类型
func mapRestFingerToProto(finger int) (corepb.Finger, error) {
	switch finger {
	case FingerLeft1:
		return corepb.Finger_FINGER_LEFT_1, nil
	case FingerLeft2:
		return corepb.Finger_FINGER_LEFT_2, nil
	case FingerLeft3:
		return corepb.Finger_FINGER_LEFT_3, nil
	case FingerLeft4:
		return corepb.Finger_FINGER_LEFT_4, nil
	case FingerLeft5:
		return corepb.Finger_FINGER_LEFT_5, nil
	case FingerRight1:
		return corepb.Finger_FINGER_RIGHT_1, nil
	case FingerRight2:
		return corepb.Finger_FINGER_RIGHT_2, nil
	case FingerRight3:
		return corepb.Finger_FINGER_RIGHT_3, nil
	case FingerRight4:
		return corepb.Finger_FINGER_RIGHT_4, nil
	case FingerRight5:
		return corepb.Finger_FINGER_RIGHT_5, nil
	}
	return corepb.Finger_FINGER_INVALID, fmt.Errorf("invalid int32 finger %d", finger)
}

func mapProtoFingerToRest(protoFinger corepb.Finger) (int, error) {
	switch protoFinger {
	case corepb.Finger_FINGER_INVALID:
		return FingerInvalid, fmt.Errorf("invalid proto finger %d", protoFinger)
	case corepb.Finger_FINGER_UNSET:
		return FingerInvalid, fmt.Errorf("invalid proto finger %d", protoFinger)
	case corepb.Finger_FINGER_LEFT_1:
		return FingerLeft1, nil
	case corepb.Finger_FINGER_LEFT_2:
		return FingerLeft2, nil
	case corepb.Finger_FINGER_LEFT_3:
		return FingerLeft3, nil
	case corepb.Finger_FINGER_LEFT_4:
		return FingerLeft4, nil
	case corepb.Finger_FINGER_LEFT_5:
		return FingerLeft5, nil
	case corepb.Finger_FINGER_RIGHT_1:
		return FingerRight1, nil
	case corepb.Finger_FINGER_RIGHT_2:
		return FingerRight2, nil
	case corepb.Finger_FINGER_RIGHT_3:
		return FingerRight3, nil
	case corepb.Finger_FINGER_RIGHT_4:
		return FingerRight4, nil
	case corepb.Finger_FINGER_RIGHT_5:
		return FingerRight5, nil
	}
	return -1, fmt.Errorf("invalid proto finger %d", protoFinger)
}

func mapRestLanguageToAnalysisProto(protoFinger corepb.Finger) (int, error) {
	switch protoFinger {
	case corepb.Finger_FINGER_INVALID:
		return FingerInvalid, fmt.Errorf("invalid proto finger %d", protoFinger)
	case corepb.Finger_FINGER_UNSET:
		return FingerInvalid, fmt.Errorf("invalid proto finger %d", protoFinger)
	case corepb.Finger_FINGER_LEFT_1:
		return FingerLeft1, nil
	case corepb.Finger_FINGER_LEFT_2:
		return FingerLeft2, nil
	case corepb.Finger_FINGER_LEFT_3:
		return FingerLeft3, nil
	case corepb.Finger_FINGER_LEFT_4:
		return FingerLeft4, nil
	case corepb.Finger_FINGER_LEFT_5:
		return FingerLeft5, nil
	case corepb.Finger_FINGER_RIGHT_1:
		return FingerRight1, nil
	case corepb.Finger_FINGER_RIGHT_2:
		return FingerRight2, nil
	case corepb.Finger_FINGER_RIGHT_3:
		return FingerRight3, nil
	case corepb.Finger_FINGER_RIGHT_4:
		return FingerRight4, nil
	case corepb.Finger_FINGER_RIGHT_5:
		return FingerRight5, nil
	}
	return -1, fmt.Errorf("invalid proto finger %d", protoFinger)
}
