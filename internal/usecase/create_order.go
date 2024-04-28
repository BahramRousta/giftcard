package usecase

import (
	"giftCard/internal/adaptor/giftcard"
)

func CreateOrderUseCase(productList []map[string]any) (map[string]any, error) {
	gf := adaptor.NewGiftCard()
	data, err := gf.CreateOrder(productList)
	if err != nil {
		return map[string]any{}, err
	}
	return data, nil
}
