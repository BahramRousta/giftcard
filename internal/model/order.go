package model

import (
	"time"
)

type Order struct {
	ID          uint `gorm:"primaryKey"`
	SKU         string
	OrderID     string
	ProductType string
	Quote       uint
	Quantity    uint
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
