package usecase

import "giftCard/internal/adaptor/giftcard"

type ICustomerUseCase interface {
	GetCustomerInfoService() (giftcard.CustomerInfoResponse, error)
}
