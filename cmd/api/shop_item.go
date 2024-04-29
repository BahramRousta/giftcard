package api

import (
	"errors"
	adaptor "giftCard/internal/adaptor/giftcard"
	"giftCard/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ShopItemHandler struct {
	shopItemService *service.ShopItemService
}

func NewShopItemHandler(s *service.ShopItemService) *ShopItemHandler {
	return &ShopItemHandler{
		shopItemService: s,
	}
}

func (h *ShopItemHandler) ShopItem(c echo.Context) error {
	productId := c.QueryParam("productId")
	if productId == "" {
		return c.String(http.StatusBadRequest, "productId is required")
	}

	data, err := h.shopItemService.GetShopItemService(productId)
	if err != nil {
		if err != nil {
			var shopItemErr *adaptor.ShopItemError
			if errors.As(err, &shopItemErr) {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"message": shopItemErr.Response,
					"success": false,
					"data":    map[string]any{},
				})
			}
			return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Something went wrong"})
		}
	}
	return c.JSON(http.StatusCreated, map[string]any{"data": data, "message": "", "success": true})
}
