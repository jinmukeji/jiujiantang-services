package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/go-pkg/age"

	"github.com/jinmukeji/jiujiantang-services/pkg/rpc"
	"github.com/jinmukeji/jiujiantang-services/service/auth"
	"github.com/jinmukeji/jiujiantang-services/service/mysqldb"
	jinmuidpb "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	corepb "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	subscriptionpb "github.com/jinmukeji/proto/gen/micro/idl/jm/subscription/v1"
)

const (
	// TrialSubscriptionExpirationDuration 试用订阅期限是1个月
	TrialSubscriptionExpirationDuration = time.Hour * 24 * 30
	// TrialSubscriptionMaxUserLimit 试用订阅最大用户上限
	TrialSubscriptionMaxUserLimit = 2
	// Inactive 未激活
	Inactive int = 0
	// Activated 已激活
	Activated int = 1
	// maxUserQuerySize 用户最大区间数
	maxUserQuerySize = 100
	// minUserQuerySize 用户最小区间数
	minUserQuerySize = 1
)

// OwnerCreateOrganization 创建组织
func (j *JinmuHealth) OwnerCreateOrganization(ctx context.Context, req *corepb.OwnerCreateOrganizationRequest, resp *corepb.OwnerCreateOrganizationResponse) error {
	l := rpc.ContextLogger(ctx)

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return NewError(ErrInvalidUser, errors.New("failed to get userID from context"))
	}
	organizationCount, err := j.datastore.GetOrganizationCountByOwnerID(ctx, int(userID))
	if err != nil {
		l.WithError(err).Warn("database error")
		return NewError(ErrDatabase, fmt.Errorf("failed to get organization count by ownerID: %d :%s", int(userID), err.Error()))
	}
	if organizationCount > 0 {
		return NewError(ErrOrganizationCountExceedsMaxLimits, fmt.Errorf("count of organization of owner %d exceeds the max limits", int(userID)))
	}
	now := time.Now()
	owner, _ := j.datastore.FindUserByUserID(ctx, int(userID))
	o := &mysqldb.Organization{
		Name:           req.Profile.Name,
		State:          req.Profile.Address.State,
		City:           req.Profile.Address.City,
		Street:         req.Profile.Address.Street,
		Phone:          req.Profile.Phone,
		Contact:        req.Profile.Contact,
		Type:           req.Profile.Type,
		Country:        req.Profile.Address.Country,
		District:       req.Profile.Address.District,
		PostalCode:     req.Profile.Address.PostalCode,
		Email:          req.Profile.Email,
		IsValid:        1,
		CustomizedCode: owner.CustomizedCode,
	}
	// TODO: CreateOrganization，CreateOrganizationOwner 要在同一个事务
	if errCreateOrganization := j.datastore.CreateOrganization(ctx, o); errCreateOrganization != nil {
		l.WithError(errCreateOrganization).Warn("database error")
		return NewError(ErrDatabase, fmt.Errorf("failed to create organization: %s", errCreateOrganization.Error()))
	}
	if errCreateOrganizationOwner := j.datastore.CreateOrganizationOwner(ctx, &mysqldb.OrganizationOwner{
		OrganizationID: o.OrganizationID,
		OwnerID:        int(userID),
		CreatedAt:      now,
		UpdatedAt:      now,
	}); errCreateOrganizationOwner != nil {
		l.WithError(errCreateOrganizationOwner).Warn("database error")
		return NewError(ErrDatabase, fmt.Errorf("failed to create organizationOwner: %s", errCreateOrganizationOwner.Error()))
	}
	users := make([]*mysqldb.OrganizationUser, 1)
	users[0] = &mysqldb.OrganizationUser{
		OrganizationID: o.OrganizationID,
		UserID:         int(userID),
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	errCreateOrganizationUsers := j.datastore.CreateOrganizationUsers(ctx, users)
	if errCreateOrganizationUsers != nil {
		l.WithError(errCreateOrganizationUsers).Warn("database error")
		return NewError(ErrDatabase, fmt.Errorf("failed to create organizationUsers: %s", errCreateOrganizationUsers.Error()))
	}
	resp.OrganizationId = int32(o.OrganizationID)
	resp.Profile = req.Profile
	return nil
}

