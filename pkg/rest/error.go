package rest

import (
	"fmt"
)

// Error is the error object
type Error struct {
	Code          int    `json:"code"`           // 错误号 一体机要求code必须是小写
	Message       string `json:"message"`        // 错误消息
	Cause         string `json:"cause"`          // 原因
	InternalError error  `json:"internal_error"` // 内部错误
}

// NewError 构建一个新的 Error
func NewError(code int, msg, cause string) Error {
	return Error{
		Code:    code,
		Message: msg,
		Cause:   cause,
	}
}

// NewErrorWithError 构建一个新的 Error
func NewErrorWithError(code int, msg, cause string, err error) Error {
	return Error{
		Code:          code,
		Message:       msg,
		Cause:         cause,
		InternalError: err,
	}
}

// Error is for the error interface
func (e Error) Error() string {
	if e.Cause != "" {
		return fmt.Sprintf("[err:%d] %s (%s)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[err:%d] %s", e.Code, e.Message)
}
