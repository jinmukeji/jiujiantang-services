package wechat

import (
	"encoding/json"
	"fmt"
	"time"

	template "gopkg.in/chanxuehong/wechat.v2/mp/message/template"
)

const (
	// FirstValue 第一个值
	FirstValue = "尊敬的客户您好，您的健康报告已经生成，请及时查看。"
	// ReportType 报告类型
	ReportType = "健康报告"
	// ReportRemark 报告备注
	ReportRemark = "点击详情即可查看"
	// TextColor 文字颜色
	TextColor = "#173177"
	// timeFormat 时间格式
	timeFormat = "2006-01-02 15:04:05"
)

// TemplateData 模块的数据
type TemplateData struct {
	TemplateData Data `json:"data"`
}

// Data 数据的格式
type Data struct {
	First    DataItem `json:"first"`
	Keyword1 DataItem `json:"keyword1"`
	Keyword2 DataItem `json:"keyword2"`
	Keyword3 DataItem `json:"keyword3"`
	Remark   DataItem `json:"remark"`
}

// DataItem Data中每个字段的格式
type DataItem struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

// SendViewReportTemplateMessage 发送查看分析报告模块消息
func (u *Wxmp) SendViewReportTemplateMessage(openID string, recordID int, nickname string, reportAt time.Time) error {
	loc, errLoadLocation := time.LoadLocation("Asia/Shanghai")
	if errLoadLocation != nil {
		return errLoadLocation
	}
	d := &Data{
		First: DataItem{
			Value: FirstValue,
			Color: TextColor,
		},
		Keyword1: DataItem{
			Value: nickname,
			Color: TextColor,
		},
		Keyword2: DataItem{
			Value: reportAt.In(loc).Format(timeFormat),
			Color: TextColor,
		},
		Keyword3: DataItem{
			Value: ReportType,
			Color: TextColor,
		},
		Remark: DataItem{
			Value: ReportRemark,
			Color: TextColor,
		},
	}
	data, errMarshal := json.Marshal(d)
	if errMarshal != nil {
		return errMarshal
	}
	path := fmt.Sprintf("%s/app.html#/analysisreport?record_id=%d", u.Options.WxH5ServerBase, recordID)
	msg := &template.TemplateMessage{
		ToUser:     openID,
		TemplateId: u.Options.WxTemplateID,
		URL:        path,
		Data:       data,
	}

	_, err := template.Send(u.WechatClient, msg)
	if err != nil {
		return err
	}
	return nil
}
