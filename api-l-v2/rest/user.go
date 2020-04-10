package rest

import (
	"fmt"
	"time"

	"errors"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/gf-api2/pkg/rest"
	"github.com/jinmukeji/go-pkg/age"
	corepb "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	generalpb "github.com/jinmukeji/proto/gen/micro/idl/ptypes/v2"
	"github.com/kataras/iris/v12"
)

const (
	// GenderMale 男性
	GenderMale string = "M"
	// GenderFemale 女性
	GenderFemale string = "F"
	// Int32GenderMale 男性
	Int32GenderMale = 0
	// Int32GenderFemale 女性
	Int32GenderFemale = 1
	// Int32GenderInvalid 性别有误
	Int32GenderInvalid = -1
)

// ModifyUserProfile 修改用户个人档案
type ModifyUserProfile struct {
	UserProfile UserProfile `json:"profile"`
}

// User 用户信息
type User struct {
	UserID             int32       `json:"user_id"`
	Username           string      `json:"username"`
	RegisterType       string      `json:"register_type"`
	Profile            UserProfile `json:"profile"`
	IsProfileCompleted bool        `json:"is_profile_completed"`
	IsRemovable        bool        `json:"is_removable"`
}

// UserProfile 用户信息
type UserProfile struct {
	Nickname        string    `json:"nickname"`
	Birthday        time.Time `json:"birthday"`
	Age             int       `json:"age"`
	Gender          int32     `json:"gender"`
	Height          int32     `json:"height"`
	Weight          int32     `json:"weight"`
	Phone           string    `json:"phone"`
	Email           string    `json:"email"`
	Remark          string    `json:"remark"`
	State           string    `json:"state"`
	City            string    `json:"city"`
	Street          string    `json:"street"`
	Country         string    `json:"country"`
	UserDefinedCode string    `json:"user_defined_code"`
}

// 修改个人档案
func (h *v2Handler) JinmuLModifyUserProfile(ctx iris.Context) {
	userID, err := ctx.Params().GetInt("user_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	req := new(corepb.JinmuLModifyUserProfileRequest)
	req.UserId = int32(userID)
	var modifyUserProfile ModifyUserProfile
	errReadJSON := ctx.ReadJSON(&modifyUserProfile)
	if errReadJSON != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", errReadJSON), false)
		return
	}
	birthday, _ := ptypes.TimestampProto(modifyUserProfile.UserProfile.Birthday)
	protoGender, errmapRestGenderToProto := mapRestGenderToProto(modifyUserProfile.UserProfile.Gender)
	if errmapRestGenderToProto != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errmapRestGenderToProto), false)
		return
	}
	req.UserProfile = &corepb.UserProfile{
		Nickname:        modifyUserProfile.UserProfile.Nickname,
		BirthdayTime:    birthday,
		Gender:          protoGender,
		Height:          modifyUserProfile.UserProfile.Height,
		Weight:          modifyUserProfile.UserProfile.Weight,
		Phone:           modifyUserProfile.UserProfile.Phone,
		Email:           modifyUserProfile.UserProfile.Email,
		Remark:          modifyUserProfile.UserProfile.Remark,
		UserDefinedCode: modifyUserProfile.UserProfile.UserDefinedCode,
		State:           modifyUserProfile.UserProfile.State,
		City:            modifyUserProfile.UserProfile.City,
		Street:          modifyUserProfile.UserProfile.Street,
		Country:         modifyUserProfile.UserProfile.Country,
	}
	errValidateUserProfile := validateUserProfile(*req.UserProfile)
	if errValidateUserProfile != nil {
		writeRpcInternalError(ctx, errValidateUserProfile, false)
		return
	}

	resp, errResp := h.rpcSvc.JinmuLModifyUserProfile(
		newRPCContext(ctx), req,
	)
	if errResp != nil {
		writeRpcInternalError(ctx, errResp, false)
		return
	}
	birthdayResp, _ := ptypes.Timestamp(resp.User.Profile.BirthdayTime)
	gender, errMapProtoGenderToRest := mapProtoGenderToRest(resp.User.Profile.Gender)
	if errMapProtoGenderToRest != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errMapProtoGenderToRest), false)
		return
	}
	rest.WriteOkJSON(ctx, User{
		UserID:             resp.User.UserId,
		Username:           resp.User.Username,
		RegisterType:       resp.User.RegisterType,
		IsProfileCompleted: resp.User.IsProfileCompleted,
		IsRemovable:        resp.User.IsRemovable,
		Profile: UserProfile{
			Nickname:        resp.User.Profile.Nickname,
			Birthday:        birthdayResp.UTC(),
			Age:             age.Age(birthdayResp),
			Height:          resp.User.Profile.Height,
			Gender:          gender,
			Weight:          resp.User.Profile.Weight,
			Phone:           resp.User.Profile.Phone,
			Email:           resp.User.Profile.Email,
			Remark:          resp.User.Profile.Remark,
			UserDefinedCode: resp.User.Profile.UserDefinedCode,
			State:           resp.User.Profile.State,
			City:            resp.User.Profile.City,
			Street:          resp.User.Profile.State,
			Country:         resp.User.Profile.Country,
		},
	})
}

