package main

import (
	"giftCard/cmd/api"
	"giftCard/cmd/server"
	"giftCard/config"
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

	walletRepo := repository.NewWalletRepository(gorm)
	exchangeRepo := repository.NewExchangeRepository(gorm)
	customerService := service.NewCustomerService(walletRepo, exchangeRepo)
	customerHandler := api.NewCustomerInfoHandler(customerService)

	server.SetupRoutes(customerHandler)

	server.Logger.Fatal(server.Start(":8000"))

}
