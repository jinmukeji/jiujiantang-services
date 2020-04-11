package handler

import (
	"context"

	"github.com/jinmukeji/go-pkg/age"

	"github.com/golang/protobuf/ptypes"
	corepb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	jinmuidpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
)

// JinmuLGetUserProfile 查看用户个人档案
func (j *JinmuHealth) JinmuLGetUserProfile(ctx context.Context, req *corepb.JinmuLGetUserProfileRequest, resp *corepb.JinmuLGetUserProfileResponse) error {
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
