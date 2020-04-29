package handler

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	generalpb "github.com/jinmukeji/proto/v3/gen/micro/idl/ptypes/v2"
	jinmuidpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// UserSignUpSuite 用户测试
type UserSignUpSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

// SetupSuite 设置测试环境
func (suite *UserSignUpSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestUserSignUpByPhone 手机验证码注册
func (suite *UserSignUpSuite) TestUserSignUpByPhone() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	verificationNumber := getSignUpVerificationNumber(suite.JinmuIDService, *suite.Account)
	respUserSignUpByPhone := new(jinmuidpb.UserSignUpByPhoneResponse)
	birthday, _ := ptypes.TimestampProto(time.Now())
	err := suite.JinmuIDService.UserSignUpByPhone(ctx, &jinmuidpb.UserSignUpByPhoneRequest{
		Phone:              suite.Account.phone,
		Language:           generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		VerificationNumber: verificationNumber,
		PlainPassword:      suite.Account.password,
		NationCode:         suite.Account.nationCode,
		Profile: &jinmuidpb.UserProfile{
			Nickname:     suite.Account.nickname,
			Gender:       suite.Account.gender,
			BirthdayTime: birthday,
			Height:       suite.Account.height,
			Weight:       suite.Account.weight,
		},
	}, respUserSignUpByPhone)
	assert.NoError(t, err)
}

// TestUserSignUpByHKPhone 香港手机号码验证码注册
func (suite *UserSignUpSuite) TestUserSignUpByHKPhone() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	verificationNumber := getSignUpVerificationNumberHK(suite.JinmuIDService, *suite.Account)
	respUserSignUpByPhone := new(jinmuidpb.UserSignUpByPhoneResponse)
	birthday, _ := ptypes.TimestampProto(time.Now())
	err := suite.JinmuIDService.UserSignUpByPhone(ctx, &jinmuidpb.UserSignUpByPhoneRequest{
		Phone:              suite.Account.phoneHK,
		Language:           generalpb.Language_LANGUAGE_ENGLISH,
		VerificationNumber: verificationNumber,
		PlainPassword:      suite.Account.password,
		NationCode:         suite.Account.nationCodeHK,
		Profile: &jinmuidpb.UserProfile{
			Nickname:     suite.Account.nickname,
			Gender:       suite.Account.gender,
			BirthdayTime: birthday,
			Height:       suite.Account.height,
			Weight:       suite.Account.weight,
		},
	}, respUserSignUpByPhone)
	assert.NoError(t, err)
}

// TestUserSignUpByUSPhone  美国手机号码验证注册
func (suite *UserSignUpSuite) TestUserSignUpByUSPhone() {
	t := suite.T()
	ctx := context.Background()
	//发送通知
	verificationNumber := getSignUpVerificationNumberUS(suite.JinmuIDService, *suite.Account)
	respUserSignUpByPhone := new(jinmuidpb.UserSignUpByPhoneResponse)
	birthday, _ := ptypes.TimestampProto(time.Now())
	err := suite.JinmuIDService.UserSignUpByPhone(ctx, &jinmuidpb.UserSignUpByPhoneRequest{
		Phone:              suite.Account.phoneUS,
		Language:           generalpb.Language_LANGUAGE_ENGLISH,
		VerificationNumber: verificationNumber,
		PlainPassword:      suite.Account.password,
		NationCode:         suite.Account.nationCodeUSA,
		Profile: &jinmuidpb.UserProfile{
			Nickname:     suite.Account.nickname,
			Gender:       suite.Account.gender,
			BirthdayTime: birthday,
			Height:       suite.Account.height,
			Weight:       suite.Account.weight,
		},
	}, respUserSignUpByPhone)
	assert.NoError(t, err)
}

// TestUserSignUpByTWPhone  台湾手机号码验证注册
func (suite *UserSignUpSuite) TestUserSignUpByTWPhone() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	verificationNumber := getSignUpVerificationNumberTW(suite.JinmuIDService, *suite.Account)
	respUserSignUpByPhone := new(jinmuidpb.UserSignUpByPhoneResponse)
	birthday, _ := ptypes.TimestampProto(time.Now())
	err := suite.JinmuIDService.UserSignUpByPhone(ctx, &jinmuidpb.UserSignUpByPhoneRequest{
		Phone:              suite.Account.phoneTW,
		Language:           generalpb.Language_LANGUAGE_ENGLISH,
		VerificationNumber: verificationNumber,
		PlainPassword:      suite.Account.password,
		NationCode:         suite.Account.nationCodeTW,
		Profile: &jinmuidpb.UserProfile{
			Nickname:     suite.Account.nickname,
			Gender:       suite.Account.gender,
			BirthdayTime: birthday,
			Height:       suite.Account.height,
			Weight:       suite.Account.weight,
		},
	}, respUserSignUpByPhone)
	assert.NoError(t, err)
}