// OwnerGetOrganizations 查看拥有的组织
func (j *JinmuHealth) OwnerGetOrganizations(ctx context.Context, req *corepb.OwnerGetOrganizationsRequest, repl *corepb.OwnerGetOrganizationsResponse) error {
	l := rpc.ContextLogger(ctx)

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return NewError(ErrInvalidUser, errors.New("failed to get userID from context"))
	}
	organizations, err := j.datastore.FindOrganizationsByOwner(ctx, int(userID))
	if err != nil {
		l.WithError(err).Warn("database error")
		return NewError(ErrDatabase, fmt.Errorf("failed to find organization by owner:%d: %s", int(userID), err.Error()))
	}
	repl.Organizations = make([]*corepb.Organization, len(organizations))
	for i, o := range organizations {
		subscriptions, _ := j.datastore.FindSubscriptionsByOrganizationID(ctx, o.OrganizationID)
		repl.Organizations[i] = &corepb.Organization{
			OrganizationId: int32(o.OrganizationID),
			Profile: &corepb.OrganizationProfile{
				Name: o.Name,
				Address: &corepb.Address{
					Street:     o.Street,
					City:       o.City,
					State:      o.State,
					Country:    o.Country,
					District:   o.District,
					PostalCode: o.PostalCode,
				},
				Phone:   o.Phone,
				Contact: o.Contact,
				Type:    o.Type,
				Email:   o.Email,
			},
		}
		if len(subscriptions) != 0 {
			subscription := subscriptions[0]
			activatedAt, _ := ptypes.TimestampProto(subscription.ActivatedAt)
			expiredAt, _ := ptypes.TimestampProto(subscription.ExpiredAt)
			totalUserCount, _ := j.datastore.GetExistingUserCountByOrganizationID(ctx, o.OrganizationID)
			createdAt, _ := ptypes.TimestampProto(subscription.CreatedAt)
			repl.Organizations[i].Subscription = &corepb.Subscription{
				SubscriptionType: subscriptionpb.SubscriptionType(subscription.SubscriptionType),
				ActiveTime:       activatedAt,
				ExpiredTime:      expiredAt,
				Active:           subscription.Active == 1,
				TotalUserCount:   int32(totalUserCount),
				MaxUserLimits:    int32(subscription.MaxUserLimits),
				CreatedTime:      createdAt,
			}
		}
	}
	return nil
}

// OwnerGetOrganizationSubscription Owner查看组织订阅
func (j *JinmuHealth) OwnerGetOrganizationSubscription(ctx context.Context, req *corepb.OwnerGetOrganizationSubscriptionRequest, repl *corepb.OwnerGetOrganizationSubscriptionResponse) error {
	l := rpc.ContextLogger(ctx)

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return NewError(ErrInvalidUser, errors.New("failed to find userID from context"))
	}
	ok, err := j.datastore.CheckOrganizationOwner(ctx, int(userID), int(req.OrganizationId))
	if err != nil || !ok {
		l.WithError(err).Warn("database error")
		return NewError(ErrDatabase, fmt.Errorf("failed to check organization %d and owner %d: %s", req.OrganizationId, userID, err.Error()))
	}
	subscriptions, errFindSubscriptionsByOrganizationID := j.datastore.FindSubscriptionsByOrganizationID(ctx, int(req.OrganizationId))
	if errFindSubscriptionsByOrganizationID != nil {
		l.WithError(errFindSubscriptionsByOrganizationID).Warn("database error")
		return NewError(ErrDatabase, fmt.Errorf("failed to find subscription by organizationID: %d: %s", req.OrganizationId, errFindSubscriptionsByOrganizationID.Error()))
	}
	if len(subscriptions) != 0 {
		subscription := subscriptions[0]
		activatedAt, _ := ptypes.TimestampProto(subscription.ActivatedAt)
		expiredAt, _ := ptypes.TimestampProto(subscription.ExpiredAt)
		totalUserCount, _ := j.datastore.GetExistingUserCountByOrganizationID(ctx, int(req.OrganizationId))
		createdAt, _ := ptypes.TimestampProto(subscription.CreatedAt)
		repl.Subscription = &corepb.Subscription{
			ActiveTime:       activatedAt,
			ExpiredTime:      expiredAt,
			Active:           subscription.Active == 1,
			TotalUserCount:   int32(totalUserCount),
			MaxUserLimits:    int32(subscription.MaxUserLimits),
			CreatedTime:      createdAt,
			SubscriptionType: subscriptionpb.SubscriptionType(subscription.SubscriptionType),
		}
	}
	return nil
}

