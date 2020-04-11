package rest

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	generalpb "github.com/jinmukeji/proto/v3/gen/micro/idl/ptypes/v2"
	jinmuidpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	"github.com/kataras/iris/v12"
)

const (
	// GenderMale 男性
	GenderMale = 0
	// GenderFemale 女性
	GenderFemale = 1
	// GenderInvalid 非法的性别
	GenderInvalid = -1
)

// User 用户
type User struct {
	UserProfile
	HasSetEmail           bool   `json:"has_set_email"`
	HasSetPhone           bool   `json:"has_set_phone"`
	HasSetUsername        bool   `json:"has_set_username"`
	HasSetPassword        bool   `json:"has_set_password"`
	HasSetRegion          bool   `json:"has_set_region"`
	HasSetSecureQuestions bool   `json:"has_set_secure_questions"`
	HasSetUserProfile     bool   `json:"has_set_user_profile"`
	SigninUsername        string `json:"signin_username"`
	SecureEmail           string `json:"secure_email"`
	SigninPhone           string `json:"signin_phone"`
	Region                string `json:"region"`
	Language              string `json:"language"`
	IsProfileCompleted    bool   `json:"is_profile_completed"`
}

// UserProfile 用户档案
type UserProfile struct {
	Nickname string    `json:"nickname"` // 用户名
	Gender   int32     `json:"gender"`   // 用户性别 0 男 1 女
	Birthday time.Time `json:"birthday"` // 用户生日
	Height   int32     `json:"height"`   // 用户身高 单位厘米
	Weight   int32     `json:"weight"`   // 用户体重 单位千克
}

// GetUserProfile 得到UserProfile
func (h *webHandler) GetUserProfile(ctx iris.Context) {
	userID, err := ctx.Params().GetInt("user_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	req := new(jinmuidpb.GetUserProfileRequest)
	req.UserId = int32(userID)
	resp, errGetUserProfile := h.rpcSvc.GetUserProfile(newRPCContext(ctx), req)
	if errGetUserProfile != nil {
		writeRpcInternalError(ctx, errGetUserProfile, false)
		return
	}
	birthday, _ := ptypes.Timestamp(resp.Profile.BirthdayTime)
	protoGender, errmapProtoGenderToRest := mapProtoGenderToRest(resp.Profile.Gender)
	if errmapProtoGenderToRest != nil {
		writeRpcInternalError(ctx, errmapProtoGenderToRest, false)
	}
	stringLanguage, ermapProtoLanguageToRest := mapProtoLanguageToRest(resp.Language)
	if ermapProtoLanguageToRest != nil {
		// 默认使用简体中文
		stringLanguage = SimpleChinese
	}
	rest.WriteOkJSON(ctx, User{
		UserProfile: UserProfile{
			Nickname: resp.Profile.Nickname,
			Gender:   protoGender,
			Birthday: birthday.UTC(),
			Height:   resp.Profile.Height,
			Weight:   resp.Profile.Weight,
		},
		HasSetEmail:           resp.HasSetEmail,
		HasSetPassword:        resp.HasSetPassword,
		HasSetRegion:          resp.HasSetRegion,
		HasSetUsername:        resp.HasSetUsername,
		HasSetPhone:           resp.HasSetPhone,
		HasSetSecureQuestions: resp.HasSetSecureQuestions,
		HasSetUserProfile:     resp.HasSetUserProfile,
		SecureEmail:           resp.SecureEmail,
		SigninPhone:           resp.SigninPhone,
		SigninUsername:        resp.SigninUsername,
		Region:                resp.Region,
		Language:              stringLanguage,
		IsProfileCompleted:    resp.IsProfileCompleted,
	})
}

// ModifyUserProfile 修改UserProfile
func (h *webHandler) ModifyUserProfile(ctx iris.Context) {
	userID, err := ctx.Params().GetInt("user_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	req := new(jinmuidpb.ModifyUserProfileRequest)
	req.UserId = int32(userID)
	var profile UserProfile
	errReadJSON := ctx.ReadJSON(&profile)
	if errReadJSON != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", errReadJSON), false)
		return
	}
	errValidateUserProfile := ValidateUserProfile(profile)
	if errValidateUserProfile != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errValidateUserProfile), false)
		return
	}
	birthday, _ := ptypes.TimestampProto(profile.Birthday)
	protoGender, errmapRestGenderToProto := mapRestGenderToProto(profile.Gender)
	if errmapRestGenderToProto != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errmapRestGenderToProto), false)
		return
	}
	req.UserProfile = &jinmuidpb.UserProfile{
		Nickname:     profile.Nickname,
		Gender:       protoGender,
		BirthdayTime: birthday,
		Weight:       profile.Weight,
		Height:       profile.Height,
	}
	resp, errGetUserProfile := h.rpcSvc.ModifyUserProfile(newRPCContext(ctx), req)
	if errGetUserProfile != nil {
		writeRpcInternalError(ctx, errGetUserProfile, false)
		return
	}
	int32Gender, errmapProtoGenderToRest := mapProtoGenderToRest(resp.UserProfile.Gender)
	if errmapProtoGenderToRest != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errmapProtoGenderToRest), false)
		return
	}
	birthdayResp, _ := ptypes.Timestamp(resp.UserProfile.BirthdayTime)
	rest.WriteOkJSON(ctx, UserProfile{
		Nickname: resp.UserProfile.Nickname,
		Gender:   int32Gender,
		Birthday: birthdayResp.UTC(),
		Height:   resp.UserProfile.Height,
		Weight:   resp.UserProfile.Weight,
	})
}

