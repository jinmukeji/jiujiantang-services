package handler

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/jinmukeji/jiujiantang-services/jinmuid/mysqldb"
	subscriptionpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/subscription/v1"

	"github.com/golang/protobuf/ptypes"
	jinmuidpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	bizcorepb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	generalpb "github.com/jinmukeji/proto/v3/gen/micro/idl/ptypes/v2"
	"github.com/mozillazg/go-pinyin"
)

const (
	// 注册方式为微信
	WechatRegisterType = "wechat"
)

// GetUserProfile 得到用户档案
func (j *JinmuIDService) GetUserProfile(ctx context.Context, req *jinmuidpb.GetUserProfileRequest, resp *jinmuidpb.GetUserProfileResponse) error {
	// 如果需要token验证
	if !req.IsSkipVerifyToken {
		token, ok := TokenFromContext(ctx)
		if !ok {
			return NewError(ErrInvalidUser, errors.New("failed to get token from context"))
		}
		userID, err := j.datastore.FindUserIDByToken(ctx, token)
		if err != nil {
			return NewError(ErrUserUnauthorized, errors.New("failed to get userId by token"))
		}
		if userID != req.UserId {
			return NewError(ErrInvalidUser, fmt.Errorf("user %d from request and user %d from token are inconsistent", req.UserId, userID))
		}
	}
	userProfile, err := j.datastore.FindUserProfile(ctx, req.UserId)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find user profile of userId %d: %s", req.UserId, err.Error()))
	}
	user, errFindUserByUserID := j.datastore.FindUserByUserID(ctx, req.UserId)
	if errFindUserByUserID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find user by userId %d: %s", req.UserId, errFindUserByUserID.Error()))
	}
	birthday, _ := ptypes.TimestampProto(userProfile.Birthday)
	protoGender, errmapDBGenderToProto := mapDBGenderToProto(userProfile.Gender)
	if errmapDBGenderToProto != nil {
		return NewError(ErrInvalidGender, errmapDBGenderToProto)
	}
	resp.Profile = &jinmuidpb.UserProfile{
		Nickname:        userProfile.Nickname,
		NicknameInitial: userProfile.NicknameInitial,
		BirthdayTime:    birthday,
		Gender:          protoGender,
		Weight:          userProfile.Weight,
		Height:          userProfile.Height,
	}
	resp.HasSetEmail = user.HasSetEmail
	resp.HasSetPassword = user.HasSetPassword
	resp.HasSetPhone = user.HasSetPhone
	resp.HasSetRegion = user.HasSetRegion
	resp.HasSetUsername = user.HasSetUsername
	resp.HasSetSecureQuestions = user.HasSetSecureQuestions
	resp.SigninUsername = user.SigninUsername
	resp.SecureEmail = safeDealEmail(user.SecureEmail)
	resp.SigninPhone = safeDealPhone(user.SigninPhone)
	resp.Region = string(user.Region)
	if user.HasSetLanguage {
		protoLanguage, errMapDBLanguageToProto := mapDBLanguageToProto(string(user.Language))
		if errMapDBLanguageToProto != nil {
			return NewError(ErrInvalidUser, errMapDBLanguageToProto)
		}

		resp.Language = protoLanguage
	} else {
		resp.Language = generalpb.Language_LANGUAGE_UNSET
	}
	resp.HasSetUserProfile = user.HasSetUserProfile
	resp.CustomizedCode = user.CustomizedCode
	resp.UserDefinedCode = user.UserDefinedCode
	resp.IsProfileCompleted = user.IsProfileCompleted
	return nil
}

