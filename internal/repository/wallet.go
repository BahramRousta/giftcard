package repository

import (
	"giftCard/internal/model"
	"gorm.io/gorm"
)

type WalletRepository struct {
	DB *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{
		db,
	}
}

func (repo *WalletRepository) InsertWallet(wallet *model.Wallet) error {
	return repo.DB.Create(wallet).Error
}
