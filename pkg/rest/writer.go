package rest

import (
	"fmt"

	"github.com/kataras/iris/v12"
)

const (
	// HTTPHeaderResponseBehavior Response 行为
	HTTPHeaderResponseBehavior = "X-Response-Behavior"
	// responseBehaviorV3 Response版本v3
	responseBehaviorV3 = "V3"
)

const (
	// errCodeHTTPError code http 错误
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
	WriteErrorJSON(ctx, err, false)
}

// InternalServerError 输出内部服务器错误的响应
func InternalServerError(ctx iris.Context) {
	err := NewError(errCodeHTTPError, "Internal Server Error", "")
	WriteErrorJSON(ctx, err, false)
}

var (
	// JSON: {}
	emptyData = iris.Map{}

	// JSON: []
	emptyArrayData = make([]struct{}, 0)

	// JSON: {"code": 0, "msg": ""}
	emptyErr = iris.Map{
		"code": 0,
		"msg":  "",
	}
)

func getEmptyData(shouldBeArrayData bool) interface{} {
	if shouldBeArrayData {
		return emptyArrayData
	}
	return emptyData
}

func dataOrEmpty(data interface{}) interface{} {
	if data == nil {
		return emptyData
	}
	return data
}

// WriteErrorJSON 输出错误响应
func WriteErrorJSON(ctx iris.Context, err Error, shouldBeArrayData bool) {
	cid := GetCidFromContext(ctx)
	ret := iris.Map{
		"cid": cid,
		"ok":  false,
		"error": iris.Map{
			"code": err.Code,
			"msg":  err.Message,
		},
	}

	if ctx.GetHeader(HTTPHeaderResponseBehavior) != responseBehaviorV3 {
		ret["data"] = getEmptyData(shouldBeArrayData)
	}

	// nolint: errcheck, gas
	ctx.JSON(ret)

	SetContextMessage(ctx, fmt.Sprintf("API ERR: %s", err.Error()))
}

// WriteOkJSON 输出正确响应
func WriteOkJSON(ctx iris.Context, data interface{}) {
	cid := GetCidFromContext(ctx)
	ret := iris.Map{
		"cid": cid,
		"ok":  true,
	}

	if ctx.GetHeader(HTTPHeaderResponseBehavior) == responseBehaviorV3 {
		ret["data"] = data
	} else {
		ret["data"] = dataOrEmpty(data)
		ret["error"] = emptyErr
	}

	// nolint: errcheck, gas
	ctx.JSON(ret)

	SetContextMessage(ctx, "API OK")
}
