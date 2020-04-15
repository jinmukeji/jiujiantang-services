package handler

import (
	"fmt"
)

const (
	// ErrDatabase 数据库错误
	ErrDatabase = 10001
	// ErrInvalidDevice 无效的设备
	ErrInvalidDevice = 11001
	// ErrBindedDevice 设备已经被关联
	ErrBindedDevice = 11002
	// ErrDeviceValidationFailure 验证设备失败
	ErrDeviceValidationFailure = 11003
	// ErrDeviceNotBelongToOrganization 设备不属于组织
	ErrDeviceNotBelongToOrganization = 11004
	// ErrDeviceNotFoundInOrganization 没有查到设备关联信息
	ErrDeviceNotFoundInOrganization = 11005
	// ErrDeviceNotBelongToClient 设备不属于当前Client
	ErrDeviceNotBelongToClient = 11006
	// ErrNullMachineUUID 机器UUID是NULL
	ErrNullMachineUUID = 11007

	// ErrInvalidAccessToken access-token错误
	ErrInvalidAccessToken = 12001
	// ErrIncorrectRegisterType 注册类型不正确
	ErrIncorrectRegisterType = 12002

	// ErrAddUserFailure 自定义用户加入组织的错误
	ErrAddUserFailure = 12004
	// ErrGetAccessTokenFailure 失败去得到Access-Token
	ErrGetAccessTokenFailure = 12005

	// ErrCreateTokenFailure 创建Token失败
	ErrCreateTokenFailure = 12010
	// ErrGetUserFailure 得到User失败
	ErrGetUserFailure = 12011
	// ErrNotFoundUser 没有发现用户
	ErrNotFoundUser = 12012
	// ErrCannotDeleteOwner 不能删除拥有者
	ErrCannotDeleteOwner = 12013
	// ErrInactivatedUser 用户没有激活
	ErrInactivatedUser = 12014
	// ErrInvalidPassword 无效的密码
	ErrInvalidPassword = 12015

	// ErrInvalidOrganization 无效的组织
	ErrInvalidOrganization = 13001
	// ErrGetOrganizationUsersFailure 无法得到组织的用户
	ErrGetOrganizationUsersFailure = 13002
	// ErrOrganizationQueryUserExceedsLimit 查询用户超出限制
	ErrOrganizationQueryUserExceedsLimit = 13003
	// ErrUserNotOrganizationOwner 用户不能组织拥有者
	ErrUserNotOrganizationOwner = 13004
	// ErrOrganizationCountExceedsMaxLimits 组织数量超出限制
	ErrOrganizationCountExceedsMaxLimits = 13005
	// ErrNotSupportDeleteOrganization 删除组织暂不支持
	ErrNotSupportDeleteOrganization = 13006

	// ErrProtoConversionFailure proto转化错误
	ErrProtoConversionFailure = 14001
	// ErrJSONUnmarshalFailure json解析失败
	ErrJSONUnmarshalFailure = 14002

	// ErrGetMeasureResultFailure 得到测量结果失败
	ErrGetMeasureResultFailure = 15001
	// ErrGetRecordFailure 得到记录失败
	ErrGetRecordFailure = 15002
	// ErrValidateSearchHistoryRequestFailure 验证历史记录请求失败
	ErrValidateSearchHistoryRequestFailure = 15003
	// ErrMeasurementDataLengthNotMatch  测量数据长度不匹配
	ErrMeasurementDataLengthNotMatch = 15004
	// ErrBuildAlgorithmRequestFailure 建立算法服务器请求失败
	ErrBuildAlgorithmRequestFailure = 15005
	// ErrInvokeAlgorithmServerFailure 连接算法服务器失败
	ErrInvokeAlgorithmServerFailure = 15006
	// ErrUploadWavedataToAWSFailure 上传波形数据到AWS失败
	ErrUploadWavedataToAWSFailure = 15007
	// ErrSetWavedataFailure 设置波形数据失败
	ErrSetWavedataFailure = 15008
	// ErrNoPermissionSubmitRemark 没有权限去提交备注
	ErrNoPermissionSubmitRemark = 15009
	// ErrRunAnalysisEngineFailure 运行引擎失败
	ErrRunAnalysisEngineFailure = 15010
	// ErrNoPermissionGetShareToken 没有权限去得到分析token
	ErrNoPermissionGetShareToken = 15011
	// ErrGetTransactionNumberFailure 得到流水失败
	ErrGetTransactionNumberFailure = 15012
	// ErrValidateSubmitFeedbackRequestFailure 验证提交意见反馈请求失败
	ErrValidateSubmitFeedbackRequestFailure = 15013
	// ErrNULLCommentOrContactWay 评论或者联系方式为NULL
	ErrNULLCommentOrContactWay = 15014

	// ErrGetWxJsSdkConfigFaliure 得到微信JSSDK配置失败
	ErrGetWxJsSdkConfigFaliure = 16001
	// ErrGetTempQrCodeURLFaliure 得到临时二维码URL失败
	ErrGetTempQrCodeURLFaliure = 16002
	// ErrSendTemplateMessageFailure 发送微信模版失败
	ErrSendTemplateMessageFailure = 16003

	// ErrLocalNotification 本地通知错误
	ErrLocalNotification = 17001
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
	// ErrInactivatedSubscription 未激活的订阅
	ErrInactivatedSubscription = 1800
	// ErrExceedSubscriptionUserQuotaLimit 订阅下的用户数量已经达到规定
	ErrExceedSubscriptionUserQuotaLimit = 1900
	// ErrSubscriptionNotBelongToUser 订阅不属于用户
	ErrSubscriptionNotBelongToUser = 15000
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
	// ErrUserNotInOrganization 用户不属于组织
	ErrUserNotInOrganization = 81000
	// ErrInvalidAnalysisStatus 无效的分析状态
	ErrInvalidAnalysisStatus = 82000
	// ErrMultiOwnersOfOrganization 组织有多个拥有者
	ErrMultiOwnersOfOrganization = 86000
	// ErrNonexistentOwnerOfOrganization 组织没有拥有者
	ErrNonexistentOwnerOfOrganization = 87000
	// ErrInvalidFinger 非法的手指
	ErrInvalidFinger = 88000
	// ErrCreateRecordFailure 创建记录失败
	ErrCreateRecordFailure = 89000
	// ErrInvalidGender 性别非法
	ErrInvalidGender = 90000
	// ErrGenRandomString 生成随机字符串失败
	ErrGenRandomString = 91000
	// ErrInvalidPayload 解码失败，即 payload 非法
	ErrInvalidPayload = 9200
)

// NewError 构建一个新的 Error
func NewError(code int, err error) error {
	return fmt.Errorf("[errcode:%d] %s", code, err.Error())
}

// NewErrorCause 构建一个新的 Error
func NewErrorCause(code int, err error, cause string) error {
	return NewError(code, err)
}
