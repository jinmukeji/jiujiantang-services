package handler

import (
	"fmt"
)

const (
	// ErrDatabase 数据库错误
	ErrDatabase = 10001
	// ErrClientUnauthorized 未授权
	ErrClientUnauthorized = 1000
	// ErrUserUnauthorized 未授权
	ErrUserUnauthorized = 1100
	// ErrInvalidUser 无效的用户
	ErrInvalidUser = 1200
	// ErrExpiredActivationCode 激活码过期
	ErrExpiredActivationCode = 10000
	// ErrInvalidActivationCode 激活码无效
	ErrInvalidActivationCode = 11000
	// ErrNotSoldActivationCode 激活码没有售出
	ErrNotSoldActivationCode = 12000
	// ErrActivationCodeWrongChecksum 校验错误
	ErrActivationCodeWrongChecksum = 13000
	// ErrActivatedActivationCode 激活码已经激活
	ErrActivatedActivationCode = 14000
	// ErrSubscriptionNotBelongToUser 订阅不属于用户
	ErrSubscriptionNotBelongToUser = 15000
	// ErrForbidToRemoveSubscriptionOwner 不能删除订阅的拥有者
	ErrForbidToRemoveSubscriptionOwner = 19000
	// ErrNoneExistSubscription 用户没有订阅
	ErrNoneExistSubscription = 20000
	// ErrUserNotShareSubscription 用户不在订阅的分享用户里
	ErrUserNotShareSubscription = 21000
)

// NewError 构建一个新的 Error
func NewError(code int, err error) error {
	return fmt.Errorf("[errcode:%d] %s", code, err.Error())
}

// NewErrorCause 构建一个error
func NewErrorCause(code int, err error, cause string) error {
	return NewError(code, err)
}
