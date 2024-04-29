package model

type Variant struct {
	ID           uint `gorm:"primaryKey"`
	ProductID    uint
	VariantID    string
	VariantName  string
	VariantQuote float64
	VariantSKU   string
	VariantState string
}