// TestUserSignUpByMacaoPhone  澳门手机号码验证注册
func (suite *UserSignUpSuite) TestUserSignUpByMacaoPhone() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	verificationNumber := getSignUpVerificationNumberMacao(suite.JinmuIDService, *suite.Account)
	respUserSignUpByPhone := new(jinmuidpb.UserSignUpByPhoneResponse)
	birthday, _ := ptypes.TimestampProto(time.Now())
	err := suite.JinmuIDService.UserSignUpByPhone(ctx, &jinmuidpb.UserSignUpByPhoneRequest{
		Phone:              suite.Account.phoneMacao,
		Language:           generalpb.Language_LANGUAGE_ENGLISH,
		VerificationNumber: verificationNumber,
		PlainPassword:      suite.Account.password,
		NationCode:         suite.Account.nationCodeMacao,
		Profile: &jinmuidpb.UserProfile{
			Nickname:     suite.Account.nickname,
			Gender:       suite.Account.gender,
			BirthdayTime: birthday,
			Height:       suite.Account.height,
			Weight:       suite.Account.weight,
		},
	}, respUserSignUpByPhone)
	assert.NoError(t, err)
}

// TestUserSignUpByCanadaPhone  加拿大手机号验证注册
func (suite *UserSignUpSuite) TestUserSignUpByCanadaPhone() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	verificationNumber := getSignUpVerificationNumberCanada(suite.JinmuIDService, *suite.Account)
	respUserSignUpByPhone := new(jinmuidpb.UserSignUpByPhoneResponse)
	birthday, _ := ptypes.TimestampProto(time.Now())
	err := suite.JinmuIDService.UserSignUpByPhone(ctx, &jinmuidpb.UserSignUpByPhoneRequest{
		Phone:              suite.Account.phoneCanada,
		Language:           generalpb.Language_LANGUAGE_ENGLISH,
		VerificationNumber: verificationNumber,
		PlainPassword:      suite.Account.password,
		NationCode:         suite.Account.nationCodeUSA,
		Profile: &jinmuidpb.UserProfile{
			Nickname:     suite.Account.nickname,
			Gender:       suite.Account.gender,
			BirthdayTime: birthday,
			Height:       suite.Account.height,
			Weight:       suite.Account.weight,
		},
	}, respUserSignUpByPhone)
	assert.NoError(t, err)
}

// TestUserSignUpByUKPhone  英国手机号码注册
func (suite *UserSignUpSuite) TestUserSignUpByUKPhone() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	verificationNumber := getSignUpVerificationNumberUK(suite.JinmuIDService, *suite.Account)
	respUserSignUpByPhone := new(jinmuidpb.UserSignUpByPhoneResponse)
	birthday, _ := ptypes.TimestampProto(time.Now())
	err := suite.JinmuIDService.UserSignUpByPhone(ctx, &jinmuidpb.UserSignUpByPhoneRequest{
		Phone:              suite.Account.phoneUK,
		Language:           generalpb.Language_LANGUAGE_ENGLISH,
		VerificationNumber: verificationNumber,
		PlainPassword:      suite.Account.password,
		NationCode:         suite.Account.nationCodeUK,
		Profile: &jinmuidpb.UserProfile{
			Nickname:     suite.Account.nickname,
			Gender:       suite.Account.gender,
			BirthdayTime: birthday,
			Height:       suite.Account.height,
			Weight:       suite.Account.weight,
		},
	}, respUserSignUpByPhone)
	assert.NoError(t, err)
}

// TestUserSignUpByJPPhone  日本手机号码注册
func (suite *UserSignUpSuite) TestUserSignUpByJPPhone() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	verificationNumber := getSignUpVerificationNumberJP(suite.JinmuIDService, *suite.Account)
	respUserSignUpByPhone := new(jinmuidpb.UserSignUpByPhoneResponse)
	birthday, _ := ptypes.TimestampProto(time.Now())
	err := suite.JinmuIDService.UserSignUpByPhone(ctx, &jinmuidpb.UserSignUpByPhoneRequest{
		Phone:              suite.Account.phoneJP,
		Language:           generalpb.Language_LANGUAGE_ENGLISH,
		VerificationNumber: verificationNumber,
		PlainPassword:      suite.Account.password,
		NationCode:         suite.Account.nationCodeJP,
		Profile: &jinmuidpb.UserProfile{
			Nickname:     suite.Account.nickname,
			Gender:       suite.Account.gender,
			BirthdayTime: birthday,
			Height:       suite.Account.height,
			Weight:       suite.Account.weight,
		},
	}, respUserSignUpByPhone)
	assert.NoError(t, err)
}

