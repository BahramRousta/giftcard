package service

import (
	adaptor "giftCard/internal/adaptor/giftcard"
	"giftCard/internal/model"
	"giftCard/internal/repository"
)

type ShopItemService struct {
	productRepo *repository.ProductRepository
	variantRepo *repository.VariantRepository
}

func NewShopItemService(productRepo *repository.ProductRepository, variantRepo *repository.VariantRepository) *ShopItemService {
	return &ShopItemService{
		productRepo: productRepo,
		variantRepo: variantRepo,
	}
}

func (s *ShopItemService) GetShopItemService(productId string) (adaptor.ProductResponse, error) {
	gf := adaptor.NewGiftCard()
	data, err := gf.ShopItem(productId)
	if err != nil {
		return adaptor.ProductResponse{}, err
	}

	product := &model.Product{
		BaseCurrency:                       data.Data.BaseCurrency,
		Country:                            data.Data.Country,
		CustomerDealDiscountAdjustmentMode: data.Data.CustomerDeal.Discount.AdjustmentMode,
		CustomerDealDiscountAmount:         float64(data.Data.CustomerDeal.Discount.Amount),
		CustomerDealFeeAmount:              float64(data.Data.CustomerDeal.Fee.Amount),
		Name:                               data.Data.Name,
		ProductID:                          data.Data.ProductID,
		ProductType:                        data.Data.ProductType,
		ParentID:                           data.Data.ParentID,
		Region:                             data.Data.Region,
	}
	if err := s.productRepo.InsertProduct(product); err != nil {
		return adaptor.ProductResponse{}, err
	}

	for _, v := range data.Data.Variants {
		variant := &model.Variant{
			ProductID:    product.ID,
			VariantName:  v.Name,
			VariantSKU:   v.SKU,
			VariantQuote: float64(v.Quote),
			VariantState: v.State,
		}
		if err := s.variantRepo.InsertVariant(variant); err != nil {
			return adaptor.ProductResponse{}, err
		}
	}
	return data, nil
}
