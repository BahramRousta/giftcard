package config

import (
	"fmt"
	"giftCard/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var database *gorm.DB
var e error

func DatabaseInit() {

	config, err := LoadConfig()
	if err != nil {
		log.Fatal("cannot load config")
	}

	host := config.DataBase.Host
	user := config.DataBase.Username
	password := config.DataBase.Password
	dbName := config.DataBase.Schema
	port := config.DataBase.Port
	sslMode := config.DataBase.SSLMode
	timeZone := config.DataBase.TimeZone

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s", host, user, password, dbName, port, sslMode, timeZone)
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
