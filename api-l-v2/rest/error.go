package rest

import (
	"regexp"
	"strconv"

	r "github.com/jinmukeji/gf-api2/pkg/rest"
	"github.com/kataras/iris/v12"
)

var codeToMsg = map[int]string{
	ErrOK:      "OK",
	ErrUnknown: "Unknown error",

	ErrClientUnauthorized: "Unauthorized Client",
	ErrUserUnauthorized:   "Unauthorized User",

	ErrParsingRequestFailed: "Parsing request body failed",
	ErrValueRequired:        "Value is required",
	ErrInvalidValue:         "Invalid value",

	ErrRPCInternal: "RPC internal error",
	ErrRPCTimeout:  "RPC request is timeout",

	ErrClientInternal: "Client internal error",

	ErrUsernamePasswordNotMatch: "Username and Password doesn't match",
	ErrNullClientID:             "ClientID is null",
	ErrIncorrectClientID:        "ClientID is incorrect",
	ErrInvalidSecretKey:         "SecretKey is invalid",
	ErrInvalidUser:              "User is invalid",
	ErrBuildJwtToken:            "Failed to build JWT token",
}

const (
	// 错误码定义清单

	// ErrOK OK. Not used.
	ErrOK = 0
	// ErrUnknown Unknown error
	ErrUnknown = 1

	// 授权、身份验证、权限等错误

	// ErrClientUnauthorized Client 未授权
	ErrClientUnauthorized = 1000
	// ErrUserUnauthorized User 未授权
	ErrUserUnauthorized = 1100
	// ErrUsernamePasswordNotMatch 用户名密码错误
	ErrUsernamePasswordNotMatch = 1200
	// ErrNullClientID 空的客户端ID
	ErrNullClientID = 1300
	// ErrIncorrectClientID 客户端ID不正确
	ErrIncorrectClientID = 1400
	// ErrInvalidSecretKey secretkey错误
	ErrInvalidSecretKey = 1500
	// ErrInvalidUser 无效的用户
	ErrInvalidUser = 1600

	// Request 数据错误

	// ErrParsingRequestFailed 解析请求失败
	ErrParsingRequestFailed = 2000
	// ErrValueRequired 请求值错误
	ErrValueRequired = 2001
	// ErrInvalidValue 无效的值
	ErrInvalidValue = 2002

	// RPC 请求相关

	// ErrRPCInternal PRC内部错误
	ErrRPCInternal = 3000
	// ErrRPCTimeout PRC超时
	ErrRPCTimeout = 3001

	// ErrBuildJwtToken JWT Token 生成错误
	ErrBuildJwtToken = 4001
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

var (
	regErrorCode = regexp.MustCompile(`\[errcode:(\d+)\]`)
)

func rpcErrorCode(rpcErr error) (int, bool) {
	regSubmatches := regErrorCode.FindStringSubmatch(rpcErr.Error())
	if len(regSubmatches) >= 2 {
		if rpcCode, err := strconv.Atoi(regSubmatches[1]); err == nil {
			return rpcCode, true
		}
	}
	return 0, false
}

func writeRpcInternalError(ctx iris.Context, err error, shouldBeArrayData bool) {
	if code, ok := rpcErrorCode(err); ok {
		if _, ok := codeToMsg[code]; ok {
			writeError(ctx, wrapError(code, "", err), shouldBeArrayData)
			return
		}
	}

	writeError(ctx, wrapError(ErrRPCInternal, "", err), shouldBeArrayData)
}
