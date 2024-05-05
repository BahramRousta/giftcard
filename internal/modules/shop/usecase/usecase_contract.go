package usecase

import (
	"context"
	"giftCard/internal/adaptor/giftcard"
)

type IShopUseCase interface {
	GetShopList(ctx context.Context, pageSize int, pageToken string) (map[string]any, error)
	GetShopItem(ctx context.Context, productId string) (giftcard.ProductResponse, error)
}
