package usecase

import "giftCard/internal/adaptor/giftcard"

type IOrderUseCase interface {
	GetOrderStatus(orderId string) (map[string]any, error)
	CreateOrder(productList []map[string]any) (giftcard.OrderResponse, error)
	ConfirmOrder(orderId string) (map[string]any, error)
}
