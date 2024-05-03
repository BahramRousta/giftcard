package customer

import (
	"giftCard/internal/modules/customer/delivery/http"
	"giftCard/internal/modules/customer/repository"
	"giftCard/internal/modules/customer/usecase"
	"go.uber.org/fx"
)

var Module = fx.Module("customer",
	fx.Provide(usecase.NewCustomerUseCase),
	fx.Provide(delivery.NewCustomerInfoHandler),
	fx.Provide(repository.NewWalletRepository),
)
