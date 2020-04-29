package handler

import (
	"io/ioutil"
	"testing"

	generalpb "github.com/jinmukeji/proto/v3/gen/micro/idl/ptypes/v2"
	"github.com/stretchr/testify/suite"
)

// SubmitMeasurementInfoTestSuite 测试提交测量数据
type SubmitMeasurementInfoTestSuite struct {
	suite.Suite
	jinmuHealth *JinmuHealth
	Account     *Account
}

/*
// SetupSuite 初始化测试类
func (suite *SubmitMeasurementInfoTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-biz-core.env")
	_ = godotenv.Load(envFilepath)
	db, _ := newTestingDbClientFromEnvFile(envFilepath)
	awsClient, _ := aws.NewClient(
		aws.BucketName(os.Getenv("X_AWS_BUCKET_NAME")),
		aws.SecretKey(os.Getenv("X_AWS_SECRET_KEY")),
		aws.AccessKeyID(os.Getenv("X_AWS_ACCESS_KEY")),
		aws.Region(os.Getenv("X_AWS_REGION")),
		aws.PulseTestRawDataEnvironmentS3KeyPrefix(os.Getenv("X_WAVE_DATA_KEY_PREFIX")),
	)
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
	configDoc, _ := blocker.LoadConfig("../../pkg/blocker/testdata/config_doc.yml")
	blockerPool, _ := blocker.NewBlockerPool(configDoc, "../../pkg/blocker/data/GeoLite2-Country.mmdb.gz")
	suite.jinmuHealth = NewJinmuHealth(db, nil, algorithmClient, awsClient, nil, nil, blockerPool)
}

// TestSubmitMeasurementInfoSuccess 测试提交测量数据成功
func (suite *SubmitMeasurementInfoTestSuite) TestSubmitMeasurementInfoSuccess() {
	t := suite.T()
	var algorithmReq algorithm.C2Request
	assert.NoError(t, json.Unmarshal(getTestMeasureData(), &algorithmReq))
	ctx := context.Background()
	ctx = mockAuth(ctx, suite.Account.clientID, suite.Account.name, suite.Account.zone)
	ctx = auth.AddContextUserID(ctx, suite.Account.userID)
	const registerType = "username"
	ctx, err := mockSignin(ctx, suite.jinmuHealth, suite.Account.userName, suite.Account.passwordHash, registerType, proto.SignInMethod_SIGN_IN_METHOD_GENERAL)
	assert.NoError(t, err)

	ctx = addContextClient(ctx, metaClient{
		RemoteClientIP: "110.165.32.0",
		ClientID:       suite.Account.clientID,
		Name:           suite.Account.name,
		Zone:           suite.Account.zone,
		CustomizedCode: "",
	})
	req, resp := new(proto.SubmitMeasurementInfoRequest), new(proto.SubmitMeasurementInfoResponse)
	req.UserId = suite.Account.userID
	req.Mac = algorithmReq.MAC
	setTestSubmitMeasurementInfoRequest(req, int(suite.Account.userID), &algorithmReq)
	assert.NoError(t, suite.jinmuHealth.SubmitMeasurementInfo(ctx, req, resp))
	assert.NotZero(t, resp.Hr)
}
*/
/*
一、测试成功
    a、mac未被拒，同时mac对应创建时间在2019/06/01之前
        1、clientID=dengyun-10001
            (1)测试标题
                mac未被拒，即zone在范围内；mac对应的创建时间在2019/06/01日之前，测试成功
            (2)预制条件
                mac创建时间：2018-08-16 20:50:06
                zone:CN-X
                customized_code:custom_dengyun
            (3)测试输入
                clientID:dengyun-10001
                ip:223.104.145.236
                mac:304511440CCF
            (4)预计结果
                测试成功
            (5)实际结果
                测试成功
        2、clientID=kangmei-10001
            (1)测试标题
                mac未被拒，即zone在范围内；mac对应的创建时间在2019/06/01日之前，测试成功
            (2)预制条件
                mac创建时间：2018-08-16 20:50:06
                zone:CN-X
                customized_code:custom_dengyun
            (3)测试输入
                clientID:kangmei-10001
                ip:223.104.145.236
                mac:304511440CCF
            (4)预计结果
                测试成功
            (5)实际结果
                测试成功
        3、clientID=jm-10001
            (1)测试标题
                mac未被拒，即zone在范围内；mac对应的创建时间在2019/06/01日之前，测试成功
            (2)预制条件
                mac创建时间：2018-08-16 20:50:06
                zone:CN
            (3)测试输入
                clientID:jm-10002
                ip:223.104.145.236
                mac:304511440CCF
            (4)预计结果
                测试成功
            (5)实际结果
                测试成功
        4、clientID=jm-10002
            (1)测试标题
                mac未被拒，即zone在范围内；mac对应的创建时间在2019/06/01日之前，测试成功
            (2)预制条件
                mac创建时间：2018-08-16 20:50:06
                zone:CN
            (3)测试输入
                clientID:jm-10003
                ip:223.104.145.236
                mac:304511440CCF
            (4)预计结果
                测试成功
            (5)实际结果
                测试成功
        6、clientID=jm-10004
            (1)测试标题
                mac未被拒，即zone在范围内；mac对应的创建时间在2019/06/01日之前，测试成功
            (2)预制条件
                mac创建时间：2018-08-16 20:50:06
                zone:CN
            (3)测试输入
                clientID:jm-10004
                ip:223.104.145.236
                mac:304511440CCF
            (4)预计结果
                测试成功
            (5)实际结果
                测试成功
        7、clientID=jm-10005
            (1)测试标题
                mac未被拒，即zone在范围内；mac对应的创建时间在2019/06/01日之前，测试成功
            (2)预制条件
                mac创建时间：2020-08-16 20:50:06
                zone:CN
            (3)测试输入
                clientID:jm-10005
                ip:223.104.145.236
                mac:304511440CCF
            (4)预计结果
                测试成功
            (5)实际结果
                测试成功
    b、mac未被拒，同时mac对应创建时间在2019/06/01之后，看ip过滤状态，允许任意ip通过
        8、clientID=dengyun-10001
            (1)测试标题
                mac未被拒，即zone在范围内；mac对应的创建时间在2019/06/01日之后，美国ip，测试成功
            (2)预制条件
                mac创建时间：2020-08-16 20:50:06
                zone:CN-X
            (3)测试输入
                clientID:dengyun-10001
                ip:67.220.91.30
                mac:304511440CCF
            (4)预计结果
                测试成功
            (5)实际结果
                测试成功
        9、clientID=kangmei-10001
            (1)测试标题
                mac未被拒，即zone在范围内；mac对应的创建时间在2019/06/01日之后，美国ip，测试成功
            (2)预制条件
                mac创建时间：2020-08-16 20:50:06
                zone:CN-X
            (3)测试输入
                clientID:kangmei-10001
                ip:67.220.91.30
                mac:304511440CCF
            (4)预计结果
                测试成功
            (5)实际结果
                测试成功
        10、clientID=jm-10002
            (1)测试标题
                mac未被拒，即zone在范围内；mac对应的创建时间在2019/06/01日之后，美国ip，测试成功
            (2)预制条件
                mac创建时间：2018-08-16 20:50:06
                zone:CN
            (3)测试输入
                clientID:jm-10002
                ip:67.220.91.30
                mac:304511440CCF
            (4)预计结果
                测试成功
            (5)实际结果
                测试成功
        11、clientID=jm-10003
            (1)测试标题
                mac未被拒，即zone在范围内；mac对应的创建时间在2019/06/01日之后，美国ip，测试成功
            (2)预制条件
                mac创建时间：2018-08-16 20:50:06
                zone:CN
            (3)测试输入
                clientID:jm-10003
                ip:67.220.91.30
                mac:304511440CCF
            (4)预计结果
                测试成功
            (5)实际结果
                测试成功
    c、mac未被拒，同时mac对应创建时间在2019/06/01之后，看ip过滤状态，允许某些特定城市ip通过
        12、clientID=jm-10001
            (1)测试标题
                mac未被拒，即zone在范围内；mac对应的创建时间在2019/06/01日之后，香港ip,测试成功
            (2)预制条件
                mac创建时间：2020-08-16 20:50:06
                zone:CN
            (3)测试输入
                clientID:jm-10001
                ip:110.165.32.0
                mac:304511440CCF
            (4)预计结果
                测试成功
            (5)实际结果
                测试成功
        13、clientID=jm-10004
            (1)测试标题
                mac未被拒，即zone在范围内；mac对应的创建时间在2019/06/01日之后，香港ip,测试成功
            (2)预制条件
                mac创建时间：2020-08-16 20:50:06
                zone:CN
            (3)测试输入
                clientID:jm-10004
                ip:110.165.32.0
                mac:304511440CCF
            (4)预计结果
                测试成功
            (5)实际结果
                测试成功
        14、clientID=jm-10005
            (1)测试标题
                mac未被拒，即zone在范围内；mac对应的创建时间在2019/06/01日之后，香港ip,测试成功
            (2)预制条件
                mac创建时间：2020-08-16 20:50:06
                zone:CN
            (3)测试输入
                clientID:jm-10005
                ip:110.165.32.0
                mac:304511440CCF
            (4)预计结果
                测试成功
            (5)实际结果
                测试成功
二、测试失败
    1、mac过滤
        (1)测试标题
            mac被拒，即zone不在范围内；mac对应的创建时间在2019/06/01日之前，测试失败
        (2)预制条件
            mac创建时间：2018-08-16 20:50:06
            zone:CN
        (3)测试输入
            clientID:dengyun-10001
            ip:223.104.145.236
            mac:304511440CCF
        (4)预计结果
            测试失败，[errcode:77000] blocked mac
        (5)实际结果
            测试失败，[errcode:77000] blocked mac
    2、mac未被拒，创建时间在2019/06/01之后，ip过滤
        a、clientID:jm-10001
            (1)测试标题
                mac未被拒，即zone在范围内；mac对应的创建时间在2019/06/01日之后，ip过滤失败，测试失败
            (2)预制条件
                mac创建时间：2020-08-16 20:50:06
                zone:CN
            (3)测试输入
                clientID:jm-10001
                ip:67.220.91.30
                mac:304511440CCF
            (4)预计结果
                测试失败，[errcode:78000] blocked ip
            (5)实际结果
                测试失败，[errcode:78000] blocked ip
        b、clientID:jm-10004
            (1)测试标题
                mac未被拒，即zone在范围内；mac对应的创建时间在2019/06/01日之后，ip过滤失败，测试失败
            (2)预制条件
                mac创建时间：2020-08-16 20:50:06
                zone:CN
            (3)测试输入
                clientID:jm-10004
                ip:67.220.91.30
                mac:304511440CCF
            (4)预计结果
                测试失败，[errcode:78000] blocked ip
            (5)实际结果
                测试失败，[errcode:78000] blocked ip
        c、clientID:jm-10005
            (1)测试标题
                mac未被拒，即zone在范围内；mac对应的创建时间在2019/06/01日之后，ip过滤失败，测试失败
            (2)预制条件
                mac创建时间：2020-08-16 20:50:06
                zone:CN
            (3)测试输入
                clientID:jm-10005
                ip:67.220.91.30
                mac:304511440CCF
            (4)预计结果
                测试失败，[errcode:78000] blocked ip
            (5)实际结果
                测试失败，[errcode:78000] blocked ip
*/

