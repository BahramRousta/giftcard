package api

import (
	"giftCard/api/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	*echo.Echo
}

func NewServer() *Server {
	e := echo.New()
	e.Use(middleware.Logger())
	return &Server{e}
}

func (s *Server) SetupRoutes() {
	s.GET("/customer/info", handlers.CustomerInfo)
	s.GET("/shop/products", handlers.ShopList)
	s.GET("/shop/product/info", handlers.ShopItem)
	s.POST("/order/create", handlers.CreateOrder)
	s.POST("/order/confirm", handlers.ConfirmOrder)
	s.GET("/order/get", handlers.RetrieveOrder)
}
