package adaptor

import (
	"fmt"
)

// RetrieveOrder use to get latest status of order
func (g *GiftCard) RetrieveOrder(orderId string) (map[string]any, error) {
	url := g.BaseUrl + fmt.Sprintf("/order/get?orderId=%s", orderId)
	method := "GET"

	data, err := g.ProcessRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type RetrieveOrderError struct {
	ErrorMsg string
	Response map[string]interface{}
}

func (e *RetrieveOrderError) Error() string {
	return e.ErrorMsg
}
