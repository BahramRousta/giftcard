package routes

import (
	shopDelivery "giftCard/internal/modules/shop/delivery/http"
	"github.com/labstack/echo/v4"
)

func MapShopHandler(g *echo.Group, d *shopDelivery.ShopHandler) {
	g.GET("/shop/products", d.ShopList)
	g.GET("/shop/product/info", d.ShopItem)
}