// OwnerAddOrganizationUsers 拥有者向组织下添加用户
func (j *JinmuHealth) OwnerAddOrganizationUsers(ctx context.Context, req *corepb.OwnerAddOrganizationUsersRequest, repl *corepb.OwnerAddOrganizationUsersResponse) error {
	l := rpc.ContextLogger(ctx)

	// TODO 这里的添加CreateOrganizationUsers和AddUsersIntoSubscription应该写在同一个事务中
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return NewError(ErrInvalidUser, errors.New("failed to find userID from context"))
	}
	ok, err := j.datastore.CheckOrganizationOwner(ctx, int(userID), int(req.OrganizationId))
	if err != nil || !ok {
		l.WithError(err).Warn("database error")
		return NewError(ErrDatabase, fmt.Errorf("failed to check organization %d and owner %d: %s", int(req.OrganizationId), int(userID), err.Error()))
	}
	users := make([]*mysqldb.OrganizationUser, len(req.UserIdList))
	now := time.Now()
	for idx, uid := range req.UserIdList {
		users[idx] = &mysqldb.OrganizationUser{
			OrganizationID: int(req.OrganizationId),
			UserID:         int(uid),
			CreatedAt:      now,
			UpdatedAt:      now,
		}
	}
	// TODO: CreateOrganizationUsers,AddUsersIntoSubscription 要在同一个事务
	errCreateOrganizationUsers := j.datastore.CreateOrganizationUsers(ctx, users)
	if errCreateOrganizationUsers != nil {
		l.WithError(errCreateOrganizationUsers).Warn("database error")
		return NewError(ErrDatabase, fmt.Errorf("failed to create organizationUsers: %s", errCreateOrganizationUsers.Error()))
	}
	reqGetUserSubscriptions := new(subscriptionpb.GetUserSubscriptionsRequest)
	reqGetUserSubscriptions.UserId = userID
	respGetUserSubscriptions, errGetUserSubscriptions := j.subscriptionSvc.GetUserSubscriptions(ctx, reqGetUserSubscriptions)
	if errGetUserSubscriptions != nil {
		return errGetUserSubscriptions
	}
	// 获取当前正在使用的订阅
	selectedSubscription := new(subscriptionpb.Subscription)
	for _, item := range respGetUserSubscriptions.Subscriptions {
		if item.IsSelected {
			selectedSubscription = item
		}
	}
	// 将用户添加到订阅中
	reqAddUsersIntoSubscription := new(subscriptionpb.AddUsersIntoSubscriptionRequest)
	reqAddUsersIntoSubscription.OwnerId = userID
	reqAddUsersIntoSubscription.SubscriptionId = selectedSubscription.SubscriptionId
	reqAddUsersIntoSubscription.UserIdList = req.UserIdList
	_, errAddUsersIntoSubscription := j.subscriptionSvc.AddUsersIntoSubscription(ctx, reqAddUsersIntoSubscription)
	if errAddUsersIntoSubscription != nil {
		return errAddUsersIntoSubscription
	}
	usersList, errFindOrganizationUsersIDList := j.datastore.FindOrganizationUsersIDList(ctx, int(req.OrganizationId))
	if errFindOrganizationUsersIDList != nil {
		l.WithError(errFindOrganizationUsersIDList).Warn("database error")
		return NewError(ErrDatabase, fmt.Errorf("failed to get user lists of the organization: %s", errFindOrganizationUsersIDList.Error()))
	}
	repl.User = make([]*corepb.User, len(usersList))

	for idx, u := range usersList {
		// 获取用户的所有信息
		reqGetUserAndProfileInformation := new(jinmuidpb.GetUserAndProfileInformationRequest)
		reqGetUserAndProfileInformation.UserId = int32(u)
		reqGetUserAndProfileInformation.IsSkipVerifyToken = true
		respGetUserAndProfileInformation, errGetUserAndProfileInformation := j.jinmuidSvc.GetUserAndProfileInformation(ctx, reqGetUserAndProfileInformation)
		if errGetUserAndProfileInformation != nil {
			return errGetUserAndProfileInformation
		}
		birthday, _ := ptypes.Timestamp(respGetUserAndProfileInformation.Profile.BirthdayTime)
		repl.User[idx] = &corepb.User{
			UserId:             int32(u),
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
				Email:           respGetUserAndProfileInformation.SecureEmail,
				Phone:           respGetUserAndProfileInformation.SigninPhone,
				Remark:          respGetUserAndProfileInformation.Remark,
				UserDefinedCode: respGetUserAndProfileInformation.UserDefinedCode,
			},
			IsRemovable: !(int32(u) == userID),
		}

	}

	return nil
}

