package delivery

import (
	"errors"
	gftErr "giftCard/internal/adaptor/giftcard"
	"giftCard/internal/modules/customer/usecase"
	"giftCard/pkg/responser"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
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
	data, err := h.customerUseCase.GetCustomerInfoService()

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
	return c.JSON(http.StatusOK, responser.Response{Message: "", Success: true, Data: data.Data})
}
