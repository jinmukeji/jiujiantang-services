package rest

import (
	"errors"
	"time"

	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	age "github.com/jinmukeji/go-pkg/age"
	jinmuidpb "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	corepb "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	generalpb "github.com/jinmukeji/proto/gen/micro/idl/ptypes/v2"
	"github.com/kataras/iris/v12"
)

// UserIDList user_id的集合
type UserIDList struct {
	UserIDList []int32 `json:"user_id_list"`
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

// UserSignUp 用户注册
type UserSignUp struct {
	UserProfile  UserProfile `json:"profile"`
	Password     string      `json:"password"`
	RegisterType string      `json:"register_type"`
	Username     string      `json:"username"`
	ClientID     string      `json:"client_id"`
}

// ModifyUserProfile 修改用户个人档案
type ModifyUserProfile struct {
	UserProfile UserProfile `json:"profile"`
}

// UserSignIn 用户登陆
type UserSignIn struct {
	SignInKey    string `json:"sign_in_key"`
	PasswordHash string `json:"password_hash"`
	RegisterType string `json:"register_type"`
}

// UserSignInResponse 用户登陆返回
type UserSignInResponse struct {
	AccessToken          string    `json:"access_token"`
	RemainDays           int32     `json:"remain_days"`
	ExpiredAt            time.Time `json:"expired_at"`
	UserID               int32     `json:"user_id"`
	AccessTokenExpiredAt time.Time `json:"access_token_expired_at"`
}

// 用户登录
func (h *v2Handler) UserSignIn(ctx iris.Context) {
	var userSignIn UserSignIn
	err := ctx.ReadJSON(&userSignIn)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}
	req := new(corepb.UserSignInRequest)
	req.SignInKey = userSignIn.SignInKey
	req.PasswordHash = userSignIn.PasswordHash
	req.RegisterType = userSignIn.RegisterType
	req.Ip = ctx.RemoteAddr()
	resp, err := h.rpcSvc.UserSignIn(
		newRPCContext(ctx), req,
	)
	if err != nil {
		writeRPCInternalError(ctx, err, false)
		return
	}
	expiredAt, _ := ptypes.Timestamp(resp.ExpireTime)
	accessTokenExpiredAt, _ := ptypes.Timestamp(resp.AccessTokenExpiredTime)
	rest.WriteOkJSON(ctx, UserSignInResponse{
		AccessToken:          resp.AccessToken,
		RemainDays:           resp.RemainDays,
		ExpiredAt:            expiredAt.UTC(),
		UserID:               resp.UserId,
		AccessTokenExpiredAt: accessTokenExpiredAt.UTC(),
	})
}

