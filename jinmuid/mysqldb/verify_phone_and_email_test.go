package mysqldb

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// VerifyPhoneAndEmailTestSuite 是 验证手机和邮箱 的 testSuite
type VerifyPhoneAndEmailTestSuite struct {
	suite.Suite
	db *DbClient
}

// SetupSuite 准备设置 Test Suite 执行
func (suite *VerifyPhoneAndEmailTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	db, _ := newTestingDbClientFromEnvFile(envFilepath)
	suite.db = db
}

// VerifyVerificationNumber 验证 VerificationNumber是否有效
func (suite *VerifyPhoneAndEmailTestSuite) TestVerifyVerificationNumber() {
	t := suite.T()
	ctx := context.Background()
	verificationType := VerificationPhone
	verificationNumber := os.Getenv("X_TEST_VERIFICATION_NUMBER")
	userID, _ := strconv.Atoi(os.Getenv("X_TEST_USER_ID"))
	isValid, err := suite.db.GetDB(ctx).VerifyVerificationNumber(ctx, verificationType, verificationNumber, int32(userID))
	assert.NoError(t, err)
	assert.Equal(t, false, isValid)
}

// SetVerificationNumberAsUsed 设置VerificationNumber已经使用
func (suite *VerifyPhoneAndEmailTestSuite) TestSetVerificationNumberAsUsed() {
	t := suite.T()
	ctx := context.Background()
	verificationType := VerificationPhone
	verificationNumber := os.Getenv("X_TEST_VERIFICATION_NUMBER")
	err := suite.db.GetDB(ctx).SetVerificationNumberAsUsed(ctx, verificationType, verificationNumber)
	assert.NoError(t, err)
}

// TestVerifyVerificationNumberByPhone 测试 手机号验证 VerificationNumber是否有效
func (suite *VerifyPhoneAndEmailTestSuite) TestVerifyVerificationNumberByPhone() {
	t := suite.T()
	ctx := context.Background()
	var verificationNumber = os.Getenv("X_TEST_VERIFICATION_NUMBER")
	var phone = os.Getenv("X_TEST_PHONE")
	var nationCode = os.Getenv("X_TEST_NATION_CODE")
	isValid, err := suite.db.GetDB(ctx).VerifyVerificationNumberByPhone(ctx, verificationNumber, phone, nationCode)
	assert.NoError(t, err)
	assert.Equal(t, false, isValid)
}

// VerifyVerificationNumberByEmail 邮箱 VerificationNumber是否有效
func (suite *VerifyPhoneAndEmailTestSuite) TestVerifyVerificationNumberByEmail() {
	t := suite.T()
	ctx := context.Background()
	var verificationNumber = os.Getenv("X_TEST_VERIFICATION_NUMBER")
	var email = os.Getenv("X_TEST_EMAIL")
	isValid, err := suite.db.GetDB(ctx).VerifyVerificationNumberByEmail(ctx, verificationNumber, email)
	assert.NoError(t, err)
	assert.Equal(t, false, isValid)
}
func TestVerifyPhoneAndEmailTestSuite(t *testing.T) {
	suite.Run(t, new(VerifyPhoneAndEmailTestSuite))
}
