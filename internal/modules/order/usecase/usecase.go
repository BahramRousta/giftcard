package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"giftcard/internal/adaptor/giftcard"
	"giftcard/internal/adaptor/trace"
	"giftcard/internal/modules/order/repository"
	"giftcard/model"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/fx"
	"go.uber.org/zap"
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

func (us giftCardOrderUseCase) GetOrderStatus(ctx context.Context, orderId string) (map[string]any, error) {
	span, spannedContext := trace.T.SpanFromContext(
		ctx,
		"GetOrderStatusUseCase",
		"UseCase")
	defer span.End()

	uniqueID, _ := ctx.Value("tracer").(string)

	logger := zap.L().With(
		zap.String("tracer", uniqueID),
	)

	order, err := us.repo.GetOrder(orderId)
	if err != nil {
		logger.Error("error while get order from DB",
			zap.String("error", err.Error()),
		)
		span.SetAttributes(attribute.String("error", err.Error()))
		return nil, err
	}

	if order == nil {
		logger.Error("error while get order from DB",
			zap.String("error", fmt.Sprintf("order with ID %s not found", orderId)))
		span.SetAttributes(attribute.String("error", fmt.Sprintf("order with ID %s not found", orderId)))
		return nil, fmt.Errorf("order with ID %s not found", orderId)
	}

	data, err := us.gf.RetrieveOrder(spannedContext, orderId)
	if err != nil {
		logger.Error("error while processing gift card retrieve order",
			zap.String("error", err.Error()),
		)
		span.SetAttributes(attribute.String("error", err.Error()))
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
		logger.Error("error while update order status from DB",
			zap.String("error", err.Error()),
		)
		span.SetAttributes(attribute.String("error", err.Error()))
		return nil, err
	}
	jsonData, err := json.Marshal(data)
	span.SetAttributes(attribute.String("data", string(jsonData)))
	return data, nil
}

func (us giftCardOrderUseCase) CreateOrder(ctx context.Context, productList []map[string]any) (giftcard.OrderResponse, error) {
	span, spannedContext := trace.T.SpanFromContext(
		ctx,
		"CreateOrderUseCase",
		"UseCase")
	defer span.End()

	uniqueID, _ := ctx.Value("tracer").(string)

	logger := zap.L().With(
		zap.String("tracer", uniqueID),
	)

	data, err := us.gf.CreateOrder(spannedContext, productList)
	if err != nil {
		logger.Error("error while processing gift card create order",
			zap.String("error", err.Error()),
		)
		span.SetAttributes(attribute.String("error", err.Error()))
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
		logger.Error("error while insert new order from DB",
			zap.String("error", err.Error()),
		)
		span.SetAttributes(attribute.String("error", err.Error()))
		return giftcard.OrderResponse{}, err
	}

	jsonData, err := json.Marshal(data)
	span.SetAttributes(attribute.String("data", string(jsonData)))
	return data, nil
}
func (us giftCardOrderUseCase) ConfirmOrder(ctx context.Context, orderId string) (map[string]any, error) {
	span, spannedContext := trace.T.SpanFromContext(
		ctx,
		"GetOrderStatusUseCase",
		"UseCase")
	defer span.End()

	uniqueID, _ := ctx.Value("tracer").(string)

	logger := zap.L().With(
		zap.String("tracer", uniqueID),
	)

	order, err := us.repo.GetOrder(orderId)
	if err != nil {
		logger.Error("error while get order from DB",
			zap.String("error", err.Error()),
		)
		span.SetAttributes(attribute.String("error", err.Error()))
		return nil, err
	}

	data, err := us.gf.ConfirmOrder(spannedContext, orderId)
	if err != nil {
		logger.Error("error while processing gift card confirm order",
			zap.String("error", err.Error()),
		)
		span.SetAttributes(attribute.String("error", err.Error()))
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
		logger.Error("error while update order status from DB",
			zap.String("error", err.Error()),
		)
		span.SetAttributes(attribute.String("error", err.Error()))
		return nil, err
	}

	jsonData, err := json.Marshal(data)
	span.SetAttributes(attribute.String("data", string(jsonData)))
	return data, nil
}
