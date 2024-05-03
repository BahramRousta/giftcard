package model

import (
	"gorm.io/gorm"
	"time"
)

type Wallet struct {
	gorm.Model
	ID            uint      `gorm:"primaryKey"`
	Currency      string    `gorm:"column:currency;type:varchar(3);not null"`
	Balance       float64   `gorm:"column:balance;type:numeric;not null"`
	CreditBalance float64   `gorm:"column:credit_balance;type:numeric;not null"`
	FrozenBalance float64   `gorm:"column:frozen_balance;type:numeric;not null"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

func (Wallet) TableName() string {
	return "Wallet"
}
