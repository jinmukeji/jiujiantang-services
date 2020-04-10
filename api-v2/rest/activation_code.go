package rest

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	corepb "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	subscriptionpb "github.com/jinmukeji/proto/gen/micro/idl/jm/subscription/v1"
	"github.com/kataras/iris/v12"
)

// ActivationCodeInfo 激活码信息
type ActivationCodeInfo struct {
	MaxUserCount     int32  `json:"max_user_count"`
	ContractYear     int32  `json:"contract_year"`
	SubscriptionType string `json:"subscription_type"`
}

// ActivationCode 激活码
type ActivationCode struct {
	Code string `json:"code"`
}

// GetActivationCodeContent 得到激活码内容
func (h *v2Handler) GetActivationCodeInfo(ctx iris.Context) {
	req := new(subscriptionpb.GetSubscriptionActivationCodeInfoRequest)
	var body ActivationCode
	errReadJSON := ctx.ReadJSON(&body)
	if errReadJSON != nil {
		// writeError(ctx, Wrap)
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", errReadJSON), false)
		return
	}
	req.Code = body.Code
	resp, err := h.rpcSubscriptionManagerSvc.GetSubscriptionActivationCodeInfo(newRPCContext(ctx), req)
	if err != nil {
		writeRPCInternalError(ctx, err, false)
		return
	}
	subscriptionType, errmapProtoSubscriptionTypeToRest := mapProtoSubscriptionTypeToRest(resp.Type)
	if errmapProtoSubscriptionTypeToRest != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errmapProtoSubscriptionTypeToRest), false)
		return
	}
	rest.WriteOkJSON(ctx, ActivationCodeInfo{
		MaxUserCount:     resp.MaxUserCount,
		ContractYear:     resp.ContractYear,
		SubscriptionType: subscriptionType,
	})
}

// UseActivationCodeInfo 使用激活码返回内容
type UseActivationCodeInfo struct {
	UserSubscriptionStatus SubscriptionStateDescription `json:"user_subscription_status"` // 订阅状态
	TotalUserCount         int32                        `json:"total_user_count"`         // 添加人数
	MaxUserCount           int32                        `json:"max_user_count"`           // 最大人数
	ExpiredAt              time.Time                    `json:"expired_at"`               // 到期时间
	SubscriptionType       string                       `json:"subscription_type"`        // 会员类型
}

