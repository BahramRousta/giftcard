package giftcard

import (
	"context"
	"encoding/json"
)

type ExchangeRate struct {
	BaseCurrency string `json:"baseCurrency"`
	ModifiedDate struct {
		Nanoseconds int64 `json:"_nanoseconds"`
		Seconds     int64 `json:"_seconds"`
	} `json:"modifiedDate"`
	Rate           float64 `json:"rate"`
	TargetCurrency string  `json:"targetCurrency"`
}

type Wallet struct {
	EUR struct {
		Balance       float64 `json:"balance"`
		CreditBalance int     `json:"creditBalance"`
		FrozenBalance int     `json:"frozenBalance"`
	} `json:"EUR"`
}

type CustomerInfoResponse struct {
	Data struct {
		ExchangeRates []ExchangeRate `json:"exchangeRates"`
		Name          string         `json:"name"`
		Wallet        Wallet         `json:"wallet"`
	} `json:"data"`
}

func (g *GiftCard) CustomerInfo(ctx context.Context) (CustomerInfoResponse, error) {
	url := g.BaseUrl + "/customer/info"
	method := "GET"

	data, err := g.ProcessRequest(ctx, method, url, nil)
	if err != nil {
		return CustomerInfoResponse{}, err
	}
	if err != nil {
		return CustomerInfoResponse{}, err
	}

	jsonData, err := json.Marshal(data)

	var responseData CustomerInfoResponse
	err = json.Unmarshal(jsonData, &responseData)
	if err != nil {
		return CustomerInfoResponse{}, err
	}
	return responseData, nil
}

type CustomerInfoError struct {
	ErrorMsg   string
	StatusCode int
	Response   []byte
}

func (e *CustomerInfoError) Error() string {
	return e.ErrorMsg
}
