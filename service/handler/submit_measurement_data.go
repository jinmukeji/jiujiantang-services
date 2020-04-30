package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/go-pkg/v2/age"
	idgen "github.com/jinmukeji/go-pkg/v2/id-gen"
	"github.com/jinmukeji/go-pkg/v2/mac"
	"github.com/jinmukeji/jiujiantang-services/ptcodec"

	"github.com/jinmukeji/go-pkg/v2/mathutil"
	"github.com/jinmukeji/jiujiantang-services/service/auth"
	"github.com/jinmukeji/jiujiantang-services/service/mysqldb"
	corepb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	pulsetestinfopb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/pulsetestinfo/v1"
	subscriptionpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/subscription/v1"
	jinmuidpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	calcpb "github.com/jinmukeji/proto/v3/gen/micro/idl/platform/calc/v2"
	generalpb "github.com/jinmukeji/proto/v3/gen/micro/idl/ptypes/v2"
	microclient "github.com/micro/go-micro/v2/client"
)

const (
	// 当前记录类型默认值为 9，表示第三代数据结构规范
	CurrentRecordType = 9
	// 当前的测试数据版本是2
	CurrentPulseTestDataSpec = 2
	//  左手指
	LeftFinger = "L"
	//  右手指
	RightFinger = "R"
)

const (
	// DeviceModelForAlgorithm 指环使用的算法服务器请求的 Model 值
	RingDeviceModelForAlgorithm = "FYI"
	// 提交给算法服务器请求的 sample 中的 Type，为红外光数据
	algorithmSampleType = "RING"
)

const (
	// 最高最低心率修正的阈值
	HeartRateThreshold = 0.02
)

