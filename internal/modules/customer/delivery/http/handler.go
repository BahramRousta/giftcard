package delivery

import (
	"bytes"
	"context"
	"errors"
	gftErr "giftcard/internal/adaptor/giftcard"
	"giftcard/internal/adaptor/trace"
	"giftcard/internal/exceptions"
	"giftcard/internal/modules/customer/usecase"
	"giftcard/pkg/requester"
	"giftcard/pkg/responser"
	"giftcard/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	_ "go.opentelemetry.io/otel/attribute"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type CustomerInfoHandler struct {
	customerUseCase usecase.ICustomerUseCase
}

type CustomerInfoHandlerParams struct {
	fx.In
	CustomerUseCase *usecase.CustomerUseCase
}

func NewCustomerInfoHandler(params CustomerInfoHandlerParams) *CustomerInfoHandler {
	return &CustomerInfoHandler{
		customerUseCase: params.CustomerUseCase,
	}
}

func (h CustomerInfoHandler) CustomerInfo(c echo.Context) error {

	span, spannedContext := trace.T.SpanFromContext(
		utils.GetRequestCtx(c),
		"CustomerInfo[CustomerDelivery]",
		"delivery")
	defer span.End()

	requestBody := ""
	if c.Request().Body != nil {
		body, err := io.ReadAll(c.Request().Body)
		if err == nil {
			requestBody = string(body)
			c.Request().Body = io.NopCloser(bytes.NewBuffer(body))
		}
	}

	uniqueID := c.Response().Header().Get(echo.HeaderXRequestID)
	logger := zap.L().With(
		zap.String("tracer", uniqueID),
	)

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
	logger.Info("Request from client", zap.Any("data", request))
	span.SetAttributes(attribute.String("Request from client", utils.Marshal(request)))

	ctx := context.WithValue(spannedContext, "tracer", uniqueID)
	data, err := h.customerUseCase.GetCustomerInfoUseCase(ctx)

	if err != nil {
		var forbiddenErr *gftErr.ForbiddenErr
		if errors.As(err, &forbiddenErr) {
			logger.Info("Response to client", zap.Any("error", forbiddenErr.ErrMsg))
			span.SetAttributes(attribute.String(exceptions.StatusForbidden, forbiddenErr.ErrMsg))
			return c.JSON(http.StatusForbidden, responser.Response{
				Message: forbiddenErr.ErrMsg,
				Data:    "",
				Success: false,
			})
		}

		var reqErr *gftErr.RequestErr
		if errors.As(err, &reqErr) {
			logger.Info("Response to client", zap.Any("error", reqErr.Error))
			span.SetAttributes(attribute.String(exceptions.StatusBadRequest, reqErr.ErrMsg))
			return c.JSON(http.StatusBadRequest, responser.Response{
				Message: reqErr.ErrMsg,
				Data:    reqErr.Response,
				Success: false,
			})
		}
		logger.Info("Response to client", zap.Any("error", err.Error()))
		span.SetAttributes(attribute.String(exceptions.InternalServerError, err.Error()))
		return c.JSON(http.StatusInternalServerError, responser.Response{
			Message: exceptions.InternalServerError,
			Data:    "",
			Success: false})
	}

	response := responser.Response{
		Message: "",
		Success: true,
		Data:    data.Data,
	}
	logger.Info("Response to client", zap.Any("data", response))
	span.SetAttributes(attribute.String("Request to client", utils.Marshal(response)))
	return c.JSON(http.StatusOK, response)
}
