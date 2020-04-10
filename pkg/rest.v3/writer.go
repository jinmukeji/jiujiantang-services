package rest

import (
	"fmt"

	"github.com/kataras/iris/v12"
)

const (
	errCodeHTTPError = -1
)

// SetContextMessage 设置context消息
func SetContextMessage(ctx iris.Context, msg string) {
	ctx.Values().Set("message", msg)
}

// GetContextMessage 得到context消息
func GetContextMessage(ctx iris.Context) string {
	return ctx.Values().GetString("message")
}

// NotFound 输出 StatusNotFound 的响应
func NotFound(ctx iris.Context) {
	err := NewError(errCodeHTTPError, "Not Found", "")
	WriteErrorJSON(ctx, err)
}

// InternalServerError 输出内部服务器错误的响应
func InternalServerError(ctx iris.Context) {
	err := NewError(errCodeHTTPError, "Internal Server Error", "")
	WriteErrorJSON(ctx, err)
}

// WriteErrorJSON 输出错误响应
func WriteErrorJSON(ctx iris.Context, err Error) {
	cid := GetCidFromContext(ctx)
	ret := iris.Map{
		"cid": cid,
		"ok":  false,
		"error": iris.Map{
			"code": err.Code,
			"msg":  err.Message,
		},
	}

	// nolint: errcheck, gas
	ctx.JSON(ret)

	SetContextMessage(ctx, fmt.Sprintf("API ERR: %s", err.Error()))
}

// WriteOkJSON 输出正确响应
func WriteOkJSON(ctx iris.Context, data interface{}) {
	cid := GetCidFromContext(ctx)
	ret := iris.Map{
		"cid":  cid,
		"ok":   true,
		"data": data,
	}

	// nolint: errcheck, gas
	ctx.JSON(ret)

	SetContextMessage(ctx, "API OK")
}
