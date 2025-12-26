package model

import (
	"time"
)

type Order struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int64     `gorm:"index;not null" json:"user_id"`
	OrderNo     string    `gorm:"type:varchar(64);unique;not null" json:"order_no"`
	Amount      int64     `gorm:"not null" json:"amount"` // In cents
	Description string    `gorm:"type:varchar(255)" json:"description"`
	Status      string    `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, paid, failed, cancelled
	PlanID      string    `gorm:"type:varchar(50)" json:"plan_id"`
	PaidAt      *time.Time `json:"paid_at,omitempty"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
