package service

import (
	adaptor "giftCard/internal/adaptor/giftcard"
)

type ShopListService struct {
}

func NewShopListService() *ShopListService {
	return &ShopListService{}
}

func (s *ShopListService) GetShopListService(pageSize int, pageToken string) (map[string]any, error) {
	gf := adaptor.NewGiftCard()
	data, err := gf.ShopList(pageSize, pageToken)
	if err != nil {
		return map[string]any{}, err
	}
	return data, nil
}
