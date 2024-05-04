package giftcard

import (
	"encoding/json"
)

func (g *GiftCard) ConfirmOrder(orderId string) (map[string]any, error) {
	url := g.BaseUrl + "/order/confirm"
	method := "PUT"

	payload := map[string]any{
		"orderId": orderId,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	data, err := g.ProcessRequest(method, url, &payloadBytes)

	if err != nil {
		return nil, err
	}
	return data, nil
}

type ConfirmOrderError struct {
	ErrorMsg string
	Response map[string]interface{}
}

func (e *ConfirmOrderError) Error() string {
	return e.ErrorMsg
}
