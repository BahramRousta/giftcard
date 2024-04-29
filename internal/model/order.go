package model

import "time"

type Order struct {
	ID          uint `gorm:"primaryKey"`
	SKU         string
	OrderID     string
	ProductType string
	Quote       uint
	Quantity    uint
	CreatedAt   time.Time
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
