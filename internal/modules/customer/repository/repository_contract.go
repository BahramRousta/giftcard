package repository

import "giftCard/model"

type IWalletRepository interface {
	InsertWallet(wallet *model.Wallet) error
}
