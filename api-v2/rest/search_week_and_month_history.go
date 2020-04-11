package rest

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	"github.com/kataras/iris/v12"
	"github.com/micro/go-micro/client"
)

// StatData 统计数据
type StatData struct {
	// 经络数据列表
	MeridianList []*CData `json:"meridian_list"`
	// 平均经络值，其中 test_time 这里时间字段返回本周的开始时间
	AverageMeridian *CData `json:"average_meridian"`
	// 疾病统计
	DiseaseStatistics *DiseaseStatistics `json:"disease_statistics"`
	// 异常提示
	ErrorMessage string `json:"error_message"`
	// 统计数据的开始时间
	StartTime time.Time `json:"start_time"`
	// 统计数据的结束时间
	EndTime time.Time `json:"end_time"`

	PhysicalDialectics []string `json:"physical_dialectics"`
}

// StressStatus 应激态状态
type StressStatus map[string]bool

// CData 存放C0-C7的结构体
type CData struct {
	C0       int32     `json:"c0"`
	C1       int32     `json:"c1"`
	C2       int32     `json:"c2"`
	C3       int32     `json:"c3"`
	C4       int32     `json:"c4"`
	C5       int32     `json:"c5"`
	C6       int32     `json:"c6"`
	C7       int32     `json:"c7"`
	RecordId int32     `json:"record_id"`
	TestTime time.Time `json:"test_time"` // 创建日期
}

// DiseaseStatistics 趋势的疾病统计.
type DiseaseStatistics struct {
	// 疾病名称和疾病对应数量的映射
	// 其中疾病名称只能为如下值且含义如下
	// abdominal_pain 腹痛腹泻
	// done_sports 运动
	// drinked_wine 饮酒
	// had_cold 感冒
	// lactation 哺乳期
	// ovulation 排卵期
	// pregnant 怀孕
	// rhinitis_episode 鼻炎发作
	// physiological_period 生理期判断
	// viral_infection 既往病毒感染
	Counters map[string]int32 `json:"counters"`
}

func getTip(dates map[int]bool, cdataLength int, typeStatData TypeStatData) string {
	if typeStatData == WeekStatData {
		if len(dates) < 3 {
			return "近7天内测量天数不足3天，无法查看周报告"
		}
		if cdataLength < 5 {
			return "近7天内测量次数不足5次，无法查看周报告"
		}
	}
	if typeStatData == MonthStatData {
		if len(dates) < 10 {
			return "近30天内测量天数不足10天，无法查看月报告"
		}
		if cdataLength < 15 {
			return "近30天内测量次数不足15次，无法查看月报告"
		}
	}
	return ""
}

// TypeStatData 枚举，周还是月
type TypeStatData int32

const (
	// MonthStatData 月统计数据
	MonthStatData TypeStatData = 0
	// WeekStatData 周统计数据
	WeekStatData TypeStatData = 1
)

