package usecase

import (
	"context"
	"giftcard/internal/adaptor/giftcard"
	"giftcard/internal/adaptor/trace"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type giftCardShopUseCase struct {
	GiftCard *giftcard.GiftCard
}

type GiftCardShopUseCaseParams struct {
	fx.In
	GiftCard *giftcard.GiftCard
}

func NewShopUseCase(params GiftCardShopUseCaseParams) IShopUseCase {
	return &giftCardShopUseCase{
		GiftCard: params.GiftCard,
	}
}

func (us giftCardShopUseCase) GetShopList(ctx context.Context, pageSize int, pageToken string) (map[string]any, error) {
	span, spannedContext := trace.T.SpanFromContext(
		ctx,
		"GetShopListUseCase",
		"UseCase")
	defer span.End()

	uniqueID, _ := ctx.Value("tracer").(string)

	logger := zap.L().With(
		zap.String("tracer", uniqueID),
	)

	data, err := us.GiftCard.ShopList(spannedContext, pageSize, pageToken)
	if err != nil {
		logger.Error("error while processing gift card shop list",
			zap.String("error", err.Error()),
		)
		span.SetAttributes(attribute.String("error", err.Error()))
		return map[string]any{}, err
	}
	return data, nil
}

func (us giftCardShopUseCase) GetShopItem(ctx context.Context, productId string) (giftcard.ProductResponse, error) {
	span, spannedContext := trace.T.SpanFromContext(
		ctx,
		"GetShopItemUseCase",
		"UseCase")
	defer span.End()

	uniqueID, _ := ctx.Value("tracer").(string)

	logger := zap.L().With(
		zap.String("tracer", uniqueID),
	)

	data, err := us.GiftCard.ShopItem(spannedContext, productId)
	if err != nil {
		logger.Error("error while processing gift card shop item",
			zap.String("error", err.Error()),
		)
		span.SetAttributes(attribute.String("error", err.Error()))
		return giftcard.ProductResponse{}, err
	}
	return data, nil
}
