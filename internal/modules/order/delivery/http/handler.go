package delivery

import (
	"context"
	"errors"
	"fmt"
	gftErr "giftcard/internal/adaptor/giftcard"
	"giftcard/internal/modules/order/usecase"
	"giftcard/pkg/responser"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
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
	us     usecase.IOrderUseCase
	Logger *zap.Logger
}

type OrderHandlerParams struct {
	fx.In
	Us usecase.IOrderUseCase
	//Logger *zap.Logger
}

func NewOrderHandler(params OrderHandlerParams) *OrderHandler {
	return &OrderHandler{
		us: params.Us,
		//Logger: params.Logger,
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

	uniqueID := uuid.New().String()

	logger := zap.L().With(zap.String("tracer", uniqueID))
	logger.Info("confirm order",
		zap.Any("request_body", requestBody),
		zap.String("user_ip", c.RealIP()),
		zap.String("uri", c.Path()),
		zap.String("method", c.Request().Method),
		zap.String("host", c.Request().Host),
		zap.Any("headers", c.Request().Header),
	)

	ctx := context.WithValue(c.Request().Context(), "tracer", uniqueID)
	data, err := h.us.ConfirmOrder(ctx, requestBody.OrderId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Info("response confirm order",
				zap.Any("error", gorm.ErrRecordNotFound),
			)

			return c.JSON(http.StatusBadRequest, responser.Response{
				Message: err.Error(),
				Data:    "",
				Success: false,
			})
		}

		var forbiddenErr *gftErr.ForbiddenErr
		if errors.As(err, &forbiddenErr) {

			logger.Info("response confirm order",
				zap.Any("error", forbiddenErr.ErrMsg),
			)

			return c.JSON(http.StatusForbidden, responser.Response{
				Message: forbiddenErr.ErrMsg,
				Data:    "",
				Success: false,
			})
		}
		var reqErr *gftErr.RequestErr
		if errors.As(err, &reqErr) {
			logger.Info("response confirm order",
				zap.Any("error", reqErr.ErrMsg),
				zap.Any("message", reqErr.Response),
			)
			return c.JSON(http.StatusBadRequest, responser.Response{
				Message: reqErr.ErrMsg,
				Data:    reqErr.Response,
				Success: false,
			})
		}

		logger.Info("response confirm order",
			zap.Any("error", "internal error"),
		)

		return c.JSON(http.StatusInternalServerError, responser.Response{
			Message: "something went wrong",
			Data:    "",
			Success: false,
		})
	}

	logger.Info("confirm order response",
		zap.Any("data", data["data"]))

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

	uniqueID := uuid.New().String()
	ctx := context.WithValue(c.Request().Context(), "tracer", uniqueID)

	data, err := h.us.CreateOrder(ctx, productList)

	logger := zap.L().With(zap.String("tracer", uniqueID))
	logger.Info("create order",
		zap.Any("request_body", requestBody),
		zap.String("user_ip", c.RealIP()),
		zap.String("uri", c.Path()),
		zap.String("method", c.Request().Method),
		zap.String("host", c.Request().Host),
		zap.Any("headers", c.Request().Header),
	)

	if err != nil {
		var forbiddenErr *gftErr.ForbiddenErr
		if errors.As(err, &forbiddenErr) {
			logger.Info("response create order",
				zap.String("message", forbiddenErr.ErrMsg),
			)
			return c.JSON(http.StatusForbidden, responser.Response{
				Message: forbiddenErr.ErrMsg,
				Data:    "",
				Success: false,
			})
		}
		var reqErr *gftErr.RequestErr
		if errors.As(err, &reqErr) {

			logger.Info("response create order",
				zap.String("message", reqErr.ErrMsg),
				zap.Any("data", reqErr.Response),
			)

			return c.JSON(http.StatusBadRequest, responser.Response{
				Message: fmt.Sprintf("%s", reqErr.Response["message"]),
				Data:    "",
				Success: false,
			})
		}
		logger.Info("response create order",
			zap.String("message", "internal error"),
		)
		return c.JSON(http.StatusInternalServerError, responser.Response{
			Data:    "",
			Message: "something went wrong",
			Success: false,
		})
	}

	logger.Info("response create order",
		zap.Any("data", data.Data))

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
	queryParams := c.QueryParams()
	uniqueID := uuid.New().String()

	logger := zap.L().With(zap.String("tracer", uniqueID))
	logger.Info("get order status",
		zap.String("user_ip", c.RealIP()),
		zap.String("uri", c.Path()),
		zap.Any("params", queryParams),
		zap.String("method", c.Request().Method),
		zap.String("host", c.Request().Host),
		zap.Any("headers", c.Request().Header),
	)

	ctx := context.WithValue(c.Request().Context(), "tracer", uniqueID)
	data, err := h.us.GetOrderStatus(ctx, orderId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			logger.Info("response get order status",
				zap.Any("message", gorm.ErrRecordNotFound))

			return c.JSON(http.StatusBadRequest, responser.Response{
				Message: "record not found",
				Data:    "",
				Success: false,
			})
		}
		var forbiddenErr *gftErr.ForbiddenErr
		if errors.As(err, &forbiddenErr) {

			logger.Info("response get order status",
				zap.String("message", forbiddenErr.ErrMsg),
			)

			return c.JSON(http.StatusForbidden, responser.Response{
				Message: forbiddenErr.ErrMsg,
				Data:    "",
				Success: false,
			})
		}
		var reqErr *gftErr.RequestErr
		if errors.As(err, &reqErr) {

			logger.Info("response get order status",
				zap.String("message", reqErr.ErrMsg),
				zap.Any("data", reqErr.Response),
			)

			return c.JSON(http.StatusBadRequest, responser.Response{
				Message: reqErr.ErrMsg,
				Data:    reqErr.Response,
				Success: false,
			})
		}
		logger.Info("response get order status",
			zap.String("message", "internal error"),
			zap.String("data", err.Error()),
		)
		return c.JSON(http.StatusInternalServerError, responser.Response{Data: "", Message: "something went wrong", Success: false})
	}

	logger.Info("response get order status",
		zap.Any("data", data["data"]),
	)

	return c.JSON(http.StatusOK, responser.Response{Data: data["data"], Message: "", Success: true})
}
