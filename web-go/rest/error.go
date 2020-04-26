package rest

var codeToMsg = map[int]string{
	ErrOK: "OK",

	ErrRPCInternal: "RPC internal error",
}

const (
	// 错误码定义清单

	// ErrOK OK. Not used.
	ErrOK = 0

	// RPC 请求相关

	// ErrRPCInternal RPC内部错误
	ErrRPCInternal = 3000
)

// ErrorMsg 根据错误码获得标准错误消息内容
func ErrorMsg(code int) string {
	if msg, ok := codeToMsg[code]; ok {
		return msg
	}

	return ""
}
