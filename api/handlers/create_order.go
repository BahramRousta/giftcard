package handlers

import (
	adaptor "giftCard/internal/adaptor/giftcard"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type Product struct {
	Sku         string `json:"sku" validate:"required"`
	ProductType string `json:"productType" validate:"required"`
	Quote       int    `json:"quote" validate:"required,gt=0"`
	Quantity    int    `json:"quantity" validate:"required,gt=0"`
}

type RequestBody struct {
	ProductList []Product `json:"productList" validate:"required,dive"`
}

func CreateOrder(c echo.Context) error {
	var requestBody RequestBody
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
	if len(requestBody.ProductList) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ProductList cannot be empty"})
	}

	var productList []map[string]interface{}
	for _, product := range requestBody.ProductList {
		productMap := map[string]interface{}{
			"sku":         product.Sku,
			"productType": product.ProductType,
			"quote":       product.Quote,
			"quantity":    product.Quantity,
		}
		productList = append(productList, productMap)
	}

	gf := adaptor.NewGiftCard()
	data, err := gf.CreateOrder(productList)
	if err != nil {
		if orderErr, ok := err.(*adaptor.OrderError); ok {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error":   orderErr.ErrorMsg,
				"payload": orderErr.Response,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, data)
}
