package handler

import (
	"context"
	"errors"
	"fmt"

	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
)

// GetUserPreferences 得到用户的偏好
func (j *JinmuIDService) GetUserPreferences(ctx context.Context, req *proto.GetUserPreferencesRequest, resp *proto.GetUserPreferencesResponse) error {
	token, ok := TokenFromContext(ctx)
	if !ok {
		return NewError(ErrInvalidUser, errors.New("failed to get token from context"))
	}
	userID, err := j.datastore.FindUserIDByToken(ctx, token)
	if err != nil {
		return NewError(ErrUserUnauthorized, fmt.Errorf("failed to find userID by token: %s", err.Error()))
	}
	if userID != req.UserId {
		return NewError(ErrInvalidUser, fmt.Errorf("user %d from request and user %d from token are inconsistent", req.UserId, userID))
	}
	userPreferences, err := j.datastore.GetUserPreferencesByUserID(ctx, req.UserId)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to get user preferences by userID %d: %s", req.UserId, err.Error()))
	}
	resp.Preferences = &proto.Preferences{
		EnableHeartRateChart:              userPreferences.EnableHeartRateChart,
		EnablePulseWaveChart:              userPreferences.EnablePulseWaveChart,
		EnableWarmPrompt:                  userPreferences.EnableWarmPrompt,
		EnableChooseStatus:                userPreferences.EnableChooseStatus,
		EnableConstitutionDifferentiation: userPreferences.EnableConstitutionDifferentiation,
		EnableSyndromeDifferentiation:     userPreferences.EnableSyndromeDifferentiation,
		EnableWesternMedicineAnalysis:     userPreferences.EnableWesternMedicineAnalysis,
		EnableMeridianBarGraph:            userPreferences.EnableMeridianBarGraph,
		EnableComment:                     userPreferences.EnableComment,
		EnableHealthTrending:              userPreferences.EnableHealthTrending,
		EnableLocalNotification:           userPreferences.EnableLocationNotification,
	}
	return nil
}
