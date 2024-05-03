package routes

import (
	orderDelivery "giftCard/internal/modules/order/delivery/http"
	"github.com/labstack/echo/v4"
)

func MapOrderHandler(g *echo.Group, d orderDelivery.OrderHandler) {
	g.GET("/order/create", d.CreateOrder)
	g.GET("/order/confirm", d.ConfirmOrder)
	g.GET("/order/get/status", d.RetrieveOrder)
}
