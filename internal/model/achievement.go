package model

import (
	"time"
)

type Achievement struct {
	ID            int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string `gorm:"type:varchar(100)" json:"name"`
	Description   string `gorm:"type:varchar(255)" json:"description"`
	IconURL       string `gorm:"type:varchar(500)" json:"icon_url"`
	ConditionType string `gorm:"type:varchar(50)" json:"condition_type"` // e.g., "total_intake", "streak_days"
	ConditionVal  int    `json:"condition_value"`
}

func (Achievement) TableName() string {
	return "achievements"
}

type UserAchievement struct {
	ID            int64       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        int64       `gorm:"index;uniqueIndex:idx_user_achievement" json:"user_id"`
	AchievementID int64       `gorm:"index;uniqueIndex:idx_user_achievement" json:"achievement_id"`
	Achievement   Achievement `gorm:"foreignKey:AchievementID" json:"achievement"`
	UnlockedAt    time.Time   `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"unlocked_at"`
}

func (UserAchievement) TableName() string {
	return "user_achievements"
}
