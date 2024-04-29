package service

import (
	"fmt"
	adaptor "giftCard/internal/adaptor/giftcard"
	"giftCard/internal/repository"
)

type RetrieveOrderService struct {
	repo *repository.OrderRepository
}

func NewRetrieveOrderService(repo *repository.OrderRepository) *RetrieveOrderService {
	return &RetrieveOrderService{
		repo: repo,
	}
}

func (s *RetrieveOrderService) GetOrderStatusService(orderId string) (map[string]any, error) {

	order, err := s.repo.GetOrder(orderId)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, fmt.Errorf("order with ID %s not found", orderId) // Return error if order does not exist
	}

	gf := adaptor.NewGiftCard()
	data, err := gf.RetrieveOrder(orderId)
	if err != nil {
		return nil, err
	}
	return data, nil
}