// TestUserSignUpByPhoneIsExist 手机已注册
func (suite *UserSignUpSuite) TestUserSignUpByPhoneIsExist() {
	t := suite.T()
	ctx := context.Background()

	// 发送通知
	verificationNumber := getSignUpVerificationNumber(suite.JinmuIDService, *suite.Account)
	birthday, _ := ptypes.TimestampProto(time.Now())
	respUserSignUpByPhone := new(jinmuidpb.UserSignUpByPhoneResponse)
	err := suite.JinmuIDService.UserSignUpByPhone(ctx, &jinmuidpb.UserSignUpByPhoneRequest{
		Phone:              suite.Account.phone,
		Language:           generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		VerificationNumber: verificationNumber,
		PlainPassword:      suite.Account.password,
		NationCode:         suite.Account.nationCode,
		Profile: &jinmuidpb.UserProfile{
			Nickname:     suite.Account.nickname,
			Gender:       suite.Account.gender,
			BirthdayTime: birthday,
			Height:       suite.Account.height,
			Weight:       suite.Account.weight,
		},
	}, respUserSignUpByPhone)
	assert.Error(t, errors.New("[errcode:2000] phone number has been registered"), err)
}

// TestUserSignUpByPhoneIsNull手机号为空
func (suite *UserSignUpSuite) TestUserSignUpByPhoneIsNull() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	verificationNumber := getSignUpVerificationNumber(suite.JinmuIDService, *suite.Account)
	birthday, _ := ptypes.TimestampProto(time.Now())
	respUserSignUpByPhone := new(jinmuidpb.UserSignUpByPhoneResponse)
	err := suite.JinmuIDService.UserSignUpByPhone(ctx, &jinmuidpb.UserSignUpByPhoneRequest{
		Phone:              suite.Account.phoneIsNull,
		Language:           generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		VerificationNumber: verificationNumber,
		PlainPassword:      suite.Account.password,
		NationCode:         suite.Account.nationCode,
		Profile: &jinmuidpb.UserProfile{
			Nickname:     suite.Account.nickname,
			Gender:       suite.Account.gender,
			BirthdayTime: birthday,
			Height:       suite.Account.height,
			Weight:       suite.Account.weight,
		},
	}, respUserSignUpByPhone)
	assert.Error(t, errors.New("[errcode:31000] vc is invalid"), err)
}

// TestUserSignUpByPhoneFormatIsError 手机号格式不正确
func (suite *UserSignUpSuite) TestUserSignUpByPhoneFormatIsError() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	verificationNumber := getSignUpVerificationNumber(suite.JinmuIDService, *suite.Account)
	birthday, _ := ptypes.TimestampProto(time.Now())
	respUserSignUpByPhone := new(jinmuidpb.UserSignUpByPhoneResponse)
	err := suite.JinmuIDService.UserSignUpByPhone(ctx, &jinmuidpb.UserSignUpByPhoneRequest{
		Phone:              suite.Account.PhoneError,
		Language:           generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		VerificationNumber: verificationNumber,
		PlainPassword:      suite.Account.password,
		NationCode:         suite.Account.nationCode,
		Profile: &jinmuidpb.UserProfile{
			Nickname:     suite.Account.nickname,
			Gender:       suite.Account.gender,
			BirthdayTime: birthday,
			Height:       suite.Account.height,
			Weight:       suite.Account.weight,
		},
	}, respUserSignUpByPhone)
	assert.Error(t, errors.New("[errcode:31000] vc is invalid"), err)
}

