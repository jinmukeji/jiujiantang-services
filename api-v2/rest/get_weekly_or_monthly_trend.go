package rest

import (
	"math"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	"github.com/kataras/iris/v12"
	"github.com/micro/go-micro/v2/client"
)

// TrendStatData 统计数据
type TrendStatData struct {
	C0                       []TrendCData `json:"c0"`
	C1                       []TrendCData `json:"c1"`
	C2                       []TrendCData `json:"c2"`
	C3                       []TrendCData `json:"c3"`
	C4                       []TrendCData `json:"c4"`
	C5                       []TrendCData `json:"c5"`
	C6                       []TrendCData `json:"c6"`
	C7                       []TrendCData `json:"c7"`
	StartTime                time.Time    `json:"start_time"`
	EndTime                  time.Time    `json:"end_time"`
	AvgC0                    int32        `json:"c0_average"`
	AvgC1                    int32        `json:"c1_average"`
	AvgC2                    int32        `json:"c2_average"`
	AvgC3                    int32        `json:"c3_average"`
	AvgC4                    int32        `json:"c4_average"`
	AvgC5                    int32        `json:"c5_average"`
	AvgC6                    int32        `json:"c6_average"`
	AvgC7                    int32        `json:"c7_average"`
	AbdominalPainCount       int32        `json:"abdominal_pain_count"`
	DoneSportsCount          int32        `json:"done_sports_count"`
	DrinkedWineCount         int32        `json:"drinked_wine_count"`
	HadColdCount             int32        `json:"had_cold_count"`
	LactationCount           int32        `json:"lactation_count"`
	OvulationCount           int32        `json:"ovulation_count"`
	PregnantCount            int32        `json:"pregnant_count"`
	RhinitisEpisodeCount     int32        `json:"rhinitis_episode_count"`
	PhysiologicalPeriodCount int32        `json:"physiological_period_count"`
	ViralInfectionCount      int32        `json:"viral_infection_count"`
	Tip                      string       `json:"tip"`
	PhysicalDialectics       []string     `json:"physical_dialectics"`
}

// CData 存放C0-C7的结构体
type TrendCData struct {
	RecordID  int32     `json:"record_id"`
	Value     int32     `json:"value"`
	CreatedAt time.Time `json:"created_at"` // 创建日期
}

// SearchWeekHistory 获取周历史记录信息
func (h *v2Handler) SearchWeekHistory(ctx iris.Context) {
	userID, err := ctx.Params().GetInt("user_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	timeZone := getTimeZone(ctx)
	// 周历史记录信息
	location, _ := time.LoadLocation(timeZone)
	now := time.Now().In(location)
	// 从当天的23点59分59秒，往前推7天
	weekStart := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, location).AddDate(0, 0, -7)
	// 从当天的23点59分59秒算结束时间
	weekEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, location).AddDate(0, 0, 0)
	weekStatData, err := h.getWeekOrMonthTrendStatData(ctx, userID, weekStart.UTC(), weekEnd.UTC(), WeekStatData)
	if err != nil {
		writeError(ctx, wrapError(ErrRPCInternal, "", err), false)
		return
	}
	rest.WriteOkJSON(ctx, weekStatData)
}

// SearchMonthHistory 获取月历史记录信息
func (h *v2Handler) SearchMonthHistory(ctx iris.Context) {
	userID, err := ctx.Params().GetInt("user_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	timeZone := getTimeZone(ctx)
	// 周历史记录信息
	location, _ := time.LoadLocation(timeZone)
	now := time.Now().In(location)
	// 从当天的23点59分59秒，往前推30天
	monthStart := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, location).AddDate(0, 0, -30)
	// 从当天的23点59分59秒算结束时间
	monthEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, location).AddDate(0, 0, 0)
	monthStatData, err := h.getWeekOrMonthTrendStatData(ctx, userID, monthStart.UTC(), monthEnd.UTC(), MonthStatData)
	if err != nil {
		writeError(ctx, wrapError(ErrRPCInternal, "", err), false)
		return
	}
	rest.WriteOkJSON(ctx, monthStatData)
}

