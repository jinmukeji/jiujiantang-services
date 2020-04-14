package rest

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	corepb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	"github.com/kataras/iris/v12"
)

const (
	// Android android 手机
	Android = "ANDROID"
	// Iphone 苹果手机
	Iphone = "IPHONE"
)
const (
	// FingerLeft1 左小拇指
	FingerLeft1 int = 1
	// FingerLeft2 左无名指
	FingerLeft2 int = 2
	// FingerLeft3 左中指
	FingerLeft3 int = 3
	// FingerLeft4 左食指
	FingerLeft4 int = 4
	// FingerLeft5 左大拇指
	FingerLeft5 int = 5
	// FingerRight5 右大拇指
	FingerRight5 int = 6
	// FingerRight4 右食指
	FingerRight4 int = 7
	// FingerRight3 右中指
	FingerRight3 int = 8
	// FingerRight2 右无名指
	FingerRight2 int = 9
	// FingerRight1 右小拇指
	FingerRight1 int = 10
	// FingerInvalid  非法的手指
	FingerInvalid int = -1
)
const (
	// MeasurementPostureSetting 坐姿
	MeasurementPostureSetting = 0
	// MeasurementPostureStanging 站姿
	MeasurementPostureStanging = 1
	// MeasurementPostureLying 躺姿
	MeasurementPostureLying = 2
)

const (
	// SunmitMeasurementPayloadCodec 提交测量的 Payload 解码方式
	SunmitMeasurementPayloadCodec = "ring-red-raw"
)

// SubmitMeasurementDataResponse 提交测量数据的响应
type SubmitMeasurementDataResponse struct {
	RecordID            int32     `json:"record_id"`
	Cid                 int32     `json:"cid"`
	C0                  int32     `json:"c0"`
	C1                  int32     `json:"c1"`
	C2                  int32     `json:"c2"`
	C3                  int32     `json:"c3"`
	C4                  int32     `json:"c4"`
	C5                  int32     `json:"c5"`
	C6                  int32     `json:"c6"`
	C7                  int32     `json:"c7"`
	WaveData            []int32   `json:"wave_data"`
	AppHeartRate        int32     `json:"app_heart_rate"`
	CreatedAt           time.Time `json:"created_at"`
	Finger              int32     `json:"finger"`
	RecordType          int32     `json:"record_type"`
	HeartRate           int32     `json:"heart_rate"`
	AppHighestHeartRate int32     `json:"app_highest_heart_rate"` // app最高心率
	AppLowestHeartRate  int32     `json:"app_lowest_heart_rate"`  // app最低心率
}

// 指环测量的采样数据
type RingSamplePayload struct {
	Fps               int32     `json:"fps"`    // 实际指环采样工作时的 FPS
	Finger            int       `json:"finger"` // 采样使用的手指
	Count             int32     `json:"count"`  // 测量采样到的数据点的数量
	PointSize         int32     `json:"point_size"`
	Payload           string    `json:"payload"`             // 测量采样数据点的数据，字节流形式。
	SamplingStartTime time.Time `json:"sampling_start_time"` // 采样开始时间
	SamplingStopTime  time.Time `json:"sampling_stop_time"`  // 采样结束时间
}

// 提交测量数据请求
type SubmitMeasurementInfoRequest struct {
	UserId             int32              `json:"user_id"`             // UserID
	Mac                string             `json:"mac"`                 // mac地址
	MobileType         string             `json:"mobile_type"`         // 手机类型
	MeasurementPosture int32              `json:"measurement_posture"` // 测量姿态
	Payload            *RingSamplePayload `json:"payload"`             // 采样数据
	Extras             map[string]string  `json:"extras"`              // 额外的扩展上下文数据，KV 键值对
}

