package model

import "time"

type Wallet struct {
	ID            uint      `gorm:"primaryKey"`
	Currency      string    `gorm:"column:currency;type:varchar(3);not null"`
	Balance       float64   `gorm:"column:balance;type:numeric;not null"`
	CreditBalance float64   `gorm:"column:credit_balance;type:numeric;not null"`
	FrozenBalance float64   `gorm:"column:frozen_balance;type:numeric;not null"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

// TableName sets the table name for the Wallet model
func (Wallet) TableName() string {
	return "Wallet"
}
