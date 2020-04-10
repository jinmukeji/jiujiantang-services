package rest

import (
	"github.com/jinmukeji/gf-api2/pkg/rest"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
	"github.com/kataras/iris/v12"
)

// UserPreferences 用户偏好
type UserPreferences struct {
	EnableHeartRateChart              bool `json:"enable_heart_rate_chart"`
	EnablePulseWaveChart              bool `json:"enable_pulse_wave_chart"`
	EnableWarmPrompt                  bool `json:"enable_warm_prompt"`
	EnableChooseStatus                bool `json:"enable_choose_status"`
	EnableConstitutionDifferentiation bool `json:"enable_constitution_differentiation"`
	EnableSyndromeDifferentiation     bool `json:"enable_syndrome_differentiation"`
	EnableWesternMedicineAnalysis     bool `json:"enable_western_medicine_analysis"`
	EnableMeridianBarGraph            bool `json:"enable_meridian_bar_graph"`
	EnableComment                     bool `json:"enable_comment"`
	EnableHealthTrending              bool `json:"enable_health_trending"`
	EnableLocalNotification           bool `json:"enable_local_notification"`
}

// 获取用户首选项配置
func (h *webHandler) GetUserPreferences(ctx iris.Context) {
	userID, err := ctx.Params().GetInt("user_id")
	if err != nil {
		writeError(ctx, wrapError(ErrInvalidValue, "", err), false)
		return
	}
	req := new(proto.GetUserPreferencesRequest)
	req.UserId = int32(userID)
	resp, err := h.rpcSvc.GetUserPreferences(
		newRPCContext(ctx), req,
	)
	if err != nil {
		writeRpcInternalError(ctx, err, false)
		return
	}

	rest.WriteOkJSON(ctx, UserPreferences{
		EnableHeartRateChart:              resp.Preferences.EnableHeartRateChart != 0,
		EnablePulseWaveChart:              resp.Preferences.EnablePulseWaveChart != 0,
		EnableWarmPrompt:                  resp.Preferences.EnableWarmPrompt != 0,
		EnableChooseStatus:                resp.Preferences.EnableChooseStatus != 0,
		EnableConstitutionDifferentiation: resp.Preferences.EnableConstitutionDifferentiation != 0,
		EnableSyndromeDifferentiation:     resp.Preferences.EnableSyndromeDifferentiation != 0,
		EnableWesternMedicineAnalysis:     resp.Preferences.EnableWesternMedicineAnalysis != 0,
		EnableMeridianBarGraph:            resp.Preferences.EnableMeridianBarGraph != 0,
		EnableComment:                     resp.Preferences.EnableComment != 0,
		EnableHealthTrending:              resp.Preferences.EnableHealthTrending != 0,
		EnableLocalNotification:           resp.Preferences.EnableLocalNotification != 0,
	},
	)
}
