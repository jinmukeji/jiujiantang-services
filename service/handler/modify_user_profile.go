package handler

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/go-pkg/age"
	jinmuidpb "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	corepb "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
)

// ModifyUserProfile 修改用户个人档案
func (j *JinmuHealth) ModifyUserProfile(ctx context.Context, req *corepb.ModifyUserProfileRequest, resp *corepb.ModifyUserProfileResponse) error {
	reqGetUserProfile := new(jinmuidpb.GetUserProfileRequest)
	reqGetUserProfile.UserId = req.UserId
	reqGetUserProfile.IsSkipVerifyToken = true
	respGetUserProfile, errGetUserProfile := j.jinmuidSvc.GetUserProfile(ctx, reqGetUserProfile)
	if errGetUserProfile != nil {
		return errGetUserProfile
	}

	// 更新用户信息
	reqModifyUserInformation := &jinmuidpb.ModifyUserInformationRequest{
		UserId:          req.UserId,
		SigninPhone:     req.UserProfile.Phone,
		SecureEmail:     req.UserProfile.Email,
		Remark:          req.UserProfile.Remark,
		CustomizedCode:  req.UserProfile.CustomizedCode,
		UserDefinedCode: req.UserProfile.UserDefinedCode,
	}
	reqModifyUserInformation.IsSkipVerifyToken = true
	// TODO：ModifyUserInformation，ModifyUserProfile 要是同一个事务
	_, errModifyUserInformation := j.jinmuidSvc.ModifyUserInformation(ctx, reqModifyUserInformation)
	if errModifyUserInformation != nil {
		return errModifyUserInformation
	}
	// 更新user_prodile信息
	reqModifyUserProfile := new(jinmuidpb.ModifyUserProfileRequest)
	reqModifyUserProfile.IsSkipVerifyToken = true
	reqModifyUserProfile.UserId = req.UserId
	// 性别不可更改
	reqModifyUserProfile.UserProfile = &jinmuidpb.UserProfile{
		Nickname:     req.UserProfile.Nickname,
		BirthdayTime: req.UserProfile.BirthdayTime,
		Height:       req.UserProfile.Height,
		Gender:       respGetUserProfile.Profile.Gender,
		Weight:       req.UserProfile.Weight,
	}
	_, errModifyUserProfile := j.jinmuidSvc.ModifyUserProfile(ctx, reqModifyUserProfile)
	if errModifyUserProfile != nil {
		return errModifyUserProfile
	}

	reqGetUserAndProfileInformation := new(jinmuidpb.GetUserAndProfileInformationRequest)
	reqGetUserAndProfileInformation.UserId = req.UserId
	reqGetUserAndProfileInformation.IsSkipVerifyToken = true
	respGetUserAndProfileInformation, errGetUserAndProfileInformation := j.jinmuidSvc.GetUserAndProfileInformation(ctx, reqGetUserAndProfileInformation)
	if errGetUserAndProfileInformation != nil {
		return errGetUserAndProfileInformation
	}
	birthday, _ := ptypes.Timestamp(respGetUserAndProfileInformation.Profile.BirthdayTime)
	resp.User = &corepb.User{
		UserId:             int32(req.UserId),
		Username:           respGetUserAndProfileInformation.SigninUsername,
		RegisterType:       respGetUserAndProfileInformation.RegisterType,
		IsProfileCompleted: respGetUserAndProfileInformation.HasSetUserProfile,
		IsRemovable:        respGetUserAndProfileInformation.IsRemovable,
		Profile: &corepb.UserProfile{
			Nickname:        respGetUserAndProfileInformation.Profile.Nickname,
			Gender:          respGetUserAndProfileInformation.Profile.Gender,
			Age:             int32(age.Age(birthday)),
			Height:          respGetUserAndProfileInformation.Profile.Height,
			Weight:          respGetUserAndProfileInformation.Profile.Weight,
			BirthdayTime:    respGetUserAndProfileInformation.Profile.BirthdayTime,
			Email:           respGetUserAndProfileInformation.SecureEmail,
			Phone:           respGetUserAndProfileInformation.SigninPhone,
			Remark:          respGetUserAndProfileInformation.Remark,
			UserDefinedCode: respGetUserAndProfileInformation.UserDefinedCode,
		},
	}
	return nil
}
