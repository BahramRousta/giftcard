package repository

import (
	"giftCard/internal/model"
	"gorm.io/gorm"
)

type ExchangeRepository struct {
	DB *gorm.DB
}

func NewExchangeRepository(db *gorm.DB) *ExchangeRepository {
	return &ExchangeRepository{
		db,
	}
}

func (repo *ExchangeRepository) InsertExchange(exchange *model.ExchangeRate) error {
	return repo.DB.Create(exchange).Error
}
