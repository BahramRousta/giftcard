package shop

import (
	"giftCard/internal/modules/shop/delivery/http"
	"giftCard/internal/modules/shop/usecase"
	"go.uber.org/fx"
)

var Module = fx.Module("shop",
	fx.Provide(usecase.NewShopUseCase),
	fx.Provide(delivery.NewShopHandler),
)
