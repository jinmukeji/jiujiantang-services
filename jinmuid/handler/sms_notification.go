package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/jinmukeji/jiujiantang-services/jinmuid/mysqldb"
	"github.com/jinmukeji/go-pkg/crypto/rand"
	jinmuidpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
	smspb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/sms/v1"
)

// 验证码有效时间2分钟
const (
	validTime = time.Minute * 2
	// smsLimitsIn24Hours 短信限制24小时10条
	smsLimitsIn24Hours = 10
)

// SmsNotification 短信通知
func (j *JinmuIDService) SmsNotification(ctx context.Context, req *jinmuidpb.SmsNotificationRequest, resp *jinmuidpb.SmsNotificationResponse) error {
	// 验证手机格式
	if validatePhoneFormat(req.Phone, req.NationCode) != nil {
		return NewError(ErrWrongFormatPhone, fmt.Errorf("phone %s%s format is invalid: %s", req.NationCode, req.Phone, validatePhoneFormat(req.Phone, req.NationCode)))
	}
	// 已经登录状态下需要验证用户
	if req.Action == jinmuidpb.TemplateAction_TEMPLATE_ACTION_SET_PHONE_NUMBER || req.Action == jinmuidpb.TemplateAction_TEMPLATE_ACTION_MODIFY_PHONE_NUMBER {
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
			return fmt.Errorf("failed to find user by userID %d: %s", req.UserId, errFindUserByUserID.Error())
		}
	}
	// 重置密码
	if req.Action == jinmuidpb.TemplateAction_TEMPLATE_ACTION_RESET_PASSWORD {
		// 判断手机号码是否存在
		user, errFindUserByPhone := j.datastore.FindUserByPhone(ctx, req.Phone, req.NationCode)
		if errFindUserByPhone != nil || !user.HasSetPhone {
			return NewError(ErrNoneExistentPhone, fmt.Errorf("failed to find username by phone %s%s", req.NationCode, req.Phone))
		}
		// 密码是否已经被设置
		if !user.HasSetPassword {
			return NewError(ErrNotExistOldPassword, fmt.Errorf("old password of user %d does not exist", req.UserId))
		}
	}
	// 登录
	if req.Action == jinmuidpb.TemplateAction_TEMPLATE_ACTION_SIGN_IN {
		// 判断手机号码是否存在
		_, errFindUserByPhone := j.datastore.FindUserByPhone(ctx, req.Phone, req.NationCode)
		if errFindUserByPhone != nil {
			return NewError(ErrNoneExistentPhone, fmt.Errorf("failed to find username by phone %s%s: %s", req.NationCode, req.Phone, errFindUserByPhone.Error()))
		}
	}
	// 修改手机号
	if req.Action == jinmuidpb.TemplateAction_TEMPLATE_ACTION_MODIFY_PHONE_NUMBER {
		if !req.SendToNewIfModify {
			// 找到设置该手机号码的用户
			user, errFindUserByPhone := j.datastore.FindUserByPhone(ctx, req.Phone, req.NationCode)
			if errFindUserByPhone != nil {
				return NewError(ErrNoneExistentPhone, fmt.Errorf("failed to find username by phone %s%s: %s", req.NationCode, req.Phone, errFindUserByPhone.Error()))
			}
			// 判断当前手机号码的用户是否是当前用户
			if user.UserID != req.UserId {
				return NewError(ErrSignInPhoneNotBelongsToUser, fmt.Errorf("phone %s%s does not belong to current user %d", req.NationCode, req.Phone, req.UserId))
			}
		} else {
			user, errFindUserByPhone := j.datastore.FindUserByPhone(ctx, req.Phone, req.NationCode)
			if errFindUserByPhone == nil {
				if user.UserID == req.UserId {
					return NewError(ErrSamePhone, fmt.Errorf("new phone %s%s shouldn't be same as the old one", req.NationCode, req.Phone))
				} else {
					return NewError(ErrExistRegisteredPhone, fmt.Errorf("phone %s%s has been registered", req.NationCode, req.Phone))
				}
			}
		}
	}
	// 设置手机号码
	if req.Action == jinmuidpb.TemplateAction_TEMPLATE_ACTION_SET_PHONE_NUMBER {
		exsit, err := j.datastore.ExistPhone(ctx, req.Phone, req.NationCode)
		if err != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to check existence of phone %s%s: %s", req.NationCode, req.Phone, err.Error()))
		}
		if exsit {
			return NewError(ErrExistRegisteredPhone, fmt.Errorf("signin phone %s%s has been set", req.NationCode, req.Phone))
		}
	}
	// 注册要查询电话号码是否已经存在
	if req.Action == jinmuidpb.TemplateAction_TEMPLATE_ACTION_SIGN_UP {
		exsit, err := j.datastore.ExistPhone(ctx, req.Phone, req.NationCode)
		if err != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to check existence of phone %s%s: %s", req.NationCode, req.Phone, err.Error()))
		}

		if exsit {
			return NewError(ErrExistRegisteredPhone, fmt.Errorf("signin phone %s%s has been set", req.NationCode, req.Phone))
		}
	}
	// 进行发送前的判断
	sendCountIn24h, err := j.datastore.SearchVcRecordCountsIn24hours(ctx, req.Phone)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to search vc record counts of phone %s in 24 hours: %s", req.Phone, err.Error()))
	}
	// 达到24小时内限制
	if sendCountIn24h >= smsLimitsIn24Hours && !req.IsForced {
		location, _ := time.LoadLocation("Asia/Shanghai")
		now := time.Now()
		t := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location).AddDate(0, 0, 1)
		subM := t.Sub(now)
		hoursInterval := int(subM.Hours())
		minutesInterval := int(subM.Minutes())
		retryHour := hoursInterval
		retryMinute := minutesInterval - MinuteInAnHour*hoursInterval
		resp.Acknowledged = false
		resp.Message = fmt.Sprintf("您本日接收验证码的次数已用完，请%d小时%d分后重试", retryHour, retryMinute)
		return nil
	}
	// 一分钟内是否发送了多条
	sendCountIn1m, err := j.datastore.SearchVcRecordCountsIn1Minute(ctx, req.Phone)
	if err != nil {
		return NewError(ErrSendMoreSMSInOneMinute, errors.New("send more sms in one minute"))
	}
	if sendCountIn1m >= limitsIn1Minutes && !req.IsForced {
		earliestVcRecord, err := j.datastore.SearchVcRecordEarliestTimeIn1Minute(ctx, req.Phone)
		if err != nil {
			return NewError(ErrDatabase, fmt.Errorf("failed to search vc record earliest time in 1 minute of phone %s: %s", req.Phone, err.Error()))
		}
		now := time.Now()
		retrySecond := int(earliestVcRecord.CreatedAt.Add(time.Minute * 1).Sub(now).Seconds())
		resp.Acknowledged = false
		resp.Message = fmt.Sprintf("重新获取（%d）", retrySecond)
		return nil
	}
	var code string
	// 1-2分钟内是否有短信发送
	usage, errmapProtoSmsTemplateActionToDB := mapProtoSmsTemplateActionToDB(req.Action)
	if errmapProtoSmsTemplateActionToDB != nil {
		return NewError(ErrInvalidSmsTemplateAction, fmt.Errorf("send sms error: %s", errmapProtoSmsTemplateActionToDB.Error()))
	}
	vc, errSearchVcRecordFrom1MinuteTo2Mintue := j.datastore.SearchVcRecordFrom1MinuteTo2Mintue(ctx, req.Phone, usage)
	if errSearchVcRecordFrom1MinuteTo2Mintue == nil && len(vc.Code) == 6 {
		code = vc.Code
	} else {
		// 随机6位数字
		code, _ = rand.RandomStringWithMask(rand.MaskDigits, 6)
	}
	reqSendMessageRequest := new(smspb.SendMessageRequest)
	// 默认关闭强制性发送
	reqSendMessageRequest.IsForced = req.IsForced
	reqSendMessageRequest.Phone = req.Phone
	smsProtoAction, errmapJinmuidProtoTemplateActionToSmsProto := mapJinmuidProtoTemplateActionToSmsProto(req.Action)
	if errmapJinmuidProtoTemplateActionToSmsProto != nil {
		return NewError(ErrInvalidSmsTemplateAction, errmapJinmuidProtoTemplateActionToSmsProto)
	}
	reqSendMessageRequest.TemplateAction = smsProtoAction
	reqSendMessageRequest.NationCode = req.NationCode
	reqSendMessageRequest.TemplateParam = map[string]string{
		"code": code,
	}
	reqSendMessageRequest.Language = req.Language
	// TODO: SendMessage,CreateVcRecord要在同一个事务
	_, errSendMessage := j.smsSvc.SendMessage(ctx, reqSendMessageRequest)
	if errSendMessage != nil {
		return NewError(ErrSendSMS, fmt.Errorf("send sms error: %s", errSendMessage.Error()))
	}
	SN := uuid.New().String()
	expiredAt := time.Now().Add(validTime)
	vcRecord := &mysqldb.VcRecord{
		Usage:      usage,
		SN:         SN,
		Code:       code,
		SendVia:    mysqldb.Phone,
		SendTo:     req.Phone,
		ExpiredAt:  &expiredAt,
		HasUsed:    false,
		NationCode: req.NationCode,
	}
	errCreateMVCRecord := j.datastore.CreateVcRecord(ctx, vcRecord)
	if errCreateMVCRecord != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to create vc record: %s", errCreateMVCRecord.Error()))
	}
	resp.SerialNumber = SN
	resp.Acknowledged = true
	return nil
}

