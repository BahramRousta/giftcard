package delivery

import (
	"context"
	"errors"
	gftErr "giftcard/internal/adaptor/giftcard"
	"giftcard/internal/adaptor/trace"
	"giftcard/internal/exceptions"
	"giftcard/internal/modules/order/usecase"
	"giftcard/pkg/requester"
	"giftcard/pkg/responser"
	"giftcard/pkg/utils"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
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
	span, spannedContext := trace.T.SpanFromContext(
		utils.GetRequestCtx(c),
		"ConfirmOrder[OrderDelivery]",
		"delivery")
	defer span.End()

	var requestBody confirmOrderRequestBody
	if err := c.Bind(&requestBody); err != nil {
		span.SetAttributes(attribute.String(exceptions.StatusBadRequest, err.Error()))
		return c.JSON(http.StatusBadRequest, responser.Response{
			Message: exceptions.InvalidConfirmOrderInput,
			Data:    "",
			Success: false})
	}

	validate := validator.New()
	if err := validate.Struct(&requestBody); err != nil {
		span.SetAttributes(attribute.String(exceptions.StatusBadRequest, err.Error()))
		return c.JSON(http.StatusBadRequest, responser.Response{
			Message: exceptions.InvalidConfirmOrderInput,
			Data:    "",
			Success: false})
	}

	uniqueID := c.Response().Header().Get(echo.HeaderXRequestID)
	request := requester.Request{
		ID:          uniqueID,
		RequestBody: requestBody,
		UserIP:      c.RealIP(),
		Uri:         c.Path(),
		Method:      c.Request().Method,
		Host:        c.Request().Host,
		Header:      c.Request().Header,
		Params:      c.QueryParams(),
	}
	span.SetAttributes(attribute.String("Request", utils.Marshal(request)))

	ctx := context.WithValue(spannedContext, "tracer", uniqueID)
	data, err := h.us.ConfirmOrder(ctx, requestBody.OrderId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.SetAttributes(attribute.String(exceptions.StatusBadRequest, gorm.ErrRecordNotFound.Error()))
			return c.JSON(http.StatusBadRequest, responser.Response{
				Message: err.Error(),
				Data:    "",
				Success: false,
			})
		}

		var forbiddenErr *gftErr.ForbiddenErr
		if errors.As(err, &forbiddenErr) {
			span.SetAttributes(attribute.String(exceptions.StatusForbidden, forbiddenErr.ErrMsg))
			return c.JSON(http.StatusForbidden, responser.Response{
				Message: forbiddenErr.ErrMsg,
				Data:    "",
				Success: false,
			})
		}
		var reqErr *gftErr.RequestErr
		if errors.As(err, &reqErr) {
			span.SetAttributes(attribute.String(exceptions.StatusBadRequest, reqErr.ErrMsg))
			return c.JSON(http.StatusBadRequest, responser.Response{
				Message: reqErr.ErrMsg,
				Data:    reqErr.Response,
				Success: false,
			})
		}

		span.SetAttributes(attribute.String(exceptions.InternalServerError, err.Error()))
		return c.JSON(http.StatusInternalServerError, responser.Response{
			Message: exceptions.InternalServerError,
			Data:    "",
			Success: false,
		})
	}

	response := responser.Response{
		Message: "",
		Success: true,
		Data:    data["data"],
	}
	span.SetAttributes(attribute.String("Response", utils.Marshal(response)))

	return c.JSON(http.StatusCreated, response)
}

