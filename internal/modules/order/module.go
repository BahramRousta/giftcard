package order

import (
	"giftCard/internal/modules/order/delivery/http"
	"giftCard/internal/modules/order/repository"
	"giftCard/internal/modules/order/usecase"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(usecase.NewOrderUseCase),
	fx.Provide(repository.NewOrderRepository),
	fx.Provide(delivery.NewOrderHandler),
)