// OwnerDeleteOrganizationUsers 拥有者组织删除用户
func (j *JinmuHealth) OwnerDeleteOrganizationUsers(ctx context.Context, req *corepb.OwnerDeleteOrganizationUsersRequest, repl *corepb.OwnerDeleteOrganizationUsersResponse) error {
	l := rpc.ContextLogger(ctx)

	// TODO 这里的添加DeleteOrganizationUsers和DeleteUsersFromSubscription应该写在同一个事务中
	ownerID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return NewError(ErrInvalidUser, errors.New("failed to get userID from context"))
	}
	ok, err := j.datastore.CheckOrganizationOwner(ctx, int(ownerID), int(req.OrganizationId))
	if err != nil || !ok {
		l.WithError(err).Warn("database error")
		return NewError(ErrDatabase, fmt.Errorf("failed to check organization and owner: %s", err.Error()))
	}
	for _, userID := range req.UserIdList {
		if userID == ownerID {
			return NewError(ErrCannotDeleteOwner, errors.New("cannot delete organization owner"))
		}
	}
	if err := j.datastore.DeleteOrganizationUsers(ctx, req.UserIdList, req.OrganizationId); err != nil {
		l.WithError(err).Warn("database error")
		return NewError(ErrDatabase, fmt.Errorf("failed to delete connection relation of organization %d and users %d: %s", req.OrganizationId, req.UserIdList, err.Error()))
	}
	// 获取当前拥有者的订阅
	reqGetUserSubscriptions := new(subscriptionpb.GetUserSubscriptionsRequest)
	reqGetUserSubscriptions.UserId = ownerID
	respGetUserSubscriptions, errGetUserSubscriptions := j.subscriptionSvc.GetUserSubscriptions(ctx, reqGetUserSubscriptions)
	if errGetUserSubscriptions != nil {
		return errGetUserSubscriptions
	}
	// 获取当前正在使用的订阅
	selectedSubscription := new(subscriptionpb.Subscription)
	for _, item := range respGetUserSubscriptions.Subscriptions {
		if item.IsSelected {
			selectedSubscription = item
		}
	}
	reqDeleteUsersFromSubscription := new(subscriptionpb.DeleteUsersFromSubscriptionRequest)
	reqDeleteUsersFromSubscription.OwnerId = ownerID
	reqDeleteUsersFromSubscription.SubscriptionId = selectedSubscription.SubscriptionId
	reqDeleteUsersFromSubscription.UserIdList = req.UserIdList
	_, errDeleteUsersFromSubscription := j.subscriptionSvc.DeleteUsersFromSubscription(ctx, reqDeleteUsersFromSubscription)
	if errDeleteUsersFromSubscription != nil {
		return errDeleteUsersFromSubscription
	}
	repl.Tip = "delete users from Organization successful"
	return nil
}