// 用户注册
func (h *v2Handler) UserSignUp(ctx iris.Context) {
	var userSignUp UserSignUp
	err := ctx.ReadJSON(&userSignUp)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}
	registerTypes := []string{"username", "email", "phone", "wechat", "legacy"}
	if !contains(registerTypes, userSignUp.RegisterType) {
		writeError(
			ctx,
			wrapError(ErrInvalidValue, "", fmt.Errorf("invalid registerType %s", userSignUp.RegisterType)),
			false,
		)
		return
	}
	if ctx.Values().GetString(ClientIDKey) != kangmeiClient {
		errValidateUserProfile := ValidateUserProfile(userSignUp.UserProfile)
		if errValidateUserProfile != nil {
			writeError(ctx, wrapError(ErrInvalidValue, "", errValidateUserProfile), false)
			return
		}
	}
	req := new(corepb.UserSignUpRequest)
	req.Password = userSignUp.Password
	req.RegisterType = userSignUp.RegisterType
	req.ClientId = userSignUp.ClientID
	birthday, _ := ptypes.TimestampProto(userSignUp.UserProfile.Birthday)
	var reqGender, reqHeight, reqWeight int32
	if userSignUp.UserProfile.Gender != nil {
		reqGender = *userSignUp.UserProfile.Gender
	}
	if userSignUp.UserProfile.Height != nil {
		reqHeight = *userSignUp.UserProfile.Height
	}
	if userSignUp.UserProfile.Weight != nil {
		reqWeight = *userSignUp.UserProfile.Weight
	}
	protoGender, errmapRestGenderToProto := mapRestGenderToProto(reqGender)
	if errmapRestGenderToProto != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errmapRestGenderToProto), false)
		return
	}
	req.UserProfile = &corepb.UserProfile{
		Nickname:        userSignUp.UserProfile.Nickname,
		BirthdayTime:    birthday,
		Gender:          protoGender,
		Height:          reqHeight,
		Weight:          reqWeight,
		Phone:           userSignUp.UserProfile.Phone,
		Email:           userSignUp.UserProfile.Email,
		Remark:          userSignUp.UserProfile.Remark,
		UserDefinedCode: userSignUp.UserProfile.UserDefinedCode,
		State:           userSignUp.UserProfile.State,
		City:            userSignUp.UserProfile.City,
		Street:          userSignUp.UserProfile.Street,
		Country:         userSignUp.UserProfile.Country,
	}
	resp, err := h.rpcSvc.UserSignUp(
		newRPCContext(ctx), req,
	)
	if err != nil {
		writeRPCInternalError(ctx, err, false)
		return
	}
	birthdayResp, _ := ptypes.Timestamp(resp.User.Profile.BirthdayTime)
	age := int32(age.Age(birthdayResp))
	int64Gender, errMapProtoGender := mapProtoGenderToRest(resp.User.Profile.Gender)
	if errMapProtoGender != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errMapProtoGender), true)
		return
	}
	gender := int32(int64Gender)
	height := resp.User.Profile.Height
	weight := resp.User.Profile.Weight
	rest.WriteOkJSON(ctx, User{
		UserID:             resp.User.UserId,
		Username:           resp.User.Username,
		RegisterType:       resp.User.RegisterType,
		IsProfileCompleted: resp.User.IsProfileCompleted,
		IsRemovable:        resp.User.IsRemovable,
		Profile: UserProfile{
			Nickname:        resp.User.Profile.Nickname,
			NicknameInitial: resp.User.Profile.NicknameInitial,
			Birthday:        birthdayResp.UTC(),
			Age:             &age,
			Gender:          &gender,
			Height:          &height,
			Weight:          &weight,
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
func (h *v2Handler) GetUserProfile(ctx iris.Context) {
	userID, err := ctx.Params().GetInt("user_id")
	req := new(corepb.GetUserProfileRequest)
	req.UserId = int32(userID)
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	resp, errResp := h.rpcSvc.GetUserProfile(
		newRPCContext(ctx), req,
	)
	if errResp != nil {
		writeRPCInternalError(ctx, errResp, false)
		return
	}
	birthday, _ := ptypes.Timestamp(resp.User.Profile.BirthdayTime)
	age := int32(age.Age(birthday))
	int64Gender, errMapProtoGender := mapProtoGenderToRest(resp.User.Profile.Gender)
	if errMapProtoGender != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errMapProtoGender), true)
		return
	}
	height := resp.User.Profile.Height
	weight := resp.User.Profile.Weight
	gender := int32(int64Gender)
	rest.WriteOkJSON(ctx, User{
		UserID:             resp.User.UserId,
		Username:           resp.User.Username,
		RegisterType:       resp.User.RegisterType,
		IsProfileCompleted: resp.User.IsProfileCompleted,
		IsRemovable:        resp.User.IsRemovable,
		Profile: UserProfile{
			Nickname:        resp.User.Profile.Nickname,
			Birthday:        birthday.UTC(),
			Age:             &age,
			Gender:          &gender,
			Height:          &height,
			Weight:          &weight,
			Remark:          resp.User.Profile.Remark,
			UserDefinedCode: resp.User.Profile.UserDefinedCode,
			State:           resp.User.Profile.State,
			City:            resp.User.Profile.City,
			Street:          resp.User.Profile.State,
			Country:         resp.User.Profile.Country,
			NicknameInitial: resp.User.Profile.NicknameInitial,
		},
	})
}

