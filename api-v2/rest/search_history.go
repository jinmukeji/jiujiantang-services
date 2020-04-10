package rest

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	"github.com/kataras/iris/v12"
	"github.com/micro/go-micro/client"
)

const (
	defaultSize = 20
	maxSize     = 100
)

type recordHistory struct {
	RecordID     int              `json:"record_id"`      // 测量结果记录ID
	C0           int32            `json:"c0"`             // 心包经测量指标
	C1           int32            `json:"c1"`             // 肝经测量指标
	C2           int32            `json:"c2"`             // 肾经测量指标
	C3           int32            `json:"c3"`             // 脾经测量指标
	C4           int32            `json:"c4"`             // 肺经测量指标
	C5           int32            `json:"c5"`             // 胃经测量指标
	C6           int32            `json:"c6"`             // 胆经测量指标
	C7           int32            `json:"c7"`             // 膀胱经测量指标
	Remark       string           `json:"remark"`         // 测量备注
	Finger       int              `json:"finger"`         // 左右手
	WaveData     []int32          `json:"-"`              // 部分测量数据
	AppHeartRate float64          `json:"app_heart_rate"` // app测的心率
	CreatedAt    time.Time        `json:"created_at"`     // 创建日期
	RecordType   int32            `json:"record_type"`    // 记录的类型
	Cid          int32            `json:"cid"`            // 上下文ID
	HeartRate    float64          `json:"heart_rate"`     // 服务器返回的心率
	Tags         []GeneralExplain `json:"tags"`           // 智能分析的tags
	Answers      string           `json:"-"`              // 智能分析的答案
	Labels       []Label          `json:"labels"`         // 标签
}

// Label 标签
type Label struct {
	Text    string  `json:"text"`
	BgColor string  `json:"bg_color"`
	BgAlpha float64 `json:"bg_alpha"`
}

const (
	// StatusNo 状态是否
	StatusNo = 0
	// StatusYes 状态是真
	StatusYes = 1
)

// SearchHistory 获取测量历史记录信息
func (h *v2Handler) SearchHistory(ctx iris.Context) {
	// url params
	userID, _ := ctx.URLParamInt("user_id")
	// build rpc request
	req := new(proto.SearchHistoryRequest)
	req.UserId = int32(userID)
	start, _ := time.Parse(time.RFC3339, ctx.URLParam("start"))
	end, _ := time.Parse(time.RFC3339, ctx.URLParam("end"))
	req.StartTime, _ = ptypes.TimestampProto(start)
	req.EndTime, _ = ptypes.TimestampProto(end)
	offset, _ := ctx.URLParamInt("offset")
	size, _ := ctx.URLParamInt("size")
	// size 没有值默认会是-1
	if size <= 0 {
		size = defaultSize
	}
	if size > 100 {
		size = maxSize
	}
	req.Offset = int32(offset)
	req.Size = int32(size)
	// 历史记录 rpc 网络 io 操作繁重，容易超时，后续需要分页
	resp, err := h.rpcSvc.SearchHistory(newRPCContext(ctx), req, client.WithRequestTimeout(time.Second*30))
	if err != nil {
		writeRPCInternalError(ctx, err, true)
		return
	}
	recordHistories := make([]recordHistory, len(resp.RecordHistories))
	for idx, record := range resp.RecordHistories {
		createdAt, err := ptypes.Timestamp(record.CreatedTime)
		if err != nil {
			continue
		}
		labels := make([]Label, 0)
		if record.StressStatus["has_done_sports"] {
			labels = append(labels, Label{
				Text:    "运动",
				BgColor: "#50A9B5",
				BgAlpha: 1.0,
			})
		}
		if record.StressStatus["has_drinked_wine"] {
			labels = append(labels, Label{
				Text:    "饮酒",
				BgColor: "#50B56B",
				BgAlpha: 1.0,
			})
		}
		if record.StressStatus["has_had_cold"] {
			labels = append(labels, Label{
				Text:    "感冒",
				BgColor: "#B5A650",
				BgAlpha: 1.0,
			})
		}
		if record.StressStatus["has_rhinitis_episode"] {
			labels = append(labels, Label{
				Text:    "鼻炎发作",
				BgColor: "#547EA3",
				BgAlpha: 1.0,
			})
		}
		if record.StressStatus["has_abdominal_pain"] {
			labels = append(labels, Label{
				Text:    "腹痛腹泻",
				BgColor: "#5CA694",
				BgAlpha: 1.0,
			})
		}
		if record.StressStatus["has_viral_infection"] {
			labels = append(labels, Label{
				Text:    "既往临床确诊病毒感染",
				BgColor: "#AF7E4E",
				BgAlpha: 1.0,
			})
		}
		if record.StressStatus["has_physiological_period"] {
			labels = append(labels, Label{
				Text:    "生理周期",
				BgColor: "#AD563F",
				BgAlpha: 1.0,
			})
		}
		if record.StressStatus["has_ovulation"] {
			labels = append(labels, Label{
				Text:    "排卵期",
				BgColor: "#8865B3",
				BgAlpha: 1.0,
			})
		}
		if record.StressStatus["has_pregnant"] {
			labels = append(labels, Label{
				Text:    "怀孕",
				BgColor: "#505BB5",
				BgAlpha: 1.0,
			})
		}
		if record.StressStatus["has_lactation"] {
			labels = append(labels, Label{
				Text:    "哺乳期",
				BgColor: "#B54F8B",
				BgAlpha: 1.0,
			})
		}
		recordHistories[idx] = recordHistory{
			CreatedAt:    createdAt.UTC(),
			RecordID:     int(record.RecordId),
			C0:           record.C0,
			C1:           record.C1,
			C2:           record.C2,
			C3:           record.C3,
			C4:           record.C4,
			C5:           record.C5,
			C6:           record.C6,
			C7:           record.C7,
			Remark:       record.Remark,
			Finger:       int(record.Finger),
			AppHeartRate: record.AppHr,
			RecordType:   record.RecordType,
			Cid:          record.RecordId,
			HeartRate:    record.Hr,
			Tags:         MapGeneralExplains(record.Tags),
			Labels:       labels,
		}
	}
	rest.WriteOkJSON(ctx, recordHistories)
}

// DeleteRecordsBody 删除记录的body
type DeleteRecordsBody struct {
	RecordIDList []int `json:"record_id_list"` // 待删除的record_id
}

// DeleteRecords 批量删除记录
func (h *v2Handler) DeleteRecords(ctx iris.Context) {
	userID, _ := ctx.Params().GetInt("user_id")
	var body DeleteRecordsBody
	err := ctx.ReadJSON(&body)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}
	req := new(proto.DeleteRecordRequest)
	req.UserId = int32(userID)
	for _, item := range body.RecordIDList {
		req.RecordId = int32(item)
		_, errDeleteRecord := h.rpcSvc.DeleteRecord(
			newRPCContext(ctx), req,
		)
		if errDeleteRecord != nil {
			writeRPCInternalError(ctx, errDeleteRecord, false)
			return
		}
	}

	rest.WriteOkJSON(ctx, nil)
}
