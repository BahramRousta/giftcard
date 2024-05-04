package usecase

import "giftCard/internal/adaptor/giftcard"

type IShopUseCase interface {
	GetShopList(pageSize int, pageToken string) (map[string]any, error)
	GetShopItem(productId string) (giftcard.ProductResponse, error)
}
