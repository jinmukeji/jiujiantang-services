package handler

import (
	"context"
	"errors"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/jinmukeji/gf-api2/jinmuid/mysqldb"
	crypto "github.com/jinmukeji/go-pkg/crypto/encrypt/legacy"
	"github.com/jinmukeji/go-pkg/crypto/rand"
	bizcorepb "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	subscriptionpb "github.com/jinmukeji/proto/gen/micro/idl/jm/subscription/v1"

	"fmt"

	jinmuidpb "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
)

const (
	defaultRegisterSource = "手机验证码注册"
	defaultRegisterType   = "phone"
	organizationName      = "金姆健康科技有限公司"
	organizationStreet    = "江苏常州天宁区关河东路66号九洲环宇大厦1501室"
	organizationPhone     = "0519-81180075"
	organizationEmail     = "information@jinmuhealth.com"
	organizationState     = "江苏省"
	organizationCity      = "常州市"
	organizationDistrict  = "天宁区"
	organizationType      = "养生"
	// initNotificationPreference 通知初始状态
	initNotificationPreference = true
)

// UserSignUpByPhone 手机号注册
func (j *JinmuIDService) UserSignUpByPhone(ctx context.Context, req *jinmuidpb.UserSignUpByPhoneRequest, resp *jinmuidpb.UserSignUpByPhoneResponse) error {
	// TODO: 注册用户和创建组织要在同一个事务
	// 该手机是否注册过
	exsit, err := j.datastore.ExistPhone(ctx, req.Phone, req.NationCode)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to check the existence of phone %s%s", req.NationCode, req.Phone))
	}
	if exsit {
		return NewError(ErrExistRegisteredPhone, fmt.Errorf("phone %s%s doesn't exist", req.NationCode, req.Phone))
	}
	isValid, errVerifyVerificationNumberByPhone := j.datastore.VerifyVerificationNumberByPhone(ctx, req.VerificationNumber, req.Phone, req.NationCode)
	if errVerifyVerificationNumberByPhone != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to verify verification number by phone %s%s: %s", req.NationCode, req.Phone, errVerifyVerificationNumberByPhone.Error()))
	}
	if !isValid {
		return NewError(ErrInvalidVerificationNumber, fmt.Errorf("verification number is invalid"))
	}
	errSetVerificationNumberAsUsed := j.datastore.SetVerificationNumberAsUsed(ctx, mysqldb.VerificationPhone, req.VerificationNumber)
	if errSetVerificationNumberAsUsed != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to set verification number of email to the status of used %s: %s", req.VerificationNumber, errSetVerificationNumberAsUsed.Error()))
	}
	// 判断密码
	if req.PlainPassword == "" {
		return NewError(ErrEmptyPassword, errors.New("password is empty"))
	}
	if !checkPasswordFormat(req.PlainPassword) {
		return NewError(ErrWrongFormatOfPassword, errors.New("password format is wrong"))
	}
	seed, _ := rand.RandomStringWithMask(rand.MaskLetterDigits, 4)
	helper := crypto.NewPasswordCipherHelper()
	encryptedPassword := helper.Encrypt(req.PlainPassword, seed, j.encryptKey)

	// 设置语言
	language, errmapProtoLanguageToDB := mapProtoLanguageToDB(req.Language)
	if errmapProtoLanguageToDB != nil {
		return NewError(ErrInvalidLanguage, errmapProtoLanguageToDB)
	}
	mysqlLanguage, errmapLanguageToDB := mapLanguageToDB(language)
	if errmapLanguageToDB != nil {
		return NewError(ErrInvalidLanguage, errmapLanguageToDB)
	}
	now := time.Now()
	user := &mysqldb.User{
		SigninPhone:             req.Phone,
		HasSetPhone:             true,
		RegisterType:            defaultRegisterType,
		RegisterSource:          defaultRegisterSource,
		NationCode:              req.NationCode,
		RegisterTime:            now,
		LatestUpdatedPhoneAt:    &now,
		IsActivated:             true,
		Language:                mysqlLanguage,
		HasSetLanguage:          true,
		HasSetUserProfile:       true,
		IsProfileCompleted:      true,
		EncryptedPassword:       encryptedPassword,
		Seed:                    seed,
		HasSetPassword:          true,
		LatestUpdatedPasswordAt: &now,
		ActivatedAt:             &now,
		CreatedAt:               now,
		UpdatedAt:               now,
	}

	birthday, _ := ptypes.Timestamp(req.Profile.BirthdayTime)
	mysqlGender, errmapProtoGenderToDB := mapProtoGenderToDB(req.Profile.Gender)
	if errmapProtoGenderToDB != nil {
		return NewError(ErrInvalidGender, errmapProtoGenderToDB)
	}
	profile := &mysqldb.UserProfile{
		Nickname:        req.Profile.Nickname,
		NicknameInitial: getNicknameInitial(req.Profile.Nickname),
		Gender:          mysqlGender,
		Weight:          req.Profile.Weight,
		Height:          req.Profile.Height,
		Birthday:        birthday,
	}

	errCreateUserAndUserProfile := j.datastore.CreateUserAndUserProfile(ctx, user, profile)
	if errCreateUserAndUserProfile != nil {
		return NewError(ErrDatabase, errors.New("fail to create user and user profile"))
	}

	errCreateUserPreferences := j.datastore.CreateUserPreferences(ctx, user.UserID)
	if errCreateUserPreferences != nil {
		return NewError(ErrDatabase, fmt.Errorf("fail to create user and user preferences of user %d: %s", user.UserID, errCreateUserPreferences.Error()))
	}
	token := uuid.New().String()
	tk, err := j.datastore.CreateToken(ctx, token, user.UserID, TokenAvailableDuration)
	if err != nil {
		return NewError(ErrGetAccessTokenFailure, fmt.Errorf("failed to create the access token of user %d: %s", user.UserID, err.Error()))
	}
	notificationPreferences := &mysqldb.NotificationPreferences{
		UserID:                 user.UserID,
		PhoneEnabled:           initNotificationPreference,
		WechatEnabled:          initNotificationPreference,
		WeiboEnabled:           initNotificationPreference,
		PhoneEnabledUpdatedAt:  now,
		WechatEnabledUpdatedAt: now,
		WeiboEnabledUpdatedAt:  now,
		CreatedAt:              now,
		UpdatedAt:              now,
	}
	errCreateNotificationPreferences := j.datastore.CreateNotificationPreferences(ctx, notificationPreferences)
	if errCreateNotificationPreferences != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to create notification preferences of user %d: %s", user.UserID, errCreateNotificationPreferences.Error()))
	}
	resp.UserId = user.UserID
	resp.AccessToken = tk.Token
	// 创建组织
	reqCreateOrganizationRequest := new(bizcorepb.OwnerCreateOrganizationRequest)
	reqCreateOrganizationRequest.Profile = &bizcorepb.OrganizationProfile{
		Name:  organizationName,
		Phone: organizationPhone,
		Type:  organizationType,
		Email: organizationEmail,
		Address: &bizcorepb.Address{
			State:    organizationState,
			Street:   organizationStreet,
			City:     organizationCity,
			District: organizationDistrict,
		},
	}
	ctx = AddContextToken(ctx, token)
	_, errOwnerCreateOrganization := j.bizSvc.OwnerCreateOrganization(ctx, reqCreateOrganizationRequest)
	if errOwnerCreateOrganization != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to create organization: %s", errOwnerCreateOrganization.Error()))
	}

	reqGetUserSubscriptions := new(subscriptionpb.GetUserSubscriptionsRequest)
	reqGetUserSubscriptions.UserId = user.UserID
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

func mapLanguageToDB(language string) (mysqldb.Language, error) {
	switch language {
	case LanguageSimpleChinese:
		return mysqldb.LanguageSimpleChinese, nil
	case LanguageTraditionalChinese:
		return mysqldb.LanguageTraditionalChinese, nil
	case LanguageEnglish:
		return mysqldb.LanguageEnglish, nil
	}
	return mysqldb.LanguageInvalid, fmt.Errorf("invalid string language %s", language)
}
