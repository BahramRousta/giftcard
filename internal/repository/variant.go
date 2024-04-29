package repository

import (
	"giftCard/internal/model"
	"gorm.io/gorm"
)

type VariantRepository struct {
	DB *gorm.DB
}

func NewVariantRepository(db *gorm.DB) *VariantRepository {
	return &VariantRepository{
		db,
	}
}

func (repo *VariantRepository) InsertVariant(variant *model.Variant) error {
	return repo.DB.Create(variant).Error
}