// 修改个人档案
func (h *v2Handler) ModifyUserProfile(ctx iris.Context) {
	userID, err := ctx.Params().GetInt("user_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	req := new(corepb.ModifyUserProfileRequest)
	req.UserId = int32(userID)
	var modifyUserProfile ModifyUserProfile
	errReadJSON := ctx.ReadJSON(&modifyUserProfile)
	if errReadJSON != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", errReadJSON), false)
		return
	}
	if ctx.Values().GetString(ClientIDKey) != kangmeiClient {
		errValidateUserProfile := ValidateUserProfile(modifyUserProfile.UserProfile)
		if errValidateUserProfile != nil {
			writeError(ctx, wrapError(ErrInvalidValue, "", errValidateUserProfile), false)
			return
		}
	}
	birthday, _ := ptypes.TimestampProto(modifyUserProfile.UserProfile.Birthday)
	var reqHeight, reqWeight int32
	if modifyUserProfile.UserProfile.Height != nil {
		reqHeight = *modifyUserProfile.UserProfile.Height
	}
	if modifyUserProfile.UserProfile.Weight != nil {
		reqWeight = *modifyUserProfile.UserProfile.Weight
	}
	req.UserProfile = &corepb.UserProfile{
		Nickname:        modifyUserProfile.UserProfile.Nickname,
		BirthdayTime:    birthday,
		Height:          reqHeight,
		Weight:          reqWeight,
		Phone:           modifyUserProfile.UserProfile.Phone,
		Email:           modifyUserProfile.UserProfile.Email,
		Remark:          modifyUserProfile.UserProfile.Remark,
		UserDefinedCode: modifyUserProfile.UserProfile.UserDefinedCode,
		State:           modifyUserProfile.UserProfile.State,
		City:            modifyUserProfile.UserProfile.City,
		Street:          modifyUserProfile.UserProfile.Street,
		Country:         modifyUserProfile.UserProfile.Country,
	}

	resp, errResp := h.rpcSvc.ModifyUserProfile(
		newRPCContext(ctx), req,
	)
	if errResp != nil {
		writeRPCInternalError(ctx, errResp, false)
		return
	}
	birthdayResp, _ := ptypes.Timestamp(resp.User.Profile.BirthdayTime)
	age := int32(age.Age(birthdayResp))
	int64Gender, errMapProtoGender := mapProtoGenderToRest(resp.User.Profile.Gender)
	if errMapProtoGender != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errMapProtoGender), true)
		return
	}
	height := resp.User.Profile.Height
	weight := resp.User.Profile.Weight
	gender := int32(int64Gender)
	rest.WriteOkJSON(ctx, User{
		UserID:             resp.User.UserId,
		Username:           resp.User.Username,
		RegisterType:       resp.User.RegisterType,
		IsProfileCompleted: resp.User.IsProfileCompleted,
		IsRemovable:        resp.User.IsRemovable,
		Profile: UserProfile{
			Nickname:        resp.User.Profile.Nickname,
			Birthday:        birthdayResp.UTC(),
			Age:             &age,
			Gender:          &gender,
			Height:          &height,
			Weight:          &weight,
			Remark:          resp.User.Profile.Remark,
			UserDefinedCode: resp.User.Profile.UserDefinedCode,
			State:           resp.User.Profile.State,
			City:            resp.User.Profile.City,
			Street:          resp.User.Profile.State,
			Country:         resp.User.Profile.Country,
			Phone:           resp.User.Profile.Phone,
			Email:           resp.User.Profile.Email,
			NicknameInitial: resp.User.Profile.NicknameInitial,
		},
	})
}

// 注销登录
func (h *v2Handler) SignOut(ctx iris.Context) {
	req := new(corepb.UserSignOutRequest)
	req.Ip = ctx.RemoteAddr()
	_, err := h.rpcSvc.UserSignOut(
		newRPCContext(ctx), req,
	)
	if err != nil {
		writeRPCInternalError(ctx, err, false)
		return
	}
	rest.WriteOkJSON(ctx, nil)
}

func mapRestGenderToProto(gender int32) (generalpb.Gender, error) {
	switch gender {
	case GenderMale:
		return generalpb.Gender_GENDER_MALE, nil
	case GenderFemale:
		return generalpb.Gender_GENDER_FEMALE, nil
	}
	return generalpb.Gender_GENDER_INVALID, fmt.Errorf("invalid int32 gender %d", gender)
}

