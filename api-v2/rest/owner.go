package rest

import (
	"math"
	"time"

	"fmt"

	"github.com/golang/protobuf/ptypes"
	age "github.com/jinmukeji/go-pkg/age"
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	corepb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	"github.com/kataras/iris/v12"
)

const (
	// CustomizedVersion 定制化
	CustomizedVersion = "定制化"
	// TrialVersion 试用版
	TrialVersion = "试用版"
	// GoldenVersion 黄喜马把脉
	GoldenVersion = "黄喜马把脉"
	// PlatinumVersion 白喜马把脉
	PlatinumVersion = "白喜马把脉"
	// DiamondVersion 钻石姆
	DiamondVersion = "钻石姆"
	// GiftVersion 礼品版
	GiftVersion = "礼品版"
	// AlmostExpire 账号快过期
	AlmostExpire = "账号快过期"
	// AlreadyExpired 账号已经过期
	AlreadyExpired = "账号已经过期"
	// kangmeiClient 康美的ClientID
	kangmeiClient = "kangmei-10001"
	// dengyunClient 登云的client
	dengyunClient = "dengyun-10001"
	// hongjingtang 弘经堂客户端
	hongjingtangClient = "hongjingtang-10001"
	// moai 摩爱客户端
	moaiClient = "moai-10001"
)

const (
	// ExpirationStatusExpired 已过期
	ExpirationStatusExpired = iota
	// ExpirationStatusAlmostExpired 即将过期
	ExpirationStatusAlmostExpired
	// ExpirationStatusNotExpired 未过期
	ExpirationStatusNotExpired
)

// Subscription 订阅信息
type Subscription struct {
	CreatedAt        time.Time `json:"created_at"`        // 创建时间
	TotalUserCount   int32     `json:"total_user_count"`  // 添加人数
	SubscriptionType string    `json:"subscription_type"` // 会员类型
	ExpiredAt        time.Time `json:"expired_at"`        // 到期时间
	MaxUserLimits    int32     `json:"max_user_limits"`   // 最大人数
	Active           bool      `json:"active"`            // 是否激活
}

// OwnerSignUp 用户注册
type OwnerSignUp struct {
	UserProfile  UserProfile `json:"profile"`       // 用户档案
	Password     string      `json:"password"`      // 密码
	RegisterType string      `json:"register_type"` // 注册类型
	Username     string      `json:"username"`      // 用户名
	ClientID     string      `json:"client_id"`     // ClientID
}

// Organization 组织
type Organization struct {
	OrganizationID int32               `json:"organization_id"` // 组织ID
	Subscription   Subscription        `json:"subscription"`    // 订阅信息
	Profile        OrganizationProfile `json:"profile"`         // 组织档案
}

// OrganizationProfile 组织档案
type OrganizationProfile struct {
	Name    string `json:"name"`    // 名字
	State   string `json:"state"`   // 区域
	City    string `json:"city"`    // 城市
	Street  string `json:"street"`  // 街道
	Phone   string `json:"phone"`   // 手机
	Contact string `json:"contact"` // 联系人
	Type    string `json:"type"`    // 类型
	Email   string `json:"email"`   // 邮箱
	Country string `json:"country"` // 国家
}

// UserProfile 用户信息
type UserProfile struct {
	Nickname        string    `json:"nickname"`
	Birthday        time.Time `json:"birthday"`
	Age             *int32    `json:"age"`
	Gender          *int32    `json:"gender"`
	Height          *int32    `json:"height"`
	Weight          *int32    `json:"weight"`
	Phone           string    `json:"phone"`
	Email           string    `json:"email"`
	Remark          string    `json:"remark"`
	State           string    `json:"state"`
	City            string    `json:"city"`
	Street          string    `json:"street"`
	Country         string    `json:"country"`
	UserDefinedCode string    `json:"user_defined_code"`
	NicknameInitial string    `json:"nickname_initial"`
}

