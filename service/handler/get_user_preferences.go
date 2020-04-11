package handler

import (
	"context"

	corepb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
	jinmuidpb "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
)

// GetUserPreferences 得到用户的偏好
func (j *JinmuHealth) GetUserPreferences(ctx context.Context, req *corepb.GetUserPreferencesRequest, resp *corepb.GetUserPreferencesResponse) error {
	reqGetUserPreferences := new(jinmuidpb.GetUserPreferencesRequest)
	reqGetUserPreferences.UserId = req.UserId
	respGetUserPreferences, errGetUserPreferences := j.jinmuidSvc.GetUserPreferences(ctx, reqGetUserPreferences)
	if errGetUserPreferences != nil {
        log.Errorf("failed to get user preferences")
		return errGetUserPreferences
	}
	resp.Preferences = &corepb.Preferences{
		EnableHeartRateChart:              int32(respGetUserPreferences.Preferences.EnableHeartRateChart),
		EnablePulseWaveChart:              int32(respGetUserPreferences.Preferences.EnablePulseWaveChart),
		EnableWarmPrompt:                  int32(respGetUserPreferences.Preferences.EnableWarmPrompt),
		EnableChooseStatus:                int32(respGetUserPreferences.Preferences.EnableChooseStatus),
		EnableConstitutionDifferentiation: int32(respGetUserPreferences.Preferences.EnableConstitutionDifferentiation),
		EnableSyndromeDifferentiation:     int32(respGetUserPreferences.Preferences.EnableSyndromeDifferentiation),
		EnableWesternMedicineAnalysis:     int32(respGetUserPreferences.Preferences.EnableWesternMedicineAnalysis),
		EnableMeridianBarGraph:            int32(respGetUserPreferences.Preferences.EnableMeridianBarGraph),
		EnableComment:                     int32(respGetUserPreferences.Preferences.EnableComment),
		EnableHealthTrending:              int32(respGetUserPreferences.Preferences.EnableHealthTrending),
	}
	return nil
}
