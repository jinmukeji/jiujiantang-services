package mysqldb

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

// UserTestSuite 是 User 的 testSuite
type UserTestSuite struct {
	suite.Suite
	db *DbClient
}

// SetupSuite 准备设置 Test Suite 执行
func (suite *UserTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.db, _ = newTestingDbClientFromEnvFile(envFilepath)
}

// FindUserByPhone 通过电话找到用户
func (suite *UserTestSuite) TestFindUserByPhone() {
	t := suite.T()
	phone := os.Getenv("X_TEST_PHONE")
	nationCode := os.Getenv("X_TEST_NATION_CODE")
	ctx := context.Background()
	username, err := suite.db.FindUserByPhone(ctx, phone, nationCode)
	assert.NoError(t, err)
	assert.NotNil(t, username)
}

// FindUserByUsername 通过用户名找到base64密码
func (suite *UserTestSuite) TestFindUserByUsername() {
	t := suite.T()
	ctx := context.Background()
	_, err := suite.db.FindUserByUsername(ctx, os.Getenv("X_TEST_USERNAME"))
	assert.NoError(t, err)
}

// SetLanguageByUserID 通过userID设置Language
func (suite *UserTestSuite) TestSetLanguageByUserID() {
	t := suite.T()
	userID, _ := strconv.Atoi(os.Getenv("X_TEST_USER_ID"))
	ctx := context.Background()
	err := suite.db.SetLanguageByUserID(ctx, int32(userID), os.Getenv("X_TEST_LANGUAGE"))
	assert.NoError(t, err)
}

// FindLanguageByUserID 通过userID找到Language
func (suite *UserTestSuite) TestFindLanguageByUserID() {
	t := suite.T()
	userID, _ := strconv.Atoi(os.Getenv("X_TEST_USER_ID"))
	ctx := context.Background()
	_, err := suite.db.FindLanguageByUserID(ctx, int32(userID))
	assert.NoError(t, err)
}

// ExistUserByUserID 查看 user 能否存在
func (suite *UserTestSuite) TestExistUserByUserID() {
	t := suite.T()
	userID, _ := strconv.Atoi(os.Getenv("X_TEST_USER_ID"))
	ctx := context.Background()
	_, err := suite.db.ExistUserByUserID(ctx, int32(userID))
	assert.NoError(t, err)
}

// ExistPasswordByUserID 查看 password 能否存在
func (suite *UserTestSuite) TestExistPasswordByUserID() {
	t := suite.T()
	userID, _ := strconv.Atoi(os.Getenv("X_TEST_USER_ID"))
	ctx := context.Background()
	password, err := suite.db.ExistPasswordByUserID(ctx, int32(userID))
	assert.NoError(t, err)
	assert.Equal(t, false, password)
}

// SetPasswordByUserID 通过userID设置密码
func (suite *UserTestSuite) TestSetPasswordByUserID() {
	t := suite.T()
	userID, _ := strconv.Atoi(os.Getenv("X_TEST_USER_ID"))
	ctx := context.Background()
	err := suite.db.SetPasswordByUserID(ctx, int32(userID), os.Getenv("X_TEST_ENCRYPTEDPASSWORD"), os.Getenv("X_TEST_SEED"))
	assert.NoError(t, err)
}

// FindSecureQuestionByUserID 通过userID找到密保问题和答案
func (suite *UserTestSuite) TestFindSecureQuestionByUserID() {
	t := suite.T()
	userID, _ := strconv.Atoi(os.Getenv("X_TEST_USER_ID"))
	ctx := context.Background()
	_, err := suite.db.FindSecureQuestionByUserID(ctx, int32(userID))
	assert.NoError(t, err)
}

// FindSecureQuestionByPhone 通过电话号码找到找到密保问题和答案
func (suite *UserTestSuite) TestFindSecureQuestionByPhone() {
	t := suite.T()
	phone := os.Getenv("X_TEST_PHONE")
	nationCode := os.Getenv("X_TEST_NATION_CODE")
	ctx := context.Background()
	_, err := suite.db.FindSecureQuestionByPhone(ctx, phone, nationCode)
	assert.NoError(t, err)
}

// FindSecureQuestionByUsername 通过用户名找到密保问题和答案
func (suite *UserTestSuite) TestFindSecureQuestionByUsername() {
	t := suite.T()
	username := os.Getenv("X_TEST_USERNAME")
	ctx := context.Background()
	_, err := suite.db.FindSecureQuestionByUsername(ctx, username)
	assert.NoError(t, err)
}

