package handler

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/jinmukeji/jiujiantang-services/jinmuid/mysqldb"
	"github.com/jinmukeji/go-pkg/crypto/rand"
	jinmuidpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	sempb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/sem/v1"
)

var (
	// Email 邮箱格式
	validEmail = regexp.MustCompile(`^[A-Za-z0-9!#$%&'+/=?^_{|}~-]+(.[A-Za-z0-9!#$%&'+/=?^_{|}~-]+)*@([A-Za-z0-9]+(?:-[A-Za-z0-9]+)?.)+[A-Za-z0-9]+(-[A-Za-z0-9]+)?$`)
)

const (
	// 邮件的有效时间
	validEmailTime = time.Minute * 30
	// limitsIn24Hours 单人的单个模板发送邮件24小时内上限（50条)
	limitsIn24Hours = 50
	// limitsAllIn24Hours 单人的所有模板发送邮件24小时内上限（250条)
	limitsAllIn24Hours = 250
	// limitsIn1Minutes 单人一分钟内发送邮件上限(1条)
	limitsIn1Minutes = 1
	// MinuteOf24Hour 24小时的分钟数量
	MinuteOf24Hour = 60 * 24
	// MinuteInAnHour 一小时的分钟数量
	MinuteInAnHour = 60
	// HourSuffix 发送的等待时间的小时的后缀
	HourSuffix = "H"
	// MinuteSuffix 发送的等待时间的分钟后缀
	MinuteSuffix = "M"

	// FindResetPasswordDescription 验证码类型描述：找回/重置密码
	FindResetPasswordDescription = "找回/重置密码"
	// FindUsernameDescription 验证码类型描述：找回用户名
	FindUsernameDescription = "找回用户名"
	// SetSecureEmailDescription 验证码类型描述：设置安全邮箱
	SetSecureEmailDescription = "设置安全邮箱"
	// ModifySecureEmailDescription 验证码类型描述：修改安全邮箱
	ModifySecureEmailDescription = "修改安全邮箱"
	// UnsetSecureEmailDescription 验证码类型描述：解绑安全邮箱
	UnsetSecureEmailDescription = "解绑安全邮箱"

	// AllExceedLimitIn24Hour 所有模板验证码超过24小时提醒内容
	AllExceedLimitIn24Hour = "您本日接收验证码的次数已用完，请%d小时%d分后重试"
	// SingleExceedLimitIn24Hour 单个模板验证码超过24小时提醒内容
	SingleExceedLimitIn24Hour = "您的邮箱本日接收%s验证码的次数已用完，请%d小时%d分后重试"
)