func mapProtoSmsTemplateActionToDB(action jinmuidpb.TemplateAction) (mysqldb.Usage, error) {
	switch action {
	case jinmuidpb.TemplateAction_TEMPLATE_ACTION_INVALID:
		return mysqldb.Unknown, fmt.Errorf("invalid proto template action %d", action)
	case jinmuidpb.TemplateAction_TEMPLATE_ACTION_UNSET:
		return mysqldb.Unknown, fmt.Errorf("invalid proto template action %d", action)
	case jinmuidpb.TemplateAction_TEMPLATE_ACTION_SIGN_IN:
		return mysqldb.SignIn, nil
	case jinmuidpb.TemplateAction_TEMPLATE_ACTION_SIGN_UP:
		return mysqldb.SignUp, nil
	case jinmuidpb.TemplateAction_TEMPLATE_ACTION_RESET_PASSWORD:
		return mysqldb.ResetPassword, nil
	case jinmuidpb.TemplateAction_TEMPLATE_ACTION_SET_PHONE_NUMBER:
		return mysqldb.SetPhoneNumber, nil
	case jinmuidpb.TemplateAction_TEMPLATE_ACTION_MODIFY_PHONE_NUMBER:
		return mysqldb.ModifyPhoneNumber, nil
	default:
		return mysqldb.Unknown, fmt.Errorf("invalid proto template action %d", action)
	}
}