// SetPasswordByPhone 根据手机号重置密码
func (suite *UserTestSuite) TestFindSetPasswordByPhone() {
	t := suite.T()
	phone := os.Getenv("X_TEST_PHONE")
	nationCode := os.Getenv("X_TEST_NATION_CODE")
	encryptedPassword := os.Getenv("X_TEST_ENCRYPTEDPASSWORD")
	seed := os.Getenv("X_TEST_SEED")
	ctx := context.Background()
	err := suite.db.SetPasswordByPhone(ctx, phone, nationCode, encryptedPassword, seed)
	assert.NoError(t, err)
}

// SetPasswordByUsername 根据用户名重置密码
func (suite *UserTestSuite) TestSetPasswordByUsername() {
	t := suite.T()
	username := os.Getenv("X_TEST_USERNAME")
	encryptedPassword := os.Getenv("X_TEST_ENCRYPTEDPASSWORD")
	seed := os.Getenv("X_TEST_SEED")
	ctx := context.Background()
	err := suite.db.SetPasswordByUsername(ctx, username, encryptedPassword, seed)
	assert.NoError(t, err)
}

// IsPasswordSameByPhone 根据手机号判断密码是否与之前密码相同
func (suite *UserTestSuite) TestIsPasswordSameByPhone() {
	t := suite.T()
	phone := os.Getenv("X_TEST_PHONE")
	nationCode := os.Getenv("X_TEST_NATION_CODE")
	encryptedPassword := os.Getenv("X_TEST_ENCRYPTEDPASSWORD")
	ctx := context.Background()
	_, err := suite.db.IsPasswordSameByPhone(ctx, phone, nationCode, encryptedPassword)
	assert.NoError(t, err)
}

// IsPasswordSameByUsername 根据用户名判断密码是否与之前密码相同
func (suite *UserTestSuite) TestIsPasswordSameByUsername() {
	t := suite.T()
	username := os.Getenv("X_TEST_USERNAME")
	encryptedPassword := os.Getenv("X_TEST_ENCRYPTEDPASSWORD")
	ctx := context.Background()
	_, err := suite.db.IsPasswordSameByUsername(ctx, username, encryptedPassword)
	assert.NoError(t, err)
}

// FindUserIDByPhone 通过电话号码找到userID
func (suite *UserTestSuite) TestFindUserIDByPhone() {
	t := suite.T()
	phone := os.Getenv("X_TEST_PHONE")
	nationCode := os.Getenv("X_TEST_NATION_CODE")
	ctx := context.Background()
	user, err := suite.db.FindUserIDByPhone(ctx, phone, nationCode)
	assert.NoError(t, err)
	assert.Equal(t, 105546, user)
}

// FindUserIDByUsername 通过用户名找到userID
func (suite *UserTestSuite) TestFindUserIDByUsername() {
	t := suite.T()
	username := os.Getenv("X_TEST_USERNAME")
	ctx := context.Background()
	user, err := suite.db.FindUserIDByUsername(ctx, username)
	assert.NoError(t, err)
	assert.Equal(t, int32(786), user)
}

// ExistUsername 用户名是否存在
func (suite *UserTestSuite) TestExistUsername() {
	t := suite.T()
	username := os.Getenv("X_TEST_USERNAME")
	ctx := context.Background()
	user, err := suite.db.ExistUsername(ctx, username)
	assert.NoError(t, err)
	assert.Equal(t, true, user)
}

// ExistPhone 手机号是否已经存在
func (suite *UserTestSuite) TestExistPhone() {
	t := suite.T()
	phone := os.Getenv("X_TEST_PHONE")
	nationCode := os.Getenv("X_TEST_NATION_CODE")
	ctx := context.Background()
	_, err := suite.db.ExistPhone(ctx, phone, nationCode)
	assert.NoError(t, err)
}

// ExistSignInPhone 登录手机号是否已经存在
func (suite *UserTestSuite) TestExistSignInPhone() {
	t := suite.T()
	phone := os.Getenv("X_TEST_PHONE")
	nationCode := os.Getenv("X_TEST_NATION_CODE")
	ctx := context.Background()
	_, err := suite.db.ExistSignInPhone(ctx, phone, nationCode)
	assert.NoError(t, err)
}

// SecureEmailExists 当前用户是否已经设置了安全邮箱
func (suite *UserTestSuite) TestSecureEmailExists() {
	t := suite.T()
	userID, _ := strconv.Atoi(os.Getenv("X_TEST_USER_ID"))
	ctx := context.Background()
	_, err := suite.db.SecureEmailExists(ctx, int32(userID))
	assert.NoError(t, err)
}