// LoggedInEmailNotification 已登录状态下的邮件通知
func (j *JinmuIDService) LoggedInEmailNotification(ctx context.Context, req *jinmuidpb.LoggedInEmailNotificationRequest, resp *jinmuidpb.LoggedInEmailNotificationResponse) error {
	// 发送邮件的时候首先判断在24小时之内是否达到了指定验证码个数，如果达到了，则直接返回距离下次发送的时间
	// 没达到则:
	// 先判断一分钟内是否已经发送了超过指定条数，超过了则返回错误
	// 否则则进入邮件网关发送邮件，发送成功后增加审计记录
	// ErrSecureEmailAddressNotMatched
	// NotEmailOfCurrentUse
	// 验证用户
	token, ok := TokenFromContext(ctx)
	if !ok {
		return NewError(ErrInvalidUser, errors.New("failed to get token from context"))
	}
	userID, err := j.datastore.FindUserIDByToken(ctx, token)
	if err != nil {
		return NewError(ErrUserUnauthorized, fmt.Errorf("failed to get userID by token: %s", err.Error()))
	}
	if userID != req.UserId {
		return NewError(ErrInvalidUser, fmt.Errorf("user %d from request and user %d from token are inconsistent", req.UserId, userID))
	}
	_, errFindUserByUserID := j.datastore.FindUserByUserID(ctx, req.UserId)
	if errFindUserByUserID != nil {
		return NewError(ErrNoneExistUser, fmt.Errorf("failed to find user by user %d: %s", req.UserId, errFindUserByUserID.Error()))
	}
	// 判断行为
	if req.Action == jinmuidpb.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_INVALID || req.Action == jinmuidpb.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_UNSET {
		return NewError(ErrInvalidEmailNotificationAction, errors.New("unknown email notification action"))
	}

	// 验证邮件格式
	if !checkEmailFormat(req.Email) {
		return NewError(ErrInvalidEmailAddress, fmt.Errorf("wrong format of email address %s", req.Email))
	}
	// 设置安全邮箱
	if req.Action == jinmuidpb.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_SET_SECURE_EMAIL {
		// 设置安全邮箱时，判断邮箱是否已经被任何人设置
		hasSetSecureEmailByAnyone, _ := j.datastore.HasSetSecureEmailByAnyone(ctx, req.Email)
		if hasSetSecureEmailByAnyone {
			return NewError(ErrSecureEmailUsedByOthers, fmt.Errorf("failed to check if secure email %s has been set by anyone", req.Email))
		}
	}
	//  解除绑定安全邮箱
	if req.Action == jinmuidpb.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_UNSET_SECURE_EMAIL {
		// 判断邮箱是否被任何人设置
		existsEmail, _ := j.datastore.HasSecureEmailSet(ctx, req.Email)
		if !existsEmail {
			return NewError(ErrNoneExistSecureEmail, fmt.Errorf("secure email %s is not set by anyone", req.Email))
		}
		// 判断是当前绑定邮箱
		user, errFindUserBySecureEmail := j.datastore.FindUserBySecureEmail(ctx, req.Email)
		if errFindUserBySecureEmail != nil {
			return NewError(ErrNoneExistSecureEmail, fmt.Errorf("failed to find user by secure email %s: %s", req.Email, errFindUserBySecureEmail.Error()))
		}
		if user.UserID != req.UserId {
			return NewError(ErrSecureEmailAddressNotMatched, fmt.Errorf("secure email %s is not email of current user %d", req.Email, req.UserId))
		}
	}
	// 修改安全邮箱
	if req.Action == jinmuidpb.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_MODIFY_SECURE_EMAIL {
		// 向原邮箱发送验证码
		if !req.SendToNewIfModify {
			// 判断是当前绑定邮箱
			user, errFindUserBySecureEmail := j.datastore.FindUserBySecureEmail(ctx, req.Email)
			if errFindUserBySecureEmail != nil {
				return NewError(ErrNoneExistSecureEmail, fmt.Errorf("failed to find user by secure email %s: %s", req.Email, errFindUserBySecureEmail.Error()))
			}
			if user.UserID != req.UserId {
				return NewError(ErrSecureEmailAddressNotMatched, fmt.Errorf("secure email %s is not email of current user %d", req.Email, req.UserId))
			}
		} else {
			// 向新邮箱发送验证码
			// 判断新旧安全邮箱是否相同
			user, errFindUserByUserID := j.datastore.FindUserByUserID(ctx, req.UserId)
			if errFindUserByUserID != nil {
				return NewError(ErrNoneExistUser, fmt.Errorf("failed to find user by userID %d: %s", req.UserId, errFindUserByUserID.Error()))
			}
			if !user.HasSetEmail {
				return NewError(ErrSecureEmailNotSet, fmt.Errorf("email of user %d has not been set", user.UserID))
			}
			if user.SecureEmail == req.Email {
				return NewError(ErrSameEmail, fmt.Errorf("new email %s cannot be the same as old email", req.Email))
			}
			// 判断新邮箱是否已经被其他人设置
			hasSetSecureEmailByAnyone, _ := j.datastore.HasSetSecureEmailByAnyone(ctx, req.Email)
			if hasSetSecureEmailByAnyone {
				return NewError(ErrSecureEmailUsedByOthers, fmt.Errorf("secure email %s has been set", req.Email))
			}

		}
	}
	usage, errmapProtoLoggedInSemTemplateActionToDB := mapProtoLoggedInSemTemplateActionToDB(req.Action)
	if errmapProtoLoggedInSemTemplateActionToDB != nil {
		return NewError(ErrInvalidEmailNotificationAction, errmapProtoLoggedInSemTemplateActionToDB)
	}
	if usage == mysqldb.Unknown {
		return NewError(ErrInvalidEmailNotificationAction, errors.New("unknown email notification action"))
	}

	// 所有模板的邮件是否达到数量限制
	sendAllCountIn24h, err := j.datastore.SearchVcRecordCountsIn24hours(ctx, req.Email)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to search vc record counts of email %s in 24 hours: %s", req.Email, err.Error()))
	}
	// 达到24小时内限制
	if sendAllCountIn24h >= limitsAllIn24Hours {
		// 计算距离下次发送的时间
		tryAfter, errGetInterval := j.getIntervalBeforeSendAll(req.Email)
		if errGetInterval != nil {
			return errGetInterval
		}
		resp.Message = fmt.Sprintf(AllExceedLimitIn24Hour, tryAfter[0], tryAfter[1])
		return nil
	}

	// 单个模板的邮件是否达到数量限制
	sendCountIn24h, err := j.datastore.SearchSpecificVcRecordCountsIn24hours(ctx, req.Email, usage)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to search specific Vc record count of email %s when the usage is %s in 24 hours: %s", req.Email, usage, err.Error()))
	}
	// 达到24小时内限制
	if sendCountIn24h >= limitsIn24Hours {
		// 计算距离下次发送的时间
		tryAfter, errGetInterval := j.getIntervalBeforeSendSingle(ctx, req.Email, usage)
		if errGetInterval != nil {
			return errGetInterval
		}
		stringDescription, errMapLoggedInFromProtoToRemindDescription := mapLoggedInFromProtoToRemindDescription(req.Action)
		if errMapLoggedInFromProtoToRemindDescription != nil {
			return NewError(ErrInvalidEmailNotificationAction, errMapLoggedInFromProtoToRemindDescription)
		}
		resp.Message = fmt.Sprintf(SingleExceedLimitIn24Hour, stringDescription, tryAfter[0], tryAfter[1])
		return nil
	}

	// 一分钟内是否发送了多条
	sendAllCountIn1m, errSendAllCountIn1m := j.datastore.SearchVcRecordCountsIn1Minute(ctx, req.Email)
	if errSendAllCountIn1m != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to search Vc record counts of email %s in 1 minute: %s", req.Email, errSendAllCountIn1m.Error()))
	}
	if sendAllCountIn1m >= limitsIn1Minutes {
		return NewError(ErrInvalidRequestCountIn1Minute, errors.New("invalid request counts in 1 minute"))
	}
	// 进行邮件发送
	emailNotificationRequest := new(sempb.EmailNotificationRequest)
	emailNotificationRequest.ToAddress = []string{req.Email}

	protoAction, errmapDBLoggedInSemTemplateActionToProto := mapDBLoggedInSemTemplateActionToProto(req.Action)
	if errmapDBLoggedInSemTemplateActionToProto != nil {
		return NewError(ErrInvalidEmailNotificationAction, errmapDBLoggedInSemTemplateActionToProto)
	}
	emailNotificationRequest.TemplateAction = protoAction
	// 随机6位数字
	code, _ := rand.RandomStringWithMask(rand.MaskDigits, 6)
	emailNotificationRequest.TemplateParam = map[string]string{
		"code": code,
	}
	emailNotificationRequest.Language = req.Language
	// TODO: EmailNotification与CreateVcRecord要在同一个事务
	_, errEmailNotification := j.semSvc.EmailNotification(ctx, emailNotificationRequest)

	if errEmailNotification != nil {
		return NewError(ErrSendEmail, fmt.Errorf("send email error: %s", errEmailNotification.Error()))
	}
	serialNumber := generateNotificationSerialNumber()

	expireAt, _ := ptypes.TimestampProto(time.Now().Add(validEmailTime))
	expiredTime, _ := ptypes.Timestamp(expireAt)
	vcRecord := &mysqldb.VcRecord{
		Usage:     usage,
		SN:        serialNumber,
		Code:      code,
		SendVia:   mysqldb.Email,
		SendTo:    req.Email,
		ExpiredAt: &expiredTime,
		HasUsed:   false,
	}
	errCreateMVCRecord := j.datastore.CreateVcRecord(ctx, vcRecord)
	if errCreateMVCRecord != nil {
		return NewError(ErrDatabase, errors.New("failed to create Vc record"))
	}
	resp.SerialNumber = serialNumber
	resp.Acknowledged = true
	return nil

}

