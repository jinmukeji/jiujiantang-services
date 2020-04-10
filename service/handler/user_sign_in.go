package handler

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/jinmukeji/gf-api2/service/auth"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/go-pkg/age"

	corepb "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	jinmuidpb "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	subscriptionpb "github.com/jinmukeji/proto/gen/micro/idl/jm/subscription/v1"
)

const (
	// RegisterTypeUsername 用户名注册
	RegisterTypeUsername = "username"
)

// UserSignIn 用户登录
func (j *JinmuHealth) UserSignIn(ctx context.Context, req *corepb.UserSignInRequest, resp *corepb.UserSignInResponse) error {
	reqUserSignInByUsernamePassword := new(jinmuidpb.UserSignInByUsernamePasswordRequest)
	reqUserSignInByUsernamePassword.Username = req.SignInKey
	reqUserSignInByUsernamePassword.HashedPassword = req.PasswordHash
	reqUserSignInByUsernamePassword.Seed = ""
	reqUserSignInByUsernamePassword.Ip = req.Ip
	respUserSignInByUsernamePassword, errUserSignInByUsernamePassword := j.jinmuidSvc.UserSignInByUsernamePassword(ctx, reqUserSignInByUsernamePassword)
	if errUserSignInByUsernamePassword != nil {
		return errUserSignInByUsernamePassword
	}
	resp.AccessToken = respUserSignInByUsernamePassword.AccessToken
	resp.UserId = respUserSignInByUsernamePassword.UserId
	reqGetUserSubscriptions := new(subscriptionpb.GetUserSubscriptionsRequest)
	reqGetUserSubscriptions.UserId = respUserSignInByUsernamePassword.UserId
	ctx = auth.AddContextToken(ctx, respUserSignInByUsernamePassword.AccessToken)
	respGetUserSubscriptions, errGetUserSubscriptions := j.subscriptionSvc.GetUserSubscriptions(ctx, reqGetUserSubscriptions)
	if errGetUserSubscriptions != nil || len(respGetUserSubscriptions.Subscriptions) == 0 {
		resp.RemainDays = 0
		resp.ExpireTime = ptypes.TimestampNow()
		return nil
	}
	// 获取当前正在使用的订阅
	selectedSubscription := new(subscriptionpb.Subscription)
	for _, item := range respGetUserSubscriptions.Subscriptions {
		if item.IsSelected {
			selectedSubscription = item
		}
	}
	// 激活订阅
	if !selectedSubscription.Activated {
		reqActivateSubscription := new(subscriptionpb.ActivateSubscriptionRequest)
		reqActivateSubscription.SubscriptionId = selectedSubscription.SubscriptionId
		respActivateSubscription, errActivateSubscription := j.subscriptionSvc.ActivateSubscription(ctx, reqActivateSubscription)
		selectedSubscription.ExpiredTime = respActivateSubscription.ExpiredTime
		if errActivateSubscription != nil {
			return errActivateSubscription
		}
	}
	expiredAt, _ := ptypes.Timestamp(selectedSubscription.ExpiredTime)
	resp.RemainDays = getRemainDays(expiredAt.UTC())
	resp.ExpireTime = selectedSubscription.ExpiredTime
	resp.AccessTokenExpiredTime = respUserSignInByUsernamePassword.ExpiredTime
	return nil

}

// 得到剩余时间
func getRemainDays(expiredAt time.Time) int32 {
	return int32(math.Ceil(time.Until(expiredAt).Hours() / 24))
}

// GetUserByUserName 通过Username得到User
func (j *JinmuHealth) GetUserByUserName(ctx context.Context, req *corepb.GetUserByUserNameRequest, resp *corepb.GetUserByUserNameResponse) error {
	u, err := j.datastore.FindUserByUsername(ctx, req.Username)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find user by username %s: %s", req.Username, err.Error()))
	}
	birth, _ := ptypes.TimestampProto(u.Birthday)
	gender, errMapDBGenderToProto := mapDBGenderToProto(u.Gender)
	if errMapDBGenderToProto != nil {
		return NewError(ErrInvalidGender, errMapDBGenderToProto)
	}
	resp.User = &corepb.User{
		UserId:             int32(u.UserID),
		Username:           u.Username,
		RegisterType:       u.RegisterType,
		IsProfileCompleted: u.IsProfileCompleted,
		IsRemovable:        u.IsRemovable,
		Profile: &corepb.UserProfile{
			Nickname:        u.Nickname,
			BirthdayTime:    birth,
			Age:             int32(age.Age(u.Birthday)),
			Gender:          gender,
			Height:          int32(u.Height),
			Weight:          int32(u.Weight),
			Phone:           u.Phone,
			Email:           u.Email,
			UserDefinedCode: u.UserDefinedCode,
			Remark:          u.Remark,
			State:           u.State,
			City:            u.City,
			Street:          u.Street,
			Country:         u.Country,
		},
	}
	return nil
}
