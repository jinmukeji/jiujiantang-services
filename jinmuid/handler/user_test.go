package handler

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	generalpb "github.com/jinmukeji/proto/v3/gen/micro/idl/ptypes/v2"
	jinmuidpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// UserTestSuite 用户测试
type UserTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

type Account struct {
	username             string
	password             string
	userID               int32
	seed                 string
	hashedPassword       string
	nationCode           string
	phone                string
	email                string
	phoneIsNull          string
	PhoneError           string
	mvcIsNull            string
	serialNumberIsNull   string
	nationCodeIsNull     string
	nationCodeUSA        string
	seedNull             string
	phoneNotExist        string
	hashedPasswordIsNull string
	emailNoExist         string
	mvcError             string
	emailNull            string
	emailError           string
	usernameExist        string
	usernameNull         string
	questionKey1         string
	questionKey2         string
	questionKey3         string
	answer1              string
	answer2              string
	answer3              string
	questionKey1New      string
	questionKey2New      string
	questionKey3New      string
	answer1New           string
	answer2New           string
	answer3New           string
	questionKey1Same     string
	questionKey2Same     string
	questionKey3Same     string
	answer1Same          string
	answer2Same          string
	answer3Same          string
	questionKey1Null     string
	questionKey2Null     string
	questionKey3Null     string
	answer1Null          string
	answer2Null          string
	answer3Null          string
	answer1Long          string
	answer1Format        string
	questionKey1NoExist  string
	nickname             string
	gender               generalpb.Gender
	weight               int32
	height               int32
	verificationType     string
	emailNew             string
	remark               string
	customizedCode       string
	phonePassword        string
	userIDNotExist       int32
	vt                   string
	weightNull           int32
	weightLow            int32
	weightHigh           int32
	heightNull           int32
	heightLow            int32
	heightHigh           int32
	mvcFormatError       string
	phoneFormatError     string
	nationCodeHK         string
	phoneHK              string
	phoneUS              string
	phoneTW              string
	nationCodeTW         string
	phoneMacao           string
	nationCodeMacao      string
	phoneCanada          string
	phoneUK              string
	nationCodeUK         string
	phoneJP              string
	nationCodeJP         string
}

// SetupSuite 设置测试环境
func (suite *UserTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
	suite.JinmuIDService.encryptKey = newTestingEncryptKeyFromEnvFile(envFilepath)
}

