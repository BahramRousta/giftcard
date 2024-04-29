package server

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type Server struct {
	*echo.Echo
}

func loggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Log incoming request
		log.Printf("Incoming Request: method:%s uri:%s body:%s", c.Request().Method, c.Request().URL.String(), c.Request().Body)

		// Call next handler
		if err := next(c); err != nil {
			c.Error(err)
		}

		// Log outgoing response
		log.Printf("Outgoing Response: status_code:%d status_message:%s", c.Response().Status, http.StatusText(c.Response().Status))

		return nil
	}
}

func NewServer() *Server {
	e := echo.New()
	e.Use(loggerMiddleware)
	return &Server{e}
}
