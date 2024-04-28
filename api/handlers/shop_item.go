package handlers

import (
	adaptor "giftCard/internal/adaptor/giftcard"
	"github.com/labstack/echo/v4"
	"net/http"
)

func ShopItem(c echo.Context) error {
	productId := c.QueryParam("productId")
	if productId == "" {
		return c.String(http.StatusBadRequest, "productId is required")
	}

	gf := adaptor.NewGiftCard()
	data, err := gf.ShopItem(productId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, data)
}