// ModifyUserProfile 修改用户档案
func (j *JinmuIDService) ModifyUserProfile(ctx context.Context, req *jinmuidpb.ModifyUserProfileRequest, resp *jinmuidpb.ModifyUserProfileResponse) error {
	// TODO:FIXME: 57000-65000
	// 如果需要token验证
	if !req.IsSkipVerifyToken {
		token, ok := TokenFromContext(ctx)
		if !ok {
			return NewError(ErrInvalidUser, errors.New("failed to get token from context"))
		}
		userID, err := j.datastore.FindUserIDByToken(ctx, token)
		if err != nil {
			return NewError(ErrUserUnauthorized, fmt.Errorf("failed to get userID by token: %s", err.Error()))
		}
		if userID != req.UserId {
			return NewError(ErrInvalidUser, fmt.Errorf("user %d from request and user %d from token are inconsistent", req.UserId, userID))
		}
	}
	dbGender, errmapProtoGenderToDB := mapProtoGenderToDB(req.UserProfile.Gender)
	if errmapProtoGenderToDB != nil {
		return NewError(ErrInvalidGender, errmapProtoGenderToDB)
	}
	userProfile, errFindUserProfile := j.datastore.FindUserProfile(ctx, req.UserId)
	if errFindUserProfile != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find user profile of user %d: %s", req.UserId, errFindUserProfile.Error()))
	}
	birthday, _ := ptypes.Timestamp(req.UserProfile.BirthdayTime)
	profile := &mysqldb.UserProfile{
		UserID:          req.UserId,
		Nickname:        req.UserProfile.Nickname,
		Gender:          dbGender,
		Weight:          req.UserProfile.Weight,
		Height:          req.UserProfile.Height,
		Birthday:        birthday,
		NicknameInitial: getNicknameInitial(req.UserProfile.Nickname),
	}
	// TODO ModifyUserProfile，ModifyHasSetUserProfileStatus 要在同一个事务
	errModifyUserProfile := j.datastore.ModifyUserProfile(ctx, profile)
	if errModifyUserProfile != nil {
		return NewError(ErrDatabase, errors.New("failed to modify user_profile"))
	}
	errModifyHasSetUserProfileStatus := j.datastore.ModifyHasSetUserProfileStatus(ctx, req.UserId)
	if errModifyHasSetUserProfileStatus != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to modify has_set_user_profile status of user %d: %s", req.UserId, errModifyHasSetUserProfileStatus.Error()))
	}
	protoGender, errmapDBGenderToProto := mapDBGenderToProto(userProfile.Gender)
	if errmapDBGenderToProto != nil {
		return NewError(ErrInvalidGender, errmapDBGenderToProto)
	}
	timeStampBirthday, _ := ptypes.TimestampProto(userProfile.Birthday)
	resp.UserProfile = &jinmuidpb.UserProfile{
		Nickname:     userProfile.Nickname,
		BirthdayTime: timeStampBirthday,
		Gender:       protoGender,
		Weight:       userProfile.Weight,
		Height:       userProfile.Height,
	}
	user, errFindUserByUserID := j.datastore.FindUserByUserID(ctx, req.UserId)
	if errFindUserByUserID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find uer by userID %d", req.UserId))
	}
	if user.RegisterType == WechatRegisterType { // 注册方式为微信
		return nil
	}
	reqGetUserSubscriptions := new(subscriptionpb.GetUserSubscriptionsRequest)
	reqGetOrganizationIDByUserID := new(bizcorepb.GetOrganizationIDByUserIDRequest)
	reqGetOrganizationIDByUserID.UserId = req.UserId
	GetOrganizationIDByUserIDReply, err := j.bizSvc.GetOrganizationIDByUserID(ctx, reqGetOrganizationIDByUserID)
	if err != nil {
		return err
	}
	var ownerID int32
	// 如果当前用户不是组织的拥有者，那么就走if里面的逻辑
	// 首选获取当前用户的组织，然后获取该组织的拥有者
	if !GetOrganizationIDByUserIDReply.IsOwner {
		reqGetOwnerIDByOrganizationID := new(bizcorepb.GetOwnerIDByOrganizationIDRequest)
		reqGetOwnerIDByOrganizationID.OrganizationId = GetOrganizationIDByUserIDReply.OrganizationId
		GetOwnerIDByOrganizationIDReply, err := j.bizSvc.GetOwnerIDByOrganizationID(ctx, reqGetOwnerIDByOrganizationID)
		if err != nil {
			return err
		}
		ownerID = GetOwnerIDByOrganizationIDReply.OwnerId
	} else {
		ownerID = req.UserId
	}

	reqGetUserSubscriptions.UserId = ownerID
	// 2.0迁移之前老用户，有没有激活的订阅，自动激活
	respGetUserSubscriptions, errGetUserSubscriptions := j.subscriptionSvc.GetUserSubscriptions(ctx, reqGetUserSubscriptions)
	if errGetUserSubscriptions == nil {
		for _, item := range respGetUserSubscriptions.Subscriptions {
			if !item.IsMigratedActivated && !item.Activated {
				// 激活订阅
				reqActivateSubscription := new(subscriptionpb.ActivateSubscriptionRequest)
				reqActivateSubscription.SubscriptionId = item.SubscriptionId
				_, errActivateSubscription := j.subscriptionSvc.ActivateSubscription(ctx, reqActivateSubscription)
				if errActivateSubscription != nil {
					return errActivateSubscription
				}
			}
		}
	}
	return nil
}

