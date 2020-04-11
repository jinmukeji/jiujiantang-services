package handler

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/jiujiantang-services/analysis/aws"
	"github.com/jinmukeji/jiujiantang-services/pkg/rpc"
	"github.com/jinmukeji/jiujiantang-services/service/auth"
	"github.com/jinmukeji/jiujiantang-services/service/mysqldb"
	corepb "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	pulsetestinfopb "github.com/jinmukeji/proto/gen/micro/idl/jm/pulsetestinfo/v1"
	context "golang.org/x/net/context"
)

// maxRecords 测量历史记录显示数量上限
const maxRecords = 100

// minRecords 测量历史记录显示数量下限
const minRecords = 1

// Answer 答案
type Answer struct {
	QuestionKey string   `json:"question_key"`
	Values      []string `json:"values"`
}

// SearchHistory 查找测量历史记录
func (j *JinmuHealth) SearchHistory(ctx context.Context, req *corepb.SearchHistoryRequest, resp *corepb.SearchHistoryResponse) error {
	l := rpc.ContextLogger(ctx)

	accessTokenType, _ := auth.AccessTokenTypeFromContext(ctx)
	var records []mysqldb.Record
	// TODO: 消除重复代码
	if accessTokenType == AccessTokenTypeWeChatValue {
		wxUser, errFindWXUserByOpenID := j.datastore.FindWXUserByOpenID(ctx, req.OpenId)
		if errFindWXUserByOpenID != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to find wx user by openID %s: %s", req.OpenId, errFindWXUserByOpenID.Error()))
		}
		r, errFindValidPaginatedRecordsByUserID := j.datastore.FindValidPaginatedRecordsByUserID(ctx, int32(wxUser.UserID))
		if errFindValidPaginatedRecordsByUserID != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to find valid paginated records by user %d: %s", wxUser.UserID, errFindValidPaginatedRecordsByUserID.Error()))
		}
		records = r
	} else {
		_, err := validateSearchHistoryRequest(req)
		if err != nil {
			return NewError(ErrValidateSearchHistoryRequestFailure, err)
		}
		start, _ := ptypes.Timestamp(req.StartTime)
		end, _ := ptypes.Timestamp(req.EndTime)
		r, errFindValidPaginatedRecordsByDateRange := j.datastore.FindValidPaginatedRecordsByDateRange(ctx, int(req.UserId), int(req.Offset), int(req.Size), start, end)
		if errFindValidPaginatedRecordsByDateRange != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to find valid paginated records by user %d: %s", req.UserId, errFindValidPaginatedRecordsByDateRange.Error()))
		}
		records = r
	}

	var recordsReply []*corepb.RecordHistory

	for _, record := range records {
		//TODO:增加显示
		createAt, _ := ptypes.TimestampProto(record.CreatedAt)
		dataArray, errgetPulseTestDataIntArray := getPulseTestDataIntArray(record.S3Key, j.awsClient)
		if errgetPulseTestDataIntArray != nil {
			cause := fmt.Sprintf("failed to getPulseTestDataIntArray, record_id: %d, data length is %d", record.RecordID, len(dataArray))
			l.WithError(errgetPulseTestDataIntArray).Warn(cause)
			dataArray = make([]int32, 0)
		}
		waveData, errpickPartialPulseTestRawData := pickPartialPulseTestRawData(dataArray)
		if errpickPartialPulseTestRawData != nil {
			cause := fmt.Sprintf("failed to errpickPartialPulseTestRawData, record_id is %d, data length is %d", record.RecordID, len(waveData))
			l.WithError(errpickPartialPulseTestRawData).Warn(cause)
			waveData = make([]int32, 0)
		}
		var answers []Answer
		// record.Answers 是空字符串 ，不做处理，直接跳过返回默认值
		if record.Answers != "" {
			errUnmarshal := json.Unmarshal([]byte(record.Answers), &answers)
			if errUnmarshal != nil {
				return NewError(ErrJSONUnmarshalFailure, fmt.Errorf("failed to unmarshal answers %s: %s", []byte(record.Answers), errUnmarshal.Error()))
			}
		}
		tags := make([]*corepb.GeneralExplain, 0)
		for _, answer := range answers {
			for _, value := range answer.Values {
				choice := j.analysisEngine.QuestionDoc.GetQuestionChoice(answer.QuestionKey, value)
				if choice != nil && choice.Tag != nil {
					tags = append(tags, &corepb.GeneralExplain{
						Key:     choice.Tag.Key,
						Label:   choice.Tag.Label,
						Content: choice.Tag.Content,
					})
				}
			}
		}
		var stressState = make(map[string]bool)
		if record.StressState != "" {
			errUnmarshal := json.Unmarshal([]byte(record.StressState), &stressState)
			if errUnmarshal != nil {
				continue
			}
		}
		body := record.AnalyzeBody
		var physicalDialectics []string
		if body != "" {
			var analysisReportRequestBody AnalysisReportRequestBody
			errUnmarshal := json.Unmarshal([]byte(body), &analysisReportRequestBody)
			if errUnmarshal != nil {
				continue
			}
			physicalDialectics = make([]string, len(analysisReportRequestBody.PhysicalDialectics))
			for idx, pd := range analysisReportRequestBody.PhysicalDialectics {
				physicalDialectics[idx] = pd.Key
			}
		}
		protoFinger, errMapDBFingerToProto := mapDBFingerToProto(record.Finger)
		if errMapDBFingerToProto != nil {
			return NewError(ErrInvalidFinger, errMapDBFingerToProto)
		}
		recordsReply = append(recordsReply, &corepb.RecordHistory{
			C0:                 Int32ValBoundedBy10FromFloat(record.C0),
			C1:                 Int32ValBoundedBy10FromFloat(record.C1),
			C2:                 Int32ValBoundedBy10FromFloat(record.C2),
			C3:                 Int32ValBoundedBy10FromFloat(record.C3),
			C4:                 Int32ValBoundedBy10FromFloat(record.C4),
			C5:                 Int32ValBoundedBy10FromFloat(record.C5),
			C6:                 Int32ValBoundedBy10FromFloat(record.C6),
			C7:                 Int32ValBoundedBy10FromFloat(record.C7),
			G0:                 record.G0,
			G1:                 record.G1,
			G2:                 record.G2,
			G3:                 record.G3,
			G4:                 record.G4,
			G5:                 record.G5,
			G6:                 record.G6,
			G7:                 record.G7,
			RecordId:           int32(record.RecordID),
			Finger:             protoFinger,
			Info:               waveData,
			CreatedTime:        createAt,
			Hr:                 record.HeartRate,
			AppHr:              record.AppHeartRate,
			RecordType:         int32(record.RecordType),
			Remark:             record.Remark,
			IsSportOrDrunk:     corepb.Status(record.IsSportOrDrunk),
			Cold:               corepb.Status(record.Cold),
			MenstrualCycle:     corepb.Status(record.MenstrualCycle),
			OvipositPeriod:     corepb.Status(record.OvipositPeriod),
			Lactation:          corepb.Status(record.Lactation),
			Pregnancy:          corepb.Status(record.Pregnancy),
			CmStatusA:          corepb.Status(record.StatusA),
			CmStatusB:          corepb.Status(record.StatusB),
			CmStatusC:          corepb.Status(record.StatusC),
			CmStatusD:          corepb.Status(record.StatusD),
			CmStatusE:          corepb.Status(record.StatusE),
			CmStatusF:          corepb.Status(record.StatusF),
			Tags:               tags,
			HasPaid:            record.HasPaid,
			ShowFullReport:     record.ShowFullReport,
			StressStatus:       stressState,
			HasStressState:     record.HasStressState,
			PhysicalDialectics: physicalDialectics,
		})
	}
	resp.RecordHistories = recordsReply
	return nil
}