func (h *OrderHandler) CreateOrder(c echo.Context) error {
	span, spannedContext := trace.T.SpanFromContext(
		utils.GetRequestCtx(c),
		"CreateOrder[OrderDelivery]",
		"delivery")
	defer span.End()

	var requestBody RequestBody
	if err := c.Bind(&requestBody); err != nil {
		span.SetAttributes(attribute.String(exceptions.StatusBadRequest, exceptions.InvalidCreateOrderInput))
		return c.JSON(http.StatusBadRequest, responser.Response{
			Data:    "",
			Message: exceptions.InvalidInput,
			Success: false,
		})
	}

	validate := validator.New()
	if err := validate.Struct(&requestBody); err != nil {
		span.SetAttributes(attribute.String(exceptions.StatusBadRequest, err.Error()))
		return c.JSON(http.StatusBadRequest, responser.Response{
			Data:    "",
			Message: exceptions.InvalidCreateOrderInput,
			Success: false,
		})
	}

	if len(requestBody.ProductList) == 0 {
		span.SetAttributes(attribute.String(exceptions.StatusBadRequest, exceptions.EmptyProductList))
		return c.JSON(http.StatusBadRequest, responser.Response{
			Data:    "",
			Message: exceptions.EmptyProductList,
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

	uniqueID := c.Response().Header().Get(echo.HeaderXRequestID)
	request := requester.Request{
		ID:          uniqueID,
		RequestBody: requestBody,
		UserIP:      c.RealIP(),
		Uri:         c.Path(),
		Method:      c.Request().Method,
		Host:        c.Request().Host,
		Header:      c.Request().Header,
		Params:      c.QueryParams(),
	}
	span.SetAttributes(attribute.String("Request", utils.Marshal(request)))

	ctx := context.WithValue(spannedContext, "tracer", uniqueID)
	data, err := h.us.CreateOrder(ctx, productList)

	if err != nil {

		var forbiddenErr *gftErr.ForbiddenErr
		if errors.As(err, &forbiddenErr) {
			span.SetAttributes(attribute.String(exceptions.StatusForbidden, forbiddenErr.ErrMsg))
			return c.JSON(http.StatusForbidden, responser.Response{
				Message: forbiddenErr.ErrMsg,
				Data:    "",
				Success: false,
			})
		}

		var reqErr *gftErr.RequestErr
		if errors.As(err, &reqErr) {
			span.SetAttributes(attribute.String(exceptions.StatusBadRequest, reqErr.ErrMsg))
			return c.JSON(http.StatusBadRequest, responser.Response{
				Message: "",
				Data:    reqErr.Response,
				Success: false,
			})
		}

		span.SetAttributes(attribute.String(exceptions.InternalServerError, err.Error()))
		return c.JSON(http.StatusInternalServerError, responser.Response{
			Data:    "",
			Message: exceptions.InternalServerError,
			Success: false,
		})
	}

	response := responser.Response{
		Message: "",
		Success: true,
		Data:    data.Data,
	}
	span.SetAttributes(attribute.String("Response", utils.Marshal(response)))

	return c.JSON(http.StatusCreated, response)
}

func (h *OrderHandler) RetrieveOrder(c echo.Context) error {
	span, spannedContext := trace.T.SpanFromContext(
		utils.GetRequestCtx(c),
		"RetrieveOrder[OrderDelivery]",
		"delivery")
	defer span.End()

	var retOrderReq retrieveOrderRequest
	if err := c.Bind(&retOrderReq); err != nil {
		span.SetAttributes(attribute.String(exceptions.StatusBadRequest, err.Error()))
		return c.JSON(http.StatusBadRequest, responser.Response{Data: "", Message: exceptions.InvalidInput, Success: false})
	}

	orderId := retOrderReq.OrderId
	if orderId == "" {
		span.SetAttributes(attribute.String(exceptions.StatusBadRequest, "Order ID is required"))
		return c.JSON(http.StatusBadRequest, responser.Response{
			Data:    "",
			Message: exceptions.RequiredOrderID,
			Success: false})
	}

	uniqueID := c.Response().Header().Get(echo.HeaderXRequestID)
	request := requester.Request{
		ID:          uniqueID,
		RequestBody: retOrderReq,
		UserIP:      c.RealIP(),
		Uri:         c.Path(),
		Method:      c.Request().Method,
		Host:        c.Request().Host,
		Header:      c.Request().Header,
		Params:      c.QueryParams(),
	}
	span.SetAttributes(attribute.String("Request", utils.Marshal(request)))

	ctx := context.WithValue(spannedContext, "tracer", uniqueID)
	data, err := h.us.GetOrderStatus(ctx, orderId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.SetAttributes(attribute.String(exceptions.StatusBadRequest, gorm.ErrRecordNotFound.Error()))
			return c.JSON(http.StatusBadRequest, responser.Response{
				Message: exceptions.RecordNotFound,
				Data:    "",
				Success: false,
			})
		}
		var forbiddenErr *gftErr.ForbiddenErr
		if errors.As(err, &forbiddenErr) {
			span.SetAttributes(attribute.String(exceptions.StatusForbidden, forbiddenErr.ErrMsg))
			return c.JSON(http.StatusForbidden, responser.Response{
				Message: forbiddenErr.ErrMsg,
				Data:    "",
				Success: false,
			})
		}
		var reqErr *gftErr.RequestErr
		if errors.As(err, &reqErr) {
			span.SetAttributes(attribute.String(exceptions.StatusBadRequest, reqErr.ErrMsg))
			return c.JSON(http.StatusBadRequest, responser.Response{
				Message: reqErr.ErrMsg,
				Data:    reqErr.Response,
				Success: false,
			})
		}
		span.SetAttributes(attribute.String(exceptions.InternalServerError, err.Error()))
		return c.JSON(http.StatusInternalServerError, responser.Response{
			Data:    "",
			Message: exceptions.InternalServerError,
			Success: false})
	}

	response := responser.Response{
		Message: "",
		Success: true,
		Data:    data["data"],
	}
	span.SetAttributes(attribute.String("Response", utils.Marshal(response)))

	return c.JSON(http.StatusOK, response)
}
