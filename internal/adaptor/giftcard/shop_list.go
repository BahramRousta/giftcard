package giftcard

import (
	"context"
	"encoding/json"
	"fmt"
	"giftcard/internal/adaptor/trace"
	"go.opentelemetry.io/otel/attribute"
)

func (g *GiftCard) ShopList(ctx context.Context, pageSize int, pageToken string) (map[string]any, error) {
	span, spannedContext := trace.T.SpanFromContext(
		ctx,
		"ShopListAdapter",
		"GiftCardAdapter",
	)
	defer span.End()

	url := fmt.Sprintf("%s/shop/products?pageSize=%d&pageToken=%s", g.BaseUrl, pageSize, pageToken)
	method := "GET"

	data, err := g.ProcessRequest(spannedContext, method, url, nil)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		return nil, err
	}

	jsonData, _ := json.Marshal(data)
	span.SetAttributes(attribute.String("data", string(jsonData)))
	return data, nil
}

type ShopListError struct {
	ErrorMsg string
	Response map[string]any
}

func (e *ShopListError) Error() string {
	return e.ErrorMsg
}
