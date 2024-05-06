package routes

import (
	customerDelivery "giftcard/internal/modules/customer/delivery/http"
	"github.com/labstack/echo/v4"
)

func MapCustomerHandler(g *echo.Group, d *customerDelivery.CustomerInfoHandler) {
	g.GET("/customer/info", d.CustomerInfo)
}
