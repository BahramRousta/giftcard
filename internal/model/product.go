package model

import "gorm.io/gorm"

type Product struct {
	gorm.Model
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
