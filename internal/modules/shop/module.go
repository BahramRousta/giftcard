package shop

import (
	"giftcard/internal/modules/shop/delivery/http"
	"giftcard/internal/modules/shop/usecase"
	"go.uber.org/fx"
)

var Module = fx.Module("shop",
	fx.Provide(usecase.NewShopUseCase),
	fx.Provide(delivery.NewShopHandler),
)
