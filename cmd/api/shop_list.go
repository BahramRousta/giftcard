package api

import (
	"errors"
	"fmt"
	gftErr "giftCard/internal/adaptor/gft_error"
	"giftCard/internal/service"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type ShopListHandler struct {
	shopListService *service.ShopListService
}

func NewShopListHandler(service *service.ShopListService) *ShopListHandler {
	return &ShopListHandler{
		shopListService: service,
	}
}

type RequestParams struct {
	PageSize  int    `json:"pageSize" form:"pageSize" validate:"omitempty,min=5,max=50"`
	PageToken string `json:"pageToken"`
}

func (h *ShopListHandler) ShopList(c echo.Context) error {

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

	data, err := h.shopListService.GetShopListService(pageSize, pageToken)
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
