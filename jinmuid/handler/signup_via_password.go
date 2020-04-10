package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/gf-api2/jinmuid/mysqldb"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
)

const (
	// DefaultSignUpEncryptionSeed 用于加密密码的Seed，暂时定为空
	DefaultSignUpEncryptionSeed = ""
)

// SignUpUserViaUsernamePassword 通过密码注册新用户
func (j *JinmuIDService) SignUpUserViaUsernamePassword(ctx context.Context, req *proto.SignUpUserViaUsernamePasswordRequest, resp *proto.SignUpUserViaUsernamePasswordResponse) error {
	if !req.IsSkipVerifyToken {
		token, ok := TokenFromContext(ctx)
		if !ok {
			return NewError(ErrInvalidUser, errors.New("failed to get token from context"))
		}
		_, err := j.datastore.FindUserIDByToken(ctx, token)
		if err != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to get userID by token: %s", err.Error()))
		}
	}
	now := time.Now()
	u := &mysqldb.User{
		Remark:             req.Remark,
		RegisterType:       req.RegisterType,
		CustomizedCode:     req.CustomizedCode,
		Zone:               req.Zone,
		UserDefinedCode:    req.UserDefinedCode,
		RegisterTime:       now,
		HasSetUserProfile:  false,
		CreatedAt:          now,
		UpdatedAt:          now,
		IsProfileCompleted: true,
		IsActivated:        true,
		ActivatedAt:        &now,
	}
	if req.IsNeedSetProfile {
		u.IsProfileCompleted = false
		u.HasSetUserProfile = false
	}
	birthday, _ := ptypes.Timestamp(req.Profile.BirthdayTime)
	mysqlGender, errmapProtoGenderToDB := mapProtoGenderToDB(req.Profile.Gender)
	if errmapProtoGenderToDB != nil {
		return NewError(ErrInvalidGender, errmapProtoGenderToDB)
	}
	userProfile := &mysqldb.UserProfile{
		Nickname:        req.Profile.Nickname,
		NicknameInitial: getNicknameInitial(req.Profile.Nickname),
		Gender:          mysqlGender,
		Height:          req.Profile.Height,
		Weight:          req.Profile.Weight,
		Birthday:        birthday,
	}
	err := j.datastore.CreateUserAndUserProfile(ctx, u, userProfile)
	if err != nil {
		return NewError(ErrDatabase, errors.New("failed to create user and user_profile"))
	}
	resp.UserId = u.UserID
	return nil
}
