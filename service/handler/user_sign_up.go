package handler

import (
	"context"
	"errors"
	"time"

	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/gf-api2/service/auth"
	corepb "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	jinmuidpb "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
)

// kangmeiClient 康美客户端
const kangmeiClient = "kangmei-10001"

// UserSignUp 注册用户
func (j *JinmuHealth) UserSignUp(ctx context.Context, req *corepb.UserSignUpRequest, resp *corepb.UserSignUpResponse) error {
	ownerID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return NewError(ErrInvalidUser, errors.New("invalid user context"))
	}
	// 判断当前的Client_id是否和传入的Client_id一致
	client, _ := clientFromContext(ctx)
	if req.ClientId != client.ClientID {
		return NewError(ErrIncorrectClientID, fmt.Errorf("client %s from request and client %s from context are inconsistent", req.ClientId, client.ClientID))
	}
	// 获取owner的信息
	reqGetUserAndProfileInformation := new(jinmuidpb.GetUserAndProfileInformationRequest)
	reqGetUserAndProfileInformation.UserId = ownerID
	respGetUserAndProfileInformation, errGetUserAndProfileInformation := j.jinmuidSvc.GetUserAndProfileInformation(ctx, reqGetUserAndProfileInformation)
	if errGetUserAndProfileInformation != nil {
		return errGetUserAndProfileInformation
	}

	// 注册用户和user_profile信息
	birthday, _ := ptypes.Timestamp(req.UserProfile.BirthdayTime)
	now := time.Now()
	reqSignUpUserViaUsernamePassword := new(jinmuidpb.SignUpUserViaUsernamePasswordRequest)
	reqSignUpUserViaUsernamePassword.ClientId = client.ClientID
	reqSignUpUserViaUsernamePassword.RegisterType = req.RegisterType
	reqSignUpUserViaUsernamePassword.CustomizedCode = respGetUserAndProfileInformation.CustomizedCode
	reqSignUpUserViaUsernamePassword.Remark = req.UserProfile.Remark
	reqSignUpUserViaUsernamePassword.UserDefinedCode = req.UserProfile.UserDefinedCode
	reqSignUpUserViaUsernamePassword.RegisterTime, _ = ptypes.TimestampProto(now)
	birthdayProto, _ := ptypes.TimestampProto(birthday)
	reqSignUpUserViaUsernamePassword.Profile = &jinmuidpb.UserProfile{
		Nickname:     req.UserProfile.Nickname,
		BirthdayTime: birthdayProto,
		Gender:       req.UserProfile.Gender,
		Height:       req.UserProfile.Height,
		Weight:       req.UserProfile.Weight,
	}
	respSignUpUserViaUsernamePassword, errSignUpUserViaUsernamePassword := j.jinmuidSvc.SignUpUserViaUsernamePassword(ctx, reqSignUpUserViaUsernamePassword)
	if errSignUpUserViaUsernamePassword != nil {
		return errSignUpUserViaUsernamePassword
	}

	resp.User = &corepb.User{
		RegisterType:       reqSignUpUserViaUsernamePassword.RegisterType,
		IsProfileCompleted: true,
		IsRemovable:        true,
		Profile: &corepb.UserProfile{
			Nickname:        req.UserProfile.Nickname,
			BirthdayTime:    req.UserProfile.BirthdayTime,
			Gender:          req.UserProfile.Gender,
			Height:          req.UserProfile.Height,
			Weight:          req.UserProfile.Weight,
			Phone:           reqSignUpUserViaUsernamePassword.SigninPhone,
			Email:           reqSignUpUserViaUsernamePassword.SecureEmail,
			Remark:          reqSignUpUserViaUsernamePassword.Remark,
			UserDefinedCode: reqSignUpUserViaUsernamePassword.UserDefinedCode,
		},
		UserId: int32(respSignUpUserViaUsernamePassword.UserId),
	}
	return nil
}
