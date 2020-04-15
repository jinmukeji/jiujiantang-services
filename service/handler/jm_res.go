package handler

import (
	"context"

	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
)

const (
	// AppVersion1_9 app 1.9版本
	AppVersion1_9 = "1.9"
	// AppVersion2_0 app 2.0版本
	AppVersion2_0 = "2.0"
	// AppVersion2_1 app 2.版本
	AppVersion2_1 = "2.1"
	// FaqURL1_9 FaqURL1_9
	FaqURL1_9 = "https://res-cdn.jinmuhealth.com/app/jinmu_v2/faq/v1_9"
	// EntryURL1_9 EntryURL1_9
	EntryURL1_9 = "https://res-cdn.jinmuhealth.com/app/jinmu_v2/entry/v1_9"
	// FaqURL2_0 FaqURL2_0
	FaqURL2_0 = "https://res-cdn.jinmuhealth.com/app/jinmu_v2/faq/v2_0"
	// EntryURL2_0 EntryURL2_0
	EntryURL2_0 = "https://res-cdn.jinmuhealth.com/app/jinmu_v2/entry/v2_0"
	// EntryURL2_1 EntryURL2_1
	EntryURL2_1 = "https://res-cdn.jinmuhealth.com/app/jinmu_v2_1/entry"
	// FaqURL2_1 FaqURL2_1
	FaqURL2_1 = "https://res-cdn.jinmuhealth.com/app/jinmu_v2_1/faq"
)

// GetJMResBaseUrl 获取喜马把脉资源的baseURL
func (j *JinmuHealth) GetJMResBaseUrl(ctx context.Context, req *proto.GetJMResBaseUrlRequest, resp *proto.GetJMResBaseUrlResponse) error {
	switch req.AppVersion {
	case AppVersion1_9:
		resp.EntryUrl = EntryURL1_9
		resp.FaqUrl = FaqURL1_9
		resp.QuestionnaireUrl = j.wechat.Options.JinmuH5Serverbase
	case AppVersion2_0:
		resp.EntryUrl = EntryURL2_0
		resp.FaqUrl = FaqURL2_0
		resp.QuestionnaireUrl = j.wechat.Options.JinmuH5ServerbaseV2_0
	case AppVersion2_1:
		resp.EntryUrl = EntryURL2_1
		resp.FaqUrl = FaqURL2_1
		resp.QuestionnaireUrl = j.wechat.Options.JinmuH5ServerbaseV2_1
	}
	return nil
}
