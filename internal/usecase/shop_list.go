package usecase

import adaptor "giftCard/internal/adaptor/giftcard"

func ShopListUseCase() (map[string]any, error) {
	gf := adaptor.NewGiftCard()
	data, err := gf.ShopList()
	if err != nil {
		return map[string]any{}, err
	}
	return data, nil
}
