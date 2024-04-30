package api

import (
	"errors"
	gftErr "giftCard/internal/adaptor/gft_error"
	"giftCard/internal/service"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

type retrieveOrderRequest struct {
	OrderId string `json:"orderId" query:"orderId" validate:"required"`
}

type RetrieveOrderHandler struct {
	service *service.RetrieveOrderService
}

func NewRetrieveOrderHandler(service *service.RetrieveOrderService) *RetrieveOrderHandler {
	return &RetrieveOrderHandler{
		service: service,
	}
}

func (h *RetrieveOrderHandler) RetrieveOrder(c echo.Context) error {

	var request retrieveOrderRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"gft_error": "Bad Request"})
	}

	orderId := request.OrderId
	if orderId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"gft_error": "Order ID is required"})
	}

	data, err := h.service.GetOrderStatusService(orderId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": err.Error(),
				"data":    "",
				"success": false,
			})
		}
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
	return c.JSON(http.StatusOK, map[string]any{"data": data["data"], "message": "", "success": true})
}
