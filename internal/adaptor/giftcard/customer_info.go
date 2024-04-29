package adaptor

import (
	"encoding/json"
	"io"
	"net/http"
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

func (g *GiftCard) CustomerInfo() (CustomerInfoResponse, error) {
	url := "https://sandbox-api.core.hub.gift/customer/info"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return CustomerInfoResponse{}, err
	}
	token, err := g.Auth()
	if err != nil {
		return CustomerInfoResponse{}, err
	}

	req.Header.Add("Authorization", token)
	res, err := client.Do(req)
	if err != nil {
		return CustomerInfoResponse{}, err
	}

	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return CustomerInfoResponse{}, err
	}
	if res.StatusCode != http.StatusOK {
		return CustomerInfoResponse{}, &CustomerInfoError{ErrorMsg: "failed to fetch customer info", StatusCode: res.StatusCode, Response: bodyBytes}
	}

	var responseData CustomerInfoResponse
	err = json.Unmarshal(bodyBytes, &responseData)
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