// SubmitMeasurementInfo 用户提交测量数据
func (j *JinmuHealth) SubmitMeasurementInfo(ctx context.Context, req *corepb.SubmitMeasurementInfoRequest, resp *corepb.SubmitMeasurementInfoResponse) error {

	accessTokenType, _ := auth.AccessTokenTypeFromContext(ctx)
	ownerID, _ := auth.UserIDFromContext(ctx)
	if accessTokenType != AccessTokenTypeLValue {
		// 不是一体机查订阅
		reqSelectedGetUserSubscription := new(subscriptionpb.GetSelectedUserSubscriptionRequest)
		reqSelectedGetUserSubscription.UserId = req.UserId
		reqSelectedGetUserSubscription.OwnerId = ownerID
		respSelectedGetUserSubscription, errGetSelectedUserSubscription := j.subscriptionSvc.GetSelectedUserSubscription(ctx, reqSelectedGetUserSubscription)
		// 获取不到订阅或订阅是空，则提示订阅过期
		if errGetSelectedUserSubscription != nil || respSelectedGetUserSubscription.Subscription == nil {
			return NewError(ErrSubscriptionExpired, fmt.Errorf("subscription of user %d is not found", req.UserId))
		}
		// 订阅过期
		expiredAt, _ := ptypes.Timestamp(respSelectedGetUserSubscription.Subscription.ExpiredTime)
		if expiredAt.Before(time.Now()) {
			return NewError(ErrSubscriptionExpired, errors.New("subscription is expired"))
		}
		IsSameOrganization, _ := j.datastore.CheckUserOwnerBelongToSameOrganization(ctx, req.UserId, ownerID)
		if !IsSameOrganization {
			return NewError(ErrInvalidUser, fmt.Errorf("user %d and owner %d do not belong to the same organization", req.UserId, ownerID))
		}
	}
	var err error

	// 验证 Request
	if err = j.validateMeasurementDataRequest(ctx, req); err != nil {
		return err
	}

	// 获取当前提交测量数据的用户
	reqGetUserProfile := new(jinmuidpb.GetUserProfileRequest)
	reqGetUserProfile.UserId = req.UserId
	reqGetUserProfile.IsSkipVerifyToken = true
	respGetUserProfile, errGetUserProfile := j.jinmuidSvc.GetUserProfile(ctx, reqGetUserProfile)
	if errGetUserProfile != nil {
		return NewError(ErrInvalidUser, fmt.Errorf("failed to get user profile by user %d", req.UserId))
	}

	// 生成算法服务器请求数据并请求运算
	algRequest, err := buildAlgorithmRequest(req.UserId, respGetUserProfile.Profile, req)
	if err != nil {
		return NewError(ErrBuildAlgorithmRequestFailure, fmt.Errorf("Fail to build request for algorithm: %s", err.Error()))
	}

	algResp, err := j.algorithmClient.SubmitCalc(ctx, algRequest, microclient.WithAddress(j.algorithmServerAddress))
	if err != nil {
		return NewError(ErrInvokeAlgorithmServerFailure, fmt.Errorf("Failed to invoke algorithm server : %s", err.Error()))
	}
	// 构建数据库存储的记录
	record, errBuildDbRecord := buildDbRecord(ctx, req, algResp)
	if errBuildDbRecord != nil {
		return NewError(ErrCreateRecordFailure, errBuildDbRecord)
	}
	if record.HeartRate == -1 {
		// 当算法服务器 HeartRate 返回结果是 -1，则测量数据不正确
		return NewError(ErrDatabase, fmt.Errorf("heart_rate is %f", record.HeartRate))
	}
	if err = j.datastore.CreateRecord(ctx, record); err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to create record: %s", err.Error()))
	}
	if accessTokenType == AccessTokenTypeLValue {
		token, _ := auth.TokenFromContext(ctx)
		account, err := j.datastore.FindJinmuLAccountByToken(ctx, token)
		if err != nil || account == "" {
			return NewError(ErrUserUnauthorized, errors.New("登录授权已失效，请重新登录"))
		}
		errCreateAccountLRecord := j.datastore.CreateAccountLRecord(ctx, account, int32(record.RecordID))
		if errCreateAccountLRecord != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to create accountL record: %d: %s", record.RecordID, errCreateAccountLRecord.Error()))
		}
	}
	// 将原始数据保存到 s3
	pulseTestRawData := getSavingData(req.Payload.Payload, int32(record.RecordID))

	dataSendToApp, err := getPartialeWaveData(req.GetPayload().GetPayload())
	if err != nil {
		return NewError(ErrSetWavedataFailure, errors.New("Fail to get wavedata"))
	}
	if err := setSubmitMeasurementReply(req, resp, algResp, dataSendToApp); err != nil {
		return NewError(ErrSetWavedataFailure, errors.New("Fail to set wavedata send to app"))
	}

	output, err := j.awsClient.Upload(pulseTestRawData)
	if err != nil || output == nil {
		return NewError(ErrUploadWavedataToAWSFailure, errors.New("Fail to upload wavedata"))
	}
	resp.RecordId = int32(record.RecordID)
	resp.CreatedTime, _ = ptypes.TimestampProto(record.CreatedAt)
	resp.Finger = req.Payload.GetFinger()
	resp.Hr = int32(record.HeartRate)
	resp.RecordType = int32(record.RecordType)
	resp.AppHighestHr = int32(algResp.GetResult().GetHighestHeartRate())
	resp.AppLowestHr = int32(algResp.GetResult().GetLowestHeartRate())
	return nil
}

func mapDeviceModelToCodec(model string) (string, error) {
	switch model {
	case RingDeviceModelForAlgorithm:
		return ptcodec.CodecRingRedRaw, nil
	}
	return "", fmt.Errorf("invalid device model type [%s]", model)
}

// validateMeasurementData 验证 app 提交的数据
func (j *JinmuHealth) validateMeasurementDataRequest(ctx context.Context, req *corepb.SubmitMeasurementInfoRequest) error {

	// 验证 mac 地址非空
	if valid.IsNull(req.Mac) {
		return NewError(ErrInvalidDevice, errors.New("mac address is required"))
	}
	// 验证测量姿势是否合法
	_, err := mapProtoPostureToDB(req.GetMeasurementPosture())
	if err != nil {
		return fmt.Errorf("invalid posture: %w", err)
	}
	// 验证手机类型
	if req.MobileType == corepb.MobileType_MOBILE_TYPE_INVALID || req.MobileType == corepb.MobileType_MOBILE_TYPE_UNSET {
		return fmt.Errorf("invalid mobile type")
	}
	payload := req.GetPayload()
	if payload == nil {
		return fmt.Errorf("payload should not be null")
	}
	// 验证数据点是否正确
	if len(payload.GetPayload()) != int(payload.GetCount()*payload.GetPointSize()) {
		return fmt.Errorf("size of payload %d is invalid", len(payload.GetPayload()))
	}

	// 验证采样开始和结束时间是否合法
	if payload.GetSamplingStartTime() == nil {
		return fmt.Errorf("sampling_start_time of payload is null")
	}
	if payload.GetSamplingStopTime() == nil {
		return fmt.Errorf("sampling_stop_time of payload is null")
	}
	payload1StartTime, err := ptypes.Timestamp(payload.GetSamplingStartTime())
	if err != nil {
		return fmt.Errorf("failed to parse payload_1.sampling_start_time: %w", err)
	}
	payload1StopTime, err := ptypes.Timestamp(payload.GetSamplingStopTime())
	if err != nil {
		return fmt.Errorf("failed to parse payload_1.sampling_stop_time: %w", err)
	}

	if payload1StartTime.After(payload1StopTime) {
		return fmt.Errorf("sampling_start_time of payload_1 should be earlier than sampling_stop_time")
	}
	return nil
}