// ValidateUserProfile 验证UserProfile 参数是否为空
func ValidateUserProfile(profile UserProfile) error {
	if profile.Nickname == "" {
		return errors.New("nickname should not be empty")
	}
	// if profile.Gender == nil {
	// 	return errors.New("gender should not be null")
	// }
	if profile.Weight == nil {
		return errors.New("weight should not be null")
	}
	if profile.Weight != nil && (*profile.Weight < 30 || *profile.Weight > 500) {
		return fmt.Errorf("invalid weight %d", *profile.Weight)
	}
	if profile.Height == nil {
		return errors.New("height should not be null")
	}
	if profile.Height != nil && (*profile.Height < 50 || *profile.Height > 250) {
		return fmt.Errorf("invalid height %d", *profile.Height)
	}
	if profile.Birthday.IsZero() {
		return errors.New("birthday should not be null")
	}
	year1900 := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	now := time.Now()
	if !profile.Birthday.IsZero() && (profile.Birthday.Before(year1900) || profile.Birthday.After(now)) {
		return fmt.Errorf("invalid birthday %v", profile.Birthday)
	}
	if profile.Gender != nil && (*profile.Gender > 1 || *profile.Gender < 0) {
		return fmt.Errorf("invalid gender %d", *profile.Gender)
	}
	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// OwnerUserSignUpBody 添加用户的请求体
type OwnerUserSignUpBody struct {
	Nickname     string    `json:"nickname"`
	Gender       *int32    `json:"gender"`
	Birthday     time.Time `json:"birthday"`
	Height       int32     `json:"height"`
	Weight       int32     `json:"weight"`
	RegisterType string    `json:"register_type"`
	ClientID     string    `json:"client_id"`
}

// OwnerUserSignUpReply 添加用户的响应
type OwnerUserSignUpReply struct {
	UserID      int64     `json:"user_id"`
	Nickname    string    `json:"nickname"`
	IsRemovable bool      `json:"is_removable"`
	Gender      int64     `json:"gender"`
	Birthday    time.Time `json:"birthday"`
	Height      int64     `json:"height"`
	Weight      int64     `json:"weight"`
}

// 添加用户
func (h *v2Handler) OwnerUserSignUp(ctx iris.Context) {
	userID, _ := ctx.Params().GetInt("owner_id")
	reqCanOwnerAddOrganizationUserRequest := new(corepb.CanOwnerAddOrganizationUserRequest)
	canOwnerAddOrganizationUser, errCanOwnerAddOrganizationUser := h.rpcSvc.CanOwnerAddOrganizationUser(
		newRPCContext(ctx), reqCanOwnerAddOrganizationUserRequest,
	)
	if errCanOwnerAddOrganizationUser != nil {
		writeRPCInternalError(ctx, errCanOwnerAddOrganizationUser, true)
		return
	}
	if !canOwnerAddOrganizationUser.Able {
		writeError(ctx, wrapError(ErrNullClientID, "", fmt.Errorf("user count of user %d will exceeded the maximum number", userID)), true)
		return
	}
	var ownerUserSignUpBody OwnerUserSignUpBody
	err := ctx.ReadJSON(&ownerUserSignUpBody)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}
	registerTypes := []string{"username", "email", "phone", "wechat", "legacy"}
	if !contains(registerTypes, ownerUserSignUpBody.RegisterType) {
		writeError(
			ctx,
			wrapError(ErrInvalidValue, "", fmt.Errorf("invalid registerType %s", ownerUserSignUpBody.RegisterType)),
			false,
		)
		return
	}
	if ctx.Values().GetString(ClientIDKey) != kangmeiClient {
		userProfile := UserProfile{
			Nickname: ownerUserSignUpBody.Nickname,
			Birthday: ownerUserSignUpBody.Birthday,
			Gender:   ownerUserSignUpBody.Gender,
			Height:   &ownerUserSignUpBody.Height,
			Weight:   &ownerUserSignUpBody.Weight,
		}
		errValidateUserProfile := ValidateUserProfile(userProfile)
		if errValidateUserProfile != nil {
			writeError(ctx, wrapError(ErrInvalidValue, "", errValidateUserProfile), false)
			return
		}
	}
	req := new(corepb.UserSignUpRequest)
	req.RegisterType = ownerUserSignUpBody.RegisterType
	req.ClientId = ownerUserSignUpBody.ClientID
	birthday, _ := ptypes.TimestampProto(ownerUserSignUpBody.Birthday)
	var reqGender int32
	if ownerUserSignUpBody.Gender != nil {
		reqGender = *ownerUserSignUpBody.Gender
	}
	protoGender, errmapRestGenderToProto := mapRestGenderToProto(reqGender)
	if errmapRestGenderToProto != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errmapRestGenderToProto), false)
		return
	}
	req.UserProfile = &corepb.UserProfile{
		Nickname:     ownerUserSignUpBody.Nickname,
		BirthdayTime: birthday,
		Gender:       protoGender,
		Height:       ownerUserSignUpBody.Height,
		Weight:       ownerUserSignUpBody.Weight,
	}
	resp, err := h.rpcSvc.UserSignUp(
		newRPCContext(ctx), req,
	)
	if err != nil {
		writeRPCInternalError(ctx, err, false)
		return
	}
	reqGetOrganizationIDByUserID := new(corepb.GetOrganizationIDByUserIDRequest)
	reqGetOrganizationIDByUserID.UserId = int32(userID)
	respGetOrganizationIDByUserID, errGetOrganizationIDByUserID := h.rpcSvc.GetOrganizationIDByUserID(
		newRPCContext(ctx), reqGetOrganizationIDByUserID,
	)
	if errGetOrganizationIDByUserID != nil {
		writeRPCInternalError(ctx, errGetOrganizationIDByUserID, false)
		return
	}
	if !respGetOrganizationIDByUserID.IsOwner {
		writeError(
			ctx,
			wrapError(ErrInvalidUser, "", errors.New("invalid user")),
			false,
		)

		return
	}
	organizationID := respGetOrganizationIDByUserID.OrganizationId
	reqOwnerAddOrganizationUsers := new(corepb.OwnerAddOrganizationUsersRequest)
	reqOwnerAddOrganizationUsers.OrganizationId = int32(organizationID)
	reqOwnerAddOrganizationUsers.UserIdList = []int32{resp.User.UserId}
	_, errOwnerAddOrganizationUsers := h.rpcSvc.OwnerAddOrganizationUsers(
		newRPCContext(ctx), reqOwnerAddOrganizationUsers,
	)
	if errOwnerAddOrganizationUsers != nil {
		writeRPCInternalError(ctx, errOwnerAddOrganizationUsers, false)
		return
	}

	birthdayResp, _ := ptypes.Timestamp(resp.User.Profile.BirthdayTime)
	gender, errMapProtoGender := mapProtoGenderToRest(resp.User.Profile.Gender)
	if errMapProtoGender != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errMapProtoGender), true)
		return
	}
	height := resp.User.Profile.Height
	weight := resp.User.Profile.Weight
	rest.WriteOkJSON(ctx, OwnerUserSignUpReply{
		UserID:      int64(resp.User.UserId),
		IsRemovable: resp.User.IsRemovable,
		Nickname:    resp.User.Profile.Nickname,
		Gender:      int64(gender),
		Birthday:    birthdayResp.UTC(),
		Height:      int64(height),
		Weight:      int64(weight),
	})
}

