package rest

import (
	"github.com/jinmukeji/gf-api2/pkg/rest.v3"
	"github.com/kataras/iris/v12"
)

// ClientPreferencesBody 获取资源文件的请求
type ClientPreferencesBody struct {
	ClientID      string `json:"client_id"`
	SecretKeyHash string `json:"secret_key_hash"`
	Seed          string `json:"seed"`
	ClientVersion string `json:"client_version"`
	Environment   string `json:"environment"`
}

// ClientPreferencesResponse 获取资源文件的响应
type ClientPreferencesResponse struct {
	ApiURL       string `json:"api_url"`
	AppLoginURL  string `json:"app_login_url"`
	AppEntryURL  string `json:"app_entry_url"`
	AppFaqURL    string `json:"app_faq_url"`
	AppReportURL string `json:"app_report_url"`
}

// 获取资源文件
func (h *sysHandler) ClientPreferences(ctx iris.Context) {
	var clientPreferencesBody ClientPreferencesBody
	err := ctx.ReadJSON(&clientPreferencesBody)
	if err != nil {
		writeError(ctx, wrapError(ErrParsingRequestFailed, "", err))
		return
	}
	// FIXME: 需要区别不同环境进行验证，目前暂时取消验证
	// req := new(proto.ClientAuthRequest)
	// req.ClientId = clientPreferencesBody.ClientID
	// req.SecretKeyHash = clientPreferencesBody.SecretKeyHash
	// req.Seed = clientPreferencesBody.Seed
	// _, errClientAuth := h.rpcSvc.ClientAuth(
	// 	newRPCContext(ctx), req,
	// )
	// if errClientAuth != nil {
	// 	writeRPCInternalError(ctx, errClientAuth)
	// 	return
	// }

	clientPreference, err := h.clientPreferences.GetClientPreferences(clientPreferencesBody.ClientID, clientPreferencesBody.ClientVersion, clientPreferencesBody.Environment)
	if err != nil {
		writeError(ctx, wrapError(ErrGetClientPreferencesFailed, "", err))
		return
	}
	rest.WriteOkJSON(ctx, ClientPreferencesResponse{
		ApiURL:       clientPreference.ApiURL,
		AppLoginURL:  clientPreference.AppLoginURL,
		AppEntryURL:  clientPreference.AppEntryURL,
		AppFaqURL:    clientPreference.AppFaqURL,
		AppReportURL: clientPreference.AppReportURL,
	})
}