// getWeekOrMonthStatData 获取周或月的StatData
func (h *v2Handler) getWeekOrMonthStatData(ctx iris.Context, userID int, startTime time.Time, endTime time.Time, typeStatData TypeStatData) (StatData, error) {
	location, _ := time.LoadLocation("Asia/Shanghai")
	req := new(proto.SearchHistoryRequest)
	req.StartTime, _ = ptypes.TimestampProto(startTime)
	req.EndTime, _ = ptypes.TimestampProto(endTime)
	req.Size = -1
	req.UserId = int32(userID)
	req.Offset = 0
	resp, err := h.rpcSvc.SearchHistory(newRPCContext(ctx), req, client.WithRequestTimeout(time.Second*30))
	if err != nil {
		return StatData{}, err
	}
	arraysLen := 0
	var dates = make(map[int]bool)
	statData := StatData{}
	var physicalDialectic = make([]string, 0)
	for _, record := range resp.RecordHistories {
		createdAt, err := ptypes.Timestamp(record.CreatedTime)
		if err != nil {
			continue
		}
		// 应激态不统计
		if record.HasStressState {
			continue
		}
		dates[createdAt.In(location).Day()] = true
		arraysLen = arraysLen + 1
	}

	tips := getTip(dates, arraysLen, typeStatData)
	if tips != "" {
		statData.ErrorMessage = tips
		statData.MeridianList = []*CData{}
		statData.AverageMeridian = &CData{}
		statData.DiseaseStatistics = &DiseaseStatistics{}
		statData.PhysicalDialectics = []string{}
		return statData, nil
	}
	allCData := make([]*CData, arraysLen)
	idx := 0
	diseaseDiseaseStatistics := make(map[string]int32)
	for _, record := range resp.RecordHistories {
		createdAt, err := ptypes.Timestamp(record.CreatedTime)
		if err != nil {
			continue
		}
		// 应激态
		if record.HasStressState {
			if record.StressStatus["has_abdominal_pain"] {
				diseaseDiseaseStatistics["abdominal_pain"]++
			}
			if record.StressStatus["has_done_sports"] {
				diseaseDiseaseStatistics["done_sports"]++
			}
			if record.StressStatus["has_drinked_wine"] {
				diseaseDiseaseStatistics["drinked_wine"]++
			}
			if record.StressStatus["has_had_cold"] {
				diseaseDiseaseStatistics["had_cold"]++
			}
			if record.StressStatus["has_lactation"] {
				diseaseDiseaseStatistics["lactation"]++
			}
			if record.StressStatus["has_ovulation"] {
				diseaseDiseaseStatistics["ovulation"]++
			}
			if record.StressStatus["has_pregnant"] {
				diseaseDiseaseStatistics["pregnant"]++
			}
			if record.StressStatus["has_physiological_period"] {
				diseaseDiseaseStatistics["physiological_period"]++
			}
			if record.StressStatus["has_rhinitis_episode"] {
				diseaseDiseaseStatistics["rhinitis_episode"]++
			}
			if record.StressStatus["has_viral_infection"] {
				diseaseDiseaseStatistics["viral_infection"]++
			}
			// 应激态跳出循环，不统计
			continue
		}
		physicalDialectic = concatStringArray(physicalDialectic, record.PhysicalDialectics)
		allCData[idx] = &CData{
			C0:       record.C0,
			C1:       record.C1,
			C2:       record.C2,
			C3:       record.C3,
			C4:       record.C4,
			C5:       record.C5,
			C6:       record.C6,
			C7:       record.C7,
			RecordId: record.RecordId,
			TestTime: createdAt,
		}
		idx++
	}

	averageCData := reduceDataOfPtAnalysis(startTime, allCData...)
	statData.MeridianList = allCData
	statData.AverageMeridian = averageCData
	statData.DiseaseStatistics = &DiseaseStatistics{
		Counters: diseaseDiseaseStatistics,
	}
	statData.PhysicalDialectics = physicalDialectic

	statData.StartTime = startTime
	statData.EndTime = endTime
	return statData, nil
}

// reduceDataOfPtAnalysis 求多段脉搏波测量分析的c0-c7平均值
func reduceDataOfPtAnalysis(testingTime time.Time, data ...*CData) *CData {
	if len(data) == 0 {
		return &CData{
			C0: 0,
			C1: 0,
			C2: 0,
			C3: 0,
			C4: 0,
			C5: 0,
			C6: 0,
			C7: 0,
		}
	}
	length := len(data)

	ptAnalysis := &CData{
		C0:       int32(reduceInt(mapInt(data, func(d *CData) int { return int(d.C0) }), sumInt) / length),
		C1:       int32(reduceInt(mapInt(data, func(d *CData) int { return int(d.C1) }), sumInt) / length),
		C2:       int32(reduceInt(mapInt(data, func(d *CData) int { return int(d.C2) }), sumInt) / length),
		C3:       int32(reduceInt(mapInt(data, func(d *CData) int { return int(d.C3) }), sumInt) / length),
		C4:       int32(reduceInt(mapInt(data, func(d *CData) int { return int(d.C4) }), sumInt) / length),
		C5:       int32(reduceInt(mapInt(data, func(d *CData) int { return int(d.C5) }), sumInt) / length),
		C6:       int32(reduceInt(mapInt(data, func(d *CData) int { return int(d.C6) }), sumInt) / length),
		C7:       int32(reduceInt(mapInt(data, func(d *CData) int { return int(d.C7) }), sumInt) / length),
		TestTime: testingTime,
	}
	return ptAnalysis
}

func sumInt(a, b int) int {
	return a + b
}

func reduceInt(s []int, fn func(int, int) int) (r int) {
	for _, elem := range s {
		r = fn(r, elem)
	}
	return r
}

func mapInt(s []*CData, fn func(*CData) int) []int {
	r := make([]int, len(s))
	for i, elem := range s {
		r[i] = fn(elem)
	}
	return r
}

func concatStringArray(s1 []string, s2 []string) []string {
	return append(s1, s2...)
}
