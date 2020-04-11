package rest

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	"github.com/kataras/iris/v12"
	"github.com/micro/go-micro/v2/client"
)

type recordHistory struct {
	RecordID       int              `json:"record_id"`        // 测量结果记录ID
	C0             int32            `json:"c0"`               // 心包经测量指标
	C1             int32            `json:"c1"`               // 肝经测量指标
	C2             int32            `json:"c2"`               // 肾经测量指标
	C3             int32            `json:"c3"`               // 脾经测量指标
	C4             int32            `json:"c4"`               // 肺经测量指标
	C5             int32            `json:"c5"`               // 胃经测量指标
	C6             int32            `json:"c6"`               // 胆经测量指标
	C7             int32            `json:"c7"`               // 膀胱经测量指标
	Remark         string           `json:"remark"`           // 测量备注
	Finger         int              `json:"finger"`           // 左右手
	WaveData       []int32          `json:"-"`                // 部分测量数据
	AppHeartRate   float64          `json:"app_heart_rate"`   // app测的心率
	CreatedAt      time.Time        `json:"created_at"`       // 创建日期
	RecordType     int32            `json:"record_type"`      // 记录的类型
	Cid            int32            `json:"cid"`              // 上下文ID
	HeartRate      float64          `json:"heart_rate"`       // 服务器返回的心率
	Tags           []GeneralExplain `json:"tags"`             // 智能分析的tags
	Answers        string           `json:"-"`                // 智能分析的答案
	HasPaid        bool             `json:"has_paid"`         // 是否完成支付
	ShowFullReport bool             `json:"show_full_report"` // 是否显示完整报告
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
	openID := ctx.URLParam("open_id")
	// build rpc request
	req := new(proto.SearchHistoryRequest)
	req.OpenId = openID
	// 历史记录 rpc 网络 io 操作繁重，容易超时，后续需要分页
	resp, err := h.rpcSvc.SearchHistory(newRPCContext(ctx), req, client.WithRequestTimeout(time.Second*30))
	if err != nil {
		writeError(ctx, wrapError(ErrRPCInternal, "", err), false)
		return
	}
	recordHistories := make([]recordHistory, len(resp.RecordHistories))
	for idx, record := range resp.RecordHistories {
		createdAt, err := ptypes.Timestamp(record.CreatedTime)
		if err != nil {
			continue
		}
		recordHistories[idx] = recordHistory{
			CreatedAt:      createdAt.UTC(),
			RecordID:       int(record.RecordId),
			C0:             record.C0,
			C1:             record.C1,
			C2:             record.C2,
			C3:             record.C3,
			C4:             record.C4,
			C5:             record.C5,
			C6:             record.C6,
			C7:             record.C7,
			Remark:         record.Remark,
			Finger:         int(record.Finger),
			AppHeartRate:   record.AppHr,
			RecordType:     record.RecordType,
			Cid:            record.RecordId,
			HeartRate:      record.Hr,
			Tags:           MapGeneralExplains(record.Tags),
			HasPaid:        record.HasPaid,
			ShowFullReport: record.ShowFullReport,
		}
	}
	rest.WriteOkJSON(ctx, recordHistories)
}
