package model

import "gorm.io/gorm"

type Variant struct {
	gorm.Model
	ID           uint `gorm:"primaryKey"`
	ProductID    uint
	VariantID    string
	VariantName  string
	VariantQuote float64
	VariantSKU   string
	VariantState string
}