func (h *v2Handler) SubmitMeasurementData(ctx iris.Context) {

	var submitMeasurementInfoRequest SubmitMeasurementInfoRequest
	err := ctx.ReadJSON(&submitMeasurementInfoRequest)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}
	if submitMeasurementInfoRequest.UserId == 0 {
		writeError(
			ctx,
			wrapError(ErrValueRequired, "", fmt.Errorf("invalid userID %d", submitMeasurementInfoRequest.UserId)),
			false,
		)
		return
	}
	ringSamplePayload := *(submitMeasurementInfoRequest.Payload)
	if ringSamplePayload.Finger < 1 || ringSamplePayload.Finger > 10 {
		writeError(
			ctx,
			wrapError(ErrInvalidValue, "", fmt.Errorf("invalid finger %d", ringSamplePayload.Finger)),
			false,
		)
		return
	}

	if submitMeasurementInfoRequest.MobileType != Android && submitMeasurementInfoRequest.MobileType != Iphone {
		writeError(
			ctx,
			wrapError(ErrInvalidValue, "", fmt.Errorf("mobile type must be ANDROID or IPHONE,current type is %s", submitMeasurementInfoRequest.MobileType)),
			false,
		)
		return
	}

	if submitMeasurementInfoRequest.MeasurementPosture < 0 || submitMeasurementInfoRequest.MeasurementPosture > 2 {
		writeError(
			ctx,
			wrapError(ErrInvalidValue, "", fmt.Errorf("invalid measurement posture %d", submitMeasurementInfoRequest.MeasurementPosture)),
			false,
		)
		return
	}

	req := new(corepb.SubmitMeasurementInfoRequest)
	protoMobileType, errmapRestMobiletypeToProto := mapRestMobiletypeToProto(submitMeasurementInfoRequest.MobileType)
	if errmapRestMobiletypeToProto != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errmapRestMobiletypeToProto), false)
		return
	}
	protoMeasurementPosture, errmapRestMeasurementPostureToProto := mapRestMeasurementPostureToProto(submitMeasurementInfoRequest.MeasurementPosture)
	if errmapRestMeasurementPostureToProto != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errmapRestMeasurementPostureToProto), false)
		return
	}
	req.UserId = int32(submitMeasurementInfoRequest.UserId)
	req.Mac = submitMeasurementInfoRequest.Mac
	req.MobileType = protoMobileType
	req.MeasurementPosture = protoMeasurementPosture
	payload, _ := base64.StdEncoding.DecodeString(ringSamplePayload.Payload)
	protoFinger, errMapProtoFingerToProto := mapRestFingerToProto(ringSamplePayload.Finger)
	if errMapProtoFingerToProto != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errmapRestMobiletypeToProto), false)
		return
	}
	samplingStartTime, err := ptypes.TimestampProto(ringSamplePayload.SamplingStartTime)
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	samplingStopTime, err := ptypes.TimestampProto(ringSamplePayload.SamplingStopTime)
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	req.Payload = &corepb.RingSamplePayload{
		Fps:               uint32(ringSamplePayload.Fps),
		Finger:            protoFinger,
		Count:             uint32(ringSamplePayload.Count),
		PointSize:         uint32(ringSamplePayload.PointSize),
		Payload:           payload,
		PayloadCodec:      SunmitMeasurementPayloadCodec,
		SamplingStartTime: samplingStartTime,
		SamplingStopTime:  samplingStopTime,
	}
	req.Extras = submitMeasurementInfoRequest.Extras
	resp, errResp := h.rpcSvc.SubmitMeasurementInfo(
		newRPCContext(ctx), req,
	)
	if errResp != nil {
		writeRPCInternalError(ctx, errResp, false)
		return
	}

	rest.WriteOkJSON(ctx, SubmitMeasurementDataResponse{
		RecordID:            resp.RecordId,
		Cid:                 resp.RecordId,
		C0:                  resp.C0,
		C1:                  resp.C1,
		C2:                  resp.C2,
		C3:                  resp.C3,
		C4:                  resp.C4,
		C5:                  resp.C5,
		C6:                  resp.C6,
		C7:                  resp.C7,
		WaveData:            resp.PartialInfo,
		AppHeartRate:        resp.AppHr,
		CreatedAt:           time.Now().UTC(),
		Finger:              int32(resp.Finger),
		RecordType:          resp.RecordType,
		HeartRate:           resp.Hr,
		AppHighestHeartRate: resp.AppHighestHr,
		AppLowestHeartRate:  resp.AppLowestHr,
	})
}

func mapRestMobiletypeToProto(mobileType string) (corepb.MobileType, error) {
	switch strings.ToUpper(mobileType) {
	case Android:
		return corepb.MobileType_MOBILE_TYPE_ANDROID, nil
	case Iphone:
		return corepb.MobileType_MOBILE_TYPE_IPHONE, nil
	}
	return corepb.MobileType_MOBILE_TYPE_INVALID, fmt.Errorf("invalid string mobile type %s", mobileType)
}

func mapRestMeasurementPostureToProto(measurementPosture int32) (corepb.MeasurementPosture, error) {
	switch measurementPosture {
	case MeasurementPostureSetting:
		return corepb.MeasurementPosture_MEASUREMENT_POSTURE_SITTING, nil
	case MeasurementPostureStanging:
		return corepb.MeasurementPosture_MEASUREMENT_POSTURE_STANDING, nil
	case MeasurementPostureLying:
		return corepb.MeasurementPosture_MEASUREMENT_POSTURE_LYING, nil
	}
	return -1, fmt.Errorf("invalid measurement posture %d", measurementPosture)
}
