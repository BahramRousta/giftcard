package repository

import (
	"giftCard/internal/model"
	"gorm.io/gorm"
)

type CreateOrderRepository struct {
	DB *gorm.DB
}

func NewCreateOrderRepository(db *gorm.DB) *CreateOrderRepository {
	return &CreateOrderRepository{
		db,
	}
}

func (repo *CreateOrderRepository) InsertOrder(order *model.Order) error {
	return repo.DB.Create(order).Error
}
