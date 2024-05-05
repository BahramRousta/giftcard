package delivery

import (
	"bytes"
	"context"
	"errors"
	gftErr "giftCard/internal/adaptor/giftcard"
	"giftCard/internal/modules/customer/usecase"
	"giftCard/pkg/responser"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type CustomerInfoHandler struct {
	customerUseCase usecase.ICustomerUseCase
	Logger          *zap.Logger
}

type CustomerInfoHandlerParams struct {
	fx.In
	CustomerUseCase *usecase.CustomerUseCase
	Logger          *zap.Logger
}

func NewCustomerInfoHandler(params CustomerInfoHandlerParams) *CustomerInfoHandler {
	return &CustomerInfoHandler{
		customerUseCase: params.CustomerUseCase,
		Logger:          params.Logger,
	}
}

func (h CustomerInfoHandler) CustomerInfo(c echo.Context) error {

	requestBody := ""
	if c.Request().Body != nil {
		body, err := io.ReadAll(c.Request().Body)
		if err == nil {
			requestBody = string(body)
			c.Request().Body = io.NopCloser(bytes.NewBuffer(body))
		}
	}

	uniqueID := uuid.New().String()

	ctx := context.WithValue(c.Request().Context(), "tracer", uniqueID)
	data, err := h.customerUseCase.GetCustomerInfoUseCase(ctx)

	logger := h.Logger.With(
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
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": reqErr.ErrMsg,
				"data":    reqErr.Response,
				"success": false,
			})
		}
		return c.JSON(http.StatusInternalServerError, responser.Response{Message: "Something went wrong", Data: "", Success: false})
	}
	logger.Info("customer info response",
		zap.Any("data", data.Data))
	return c.JSON(http.StatusOK, responser.Response{Message: "", Success: true, Data: data.Data})
}
