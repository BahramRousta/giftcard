package api

import (
	"errors"
	"fmt"
	adaptor "giftCard/internal/adaptor/giftcard"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type confirmOrderRequestBody struct {
	OrderId string `json:"orderId" validate:"required"`
}

func ConfirmOrder(c echo.Context) error {
	var requestBody confirmOrderRequestBody
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad Request"})
	}
	validate := validator.New()
	if err := validate.Struct(&requestBody); err != nil {
		var messages []string
		for _, fieldErr := range err.(validator.ValidationErrors) {
			messages = append(messages, fieldErr.Field()+" is invalid")
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": strings.Join(messages, ", ")})
	}
	fmt.Println("req.body again", requestBody.OrderId)
	gf := adaptor.NewGiftCard()
	data, err := gf.ConfirmOrder(requestBody.OrderId)
	if err != nil {
		var orderErr *adaptor.ConfirmOrderError
		if errors.As(err, &orderErr) {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error":   orderErr.ErrorMsg,
				"payload": orderErr.Response,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Something went wrong"})
	}
	return c.JSON(http.StatusCreated, map[string]any{"data": data["data"], "message": "", "success": true})
}
