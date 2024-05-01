package main

import (
	"giftCard/cmd/api"
	"giftCard/cmd/server"
	"giftCard/internal/adaptor/db"
	"giftCard/internal/repository"
	"giftCard/internal/service"
)

func main() {
	server := server.NewServer()

	gorm := db.DB()

	dbGorm, err := gorm.DB()
	if err != nil {
		panic(err)
	}
	dbGorm.Ping()

	walletRepo := repository.NewWalletRepository(gorm)
	exchangeRepo := repository.NewExchangeRepository(gorm)
	customerService := service.NewCustomerService(walletRepo, exchangeRepo)
	customerHandler := api.NewCustomerInfoHandler(customerService)

	productRepo := repository.NewProductRepository(gorm)
	variantRepo := repository.NewVariantRepository(gorm)
	shopItemService := service.NewShopItemService(productRepo, variantRepo)
	shopItemHandler := api.NewShopItemHandler(shopItemService)

	shopListService := service.NewShopListService()
	shopListHandler := api.NewShopListHandler(shopListService)

	orderRepo := repository.NewOrderRepository(gorm)
	createOrderService := service.NewCreateOrderService(orderRepo)
	createOrderHandler := api.NewCreateOrderHandler(createOrderService)

	getOrderService := service.NewRetrieveOrderService(orderRepo)
	getOrderHandler := api.NewRetrieveOrderHandler(getOrderService)

	confirmOrderService := service.NewConfirmOrderService(orderRepo)
	confirmOrderHandler := api.NewConfirmOrderHandler(confirmOrderService)

	server.SetupRoutes(
		customerHandler,
		shopItemHandler,
		shopListHandler,
		createOrderHandler,
		getOrderHandler,
		confirmOrderHandler,
	)

	server.Logger.Fatal(server.Start(":8000"))

}
