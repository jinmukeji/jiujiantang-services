package rest

import (
	"regexp"
	"strconv"

	r "github.com/jinmukeji/jiujiantang-services/pkg/rest"
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

	ErrUsernamePasswordNotMatch:              "用户名或密码错误",
	ErrNullClientID:                          "空的客户端ID",
	ErrIncorrectClientID:                     "客户端ID不正确",
	ErrInvalidSecretKey:                      "secretkey错误",
	ErrInvalidUser:                           "无效的用户",
	ErrBuildJwtToken:                         "JWT Token 生成错误",
	ErrInvalidPassword:                       "Request 数据错误",
	ErrPhonePasswordNotMatch:                 "手机号或密码错误",
	ErrExistPassword:                         "密码已经设置",
	ErrIncorrectPassword:                     "密码不正确",
	ErrExistRegisteredPhone:                  "手机号已被注册",
	ErrInValidMVC:                            "验证码失效",
	ErrGetAccessTokenFailure:                 "登录已过期，请重新登录",
	ErrInvalidValidationValue:                "请选择一项填写",
	ErrInvalidSecureQuestionValidationMethod: "非法的安全问题验证方法",
	ErrWrongSecureQuestionCount:              "请选择3个密保问题",
	ErrSecureQuestionExist:                   "已设置过密保问题",
	ErrEmptySecureQuestion:                   "已设置过密保问题",
	ErrEmptyAnswer:                           "答案不能为空",
	ErrMismatchQuestion:                      "密保问题错误",
	ErrWrongFormatQuestion:                   "答案限制15字以内的中英文或数字",
	ErrRepeatedQuestion:                      "问题不能重复",
	ErrWrongFormatPhone:                      "手机号格式错误",
	ErrNonexistentUsername:                   "用户名不存在",
	ErrNoneExistentPhone:                     "手机号未注册",
	ErrInvalidEmailAddress:                   "邮箱格式错误",
	ErrSecureEmailExists:                     "已设置安全邮箱",
	ErrSecureEmailUsedByOthers:               "该邮箱已被其他账号绑定",
	ErrSecureEmailNotSet:                     "未设置安全邮箱",
	ErrSecureEmailAddressNotMatched:          "非当前绑定邮箱",
	ErrInvalidEmailNotificationAction:        "非法的邮件通知的操作",
	ErrUsedVcRecord:                          "验证码已失效",
	ErrExpiredVcRecord:                       "验证码已过期",
	ErrInvalidVcRecord:                       "验证码错误",
	ErrInvalidRequestCountIn1Minute:          "邮件一分钟内只能发送一次验证码",
	ErrExsitRegion:                           "区域已存在",
	ErrExsitSignInPhone:                      "手机号已被注册",
	ErrInvalidSigninPhone:                    "手机号无效",
	ErrInvalidVerificationNumber:             "验证码无效",
	ErrNotExistSigninPhone:                   "手机号不存在",
	ErrSamePhone:                             "新旧手机号相同",
	ErrWrongSendVia:                          "发送验证码的方式错误",
	ErrInvalidSendValue:                      "发送验证码的值式错误",
	ErrInvalidSigninEmail:                    "无效的安全邮箱",
	ErrInvalidValidationType:                 "非法的验证类型",
	ErrInvalidValidationMethod:               "获取方式非法",
	ErrNotExistNewSecureEmail:                "新安全邮箱不存在",
	ErrSameEmail:                             "新旧安全邮箱相同",
	ErrSameSecureQuestion:                    "问题未修改",
	ErrNotExistOldPassword:                   "您未设置过密码",
	ErrCurrentSecureQuestionsNotSet:          "未设置密保问题",
	ErrSecureEmailNotSetByAnyone:             "该邮箱未被绑定",
	ErrNationCode:                            "区号不正确",
	ErrWrongSmsNotificationType:              "短信类型不正确",
	ErrSignInPhoneNotBelongsToUser:           "手机号不属于用户",
	ErrUsernameNotSet:                        "未设置用户名",
	ErrNonexistentPassword:                   "未设置过密码",
	ErrIncorrectOldPassword:                  "旧密码不正确",
	ErrWrongFormatOfNickname:                 "昵称为1—15个字符，首位不能为特殊字符",
	ErrSensitiveWordsInNickname:              "昵称包含违禁词",
	ErrReservedWordsInNickname:               "昵称包含违禁词",
	ErrMaskWordsInNickname:                   "昵称包含违禁词",
	ErrEmptyNickname:                         "请设置昵称",
	ErrEmptyGender:                           "请选择性别",
	ErrEmptyBirthday:                         "请选择生日",
	ErrEmptyHeight:                           "请选择身高",
	ErrEmptyWeight:                           "请选择体重",
	ErrEmptyLanguage:                         "请选择语言",
	ErrEmptyRegion:                           "请选择区域",
	ErrWrongFormatOfPassword:                 "密码为8-20位，同时包含数字和字母",
	ErrEmptyPassword:                         "请设置密码",
	ErrNotEmailOfCurrentUser:                 "非当前绑定邮箱",
	ErrNoneExistUser:                         "不存在的User",
	ErrSamePassword:                          "新密码不能与原密码相同",
	ErrSendMoreSMSInOneMinute:                "短信一分钟内只能发送一次验证码（72000）",
	ErrSendSMS:                               "短信网关异常（74000）",
	ErrDeactivatedUser:                       "用户被禁用",
	ErrSendEmail:                             "邮件网关异常（75000）",
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

	// ErrPhonePasswordNotMatch 电话密码不匹配
	ErrPhonePasswordNotMatch = 1700
	// ErrExistPassword 密码已经存在
	ErrExistPassword = 1800
	// ErrIncorrectPassword 密码不正确
	ErrIncorrectPassword = 1900
	// ErrExistRegisteredPhone 注册手机号已经存在
	ErrExistRegisteredPhone = 2000
	// ErrInValidMVC 验证码失效
	ErrInValidMVC = 2100
	// ErrGetAccessTokenFailure 失败去得到Access-Token
	ErrGetAccessTokenFailure = 2200

	// Request 数据错误

	// ErrParsingRequestFailed 解析请求错误
	ErrParsingRequestFailed = 2300
	//  ErrValueRequired 请求值错误
	ErrValueRequired = 2400
	// ErrInvalidValue 无效的值
	ErrInvalidValue = 2500
	// ErrInvalidPassword 无效的密码
	ErrInvalidPassword = 2600

	// ErrDeactivatedUser 用户被禁用
	ErrDeactivatedUser = 2700
	// RPC 请求相关

	// ErrRPCInternal RPC内部错误
	ErrRPCInternal = 3000
	// ErrRPCTimeout RPC超时错误
	ErrRPCTimeout = 3001

	// ErrBuildJwtToken JWT Token 生成错误
	ErrBuildJwtToken = 4001
	// ErrClientInternal 客户端内部错误
	ErrClientInternal = 5000

	// ErrInvalidValidationValue 非法的验证方式的值
	ErrInvalidValidationValue = 10000
	// ErrInvalidSecureQuestionValidationMethod 非法的安全问题验证方法
	ErrInvalidSecureQuestionValidationMethod = 11000
	// ErrWrongSecureQuestionCount 安全问题数量不匹配
	ErrWrongSecureQuestionCount = 12000
	// ErrSecureQuestionExist 已经设置过密保
	ErrSecureQuestionExist = 13000
	// ErrEmptySecureQuestion 密保问题为空
	ErrEmptySecureQuestion = 14000
	// ErrEmptyAnswer  答案为空
	ErrEmptyAnswer = 15000
	// ErrSamePassword  新密码与旧密码相同
	ErrSamePassword = 16000
	// ErrMismatchQuestion 传入的密保问题错误
	ErrMismatchQuestion = 17000
	// ErrWrongFormatQuestion 传入的密保问题格式错误
	ErrWrongFormatQuestion = 18000
	// ErrRepeatedQuestion 问题重复
	ErrRepeatedQuestion = 19000
	// ErrWrongFormatPhone 手机号格式错误
	ErrWrongFormatPhone = 20000
	// ErrNonexistentUsername 用户名不存在
	ErrNonexistentUsername = 21000
	// ErrNoneExistentPhone 手机号不存在
	ErrNoneExistentPhone = 22000
	// ErrInvalidEmailAddress 邮箱格式错误
	ErrInvalidEmailAddress = 23000
	// ErrSecureEmailExists 用户已经设置了自己的邮箱
	ErrSecureEmailExists = 24000
	// ErrSecureEmailUsedByOthers 邮箱已经被其他人设置
	ErrSecureEmailUsedByOthers = 25000
	// ErrSecureEmailNotSet 没有设置邮箱
	ErrSecureEmailNotSet = 26000
	// ErrSecureEmailAddressNotMatched 与原邮箱不匹配
	ErrSecureEmailAddressNotMatched = 27000
	// ErrInvalidEmailNotificationAction 非法的邮件通知的操作
	ErrInvalidEmailNotificationAction = 28000
	// ErrUsedVcRecord 该记录已经被使用过
	ErrUsedVcRecord = 29000
	// ErrExpiredVcRecord 该记录已经过期
	ErrExpiredVcRecord = 30000
	// ErrInvalidVcRecord 验证码错误
	ErrInvalidVcRecord = 31000
	// ErrInvalidRequestCountIn1Minute 1分钟请求次数非法
	ErrInvalidRequestCountIn1Minute = 32000
	// ErrExsitRegion 区域已经存在
	ErrExsitRegion = 33000
	// ErrExsitSignInPhone 登录电话已经存在
	ErrExsitSignInPhone = 34000
	// ErrInvalidSigninPhone 无效的登录电话
	ErrInvalidSigninPhone = 35000
	// ErrInvalidVerificationNumber 无效的VerificationNumber
	ErrInvalidVerificationNumber = 36000
	// ErrNotExistSigninPhone 登录手机号不存在
	ErrNotExistSigninPhone = 37000
	// ErrSamePhone 新旧手机号一样
	ErrSamePhone = 38000
	// ErrWrongSendVia 发送验证码的方式错误
	ErrWrongSendVia = 39000
	// ErrInvalidSendValue 发送验证码的方式的值错误
	ErrInvalidSendValue = 40000
	// ErrInvalidSigninEmail 无效的安全邮箱
	ErrInvalidSigninEmail = 41000
	// ErrInvalidValidationType 验证邮箱验证码时非法的验证类型
	ErrInvalidValidationType = 42000
	// ErrInvalidValidationMethod 获取方式非法
	ErrInvalidValidationMethod = 43000
	// ErrNotExistNewSecureEmail 新安全邮箱不存在
	ErrNotExistNewSecureEmail = 45000
	// ErrSameEmail 新旧安全邮箱相同
	ErrSameEmail = 46000
	// ErrSameSecureQuestion 新旧密保问题相同
	ErrSameSecureQuestion = 47000
	// ErrNotExistOldPassword 旧密码不存在
	ErrNotExistOldPassword = 48000
	// ErrCurrentSecureQuestionsNotSet 当前用户的密保问题未设置
	ErrCurrentSecureQuestionsNotSet = 49000
	// ErrSecureEmailNotSetByAnyone 安全邮箱没有被任何人设置
	ErrSecureEmailNotSetByAnyone = 50000
	// ErrNationCode 区号不正确
	ErrNationCode = 51000
	// ErrWrongSmsNotificationType 短信类型不正确
	ErrWrongSmsNotificationType = 52000
	// ErrSignInPhoneNotBelongsToUser 手机号不属于用户
	ErrSignInPhoneNotBelongsToUser = 53000
	// ErrUsernameNotSet 用户名未设置
	ErrUsernameNotSet = 54000
	// ErrNonexistentPassword 密码不存在
	ErrNonexistentPassword = 55000
	// ErrIncorrectOldPassword 旧密码错误
	ErrIncorrectOldPassword = 56000
	// ErrWrongFormatOfNickname 昵称格式错误
	ErrWrongFormatOfNickname = 57000
	// ErrSensitiveWordsInNickname 昵称包含敏感词
	ErrSensitiveWordsInNickname = 58000
	// ErrReservedWordsInNickname 昵称包含保留词
	ErrReservedWordsInNickname = 59000
	// ErrMaskWordsInNickname 昵称包含屏蔽词
	ErrMaskWordsInNickname = 60000
	// ErrEmptyNickname 昵称为空
	ErrEmptyNickname = 61000
	// ErrEmptyGender 性别为空
	ErrEmptyGender = 62000
	// ErrEmptyBirthday 生日为空
	ErrEmptyBirthday = 63000
	// ErrEmptyHeight 身高为空
	ErrEmptyHeight = 64000
	// ErrEmptyWeight 体重为空
	ErrEmptyWeight = 65000
	// ErrEmptyLanguage 语言为空
	ErrEmptyLanguage = 66000
	// ErrEmptyRegion 区域为空
	ErrEmptyRegion = 67000
	// ErrWrongFormatOfPassword 密码格式错误
	ErrWrongFormatOfPassword = 68000
	// ErrEmptyPassword 密码为空
	ErrEmptyPassword = 69000
	// ErrNotEmailOfCurrentUser 非当前绑定邮箱
	ErrNotEmailOfCurrentUser = 70000
	// ErrNoneExistUser 不存在的User
	ErrNoneExistUser = 71000
	// ErrSendMoreSMSInOneMinute 一分钟内发送多条短信
	ErrSendMoreSMSInOneMinute = 72000
	// ErrSendSMS 短信发送异常
	ErrSendSMS = 74000
	// ErrSendEmail 邮件发送异常
	ErrSendEmail = 75000
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
