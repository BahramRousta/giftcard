package usecase

import (
	"context"
	"encoding/json"
	"giftcard/internal/adaptor/giftcard"
	"giftcard/internal/adaptor/trace"
	"giftcard/internal/exceptions"
	"giftcard/internal/modules/customer/repository"
	"giftcard/model"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"

	//"giftcard/pkg/utils"
	//"github.com/labstack/echo/v4"
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

func (us CustomerUseCase) GetCustomerInfoUseCase(ctx context.Context) (giftcard.CustomerInfoResponse, error) {
	span, spannedContext := trace.T.SpanFromContext(
		ctx,
		"CustomerInfoUseCase",
		"UseCase")
	defer span.End()

	uniqueID, _ := ctx.Value("tracer").(string)

	logger := zap.L().With(
		zap.String("tracer", uniqueID),
	)

	data, err := us.gf.CustomerInfo(spannedContext)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
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
		logger.Error(exceptions.DBError, zap.String("error", err.Error()))
		span.SetAttributes(attribute.String("error", err.Error()))
		return giftcard.CustomerInfoResponse{}, err
	}

	jsonData, err := json.Marshal(data.Data)
	span.SetAttributes(attribute.String("data", string(jsonData)))
	return data, nil
}
