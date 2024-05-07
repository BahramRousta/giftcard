package customer

import (
	"giftcard/internal/modules/customer/delivery/http"
	"giftcard/internal/modules/customer/repository"
	"giftcard/internal/modules/customer/usecase"
	"go.uber.org/fx"
)

var Module = fx.Module("customer",
	fx.Provide(usecase.NewCustomerUseCase),
	fx.Provide(delivery.NewCustomerInfoHandler),
	fx.Provide(repository.NewWalletRepository),
)
