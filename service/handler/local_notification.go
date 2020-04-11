package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	valid "github.com/asaskevich/govalidator"
	"github.com/jinmukeji/jiujiantang-services/service/mysqldb"
	corepb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	generalpb "github.com/jinmukeji/proto/v3/gen/micro/idl/ptypes/v2"
)

// CreateLocalNotification 新建本地消息推送
func (j *JinmuHealth) CreateLocalNotification(ctx context.Context, req *corepb.CreateLocalNotificationRequest, repl *corepb.CreateLocalNotificationResponse) error {
	if ok, err := validateCreateLocalNotificationRequest(req); !ok || err != nil {
		return NewError(ErrLocalNotification, errors.New("none-empty schedule required"))
	}
	notice := &mysqldb.LocalNotification{
		Title:         req.LocalNotification.Title,
		Content:       req.LocalNotification.Content,
		EventHappenAt: req.LocalNotification.Schedule.EventHappenAt,
		Timezone:      req.LocalNotification.Schedule.Timezone,
	}
	if req.LocalNotification.Schedule.Repeat != nil {
		mysqlFrequency, errmapProtoFrequencyToDB := mapProtoFrequencyToDB(req.LocalNotification.Schedule.Repeat.Frequency)
		if errmapProtoFrequencyToDB != nil {
			return errmapProtoFrequencyToDB
		}
		notice.Frequency = mysqlFrequency
		notice.Interval = req.LocalNotification.Schedule.Repeat.Interval
		notice.HasWeekdays = req.LocalNotification.Schedule.Repeat.HasWeekdays
		notice.HasMonthDays = req.LocalNotification.Schedule.Repeat.HasMonthDays
		notice.MaxNotificationTimes = ifHasMaxNotificationTimes(req)
		notice.EndAt = ifHasEndAt(req)
	}
	weekdays := req.LocalNotification.Schedule.Repeat.Weekdays
	if len(weekdays) == 0 {
		notice.Weekdays = ""
	} else {
		json, errMarshal := json.Marshal(&weekdays)
		if errMarshal != nil {
			return fmt.Errorf("failed to marshal weekdays: %s", errMarshal.Error())
		}
		notice.Weekdays = string(json)
	}

	monthdays := req.LocalNotification.Schedule.Repeat.MonthDays
	if len(monthdays) == 0 {
		notice.MonthDays = ""
	} else {
		json, errMarshal := json.Marshal(&monthdays)
		if errMarshal != nil {
			return fmt.Errorf("failed to marshal monthdays: %s", errMarshal.Error())
		}
		notice.MonthDays = string(json)
	}

	if err := j.datastore.CreateLocalNotification(ctx, notice); err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to create local notification: %s", err.Error()))
	}
	return nil
}

// GetLocalNotifications 获取本地消息推送
func (j *JinmuHealth) GetLocalNotifications(ctx context.Context, req *corepb.GetLocalNotificationsRequest, repl *corepb.GetLocalNotificationsResponse) error {
	localNotification, err := j.datastore.GetLocalNotifications(ctx)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to get local notification: %s", err.Error()))
	}
	items := make([]*corepb.LocalNotification, len(localNotification))
	for idx, item := range localNotification {
		items[idx] = &corepb.LocalNotification{
			Title:    item.Title,
			Content:  item.Content,
			Schedule: &corepb.Schedule{},
		}
		items[idx].Schedule.EventHappenAt = localNotification[idx].EventHappenAt
		items[idx].Schedule.Timezone = localNotification[idx].Timezone

		if localNotification[idx].Interval == 0 && !localNotification[idx].HasWeekdays && !localNotification[idx].HasMonthDays {
			items[idx].Schedule.Repeat = nil
		} else {
			items[idx].Schedule.Repeat = &corepb.RepeatSchedule{}
			protoFrequency, errmapDBFrequencyToProto := mapDBFrequencyToProto(localNotification[idx].Frequency)
			if errmapDBFrequencyToProto != nil {
				return errmapDBFrequencyToProto
			}
			items[idx].Schedule.Repeat.Frequency = protoFrequency
			items[idx].Schedule.Repeat.Interval = item.Interval
			items[idx].Schedule.Repeat.HasWeekdays = item.HasWeekdays
			stb := []generalpb.Weekday{}
			if localNotification[idx].HasWeekdays {
				err = json.Unmarshal([]byte(localNotification[idx].Weekdays), &stb)
				if err != nil {
					return fmt.Errorf("failed to unmarshal weekdays: %s", err.Error())
				}
				items[idx].Schedule.Repeat.Weekdays = stb
			}
			items[idx].Schedule.Repeat.HasMonthDays = item.HasMonthDays
			if localNotification[idx].HasMonthDays {
				items[idx].Schedule.Repeat.MonthDays, err = mapFromStringToIntArray(localNotification[idx].MonthDays)
				if err != nil {
					return fmt.Errorf("failed to unmarshal monthdays: %s", err.Error())
				}
			}
			if item.MaxNotificationTimes != nil {
				items[idx].Schedule.Repeat.MaxNotificationTimes = *item.MaxNotificationTimes
			}
			if item.EndAt != nil {
				items[idx].Schedule.Repeat.EndAt = *item.EndAt
			}
		}
	}
	repl.LocalNotifications = items
	return nil
}

