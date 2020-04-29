package rest

import (
	"github.com/jinmukeji/jiujiantang-services/pkg/rest"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	"github.com/kataras/iris/v12"
)

// 获取用户首选项配置
func (h *v2Handler) OwnerGetUserPreferences(ctx iris.Context) {
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
		writeRPCInternalError(ctx, err, false)
		return
	}

	// FIXME: 使用强类型替代 iris.Map
	rest.WriteOkJSON(ctx, iris.Map{
		"enable_heart_rate_chart":             resp.Preferences.EnableHeartRateChart != 0,
		"enable_pulse_wave_chart":             resp.Preferences.EnablePulseWaveChart != 0,
		"enable_warm_prompt":                  resp.Preferences.EnableWarmPrompt != 0,
		"enable_choose_status":                resp.Preferences.EnableChooseStatus != 0,
		"enable_constitution_differentiation": resp.Preferences.EnableConstitutionDifferentiation != 0,
		"enable_syndrome_differentiation":     resp.Preferences.EnableSyndromeDifferentiation != 0,
		"enable_western_medicine_analysis":    resp.Preferences.EnableWesternMedicineAnalysis != 0,
		"enable_meridian_bar_graph":           resp.Preferences.EnableMeridianBarGraph != 0,
		"enable_comment":                      resp.Preferences.EnableComment != 0,
		"enable_health_trending":              resp.Preferences.EnableHealthTrending != 0,
	},
	)
}
