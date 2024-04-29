package config

import (
	"fmt"
	"giftCard/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var database *gorm.DB
var e error

func DatabaseInit() {
	host := "localhost"
	user := "root"
	password := "password"
	dbName := "gift_card_db"
	port := 5433

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta", host, user, password, dbName, port)
	database, e = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	database.AutoMigrate(&model.Wallet{})
	database.AutoMigrate(&model.ExchangeRate{})
	database.AutoMigrate(&model.Product{}, &model.Variant{})
	database.AutoMigrate(&model.Order{})

	if e != nil {
		panic(e)
	}
}

func DB() *gorm.DB {
	return database
}