func mapProtoFrequencyToDB(frequency corepb.Frequency) (mysqldb.Frequency, error) {
	switch frequency {
	case corepb.Frequency_FREQUENCY_DAILY:
		return mysqldb.FrequencyDaily, nil
	case corepb.Frequency_FREQUENCY_WEEKLY:
		return mysqldb.FrequencyWeekly, nil
	case corepb.Frequency_FREQUENCY_MONTHLY:
		return mysqldb.FrequencyMonthly, nil
	case corepb.Frequency_FREQUENCY_INVALID:
		return mysqldb.FrequencyDaily, fmt.Errorf("invalid proto frequency %d", corepb.Frequency_FREQUENCY_INVALID)
	case corepb.Frequency_FREQUENCY_UNSET:
		return mysqldb.FrequencyDaily, fmt.Errorf("invalid proto frequency %d", corepb.Frequency_FREQUENCY_UNSET)
	}
	return mysqldb.FrequencyDaily, fmt.Errorf("invalid proto frequency")
}

func mapDBFrequencyToProto(frequency mysqldb.Frequency) (corepb.Frequency, error) {
	switch frequency {
	case mysqldb.FrequencyDaily:
		return corepb.Frequency_FREQUENCY_DAILY, nil
	case mysqldb.FrequencyWeekly:
		return corepb.Frequency_FREQUENCY_WEEKLY, nil
	case mysqldb.FrequencyMonthly:
		return corepb.Frequency_FREQUENCY_MONTHLY, nil
	}
	return corepb.Frequency_FREQUENCY_INVALID, fmt.Errorf("invalid mysql frequency %s", frequency)
}

// mapFromStringToIntArray json字符串转int32数组
func mapFromStringToIntArray(jsonStr string) ([]int32, error) {
	stb := []int32{}
	err := json.Unmarshal([]byte(jsonStr), &stb)
	if err != nil {
		return stb, err
	}
	return stb, nil
}

// ifHasMaxNotificationTimes 获取最大推送次数
func ifHasMaxNotificationTimes(req *corepb.CreateLocalNotificationRequest) *int32 {
	if !req.LocalNotification.Schedule.Repeat.HasMaxNotificationTimes {
		return nil
	}
	return &req.LocalNotification.Schedule.Repeat.MaxNotificationTimes
}

// ifHasEndAt 获取结束时间
func ifHasEndAt(req *corepb.CreateLocalNotificationRequest) *string {
	if !req.LocalNotification.Schedule.Repeat.HasEndAt {
		return nil
	}
	return &req.LocalNotification.Schedule.Repeat.EndAt
}

func validateCreateLocalNotificationRequest(req *corepb.CreateLocalNotificationRequest) (bool, error) {
	if valid.IsNull(req.LocalNotification.Schedule.EventHappenAt) {
		return false, nil
	}
	return true, nil
}
