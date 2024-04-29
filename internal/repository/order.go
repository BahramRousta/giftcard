package repository

import (
	"errors"
	"giftCard/internal/model"
	"gorm.io/gorm"
)

type OrderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		db,
	}
}

func (repo *OrderRepository) InsertOrder(order *model.Order) error {
	return repo.DB.Create(order).Error
}

func (repo *OrderRepository) GetOrder(orderId string) (*model.Order, error) {
	var order model.Order
	if err := repo.DB.Where("order_id = ?", orderId).First(&order).Error; err != nil {
		return nil, errors.New("order not found")
	}
	return &order, nil
}

func (repo *OrderRepository) UpdateOrder(order *model.Order, newStatus string) error {
	order.Status = newStatus
	if err := repo.DB.Save(order).Error; err != nil {
		return err // Return error if there's a problem saving the order
	}
	return nil
}
