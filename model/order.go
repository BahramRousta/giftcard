package model

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	gorm.Model
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
