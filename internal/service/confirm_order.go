package service

import (
	"errors"
	"fmt"
	adaptor "giftCard/internal/adaptor/giftcard"
	"giftCard/internal/repository"
)

type ConfirmOrderService struct {
	repo *repository.OrderRepository
}

func NewConfirmOrderService(repo *repository.OrderRepository) *ConfirmOrderService {
	return &ConfirmOrderService{
		repo: repo,
	}
}

func (s *ConfirmOrderService) OrderConfirmService(orderId string) (map[string]any, error) {
	order, err := s.repo.GetOrder(orderId)
	if err != nil {
		return nil, errors.New("order not found")
	}

	gf := adaptor.NewGiftCard()
	data, err := gf.ConfirmOrder(orderId)

	if err != nil {
		return nil, err
	}

	dataMap, ok := data["data"].(map[string]interface{})
	if !ok {
		fmt.Println("Data not found")
		return nil, err
	}

	state, ok := dataMap["state"].(string)
	if !ok {
		fmt.Println("State not found")
		return nil, err
	}

	err = s.repo.UpdateOrder(order, state)
	if err != nil {
		return nil, err
	}
	return data, nil
}
