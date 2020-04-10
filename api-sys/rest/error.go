package rest

import (
	r "github.com/jinmukeji/jiujiantang-services/pkg/rest.v3"
	"github.com/kataras/iris/v12"
)

var codeToMsg = map[int]string{
	// 错误码定义清单
	ErrClientUnauthorized: "Unauthorized Client",

	ErrParsingRequestFailed:       "Parsing request body failed",
	ErrGetClientPreferencesFailed: "Failed to get client preferences",

	ErrRPCInternal:      "RPC internal error",
	ErrInvalidSecretKey: "SecretKey is invalid",
}

const (
	// 错误码定义清单

	// 授权、身份验证、权限等错误

	// ErrClientUnauthorized Client未授权
	ErrClientUnauthorized = 1000

	// ErrInvalidSecretKey secretkey错误
	ErrInvalidSecretKey = 1500

	// ErrParsingRequestFailed 数据错误
	ErrParsingRequestFailed = 2000
	// ErrGetClientPreferencesFailed 获取客户端对应的URL失败
	ErrGetClientPreferencesFailed = 2001

	// ErrRPCInternal RPC 请求相关
	ErrRPCInternal = 3000
)

// ErrorMsg 根据错误码获得标准错误消息内容
func ErrorMsg(code int) string {
	if msg, ok := codeToMsg[code]; ok {
		return msg
	}

	return ""
}

// WrapError 包装一个 Error
func wrapError(code int, cause string, err error) r.Error {
	return r.NewErrorWithError(code, ErrorMsg(code), cause, err)
}

// WriteError 向 Response 中写入 Error
func writeError(ctx iris.Context, err r.Error) {
	l := r.ContextLogger(ctx)
	if err.InternalError != nil {
		l.WithError(err.InternalError).Warn(err.Error())
	} else {
		l.Warn(err.Error())
	}
	r.WriteErrorJSON(ctx, err)
}
