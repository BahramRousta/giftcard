package model

type Product struct {
	ID                                 uint `gorm:"primaryKey"`
	BaseCurrency                       string
	Country                            string
	CustomerDealDiscountAdjustmentMode string
	CustomerDealDiscountAmount         float64
	CustomerDealFeeAdjustmentMode      string
	CustomerDealFeeAmount              float64
	Name                               string
	ParentID                           string
	ProductID                          string
	ProductType                        string
	Region                             string
	Variants                           []Variant
}
