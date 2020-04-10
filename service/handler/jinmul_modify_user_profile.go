package handler

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/go-pkg/age"
	corepb "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	jinmuidpb "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
)

// JinmuLModifyUserProfile 修改用户个人档案
func (j *JinmuHealth) JinmuLModifyUserProfile(ctx context.Context, req *corepb.JinmuLModifyUserProfileRequest, resp *corepb.JinmuLModifyUserProfileResponse) error {
	// 获取user_profile信息
	reqGetUserProfile := new(jinmuidpb.GetUserProfileRequest)
	reqGetUserProfile.UserId = req.UserId
	reqGetUserProfile.IsSkipVerifyToken = true
	respGetUserProfile, errGetUserProfile := j.jinmuidSvc.GetUserProfile(ctx, reqGetUserProfile)
	if errGetUserProfile != nil {
		return errGetUserProfile
	}

	wxUserByUserID, errFindWXUserByUserID := j.datastore.FindWXUserByUserID(ctx, req.UserId)
	if errFindWXUserByUserID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to get wechat user by userID %d: %s", req.UserId, errFindWXUserByUserID.Error()))
	}
	// 更新用户信息
	reqModifyUserInformation := &jinmuidpb.ModifyUserInformationRequest{
		IsSkipVerifyToken: true,
		UserId:            req.UserId,
		SigninPhone:       req.UserProfile.Phone,
		SecureEmail:       req.UserProfile.Email,
		Remark:            req.UserProfile.Remark,
		CustomizedCode:    req.UserProfile.CustomizedCode,
		UserDefinedCode:   req.UserProfile.UserDefinedCode,
	}
	_, errModifyUserInformation := j.jinmuidSvc.ModifyUserInformation(ctx, reqModifyUserInformation)
	if errModifyUserInformation != nil {
		return errModifyUserInformation
	}
	// 更新user_prodile信息
	reqModifyUserProfile := new(jinmuidpb.ModifyUserProfileRequest)
	reqModifyUserProfile.UserId = req.UserId
	reqModifyUserProfile.IsSkipVerifyToken = true
	if !respGetUserProfile.IsProfileCompleted {
		reqModifyUserProfile.UserProfile = &jinmuidpb.UserProfile{
			Nickname:     wxUserByUserID.Nickname,
			BirthdayTime: req.UserProfile.BirthdayTime,
			Gender:       req.UserProfile.Gender,
			Height:       req.UserProfile.Height,
			Weight:       req.UserProfile.Weight,
		}
	} else {
		reqModifyUserProfile.UserProfile = &jinmuidpb.UserProfile{
			Nickname:     wxUserByUserID.Nickname,
			BirthdayTime: req.UserProfile.BirthdayTime,
			Height:       req.UserProfile.Height,
			Gender:       req.UserProfile.Gender,
			Weight:       req.UserProfile.Weight,
		}
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
		UserId:             req.UserId,
		Username:           respGetUserAndProfileInformation.SigninUsername,
		RegisterType:       respGetUserAndProfileInformation.RegisterType,
		IsProfileCompleted: respGetUserAndProfileInformation.IsProfileCompleted,
		Profile: &corepb.UserProfile{
			Nickname:        respGetUserAndProfileInformation.Profile.Nickname,
			Gender:          respGetUserAndProfileInformation.Profile.Gender,
			Age:             int32(age.Age(birthday)),
			Height:          respGetUserAndProfileInformation.Profile.Height,
			Weight:          respGetUserAndProfileInformation.Profile.Weight,
			BirthdayTime:    respGetUserAndProfileInformation.Profile.BirthdayTime,
			UserDefinedCode: respGetUserAndProfileInformation.UserDefinedCode,
			Remark:          respGetUserAndProfileInformation.Remark,
		},
	}
	return nil
}
