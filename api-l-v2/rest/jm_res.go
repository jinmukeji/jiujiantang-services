package rest

import (
	"github.com/jinmukeji/gf-api2/pkg/rest"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	"github.com/kataras/iris/v12"
)

// JMRes 金姆资源的获取条件
type JMRes struct {
	ClientID   string `json:"client_id"`
	AppVersion string `json:"app_version"`
	MobileType int32  `json:"mobile_type"`
}

// BaseURL 获取的baseURL
type BaseURL struct {
	EntryURL         string `json:"entry_url"`
	FaqURL           string `json:"faq_url"`
	QuestionnaireURL string `json:"questionnaire_url"`
}

// GetJMResBaseURL 获取金姆资源的baseURL
func (h *v2Handler) GetJMResBaseURL(ctx iris.Context) {
	var jMRes JMRes
	err := ctx.ReadJSON(&jMRes)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err), false)
		return
	}
	req := new(proto.GetJMResBaseUrlRequest)
	req.AppVersion = jMRes.AppVersion
	req.ClientId = jMRes.ClientID
	req.MobileType = proto.MobileType(jMRes.MobileType)
	resp, errGetJMResBaseURL := h.rpcSvc.GetJMResBaseUrl(
		newRPCContext(ctx), req,
	)
	if errGetJMResBaseURL != nil {
		writeRpcInternalError(ctx, errGetJMResBaseURL, false)
		return
	}
	rest.WriteOkJSON(ctx, BaseURL{
		EntryURL:         resp.EntryUrl,
		FaqURL:           resp.FaqUrl,
		QuestionnaireURL: resp.QuestionnaireUrl,
	})
}
