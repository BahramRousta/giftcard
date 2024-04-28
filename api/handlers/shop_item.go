package handlers

import (
	"errors"
	adaptor "giftCard/internal/adaptor/giftcard"
	"giftCard/internal/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
)

func ShopItem(c echo.Context) error {
	productId := c.QueryParam("productId")
	if productId == "" {
		return c.String(http.StatusBadRequest, "productId is required")
	}

	data, err := usecase.ShopItemUseCase(productId)
	if err != nil {
		if err != nil {
			var shopItemErr *adaptor.ShopItemError
			if errors.As(err, &shopItemErr) {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"message": shopItemErr.Response["message"],
					"success": false,
					"data":    map[string]any{},
				})
			}
			return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Something went wrong"})
		}
	}
	return c.JSON(http.StatusCreated, map[string]any{"data": data["data"], "message": "", "success": true})
}
