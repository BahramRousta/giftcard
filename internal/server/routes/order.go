package routes

import (
	orderDelivery "giftcard/internal/modules/order/delivery/http"
	"github.com/labstack/echo/v4"
)

func MapOrderHandler(g *echo.Group, d *orderDelivery.OrderHandler) {
	g.POST("/order/create", d.CreateOrder)
	g.POST("/order/confirm", d.ConfirmOrder)
	g.GET("/order/get/status", d.RetrieveOrder)
}