// buildTxId 生成流水号
func buildTxId() string {
	return idgen.NewXid()
}

func mapProtoGenderToCalc(gender generalpb.Gender) (calcpb.Gender, error) {
	switch gender {
	case generalpb.Gender_GENDER_FEMALE:
		return calcpb.Gender_GENDER_FEMALE, nil
	case generalpb.Gender_GENDER_MALE:
		return calcpb.Gender_GENDER_MALE, nil
	case generalpb.Gender_GENDER_INVALID:
		return calcpb.Gender_GENDER_INVALID, fmt.Errorf("invalid proto gender %d", generalpb.Gender_GENDER_INVALID)
	case generalpb.Gender_GENDER_UNSET:
		return calcpb.Gender_GENDER_UNSET, fmt.Errorf("invalid proto gender %d", generalpb.Gender_GENDER_UNSET)
	}
	return calcpb.Gender_GENDER_INVALID, errors.New("invalid proto gender")
}

func mapCoreFingerToCalc(finger corepb.Finger) (calcpb.Finger, error) {
	switch finger {
	case corepb.Finger_FINGER_LEFT_1:
		return calcpb.Finger_FINGER_LEFT_1, nil
	case corepb.Finger_FINGER_LEFT_2:
		return calcpb.Finger_FINGER_LEFT_2, nil
	case corepb.Finger_FINGER_LEFT_3:
		return calcpb.Finger_FINGER_LEFT_3, nil
	case corepb.Finger_FINGER_LEFT_4:
		return calcpb.Finger_FINGER_LEFT_4, nil
	case corepb.Finger_FINGER_LEFT_5:
		return calcpb.Finger_FINGER_LEFT_5, nil
	case corepb.Finger_FINGER_RIGHT_1:
		return calcpb.Finger_FINGER_RIGHT_1, nil
	case corepb.Finger_FINGER_RIGHT_2:
		return calcpb.Finger_FINGER_RIGHT_1, nil
	case corepb.Finger_FINGER_RIGHT_3:
		return calcpb.Finger_FINGER_RIGHT_3, nil
	case corepb.Finger_FINGER_RIGHT_4:
		return calcpb.Finger_FINGER_RIGHT_4, nil
	case corepb.Finger_FINGER_RIGHT_5:
		return calcpb.Finger_FINGER_RIGHT_5, nil
	}
	return calcpb.Finger_FINGER_INVALID, fmt.Errorf("invalid finger [%d]", finger)
}

