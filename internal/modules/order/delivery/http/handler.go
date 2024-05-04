package delivery

import (
	"errors"
	"fmt"
	gftErr "giftCard/internal/adaptor/giftcard"
	"giftCard/internal/modules/order/usecase"
	"giftCard/pkg/responser"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type confirmOrderRequestBody struct {
	OrderId string `json:"orderId" validate:"required"`
}

type retrieveOrderRequest struct {
	OrderId string `json:"orderId" query:"orderId" validate:"required"`
}

type Product struct {
	Sku         string `json:"sku" validate:"required"`
	ProductType string `json:"productType" validate:"required"`
	Quote       uint   `json:"quote" validate:"required,gt=0"`
	Quantity    uint   `json:"quantity" validate:"required,gt=0"`
}

type RequestBody struct {
	ProductList []Product `json:"productList" validate:"required,dive"`
}

type OrderHandler struct {
	us usecase.IOrderUseCase
}

type OrderHandlerParams struct {
	fx.In
	Us usecase.IOrderUseCase
}

func NewOrderHandler(params OrderHandlerParams) *OrderHandler {
	return &OrderHandler{
		us: params.Us,
	}
}

func (h *OrderHandler) ConfirmOrder(c echo.Context) error {
	var requestBody confirmOrderRequestBody
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, responser.Response{Message: "Bad Request", Data: "", Success: false})
	}
	validate := validator.New()
	if err := validate.Struct(&requestBody); err != nil {
		var messages []string
		for _, fieldErr := range err.(validator.ValidationErrors) {
			messages = append(messages, fieldErr.Field()+" is invalid")
		}
		return c.JSON(http.StatusBadRequest, responser.Response{Message: strings.Join(messages, ", "), Data: "", Success: false})
	}

	data, err := h.us.ConfirmOrder(requestBody.OrderId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusBadRequest, responser.Response{
				Message: err.Error(),
				Data:    "",
				Success: false,
			})
		}

		var forbiddenErr *gftErr.ForbiddenErr
		if errors.As(err, &forbiddenErr) {
			return c.JSON(http.StatusForbidden, responser.Response{
				Message: forbiddenErr.ErrMsg,
				Data:    "",
				Success: false,
			})
		}
		var reqErr *gftErr.RequestErr
		if errors.As(err, &reqErr) {
			return c.JSON(http.StatusBadRequest, responser.Response{
				Message: reqErr.ErrMsg,
				Data:    reqErr.Response,
				Success: false,
			})
		}
		return c.JSON(http.StatusInternalServerError, responser.Response{
			Message: "something went wrong",
			Data:    "",
			Success: false,
		})
	}
	return c.JSON(http.StatusCreated, responser.Response{Data: data["data"], Message: "", Success: true})
}

func (h *OrderHandler) CreateOrder(c echo.Context) error {
	var requestBody RequestBody
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, responser.Response{
			Data:    "",
			Message: "bad request",
			Success: false,
		})
	}

	validate := validator.New()
	if err := validate.Struct(&requestBody); err != nil {
		var messages []string
		for _, fieldErr := range err.(validator.ValidationErrors) {
			messages = append(messages, fieldErr.Field()+" is invalid")
		}
		return c.JSON(http.StatusBadRequest, responser.Response{
			Data:    "",
			Message: strings.Join(messages, ", "),
			Success: false,
		})
	}
	if len(requestBody.ProductList) == 0 {
		return c.JSON(http.StatusBadRequest, responser.Response{
			Data:    "",
			Message: "product list can not be empty",
			Success: false,
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

	data, err := h.us.CreateOrder(productList)

	if err != nil {
		var forbiddenErr *gftErr.ForbiddenErr
		if errors.As(err, &forbiddenErr) {
			return c.JSON(http.StatusForbidden, responser.Response{
				Message: forbiddenErr.ErrMsg,
				Data:    "",
				Success: false,
			})
		}
		var reqErr *gftErr.RequestErr
		if errors.As(err, &reqErr) {
			return c.JSON(http.StatusBadRequest, responser.Response{
				Message: fmt.Sprintf("%s", reqErr.Response["message"]),
				Data:    "",
				Success: false,
			})
		}
		return c.JSON(http.StatusInternalServerError, responser.Response{
			Data:    "",
			Message: "something went wrong",
			Success: false,
		})
	}
	return c.JSON(http.StatusCreated, responser.Response{Data: data.Data, Message: "", Success: true})
}

func (h *OrderHandler) RetrieveOrder(c echo.Context) error {

	var request retrieveOrderRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, responser.Response{Data: "", Message: "bas request", Success: false})
	}

	orderId := request.OrderId
	if orderId == "" {
		return c.JSON(http.StatusBadRequest, responser.Response{Data: "", Message: "Order ID is required", Success: false})
	}

	data, err := h.us.GetOrderStatus(orderId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusBadRequest, responser.Response{
				Message: "record not found",
				Data:    "",
				Success: false,
			})
		}
		var forbiddenErr *gftErr.ForbiddenErr
		if errors.As(err, &forbiddenErr) {
			return c.JSON(http.StatusForbidden, responser.Response{
				Message: forbiddenErr.ErrMsg,
				Data:    "",
				Success: false,
			})
		}
		var reqErr *gftErr.RequestErr
		if errors.As(err, &reqErr) {
			return c.JSON(http.StatusBadRequest, responser.Response{
				Message: reqErr.ErrMsg,
				Data:    reqErr.Response,
				Success: false,
			})
		}
		return c.JSON(http.StatusInternalServerError, responser.Response{Data: "", Message: "something went wrong", Success: false})
	}
	return c.JSON(http.StatusOK, responser.Response{Data: data["data"], Message: "", Success: true})
}