// GetMeasurementRecord 获取测量记录
func (j *JinmuHealth) GetMeasurementRecord(ctx context.Context, req *corepb.GetMeasurementRecordRequest, resp *corepb.GetMeasurementRecordResponse) error {
	recordID := req.RecordId
	isExsit, err := j.datastore.ExistRecordByRecordID(ctx, recordID)
	if err != nil || !isExsit {
		return NewError(ErrGetRecordFailure, fmt.Errorf("failed to check record existence by recordID %d: %s", recordID, err.Error()))
	}
	accessTokenType, _ := auth.AccessTokenTypeFromContext(ctx)
	// TODO: 以后修改
	if accessTokenType != AccessTokenTypeLValue && accessTokenType != AccessTokenTypeWeChatValue {
		userID, errGetUserIDByRecordID := j.datastore.GetUserIDByRecordID(ctx, recordID)
		if errGetUserIDByRecordID != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to get userID by record %d: %s", recordID, errGetUserIDByRecordID.Error()))
		}
		ownerID, ok := auth.UserIDFromContext(ctx)
		if !ok {
			return NewError(ErrInvalidUser, errors.New("failed to get userID from context"))
		}
		userOrganization, _ := j.datastore.FindOrganizationByUserID(ctx, int(userID))
		ownerOrganization, _ := j.datastore.FindOrganizationByUserID(ctx, int(ownerID))
		if ownerOrganization.OrganizationID != userOrganization.OrganizationID {
			return NewError(ErrInvalidUser, fmt.Errorf("this user_id %d cannot get measurement by record_id %d", userID, recordID))
		}
	}
	record, errGetRecordByRecordID := j.datastore.FindRecordByID(ctx, int(recordID))
	if errGetRecordByRecordID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find record by recordID %d: %s", recordID, errGetRecordByRecordID.Error()))
	}
	resp.AppHr = record.AppHeartRate
	resp.Hr = int32(record.HeartRate)
	resp.C0 = Int32ValBoundedBy10FromFloat(record.C0)
	resp.C1 = Int32ValBoundedBy10FromFloat(record.C1)
	resp.C2 = Int32ValBoundedBy10FromFloat(record.C2)
	resp.C3 = Int32ValBoundedBy10FromFloat(record.C3)
	resp.C4 = Int32ValBoundedBy10FromFloat(record.C4)
	resp.C5 = Int32ValBoundedBy10FromFloat(record.C5)
	resp.C6 = Int32ValBoundedBy10FromFloat(record.C6)
	resp.C7 = Int32ValBoundedBy10FromFloat(record.C7)
	protoFinger, errMapDBFingerToProto := mapDBFingerToProto(record.Finger)
	if errMapDBFingerToProto != nil {
		return NewError(ErrInvalidFinger, errMapDBFingerToProto)
	}
	resp.Finger = protoFinger
	dataArray, _ := getPulseTestDataIntArray(record.S3Key, j.awsClient)
	if len(dataArray) >= 4000 {
		resp.Info = dataArray[3000:4000]
	}
	resp.Answers = record.Answers
	return nil
}