// OwnerGetOrganizations 查看拥有的组织
func (h *v2Handler) OwnerGetOrganizations(ctx iris.Context) {
	req := new(corepb.OwnerGetOrganizationsRequest)
	resp, err := h.rpcSvc.OwnerGetOrganizations(
		newRPCContext(ctx), req,
	)
	if err != nil {
		writeRPCInternalError(ctx, err, true)
		return
	}
	organizations := make([]Organization, len(resp.Organizations))
	for idx, o := range resp.Organizations {
		organizations[idx] = Organization{
			OrganizationID: o.OrganizationId,
			Profile: OrganizationProfile{
				Name:    o.Profile.Name,
				State:   o.Profile.Address.State,
				City:    o.Profile.Address.City,
				Street:  o.Profile.Address.Street,
				Phone:   o.Profile.Phone,
				Contact: o.Profile.Contact,
				Type:    o.Profile.Type,
				Email:   o.Profile.Email,
				Country: o.Profile.Address.Country,
			},
		}
		if o.Subscription != nil {
			expiredAt, _ := ptypes.Timestamp(o.Subscription.ExpiredTime)
			createdAt, _ := ptypes.Timestamp(o.Subscription.CreatedTime)
			stringSubscriptionType, errmapProtoSubscriptionTypeToRest := mapProtoSubscriptionTypeToRest(o.Subscription.SubscriptionType)
			if errmapProtoSubscriptionTypeToRest != nil {
				writeError(ctx, wrapError(ErrInvalidValue, "", errmapProtoSubscriptionTypeToRest), false)
				return
			}
			organizations[idx].Subscription = Subscription{
				CreatedAt:        createdAt.UTC(),
				TotalUserCount:   o.Subscription.TotalUserCount,
				MaxUserLimits:    o.Subscription.MaxUserLimits,
				SubscriptionType: stringSubscriptionType,
				ExpiredAt:        expiredAt.UTC(),
				Active:           o.Subscription.Active,
			}
		}
	}

	rest.WriteOkJSON(ctx, organizations)
}

// OwnerAddOrganizationUsers Owner 在组织中增加一个或多个 User
func (h *v2Handler) OwnerAddOrganizationUsers(ctx iris.Context) {
	organizationID, err := ctx.Params().GetInt("organization_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), true)
		return
	}
	reqCanOwnerAddOrganizationUserRequest := new(corepb.CanOwnerAddOrganizationUserRequest)
	canOwnerAddOrganizationUser, errCanOwnerAddOrganizationUser := h.rpcSvc.CanOwnerAddOrganizationUser(
		newRPCContext(ctx), reqCanOwnerAddOrganizationUserRequest,
	)
	if errCanOwnerAddOrganizationUser != nil {
		writeRPCInternalError(ctx, errCanOwnerAddOrganizationUser, true)
		return
	}
	if !canOwnerAddOrganizationUser.Able {
		writeError(ctx, wrapError(ErrNullClientID, "", fmt.Errorf("user count of organization %d will exceeded the maximum number", organizationID)), true)
		return
	}
	var userIDList UserIDList
	err = ctx.ReadJSON(&userIDList)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), true)
		return
	}
	req := new(corepb.OwnerAddOrganizationUsersRequest)
	req.OrganizationId = int32(organizationID)
	req.UserIdList = userIDList.UserIDList
	resp, err := h.rpcSvc.OwnerAddOrganizationUsers(
		newRPCContext(ctx), req,
	)
	if err != nil {
		writeRPCInternalError(ctx, err, true)
		return
	}
	users := make([]User, len(resp.User))
	for idx, u := range resp.User {
		birthday, err := ptypes.Timestamp(u.Profile.BirthdayTime)
		if err != nil {
			continue
		}
		age := int32(age.Age(birthday))
		int64Gender, errMapProtoGender := mapProtoGenderToRest(u.Profile.Gender)
		if errMapProtoGender != nil {
			writeError(ctx, wrapError(ErrInvalidValue, "", errMapProtoGender), true)
			return
		}
		gender := int32(int64Gender)
		height := u.Profile.Height
		weight := u.Profile.Weight
		users[idx] = User{
			UserID:             u.UserId,
			Username:           u.Username,
			RegisterType:       u.RegisterType,
			IsProfileCompleted: u.IsProfileCompleted,
			IsRemovable:        u.IsRemovable,
			Profile: UserProfile{
				Nickname:        u.Profile.Nickname,
				NicknameInitial: u.Profile.NicknameInitial,
				Birthday:        birthday.UTC(),
				Age:             &age,
				Gender:          &gender,
				Height:          &height,
				Weight:          &weight,
				Phone:           u.Profile.Phone,
				Email:           u.Profile.Email,
				Remark:          u.Profile.Remark,
				State:           u.Profile.State,
				City:            u.Profile.City,
				Street:          u.Profile.Street,
				Country:         u.Profile.Country,
				UserDefinedCode: u.Profile.UserDefinedCode,
			},
		}
	}
	rest.WriteOkJSON(ctx, users)
}

