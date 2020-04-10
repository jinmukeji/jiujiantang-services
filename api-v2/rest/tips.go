package rest

import (
	"math/rand"

	"github.com/jinmukeji/gf-api2/pkg/rest"
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
		Content:  "检测前，请在舒适的环境中休息至少5分钟",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "调整呼吸，放松心情",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "避免焦虑与激动",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "测量时不要移动手臂和身体",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "测量时不要说话",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "未满16周岁，解读结果可能不准确",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "涂红色指甲油可能导致无法测量",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "手指有厚茧可能导致无法正常测量",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "手指过冷可能导致无法正常测量",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "本身脉搏波微弱可能导致无法正常测量",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "女性特殊时期测脉，通常为气血异常",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "测量过程中双脚自然落地，不得翘腿，抖动",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "左右手均测脉，一般相差不超过正负10%，若相差较大属于正常现象，即表示两侧经络情况不一致",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "用眼一段时间要注意休息，眨一眨眼睛",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "最佳睡眠时间是在晚上10点-清晨6点",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "早上醒来，先喝一杯水，预防结石",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "睡前三小时不要吃东西，会胖",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "远离充电座，人体应远离30公分以上，切忌放在床边",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "每天十杯水，膀胱癌不会来",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "白天多喝水，晚上少喝水",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "揉搓耳朵，刺激耳朵的穴位，能让脑袋更灵活",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "10种吃了会快乐的食物：深海鱼，香蕉，葡萄柚，全麦面包，菠菜，大蒜，南瓜，低脂牛奶，鸡肉，樱桃",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "帮助头发生长：多食用包心菜，蛋，豆类；少吃甜食（尤其是果糖）",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "每天一杯柠檬汁、橙汁，美白又淡化黑斑",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "女性不宜喝茶的五个时期：月经期，怀孕，临产前，产后，更年期",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "减少食用盐腌、烟熏、烧烤的食物",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "每天摄取新鲜的蔬菜与水果",
		Duration: 7,
	})

	tips = append(tips, Tip{
		Content:  "每天摄取富含高纤维的五谷类及豆类",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "每天摄取均衡的饮食，不过量",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "皮肤干燥多吃胡萝卜",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "多做思考，正是训练大脑的良方",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "多晒太阳补充体内维生素D，这对骨骼健康大有好处",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "生姜的辣味能使身体从内部生热，增强免疫功能",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "酸奶能保护胃黏膜、钙含量丰富",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "每周散步四到五次，每次30到40分钟，对身体非常有益",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "爬楼梯是种非常好的锻炼形式，对心血管有益",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "冥想可缓解紧张，治疗心脏病、关节炎等疾病",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "放声大笑能带动身体80多块肌肉活动，并释放感到幸福的激素",
		Duration: 7,
	})
	tips = append(tips, Tip{
		Content:  "维持理想体重，不过胖",
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
