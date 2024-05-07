package giftcard

import (
	"context"
	"encoding/json"
	"giftcard/internal/adaptor/trace"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

func (g *GiftCard) ConfirmOrder(ctx context.Context, orderId string) (map[string]any, error) {
	span, spannedContext := trace.T.SpanFromContext(
		ctx,
		"ConfirmOrderAdapter",
		"GiftCardAdapter",
	)
	defer span.End()

	uniqueID, _ := ctx.Value("tracer").(string)
	logger := zap.L().With(
		zap.String("tracer", uniqueID),
	)

	url := g.BaseUrl + "/order/confirm"
	method := "PUT"

	payload := map[string]any{
		"orderId": orderId,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	data, err := g.ProcessRequest(spannedContext, method, url, &payloadBytes)

	if err != nil {
		logger.Error("error while processing request to gift card provider",
			zap.String("error", err.Error()),
		)
		span.SetAttributes(attribute.String("error", err.Error()))
		return nil, err
	}

	jsonData, err := json.Marshal(data)
	span.SetAttributes(attribute.String("data", string(jsonData)))
	return data, nil
}

type ConfirmOrderError struct {
	ErrorMsg string
	Response map[string]interface{}
}

func (e *ConfirmOrderError) Error() string {
	return e.ErrorMsg
}
