package rest

import (
	"fmt"
	"time"

	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	corepb "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	generalpb "github.com/jinmukeji/proto/gen/micro/idl/ptypes/v2"
	"github.com/kataras/iris/v12"
)

// Schedule 推送时间表
type Schedule struct {
	EventHappenAt string          `json:"event_happen_at"` // 推送时间
	TimeZone      string          `json:"timezone"`        // 时区信息
	Repeat        *RepeatSchedule `json:"repeat"`          // 重复推送规律
}

// Frequency 标记推送时间间隔的基本单位
type Frequency int

const (
	// FrequencyDaily 以天为基本单位
	FrequencyDaily Frequency = iota
	// FrequencyWeekly 以周为基本单位
	FrequencyWeekly
	// FrequencyMonthly 以月为基本单位
	FrequencyMonthly
)

// RepeatSchedule 重复推送的规律
type RepeatSchedule struct {
	Frequency Frequency `json:"frequency"` // 推送时间间隔基本单位
	Interval  int32     `json:"interval"`  // 推送时间间隔

	HasWeekdays  bool           `json:"has_weekdays"`   // 是否需要以周为基本单位
	Weekdays     []time.Weekday `json:"weekdays"`       // 一周内有哪些天推送
	HasMonthdays bool           `json:"has_month_days"` // 是否需要以月为基本单位
	MonthDays    []int32        `json:"month_days"`     // 一个月内有哪些天推送

	MaxNotificationTimes *int32  `json:"max_notification_times"` // 最大推送次数
	EndAt                *string `json:"end_at"`                 // 推送结束时间
}

// LocalNotification 本地消息推送
type LocalNotification struct {
	Title    string   `json:"title"`    // 推送标题
	Content  string   `json:"content"`  // 推送内容
	Schedule Schedule `json:"schedule"` // 推送时间表
}

func mapToWeekdays(protoweekdays []generalpb.Weekday) []time.Weekday {
	var mapping = map[generalpb.Weekday]time.Weekday{
		generalpb.Weekday_WEEKDAY_SUNDAY:    time.Sunday,
		generalpb.Weekday_WEEKDAY_MONDAY:    time.Monday,
		generalpb.Weekday_WEEKDAY_TUESDAY:   time.Tuesday,
		generalpb.Weekday_WEEKDAY_WEDNESDAY: time.Wednesday,
		generalpb.Weekday_WEEKDAY_THURSDAY:  time.Thursday,
		generalpb.Weekday_WEEKDAY_FRIDAY:    time.Friday,
		generalpb.Weekday_WEEKDAY_SATURDAY:  time.Saturday,
	}
	timeweekdays := make([]time.Weekday, len(protoweekdays))
	for idx, item := range protoweekdays {
		timeweekdays[idx] = mapping[item]
	}
	return timeweekdays
}

// GetLocalNotifications 获取本地推送消息
func (h *v2Handler) GetLocalNotifications(ctx iris.Context) {
	req := new(corepb.GetLocalNotificationsRequest)
	resp, err := h.rpcSvc.GetLocalNotifications(
		newRPCContext(ctx), req,
	)
	if err != nil {
		writeRPCInternalError(ctx, err, true)
		return
	}

	respData := make([]LocalNotification, len(resp.LocalNotifications))
	for idx, notification := range resp.LocalNotifications {
		respData[idx] = LocalNotification{
			Title:   notification.Title,
			Content: notification.Content,
			Schedule: Schedule{
				EventHappenAt: notification.Schedule.EventHappenAt,
				TimeZone:      notification.Schedule.Timezone,
			},
		}
		if notification.Schedule.Repeat == nil {
			respData[idx].Schedule.Repeat = nil
		} else {
			frequency, errmapProtoFrequencyToRest := mapProtoFrequencyToRest(notification.Schedule.Repeat.Frequency)
			if errmapProtoFrequencyToRest != nil {
				writeError(ctx, wrapError(ErrInvalidValue, "", errmapProtoFrequencyToRest), false)
				return
			}
			respData[idx].Schedule.Repeat = &RepeatSchedule{
				Frequency:    frequency,
				Interval:     notification.Schedule.Repeat.Interval,
				HasWeekdays:  notification.Schedule.Repeat.HasWeekdays,
				HasMonthdays: notification.Schedule.Repeat.HasMonthDays,
				MonthDays:    notification.Schedule.Repeat.MonthDays,
			}
			if notification.Schedule.Repeat.HasWeekdays {
				respData[idx].Schedule.Repeat.Weekdays = mapToWeekdays(notification.Schedule.Repeat.Weekdays)
			}
			if !notification.Schedule.Repeat.HasMaxNotificationTimes {
				respData[idx].Schedule.Repeat.MaxNotificationTimes = nil
			} else {
				respData[idx].Schedule.Repeat.MaxNotificationTimes = &notification.Schedule.Repeat.MaxNotificationTimes
			}
			if !notification.Schedule.Repeat.HasEndAt {
				respData[idx].Schedule.Repeat.EndAt = nil
			} else {
				respData[idx].Schedule.Repeat.EndAt = &notification.Schedule.Repeat.EndAt
			}
		}
	}
	rest.WriteOkJSON(ctx, respData)
}

func mapProtoFrequencyToRest(frequency corepb.Frequency) (Frequency, error) {
	switch frequency {
	case corepb.Frequency_FREQUENCY_DAILY:
		return FrequencyDaily, nil
	case corepb.Frequency_FREQUENCY_WEEKLY:
		return FrequencyWeekly, nil
	case corepb.Frequency_FREQUENCY_MONTHLY:
		return FrequencyMonthly, nil
	case corepb.Frequency_FREQUENCY_INVALID:
		return FrequencyDaily, fmt.Errorf("invalid proto frequency %d", corepb.Frequency_FREQUENCY_INVALID)
	case corepb.Frequency_FREQUENCY_UNSET:
		return FrequencyDaily, fmt.Errorf("invalid proto frequency %d", corepb.Frequency_FREQUENCY_UNSET)
	}
	return FrequencyDaily, fmt.Errorf("invalid proto frequency")
}