// NotLoggedInEmailNotification 未登录状态下的邮件通知
func (j *JinmuIDService) NotLoggedInEmailNotification(ctx context.Context, req *jinmuidpb.NotLoggedInEmailNotificationRequest, resp *jinmuidpb.NotLoggedInEmailNotificationResponse) error {
	if req.Action == jinmuidpb.NotLoggedInSemTemplateAction_NOT_LOGGED_IN_SEM_TEMPLATE_ACTION_NOT_LOGGED_IN_UNKNOW {
		return NewError(ErrInvalidEmailNotificationAction, errors.New("unknown email notification action"))
	}
	// 判断邮箱是否被任何人设置，只有被设置了才可以发送除了设置安全邮箱之外的邮件
	if req.Action == jinmuidpb.NotLoggedInSemTemplateAction_NOT_LOGGED_IN_SEM_TEMPLATE_ACTION_FIND_RESET_PASSWORD || req.Action == jinmuidpb.NotLoggedInSemTemplateAction_NOT_LOGGED_IN_SEM_TEMPLATE_ACTION_FIND_USERNAME {
		existsEmail, _ := j.datastore.HasSecureEmailSet(ctx, req.Email)
		if !existsEmail {
			return NewError(ErrNoneExistSecureEmail, fmt.Errorf("secure email %s is not set by anyone", req.Email))
		}
	}
	if req.Action == jinmuidpb.NotLoggedInSemTemplateAction_NOT_LOGGED_IN_SEM_TEMPLATE_ACTION_FIND_RESET_PASSWORD {
		// 重置密码时判断密码是否已经被设置
		user, _ := j.datastore.FindUserBySecureEmail(ctx, req.Email)
		if !user.HasSetPassword {
			return NewError(ErrNotExistOldPassword, fmt.Errorf("old password of user whose email is %s does not exist", req.Email))
		}
	}
	// 验证邮件格式
	if !checkEmailFormat(req.Email) {
		return NewError(ErrInvalidEmailAddress, fmt.Errorf("wrong format of email address %s", req.Email))
	}
	usage, errmapProtoNotLoggedInSemTemplateActionToDB := mapProtoNotLoggedInSemTemplateActionToDB(req.Action)
	if errmapProtoNotLoggedInSemTemplateActionToDB != nil {
		return NewError(ErrInvalidEmailNotificationAction, errmapProtoNotLoggedInSemTemplateActionToDB)
	}

	if usage == mysqldb.Unknown {
		return NewError(ErrInvalidEmailNotificationAction, errors.New("unknown email notification action"))
	}
	if req.Action == jinmuidpb.NotLoggedInSemTemplateAction_NOT_LOGGED_IN_SEM_TEMPLATE_ACTION_FIND_USERNAME {
		_, errFindUsernameBySecureEmail := j.datastore.FindUsernameBySecureEmail(ctx, req.Email)
		if errFindUsernameBySecureEmail != nil {
			return NewError(ErrNonexistentUsername, fmt.Errorf("failed to find username by secure email %s: %s", req.Email, errFindUsernameBySecureEmail.Error()))
		}
	}

	// 所有模板的邮件是否达到数量限制
	sendAllCountIn24h, err := j.datastore.SearchVcRecordCountsIn24hours(ctx, req.Email)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to search Vc record count of email %s in 24 hours: %s", req.Email, err.Error()))
	}
	// 达到24小时内限制
	if sendAllCountIn24h >= limitsAllIn24Hours {
		// 计算距离下次发送的时间
		tryAfter, errGetInterval := j.getIntervalBeforeSendAll(req.Email)
		if errGetInterval != nil {
			return errGetInterval
		}
		resp.Message = fmt.Sprintf(AllExceedLimitIn24Hour, tryAfter[0], tryAfter[1])
		return nil
	}

	// 单个模板的邮件是否达到数量限制
	sendCountIn24h, err := j.datastore.SearchSpecificVcRecordCountsIn24hours(ctx, req.Email, usage)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to search specific Vc record count of email %s when the usage is %s in 24 hours: %s", req.Email, usage, err.Error()))
	}
	// 达到24小时内限制
	if sendCountIn24h >= limitsIn24Hours {
		// 计算距离下次发送的时间
		tryAfter, errGetInterval := j.getIntervalBeforeSendSingle(ctx, req.Email, usage)
		if errGetInterval != nil {
			return errGetInterval
		}
		stringAction, errMapNotLoggedInFromProtoToRemindDescription := mapNotLoggedInFromProtoToRemindDescription(req.Action)
		if errMapNotLoggedInFromProtoToRemindDescription != nil {
			return errMapNotLoggedInFromProtoToRemindDescription
		}
		resp.Message = fmt.Sprintf(SingleExceedLimitIn24Hour, stringAction, tryAfter[0], tryAfter[1])
		return nil
	}

	// 一分钟内是否发送了多条
	sendCountIn1m, errSendAllCountIn1m := j.datastore.SearchVcRecordCountsIn1Minute(ctx, req.Email)
	if errSendAllCountIn1m != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to search Vc record counts of email %s in 1 minute: %s", req.Email, errSendAllCountIn1m.Error()))
	}

	if sendCountIn1m >= limitsIn1Minutes {
		return NewError(ErrInvalidRequestCountIn1Minute, errors.New("invalid request counts in 1 minute"))
	}

	emailNotificationRequest := new(sempb.EmailNotificationRequest)
	emailNotificationRequest.ToAddress = []string{req.Email}

	semProtoAction, errmapDBNotLoggedInSemProtoActionToProto := mapDBNotLoggedInSemProtoActionToProto(req.Action)
	if errmapDBNotLoggedInSemProtoActionToProto != nil {
		return errmapDBNotLoggedInSemProtoActionToProto
	}
	emailNotificationRequest.TemplateAction = semProtoAction

	// 随机6位数字
	code, _ := rand.RandomStringWithMask(rand.MaskDigits, 6)
	emailNotificationRequest.TemplateParam = map[string]string{
		"code": code,
	}
	emailNotificationRequest.Language = req.Language
	_, errEmailNotification := j.semSvc.EmailNotification(ctx, emailNotificationRequest)
	if errEmailNotification != nil {
		return NewError(ErrSendEmail, errors.New("send email error"))
	}
	serialNumber := generateNotificationSerialNumber()

	expireAt, _ := ptypes.TimestampProto(time.Now().Add(validEmailTime))
	expiredTime, _ := ptypes.Timestamp(expireAt)

	vcRecord := &mysqldb.VcRecord{
		Usage:     usage,
		SN:        serialNumber,
		Code:      code,
		SendVia:   mysqldb.Email,
		SendTo:    req.Email,
		ExpiredAt: &expiredTime,
		HasUsed:   false,
	}
	errCreateMVCRecord := j.datastore.CreateVcRecord(ctx, vcRecord)
	if errCreateMVCRecord != nil {
		return NewError(ErrDatabase, errors.New("failed to create Vc record"))
	}
	resp.SerialNumber = serialNumber
	resp.Acknowledged = true
	return nil

}

