class const:
    # 错误码定义清单
    ErrOK = 0  # OK.Not used.
    ErrUnknown = 1  # Unknown error

    # 授权、身份验证、权限等错误
    ErrClientUnauthorized = 1000  # Client未授权
    ErrUserUnauthorized = 1100  # User未授权
    ErrUsernamePasswordNotMatch = 1200  # ErrUsernamePasswordNotMatch用户名密码错误
    ErrNullClientID = 1300  # ErrNullClientID空的客户端ID
    ErrIncorrectClientID = 1400  # ErrIncorrectClientID 客户端ID不正确
    ErrInvalidSecretKey = 1500  # ErrInvalidSecretKey secretkey错误
    ErrInvalidUser = 1600  # ErrInvalidUser  无效的用户

    # Request数据错误

    ErrParsingRequestFailed = 2000
    ErrValueRequired = 2001
    ErrInvalidValue = 2002

    # RPC请求相关

    ErrRPCInternal = 3000
    ErrRPCTimeout = 3001
    ErrClientInternal = 5000
    # ErrDatabase数据库错误
    ErrDatabase = 10001
    # ErrExpiredActivationCode激活码过期
    ErrExpiredActivationCode = 10000
    # ErrInvalidActivationCode 激活码无效
    ErrInvalidActivationCode = 11000
    # ErrNotSoldActivationCode激活码没有售出
    ErrNotSoldActivationCode = 12000
    # ErrActivationCodeWrongChecksum校验错误
    ErrActivationCodeWrongChecksum = 13000
    # ErrActivatedActivationCode激活码已经激活
    ErrActivatedActivationCode = 14000
    # ErrSubscriptionNotBelongToUser订阅不属于用户
    ErrSubscriptionNotBelongToUser = 15000
    # ErrInactivatedSubscription未激活的订阅
    ErrInactivatedSubscription = 1800
    # ErrExceedSubscriptionUserQuotaLimit订阅下的用户数量已经达到规定
    ErrExceedSubscriptionUserQuotaLimit = 1900
    # ErrForbidToRemoveSubscriptionOwner不能删除订阅的拥有者
    ErrForbidToRemoveSubscriptionOwner = 19000
    # ErrSubscriptionExpired订阅过期
    ErrSubscriptionExpired = 76000
    # ErrBlockedMacMac不可用
    ErrBlockedMac = 77000
    # ErrBlockedIPIP不可用
    ErrBlockedIP = 78000
    # ErrRecordNotBelongToUser记录不属于用户
    ErrRecordNotBelongToUser = 80000
    # ErrInvalidOrganizationCount用户没有组织或者有或个组织
    ErrInvalidOrganizationCount = 81000
