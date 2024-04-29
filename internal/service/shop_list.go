package service

import (
	adaptor "giftCard/internal/adaptor/giftcard"
)

type ShopListService struct {
}

func NewShopListService() *ShopListService {
	return &ShopListService{}
}

func (s *ShopListService) GetShopListService() (map[string]any, error) {
	gf := adaptor.NewGiftCard()
	data, err := gf.ShopList()
	if err != nil {
		return map[string]any{}, err
	}
	return data, nil
}
