package api

import (
	"errors"
	gftErr "giftCard/internal/adaptor/gft_error"
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
		return c.JSON(http.StatusBadRequest, map[string]any{
			"data":    "",
			"message": "bad request",
			"success": false,
		})
	}

	validate := validator.New()
	if err := validate.Struct(&requestBody); err != nil {
		var messages []string
		for _, fieldErr := range err.(validator.ValidationErrors) {
			messages = append(messages, fieldErr.Field()+" is invalid")
		}
		return c.JSON(http.StatusBadRequest, map[string]any{
			"data":    strings.Join(messages, ", "),
			"message": "",
			"success": false,
		})
	}
	if len(requestBody.ProductList) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"data":    "",
			"message": "product list can not be empty",
			"success": false,
		})
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
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
			"data":    "",
			"message": "something went wrong",
			"success": false,
		})
	}
	return c.JSON(http.StatusCreated, map[string]any{"data": data, "message": "", "success": true})
}
