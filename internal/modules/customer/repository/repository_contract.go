package repository

import "giftcard/model"

type IWalletRepository interface {
	InsertWallet(wallet *model.Wallet) error
}
