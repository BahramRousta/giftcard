package server

import "giftCard/cmd/api"

func (s *Server) SetupRoutes(customerHandler *api.CustomerInfoHandler,
	shopItemHandler *api.ShopItemHandler,
	shopListHandler *api.ShopListHandler,
) {
	s.GET("/customer/info", customerHandler.CustomerInfo)
	s.GET("/shop/products", shopListHandler.ShopList)
	s.GET("/shop/product/info", shopItemHandler.ShopItem)
	s.POST("/order/create", api.CreateOrder)
	s.POST("/order/confirm", api.ConfirmOrder)
	s.GET("/order/get", api.RetrieveOrder)
}
