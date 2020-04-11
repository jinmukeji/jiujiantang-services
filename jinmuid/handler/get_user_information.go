package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
)

// GetUserAndProfileInformation 得到用户和用户档案信息
func (j *JinmuIDService) GetUserAndProfileInformation(ctx context.Context, req *proto.GetUserAndProfileInformationRequest, resp *proto.GetUserAndProfileInformationResponse) error {
	// 如果需要token验证
	if !req.IsSkipVerifyToken {
		token, ok := TokenFromContext(ctx)
		if !ok {
			return NewError(ErrInvalidUser, errors.New("failed to get token from context"))
		}
		userID, err := j.datastore.FindUserIDByToken(ctx, token)
		if err != nil {
			return NewError(ErrUserUnauthorized, fmt.Errorf("failed to find userID by token: %s", err.Error()))
		}
		if userID != req.UserId {
			return NewError(ErrInvalidUser, fmt.Errorf("user %d from request and user %d from token are inconsistent", req.UserId, userID))
		}
	}
	userInformation, errFindUserByUserID := j.datastore.FindUserByUserID(ctx, req.UserId)
	if errFindUserByUserID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find user by userID %d: %s", req.UserId, errFindUserByUserID.Error()))
	}
	if userInformation.HasSetUsername {
		resp.SigninUsername = userInformation.SigninUsername
	}
	if userInformation.HasSetEmail {
		resp.SecureEmail = userInformation.SecureEmail
	}
	if userInformation.HasSetPhone {
		resp.SigninPhone = userInformation.SigninPhone
	}
	resp.HasSetUserProfile = userInformation.HasSetUserProfile

	birthday, _ := ptypes.TimestampProto(userInformation.Birthday)
	protoGender, errmapDBGenderToProto := mapDBGenderToProto(userInformation.Gender)
	if errmapDBGenderToProto != nil {
		return NewError(ErrInvalidGender, errmapDBGenderToProto)
	}
	resp.Profile = &proto.UserProfile{
		Nickname:        userInformation.Nickname,
		BirthdayTime:    birthday,
		Gender:          protoGender,
		Weight:          userInformation.Weight,
		Height:          userInformation.Height,
		NicknameInitial: userInformation.NicknameInitial,
	}

	resp.Remark = userInformation.Remark
	resp.RegisterType = userInformation.RegisterType
	resp.CustomizedCode = userInformation.CustomizedCode
	resp.UserDefinedCode = userInformation.UserDefinedCode
	resp.RegisterTime, _ = ptypes.TimestampProto(userInformation.RegisterTime)
	resp.HasSetUserProfile = userInformation.HasSetUserProfile
	resp.IsProfileCompleted = userInformation.IsProfileCompleted
	resp.IsRemovable = userInformation.IsRemovable
	return nil
}
