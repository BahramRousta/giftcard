package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func CustomerInfo(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, Echo!")
}
