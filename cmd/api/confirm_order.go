package api

import (
	"errors"
	adaptor "giftCard/internal/adaptor/giftcard"
	"giftCard/internal/service"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type confirmOrderRequestBody struct {
	OrderId string `json:"orderId" validate:"required"`
}

type ConfirmOrderHandler struct {
	service *service.ConfirmOrderService
}

func NewConfirmOrderHandler(service *service.ConfirmOrderService) *ConfirmOrderHandler {
	return &ConfirmOrderHandler{
		service: service,
	}
}

func (h *ConfirmOrderHandler) ConfirmOrder(c echo.Context) error {
	var requestBody confirmOrderRequestBody
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad Request"})
	}
	validate := validator.New()
	if err := validate.Struct(&requestBody); err != nil {
		var messages []string
		for _, fieldErr := range err.(validator.ValidationErrors) {
			messages = append(messages, fieldErr.Field()+" is invalid")
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": strings.Join(messages, ", ")})
	}

	data, err := h.service.OrderConfirmService(requestBody.OrderId)
	if err != nil {
		var orderErr *adaptor.ConfirmOrderError
		if errors.As(err, &orderErr) {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error":   orderErr.ErrorMsg,
				"payload": orderErr.Response,
			})
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, map[string]any{"data": data["data"], "message": "", "success": true})
}