// TestUserSignUpByPhonePasswordIsNull mvc为空
func (suite *UserSignUpSuite) TestUserSignUpByPhonePasswordIsNull() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	verificationNumber := getSignUpVerificationNumber(suite.JinmuIDService, *suite.Account)
	birthday, _ := ptypes.TimestampProto(time.Now())
	respUserSignUpByPhone := new(jinmuidpb.UserSignUpByPhoneResponse)
	err := suite.JinmuIDService.UserSignUpByPhone(ctx, &jinmuidpb.UserSignUpByPhoneRequest{
		Phone:              suite.Account.phone,
		Language:           generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		VerificationNumber: verificationNumber,
		PlainPassword:      suite.Account.hashedPasswordIsNull,
		NationCode:         suite.Account.nationCode,
		Profile: &jinmuidpb.UserProfile{
			Nickname:     suite.Account.nickname,
			Gender:       suite.Account.gender,
			BirthdayTime: birthday,
			Height:       suite.Account.height,
			Weight:       suite.Account.weight,
		},
	}, respUserSignUpByPhone)
	assert.Error(t, errors.New("[errcode:31000] vc is invalid"), err)
}

// TestUserSignUpByPhoneVerificationNumberIsNull  VerificationNumber为空
func (suite *UserSignUpSuite) TestUserSignUpByPhoneVerificationNumberIsNull() {
	t := suite.T()
	ctx := context.Background()
	birthday, _ := ptypes.TimestampProto(time.Now())
	respUserSignUpByPhone := new(jinmuidpb.UserSignUpByPhoneResponse)
	err := suite.JinmuIDService.UserSignUpByPhone(ctx, &jinmuidpb.UserSignUpByPhoneRequest{
		Phone:              suite.Account.phone,
		Language:           generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		VerificationNumber: "",
		PlainPassword:      suite.Account.password,
		NationCode:         suite.Account.nationCode,
		Profile: &jinmuidpb.UserProfile{
			Nickname:     suite.Account.nickname,
			Gender:       suite.Account.gender,
			BirthdayTime: birthday,
			Height:       suite.Account.height,
			Weight:       suite.Account.weight,
		},
	}, respUserSignUpByPhone)
	assert.Error(t, errors.New("[errcode:31000] vc is invalid"), err)
}

// TestUserSignUpByPhoneNationCodeIsNull  nationcode为空
func (suite *UserSignUpSuite) TestUserSignUpByPhoneNationCodeIsNull() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	verificationNumber := getSignUpVerificationNumber(suite.JinmuIDService, *suite.Account)
	birthday, _ := ptypes.TimestampProto(time.Now())
	respUserSignUpByPhone := new(jinmuidpb.UserSignUpByPhoneResponse)
	err := suite.JinmuIDService.UserSignUpByPhone(ctx, &jinmuidpb.UserSignUpByPhoneRequest{
		Phone:              suite.Account.phoneNotExist,
		Language:           generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		VerificationNumber: verificationNumber,
		PlainPassword:      suite.Account.password,
		NationCode:         suite.Account.nationCodeIsNull,
		Profile: &jinmuidpb.UserProfile{
			Nickname:     suite.Account.nickname,
			Gender:       suite.Account.gender,
			BirthdayTime: birthday,
			Height:       suite.Account.height,
			Weight:       suite.Account.weight,
		},
	}, respUserSignUpByPhone)
	assert.Error(t, errors.New("[errcode:31000] vc is invalid"), err)
}

// TestUserSignUpByPhoneNotCorrespondNationCode  手机号和验证码不一致
func (suite *UserSignUpSuite) TestUserSignUpByPhoneNotCorrespondNationCode() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	verificationNumber := getSignUpVerificationNumber(suite.JinmuIDService, *suite.Account)
	birthday, _ := ptypes.TimestampProto(time.Now())
	respUserSignUpByPhone := new(jinmuidpb.UserSignUpByPhoneResponse)
	err := suite.JinmuIDService.UserSignUpByPhone(ctx, &jinmuidpb.UserSignUpByPhoneRequest{
		Phone:              suite.Account.phone,
		Language:           generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		VerificationNumber: verificationNumber,
		PlainPassword:      suite.Account.password,
		NationCode:         suite.Account.nationCodeUSA,
		Profile: &jinmuidpb.UserProfile{
			Nickname:     suite.Account.nickname,
			Gender:       suite.Account.gender,
			BirthdayTime: birthday,
			Height:       suite.Account.height,
			Weight:       suite.Account.weight,
		},
	}, respUserSignUpByPhone)
	assert.Error(t, errors.New("[errcode:36000] verification number is invalid"), err)
}

