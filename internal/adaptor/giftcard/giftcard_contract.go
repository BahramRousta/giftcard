package giftcard

type IGiftCard interface {
	ProcessRequest(method string, url string, payload *[]byte) (map[string]any, error)
	Auth() (AuthToken, error)
	ConfirmOrder(orderId string) (map[string]any, error)
	CreateOrder(productList []map[string]any) (OrderResponse, error)
	CustomerInfo() (CustomerInfoResponse, error)
	RetrieveOrder(orderId string) (map[string]any, error)
	ShopItem(productId string) (ProductResponse, error)
	ShopList(pageSize int, pageToken string) (map[string]any, error)
}