// GetUserProfileByRecordID 通过记录id 获取用户信息
func (j *JinmuIDService) GetUserProfileByRecordID(ctx context.Context, req *jinmuidpb.GetUserProfileByRecordIDRequest, resp *jinmuidpb.GetUserProfileByRecordIDResponse) error {
	userProfile, errFindUserProfileByRecordID := j.datastore.FindUserProfileByRecordID(ctx, req.RecordId)
	if errFindUserProfileByRecordID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find user_profile by record_id %d: %s", req.RecordId, errFindUserProfileByRecordID.Error()))
	}
	// 如果需要token验证
	if !req.IsSkipVerifyToken {
		token, ok := TokenFromContext(ctx)
		if !ok {
			return NewError(ErrInvalidUser, errors.New("failed to get token from context"))
		}
		userID, err := j.datastore.FindUserIDByToken(ctx, token)
		if err != nil {
			return NewError(ErrUserUnauthorized, fmt.Errorf("failed to get userId by token: %s", err.Error()))
		}
		if userID != userProfile.UserID {
			return NewError(ErrInvalidUser, fmt.Errorf("user %d from request and user %d from token are inconsistent", userProfile.UserID, userID))
		}
	}
	resp.UserId = userProfile.UserID
	birthday, _ := ptypes.TimestampProto(userProfile.Birthday)
	protoGender, errmapDBGenderToProto := mapDBGenderToProto(userProfile.Gender)
	if errmapDBGenderToProto != nil {
		return NewError(ErrInvalidGender, errmapDBGenderToProto)
	}
	resp.UserProfile = &jinmuidpb.UserProfile{
		Nickname:        userProfile.Nickname,
		NicknameInitial: userProfile.NicknameInitial,
		BirthdayTime:    birthday,
		Gender:          protoGender,
		Weight:          userProfile.Weight,
		Height:          userProfile.Height,
	}
	return nil
}

func mapProtoGenderToDB(gender generalpb.Gender) (mysqldb.Gender, error) {
	switch gender {
	case generalpb.Gender_GENDER_FEMALE:
		return mysqldb.GenderFemale, nil
	case generalpb.Gender_GENDER_MALE:
		return mysqldb.GenderMale, nil
	case generalpb.Gender_GENDER_INVALID:
		return mysqldb.GenderMale, fmt.Errorf("invalid proto gender %d", generalpb.Gender_GENDER_INVALID)
	case generalpb.Gender_GENDER_UNSET:
		return mysqldb.GenderMale, fmt.Errorf("invalid proto gender %d", generalpb.Gender_GENDER_UNSET)
	}
	return mysqldb.GenderMale, errors.New("invalid proto gender")
}

func mapDBGenderToProto(gender mysqldb.Gender) (generalpb.Gender, error) {
	switch gender {
	case mysqldb.GenderMale:
		return generalpb.Gender_GENDER_MALE, nil
	case mysqldb.GenderFemale:
		return generalpb.Gender_GENDER_FEMALE, nil
	}
	return generalpb.Gender_GENDER_MALE, fmt.Errorf("invalid mysql gender %s", gender)
}

// safeDealPhone 安全处理手机
func safeDealPhone(phone string) string {
	if len(phone) == 0 {
		return phone
	}
	middle := ""
	middleReplace := ""
	for index, value := range phone {
		if index > 3 && index < (len(phone)-3) {
			middle = middle + string(value)
			middleReplace = middleReplace + "*"
		}
	}
	return strings.Replace(phone, middle, middleReplace, -1)
}

// safeDealEmail 安全处理邮箱
func safeDealEmail(email string) string {
	if len(email) == 0 {
		return email
	}
	idx := strings.Index(email, "@")
	if idx == -1 || idx < 3 {
		return email
	}
	middle := ""
	middleReplace := ""
	for index, value := range email {
		if index > 1 && index < (idx-1) {
			middle = middle + string(value)
			middleReplace = middleReplace + "*"
		}
	}
	return strings.Replace(email, middle, middleReplace, -1)
}

// getNicknameInitial 获得昵称的首字母
func getNicknameInitial(nickname string) string {
	nicknameRune := []rune(nickname)
	if len(nicknameRune) == 0 {
		return ""
	}
	nicknamePreffix := nicknameRune[0]
	// 是否为汉字
	if unicode.Is(unicode.Han, nicknamePreffix) {
		a := pinyin.NewArgs()
		a.Style = pinyin.FirstLetter
		nicknameInitial := pinyin.Pinyin(string(nicknamePreffix), a)
		if len(nicknameInitial) == 0 || len(nicknameInitial[0]) == 0 {
			return ""
		}
		return strings.ToUpper(nicknameInitial[0][0])
	}
	// 是否为英文字母
	if unicode.IsLetter(nicknamePreffix) && nicknamePreffix < unicode.MaxASCII {
		return strings.ToUpper(string(nicknamePreffix))
	}

	// 其他字符
	return "~"
}