// MatchSecureEmail 安全邮箱是否与原来邮箱一致
func (suite *UserTestSuite) TestMatchSecureEmail() {
	t := suite.T()
	userID, _ := strconv.Atoi(os.Getenv("X_TEST_USER_ID"))
	email := os.Getenv("X_TEST_EMAIL")
	ctx := context.Background()
	_, err := suite.db.MatchSecureEmail(ctx, email, int32(userID))
	assert.NoError(t, err)
}

// CreateUserByPhone 创建user通过电话
func (suite *UserTestSuite) TestCreateUserByPhone() {
	t := suite.T()
	now := time.Now()
	user := &User{
		SigninPhone:    os.Getenv("X_TEST_PHONE"),
		HasSetPhone:    true,
		RegisterSource: os.Getenv("X_TEST_REGISTER_SOURCE"),
		NationCode:     os.Getenv("X_TEST_NATION_CODE"),
		RegisterTime:   now.UTC(),
		CreatedAt:      now.UTC(),
		UpdatedAt:      now.UTC(),
	}
	ctx := context.Background()
	_, err := suite.db.CreateUserByPhone(ctx, user)
	assert.NoError(t, err)
}

// SetSecureEmail 设置安全邮箱
func (suite *UserTestSuite) TestSetSecureEmail() {
	t := suite.T()
	userID, _ := strconv.Atoi(os.Getenv("X_TEST_USER_ID"))
	email := os.Getenv("X_TEST_EMAIL")
	ctx := context.Background()
	err := suite.db.SetSecureEmail(ctx, email, int32(userID))
	assert.NoError(t, err)
}

// UnsetSecureEmail 解除设置安全邮箱
func (suite *UserTestSuite) TestUnsetSecureEmail() {
	t := suite.T()
	userID, _ := strconv.Atoi(os.Getenv("X_TEST_USER_ID"))
	ctx := context.Background()
	err := suite.db.UnsetSecureEmail(ctx, int32(userID))
	assert.NoError(t, err)
}

// ExistsSecureQuestion 用户是否已经设置了密保问题
func (suite *UserTestSuite) TestExistsSecureQuestion() {
	t := suite.T()
	userID, _ := strconv.Atoi(os.Getenv("X_TEST_USER_ID"))
	ctx := context.Background()
	_, err := suite.db.ExistsSecureQuestion(ctx, int32(userID))
	assert.NoError(t, err)
}

// SetSecureQuestion 设置密保问题
func (suite *UserTestSuite) TestSetSecureQuestion() {
	t := suite.T()
	userID, _ := strconv.Atoi(os.Getenv("X_TEST_USER_ID"))
	secureQuestion := make([]SecureQuestion, 3)
	secureQuestion[0] = SecureQuestion{
		SecureQuestionKey: "1",
		SecureAnswer:      "1",
	}
	secureQuestion[1] = SecureQuestion{
		SecureQuestionKey: "1",
		SecureAnswer:      "1",
	}
	secureQuestion[2] = SecureQuestion{
		SecureQuestionKey: "1",
		SecureAnswer:      "1",
	}
	ctx := context.Background()
	err := suite.db.SetSecureQuestion(ctx, int32(userID), secureQuestion)
	assert.NoError(t, err)
}

// SetSigninPhoneByUserID 通过userID设置登录手机号
func (suite *UserTestSuite) TestSetSigninPhoneByUserID() {
	t := suite.T()
	userID, _ := strconv.Atoi(os.Getenv("X_TEST_USER_ID"))
	signinPhone := os.Getenv("X_TEST_PHONE")
	nationCode := os.Getenv("X_TEST_NATION_CODE")
	ctx := context.Background()
	err := suite.db.SetSigninPhoneByUserID(ctx, int32(userID), signinPhone, nationCode)
	assert.NoError(t, err)
}

// SetUserRegion 设置用户区域
func (suite *UserTestSuite) TestSetUserRegion() {
	t := suite.T()
	userID, _ := strconv.Atoi(os.Getenv("X_TEST_USER_ID"))
	region := MainlandChina
	ctx := context.Background()
	err := suite.db.SetUserRegion(ctx, int32(userID), region)
	assert.NoError(t, err)
}

