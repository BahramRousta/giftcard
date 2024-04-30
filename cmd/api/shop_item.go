package api

import (
	"errors"
	gftErr "giftCard/internal/adaptor/gft_error"
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
		var forbiddenErr *gftErr.ForbiddenErr
		if errors.As(err, &forbiddenErr) {
			return c.JSON(http.StatusForbidden, map[string]any{
				"message": forbiddenErr.ErrMsg,
				"data":    "",
				"success": false,
			})
		}
		var reqErr *gftErr.RequestErr
		if errors.As(err, &reqErr) {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": reqErr.ErrMsg,
				"data":    reqErr.Response,
				"success": false,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]any{"data": "", "message": err.Error(), "success": false})

	}
	return c.JSON(http.StatusOK, map[string]any{"data": data, "message": "", "success": true})
}
