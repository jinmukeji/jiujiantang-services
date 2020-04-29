package mysqldb

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"context"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// VCRecord 是 验证码的 testSuite
type VCRecordTestSuite struct {
	suite.Suite
	db *DbClient
}

// SetupSuite 准备设置 Test Suite 执行
func (suite *VCRecordTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	db, _ := newTestingDbClientFromEnvFile(envFilepath)
	suite.db = db
}

// SearchVcRecordCountsIn24hours 搜索24小时内的验证码记录个数
func (suite *VCRecordTestSuite) TestSearchVcRecordCountsIn24hours() {
	t := suite.T()
	ctx := context.Background()
	_, err := suite.db.GetDB(ctx).SearchVcRecordCountsIn24hours(ctx, os.Getenv("X_TEST_EMAIL"))
	assert.NoError(t, err)
}

// SearchVcRecordCountsIn1Minute 搜索1分钟内的验证码记录
func (suite *VCRecordTestSuite) TestSearchVcRecordCountsIn1Minute() {
	t := suite.T()
	ctx := context.Background()
	_, err := suite.db.GetDB(ctx).SearchVcRecordCountsIn1Minute(ctx, os.Getenv("X_TEST_EMAIL"))
	assert.NoError(t, err)
}

// SearchVcRecordEarliestTimeIn1Minute 搜索1分钟最早的验证码记录时间
func (suite *VCRecordTestSuite) TestSearchVcRecordEarliestTimeIn1Minute() {
	t := suite.T()
	ctx := context.Background()
	_, err := suite.db.GetDB(ctx).SearchVcRecordEarliestTimeIn1Minute(ctx, os.Getenv("X_TEST_EMAIL"))
	assert.NoError(t, err)
}

// FindVcRecord 查找验证码记录
func (suite *VCRecordTestSuite) TestFindVcRecord() {
	t := suite.T()
	ctx := context.Background()
	sn := os.Getenv("X_TEST_SN")
	vc := os.Getenv("X_TEST_VC")
	sendTo := os.Getenv("X_TEST_EMAIL")
	usage := SignUp
	_, err := suite.db.GetDB(ctx).FindVcRecord(ctx, sn, vc, sendTo, usage)
	assert.NoError(t, err)
}

// HasSnExpired 判断Sn是否已经过期
func (suite *VCRecordTestSuite) TestHasSnExpired() {
	t := suite.T()
	ctx := context.Background()
	sn := os.Getenv("X_TEST_SN")
	vc := os.Getenv("X_TEST_VC")
	_, err := suite.db.GetDB(ctx).HasSnExpired(ctx, sn, vc)
	assert.NoError(t, err)
}

// ModifyVcRecordStatus 修改验证码的状态
func (suite *VCRecordTestSuite) TestModifyVcRecordStatus() {
	t := suite.T()
	ctx := context.Background()
	recordID, _ := strconv.Atoi(os.Getenv("X_TEST_RECORD_ID"))
	err := suite.db.GetDB(ctx).ModifyVcRecordStatus(ctx, int32(recordID))
	assert.NoError(t, err)
}

// VerifyMVC 验证MVC
func (suite *VCRecordTestSuite) TestVerifyMVC() {
	t := suite.T()
	ctx := context.Background()
	sn := os.Getenv("X_TEST_SN")
	vc := os.Getenv("X_TEST_VC")
	sendTo := os.Getenv("X_TEST_EMAIL")
	nationCode := os.Getenv("X_TEST_NATION_CODE")
	_, err := suite.db.GetDB(ctx).VerifyMVC(ctx, sn, vc, sendTo, nationCode)
	assert.NoError(t, err)
}

// SearchVcRecord 查找验证码记录
func (suite *VCRecordTestSuite) TestSearchVcRecord() {
	t := suite.T()
	ctx := context.Background()
	sn := os.Getenv("X_TEST_SN")
	vc := os.Getenv("X_TEST_VC")
	sendTo := os.Getenv("X_TEST_EMAIL")
	nationCode := os.Getenv("X_TEST_NATION_CODE")
	_, err := suite.db.GetDB(ctx).SearchVcRecord(ctx, sn, vc, sendTo, nationCode)
	assert.NoError(t, err)
}

