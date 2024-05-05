package usecase

import (
	"context"
	"giftCard/internal/adaptor/giftcard"
)

type IOrderUseCase interface {
	GetOrderStatus(ctx context.Context, orderId string) (map[string]any, error)
	CreateOrder(ctx context.Context, productList []map[string]any) (giftcard.OrderResponse, error)
	ConfirmOrder(ctx context.Context, orderId string) (map[string]any, error)
}