// 查看用户个人档案
func (h *v2Handler) JinmuLGetUserProfile(ctx iris.Context) {
	userID, err := ctx.Params().GetInt("user_id")
	req := new(corepb.JinmuLGetUserProfileRequest)
	req.UserId = int32(userID)
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	resp, errResp := h.rpcSvc.JinmuLGetUserProfile(
		newRPCContext(ctx), req,
	)
	if errResp != nil {
		writeRpcInternalError(ctx, errResp, false)
		return
	}
	birthday, _ := ptypes.Timestamp(resp.User.Profile.BirthdayTime)
	gender, errMapProtoGenderToRest := mapProtoGenderToRest(resp.User.Profile.Gender)
	if errMapProtoGenderToRest != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errMapProtoGenderToRest), false)
		return
	}
	rest.WriteOkJSON(ctx, User{
		UserID:             resp.User.UserId,
		Username:           resp.User.Username,
		RegisterType:       resp.User.RegisterType,
		IsProfileCompleted: resp.User.IsProfileCompleted,
		IsRemovable:        resp.User.IsRemovable,
		Profile: UserProfile{
			Nickname:        resp.User.Profile.Nickname,
			Birthday:        birthday.UTC(),
			Age:             age.Age(birthday),
			Gender:          gender,
			Weight:          resp.User.Profile.Weight,
			Height:          resp.User.Profile.Height,
			Phone:           resp.User.Profile.Phone,
			Email:           resp.User.Profile.Email,
			Remark:          resp.User.Profile.Remark,
			UserDefinedCode: resp.User.Profile.UserDefinedCode,
			State:           resp.User.Profile.State,
			City:            resp.User.Profile.City,
			Street:          resp.User.Profile.State,
			Country:         resp.User.Profile.Country,
		},
	})
}

// mapRestGenderToProto 将传入的 int32 类型的 gender 转换为 proto 类型
func mapRestGenderToProto(gender int32) (generalpb.Gender, error) {
	switch gender {
	case Int32GenderMale:
		return generalpb.Gender_GENDER_MALE, nil
	case Int32GenderFemale:
		return generalpb.Gender_GENDER_FEMALE, nil
	}
	return generalpb.Gender_GENDER_INVALID, fmt.Errorf("invalid int32 gender %d", gender)
}

// mapProtoGenderToRest 将 proto 类型的 gender 转换为Rest的 int32 类型
func mapProtoGenderToRest(gender generalpb.Gender) (int32, error) {
	switch gender {
	case generalpb.Gender_GENDER_MALE:
		return Int32GenderMale, nil
	case generalpb.Gender_GENDER_FEMALE:
		return Int32GenderFemale, nil
	}
	return Int32GenderInvalid, fmt.Errorf("invalid int32 gender %d", gender)
}

// validateUserProfile 验证UserProfile 参数是否为空
func validateUserProfile(profile corepb.UserProfile) error {
	if profile.Weight == 0 {
		return errors.New("weight should not be 0")
	}
	if profile.Height == 0 {
		return errors.New("height should not be 0")
	}
	birthday, _ := ptypes.Timestamp(profile.BirthdayTime)
	if birthday.IsZero() {
		return errors.New("birthday should not be 0")
	}
	if profile.Gender != generalpb.Gender_GENDER_MALE && profile.Gender != generalpb.Gender_GENDER_FEMALE {
		return fmt.Errorf("invalid gender %d", profile.Gender)
	}
	return nil
}