// GetSecureQuestionListToModifyByUserID 通过userID找到密保问题
func (suite *UserTestSuite) TestGetSecureQuestionListToModifyByUserID() {
	t := suite.T()
	userID, _ := strconv.Atoi(os.Getenv("X_TEST_USER_ID"))
	ctx := context.Background()
	_, err := suite.db.GetSecureQuestionListToModifyByUserID(ctx, int32(userID))
	assert.NoError(t, err)
}

// FindUserBySecureEmail 通过安全邮箱找到User
func (suite *UserTestSuite) TestFindUserBySecureEmail() {
	t := suite.T()
	email := os.Getenv("X_TEST_EMAIL")
	ctx := context.Background()
	_, err := suite.db.FindUserBySecureEmail(ctx, email)
	assert.NoError(t, err)
}

// FindUsernameBySecureEmail 通过邮箱查找用户名
func (suite *UserTestSuite) TestFindUsernameBySecureEmail() {
	t := suite.T()
	email := os.Getenv("X_TEST_EMAIL")
	ctx := context.Background()
	_, err := suite.db.FindUsernameBySecureEmail(ctx, email)
	assert.NoError(t, err)
}

// GetSecureQuestionsByPhone 根据手机号获取当前设置的密保问题
func (suite *UserTestSuite) TestGetSecureQuestionsByPhone() {
	t := suite.T()
	phone := os.Getenv("X_TEST_PHONE")
	nationCode := os.Getenv("X_TEST_NATION_CODE")
	ctx := context.Background()
	_, err := suite.db.GetSecureQuestionsByPhone(ctx, phone, nationCode)
	assert.NoError(t, err)
}

// GetSecureQuestionsByUsername 根据用户名获取当前设置的密保问题
func (suite *UserTestSuite) TestGetSecureQuestionsByUsername() {
	t := suite.T()
	username := os.Getenv("X_TEST_USERNAME")
	ctx := context.Background()
	_, err := suite.db.GetSecureQuestionsByUsername(ctx, username)
	assert.NoError(t, err)
}

// FindUserByEmail 通过邮箱找到User
func (suite *UserTestSuite) TestFindUserByEmail() {
	t := suite.T()
	email := os.Getenv("X_TEST_EMAIL")
	ctx := context.Background()
	_, err := suite.db.FindUserByEmail(ctx, email)
	assert.NoError(t, err)
}

// SetSecureEmailByUserID 根据userID重置安全邮箱
func (suite *UserTestSuite) TestSetSecureEmailByUserID() {
	t := suite.T()
	email := os.Getenv("X_TEST_EMAIL")
	userID, _ := strconv.Atoi(os.Getenv("X_TEST_USER_ID"))
	ctx := context.Background()
	err := suite.db.SetSecureEmailByUserID(ctx, int32(userID), email)
	assert.NoError(t, err)
}

// ModifyHasSetUserProfileStatus 修改HasSetUserProfile状态
func (suite *UserTestSuite) TestModifyHasSetUserProfileStatus() {
	t := suite.T()
	userID, _ := strconv.Atoi(os.Getenv("X_TEST_USER_ID"))
	ctx := context.Background()
	err := suite.db.ModifyHasSetUserProfileStatus(ctx, int32(userID))
	assert.NoError(t, err)
}

// HasSecureEmailSet 当前安全邮箱是否被任何人设置
func (suite *UserTestSuite) TestHasSecureEmailSets() {
	t := suite.T()
	email := os.Getenv("X_TEST_EMAIL")
	ctx := context.Background()
	_, err := suite.db.HasSecureEmailSet(ctx, email)
	assert.NoError(t, err)
}

// FindUserByUserID 获取用户和用户档案信息
func (suite *UserTestSuite) TestFindUserByUserID() {
	t := suite.T()
	userID, _ := strconv.Atoi(os.Getenv("X_TEST_USER_ID"))
	ctx := context.Background()
	user, err := suite.db.FindUserByUserID(ctx, int32(userID))
	assert.NoError(t, err)
	assert.Equal(t, "13221058643", user.SigninPhone)
}

// ModifyUser 修改用户信息
func (suite *UserTestSuite) TestModifyUser() {
	t := suite.T()
	userID, _ := strconv.Atoi(os.Getenv("X_TEST_USER_ID"))
	user := &User{
		UserID:            int32(userID),
		Remark:            os.Getenv("X_TEST_REMARK"),
		CustomizedCode:    os.Getenv("X_TEST_CUSTOMIZED_CODE"),
		HasSetUserProfile: true,
		UpdatedAt:         time.Now().UTC(),
	}
	ctx := context.Background()
	err := suite.db.ModifyUser(ctx, user)
	assert.NoError(t, err)
}
func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