// mapProtoGenderToRest 将 proto 类型的gender 转换为 Rest 的 int32 类型
func mapProtoGenderToRest(gender generalpb.Gender) (int32, error) {
	switch gender {
	case generalpb.Gender_GENDER_INVALID:
		return GenderInvalid, fmt.Errorf("invalid proto gender %d", generalpb.Gender_GENDER_INVALID)
	case generalpb.Gender_GENDER_UNSET:
		return GenderInvalid, fmt.Errorf("invalid proto gender %d", generalpb.Gender_GENDER_UNSET)
	case generalpb.Gender_GENDER_MALE:
		return GenderMale, nil
	case generalpb.Gender_GENDER_FEMALE:
		return GenderFemale, nil
	}
	return GenderInvalid, fmt.Errorf("invalid proto gender")
}

// mapRestGenderToProto 将 Rest 的 int32 类型的gender 转换为 proto 类型
func mapRestGenderToProto(gender int32) (generalpb.Gender, error) {
	switch gender {
	case GenderMale:
		return generalpb.Gender_GENDER_MALE, nil
	case GenderFemale:
		return generalpb.Gender_GENDER_FEMALE, nil
	}
	return generalpb.Gender_GENDER_INVALID, fmt.Errorf("invalid int32 gender %d", gender)
}

// ValidateUserProfile 验证UserProfile 参数是否合法
func ValidateUserProfile(profile UserProfile) error {
	if profile.Nickname == "" {
		return errors.New("nickname should not be empty")
	}
	if profile.Gender > 1 || profile.Gender < 0 {
		return fmt.Errorf("invalid gender %d", profile.Gender)
	}
	if profile.Birthday.IsZero() {
		return errors.New("birthday should not be null")
	}
	year1900 := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	now := time.Now()
	if !profile.Birthday.IsZero() && (profile.Birthday.Before(year1900) || profile.Birthday.After(now)) {
		return fmt.Errorf("invalid birthday %v", profile.Birthday)
	}
	if profile.Weight < 30 || profile.Weight > 500 {
		return fmt.Errorf("invalid weight %d", profile.Weight)
	}

	if profile.Height < 50 || profile.Height > 250 {
		return fmt.Errorf("invalid height %d", profile.Height)
	}

	return nil
}