// BindOldUserResponse 绑定老用户的响应
type BindOldUserResponse struct {
	UserID      int32  `json:"user_id"`
	AccessToken string `json:"access_token"`
}

// BindOldUserBody 绑定老用户的请求
type BindOldUserBody struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Seed         string `json:"seed"`
}

// 绑定老用户
func (h *v2Handler) BindOldUser(ctx iris.Context) {
	userID, _ := ctx.Params().GetInt("user_id")
	var bindOldUserBody BindOldUserBody
	errReadJSON := ctx.ReadJSON(&bindOldUserBody)
	if errReadJSON != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", errReadJSON), false)
		return
	}
	req := new(jinmuidpb.BindOldUserRequest)
	req.Username = bindOldUserBody.Username
	req.UserId = int32(userID)
	req.PasswordHash = bindOldUserBody.PasswordHash
	req.Seed = bindOldUserBody.Seed

	resp, errBindOldUser := h.rpcJinmuidSvc.BindOldUser(
		newRPCContext(ctx), req,
	)
	if errBindOldUser != nil {
		writeRPCInternalError(ctx, errBindOldUser, false)
		return
	}

	rest.WriteOkJSON(ctx, BindOldUserResponse{
		UserID:      resp.UserId,
		AccessToken: resp.AccessToken,
	})

}
