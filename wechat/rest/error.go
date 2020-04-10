package rest

import (

    "github.com/kataras/iris/v12"
    r "github.com/jinmukeji/jiujiantang-services/pkg/rest"
)

var codeToMsg = map[int]string{
	ErrOK:      "OK",
	ErrUnknown: "Unknown error",

	ErrClientUnauthorized:  "Unauthorized Client",
	ErrUserUnauthorized:    "Unauthorized User",
	ErrSessionUnauthorized: "Unauthorized Session",

	ErrParsingRequestFailed: "Parsing request body failed",
	ErrValueRequired:        "Value is required",
	ErrInvalidValue:         "Invalid value",

	ErrRPCInternal: "RPC internal error",
	ErrRPCTimeout:  "RPC request is timeout",

	ErrWxOAuth:       "WX OAuth error",
	ErrWxJsSdkTicket: "WX JS-SDK Ticket Server error",

	ErrClientInternal: "Client internal error",
}

const (
	// 错误码定义清单

	// ErrOK OK. Not used.
	ErrOK      = 0 
	// ErrUnknown Unknown error
	ErrUnknown = 1 

	// 授权、身份验证、权限等错误

	// ErrClientUnauthorized Client 未授权
	ErrClientUnauthorized  = 1000 
	// ErrUserUnauthorized User 未授权
	ErrUserUnauthorized    = 1100 
	// ErrSessionUnauthorized Session 未授权
	ErrSessionUnauthorized = 1200 

	// Request 数据错误

	// ErrParsingRequestFailed 解析请求失败
	ErrParsingRequestFailed = 2000
	// ErrValueRequired 请求值错误
	ErrValueRequired        = 2001
	// ErrInvalidValue 无效的值
	ErrInvalidValue         = 2002

	// RPC 请求相关
	// ErrRPCInternal RPC内部错误
	ErrRPCInternal = 3000
	// ErrRPCTimeout RPC超时
	ErrRPCTimeout  = 3001

	// 微信请求相关
	
	// ErrWxOAuth 微信OAuth错误
	ErrWxOAuth       = 60000
	// ErrWxJsSdkTicket WX JS-SDK Ticket Server错误
	ErrWxJsSdkTicket = 60001
	// ErrClientInternal 客户端内部错误
	ErrClientInternal = 5000
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
func writeError(ctx iris.Context, err r.Error, shouldBeArrayData bool) {
    l := r.ContextLogger(ctx)
    if err.InternalError != nil {
        l.WithError(err.InternalError).Warn(err.Error())
	} else {
		l.Warn(err.Error())
	}
    r.WriteErrorJSON(ctx, err, shouldBeArrayData)
}


