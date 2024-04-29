package service

import (
	adaptor "giftCard/internal/adaptor/giftcard"
	"giftCard/internal/model"
	"giftCard/internal/repository"
)

type CustomerService struct {
	walletRepo   *repository.WalletRepository
	exchangeRepo *repository.ExchangeRepository
}

func NewCustomerService(walletRepo *repository.WalletRepository, exchangeRepo *repository.ExchangeRepository) *CustomerService {
	return &CustomerService{
		walletRepo:   walletRepo,
		exchangeRepo: exchangeRepo,
	}
}

func (s CustomerService) GetCustomerInfoService() (adaptor.CustomerInfoResponse, error) {
	gf := adaptor.NewGiftCard()
	data, err := gf.CustomerInfo()
	if err != nil {
		return adaptor.CustomerInfoResponse{}, err
	}
	currency := "EUR"
	wallet := &model.Wallet{
		Currency:      currency,
		Balance:       data.Data.Wallet.EUR.Balance,
		CreditBalance: float64(data.Data.Wallet.EUR.CreditBalance),
		FrozenBalance: float64(data.Data.Wallet.EUR.FrozenBalance),
	}

	if err := s.walletRepo.InsertWallet(wallet); err != nil {
		return adaptor.CustomerInfoResponse{}, err
	}

	return data, nil
}
