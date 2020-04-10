package handler

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	phone = "1"
	email = "2"
)

// UserResetPasswordTestSuite 修改安全邮箱测试
type UserResetPasswordTestSuite struct {
	suite.Suite
	JinmuIDService *JinmuIDService
	Account        *Account
}

// UserResetPasswordTestSuite 设置测试环境
func (suite *UserResetPasswordTestSuite) SetupSuite() {
	envFilepath := filepath.Join("testdata", "local.svc-jinmuid.env")
	suite.JinmuIDService = newJinmuIDServiceForTest()
	suite.Account = newTestingAccountFromEnvFile(envFilepath)
}

// TestUserResetPassword 修改密码
func (suite *UserResetPasswordTestSuite) TestUserResetPassword() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	// 获取验证码
	mvc := getEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	protoVerificationType, errmapDBProtoVerificationTypeToProto := mapDBProtoVerificationTypeToProto(suite.Account.vt)
	assert.NoError(t, errmapDBProtoVerificationTypeToProto)
	req := &proto.UserResetPasswordRequest{
		UserId:             userID,
		VerificationNumber: mvc,
		PlainPassword:      suite.Account.password,
		VerificationType:   protoVerificationType,
	}

	resp := new(proto.UserResetPasswordResponse)
	err = suite.JinmuIDService.UserResetPassword(ctx, req, resp)
	assert.NoError(t, err)
}

// TestUserResetPasswordVerificationIsInvalid 修改密码
func (suite *UserResetPasswordTestSuite) TestUserResetPasswordVerificationIsInvalid修改密码V() {
	t := suite.T()
	ctx := context.Background()
	ctx, userID, err := mockSigninByPhonePassword(ctx, suite.JinmuIDService, suite.Account.phone, suite.Account.phonePassword, suite.Account.seed, suite.Account.nationCode)
	assert.NoError(t, err)
	// 获取验证码
	mvc := getEmailVerificationCode(suite.JinmuIDService, *suite.Account)
	protoVerificationType, errmapDBProtoVerificationTypeToProto := mapDBProtoVerificationTypeToProto(suite.Account.vt)
	assert.NoError(t, errmapDBProtoVerificationTypeToProto)
	req := &proto.UserResetPasswordRequest{
		UserId:             userID,
		VerificationNumber: mvc,
		PlainPassword:      suite.Account.password,
		VerificationType:   protoVerificationType,
	}

	resp := new(proto.UserResetPasswordResponse)
	err = suite.JinmuIDService.UserResetPassword(ctx, req, resp)
	assert.Error(t, errors.New("[errcode:36000] verification number is invalid"), err)
}

// 读type
func mapDBProtoVerificationTypeToProto(verificationType string) (proto.VerificationType, error) {
	switch verificationType {
	case email:
		return proto.VerificationType_VERIFICATION_TYPE_EMAIL, nil
	case phone:
		return proto.VerificationType_VERIFICATION_TYPE_PHONE, nil
	}
	return proto.VerificationType_VERIFICATION_TYPE_INVALID, fmt.Errorf("invalid string verification type %s", verificationType)
}

func (suite *UserResetPasswordTestSuite) TearDownSuite() {
	ctx := context.Background()
	suite.JinmuIDService.datastore.SafeCloseDB(ctx)
}

func TestUserResetPasswordTestSuite(t *testing.T) {
	suite.Run(t, new(UserResetPasswordTestSuite))
}
