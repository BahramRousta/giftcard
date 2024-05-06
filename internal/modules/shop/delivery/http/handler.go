package delivery

import (
	"context"
	"errors"
	"fmt"
	gftErr "giftcard/internal/adaptor/giftcard"
	"giftcard/internal/adaptor/trace"
	"giftcard/internal/modules/shop/usecase"
	"giftcard/pkg/responser"
	"giftcard/pkg/utils"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type ShopHandler struct {
	us usecase.IShopUseCase
	//Logger *zap.Logger
}
type ShopHandlerParams struct {
	fx.In
	Us usecase.IShopUseCase
	//Logger *zap.Logger
}

func NewShopHandler(params ShopHandlerParams) *ShopHandler {
	return &ShopHandler{
		us: params.Us,
		//Logger: params.Logger,
	}
}

func (h *ShopHandler) ShopItem(c echo.Context) error {
	span, spannedContext := trace.T.SpanFromContext(
		utils.GetRequestCtx(c),
		"CustomerInfo[CustomerDelivery]",
		"delivery")
	defer span.End()

	productId := c.QueryParam("productId")
	if productId == "" {
		return c.String(http.StatusBadRequest, "productId is required")
	}

	uniqueID := uuid.New().String()

	logger := zap.L().With(
		zap.String("tracer", uniqueID),
	)
	logger.Info("shop item request",
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
			logger.Info("shop item response",
				zap.Any("data", forbiddenErr.ErrMsg),
			)
			return c.JSON(http.StatusForbidden, responser.Response{
				Message: forbiddenErr.ErrMsg,
				Data:    "",
				Success: false,
			})
		}
		var reqErr *gftErr.RequestErr
		if errors.As(err, &reqErr) {
			logger.Info("shop item response",
				zap.Any("error", reqErr.ErrMsg),
				zap.Any("data", reqErr.Response),
			)
			return c.JSON(http.StatusBadRequest, responser.Response{
				Message: reqErr.ErrMsg,
				Data:    reqErr.Response,
				Success: false,
			})
		}
		logger.Info("shop item response",
			zap.Any("error", "internal error"),
		)
		return c.JSON(http.StatusInternalServerError, responser.Response{
			Message: "something went wrong",
			Data:    "",
			Success: false,
		})
	}
	logger.Info("shop item response",
		zap.Any("data", data.Data),
	)
	span.SetAttributes(attribute.String("response", ""))
	return c.JSON(http.StatusOK, responser.Response{
		Message: "",
		Data:    data.Data,
		Success: true,
	})
}

func (h *ShopHandler) ShopList(c echo.Context) error {

	pageSizeHeader := c.Request().Header.Get("PageSize")
	pageSize, err := strconv.Atoi(pageSizeHeader)

	if err != nil {
		return c.JSON(http.StatusBadRequest, responser.Response{
			Message: "pageSize must be an integer",
			Success: false,
			Data:    "",
		})
	}

	validate := validator.New()
	if err := validate.Var(pageSize, "min=5,max=50"); err != nil {
		return c.JSON(http.StatusBadRequest, responser.Response{
			Message: fmt.Sprintf("page size must between 5 to 50."),
			Success: false,
			Data:    "",
		})
	}

	pageToken := c.Request().Header.Get("PageToken")
	uniqueID := uuid.New().String()
	ctx := context.WithValue(c.Request().Context(), "tracer", uniqueID)

	logger := zap.L().With(
		zap.String("tracer", uniqueID),
	)
	logger.Info("shop list request",
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
			logger.Info("shop list response",
				zap.Any("data", forbiddenErr.ErrMsg),
			)

			return c.JSON(http.StatusForbidden, responser.Response{
				Message: forbiddenErr.ErrMsg,
				Data:    "",
				Success: false,
			})
		}
		var reqErr *gftErr.RequestErr
		if errors.As(err, &reqErr) {
			logger.Info("shop list response",
				zap.Any("error", reqErr.ErrMsg),
				zap.Any("data", reqErr.Response),
			)

			return c.JSON(http.StatusBadRequest, responser.Response{
				Message: reqErr.ErrMsg,
				Data:    reqErr.Response,
				Success: false,
			})
		}
		logger.Info("shop list response",
			zap.Any("error", reqErr.ErrMsg),
			zap.Any("data", reqErr.Response),
		)
		return c.JSON(http.StatusInternalServerError, responser.Response{
			Data:    "",
			Message: "Something went wrong",
			Success: false,
		})
	}
	logger.Info("shop list response",
		zap.Any("data", data["data"]),
	)
	return c.JSON(http.StatusOK, responser.Response{Data: data["data"], Message: "", Success: true})
}
