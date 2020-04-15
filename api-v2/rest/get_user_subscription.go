package rest

import (
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/subscription/v1"
	"github.com/kataras/iris/v12"
)

// SubscriptionStateDescription 用户的订阅情况
type SubscriptionStateDescription int

const (
	// Unsubscribed 未订阅
	Unsubscribed SubscriptionStateDescription = iota
	// NotExpired 订阅未过期
	NotExpired
	// Expired 订阅已过期
	Expired
)

// UsingSubscriptions 用户正在使用的订阅的信息
type UsingSubscriptions struct {
	Subscription AllSubscriptionStatus `json:"subscriptions"`
}

// AllSubscriptionStatus 所有相关的订阅状态
type AllSubscriptionStatus struct {
	SubscriptionStatus SubscriptionStateDescription `json:"user_subscription_status"`
	TotalUserCount     int32                        `json:"total_user_count"`
	MaxUserCount       int32                        `json:"max_user_count"`
	ExpiredAt          time.Time                    `json:"expired_at"`
	SubscriptionType   string                       `json:"subscription_type"`
}

func (h *v2Handler) GetUserSubscriptions(ctx iris.Context) {
	userID, err := ctx.Params().GetInt("user_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	respData := new(UsingSubscriptions)
	reqGetUserSubscriptions := new(proto.GetUserSubscriptionsRequest)
	reqGetUserSubscriptions.UserId = int32(userID)
	respGetUserSubscriptions, errGetUserSubscriptions := h.rpcSubscriptionManagerSvc.GetUserSubscriptions(newRPCContext(ctx), reqGetUserSubscriptions)
	if errGetUserSubscriptions != nil {
		code, _ := strconv.Atoi(regErrorCode.FindStringSubmatch(errGetUserSubscriptions.Error())[1])
		if code == 20000 {
			respData.Subscription.SubscriptionStatus = Unsubscribed
			rest.WriteOkJSON(ctx, respData)
			return
		}
		writeRPCInternalError(ctx, errGetUserSubscriptions, false)
		return
	}
	selectedSubscription := new(proto.Subscription)
	for _, item := range respGetUserSubscriptions.Subscriptions {
		if item.IsSelected {
			selectedSubscription = item
		}
	}
	expiredAt, err := ptypes.Timestamp(selectedSubscription.ExpiredTime)
	if err != nil {
		writeRPCInternalError(ctx, err, false)
	}
	if expiredAt.Before(time.Now()) {
		respData.Subscription.SubscriptionStatus = Expired
	}
	if expiredAt.After(time.Now()) {
		respData.Subscription.SubscriptionStatus = NotExpired
	}

	reqGetUserCountOfSubscription := new(proto.GetUserCountOfSubscriptionRequest)
	reqGetUserCountOfSubscription.SubscriptionId = selectedSubscription.SubscriptionId
	respGetUserCountOfSubscription, errGetUserCountOfSubscription := h.rpcSubscriptionManagerSvc.GetUserCountOfSubscription(newRPCContext(ctx), reqGetUserCountOfSubscription)
	if errGetUserCountOfSubscription != nil {
		writeRPCInternalError(ctx, errGetUserCountOfSubscription, false)
	}
	respData.Subscription.TotalUserCount = respGetUserCountOfSubscription.UserCount
	respData.Subscription.MaxUserCount = selectedSubscription.MaxUserLimits

	respData.Subscription.ExpiredAt = expiredAt
	stringSubscription, errmapProtoSubscriptionTypeToRest := mapProtoSubscriptionTypeToRest(selectedSubscription.SubscriptionType)
	if errmapProtoSubscriptionTypeToRest != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", errmapProtoSubscriptionTypeToRest), false)
		return
	}

	respData.Subscription.SubscriptionType = stringSubscription

	rest.WriteOkJSON(ctx, respData)

}
