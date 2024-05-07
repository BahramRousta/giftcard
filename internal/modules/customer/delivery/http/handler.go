package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	gftErr "giftcard/internal/adaptor/giftcard"
	"giftcard/internal/adaptor/trace"
	"giftcard/internal/modules/customer/usecase"
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
	logger.Info("customer info request",
		zap.String("request_body", requestBody),
		zap.String("user_ip", c.RealIP()),
		zap.String("uri", c.Path()),
		zap.String("method", c.Request().Method),
		zap.String("host", c.Request().Host),
		zap.Any("header", c.Request().Header),
	)

	ctx := context.WithValue(spannedContext, "tracer", uniqueID)
	data, err := h.customerUseCase.GetCustomerInfoUseCase(ctx)

	if err != nil {
		var forbiddenErr *gftErr.ForbiddenErr
		if errors.As(err, &forbiddenErr) {
			logger.Info("response to client",
				zap.Any("data", forbiddenErr.ErrMsg),
			)
			span.SetAttributes(attribute.String("err", forbiddenErr.ErrMsg))
			return c.JSON(http.StatusForbidden, responser.Response{
				Message: forbiddenErr.ErrMsg,
				Data:    "",
				Success: false,
			})
		}

		var reqErr *gftErr.RequestErr
		if errors.As(err, &reqErr) {
			logger.Info("response to client",
				zap.Any("error", reqErr.ErrMsg),
				zap.Any("data", reqErr.Response),
			)
			span.SetAttributes(attribute.String("err", reqErr.ErrMsg))
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": reqErr.ErrMsg,
				"data":    reqErr.Response,
				"success": false,
			})
		}
		logger.Error("response to client",
			zap.Any("error", "internal error"),
		)
		span.SetAttributes(attribute.String("err", err.Error()))
		return c.JSON(http.StatusInternalServerError, responser.Response{Message: "Something went wrong", Data: "", Success: false})
	}

	//logger.Info("customer info response", zap.Any("data", data.Data))

	dataJSON, err := json.Marshal(data.Data)
	span.SetAttributes(
		attribute.String("message", "get customer info successfully passed."),
		attribute.String("data", string(dataJSON)),
	)
	return c.JSON(http.StatusOK, responser.Response{Message: "", Success: true, Data: data.Data})
}
