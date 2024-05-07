package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	gftErr "giftcard/internal/adaptor/giftcard"
	"giftcard/internal/adaptor/trace"
	"giftcard/internal/modules/shop/usecase"
	"giftcard/pkg/responser"
	"giftcard/pkg/utils"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type ShopHandler struct {
	us usecase.IShopUseCase
}
type ShopHandlerParams struct {
	fx.In
	Us usecase.IShopUseCase
}

func NewShopHandler(params ShopHandlerParams) *ShopHandler {
	return &ShopHandler{
		us: params.Us,
	}
}

func (h *ShopHandler) ShopItem(c echo.Context) error {
	span, spannedContext := trace.T.SpanFromContext(
		utils.GetRequestCtx(c),
		"ShopItem[ShopDelivery]",
		"delivery")
	defer span.End()

	productId := c.QueryParam("productId")
	if productId == "" {
		return c.String(http.StatusBadRequest, "productId is required")
	}

	uniqueID := c.Response().Header().Get(echo.HeaderXRequestID)

	logger := zap.L().With(
		zap.String("tracer", uniqueID),
	)
	logger.Info("request to provider",
		zap.String("user_ip", c.RealIP()),
		zap.String("uri", c.Path()),
		zap.String("method", c.Request().Method),
		zap.String("host", c.Request().Host),
		zap.Any("header", c.Request().Header),
	)

	ctx := context.WithValue(spannedContext, "tracer", uniqueID)

	data, err := h.us.GetShopItem(ctx, productId)
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
			return c.JSON(http.StatusBadRequest, responser.Response{
				Message: reqErr.ErrMsg,
				Data:    reqErr.Response,
				Success: false,
			})
		}
		logger.Info("response to client",
			zap.Any("error", err.Error()),
		)
		span.SetAttributes(attribute.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, responser.Response{
			Message: "something went wrong",
			Data:    "",
			Success: false,
		})
	}

	//logger.Info("response to client",
	//	zap.Any("data", data.Data),
	//)

	dataJSON, err := json.Marshal(data.Data)
	span.SetAttributes(
		attribute.String("message", "get shop item info successfully passed."),
		attribute.String("data", string(dataJSON)),
	)
	return c.JSON(http.StatusOK, responser.Response{
		Message: "",
		Data:    data.Data,
		Success: true,
	})
}

func (h *ShopHandler) ShopList(c echo.Context) error {
	span, spannedContext := trace.T.SpanFromContext(
		utils.GetRequestCtx(c),
		"ShopList[ShopDelivery]",
		"delivery")
	defer span.End()

	uniqueID := c.Response().Header().Get(echo.HeaderXRequestID)
	ctx := context.WithValue(spannedContext, "tracer", uniqueID)

	logger := zap.L().With(
		zap.String("tracer", uniqueID),
	)

	pageSizeHeader := c.Request().Header.Get("PageSize")
	pageSize, err := strconv.Atoi(pageSizeHeader)

	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, responser.Response{
			Message: "pageSize must be an integer",
			Success: false,
			Data:    "",
		})
	}

	validate := validator.New()
	if err := validate.Var(pageSize, "min=5,max=50"); err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, responser.Response{
			Message: fmt.Sprintf("page size must between 5 to 50."),
			Success: false,
			Data:    "",
		})
	}

	pageToken := c.Request().Header.Get("PageToken")

	logger.Info("request to provider",
		zap.String("user_ip", c.RealIP()),
		zap.String("uri", c.Path()),
		zap.String("method", c.Request().Method),
		zap.String("host", c.Request().Host),
		zap.Any("header", c.Request().Header),
	)

	data, err := h.us.GetShopList(ctx, pageSize, pageToken)

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
			span.SetAttributes(attribute.String("err", forbiddenErr.ErrMsg))
			return c.JSON(http.StatusBadRequest, responser.Response{
				Message: reqErr.ErrMsg,
				Data:    reqErr.Response,
				Success: false,
			})
		}
		logger.Info("response to client",
			zap.Any("error", err.Error()),
		)
		span.SetAttributes(attribute.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, responser.Response{
			Data:    "",
			Message: "Something went wrong",
			Success: false,
		})
	}
	//logger.Info("shop list response",
	//	zap.Any("data", data["data"]),
	//)
	dataJSON, err := json.Marshal(data)
	span.SetAttributes(
		attribute.String("message", "get shop list info successfully passed."),
		attribute.String("data", string(dataJSON)),
	)
	return c.JSON(http.StatusOK, responser.Response{Data: data["data"], Message: "", Success: true})
}
