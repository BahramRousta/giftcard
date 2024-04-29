package repository

import (
	"giftCard/internal/model"
	"gorm.io/gorm"
)

type ProductRepository struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		db,
	}
}

func (repo *ProductRepository) InsertProduct(product *model.Product) error {
	return repo.DB.Create(product).Error
}
