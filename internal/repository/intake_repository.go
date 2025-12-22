package repository

import (
	"miniapp-backend/internal/model"

	"gorm.io/gorm"
)

type IntakeRepository struct {
	db *gorm.DB
}

func NewIntakeRepository(db *gorm.DB) *IntakeRepository {
	return &IntakeRepository{db: db}
}

func (r *IntakeRepository) Create(record *model.IntakeRecord) error {
	return r.db.Create(record).Error
}

func (r *IntakeRepository) Delete(id int64, userID int64) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.IntakeRecord{}).Error
}

func (r *IntakeRepository) GetByDate(userID int64, date string) ([]model.IntakeRecord, error) {
	var records []model.IntakeRecord
	err := r.db.Where("user_id = ? AND date = ?", userID, date).Order("recorded_at desc").Find(&records).Error
	return records, err
}

func (r *IntakeRepository) GetTotalIntakeByDate(userID int64, date string) (int, error) {
	var total int
	err := r.db.Model(&model.IntakeRecord{}).
		Where("user_id = ? AND date = ?", userID, date).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

func (r *IntakeRepository) GetWeeklyStats(userID int64, endDate string) ([]map[string]interface{}, error) {
	// Calculate start date (6 days ago)
	// Note: In a real implementation, you'd want to handle date parsing more robustly
	// and possibly generate a full date series to fill in missing days with 0.
	// For now, we'll just query the database.
	
	var results []map[string]interface{}
	// This is a simplified query. Postgres specific date arithmetic could be used.
	err := r.db.Model(&model.IntakeRecord{}).
		Select("date, SUM(amount) as total").
		Where("user_id = ? AND date > (DATE(?) - INTERVAL '7 days') AND date <= ?", userID, endDate, endDate).
		Group("date").
		Order("date asc").
		Find(&results).Error
	
	return results, err
}
