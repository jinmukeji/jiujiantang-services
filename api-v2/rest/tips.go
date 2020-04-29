package rest

import (
	"math/rand"

	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	"github.com/kataras/iris/v12"
)

// Tip 提醒
type Tip struct {
	Content  string `json:"content"`  // 内容
	Duration int    `json:"duration"` // 期间 “显示时间，单位秒”
}

var tips = make([]Tip, 0)

func init() {
	tips = append(tips, Tip{
		Content:  "为了更精确的测量脉搏数据，测量时间可能需要持续4～5分钟，请耐心等待",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "测量过程中请保持心情平静，保持姿态静止，切勿晃动身体和手",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "请保持安静，不要说话，说话将会导致脉搏波变化改变测量结果",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "如果一直无法测量或提示测量异常，可以调整指环佩戴的位置",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "佩戴指环时，红灯应照在指腹位置",
		Duration: 7,
	})
}

// 获取tips
func (h *v2Handler) GetTips(ctx iris.Context) {
	rest.WriteOkJSON(ctx, DisOrder(tips))
}

// DisOrder 随机打乱数组
func DisOrder(arr []Tip) []Tip {
	count := len(arr)
	for index := 0; index < count; index++ {
		start := (int)(rand.Intn(count))
		arr[start], arr[index] = arr[index], arr[start]
	}
	return arr
}
