package model

import (
	"time"
)

type IntakeRecord struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int64     `gorm:"index;not null" json:"user_id"`
	Amount     int       `gorm:"not null" json:"amount"` // ml
	RecordedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"recorded_at"`
	Date       string    `gorm:"type:date;index" json:"date"` // YYYY-MM-DD for easy grouping
}

func (IntakeRecord) TableName() string {
	return "intake_records"
}