// getTestMeasureData 读取测量数据文件
func getTestMeasureData() []byte {
	jsondata, _ := ioutil.ReadFile("testdata/test.json")
	return jsondata
}

/*

// setTestSubmitMeasurementInfoRequest 从测试 json 数据设置一次 mock rpc 请求
func setTestSubmitMeasurementInfoRequest(req *proto.SubmitMeasurementInfoRequest, subjectID int, algorithmReq *algorithm.C2Request) {
	dataNum, _ := strconv.ParseInt(algorithmReq.Bit, 10, 0)
	req.InfoNum = int32(dataNum)
	req.MobileType = proto.MobileType_MOBILE_TYPE_ANDROID
	req.Mac = algorithmReq.MAC
	protoGender := mapAlgorithmGenderToProto(algorithmReq.Gender)
	req.Gender = protoGender
	weight, _ := strconv.Atoi(algorithmReq.Weight)
	height, _ := strconv.Atoi(algorithmReq.Height)
	req.Age = int32(algorithmReq.Age)
	req.Info0 = new(proto.BluetoothInfo)
	req.Info0.Ir5160, _ = base64.StdEncoding.DecodeString(algorithmReq.Data0.IR5160)
	req.Info0.Ir5160Md5 = algorithmReq.Data0.IR5160MD5
	req.Info0.R5164 = []byte(algorithmReq.Data0.R5164)
	req.Info0.R5164Md5 = algorithmReq.Data0.R5164MD5
	req.Info1 = new(proto.BluetoothInfo)
	req.Info1.Ir5160, _ = base64.StdEncoding.DecodeString(algorithmReq.Data1.IR5160)
	req.Info1.Ir5160Md5 = algorithmReq.Data1.IR5160MD5
	req.Info1.R5164 = []byte(algorithmReq.Data1.R5164)
	req.Info1.R5164Md5 = algorithmReq.Data1.R5164MD5
	req.Weight = int32(weight)
	req.Height = int32(height)
}
*/
func TestSubmitMeasurementInfoTestSuite(t *testing.T) {
	suite.Run(t, new(SubmitMeasurementInfoTestSuite))
}

func mapAlgorithmGenderToProto(gender string) generalpb.Gender {
	switch gender {
	case "M":
		return generalpb.Gender_GENDER_MALE
	case "F":
		return generalpb.Gender_GENDER_FEMALE
	}
	return generalpb.Gender_GENDER_MALE
}