// FindLatestVcRecord 查找最新验证码记录
func (suite *VCRecordTestSuite) TestFindLatestVcRecord() {
	t := suite.T()
	ctx := context.Background()
	sendTo := os.Getenv("X_TEST_EMAIL")
	usage := SignUp
	_, err := suite.db.GetDB(ctx).FindLatestVcRecord(ctx, sendTo, usage)
	assert.NoError(t, err)
}

// SearchSpecificVcRecordCountsIn24hours 搜索24小时内的指定模板的验证码记录个数
func (suite *VCRecordTestSuite) TestSearchSpecificVcRecordCountsIn24hours() {
	t := suite.T()
	ctx := context.Background()
	sendTo := os.Getenv("X_TEST_EMAIL")
	usage := SignUp
	_, err := suite.db.GetDB(ctx).SearchSpecificVcRecordCountsIn24hours(ctx, sendTo, usage)
	assert.NoError(t, err)
}

// SearchSpecificVcRecordEarliestTimeIn24hours 搜索24小时指定模板最早的验证码记录
func (suite *VCRecordTestSuite) TestSearchSpecificVcRecordEarliestTimeIn24hours() {
	t := suite.T()
	ctx := context.Background()
	sendTo := os.Getenv("X_TEST_EMAIL")
	usage := SignUp
	_, err := suite.db.GetDB(ctx).SearchSpecificVcRecordEarliestTimeIn24hours(ctx, sendTo, usage)
	assert.NoError(t, err)
}

// SearchLatestPhoneVerificationCode 搜索最新的电话验证码
func (suite *VCRecordTestSuite) TestSearchLatestPhoneVerificationCode() {
	t := suite.T()
	ctx := context.Background()
	sendTo := os.Getenv("X_TEST_EMAIL")
	nationCode := os.Getenv("X_TEST_NATION_CODE")
	_, err := suite.db.GetDB(ctx).SearchLatestPhoneVerificationCode(ctx, sendTo, nationCode)
	assert.NoError(t, err)
}

// SearchLatestEmailVerificationCode 搜索最新的邮件验证码
func (suite *VCRecordTestSuite) TestSearchLatestEmailVerificationCode() {
	t := suite.T()
	ctx := context.Background()
	sendTo := os.Getenv("X_TEST_EMAIL")
	_, err := suite.db.GetDB(ctx).SearchLatestEmailVerificationCode(ctx, sendTo)
	assert.NoError(t, err)
}

// VerifyMVCBySecureEmail 根据安全邮箱验证MVC
func (suite *VCRecordTestSuite) TestVerifyMVCBySecureEmail() {
	t := suite.T()
	ctx := context.Background()
	sn := os.Getenv("X_TEST_SN")
	vc := os.Getenv("X_TEST_VC")
	email := os.Getenv("X_TEST_EMAIL")
	_, err := suite.db.GetDB(ctx).VerifyMVCBySecureEmail(ctx, sn, vc, email)
	assert.NoError(t, err)
}

// ModifyVcRecordStatusByEmail 根据安全邮箱修改验证码的状态
func (suite *VCRecordTestSuite) TestModifyVcRecordStatusByEmail() {
	t := suite.T()
	ctx := context.Background()
	verificationCode := os.Getenv("X_TEST_VC")
	serialNumber := os.Getenv("X_TEST_SN")
	email := os.Getenv("X_TEST_EMAIL")
	err := suite.db.GetDB(ctx).ModifyVcRecordStatusByEmail(ctx, email, verificationCode, serialNumber)
	assert.NoError(t, err)
}

// SetVcAsUsed 设置vc为使用过的
func (suite *VCRecordTestSuite) TestSetVcAsUsed() {
	t := suite.T()
	ctx := context.Background()
	sn := os.Getenv("X_TEST_SN")
	vc := os.Getenv("X_TEST_VC")
	email := os.Getenv("X_TEST_EMAIL")
	nationCode := os.Getenv("X_TEST_NATION_CODE")
	err := suite.db.GetDB(ctx).SetVcAsUsed(ctx, sn, vc, email, nationCode)
	assert.NoError(t, err)
}

func TestVCRecordTestSuite(t *testing.T) {
	suite.Run(t, new(VCRecordTestSuite))
}