// validatSearchHistoryRequest 验证请求数据
func validateSearchHistoryRequest(req *corepb.SearchHistoryRequest) (bool, error) {
	if _, err := ptypes.Timestamp(req.StartTime); err != nil {
		return false, fmt.Errorf("failed to parse timestamp of start time %s: %s", req.StartTime, err.Error())
	}
	if _, err := ptypes.Timestamp(req.EndTime); err != nil {
		return false, fmt.Errorf("failed to parse timestamp of end time %s: %s", req.EndTime, err.Error())
	}
	if req.Size != -1 && (req.Size > maxRecords || req.Size < minRecords) {
		return false, fmt.Errorf("size %d exceeds the maximum or minimum limit", req.Size)
	}
	return true, nil
}

// pickPartialPulseTestRawData 从S3上得到特殊的波形数据
func pickPartialPulseTestRawData(dataArray []int32) ([]int32, error) {
	if len(dataArray) < appWaveDataStart+appWaveDataLength {
		return nil, fmt.Errorf("wave data length %d too short", len(dataArray))
	}
	return dataArray[appWaveDataStart : appWaveDataStart+appWaveDataLength], nil
}

// getPulseTestDataIntArray 从aws上得到波形数据
func getPulseTestDataIntArray(s3Key string, client aws.PulseTestRawDataS3Client) ([]int32, error) {
	pulseTestRawData, err := client.Download(s3Key)
	if err != nil {
		return []int32{}, fmt.Errorf("failed to download raw data of s3key %s: %s", s3Key, err.Error())
	}
	return ParsePayload(pulseTestRawData)
}

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

// DeleteRecord 删除记录
func (j *JinmuHealth) DeleteRecord(ctx context.Context, req *corepb.DeleteRecordRequest, resp *corepb.DeleteRecordResponse) error {
	userID, _ := j.datastore.GetUserIDByRecordID(ctx, req.RecordId)
	if userID != req.UserId {
		return NewError(ErrRecordNotBelongToUser, fmt.Errorf("failed to get user by record %d", req.RecordId))
	}
	return j.datastore.DeleteRecord(ctx, req.RecordId)
}
