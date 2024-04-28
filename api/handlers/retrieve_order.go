package handlers

import (
	adaptor "giftCard/internal/adaptor/giftcard"
	"github.com/labstack/echo/v4"
	"net/http"
)

type retrieveOrderRequest struct {
	OrderId string `json:"orderId" query:"orderId" validate:"required"`
}

func RetrieveOrder(c echo.Context) error {

	var request retrieveOrderRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad Request"})
	}

	orderId := request.OrderId
	if orderId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Order ID is required"})
	}

	gf := adaptor.NewGiftCard()
	data, err := gf.RetrieveOrder(orderId)
	if err != nil {
		if retrieveErr, ok := err.(*adaptor.RetrieveOrderError); ok {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"error":   retrieveErr.ErrorMsg,
				"message": retrieveErr.Response,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, data)
}