// getWeekOrMonthTrendStatData 获取周或月的StatData
func (h *v2Handler) getWeekOrMonthTrendStatData(ctx iris.Context, userID int, start time.Time, end time.Time, typeStatData TypeStatData) (TrendStatData, error) {
	location, _ := time.LoadLocation("Asia/Shanghai")
	req := new(proto.SearchHistoryRequest)
	req.StartTime, _ = ptypes.TimestampProto(start)
	req.EndTime, _ = ptypes.TimestampProto(end)
	req.Size = -1
	req.UserId = int32(userID)
	req.Offset = 0
	resp, err := h.rpcSvc.SearchHistory(newRPCContext(ctx), req, client.WithRequestTimeout(time.Second*30))
	if err != nil {
		return TrendStatData{}, err
	}
	arraysLen := 0
	for _, record := range resp.RecordHistories {
		_, err := ptypes.Timestamp(record.CreatedTime)
		if err != nil {
			continue
		}
		if record.HasStressState {
			continue
		}
		arraysLen = arraysLen + 1
	}
	statData := TrendStatData{
		C0:        make([]TrendCData, arraysLen),
		C1:        make([]TrendCData, arraysLen),
		C2:        make([]TrendCData, arraysLen),
		C3:        make([]TrendCData, arraysLen),
		C4:        make([]TrendCData, arraysLen),
		C5:        make([]TrendCData, arraysLen),
		C6:        make([]TrendCData, arraysLen),
		C7:        make([]TrendCData, arraysLen),
		StartTime: start,
		EndTime:   end,
	}

	var abdominalPainCount int32
	var doneSportsCount int32
	var drinkedWineCount int32
	var hadColdCount int32
	var lactationCount int32
	var ovulationCount int32
	var pregnantCount int32
	var rhinitisEpisodeCount int32
	var physiologicalPeriodCount int32
	var viralInfectionCount int32
	var dates = make(map[int]bool)
	var count = 0
	var physicalDialectic = make([]string, 0)
	for _, record := range resp.RecordHistories {
		createdAt, err := ptypes.Timestamp(record.CreatedTime)
		if err != nil {
			continue
		}
		if record.HasStressState {
			if record.StressStatus["has_abdominal_pain"] {
				abdominalPainCount = abdominalPainCount + 1
			}
			if record.StressStatus["has_done_sports"] {
				doneSportsCount = doneSportsCount + 1
			}
			if record.StressStatus["has_drinked_wine"] {
				drinkedWineCount = drinkedWineCount + 1
			}
			if record.StressStatus["has_had_cold"] {
				hadColdCount = hadColdCount + 1
			}
			if record.StressStatus["has_lactation"] {
				lactationCount = lactationCount + 1
			}
			if record.StressStatus["has_ovulation"] {
				ovulationCount = ovulationCount + 1
			}
			if record.StressStatus["has_pregnant"] {
				pregnantCount = pregnantCount + 1
			}
			if record.StressStatus["has_physiological_period"] {
				physiologicalPeriodCount = physiologicalPeriodCount + 1
			}
			if record.StressStatus["has_rhinitis_episode"] {
				rhinitisEpisodeCount = rhinitisEpisodeCount + 1
			}
			if record.StressStatus["has_viral_infection"] {
				viralInfectionCount = viralInfectionCount + 1
			}
			// 应激态跳出循环
			continue
		}
		dates[createdAt.In(location).Day()] = true
		statData = setCToStatData(statData, arraysLen, count, record)
		physicalDialectic = concatStringArray(physicalDialectic, record.PhysicalDialectics)
		count++
	}
	statData.Tip = getTrendTip(dates, statData.C0, typeStatData)
	statData.AbdominalPainCount = abdominalPainCount
	statData.DoneSportsCount = doneSportsCount
	statData.DrinkedWineCount = drinkedWineCount
	statData.HadColdCount = hadColdCount
	statData.LactationCount = lactationCount
	statData.OvulationCount = ovulationCount
	statData.PhysiologicalPeriodCount = physiologicalPeriodCount
	statData.PregnantCount = pregnantCount
	statData.RhinitisEpisodeCount = rhinitisEpisodeCount
	statData.ViralInfectionCount = viralInfectionCount
	statData.AvgC0 = getTrendAverage(statData.C0)
	statData.AvgC1 = getTrendAverage(statData.C1)
	statData.AvgC2 = getTrendAverage(statData.C2)
	statData.AvgC3 = getTrendAverage(statData.C3)
	statData.AvgC4 = getTrendAverage(statData.C4)
	statData.AvgC5 = getTrendAverage(statData.C5)
	statData.AvgC6 = getTrendAverage(statData.C6)
	statData.AvgC7 = getTrendAverage(statData.C7)
	statData.PhysicalDialectics = physicalDialectic
	return statData, nil
}

