package repository

import (
	"giftcard/model"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

type OrderRepositoryParams struct {
	fx.In
	Db *gorm.DB
}

func NewOrderRepository(params OrderRepositoryParams) IOrderRepository {
	return &OrderRepository{
		db: params.Db,
	}
}

func (repo *OrderRepository) InsertOrder(order *model.Order) error {
	return repo.db.Create(order).Error
}

func (repo *OrderRepository) GetOrder(orderId string) (*model.Order, error) {
	var order model.Order
	if err := repo.db.Where("order_id = ?", orderId).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (repo *OrderRepository) UpdateOrder(order *model.Order, newStatus string) error {
	order.Status = newStatus
	if err := repo.db.Save(order).Error; err != nil {
		return err
	}
	return nil
}