// TestUserSignUpByPhoneNicknameIsNull  nickname为空
func (suite *UserSignUpSuite) TestUserSignUpByPhoneNicknameIsNull() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	verificationNumber := getSignUpVerificationNumber(suite.JinmuIDService, *suite.Account)
	respUserSignUpByPhone := new(jinmuidpb.UserSignUpByPhoneResponse)
	birthday, _ := ptypes.TimestampProto(time.Now())
	err := suite.JinmuIDService.UserSignUpByPhone(ctx, &jinmuidpb.UserSignUpByPhoneRequest{
		Phone:              suite.Account.phone,
		Language:           generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		VerificationNumber: verificationNumber,
		PlainPassword:      suite.Account.password,
		NationCode:         suite.Account.nationCode,
		Profile: &jinmuidpb.UserProfile{
			Nickname:     "",
			Gender:       suite.Account.gender,
			BirthdayTime: birthday,
			Height:       suite.Account.height,
			Weight:       suite.Account.weight,
		},
	}, respUserSignUpByPhone)
	assert.NoError(t, err)
}

// TestUserSignUpByPhoneHeightIsNull   height为空
func (suite *UserSignUpSuite) TestUserSignUpByPhoneHeightIsNull() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	verificationNumber := getSignUpVerificationNumber(suite.JinmuIDService, *suite.Account)
	respUserSignUpByPhone := new(jinmuidpb.UserSignUpByPhoneResponse)
	birthday, _ := ptypes.TimestampProto(time.Now())
	err := suite.JinmuIDService.UserSignUpByPhone(ctx, &jinmuidpb.UserSignUpByPhoneRequest{
		Phone:              suite.Account.phone,
		Language:           generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		VerificationNumber: verificationNumber,
		PlainPassword:      suite.Account.password,
		NationCode:         suite.Account.nationCode,
		Profile: &jinmuidpb.UserProfile{
			Nickname:     suite.Account.nickname,
			Gender:       suite.Account.gender,
			BirthdayTime: birthday,
			Height:       suite.Account.heightNull,
			Weight:       suite.Account.weight,
		},
	}, respUserSignUpByPhone)
	assert.NoError(t, err)
}

// TestUserSignUpByPhoneHeightIsLow  height小于标准值
func (suite *UserSignUpSuite) TestUserSignUpByPhoneHeightIsLow() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	verificationNumber := getSignUpVerificationNumber(suite.JinmuIDService, *suite.Account)
	respUserSignUpByPhone := new(jinmuidpb.UserSignUpByPhoneResponse)
	birthday, _ := ptypes.TimestampProto(time.Now())
	err := suite.JinmuIDService.UserSignUpByPhone(ctx, &jinmuidpb.UserSignUpByPhoneRequest{
		Phone:              suite.Account.phone,
		Language:           generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		VerificationNumber: verificationNumber,
		PlainPassword:      suite.Account.password,
		NationCode:         suite.Account.nationCode,
		Profile: &jinmuidpb.UserProfile{
			Nickname:     suite.Account.nickname,
			Gender:       suite.Account.gender,
			BirthdayTime: birthday,
			Height:       suite.Account.heightLow,
			Weight:       suite.Account.weight,
		},
	}, respUserSignUpByPhone)
	assert.NoError(t, err)
}

// TestUserSignUpByPhoneHeightIsHigh   height大于标准值
func (suite *UserSignUpSuite) TestUserSignUpByPhoneHeightIsHigh() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	verificationNumber := getSignUpVerificationNumber(suite.JinmuIDService, *suite.Account)
	respUserSignUpByPhone := new(jinmuidpb.UserSignUpByPhoneResponse)
	birthday, _ := ptypes.TimestampProto(time.Now())
	err := suite.JinmuIDService.UserSignUpByPhone(ctx, &jinmuidpb.UserSignUpByPhoneRequest{
		Phone:              suite.Account.phone,
		Language:           generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		VerificationNumber: verificationNumber,
		PlainPassword:      suite.Account.password,
		NationCode:         suite.Account.nationCode,
		Profile: &jinmuidpb.UserProfile{
			Nickname:     suite.Account.nickname,
			Gender:       suite.Account.gender,
			BirthdayTime: birthday,
			Height:       suite.Account.heightHigh,
			Weight:       suite.Account.weight,
		},
	}, respUserSignUpByPhone)
	assert.NoError(t, err)
}

