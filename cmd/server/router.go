package server

import "giftCard/cmd/api"

func (s *Server) SetupRoutes(customerHandler *api.CustomerInfoHandler) {
	s.GET("/customer/info", customerHandler.CustomerInfo)
	s.GET("/shop/products", api.ShopList)
	s.GET("/shop/product/info", api.ShopItem)
	s.POST("/order/create", api.CreateOrder)
	s.POST("/order/confirm", api.ConfirmOrder)
	s.GET("/order/get", api.RetrieveOrder)
}
