package api

import (
	"errors"
	adaptor "giftCard/internal/adaptor/giftcard"
	"giftCard/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ShopListHandler struct {
	shopListService *service.ShopListService
}

func NewShopListHandler(service *service.ShopListService) *ShopListHandler {
	return &ShopListHandler{
		shopListService: service,
	}
}

func (h *ShopListHandler) ShopList(c echo.Context) error {
	data, err := h.shopListService.GetShopListService()
	if err != nil {
		if err != nil {
			var shopListErr *adaptor.ShopListError
			if errors.As(err, &shopListErr) {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"message": shopListErr.Response["message"],
					"success": false,
					"data":    map[string]any{},
				})
			}
			return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Something went wrong"})
		}
	}
	return c.JSON(http.StatusCreated, map[string]any{"data": data["data"], "message": "", "success": true})
}
