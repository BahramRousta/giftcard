package api

import (
	"errors"
	gftErr "giftCard/internal/adaptor/gft_error"
	"giftCard/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CustomerInfoHandler struct {
	customerService *service.CustomerService
}

func NewCustomerInfoHandler(s *service.CustomerService) *CustomerInfoHandler {
	return &CustomerInfoHandler{
		customerService: s,
	}
}

func (h CustomerInfoHandler) CustomerInfo(c echo.Context) error {
	data, err := h.customerService.GetCustomerInfoService()

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
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{"data": "", "message": "something went wrong", "success": false})
	}
	return c.JSON(http.StatusOK, map[string]any{"data": data.Data, "message": "", "success": true})
}
