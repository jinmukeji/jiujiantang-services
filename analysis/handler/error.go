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
	// ErrInvalidAnalysisStatus 无效的分析状态
	ErrInvalidAnalysisStatus = 82000
	// ErrInvalidAge 无效的年龄
	ErrInvalidAge = 83000
	// ErrInvalidWeight 无效的体重
	ErrInvalidWeight = 84000
	// ErrInvalidHeight 无效的身高
	ErrInvalidHeight = 85000
	// ErrInvalidGender 无效的性别
	ErrInvalidGender = 86000
	// ErrAEError AE错误
	ErrAEError = 88000
	// ErrClientID 错误的客户端 ID
	ErrClientID = 89000
	// ErrInvalidLanguage 错误的语言
	ErrInvalidLanguage = 90000
	// ErrLoadConfig 加载配置失败
	ErrLoadConfig = 91000
	// ErrInvalidFinger 无效的手指
	ErrInvalidFinger = 92000
	// ErrBuildReturnModule 构建返回的模块失败
	ErrBuildReturnModule = 93000
)

// NewError 构建一个新的 Error
func NewError(code int, err error) error {
	return fmt.Errorf("[errcode:%d] %s", code, err.Error())
}

// NewErrorCause 构建一个error
func NewErrorCause(code int, err error, cause string) error {
	return NewError(code, err)
}
