package giftcard

import (
	"context"
	"encoding/json"
	"fmt"
	"giftcard/internal/adaptor/trace"
	"go.opentelemetry.io/otel/attribute"
)

// RetrieveOrder use to get latest status of order
func (g *GiftCard) RetrieveOrder(ctx context.Context, orderId string) (map[string]any, error) {
	span, spannedContext := trace.T.SpanFromContext(
		ctx,
		"RetrieveOrderAdapter",
		"GiftCardAdapter",
	)
	defer span.End()

	url := g.BaseUrl + fmt.Sprintf("/order/get?orderId=%s", orderId)
	method := "GET"

	data, err := g.ProcessRequest(spannedContext, method, url, nil)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		return nil, err
	}

	jsonData, err := json.Marshal(data)
	span.SetAttributes(attribute.String("data", string(jsonData)))
	return data, nil
}

type RetrieveOrderError struct {
	ErrorMsg string
	Response map[string]interface{}
}

func (e *RetrieveOrderError) Error() string {
	return e.ErrorMsg
}
