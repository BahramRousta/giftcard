package server

import (
	"context"
	"fmt"
	CustomerHttp "giftCard/internal/modules/customer/delivery/http"
	OrderHttp "giftCard/internal/modules/order/delivery/http"
	ShopHttp "giftCard/internal/modules/shop/delivery/http"
	"giftCard/internal/server/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
	"log"
	"net/http"
	"os"
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
	s.srv.Use(middleware.Logger())
	s.srv.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		logFile, err := os.OpenFile("server.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			s.srv.Logger.Error("Failed to open log file for body dump:", err)
			return
		}
		defer logFile.Close()

		currentTime := time.Now()

		if _, err := fmt.Fprintf(
			logFile,
			"[%s] - Method: %s, "+
				"URI: %s\n"+
				"Request Body: %s\n"+
				"Response Body: %s\n",
			currentTime.Format(time.RFC3339),
			c.Request().Method,
			c.Request().RequestURI,
			string(reqBody),
			string(resBody),
		); err != nil {
			s.srv.Logger.Error("Failed to write request and response payloads to log file:", err)
			return
		}
	}))
	v1 := s.srv.Group("/v1")
	routes.MapShopHandler(v1, container.ShopHandler)
	routes.MapCustomerHandler(v1, container.CustomerHandler)
	routes.MapOrderHandler(v1, *container.OrderHandler)

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
