package handlers

import (
	adaptor "giftCard/internal/adaptor/giftcard"
	"github.com/labstack/echo/v4"
	"net/http"
)

func ShopList(c echo.Context) error {
	gf := adaptor.NewGiftCard()
	data, err := gf.ShopList()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, data)
}
