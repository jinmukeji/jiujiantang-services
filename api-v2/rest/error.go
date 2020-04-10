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

	ErrClientUnauthorized: "Client 未授权",
	ErrUserUnauthorized:   "User 未授权",

	ErrParsingRequestFailed: "Request 数据错误",
	ErrValueRequired:        "Request 数据错误",
	ErrInvalidValue:         "Request 数据错误",

	ErrRPCInternal: "RPC请求错误",
	ErrRPCTimeout:  "RPC请求错误",

	ErrClientInternal: "JWT Token 生成错误",

	ErrUsernamePasswordNotMatch: "用户名或密码错误",
	ErrNullClientID:             "空的客户端ID",
	ErrIncorrectClientID:        "客户端ID不正确",
	ErrInvalidSecretKey:         "secretkey错误",
	ErrInvalidUser:              "无效的用户",
	ErrBuildJwtToken:            "JWT Token 生成错误",

	ErrExpiredActivationCode:            "激活码已过期",
	ErrInvalidActivationCode:            "激活码已失效",
	ErrNotSoldActivationCode:            "激活码不可用",
	ErrActivationCodeWrongChecksum:      "激活码错误",
	ErrActivatedActivationCode:          "激活码已失效",
	ErrSubscriptionNotBelongToUser:      "订阅不属于用户",
	ErrInactivatedSubscription:          "订阅未激活",
	ErrExceedSubscriptionUserQuotaLimit: "您的用户数已达上限",
	ErrForbidToRemoveSubscriptionOwner:  "默认用户不能删除",
	ErrDeactivatedUser:                  "用户被禁用",
	ErrDeniedToAccessAPI:                "权限被拒",
	ErrSubscriptionExpired:              "订阅已过期",
	ErrBlockedMac:                       "Mac不可用",
	ErrBlockedIP:                        "IP不可用",
	ErrRecordNotBelongToUser:            "记录不属于用户",
	ErrInvalidAnalysisStatus:            "无效的分析状态",
	ErrInvalidAge:                       "无效的年龄",
	ErrInvalidWeight:                    "无效的体重",
	ErrInvalidHeight:                    "无效的身高",
	ErrMultiOwnersOfOrganization:        "组织有多个拥有者",
	ErrNonexistentOwnerOfOrganization:   "组织没有拥有者",
	ErrAEError:                          "分析错误",
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

	// ErrDeactivatedUser 用户被禁用
	ErrDeactivatedUser = 2700

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
	// ErrInactivatedSubscription 未激活的订阅
	ErrInactivatedSubscription = 1800
	// ErrExceedSubscriptionUserQuotaLimit 订阅下的用户数量已经达到规定
	ErrExceedSubscriptionUserQuotaLimit = 1900
	// ErrForbidToRemoveSubscriptionOwner 不能删除订阅的拥有者
	ErrForbidToRemoveSubscriptionOwner = 19000
	// ErrSubscriptionExpired 订阅过期
	ErrSubscriptionExpired = 76000
	// ErrBlockedMac Mac不可用
	ErrBlockedMac = 77000
	// ErrBlockedIP IP不可用
	ErrBlockedIP = 78000
	// ErrRecordNotBelongToUser 记录不属于用户
	ErrRecordNotBelongToUser = 80000
	// ErrInvalidOrganizationCount 用户没有组织或者有或个组织
	ErrInvalidOrganizationCount = 81000
	// ErrInvalidAnalysisStatus 无效的分析状态
	ErrInvalidAnalysisStatus = 82000
	// Request 数据错误

	// ErrParsingRequestFailed 解析请求失败
	ErrParsingRequestFailed = 2000
	// ErrValueRequired 请求值错误
	ErrValueRequired = 2001
	// ErrInvalidValue 无效的值
	ErrInvalidValue = 2002

	// RPC 请求相关

	// ErrRPCInternal RPC内部错误
	ErrRPCInternal = 3000
	// ErrRPCTimeout RPCTimeout
	ErrRPCTimeout = 3001

	//ErrBuildJwtToken JWT Token 生成错误
	ErrBuildJwtToken = 4001

	// ErrClientInternal 客户端内部错误
	ErrClientInternal = 5000

	// ErrDeniedToAccessAPI 权限被拒
	ErrDeniedToAccessAPI = 6000
	// ErrInvalidAge 无效的年龄
	ErrInvalidAge = 83000
	// ErrInvalidWeight 无效的体重
	ErrInvalidWeight = 84000
	// ErrInvalidHeight 无效的身高
	ErrInvalidHeight = 85000
	// ErrMultiOwnersOfOrganization 组织有多个拥有者
	ErrMultiOwnersOfOrganization = 86000
	// ErrNonexistentOwnerOfOrganization 组织没有拥有者
	ErrNonexistentOwnerOfOrganization = 87000
	// ErrAEError 分析错误
	ErrAEError = 88000
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

func writeRPCInternalError(ctx iris.Context, err error, shouldBeArrayData bool) {
	if code, ok := rpcErrorCode(err); ok {
		if _, ok := codeToMsg[code]; ok {
			writeError(ctx, wrapError(code, "", err), shouldBeArrayData)
			return
		}
	}

	writeError(ctx, wrapError(ErrRPCInternal, "", err), shouldBeArrayData)
}
