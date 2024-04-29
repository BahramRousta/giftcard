package main

import (
	"giftCard/cmd/api"
	"giftCard/cmd/server"
	"giftCard/config"
	"giftCard/internal/adaptor/redis"
	"giftCard/internal/repository"
	"giftCard/internal/service"
)

func main() {
	server := server.NewServer()

	config.DatabaseInit()
	gorm := config.DB()

	dbGorm, err := gorm.DB()
	if err != nil {
		panic(err)
	}

	dbGorm.Ping()

	// initial redis
	//redisConfig := config.Redis{
	//	Password: "password",
	//	DB:       0,
	//	Host:     "localhost",
	//	Port:     6379,
	//}

	redis.RedisInit()

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