// buildAlgorithmRequest 从 app 提交的数据生成算法服务器请求数据
func buildAlgorithmRequest(userID int32, userProfile *jinmuidpb.UserProfile, req *corepb.SubmitMeasurementInfoRequest) (*calcpb.SubmitCalcRequest, error) {
	txId := buildTxId()

	gender, err := mapProtoGenderToCalc(userProfile.GetGender())
	if err != nil {
		return nil, err
	}

	birthday, err := ptypes.Timestamp(userProfile.GetBirthdayTime())
	if err != nil {
		return nil, err
	}

	profile := &calcpb.SubjectProfile{
		Gender: gender,
		Height: userProfile.GetHeight(),
		Weight: userProfile.GetWeight(),
		Age:    int32(age.Age(birthday.UTC())),
	}

	device := &calcpb.PulseTestDevice{
		// DeviceId 传空
		DeviceId: "",
		//  十六进制,返回标准的 MAC 地址
		Mac:   mac.NormalizeMac(req.Mac),
		Model: RingDeviceModelForAlgorithm,
	}

	// 获得解码后的原始数据的数组
	decoder := ptcodec.NewDecoder(req.Payload.GetPayloadCodec())
	if decoder == nil {
		return nil, fmt.Errorf("failed to get decoder of payload codec [%s]", req.Payload.GetPayloadCodec())
	}
	// 请求算法服务器用第二段蓝牙数据计算
	payload := req.GetPayload().GetPayload()

	payloadArray, err := decoder.Decode(payload)
	if err != nil {
		return nil, errors.New("failed to decode payload of payload codec")
	}

	finger, err := mapCoreFingerToCalc(req.GetPayload().GetFinger())
	if err != nil {
		return nil, fmt.Errorf("failed to map core finger to pulse test: %s", err.Error())
	}
	now := time.Now()
	start, _ := ptypes.TimestampProto(now.Add(-1 * time.Minute))
	end, _ := ptypes.TimestampProto(now)
	samplePayload := &calcpb.SamplePayload{
		Type:              algorithmSampleType,
		Finger:            finger,
		SampleRate:        req.Payload.GetFps(),
		Payload:           payloadArray,
		Params:            req.Payload.GetSampleContext(),
		SamplingStartTime: start,
		SamplingStopTime:  end,
	}

	algorithmRequest := &calcpb.SubmitCalcRequest{
		TxId:    txId,
		Subject: profile,
		Device:  device,
		Sample:  samplePayload,
	}

	return algorithmRequest, nil
}