// 根据Proto获取已登录状态下在数据库中对应的Usage
func mapProtoLoggedInSemTemplateActionToDB(action jinmuidpb.LoggedInSemTemplateAction) (mysqldb.Usage, error) {
	switch action {
	case jinmuidpb.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_INVALID:
		return mysqldb.SetSecureEmail, fmt.Errorf("invalid proto logged in sem template action %d", action)
	case jinmuidpb.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_UNSET:
		return mysqldb.ModifySecureEmail, fmt.Errorf("invalid proto logged in sem template action %d", action)
	case jinmuidpb.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_SET_SECURE_EMAIL:
		return mysqldb.SetSecureEmail, nil
	case jinmuidpb.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_MODIFY_SECURE_EMAIL:
		return mysqldb.ModifySecureEmail, nil
	case jinmuidpb.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_UNSET_SECURE_EMAIL:
		return mysqldb.UnsetSecureEmail, nil
	}
	return mysqldb.SetSecureEmail, fmt.Errorf("invalid proto logged in sem template action %d", action)
}

// 根据Proto获取未登录状态下在数据库中对应的Usage
func mapProtoNotLoggedInSemTemplateActionToDB(action jinmuidpb.NotLoggedInSemTemplateAction) (mysqldb.Usage, error) {
	switch action {
	case jinmuidpb.NotLoggedInSemTemplateAction_NOT_LOGGED_IN_SEM_TEMPLATE_ACTION_INVALID:
		return mysqldb.Unknown, fmt.Errorf("invalid proto not logged in sem template action %d", action)
	case jinmuidpb.NotLoggedInSemTemplateAction_NOT_LOGGED_IN_SEM_TEMPLATE_ACTION_UNSET:
		return mysqldb.Unknown, fmt.Errorf("invalid proto not logged in sem template action %d", action)
	case jinmuidpb.NotLoggedInSemTemplateAction_NOT_LOGGED_IN_SEM_TEMPLATE_ACTION_NOT_LOGGED_IN_UNKNOW:
		return mysqldb.FindUsername, nil
	case jinmuidpb.NotLoggedInSemTemplateAction_NOT_LOGGED_IN_SEM_TEMPLATE_ACTION_FIND_RESET_PASSWORD:
		return mysqldb.FindResetPassword, nil
	case jinmuidpb.NotLoggedInSemTemplateAction_NOT_LOGGED_IN_SEM_TEMPLATE_ACTION_FIND_USERNAME:
		return mysqldb.FindUsername, nil
	}
	return mysqldb.Unknown, fmt.Errorf("invalid proto not logged in sem template action %d", action)
}

