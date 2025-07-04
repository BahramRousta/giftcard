package delivery

import (
	"context"
	"errors"
	gftErr "giftcard/internal/adaptor/giftcard"
	"giftcard/internal/adaptor/trace"
	"giftcard/internal/exceptions"
	"giftcard/internal/modules/shop/usecase"
	"giftcard/pkg/requester"
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
	uniqueID := c.Response().Header().Get(echo.HeaderXRequestID)
	logger := zap.L().With(
		zap.String("tracer", uniqueID),
	)

	request := requester.Request{
		ID:          uniqueID,
		RequestBody: "",
		UserIP:      c.RealIP(),
		Uri:         c.Path(),
		Method:      c.Request().Method,
		Host:        c.Request().Host,
		Header:      c.Request().Header,
		Params:      c.QueryParams(),
	}
	logger.Info("Request from client", zap.Any("data", request))
	span.SetAttributes(attribute.String("Request", utils.Marshal(request)))

	productId := c.QueryParam("productId")
	if productId == "" {
		logger.Info("Response to client", zap.Any("error", exceptions.ProductIDError))
		span.SetAttributes(attribute.String(exceptions.StatusBadRequest, exceptions.ProductIDError))
		return c.String(http.StatusBadRequest, exceptions.ProductIDError)
	}

	ctx := context.WithValue(spannedContext, "tracer", uniqueID)

	data, err := h.us.GetShopItem(ctx, productId)
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
			logger.Info("Response to client", zap.Any("error", reqErr.ErrMsg))
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
			Success: false,
		})
	}

	response := responser.Response{
		Message: "",
		Success: true,
		Data:    data.Data,
	}
	logger.Info("Response to client", zap.Any("data", response))
	span.SetAttributes(attribute.String("Response", utils.Marshal(response)))
	return c.JSON(http.StatusOK, response)
}

func (h *ShopHandler) ShopList(c echo.Context) error {
	span, spannedContext := trace.T.SpanFromContext(
		utils.GetRequestCtx(c),
		"ShopList[ShopDelivery]",
		"delivery")
	defer span.End()

	uniqueID := c.Response().Header().Get(echo.HeaderXRequestID)
	logger := zap.L().With(
		zap.String("tracer", uniqueID),
	)

	request := requester.Request{
		ID:          uniqueID,
		RequestBody: "",
		UserIP:      c.RealIP(),
		Uri:         c.Path(),
		Method:      c.Request().Method,
		Host:        c.Request().Host,
		Header:      c.Request().Header,
		Params:      c.QueryParams(),
	}
	logger.Info("Request from client", zap.Any("data", request))
	span.SetAttributes(attribute.String("Request", utils.Marshal(request)))

	ctx := context.WithValue(spannedContext, "tracer", uniqueID)

	pageSizeHeader := c.Request().Header.Get("PageSize")
	pageSize, err := strconv.Atoi(pageSizeHeader)

	if err != nil {
		logger.Info("Response to client", zap.Any("error", exceptions.PageSizeError))
		span.SetAttributes(attribute.String(exceptions.StatusBadRequest, err.Error()))
		return c.JSON(http.StatusBadRequest, responser.Response{
			Message: exceptions.PageSizeError,
			Success: false,
			Data:    "",
		})
	}

	validate := validator.New()
	if err := validate.Var(pageSize, "min=5,max=50"); err != nil {
		logger.Info("Response to client", zap.Any("error", exceptions.PageSizeError))
		span.SetAttributes(attribute.String(exceptions.StatusBadRequest, err.Error()))
		return c.JSON(http.StatusBadRequest, responser.Response{
			Message: exceptions.PageSizeError,
			Success: false,
			Data:    "",
		})
	}

	pageToken := c.Request().Header.Get("PageToken")

	data, err := h.us.GetShopList(ctx, pageSize, pageToken)

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
			logger.Info("Response to client", zap.Any("error", reqErr.ErrMsg))
			span.SetAttributes(attribute.String(exceptions.StatusForbidden, reqErr.ErrMsg))
			return c.JSON(http.StatusBadRequest, responser.Response{
				Message: reqErr.ErrMsg,
				Data:    reqErr.Response,
				Success: false,
			})
		}
		logger.Info("Response to client", zap.Any("error", err.Error()))
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
		Data:    data["data"],
	}
	logger.Info("Response to client", zap.Any("data", response))
	span.SetAttributes(attribute.String("Response", utils.Marshal(response)))
	return c.JSON(http.StatusOK, response)
}