// generateAlgCid 生成算法服务调用的 cid
func generateAlgCid() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// buildDbRecord 从算法服务器返回和 app 提交数据生成 record
func buildDbRecord(ctx context.Context, req *corepb.SubmitMeasurementInfoRequest, resp *calcpb.SubmitCalcResponse) (*mysqldb.Record, error) {

	now := time.Now()

	client, _ := clientFromContext(ctx)
	clientID := client.ClientID

	// 设置 C0-C7 G0-G7 C0CV-C7CV
	C0 := float64(resp.GetResult().GetC0())
	C1 := float64(resp.GetResult().GetC1())
	C2 := float64(resp.GetResult().GetC2())
	C3 := float64(resp.GetResult().GetC3())
	C4 := float64(resp.GetResult().GetC4())
	C5 := float64(resp.GetResult().GetC5())
	C6 := float64(resp.GetResult().GetC6())
	C7 := float64(resp.GetResult().GetC7())
	G0 := resp.GetResult().GetG0()
	G1 := resp.GetResult().GetG1()
	G2 := resp.GetResult().GetG2()
	G3 := resp.GetResult().GetG3()
	G4 := resp.GetResult().GetG4()
	G5 := resp.GetResult().GetG5()
	G6 := resp.GetResult().GetG6()
	G7 := resp.GetResult().GetG7()
	C0CV := float64(resp.GetResult().GetC0Cv())
	C1CV := float64(resp.GetResult().GetC1Cv())
	C2CV := float64(resp.GetResult().GetC2Cv())
	C3CV := float64(resp.GetResult().GetC3Cv())
	C4CV := float64(resp.GetResult().GetC4Cv())
	C5CV := float64(resp.GetResult().GetC5Cv())
	C6CV := float64(resp.GetResult().GetC6Cv())
	C7CV := float64(resp.GetResult().GetC7Cv())
	heartRate := req.GetAppHeartRate()
	// 根据阈值处理最高最低心率
	// 小数直接取整
	lowestHr := int32(float64(heartRate) * (1 - HeartRateThreshold))
	highestHr := int32(float64(heartRate) * (1 + HeartRateThreshold))
	measurementPosture, errmapProtoPostureToDB := mapProtoPostureToDB(req.MeasurementPosture)
	if errmapProtoPostureToDB != nil {
		return nil, errmapProtoPostureToDB
	}
	mysqlFinger, errMapProtoFingerToDB := mapProtoFingerToDB(req.Payload.GetFinger())
	if errMapProtoFingerToDB != nil {
		return nil, errMapProtoFingerToDB
	}
	record := &mysqldb.Record{
		AppHeartRate:              float64(heartRate),
		HeartRate:                 float64(heartRate),
		HeartRateCV:               float32(resp.GetResult().GetHeartRateCv()),
		AlgorithmHighestHeartRate: highestHr,
		AlgorithmLowestHeartRate:  lowestHr,
		ClientID:                  clientID,
		UserID:                    int(req.UserId),
		DeviceID:                  0,
		RecordType:                CurrentRecordType,

		Finger: mysqlFinger,

		C0:   C0,
		C1:   C1,
		C2:   C2,
		C3:   C3,
		C4:   C4,
		C5:   C5,
		C6:   C6,
		C7:   C7,
		G0:   int32(G0),
		G1:   int32(G1),
		G2:   int32(G2),
		G3:   int32(G3),
		G4:   int32(G4),
		G5:   int32(G5),
		G6:   int32(G6),
		G7:   int32(G7),
		C0CV: C0CV,
		C1CV: C1CV,
		C2CV: C2CV,
		C3CV: C3CV,
		C4CV: C4CV,
		C5CV: C5CV,
		C6CV: C6CV,
		C7CV: C7CV,

		IsValid:            mysqldb.DbValidValue,
		SNR:                float32(resp.GetResult().GetSnr()),
		MeasurementPosture: measurementPosture,
		AnalyzeStatus:      mysqldb.AnalysisStatusPending, // pending
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	return record, nil
}

func mapProtoPostureToDB(measurementPosture corepb.MeasurementPosture) (mysqldb.MeasurementPosture, error) {
	switch measurementPosture {
	case corepb.MeasurementPosture_MEASUREMENT_POSTURE_INVALID:
		return mysqldb.MeasurementPostureInvalid, fmt.Errorf("invalid proto posture %d", corepb.MeasurementPosture_MEASUREMENT_POSTURE_INVALID)
	case corepb.MeasurementPosture_MEASUREMENT_POSTURE_UNSET:
		return mysqldb.MeasurementPostureInvalid, fmt.Errorf("invalid proto posture %d", corepb.MeasurementPosture_MEASUREMENT_POSTURE_UNSET)
	case corepb.MeasurementPosture_MEASUREMENT_POSTURE_SITTING:
		return mysqldb.MeasurementPostureSetting, nil
	case corepb.MeasurementPosture_MEASUREMENT_POSTURE_STANDING:
		return mysqldb.MeasurementPostureStanging, nil
	case corepb.MeasurementPosture_MEASUREMENT_POSTURE_LYING:
		return mysqldb.MeasurementPostureLying, nil
	}
	return mysqldb.MeasurementPostureInvalid, errors.New("invalid proto posture")
}

// setSubmitMeasurementReply 设置 SubmitMeasurementReply 返回数据
func setSubmitMeasurementReply(req *corepb.SubmitMeasurementInfoRequest, reply *corepb.SubmitMeasurementInfoResponse, resp *calcpb.SubmitCalcResponse, dataSendToApp []int32) error {
	hr := resp.GetResult().GetAverageHeartRate()
	//  处理站姿C0-C7
	reply.Hr = int32(hr)
	reply.PartialInfo = dataSendToApp

	// 设置 C0-C7 G0-G7
	C0 := int(resp.GetResult().GetC0())
	C1 := int(resp.GetResult().GetC1())
	C2 := int(resp.GetResult().GetC2())
	C3 := int(resp.GetResult().GetC3())
	C4 := int(resp.GetResult().GetC4())
	C5 := int(resp.GetResult().GetC5())
	C6 := int(resp.GetResult().GetC6())
	C7 := int(resp.GetResult().GetC7())

	G0 := resp.GetResult().GetG0()
	G1 := resp.GetResult().GetG1()
	G2 := resp.GetResult().GetG2()
	G3 := resp.GetResult().GetG3()
	G4 := resp.GetResult().GetG4()
	G5 := resp.GetResult().GetG5()
	G6 := resp.GetResult().GetG6()
	G7 := resp.GetResult().GetG7()
	if req.MeasurementPosture == corepb.MeasurementPosture_MEASUREMENT_POSTURE_STANDING {
		C0, C1, C2, C3, C4, C5, C6, C7 = ConvertToStandingCValues(req.Gender, C0, C1, C2, C3, C4, C5, C6, C7)
	}
	reply.C0 = Int32ValBoundedBy10FromInt(C0)
	reply.C1 = Int32ValBoundedBy10FromInt(C1)
	reply.C2 = Int32ValBoundedBy10FromInt(C2)
	reply.C3 = Int32ValBoundedBy10FromInt(C3)
	reply.C4 = Int32ValBoundedBy10FromInt(C4)
	reply.C5 = Int32ValBoundedBy10FromInt(C5)
	reply.C6 = Int32ValBoundedBy10FromInt(C6)
	reply.C7 = Int32ValBoundedBy10FromInt(C7)
	reply.G0 = int32(G0)
	reply.G1 = int32(G1)
	reply.G2 = int32(G2)
	reply.G3 = int32(G3)
	reply.G4 = int32(G4)
	reply.G5 = int32(G5)
	reply.G6 = int32(G6)
	reply.G7 = int32(G7)
	return nil
}

// getSavingData 生成发送到 AWS 存储桶的数据 (仅 Ir5160 数据)
func getSavingData(payload []byte, recordID int32) pulsetestinfopb.PulseTestRawInfo {
	var b bytes.Buffer
	b.Write(payload)
	return pulsetestinfopb.PulseTestRawInfo{
		Spec:     CurrentPulseTestDataSpec,
		RecordId: uint32(recordID),
		Payloads: b.Bytes(),
	}
}

// getPartialWaveData 1000条波形数据
func getPartialeWaveData(payload []byte) ([]int32, error) {
	const (
		// App 接收波形数据的开始位置
		appWaveDataStart = 3000
		// App 接收波形数据长度
		appWaveDataLength = 1000
	)

	// 获得解码后的原始数据的数组
	decodeModel, err := mapDeviceModelToCodec(RingDeviceModelForAlgorithm)
	if err != nil {
		return []int32{}, NewError(ErrInvalidDevice, fmt.Errorf("invalid device model %s", err.Error()))
	}
	decoder := ptcodec.NewDecoder(decodeModel)
	int32Array, err := decoder.Decode(payload)
	if err != nil {
		return []int32{}, NewError(ErrInvalidPayload, fmt.Errorf("failed to decode payload %s", err.Error()))
	}

	partialData := make([]int32, appWaveDataLength)
	for k, v := range int32Array[appWaveDataStart : appWaveDataStart+appWaveDataLength] {
		partialData[k] = int32(v)
	}
	return partialData, nil
}

// Int32ValBoundedBy10 返回 -10 到 10 之间的整数
func Int32ValBoundedBy10FromInt(val int) int32 {
	if val < -10 {
		return -10
	}
	if val > 10 {
		return 10
	}
	return int32(val)
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

// IntValBoundedBy10 返回 -10 到 10 之间的整数
func IntValBoundedBy10(val int) int {
	if val < -10 {
		return -10
	}
	if val > 10 {
		return 10
	}
	return val
}

// getHandByFinger 获取手指处于左手还是右手
func getHandByFinger(protoFinger corepb.Finger) (string, error) {
	switch protoFinger {
	case corepb.Finger_FINGER_INVALID:
		return "", fmt.Errorf("invalid proto finger %d", protoFinger)
	case corepb.Finger_FINGER_UNSET:
		return "", fmt.Errorf("invalid proto finger %d", protoFinger)
	case corepb.Finger_FINGER_LEFT_1:
		return LeftFinger, nil
	case corepb.Finger_FINGER_LEFT_2:
		return LeftFinger, nil
	case corepb.Finger_FINGER_LEFT_3:
		return LeftFinger, nil
	case corepb.Finger_FINGER_LEFT_4:
		return LeftFinger, nil
	case corepb.Finger_FINGER_LEFT_5:
		return LeftFinger, nil
	case corepb.Finger_FINGER_RIGHT_1:
		return RightFinger, nil
	case corepb.Finger_FINGER_RIGHT_2:
		return RightFinger, nil
	case corepb.Finger_FINGER_RIGHT_3:
		return RightFinger, nil
	case corepb.Finger_FINGER_RIGHT_4:
		return RightFinger, nil
	case corepb.Finger_FINGER_RIGHT_5:
		return RightFinger, nil
	}
	return "", fmt.Errorf("invalid proto finger %d", protoFinger)
}
