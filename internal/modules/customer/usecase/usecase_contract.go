package usecase

import (
	"context"
	"giftcard/internal/adaptor/giftcard"
)

type ICustomerUseCase interface {
	GetCustomerInfoUseCase(ctx context.Context) (giftcard.CustomerInfoResponse, error)
}
