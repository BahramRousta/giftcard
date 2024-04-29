package service

import (
	"fmt"
	"giftCard/internal/adaptor/giftcard"
	"giftCard/internal/model"
	"giftCard/internal/repository"
	"time"
)

type OrderService struct {
	repo *repository.OrderRepository
}

func NewCreateOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

func (s *OrderService) CreateOrderService(productList []map[string]any) (adaptor.OrderResponse, error) {
	gf := adaptor.NewGiftCard()
	data, err := gf.CreateOrder(productList)
	if err != nil {
		return adaptor.OrderResponse{}, err
	}
	fmt.Println("data", data)
	fmt.Println("productList", productList)

	order := &model.Order{
		OrderID:     data.Data.ID,
		SKU:         productList[0]["sku"].(string),
		ProductType: productList[0]["productType"].(string),
		Quote:       productList[0]["quote"].(uint),
		Quantity:    productList[0]["quantity"].(uint),
		Status:      data.Data.Invoice.Status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	s.repo.InsertOrder(order)
	return data, nil
}
