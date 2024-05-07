package repository

import (
	"giftcard/model"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type WalletRepository struct {
	db *gorm.DB
}

type WalletRepositoryParams struct {
	fx.In
	Db *gorm.DB
}

func NewWalletRepository(params WalletRepositoryParams) IWalletRepository {
	return &WalletRepository{
		params.Db,
	}
}

func (repo *WalletRepository) InsertWallet(wallet *model.Wallet) error {
	return repo.db.Create(wallet).Error
}