// TestUserSignUpByPhoneWeightIsNull weight为空
func (suite *UserSignUpSuite) TestUserSignUpByPhoneWeightIsNull() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	verificationNumber := getSignUpVerificationNumber(suite.JinmuIDService, *suite.Account)
	respUserSignUpByPhone := new(jinmuidpb.UserSignUpByPhoneResponse)
	birthday, _ := ptypes.TimestampProto(time.Now())
	err := suite.JinmuIDService.UserSignUpByPhone(ctx, &jinmuidpb.UserSignUpByPhoneRequest{
		Phone:              suite.Account.phone,
		Language:           generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		VerificationNumber: verificationNumber,
		PlainPassword:      suite.Account.password,
		NationCode:         suite.Account.nationCode,
		Profile: &jinmuidpb.UserProfile{
			Nickname:     suite.Account.nickname,
			Gender:       suite.Account.gender,
			BirthdayTime: birthday,
			Height:       suite.Account.height,
			Weight:       suite.Account.weightNull,
		},
	}, respUserSignUpByPhone)
	assert.NoError(t, err)
}

// TestUserSignUpByPhoneWeightIsLow weight为空
func (suite *UserSignUpSuite) TestUserSignUpByPhoneWeightIsLow() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	verificationNumber := getSignUpVerificationNumber(suite.JinmuIDService, *suite.Account)
	respUserSignUpByPhone := new(jinmuidpb.UserSignUpByPhoneResponse)
	birthday, _ := ptypes.TimestampProto(time.Now())
	err := suite.JinmuIDService.UserSignUpByPhone(ctx, &jinmuidpb.UserSignUpByPhoneRequest{
		Phone:              suite.Account.phone,
		Language:           generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		VerificationNumber: verificationNumber,
		PlainPassword:      suite.Account.password,
		NationCode:         suite.Account.nationCode,
		Profile: &jinmuidpb.UserProfile{
			Nickname:     suite.Account.nickname,
			Gender:       suite.Account.gender,
			BirthdayTime: birthday,
			Height:       suite.Account.height,
			Weight:       suite.Account.weightLow,
		},
	}, respUserSignUpByPhone)
	assert.NoError(t, err)
}

// TestUserSignUpByPhoneWeightIsHigh weight为空
func (suite *UserSignUpSuite) TestUserSignUpByPhoneWeightIsHigh() {
	t := suite.T()
	ctx := context.Background()
	// 发送通知
	verificationNumber := getSignUpVerificationNumber(suite.JinmuIDService, *suite.Account)
	respUserSignUpByPhone := new(jinmuidpb.UserSignUpByPhoneResponse)
	birthday, _ := ptypes.TimestampProto(time.Now())
	err := suite.JinmuIDService.UserSignUpByPhone(ctx, &jinmuidpb.UserSignUpByPhoneRequest{
		Phone:              suite.Account.phone,
		Language:           generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		VerificationNumber: verificationNumber,
		PlainPassword:      suite.Account.password,
		NationCode:         suite.Account.nationCode,
		Profile: &jinmuidpb.UserProfile{
			Nickname:     suite.Account.nickname,
			Gender:       suite.Account.gender,
			BirthdayTime: birthday,
			Height:       suite.Account.height,
			Weight:       suite.Account.weightHigh,
		},
	}, respUserSignUpByPhone)
	assert.NoError(t, err)
}

// TestUserSignUpByPhoneError  VerificationNumber为空
func (suite *UserSignUpSuite) TestUserSignUpByPhoneError() {
	t := suite.T()
	ctx := context.Background()
	birthday, _ := ptypes.TimestampProto(time.Now())
	respUserSignUpByPhone := new(jinmuidpb.UserSignUpByPhoneResponse)
	verificationNumber := getSignUpVerificationNumber(suite.JinmuIDService, *suite.Account)
	err := suite.JinmuIDService.UserSignUpByPhone(ctx, &jinmuidpb.UserSignUpByPhoneRequest{
		Phone:              suite.Account.phone,
		Language:           generalpb.Language_LANGUAGE_SIMPLIFIED_CHINESE,
		VerificationNumber: verificationNumber,
		PlainPassword:      suite.Account.password,
		NationCode:         suite.Account.nationCode,
		Profile: &jinmuidpb.UserProfile{
			Nickname:     suite.Account.nickname,
			Gender:       suite.Account.gender,
			BirthdayTime: birthday,
			Height:       suite.Account.height,
			Weight:       suite.Account.weight,
		},
	}, respUserSignUpByPhone)
	assert.Error(t, errors.New("[errcode:31000] vc is invalid"), err)
}

func (suite *UserSignUpSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}

func TestUserSignUpTestSuite(t *testing.T) {
	suite.Run(t, new(UserSignUpSuite))
}
