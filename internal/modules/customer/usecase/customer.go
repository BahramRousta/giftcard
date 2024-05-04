package usecase

import (
	"giftCard/internal/adaptor/giftcard"
	"giftCard/internal/modules/customer/repository"
	"giftCard/model"
	"go.uber.org/fx"
)

type CustomerUseCase struct {
	walletRepo repository.IWalletRepository
	gf         giftcard.IGiftCard
}

type CustomerUseCaseParam struct {
	fx.In
	WalletRepo repository.IWalletRepository
	Gf         *giftcard.GiftCard
}

func NewCustomerUseCase(param CustomerUseCaseParam) *CustomerUseCase {
	return &CustomerUseCase{
		walletRepo: param.WalletRepo,
		gf:         param.Gf,
	}
}

func (us CustomerUseCase) GetCustomerInfoService() (giftcard.CustomerInfoResponse, error) {

	data, err := us.gf.CustomerInfo()
	if err != nil {
		return giftcard.CustomerInfoResponse{}, err
	}
	currency := "EUR"
	wallet := &model.Wallet{
		Currency:      currency,
		Balance:       data.Data.Wallet.EUR.Balance,
		CreditBalance: float64(data.Data.Wallet.EUR.CreditBalance),
		FrozenBalance: float64(data.Data.Wallet.EUR.FrozenBalance),
	}

	if err := us.walletRepo.InsertWallet(wallet); err != nil {
		return giftcard.CustomerInfoResponse{}, err
	}

	return data, nil
}
