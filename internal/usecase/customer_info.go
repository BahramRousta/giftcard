package usecase

import adaptor "giftCard/internal/adaptor/giftcard"

func CustomerInfoUseCase() (map[string]any, error) {
	gf := adaptor.NewGiftCard()
	data, err := gf.CustomerInfo()
	if err != nil {
		return map[string]any{}, err
	}
	return data, nil
}
