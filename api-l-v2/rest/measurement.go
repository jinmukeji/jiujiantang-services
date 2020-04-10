package rest

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	corepb "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	"github.com/kataras/iris/v12"
)

const (
	// Android android 手机
	Android = "ANDROID"
	// Iphone 苹果手机
	Iphone = "IPHONE"
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

// Measurement 测量
type Measurement struct {
	MeasurementData MeasurementData `json:"measurement"`
}

// MeasurementData 测量数据
type MeasurementData struct {
	UserID              int    `json:"user_id"`                // 用户ID
	Data0               string `json:"data0"`                  // 第一条数据
	Data1               string `json:"data1"`                  // 第二条数据
	Mac                 string `json:"mac"`                    // mac
	MobileType          string `json:"mobile_type"`            // 手机类型
	AppHeartRate        int    `json:"app_heart_rate"`         // 心率
	Finger              int    `json:"finger"`                 // 测量手指
	RecordType          int    `json:"record_type"`            // 记录类型
	AppHighestHeartRate int32  `json:"app_highest_heart_rate"` // app最高心率
	AppLowestHeartRate  int32  `json:"app_lowest_heart_rate"`  // app最低心率
	MeasurementPosture  int32  `json:"measurement_posture"`    // 测量姿态
}

// SubmitMeasurementData 提交测量数据
type SubmitMeasurementData struct {
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

func (h *v2Handler) SubmitMeasurementData(ctx iris.Context) {
	var measurement Measurement
	err := ctx.ReadJSON(&measurement)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}
	if measurement.MeasurementData.UserID == 0 {
		writeError(ctx, wrapError(ErrValueRequired, "", errors.New("missing user_id")), false)
		return
	}
	if measurement.MeasurementData.AppHeartRate < 0 {
		writeError(ctx, wrapError(ErrInvalidValue, "", fmt.Errorf("invalid app heart rate %d", measurement.MeasurementData.AppHeartRate)), false)
		return
	}

	if measurement.MeasurementData.Finger < 1 || measurement.MeasurementData.Finger > 10 {
		writeError(ctx, wrapError(ErrInvalidValue, "", fmt.Errorf("invalid finger %d", measurement.MeasurementData.Finger)), false)
		return
	}

	if measurement.MeasurementData.MobileType != Android && measurement.MeasurementData.MobileType != Iphone {
		writeError(ctx, wrapError(ErrInvalidValue, "", fmt.Errorf("mobile type must be ANDROID or IPHONE,current type is %s", measurement.MeasurementData.MobileType)), false)
		return
	}
	req := new(corepb.SubmitMeasurementInfoRequest)
	req.UserId = int32(measurement.MeasurementData.UserID)
	protoMobileType, errmapRestMobileTypeToProto := mapRestMobileTypeToProto(measurement.MeasurementData.MobileType)
	if errmapRestMobileTypeToProto != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errmapRestMobileTypeToProto), false)
		return
	}
	req.MobileType = protoMobileType
	req.AppHighestHr = measurement.MeasurementData.AppHighestHeartRate
	req.AppLowestHr = measurement.MeasurementData.AppLowestHeartRate
	protoMeasurementPosture, errmapRestMeasurementPostureToProto := mapRestMeasurementPostureToProto(measurement.MeasurementData.MeasurementPosture)
	if errmapRestMeasurementPostureToProto != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errmapRestMeasurementPostureToProto), false)
		return
	}
	req.MeasurementPosture = protoMeasurementPosture
	data0, _ := base64.StdEncoding.DecodeString(measurement.MeasurementData.Data0)
	req.Info0 = &corepb.BluetoothInfo{
		Ir5160: data0,
	}

	data1, _ := base64.StdEncoding.DecodeString(measurement.MeasurementData.Data1)
	req.Info1 = &corepb.BluetoothInfo{
		Ir5160: data1,
	}
	req.Mac = measurement.MeasurementData.Mac
	req.AppHr = int32(measurement.MeasurementData.AppHeartRate)
	protoFinger, errMapRestFingerToProto := mapRestFingerToProto(measurement.MeasurementData.Finger)
	if errMapRestFingerToProto != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errMapRestFingerToProto), false)
		return
	}
	req.Finger = protoFinger
	resp, errResp := h.rpcSvc.SubmitMeasurementInfo(
		newRPCContext(ctx), req,
	)
	if errResp != nil {
		writeRpcInternalError(ctx, errResp, false)
		return
	}

	rest.WriteOkJSON(ctx, SubmitMeasurementData{
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

func mapRestMobileTypeToProto(mobileType string) (corepb.MobileType, error) {
	switch strings.ToUpper(mobileType) {
	case Android:
		return corepb.MobileType_MOBILE_TYPE_ANDROID, nil
	case Iphone:
		return corepb.MobileType_MOBILE_TYPE_IPHONE, nil
	}
	return corepb.MobileType_MOBILE_TYPE_INVALID, fmt.Errorf("invalid string mobile type %s", mobileType)
}

// mapRestMeasurementPostureToProto 将rest曾传入的 int32 格式的 measurementPosture 转换为 proto 类型
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
