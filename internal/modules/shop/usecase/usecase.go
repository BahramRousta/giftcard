package usecase

import (
	"giftCard/internal/adaptor/giftcard"
	"go.uber.org/fx"
)

type giftCardShopUseCase struct {
	GiftCard *giftcard.GiftCard
}

type GiftCardShopUseCaseParams struct {
	fx.In
	GiftCard *giftcard.GiftCard
}

func NewShopUseCase(params GiftCardShopUseCaseParams) IShopUseCase {
	return &giftCardShopUseCase{
		GiftCard: params.GiftCard,
	}
}

func (us giftCardShopUseCase) GetShopList(pageSize int, pageToken string) (map[string]any, error) {
	data, err := us.GiftCard.ShopList(pageSize, pageToken)
	if err != nil {
		return map[string]any{}, err
	}
	return data, nil
}

func (us giftCardShopUseCase) GetShopItem(productId string) (giftcard.ProductResponse, error) {
	data, err := us.GiftCard.ShopItem(productId)
	if err != nil {
		return giftcard.ProductResponse{}, err
	}
	return data, nil
}
