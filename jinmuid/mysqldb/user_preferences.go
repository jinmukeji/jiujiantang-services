package mysqldb

import (
	"context"
	"time"
)

// UserPreferences 用户偏好
type UserPreferences struct {
	UserID                            int32      `gorm:"primary_key;column:user_id"`
	EnableHeartRateChart              int32      // 是否开启心率扇形图
	EnablePulseWaveChart              int32      // 是否开启波形图
	EnableWarmPrompt                  int32      // 是否开启温馨提示
	EnableChooseStatus                int32      // 是否开启选择状态
	EnableConstitutionDifferentiation int32      // 是否开启中医体质判读
	EnableSyndromeDifferentiation     int32      // 是否开启中医脏腑判读
	EnableWesternMedicineAnalysis     int32      // 是否开启西医判读
	EnableMeridianBarGraph            int32      // 是否开启柱状图
	EnableComment                     int32      // 是否开启备注
	EnableHealthTrending              int32      // 开启健康趋势
	EnableLocationNotification        int32      // 是否开启本地通知
	CreatedAt                         time.Time  // 创建时间
	UpdatedAt                         time.Time  // 更新时间
	DeletedAt                         *time.Time // 删除时间
}

// TableName 返回 User 所在的表名
func (u UserPreferences) TableName() string {
	return "user_preferences"
}

// GetUserPreferencesByUserID 返回数据库中的用户偏好
func (db *DbClient) GetUserPreferencesByUserID(ctx context.Context, userID int32) (*UserPreferences, error) {
	var u UserPreferences
	if err := db.GetDB(ctx).First(&u, "user_id = ? ", userID).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// CreateUserPreferences 创建UserPreferences
func (db *DbClient) CreateUserPreferences(ctx context.Context, userID int32) error {
	now := time.Now()
	userPreferences := UserPreferences{
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return db.GetDB(ctx).Create(&userPreferences).Error
}
