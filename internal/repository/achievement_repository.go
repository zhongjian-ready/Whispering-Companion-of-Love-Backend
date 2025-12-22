package repository

import (
	"miniapp-backend/internal/model"

	"gorm.io/gorm"
)

type AchievementRepository struct {
	db *gorm.DB
}

func NewAchievementRepository(db *gorm.DB) *AchievementRepository {
	return &AchievementRepository{db: db}
}

func (r *AchievementRepository) FindAll() ([]model.Achievement, error) {
	var achievements []model.Achievement
	err := r.db.Find(&achievements).Error
	return achievements, err
}

func (r *AchievementRepository) FindUserAchievements(userID int64) ([]model.UserAchievement, error) {
	var userAchievements []model.UserAchievement
	err := r.db.Preload("Achievement").Where("user_id = ?", userID).Find(&userAchievements).Error
	return userAchievements, err
}

func (r *AchievementRepository) Unlock(userID, achievementID int64) error {
	userAchievement := model.UserAchievement{
		UserID:        userID,
		AchievementID: achievementID,
	}
	return r.db.Create(&userAchievement).Error
}
