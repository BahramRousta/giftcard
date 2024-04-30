package api

import (
	"errors"
	gftErr "giftCard/internal/adaptor/gft_error"
	"giftCard/internal/service"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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
		return c.JSON(http.StatusBadRequest, map[string]string{"gft_error": "Bad Request"})
	}
	validate := validator.New()
	if err := validate.Struct(&requestBody); err != nil {
		var messages []string
		for _, fieldErr := range err.(validator.ValidationErrors) {
			messages = append(messages, fieldErr.Field()+" is invalid")
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"gft_error": strings.Join(messages, ", ")})
	}

	data, err := h.service.OrderConfirmService(requestBody.OrderId)
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
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"data":    "",
			"message": "something went wrong",
			"success": false,
		})
	}
	return c.JSON(http.StatusCreated, map[string]any{"data": data["data"], "message": "", "success": true})
}
