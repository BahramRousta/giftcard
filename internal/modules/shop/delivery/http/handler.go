package delivery

import (
	"errors"
	"fmt"
	gftErr "giftCard/internal/adaptor/giftcard"
	"giftCard/internal/modules/shop/usecase"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
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
	productId := c.QueryParam("productId")
	if productId == "" {
		return c.String(http.StatusBadRequest, "productId is required")
	}

	data, err := h.us.GetShopItem(productId)
	if err != nil {
		var forbiddenErr *gftErr.ForbiddenErr
		if errors.As(err, &forbiddenErr) {
			return c.JSON(http.StatusForbidden, map[string]any{
				"message": forbiddenErr.ErrMsg,
				"data":    "",
				"success": false,
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
		return c.JSON(http.StatusInternalServerError, map[string]any{"data": "", "message": err.Error(), "success": false})

	}
	return c.JSON(http.StatusOK, map[string]any{"data": data, "message": "", "success": true})
}

type RequestParams struct {
	PageSize  int    `json:"pageSize" form:"pageSize" validate:"omitempty,min=5,max=50"`
	PageToken string `json:"pageToken"`
}

func (h *ShopHandler) ShopList(c echo.Context) error {

	pageSizeHeader := c.Request().Header.Get("PageSize")
	pageSize, err := strconv.Atoi(pageSizeHeader)

	if err != nil {
		c.Logger().Error("Error in page size", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "pageSize must be an integer",
			"success": false,
			"data":    map[string]interface{}{},
		})
	}

	validate := validator.New()
	if err := validate.Var(pageSize, "required,min=5,max=50"); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": fmt.Sprintf("Validation error: %s", err.Error()),
			"success": false,
			"data":    map[string]interface{}{},
		})
	}

	pageToken := c.Request().Header.Get("PageToken")

	data, err := h.us.GetShopList(pageSize, pageToken)
	if err != nil {
		var forbiddenErr *gftErr.ForbiddenErr
		if errors.As(err, &forbiddenErr) {
			return c.JSON(http.StatusForbidden, map[string]any{
				"message": forbiddenErr.ErrMsg,
				"data":    "",
				"success": false,
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
		return c.JSON(http.StatusInternalServerError, map[string]any{"data": "", "message": "Something went wrong", "success": false})
	}
	return c.JSON(http.StatusOK, map[string]any{"data": data["data"], "message": "", "success": true})
}