// UseSubscriptionActivationCode 使用激活码
func (h *v2Handler) UseSubscriptionActivationCode(ctx iris.Context) {
	userID, err := ctx.Params().GetInt("user_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	req := new(subscriptionpb.UseSubscriptionActivationCodeRequest)
	req.UserId = int32(userID)
	var body ActivationCode
	errReadJSON := ctx.ReadJSON(&body)
	if errReadJSON != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", errReadJSON), false)
		return
	}
	req.Code = body.Code
	_, errUseSubscriptionActivationCode := h.rpcSubscriptionManagerSvc.UseSubscriptionActivationCode(newRPCContext(ctx), req)
	if errUseSubscriptionActivationCode != nil {
		writeRPCInternalError(ctx, errUseSubscriptionActivationCode, false)
		return
	}
	// 获取当前用户的组织ID
	reqOwnerGetOrganizations := new(corepb.OwnerGetOrganizationsRequest)
	respOwnerGetOrganizations, errOwnerGetOrganizations := h.rpcSvc.OwnerGetOrganizations(
		newRPCContext(ctx), reqOwnerGetOrganizations,
	)
	if errOwnerGetOrganizations != nil {
		writeRPCInternalError(ctx, err, true)
		return
	}
	if len(respOwnerGetOrganizations.Organizations) != 1 {
		writeError(
			ctx,
			wrapError(ErrInvalidOrganizationCount, "", errors.New("invalid organization count of current user")),
			false,
		)
		return
	}

	organizationID := respOwnerGetOrganizations.Organizations[0].OrganizationId
	// 获取组织的订阅
	reqOwnerGetOrganizationSubscription := new(corepb.OwnerGetOrganizationSubscriptionRequest)
	reqOwnerGetOrganizationSubscription.OrganizationId = organizationID
	respOwnerGetOrganizationSubscription, errOwnerGetOrganizationSubscription := h.rpcSvc.OwnerGetOrganizationSubscription(
		newRPCContext(ctx), reqOwnerGetOrganizationSubscription,
	)
	if errOwnerGetOrganizationSubscription != nil {
		writeRPCInternalError(ctx, errOwnerGetOrganizationSubscription, true)
		return
	}
	// 获取订阅状态
	var userSubscriptionStatus SubscriptionStateDescription
	reqGetUserSubscriptions := new(subscriptionpb.GetUserSubscriptionsRequest)
	reqGetUserSubscriptions.UserId = int32(userID)
	respGetUserSubscriptions, errGetUserSubscriptions := h.rpcSubscriptionManagerSvc.GetUserSubscriptions(newRPCContext(ctx), reqGetUserSubscriptions)
	if errGetUserSubscriptions != nil {
		if errGetUserSubscriptions.Error() == "[errcode:20000] failed to get subscription" {
			userSubscriptionStatus = Unsubscribed
		} else {
			writeRPCInternalError(ctx, errGetUserSubscriptions, false)
			return
		}
	}
	selectedSubscription := new(subscriptionpb.Subscription)
	for _, item := range respGetUserSubscriptions.Subscriptions {
		if item.IsSelected {
			selectedSubscription = item
		}
	}
	expiredAt, _ := ptypes.Timestamp(selectedSubscription.ExpiredTime)
	currentExpiredAt, _ := ptypes.Timestamp(respOwnerGetOrganizationSubscription.Subscription.ExpiredTime)
	if expiredAt.Before(time.Now()) {
		userSubscriptionStatus = Expired
	}
	if expiredAt.After(time.Now()) {
		userSubscriptionStatus = NotExpired
	}
	stringSubscription, errmapProtoSubscriptionTypeToRest := mapProtoSubscriptionTypeToRest(respOwnerGetOrganizationSubscription.Subscription.SubscriptionType)
	if errmapProtoSubscriptionTypeToRest != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errmapProtoSubscriptionTypeToRest), false)
		return
	}
	rest.WriteOkJSON(ctx, UseActivationCodeInfo{
		UserSubscriptionStatus: userSubscriptionStatus,
		TotalUserCount:         respOwnerGetOrganizationSubscription.Subscription.TotalUserCount,
		MaxUserCount:           respOwnerGetOrganizationSubscription.Subscription.MaxUserLimits,
		ExpiredAt:              currentExpiredAt,
		SubscriptionType:       stringSubscription,
	})
}

// mapProtoSubscriptionTypeToRest 把 proto 文件中订阅信息类型转化成 rest 的 string 类型
func mapProtoSubscriptionTypeToRest(subscriptionType subscriptionpb.SubscriptionType) (string, error) {
	switch subscriptionType {
	case subscriptionpb.SubscriptionType_SUBSCRIPTION_TYPE_INVALID:
		return "", fmt.Errorf("invalid proto subscription type %d", subscriptionType)
	case subscriptionpb.SubscriptionType_SUBSCRIPTION_TYPE_UNSET:
		return "", fmt.Errorf("invalid proto subscription type %d", subscriptionType)
	case subscriptionpb.SubscriptionType_SUBSCRIPTION_TYPE_CUSTOMIZED_VERSION:
		return CustomizedVersion, nil
	case subscriptionpb.SubscriptionType_SUBSCRIPTION_TYPE_TRIAL_VERSION:
		return TrialVersion, nil
	case subscriptionpb.SubscriptionType_SUBSCRIPTION_TYPE_GOLDEN_VERSION:
		return GoldenVersion, nil
	case subscriptionpb.SubscriptionType_SUBSCRIPTION_TYPE_PLATINUM_VERSION:
		return PlatinumVersion, nil
	case subscriptionpb.SubscriptionType_SUBSCRIPTION_TYPE_DIAMOND_VERSION:
		return DiamondVersion, nil
	case subscriptionpb.SubscriptionType_SUBSCRIPTION_TYPE_GIFT_VERSION:
		return GiftVersion, nil
	}
	return "", fmt.Errorf("invalid proto subscription type %d", subscriptionType)
}
