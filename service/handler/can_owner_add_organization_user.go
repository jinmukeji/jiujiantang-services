package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/jinmukeji/jiujiantang-services/service/auth"
	corepb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	subscriptionpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/subscription/v1"
	jinmuidpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
)

const (
	// maxSanshuiUser sanshui项目下最大用户数
	maxSanshuiUser = 600000
	// customSanshui sanshui项目下帐号标记
	customSanshui = "custom_sanshui"

	// maxDengyunUser dengyun项目下最大用户数
	maxDengyunUser = 500000
	// customDengyun dengyun项目下帐号标记
	customDengyun = "custom_dengyun"
)

// CanOwnerAddOrganizationUser 是否能够加入用户到组织
func (j *JinmuHealth) CanOwnerAddOrganizationUser(ctx context.Context, req *corepb.CanOwnerAddOrganizationUserRequest, resp *corepb.CanOwnerAddOrganizationUserResponse) error {
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return NewError(ErrInvalidUser, errors.New("failed to get userID from context"))
	}
	// 获取用户信息
	reqGetUserAndProfileInformation := new(jinmuidpb.GetUserAndProfileInformationRequest)
	reqGetUserAndProfileInformation.UserId = userID
	respGetUserAndProfileInformation, errGetUserAndProfileInformation := j.jinmuidSvc.GetUserAndProfileInformation(ctx, reqGetUserAndProfileInformation)
	if errGetUserAndProfileInformation != nil {
		return errGetUserAndProfileInformation
	}
	o, errFindOrganizationByUserID := j.datastore.FindOrganizationByUserID(ctx, int(userID))
	if errFindOrganizationByUserID != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to find organization by user id [%d]: %s", userID, errFindOrganizationByUserID.Error()))
	}

	switch respGetUserAndProfileInformation.CustomizedCode {
	case customDengyun:
		if ableDengyun, err := j.CanCustomizedOwnerAddOrganizationUser(ctx, respGetUserAndProfileInformation.CustomizedCode, customDengyun, maxDengyunUser); err == nil {
			resp.Able = ableDengyun
		} else {
			return NewError(ErrAddUserFailure, err)
		}
	case customSanshui:
		if ableSanshui, err := j.CanCustomizedOwnerAddOrganizationUser(ctx, respGetUserAndProfileInformation.CustomizedCode, customSanshui, maxSanshuiUser); err == nil {
			resp.Able = ableSanshui
		} else {
			return NewError(ErrAddUserFailure, err)
		}
	default:
		// 获取组织下的用户数量
		count, errGetExistingUserCountByOrganizationID := j.datastore.GetExistingUserCountByOrganizationID(ctx, o.OrganizationID)
		if errGetExistingUserCountByOrganizationID != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to get existing user count by organizationID %d: %s", o.OrganizationID, errGetExistingUserCountByOrganizationID.Error()))
		}
		// 获取当前拥有者的订阅
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

		// 获取该订阅下的用户数量
		reqGetUserCountOfSubscription := new(subscriptionpb.GetUserCountOfSubscriptionRequest)
		reqGetUserCountOfSubscription.SubscriptionId = selectedSubscription.SubscriptionId
		respGetUserCountOfSubscription, errGetUserCount := j.subscriptionSvc.GetUserCountOfSubscription(ctx, reqGetUserCountOfSubscription)
		if errGetUserCount != nil {
			return NewError(ErrInactivatedSubscription, errGetUserCount)
		}
		if respGetUserCountOfSubscription.UserCount >= selectedSubscription.MaxUserLimits {
			return NewError(ErrExceedSubscriptionUserQuotaLimit, fmt.Errorf("current user count exceed the count of subscription %d limit", reqGetUserCountOfSubscription.SubscriptionId))
		}
		// 判断用户数目是否达到用户的最大限制
		if selectedSubscription.MaxUserLimits > int32(count) {
			resp.Able = true
		}

	}
	return nil
}

// CanCustomizedOwnerAddOrganizationUser 自定义用户是否能够加入用户到组织
func (j *JinmuHealth) CanCustomizedOwnerAddOrganizationUser(ctx context.Context, ownerCustomizedCode string, customizedCode string, maxUser int) (bool, error) {
	if ownerCustomizedCode == customizedCode {
		count, errCountUserByCustomizedCode := j.datastore.CountUserByCustomizedCode(ctx, customizedCode)
		if errCountUserByCustomizedCode != nil {
			return false, NewError(ErrDatabase, fmt.Errorf("failed to get count user by cutomized code %s: %s", customizedCode, errCountUserByCustomizedCode.Error()))
		}
		if count < maxUser {
			return true, nil
		}
		return false, nil
	}
	return false, nil
}
