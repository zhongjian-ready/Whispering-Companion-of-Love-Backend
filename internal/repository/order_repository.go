package repository

import (
	"miniapp-backend/internal/model"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) FindByOrderNo(orderNo string) (*model.Order, error) {
	var order model.Order
	err := r.db.Where("order_no = ?", orderNo).First(&order).Error
	return &order, err
}

func (r *OrderRepository) UpdateStatus(orderNo string, status string) error {
	return r.db.Model(&model.Order{}).Where("order_no = ?", orderNo).Update("status", status).Error
}