// OwnerGetOrganizationUsers Owner 获取组织中的 User
func (h *v2Handler) OwnerGetOrganizationUsers(ctx iris.Context) {
	organizationID, err := ctx.Params().GetInt("organization_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), true)
		return
	}
	req := new(corepb.OwnerGetOrganizationUsersRequest)
	req.OrganizationId = int32(organizationID)
	offset, _ := ctx.URLParamInt("offset")
	size, _ := ctx.URLParamInt("size")
	keyword := ctx.URLParam("keyword")
	req.Keyword = keyword
	req.Size = int32(size)
	req.Offset = int32(offset)
	resp, err := h.rpcSvc.OwnerGetOrganizationUsers(
		newRPCContext(ctx), req,
	)
	if err != nil {
		writeRPCInternalError(ctx, err, true)
		return
	}
	users := make([]User, 0)
	for _, u := range resp.UserList {
		birthday, err := ptypes.Timestamp(u.Profile.BirthdayTime)
		if err != nil {
			continue
		}
		age := int32(age.Age(birthday))
		int64Gender, errMapProtoGender := mapProtoGenderToRest(u.Profile.Gender)
		if errMapProtoGender != nil {
			writeError(ctx, wrapError(ErrInvalidValue, "", errMapProtoGender), true)
			return
		}
		gender := int32(int64Gender)
		height := u.Profile.Height
		weight := u.Profile.Weight
		users = append(users, User{
			UserID:             u.UserId,
			Username:           u.Username,
			RegisterType:       u.RegisterType,
			IsProfileCompleted: u.IsProfileCompleted,
			IsRemovable:        u.IsRemovable,
			Profile: UserProfile{
				Nickname:        u.Profile.Nickname,
				NicknameInitial: u.Profile.NicknameInitial,
				Birthday:        birthday.UTC(),
				Age:             &age,
				Gender:          &gender,
				Height:          &height,
				Weight:          &weight,
				Remark:          u.Profile.Remark,
				State:           u.Profile.State,
				City:            u.Profile.City,
				Street:          u.Profile.Street,
				Country:         u.Profile.Country,
				UserDefinedCode: u.Profile.UserDefinedCode,
				Phone:           u.Profile.Phone,
				Email:           u.Profile.Email,
			},
		})
	}
	rest.WriteOkJSON(ctx, users)
}

// OwnerGetOrganizationUsersByOwnerID OwnerID获取组织下用户
func (h *v2Handler) OwnerGetOrganizationUsersByOwnerID(ctx iris.Context) {
	reqGetOrganizations := new(corepb.OwnerGetOrganizationsRequest)
	respOwnerGetOrganizations, err := h.rpcSvc.OwnerGetOrganizations(
		newRPCContext(ctx), reqGetOrganizations,
	)
	if err != nil {
		writeRPCInternalError(ctx, err, true)
		return
	}
	req := new(corepb.OwnerGetOrganizationUsersRequest)
	req.OrganizationId = int32(respOwnerGetOrganizations.Organizations[0].OrganizationId)
	offset, _ := ctx.URLParamInt("offset")
	size, _ := ctx.URLParamInt("size")
	keyword := ctx.URLParam("keyword")
	req.Keyword = keyword
	req.Size = int32(size)
	req.Offset = int32(offset)
	resp, err := h.rpcSvc.OwnerGetOrganizationUsers(
		newRPCContext(ctx), req,
	)
	if err != nil {
		writeRPCInternalError(ctx, err, true)
		return
	}
	users := make([]User, 0)
	for _, u := range resp.UserList {
		birthday, err := ptypes.Timestamp(u.Profile.BirthdayTime)
		if err != nil {
			continue
		}
		age := int32(age.Age(birthday))
		int64Gender, errMapProtoGender := mapProtoGenderToRest(u.Profile.Gender)
		if errMapProtoGender != nil {
			writeError(ctx, wrapError(ErrInvalidValue, "", errMapProtoGender), true)
			return
		}
		gender := int32(int64Gender)
		height := u.Profile.Height
		weight := u.Profile.Weight
		users = append(users, User{
			UserID:             u.UserId,
			Username:           u.Username,
			RegisterType:       u.RegisterType,
			IsProfileCompleted: u.IsProfileCompleted,
			IsRemovable:        u.IsRemovable,
			Profile: UserProfile{
				Nickname:        u.Profile.Nickname,
				NicknameInitial: u.Profile.NicknameInitial,
				Birthday:        birthday.UTC(),
				Age:             &age,
				Gender:          &gender,
				Height:          &height,
				Weight:          &weight,
				Remark:          u.Profile.Remark,
				State:           u.Profile.State,
				City:            u.Profile.City,
				Street:          u.Profile.Street,
				Country:         u.Profile.Country,
				UserDefinedCode: u.Profile.UserDefinedCode,
				Phone:           u.Profile.UserDefinedCode,
				Email:           u.Profile.Email,
			},
		})
	}
	rest.WriteOkJSON(ctx, users)
}

