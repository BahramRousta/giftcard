package usecase

import adaptor "giftCard/internal/adaptor/giftcard"

func ShopItemUseCase(productId string) (map[string]any, error) {
	gf := adaptor.NewGiftCard()
	data, err := gf.ShopItem(productId)
	if err != nil {
		return map[string]any{}, err
	}
	return data, nil
}
