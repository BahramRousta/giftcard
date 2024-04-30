package model

import (
	"gorm.io/gorm"
	"time"
)

type ExchangeRate struct {
	gorm.Model
	ID             uint      `gorm:"primaryKey"`
	BaseCurrency   string    `gorm:"column:base_currency;type:varchar(3);not null"`
	TargetCurrency string    `gorm:"column:target_currency;type:varchar(3);not null"`
	Rate           float64   `gorm:"column:rate;type:numeric;not null"`
	ModifiedDate   time.Time `gorm:"column:modified_date;not null"`
}

// TableName sets the table name for the ExchangeRate model
func (ExchangeRate) TableName() string {
	return "ExchangeRates"
}