// OwnerSignUp Owner 注册 User
func (h *v2Handler) OwnerSignUp(ctx iris.Context) {
	var ownerSignUp OwnerSignUp
	err := ctx.ReadJSON(&ownerSignUp)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}
	registerTypes := []string{"username", "email", "phone", "wechat", "legacy"}
	if !contains(registerTypes, ownerSignUp.RegisterType) {
		writeError(
			ctx,
			wrapError(ErrInvalidValue, "", fmt.Errorf("invalid registerType %s", ownerSignUp.RegisterType)),
			false,
		)
		return
	}
	if ctx.Values().GetString(ClientIDKey) != kangmeiClient && ctx.Values().GetString(ClientIDKey) != dengyunClient && ctx.Values().GetString(ClientIDKey) != hongjingtangClient && ctx.Values().GetString(ClientIDKey) != moaiClient {
		errValidateUserProfile := ValidateUserProfile(ownerSignUp.UserProfile)
		if errValidateUserProfile != nil {
			writeError(ctx, wrapError(ErrInvalidValue, "", errValidateUserProfile), false)
			return
		}
	}
	req := new(corepb.UserSignUpRequest)
	req.Password = ownerSignUp.Password
	req.RegisterType = ownerSignUp.RegisterType
	req.ClientId = ownerSignUp.ClientID
	req.Username = ownerSignUp.Username
	birthday, _ := ptypes.TimestampProto(ownerSignUp.UserProfile.Birthday)
	var reqGender, reqHeight, reqWeight int32
	if ownerSignUp.UserProfile.Gender != nil {
		reqGender = *ownerSignUp.UserProfile.Gender
	}
	if ownerSignUp.UserProfile.Height != nil {
		reqHeight = *ownerSignUp.UserProfile.Height
	}
	if ownerSignUp.UserProfile.Weight != nil {
		reqWeight = *ownerSignUp.UserProfile.Weight
	}
	protoGender, errmapRestGenderToProto := mapRestGenderToProto(reqGender)
	if errmapRestGenderToProto != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errmapRestGenderToProto), false)
		return
	}
	req.UserProfile = &corepb.UserProfile{
		Nickname:        ownerSignUp.UserProfile.Nickname,
		BirthdayTime:    birthday,
		Gender:          protoGender,
		Height:          reqHeight,
		Weight:          reqWeight,
		Phone:           ownerSignUp.UserProfile.Phone,
		Email:           ownerSignUp.UserProfile.Email,
		Remark:          ownerSignUp.UserProfile.Remark,
		State:           ownerSignUp.UserProfile.State,
		City:            ownerSignUp.UserProfile.City,
		Street:          ownerSignUp.UserProfile.Street,
		Country:         ownerSignUp.UserProfile.Country,
		UserDefinedCode: ownerSignUp.UserProfile.UserDefinedCode,
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

// OwnerDeleteOrganizationUsers
func (h *v2Handler) OwnerDeleteOrganizationUsers(ctx iris.Context) {
	var userIDList UserIDList
	err := ctx.ReadJSON(&userIDList)
	if err != nil {
		writeError(ctx, wrapError(ErrRPCInternal, "", err), false)
		return
	}
	organizationID, _ := ctx.Params().GetInt("organization_id")
	req := new(corepb.OwnerDeleteOrganizationUsersRequest)
	req.OrganizationId = int32(organizationID)
	req.UserIdList = userIDList.UserIDList
	_, errOwnerDeleteOrganizationUsers := h.rpcSvc.OwnerDeleteOrganizationUsers(
		newRPCContext(ctx), req,
	)
	if errOwnerDeleteOrganizationUsers != nil {
		writeRPCInternalError(ctx, errOwnerDeleteOrganizationUsers, false)
		return
	}

	rest.WriteOkJSON(ctx, nil)
}

// OrganizationSubscription 组织订阅
type OrganizationSubscription struct {
	CreatedAt        time.Time `json:"created_at"`
	TotalUserCount   int32     `json:"total_user_count"`
	MaxUserLimits    int32     `json:"max_user_limits"`
	SubscriptionType string    `json:"subscription_type"`
	ExpiredAt        time.Time `json:"expired_at"`
	Message          string    `json:"message"`
	ExpirationStatus int       `json:"expiration_status"`
}

// GetOrganizationSubscription 查看订阅
func (h *v2Handler) GetOrganizationSubscription(ctx iris.Context) {
	organizationID, _ := ctx.Params().GetInt("organization_id")
	req := new(corepb.OwnerGetOrganizationSubscriptionRequest)
	req.OrganizationId = int32(organizationID)
	resp, err := h.rpcSvc.OwnerGetOrganizationSubscription(
		newRPCContext(ctx), req,
	)
	if err != nil {
		writeRPCInternalError(ctx, err, false)
		return
	}
	expiredAt, _ := ptypes.Timestamp(resp.Subscription.ExpiredTime)
	createdAt, _ := ptypes.Timestamp(resp.Subscription.CreatedTime)
	day := int32(math.Ceil(time.Until(expiredAt.UTC()).Hours() / 24))
	if day <= 0 {
		stringSubscriptionType, errmapProtoSubscriptionTypeToRest := mapProtoSubscriptionTypeToRest(resp.Subscription.SubscriptionType)
		if errmapProtoSubscriptionTypeToRest != nil {
			writeError(ctx, wrapError(ErrInvalidValue, "", errmapProtoSubscriptionTypeToRest), false)
			return
		}
		rest.WriteOkJSON(ctx, OrganizationSubscription{
			CreatedAt:        createdAt.UTC(),
			TotalUserCount:   resp.Subscription.TotalUserCount,
			MaxUserLimits:    resp.Subscription.MaxUserLimits,
			SubscriptionType: stringSubscriptionType,
			ExpiredAt:        expiredAt.UTC(),
			Message:          AlreadyExpired,
			ExpirationStatus: ExpirationStatusExpired,
		})
		return
	}

	// 一个月以内提示快用过期
	if day > 0 && day < 30 {
		stringSubscriptionType, errmapProtoSubscriptionTypeToRest := mapProtoSubscriptionTypeToRest(resp.Subscription.SubscriptionType)
		if errmapProtoSubscriptionTypeToRest != nil {
			writeError(ctx, wrapError(ErrInvalidValue, "", errmapProtoSubscriptionTypeToRest), false)
			return
		}
		rest.WriteOkJSON(ctx, OrganizationSubscription{
			CreatedAt:        createdAt.UTC(),
			TotalUserCount:   resp.Subscription.TotalUserCount,
			MaxUserLimits:    resp.Subscription.MaxUserLimits,
			SubscriptionType: stringSubscriptionType,
			ExpiredAt:        expiredAt.UTC(),
			Message:          AlmostExpire,
			ExpirationStatus: ExpirationStatusAlmostExpired,
		})
		return
	}
	stringSubscriptionType, errmapProtoSubscriptionTypeToRest := mapProtoSubscriptionTypeToRest(resp.Subscription.SubscriptionType)
	if errmapProtoSubscriptionTypeToRest != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errmapProtoSubscriptionTypeToRest), false)
		return
	}
	rest.WriteOkJSON(ctx, OrganizationSubscription{
		CreatedAt:        createdAt.UTC(),
		TotalUserCount:   resp.Subscription.TotalUserCount,
		MaxUserLimits:    resp.Subscription.MaxUserLimits,
		SubscriptionType: stringSubscriptionType,
		ExpiredAt:        expiredAt.UTC(),
		Message:          "",
		ExpirationStatus: ExpirationStatusNotExpired,
	})
}

// OwnerDeleteUsers Owner删除用户
func (h *v2Handler) OwnerDeleteUsers(ctx iris.Context) {
	var userIDList UserIDList
	err := ctx.ReadJSON(&userIDList)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}
	ownerID, err := ctx.Params().GetInt("owner_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	reqOwnerDeleteUsers := new(corepb.OwnerDeleteUsersRequest)
	reqOwnerDeleteUsers.OwnerId = int32(ownerID)
	reqOwnerDeleteUsers.UserIdList = userIDList.UserIDList
	_, errOwnerDeleteUsers := h.rpcSvc.OwnerDeleteUsers(
		newRPCContext(ctx), reqOwnerDeleteUsers,
	)
	if errOwnerDeleteUsers != nil {
		writeRPCInternalError(ctx, errOwnerDeleteUsers, false)
		return
	}
	rest.WriteOkJSON(ctx, nil)

}
