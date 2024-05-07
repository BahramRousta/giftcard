package postgres

import (
	"fmt"
	"giftcard/config"
	model2 "giftcard/model"
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
		log.Fatal("failed to connect postgres")
	}
}

func MigrateModels() error {
	err := database.AutoMigrate(
		&model2.Wallet{},
		&model2.Order{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate model, %w", err)
	}
	return nil
}
func DB() *gorm.DB {
	return database
}
