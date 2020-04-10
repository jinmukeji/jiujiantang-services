class const:
    # 错误码定义清单
    ErrOK = 0  # OK.Not used.
    ErrUnknown = 1  # Unknown error

    # 授权、身份验证、权限等错误

    ErrClientUnauthorized = 1000  # Client未授权

    ErrUserUnauthorized = 1100  # User未授权

    ErrUsernamePasswordNotMatch = 1200  # ErrUsernamePasswordNotMatch用户名密码错误

    ErrNullClientID = 1300  # ErrNullClientID空的客户端ID

    ErrIncorrectClientID = 1400  # ErrIncorrectClientID客户端ID不正确

    ErrInvalidSecretKey = 1500  # ErrInvalidSecretKeysecretkey错误

    ErrInvalidUser = 1600  # ErrInvalidUser无效的用户

    # ErrPhonePasswordNotMatch电话密码不匹配
    ErrPhonePasswordNotMatch = 1700
    # ErrExistPassword密码已经存在
    ErrExistPassword = 1800
    # ErrIncorrectPassword密码不正确
    ErrIncorrectPassword = 1900
    # ErrExsitRegisteredPhone注册手机号已经存在
    ErrExsitRegisteredPhone = 2000
    # ErrInValidMVC验证码失效
    ErrInValidMVC = 2100
    # ErrGetAccessTokenFailure失败去得到Access - Token
    ErrGetAccessTokenFailure = 2200
    # Request 数据错误

    ErrParsingRequestFailed = 2300
    ErrValueRequired = 2400
    ErrInvalidValue = 2500
    ErrInvalidPassword = 2600
    # RPC请求相关

    ErrRPCInternal = 3000
    ErrRPCTimeout = 3001

    # ErrBuildJwtToken JWT Token 生成错误
    ErrBuildJwtToken = 4001

    ErrClientInternal = 5000

    # ErrInvalidValidationValue非法的验证方式的值
    ErrInvalidValidationValue = 10000
    # ErrInvalidSecureQuestionValidationMethod非法的安全问题验证方法
    ErrInvalidSecureQuestionValidationMethod = 11000
    # ErrWrongSecureQuestionCount安全问题数量不匹配
    ErrWrongSecureQuestionCount = 12000
    # ErrSecureQuestionExist已经设置过密保
    ErrSecureQuestionExist = 13000
    # ErrEmptySecureQuestion密保问题为空
    ErrEmptySecureQuestion = 14000
    # ErrEmptyAnswer答案为空
    ErrEmptyAnswer = 15000
    # ErrSamePassword新密码与旧密码相同
    ErrSamePassword = 16000
    # ErrMismatchQuestion传入的密保问题错误
    ErrMismatchQuestion = 17000
    # ErrWrongFormatQuestion传入的密保问题格式错误
    ErrWrongFormatQuestion = 18000
    # ErrRepeatedQuestion问题重复
    ErrRepeatedQuestion = 19000
    # ErrWrongFormatPhone手机号格式错误
    ErrWrongFormatPhone = 20000
    # ErrNonexistentUsername用户名不存在
    ErrNonexistentUsername = 21000
    # ErrNonexistentPhone手机号不存在
    ErrNonexistentPhone = 22000
    # ErrInvalidEmailAddress邮箱格式错误
    ErrInvalidEmailAddress = 23000
    # ErrSecureEmailExists用户已经设置了自己的邮箱
    ErrSecureEmailExists = 24000
    # ErrSecureEmailUsedByOthers邮箱已经被其他人设置
    ErrSecureEmailUsedByOthers = 25000
    # ErrSecureEmailNotSet没有设置邮箱
    ErrSecureEmailNotSet = 26000
    # ErrSecureEmailAddressNotMatched与原邮箱不匹配
    ErrSecureEmailAddressNotMatched = 27000
    # ErrInvalidEmailNotificationAction非法的邮件通知的操作
    ErrInvalidEmailNotificationAction = 28000
    # ErrUsedVcRecord该记录已经被使用过
    ErrUsedVcRecord = 29000
    # ErrExpiredVcRecord该记录已经过期
    ErrExpiredVcRecord = 30000
    # InvalidVcRecord验证码错误
    InvalidVcRecord = 31000
    # ErrInvalidRequestCountIn1Minute 1分钟请求次数非法
    ErrInvalidRequestCountIn1Minute = 32000
    # ErrExsitRegion区域已经存在
    ErrExsitRegion = 33000
    # ErrExsitSigninPhone登录电话已经存在
    ErrExsitSigninPhone = 34000
    # ErrInvalidSigninPhone无效的登录电话
    ErrInvalidSigninPhone = 35000
    # ErrInvalidVerificationNumber无效的VerificationNumber
    ErrInvalidVerificationNumber = 36000
    # ErrNotExistSigninPhone登录手机号不存在
    ErrNotExistSigninPhone = 37000
    # ErrSamePhone新旧手机号一样
    ErrSamePhone = 38000
    # ErrWrongSendVia发送验证码的方式错误
    ErrWrongSendVia = 39000
    # ErrInvalidSendValue发送验证码的方式的值错误
    ErrInvalidSendValue = 40000
    # ErrInvalidSigninEmail无效的安全邮箱
    ErrInvalidSigninEmail = 41000
    # ErrInvalidValidationType验证邮箱验证码时非法的验证类型
    ErrInvalidValidationType = 42000
    # ErrInvalidValidationMethod获取方式非法
    ErrInvalidValidationMethod = 43000
    # ErrNotExistNewSecureEmail新安全邮箱不存在
    ErrNotExistNewSecureEmail = 45000
    # ErrSameEmail新旧安全邮箱相同
    ErrSameEmail = 46000
    # ErrSameSecureQuestion新旧密保问题一样
    ErrSameSecureQuestion = 47000
    # ErrNotExistOldPassword旧密码不存在
    ErrNotExistOldPassword = 48000
    # ErrNonexistentSecureQuestions密保问题未设置
    ErrNonexistentSecureQuestions = 49000
    # ErrNoneExistSecureEmail 邮箱不存在
    ErrNoneExistSecureEmail = 50000
    # ErrNationCode区号不正确
    ErrNationCode = 51000
    # ErrWrongSmsNotificationType短信类型不正确
    ErrWrongSmsNotificationType = 52000
    # ErrSignInPhoneNotBelongsToUser手机号不属于用户
    ErrSignInPhoneNotBelongsToUser = 53000
    # ErrUsernameNotSet用户名未设置
    ErrUsernameNotSet = 54000
    # ErrNonexistentPassword密码不存在
    ErrNonexistentPassword = 55000
    # ErrIncorrectOldPassword旧密码错误
    ErrIncorrectOldPassword = 56000
    ErrIncorrectOldPassword = 56000
    # ErrWrongFormatOfNickname昵称格式错误
    ErrWrongFormatOfNickname = 57000
    # ErrSensitiveWordsInNickname 昵称包含敏感词
    ErrSensitiveWordsInNickname = 58000
    # ErrReservedWordsInNickname 昵称包含保留词
    ErrReservedWordsInNickname = 59000
    # ErrMaskWordsInNickname 昵称包含屏蔽词
    ErrMaskWordsInNickname = 60000
    # ErrEmptyNickname 昵称为空
    ErrEmptyNickname = 61000
    # ErrEmptyGender性别为空
    ErrEmptyGender = 62000
    # ErrEmptyBirthday 生日为空
    ErrEmptyBirthday = 63000
    # ErrEmptyHeight 身高为空
    ErrEmptyHeight = 64000
    # ErrEmptyWeight 体重为空
    ErrEmptyWeight = 65000
    # ErrEmptyLanguage 语言为空
    ErrEmptyLanguage = 66000
    # ErrEmptyRegion 区域为空
    ErrEmptyRegion = 67000
    # ErrWrongFormatOfPassword 密码格式错误
    ErrWrongFormatOfPassword = 68000
    # ErrEmptyPassword 密码为空
    ErrEmptyPassword = 69000
    # ErrNotEmailOfCurrentUser  非当前绑定邮箱
    ErrNotEmailOfCurrentUser = 70000
    # ErrNoneExistUser 不存在的User
    ErrNoneExistUser = 71000
    # ErrSendMoreSMSInOneMinute 一分钟内发送多条短信
    ErrSendMoreSMSInOneMinute = 72000
    # ErrSendSMS 短信发送异常
    ErrSendSMS = 74000
    # ErrSendEmail   邮件发送异常
    ErrSendEmail = 75000