func mapJinmuidProtoTemplateActionToSmsProto(action jinmuidpb.TemplateAction) (smspb.TemplateAction, error) {
	switch action {
	case jinmuidpb.TemplateAction_TEMPLATE_ACTION_INVALID:
		return smspb.TemplateAction_TEMPLATE_ACTION_INVALID, fmt.Errorf("invalid proto template action %d", action)
	case jinmuidpb.TemplateAction_TEMPLATE_ACTION_UNSET:
		return smspb.TemplateAction_TEMPLATE_ACTION_UNSET, fmt.Errorf("invalid proto template action %d", action)
	case jinmuidpb.TemplateAction_TEMPLATE_ACTION_SIGN_IN:
		return smspb.TemplateAction_TEMPLATE_ACTION_SIGN_IN, nil
	case jinmuidpb.TemplateAction_TEMPLATE_ACTION_SIGN_UP:
		return smspb.TemplateAction_TEMPLATE_ACTION_SIGN_UP, nil
	case jinmuidpb.TemplateAction_TEMPLATE_ACTION_RESET_PASSWORD:
		return smspb.TemplateAction_TEMPLATE_ACTION_RESET_PASSWORD, nil
	case jinmuidpb.TemplateAction_TEMPLATE_ACTION_SET_PHONE_NUMBER:
		return smspb.TemplateAction_TEMPLATE_ACTION_SET_PHONE_NUMBER, nil
	case jinmuidpb.TemplateAction_TEMPLATE_ACTION_MODIFY_PHONE_NUMBER:
		return smspb.TemplateAction_TEMPLATE_ACTION_MODIFY_PHONE_NUMBER, nil
	}
	return smspb.TemplateAction_TEMPLATE_ACTION_INVALID, fmt.Errorf("invalid proto template action %d", action)
}
