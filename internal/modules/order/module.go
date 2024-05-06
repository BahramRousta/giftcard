package order

import (
	"giftcard/internal/modules/order/delivery/http"
	"giftcard/internal/modules/order/repository"
	"giftcard/internal/modules/order/usecase"
	"go.uber.org/fx"
)

var Module = fx.Module("order",
	fx.Provide(usecase.NewOrderUseCase),
	fx.Provide(delivery.NewOrderHandler),
	fx.Provide(repository.NewOrderRepository),
)
