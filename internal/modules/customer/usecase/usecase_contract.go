package usecase

import (
	"context"
	"giftCard/internal/adaptor/giftcard"
)

type ICustomerUseCase interface {
	GetCustomerInfoUseCase(ctx context.Context) (giftcard.CustomerInfoResponse, error)
}
