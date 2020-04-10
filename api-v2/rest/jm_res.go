package rest

import (
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	"github.com/kataras/iris/v12"
)

// JMRes 喜马把脉资源的获取条件
type JMRes struct {
	AppVersion string `json:"app_version"`
}

// BaseURL 获取的baseURL
type BaseURL struct {
	EntryURL         string `json:"entry_url"`
	FaqURL           string `json:"faq_url"`
	QuestionnaireURL string `json:"questionnaire_url"`
}

// GetJMResBaseURL 获取喜马把脉资源的baseURL
func (h *v2Handler) GetJMResBaseURL(ctx iris.Context) {
	var jMRes JMRes
	err := ctx.ReadJSON(&jMRes)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}
	req := new(proto.GetJMResBaseUrlRequest)
	req.AppVersion = jMRes.AppVersion
	resp, errGetJMResBaseURL := h.rpcSvc.GetJMResBaseUrl(
		newRPCContext(ctx), req,
	)
	if errGetJMResBaseURL != nil {
		writeRPCInternalError(ctx, errGetJMResBaseURL, false)
		return
	}
	rest.WriteOkJSON(ctx, BaseURL{
		EntryURL:         resp.EntryUrl,
		FaqURL:           resp.FaqUrl,
		QuestionnaireURL: resp.QuestionnaireUrl,
	})
}
