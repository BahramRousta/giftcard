package repository

import "giftcard/model"

type IOrderRepository interface {
	InsertOrder(order *model.Order) error
	GetOrder(orderId string) (*model.Order, error)
	UpdateOrder(order *model.Order, newStatus string) error
}
