package giftcard

import "context"

type IGiftCard interface {
	ProcessRequest(ctx context.Context, method string, url string, payload *[]byte) (map[string]any, error)
	Auth(ctx context.Context) (AuthToken, error)
	ConfirmOrder(ctx context.Context, orderId string) (map[string]any, error)
	CreateOrder(ctx context.Context, productList []map[string]any) (OrderResponse, error)
	CustomerInfo(ctx context.Context) (CustomerInfoResponse, error)
	RetrieveOrder(ctx context.Context, orderId string) (map[string]any, error)
	ShopItem(ctx context.Context, productId string) (ProductResponse, error)
	ShopList(ctx context.Context, pageSize int, pageToken string) (map[string]any, error)
}
