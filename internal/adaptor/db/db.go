package db

import (
	"fmt"
	"giftCard/config"
	"giftCard/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var database *gorm.DB
var e error

func init() {
	config.C()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		config.C().DataBase.Host,
		config.C().DataBase.Username,
		config.C().DataBase.Password,
		config.C().DataBase.Schema,
		config.C().DataBase.Port,
		config.C().DataBase.SSLMode,
		config.C().DataBase.TimeZone,
	)

	database, e = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if e != nil {
		log.Fatal("failed to connect db")
	}

	migrateModels()
}

func migrateModels() {
	err := database.AutoMigrate(
		&model.Wallet{},
		&model.ExchangeRate{},
		&model.Variant{},
		&model.Product{},
		&model.Order{},
	)
	if err != nil {
		log.Fatalf("failed to migrate model")
	}
}
func DB() *gorm.DB {
	return database
}
