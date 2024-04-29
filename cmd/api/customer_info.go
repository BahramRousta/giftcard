package api

import (
	"errors"
	"fmt"
	adaptor "giftCard/internal/adaptor/giftcard"
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
	fmt.Println("data", data)
	if err != nil {
		if err != nil {
			var customerErr *adaptor.CustomerInfoError
			if errors.As(err, &customerErr) {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"message": customerErr.Response,
					"success": false,
					"data":    map[string]any{},
				})
			}
			return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{"error": "Something went wrong"})
		}
	}
	return c.JSON(http.StatusCreated, map[string]any{"data": data, "message": "", "success": true})
}
