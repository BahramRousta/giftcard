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

type Product struct {
	Sku         string `json:"sku" validate:"required"`
	ProductType string `json:"productType" validate:"required"`
	Quote       uint   `json:"quote" validate:"required,gt=0"`
	Quantity    uint   `json:"quantity" validate:"required,gt=0"`
}

type RequestBody struct {
	ProductList []Product `json:"productList" validate:"required,dive"`
}

type CreateOrderHandler struct {
	service *service.OrderService
}

func NewCreateOrderHandler(service *service.OrderService) *CreateOrderHandler {
	return &CreateOrderHandler{
		service: service,
	}
}

func (h *CreateOrderHandler) CreateOrder(c echo.Context) error {
	var requestBody RequestBody
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
	if len(requestBody.ProductList) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ProductList cannot be empty"})
	}

	var productList []map[string]interface{}
	for _, product := range requestBody.ProductList {
		productMap := map[string]interface{}{
			"sku":         product.Sku,
			"productType": product.ProductType,
			"quote":       product.Quote,
			"quantity":    product.Quantity,
		}
		productList = append(productList, productMap)
	}

	data, err := h.service.CreateOrderService(productList)

	if err != nil {
		if err != nil {
			var orderErr *adaptor.CreateOrderError
			if errors.As(err, &orderErr) {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"message": orderErr.Response,
					"success": false,
					"data":    map[string]any{},
				})
			}
			return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Something went wrong"})
		}
	}
	return c.JSON(http.StatusCreated, map[string]any{"data": data, "message": "", "success": true})
}