// OwnerGetOrganizationUsers Owner 获取组织中的 User
func (j *JinmuHealth) OwnerGetOrganizationUsers(ctx context.Context, req *corepb.OwnerGetOrganizationUsersRequest, repl *corepb.OwnerGetOrganizationUsersResponse) error {
	l := rpc.ContextLogger(ctx)

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return NewError(ErrInvalidUser, errors.New("failed to find userID from context"))
	}
	ok, err := j.datastore.CheckOrganizationOwner(ctx, int(userID), int(req.OrganizationId))
	if err != nil || !ok {
		l.WithError(err).Warn("database error")
		return NewError(ErrDatabase, fmt.Errorf("failed to check organization %d and owner %d: %s", int(req.OrganizationId), int(userID), err.Error()))
	}
	isValid := j.datastore.CheckOrganizationIsValid(ctx, int(req.OrganizationId))
	if !isValid {
		return NewError(ErrInvalidOrganization, fmt.Errorf("organization %d is invalid", req.OrganizationId))
	}
	offset := req.Offset
	size := req.Size
	keyword := req.Keyword
	if size != -1 && (size > maxUserQuerySize || size < minUserQuerySize) {
		return NewError(ErrOrganizationQueryUserExceedsLimit, errors.New("size exceeds the maximum or minimum limit"))
	}
	users, err := j.searchUserByKeyword(ctx, req.OrganizationId, keyword, size, offset)
	if err != nil {
		l.WithError(err).Warn("database error")
		return NewErrorCause(ErrDatabase, errors.New("database error"), err.Error())
	}
	repl.UserList = make([]*corepb.User, len(users))
	for idx, u := range users {
		birth, err := ptypes.TimestampProto(u.Birthday)
		if err != nil {
			birth = nil
		}
		gender, errMapDBGenderToProto := mapDBGenderToProto(u.Gender)
		if errMapDBGenderToProto != nil {
			return NewError(ErrInvalidGender, errMapDBGenderToProto)
		}
		repl.UserList[idx] = &corepb.User{
			UserId:             int32(u.UserID),
			Username:           u.Username,
			RegisterType:       u.RegisterType,
			IsProfileCompleted: u.IsProfileCompleted,
			IsRemovable:        u.IsRemovable,
			Profile: &corepb.UserProfile{
				Nickname:        u.Nickname,
				NicknameInitial: u.NicknameInitial,
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
	}
	return nil
}

// searchUserByKeyword 通过关键字搜索用户
func (j *JinmuHealth) searchUserByKeyword(ctx context.Context, organizationID int32, keyword string, size int32, offset int32) ([]*mysqldb.User, error) {
	if keyword != "" {
		users, err := j.datastore.FindOrganizationUsersByKeyword(ctx, organizationID, keyword, size, offset)
		if err != nil {
			return []*mysqldb.User{}, NewError(ErrDatabase, fmt.Errorf("failed to find users of organization %d by keyword %s,size %d, offset %d: %s", organizationID, keyword, size, offset, err.Error()))
		}

		return users, nil
	}
	users, err := j.datastore.FindOrganizationUsersByOffset(ctx, organizationID, size, offset)
	if err != nil {
		return []*mysqldb.User{}, NewError(ErrDatabase, fmt.Errorf("failed to find users of organization %d by size %d, offset %d: %s", organizationID, size, offset, err.Error()))
	}
	return users, nil
}

// GetOrganizationIDByUserID 根据userID查找对应的organizationID
func (j *JinmuHealth) GetOrganizationIDByUserID(ctx context.Context, req *corepb.GetOrganizationIDByUserIDRequest, repl *corepb.GetOrganizationIDByUserIDResponse) error {
	userID := req.UserId
	o, err := j.datastore.FindOrganizationByUserID(ctx, int(userID))
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find organization by userID %d: %s", int(userID), err.Error()))
	}
	repl.OrganizationId = int32(o.OrganizationID)
	isOwner, err := j.datastore.CheckOrganizationOwner(ctx, int(userID), o.OrganizationID)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to check organization %d and ownerID %d: %s", o.OrganizationID, int(userID), err.Error()))
	}
	repl.IsOwner = isOwner
	return nil
}
