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

func DatabaseInit(cfn *config.Config) {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfn.DataBase.Host,
		cfn.DataBase.Username,
		cfn.DataBase.Password,
		cfn.DataBase.Schema,
		cfn.DataBase.Port,
		cfn.DataBase.SSLMode,
		cfn.DataBase.TimeZone,
	)

	database, e = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if e != nil {
		log.Fatal("failed to connect db")
	}

	database.AutoMigrate(
		&model.Wallet{},
		&model.ExchangeRate{},
		&model.Variant{},
		&model.Product{},
		&model.Order{},
	)
}

func DB() *gorm.DB {
	return database
}
