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
	"github.com/google/uuid"
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
	Logger          *zap.Logger
}

type CustomerInfoHandlerParams struct {
	fx.In
	CustomerUseCase *usecase.CustomerUseCase
	//Logger          *zap.Logger
}

func NewCustomerInfoHandler(params CustomerInfoHandlerParams) *CustomerInfoHandler {
	return &CustomerInfoHandler{
		customerUseCase: params.CustomerUseCase,
		//Logger:          params.Logger,
	}
}

func (h CustomerInfoHandler) CustomerInfo(c echo.Context) error {

	span, spannedContext := trace.T.SpanFromContext(
		utils.GetRequestCtx(c),
		"CustomerInfo[CustomerDelivery]",
		"delivery")
	defer span.End()
	span.SetAttributes(attribute.String("msg", "start fetch customer info"))
	requestBody := ""
	if c.Request().Body != nil {
		body, err := io.ReadAll(c.Request().Body)
		if err == nil {
			requestBody = string(body)
			c.Request().Body = io.NopCloser(bytes.NewBuffer(body))
		}
	}

	uniqueID := uuid.New().String()

	ctx := context.WithValue(spannedContext, "tracer", uniqueID)
	data, err := h.customerUseCase.GetCustomerInfoUseCase(ctx)

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

	if err != nil {
		var forbiddenErr *gftErr.ForbiddenErr
		if errors.As(err, &forbiddenErr) {
			span.SetAttributes(attribute.String("err", forbiddenErr.ErrMsg))
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
	dataJSON, err := json.Marshal(data.Data)
	span.SetAttributes(attribute.String("data", string(dataJSON)))

	span.SetAttributes(
		attribute.String("msg", ""),
		attribute.Bool("success", true),
		attribute.String("data", string(dataJSON)), // Convert map to string
	)
	return c.JSON(http.StatusOK, responser.Response{Message: "", Success: true, Data: data.Data})
}
