package usecase

import (
	"fmt"
	"giftCard/internal/adaptor/giftcard"
	"giftCard/internal/modules/order/repository"
	"giftCard/model"
	"go.uber.org/fx"
	"time"
)

type giftCardOrderUseCase struct {
	repo repository.IOrderRepository
	gf   giftcard.IGiftCard
}

type GiftCardOrderUseCaseParams struct {
	fx.In
	Repo repository.IOrderRepository
	Gf   *giftcard.GiftCard
}

func NewOrderUseCase(params GiftCardOrderUseCaseParams) IOrderUseCase {
	return &giftCardOrderUseCase{
		repo: params.Repo,
		gf:   params.Gf,
	}
}

func (us giftCardOrderUseCase) GetOrderStatus(orderId string) (map[string]any, error) {
	order, err := us.repo.GetOrder(orderId)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, fmt.Errorf("order with ID %s not found", orderId)
	}

	data, err := us.gf.RetrieveOrder(orderId)
	if err != nil {
		return nil, err
	}

	invoice, ok := data["data"].(map[string]interface{})["invoice"].(map[string]interface{})
	if !ok {
		return nil, err
	}

	status, ok := invoice["status"].(string)
	if !ok {
		return nil, err
	}
	if err := us.repo.UpdateOrder(order, status); err != nil {
		return nil, err
	}

	return data, nil
}
func (us giftCardOrderUseCase) CreateOrder(productList []map[string]any) (giftcard.OrderResponse, error) {
	data, err := us.gf.CreateOrder(productList)
	if err != nil {
		return giftcard.OrderResponse{}, err
	}

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
	err = us.repo.InsertOrder(order)
	if err != nil {
		return giftcard.OrderResponse{}, err
	}
	return data, nil
}
func (us giftCardOrderUseCase) ConfirmOrder(orderId string) (map[string]any, error) {
	order, err := us.repo.GetOrder(orderId)
	if err != nil {
		return nil, err
	}

	data, err := us.gf.ConfirmOrder(orderId)

	if err != nil {
		return nil, err
	}

	dataMap, ok := data["data"].(map[string]interface{})
	if !ok {
		return nil, err
	}

	state, ok := dataMap["state"].(string)
	if !ok {
		return nil, err
	}

	err = us.repo.UpdateOrder(order, state)
	if err != nil {
		return nil, err
	}
	return data, nil
}
