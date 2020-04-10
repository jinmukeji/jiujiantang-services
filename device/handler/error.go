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
)

// NewError 构建一个新的 Error
func NewError(code int, err error) error {
	return fmt.Errorf("[errcode:%d] %s", code, err.Error())
}

// NewErrorCause 构建一个error
func NewErrorCause(code int, err error, cause string) error {
	return NewError(code, err)
}
