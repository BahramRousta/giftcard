package server

import "giftCard/cmd/api"

func (s *Server) SetupRoutes(
	customerHandler *api.CustomerInfoHandler,
	shopItemHandler *api.ShopItemHandler,
	shopListHandler *api.ShopListHandler,
	createOrderHandler *api.CreateOrderHandler,
	getOrderHandler *api.RetrieveOrderHandler,
	confirmOrderHandler *api.ConfirmOrderHandler,
) {
	s.GET("/customer/info", customerHandler.CustomerInfo)
	s.GET("/shop/products", shopListHandler.ShopList)
	s.GET("/shop/product/info", shopItemHandler.ShopItem)
	s.POST("/order/create", createOrderHandler.CreateOrder)
	s.POST("/order/confirm", confirmOrderHandler.ConfirmOrder)
	s.GET("/order/get", getOrderHandler.RetrieveOrder)
}