// 计算所有模板邮件达到上限后距离下次可以发送验证码的时间间隔
func (j *JinmuIDService) getIntervalBeforeSendAll(email string) ([]int, error) {
	location, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now()
	t := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location).AddDate(0, 0, 1)
	subM := t.Sub(now)
	hoursInterval := int(subM.Hours())
	minutesInterval := int(subM.Minutes())
	retryHour := hoursInterval
	retryMinute := minutesInterval - MinuteInAnHour*hoursInterval
	return []int{retryHour, retryMinute}, nil
}

// 计算特定模板验证码有距离下次可以发送验证码的时间间隔
func (j *JinmuIDService) getIntervalBeforeSendSingle(ctx context.Context, email string, usage mysqldb.Usage) ([]int, error) {
	earliestVcRecord, err := j.datastore.SearchSpecificVcRecordEarliestTimeIn24hours(ctx, email, usage)
	if err != nil {
		return nil, NewError(ErrDatabase, NewError(ErrDatabase, fmt.Errorf("failed to search specific Vc record of email %s when usage is %s earliest time in 24 hours: %s", email, usage, err.Error())))
	}
	createdAt := *earliestVcRecord.ExpiredAt
	now := time.Now()
	subM := now.Sub(createdAt)
	minuteInterval := subM.Minutes()
	retryAfterMinute := MinuteOf24Hour - minuteInterval
	retryHour := int(retryAfterMinute / MinuteInAnHour)
	retryMinute := int(int(retryAfterMinute)-MinuteInAnHour*retryHour) + 1
	return []int{retryHour, retryMinute}, nil
}

