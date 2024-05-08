package server

import (
	"context"
	"fmt"
	CustomerHttp "giftcard/internal/modules/customer/delivery/http"
	OrderHttp "giftcard/internal/modules/order/delivery/http"
	ShopHttp "giftcard/internal/modules/shop/delivery/http"
	"giftcard/internal/server/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
	"log"
	"net/http"
	"time"
)

type Server struct {
	srv       *echo.Echo
	container DeliveryContainer
}

func NewServer(p DeliveryContainer) IServer {
	server := Server{
		srv:       echo.New(),
		container: p,
	}
	return &server
}

func (s *Server) SetUpServer(container DeliveryContainer) {
	s.srv.Use(middleware.RequestID())
	v1 := s.srv.Group("/v1")
	routes.MapShopHandler(v1, container.ShopHandler)
	routes.MapCustomerHandler(v1, container.CustomerHandler)
	routes.MapOrderHandler(v1, container.OrderHandler)

	s.srv.GET("/health", func(c echo.Context) error {
		return c.String(200, fmt.Sprintf("Hi :)) i'm in healthy"))
	})

}

func (s *Server) Run() error {
	server := &http.Server{
		Addr: ":8080",
	}

	// invoice-go
	go func() {
		s.SetUpServer(s.container)

		//log.Printf("Server is listening on PORT: %s", config.C().Service.Server.Port)
		if err := s.srv.StartServer(server); err != nil {
			panic("error in starting http " + err.Error())
		}
	}()
	return nil

}

func (s *Server) Shutdown() error {
	log.Println("Shutting down............")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return s.srv.Server.Shutdown(ctx)
}

type DeliveryContainer struct {
	fx.In
	ShopHandler     *ShopHttp.ShopHandler
	OrderHandler    *OrderHttp.OrderHandler
	CustomerHandler *CustomerHttp.CustomerInfoHandler
}
