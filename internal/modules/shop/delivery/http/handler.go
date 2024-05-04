package delivery

import (
	"errors"
	"fmt"
	gftErr "giftCard/internal/adaptor/giftcard"
	"giftCard/internal/modules/shop/usecase"
	"giftCard/pkg/responser"
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
			return c.JSON(http.StatusForbidden, responser.Response{
				Message: forbiddenErr.ErrMsg,
				Data:    "",
				Success: false,
			})
		}
		var reqErr *gftErr.RequestErr
		if errors.As(err, &reqErr) {
			return c.JSON(http.StatusBadRequest, responser.Response{
				Message: reqErr.ErrMsg,
				Data:    reqErr.Response,
				Success: false,
			})
		}
		return c.JSON(http.StatusInternalServerError, responser.Response{
			Message: "something went wrong",
			Data:    "",
			Success: false,
		})
	}
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

	data, err := h.us.GetShopList(pageSize, pageToken)
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
			return c.JSON(http.StatusBadRequest, responser.Response{
				Message: reqErr.ErrMsg,
				Data:    reqErr.Response,
				Success: false,
			})
		}
		return c.JSON(http.StatusInternalServerError, responser.Response{
			Data:    "",
			Message: "Something went wrong",
			Success: false,
		})
	}
	return c.JSON(http.StatusOK, responser.Response{Data: data["data"], Message: "", Success: true})
}