// 获取已登录状态下进入邮件网关的行为
func mapDBLoggedInSemTemplateActionToProto(action jinmuidpb.LoggedInSemTemplateAction) (sempb.TemplateAction, error) {
	switch action {
	case jinmuidpb.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_INVALID:
		return sempb.TemplateAction_TEMPLATE_ACTION_INVALID, fmt.Errorf("invalid proto logged in sem template action %d", action)
	case jinmuidpb.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_UNSET:
		return sempb.TemplateAction_TEMPLATE_ACTION_INVALID, fmt.Errorf("invalid proto logged in sem template action %d", action)
	case jinmuidpb.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_SET_SECURE_EMAIL:
		return sempb.TemplateAction_TEMPLATE_ACTION_SET_SECURE_EMAIL, nil
	case jinmuidpb.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_MODIFY_SECURE_EMAIL:
		return sempb.TemplateAction_TEMPLATE_ACTION_MODIFY_SECURE_EMAIL, nil
	case jinmuidpb.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_UNSET_SECURE_EMAIL:
		return sempb.TemplateAction_TEMPLATE_ACTION_UNSET_SECURE_EMAIL, nil
	}
	return sempb.TemplateAction_TEMPLATE_ACTION_INVALID, fmt.Errorf("invalid proto logged in sem template action %d", action)
}

