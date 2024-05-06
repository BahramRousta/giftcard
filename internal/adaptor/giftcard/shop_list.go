package giftcard

import (
	"context"
	"fmt"
)

func (g *GiftCard) ShopList(ctx context.Context, pageSize int, pageToken string) (map[string]any, error) {
	url := fmt.Sprintf("%s/shop/products?pageSize=%d&pageToken=%s", g.BaseUrl, pageSize, pageToken)
	method := "GET"

	data, err := g.ProcessRequest(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type ShopListError struct {
	ErrorMsg string
	Response map[string]any
}

func (e *ShopListError) Error() string {
	return e.ErrorMsg
}