func newTestingAccountFromEnvFile(filepath string) *Account {
	_ = godotenv.Load(filepath)
	userID, _ := strconv.Atoi(os.Getenv("X_TEST_USER_ID"))
	weight, _ := strconv.Atoi(os.Getenv("X_TEST_WEIGHT"))
	height, _ := strconv.Atoi(os.Getenv("X_TEST_HEIGHT"))
	weightNull, _ := strconv.Atoi(os.Getenv("X_TEST_WEIGHT_NULL"))
	weightLow, _ := strconv.Atoi(os.Getenv("X_TEST_WEIGHT_LOW"))
	weightHigh, _ := strconv.Atoi(os.Getenv("X_TEST_WEIGHT_HIGH"))
	heightNull, _ := strconv.Atoi(os.Getenv("X_TEST_HEIGHT_NULL"))
	heightLow, _ := strconv.Atoi(os.Getenv("X_TEST_HEIGHT_LOW"))
	heightHigh, _ := strconv.Atoi(os.Getenv("X_TEST_HEIGHT_HIGH"))
	userIDNotExist, _ := strconv.Atoi(os.Getenv("X_TEST_USER_ID_NOT_EXIST"))
	protoGender, _ := mapTestGenderToProto(os.Getenv("X_TEST_ANSWER1_FORMAT"))
	return &Account{
		os.Getenv("X_TEST_USERNAME"),
		os.Getenv("X_TEST_PASSWORD"),
		int32(userID),
		os.Getenv("X_TEST_SEED"),
		os.Getenv("X_TEST_HASHED_PASSWORD"),
		os.Getenv("X_TEST_NATION_CODE"),
		os.Getenv("X_TEST_PHONE"),
		os.Getenv("X_TEST_EMAIL"),
		os.Getenv("X_TEST_PHONE_ERROR"),
		os.Getenv("X_TEST_PHONE_NULL"),
		os.Getenv("X_TEST_MVC_NULL"),
		os.Getenv("X_TEST_SERIALNUMBER_NULL"),
		os.Getenv("X_TEST_NATION_CODE_NULL"),
		os.Getenv("X_TEST_NATION_CODE_USA"),
		os.Getenv("X_TEST_SEED_NULL"),
		os.Getenv("X_TEST_PHONE_NOT_EXIST"),
		os.Getenv("X_TEST_HASHED_PASSWORD_NULL"),
		os.Getenv("X_TEST_EMAIL_NOT_EXIST"),
		os.Getenv("X_TEST_MVC_ERROR"),
		os.Getenv("X_TEST_EMAIL_NULL"),
		os.Getenv("X_TEST_EMAIL_ERROR"),
		os.Getenv("X_TEST_USERNAME_EXIST"),
		os.Getenv("X_TEST_USERNAME_NULL"),
		os.Getenv("X_TEST_QUESTIONKEY1"),
		os.Getenv("X_TEST_QUESTIONKEY2"),
		os.Getenv("X_TEST_QUESTIONKEY3"),
		os.Getenv("X_TEST_ANSWER1"),
		os.Getenv("X_TEST_ANSWER2"),
		os.Getenv("X_TEST_ANSWER3"),
		os.Getenv("X_TEST_QUESTIONKEY1_NEW"),
		os.Getenv("X_TEST_QUESTIONKEY2_NEW"),
		os.Getenv("X_TEST_QUESTIONKEY3_NEW"),
		os.Getenv("X_TEST_ANSWER1_NEW"),
		os.Getenv("X_TEST_ANSWER2_NEW"),
		os.Getenv("X_TEST_ANSWER3_NEW"),
		os.Getenv("X_TEST_QUESTIONKEY1_SAME"),
		os.Getenv("X_TEST_QUESTIONKEY2_SAME"),
		os.Getenv("X_TEST_QUESTIONKEY3_SAME"),
		os.Getenv("X_TEST_ANSWER1_SAME"),
		os.Getenv("X_TEST_ANSWER2_SAME"),
		os.Getenv("X_TEST_ANSWER3_SAME"),
		os.Getenv("X_TEST_QUESTIONKEY1_NULL"),
		os.Getenv("X_TEST_QUESTIONKEY2_NULL"),
		os.Getenv("X_TEST_QUESTIONKEY3_NULL"),
		os.Getenv("X_TEST_ANSWER1_NULL"),
		os.Getenv("X_TEST_ANSWER2_NULL"),
		os.Getenv("X_TEST_ANSWER3_NULL"),
		os.Getenv("X_TEST_ANSWER1_LONG"),
		os.Getenv("X_TEST_ANSWER1_FORMAT"),
		os.Getenv("X_TEST_QUESTIONKEY1_NOEXIST"),
		os.Getenv("X_TEST_NICKNAME"),
		protoGender,
		int32(weight),
		int32(height),
		os.Getenv("X_TEST_VERIFICATION_TYPE"),
		os.Getenv("X_TEST_NEW_EMAIL"),
		os.Getenv("X_TEST_REMARK"),
		os.Getenv("X_TEST_CUSTOMISED_CODE"),
		os.Getenv("X_TEST_PHONE_PASSWORD"),
		int32(userIDNotExist),
		os.Getenv("X_TEST_VT"),
		int32(weightNull),
		int32(weightLow),
		int32(weightHigh),
		int32(heightNull),
		int32(heightLow),
		int32(heightHigh),
		os.Getenv("X_TEST_MVC_FORMAT_ERROR"),
		os.Getenv("X_TEST_PHONE_FORMAT_ERROR"),
		os.Getenv("X_TEST_NATION_CODE_HK"),
		os.Getenv("X_TEST_PHONE_HK"),
		os.Getenv("X_TEST_PHONE_US"),
		os.Getenv("X_TEST_PHONE_TW"),
		os.Getenv("X_TEST_NATION_CODE_TW"),
		os.Getenv("X_TEST_PHONE_MACAO"),
		os.Getenv("X_TEST_NATION_CODE_MACAO"),
		os.Getenv("X_TEST_PHONE_CANADA"),
		os.Getenv("X_TEST_PHONE_UK"),
		os.Getenv("X_TEST_NATION_CODE_UK"),
		os.Getenv("X_TEST_PHONE_JP"),
		os.Getenv("X_TEST_NATION_CODE_JP"),
	}
}

// TestUserSetPassword 测试设置密码
func (suite *UserTestSuite) TestUserSetPassword() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	err = suite.JinmuIDService.UserSetPassword(ctx, &jinmuidpb.UserSetPasswordRequest{
		UserId:        userID,
		PlainPassword: suite.Account.password,
	}, nil)
	assert.NoError(t, err)
	assert.Error(t, errors.New("[errcode:1800] user password already exists"), err)
}

// TestUserModifyPassword 测试修改密码
func (suite *UserTestSuite) TestUserModifyPassword() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	err = suite.JinmuIDService.UserModifyPassword(ctx, &jinmuidpb.UserModifyPasswordRequest{
		UserId:            userID,
		NewPlainPassword:  suite.Account.password,
		OldHashedPassword: suite.Account.hashedPassword,
		Seed:              suite.Account.seed,
	}, nil)

	assert.Error(t, errors.New("[errcode:16000] new password cannot equals old password"), err)
}

// UserGetUsingServiceUserIdIsError 用户正在使用的服务
func (suite *UserTestSuite) TestUserGetUsingService() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	err = suite.JinmuIDService.UserGetUsingService(ctx, &jinmuidpb.UserGetUsingServiceRequest{
		UserId: userID,
	}, nil)
	assert.NoError(t, err)
}

// UserGetUsingServiceUserIdIsError 用户正在使用的服务UserId为空
func (suite *UserTestSuite) TestUserGetUsingServiceUserIdIsError() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	err = suite.JinmuIDService.UserGetUsingService(ctx, &jinmuidpb.UserGetUsingServiceRequest{
		UserId: userID,
	}, nil)
	assert.NoError(t, err)
}

func (suite *UserTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

// mapTestGenderToProto 将测试使用的 gender 转换为  proto 类型
func mapTestGenderToProto(gender string) (generalpb.Gender, error) {
	switch gender {
	case "M":
		return generalpb.Gender_GENDER_FEMALE, nil
	case "F":
		return generalpb.Gender_GENDER_MALE, nil
	}
	return generalpb.Gender_GENDER_INVALID, fmt.Errorf("invalid mysql gender %s", gender)
}