// 获取未登录状态下进入邮件网关的行为
func mapDBNotLoggedInSemProtoActionToProto(action jinmuidpb.NotLoggedInSemTemplateAction) (sempb.TemplateAction, error) {
	switch action {
	case jinmuidpb.NotLoggedInSemTemplateAction_NOT_LOGGED_IN_SEM_TEMPLATE_ACTION_INVALID:
		return sempb.TemplateAction_TEMPLATE_ACTION_INVALID, fmt.Errorf("invalid proto not logged in sem template action %d", action)
	case jinmuidpb.NotLoggedInSemTemplateAction_NOT_LOGGED_IN_SEM_TEMPLATE_ACTION_UNSET:
		return sempb.TemplateAction_TEMPLATE_ACTION_INVALID, fmt.Errorf("invalid proto not logged in sem template action %d", action)
	case jinmuidpb.NotLoggedInSemTemplateAction_NOT_LOGGED_IN_SEM_TEMPLATE_ACTION_NOT_LOGGED_IN_UNKNOW:
		return sempb.TemplateAction_TEMPLATE_ACTION_INVALID, fmt.Errorf("invalid proto not logged in sem template action %d", action)
	case jinmuidpb.NotLoggedInSemTemplateAction_NOT_LOGGED_IN_SEM_TEMPLATE_ACTION_FIND_RESET_PASSWORD:
		return sempb.TemplateAction_TEMPLATE_ACTION_FIND_RESET_PASSWORD, nil
	case jinmuidpb.NotLoggedInSemTemplateAction_NOT_LOGGED_IN_SEM_TEMPLATE_ACTION_FIND_USERNAME:
		return sempb.TemplateAction_TEMPLATE_ACTION_FIND_USERNAME, nil
	}
	return sempb.TemplateAction_TEMPLATE_ACTION_INVALID, fmt.Errorf("invalid proto not logged in sem template action %d", action)
}

// mapLoggedInFromProtoToRemindDescription 已登录时单个模板超过制定次数提醒时获取指定的描述
func mapLoggedInFromProtoToRemindDescription(action jinmuidpb.LoggedInSemTemplateAction) (string, error) {
	switch action {
	case jinmuidpb.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_INVALID:
		return "", fmt.Errorf("invalid proto logged in sem template action %d", action)
	case jinmuidpb.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_UNSET:
		return "", fmt.Errorf("invalid proto logged in sem template action %d", action)
	case jinmuidpb.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_SET_SECURE_EMAIL:
		return SetSecureEmailDescription, nil
	case jinmuidpb.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_MODIFY_SECURE_EMAIL:
		return ModifySecureEmailDescription, nil
	case jinmuidpb.LoggedInSemTemplateAction_LOGGED_IN_SEM_TEMPLATE_ACTION_UNSET_SECURE_EMAIL:
		return UnsetSecureEmailDescription, nil
	}
	return "", fmt.Errorf("invalid proto logged in sem template action %d", action)
}

// mapNotLoggedInFromProtoToRemindDescription 未登录时单个模板超过制定次数提醒时获取指定的描述
func mapNotLoggedInFromProtoToRemindDescription(action jinmuidpb.NotLoggedInSemTemplateAction) (string, error) {
	switch action {
	case jinmuidpb.NotLoggedInSemTemplateAction_NOT_LOGGED_IN_SEM_TEMPLATE_ACTION_INVALID:
		return "", fmt.Errorf("invalid proto not logged in sem template action %d", action)
	case jinmuidpb.NotLoggedInSemTemplateAction_NOT_LOGGED_IN_SEM_TEMPLATE_ACTION_UNSET:
		return "", fmt.Errorf("invalid proto not logged in sem template action %d", action)
	case jinmuidpb.NotLoggedInSemTemplateAction_NOT_LOGGED_IN_SEM_TEMPLATE_ACTION_NOT_LOGGED_IN_UNKNOW:
		return "", fmt.Errorf("invalid proto not logged in sem template action %d", action)
	case jinmuidpb.NotLoggedInSemTemplateAction_NOT_LOGGED_IN_SEM_TEMPLATE_ACTION_FIND_RESET_PASSWORD:
		return FindResetPasswordDescription, nil
	case jinmuidpb.NotLoggedInSemTemplateAction_NOT_LOGGED_IN_SEM_TEMPLATE_ACTION_FIND_USERNAME:
		return FindUsernameDescription, nil
	}
	return "", fmt.Errorf("invalid proto not logged in sem template action %d", action)
}

// 生成SN号
func generateNotificationSerialNumber() string {
	return uuid.New().String()
}

// checkEmailFormat 检查邮件格式
func checkEmailFormat(email string) bool {
	return validEmail.MatchString(email)
}