func getAverage(cDatas []TrendCData) int32 {
	if len(cDatas) == 0 {
		return int32(0)
	}
	var sum int
	for _, cData := range cDatas {
		sum = sum + int(cData.Value)
	}
	average := float64(sum / len(cDatas))
	return int32(math.Floor(average + 0.5))
}

func setCToStatData(statData TrendStatData, length int, idx int, record *proto.RecordHistory) TrendStatData {
	createdAt, _ := ptypes.Timestamp(record.CreatedTime)
	statData.C0[length-idx-1] = TrendCData{
		RecordID:  record.RecordId,
		Value:     record.C0,
		CreatedAt: createdAt,
	}
	statData.C1[length-idx-1] = TrendCData{
		RecordID:  record.RecordId,
		Value:     record.C1,
		CreatedAt: createdAt,
	}
	statData.C2[length-idx-1] = TrendCData{
		RecordID:  record.RecordId,
		Value:     record.C2,
		CreatedAt: createdAt,
	}
	statData.C3[length-idx-1] = TrendCData{
		RecordID:  record.RecordId,
		Value:     record.C3,
		CreatedAt: createdAt,
	}
	statData.C4[length-idx-1] = TrendCData{
		RecordID:  record.RecordId,
		Value:     record.C4,
		CreatedAt: createdAt,
	}
	statData.C5[length-idx-1] = TrendCData{
		RecordID:  record.RecordId,
		Value:     record.C5,
		CreatedAt: createdAt,
	}
	statData.C6[length-idx-1] = TrendCData{
		RecordID:  record.RecordId,
		Value:     record.C6,
		CreatedAt: createdAt,
	}
	statData.C7[length-idx-1] = TrendCData{
		RecordID:  record.RecordId,
		Value:     record.C7,
		CreatedAt: createdAt,
	}
	return statData
}

func getTrendAverage(cDatas []TrendCData) int32 {
	if len(cDatas) == 0 {
		return int32(0)
	}
	var sum int
	for _, cData := range cDatas {
		sum = sum + int(cData.Value)
	}
	average := float64(sum / len(cDatas))
	return int32(math.Floor(average + 0.5))
}

func getTrendTip(dates map[int]bool, c []TrendCData, typeStatData TypeStatData) string {
	if typeStatData == WeekStatData {
		if len(dates) < 3 {
			return "近7天内测量天数不足3天，无法查看周报告"
		}
		if len(c) < 5 {
			return "近7天内测量次数不足5次，无法查看周报告"
		}
	}
	if typeStatData == MonthStatData {
		if len(dates) < 10 {
			return "近30天内测量天数不足10天，无法查看月报告"
		}
		if len(c) < 15 {
			return "近30天内测量次数不足15次，无法查看月报告"
		}
	}
	return ""
}
