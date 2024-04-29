package adaptor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ModifiedDate struct {
	Seconds     int64 `json:"_seconds"`
	Nanoseconds int   `json:"_nanoseconds"`
}

type exchangeRate struct {
	BaseCurrency   string       `json:"baseCurrency"`
	ModifiedDate   ModifiedDate `json:"modifiedDate"`
	Rate           float64      `json:"rate"`
	TargetCurrency string       `json:"targetCurrency"`
}

type Effect struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type MetaData struct {
	BaseCurrency   string       `json:"baseCurrency,omitempty"`
	ModifiedDate   ModifiedDate `json:"modifiedDate,omitempty"`
	Rate           float64      `json:"rate,omitempty"`
	TargetCurrency string       `json:"targetCurrency,omitempty"`
	Quantity       int          `json:"quantity,omitempty"`
	Quote          int          `json:"quote,omitempty"`
}

type Item struct {
	Description string   `json:"description"`
	Effect      Effect   `json:"effect"`
	MetaData    MetaData `json:"metaData"`
	Type        string   `json:"type"`
}

type Record struct {
	Items []Item `json:"items"`
	Sku   string `json:"sku"`
	Total struct {
		DKK float64 `json:"DKK"`
		EUR float64 `json:"EUR"`
	} `json:"total"`
}

type Invoice struct {
	PaymentMethod string   `json:"paymentMethod"`
	Records       []Record `json:"records"`
	Status        string   `json:"status"`
	Total         float64  `json:"total"`
	Wallet        string   `json:"wallet"`
}

type OrderResponse struct {
	Data struct {
		ExchangeRates []exchangeRate `json:"exchangeRates"`
		ExpiresAt     int64          `json:"expiresAt"`
		ID            string         `json:"id"`
		Invoice       Invoice        `json:"invoice"`
	} `json:"data"`
}

func (g *GiftCard) CreateOrder(productList []map[string]any) (OrderResponse, error) {
	url := g.BaseUrl + "/order/create"
	method := "POST"

	client := &http.Client{}

	payload := map[string]any{
		"productList": productList,
		"wallet":      "EUR",
		"reference":   "Test Reference",
		"webhookUrl":  "YOUR WEBHOOK URL",
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return OrderResponse{}, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return OrderResponse{}, err
	}

	token, err := g.Auth()
	if err != nil {
		return OrderResponse{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token)

	res, err := client.Do(req)
	if err != nil {
		return OrderResponse{}, err
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return OrderResponse{}, err
	}

	if res.StatusCode != http.StatusOK {
		return OrderResponse{}, &CreateOrderError{ErrorMsg: "failed to fetch customer info", Response: bodyBytes}
	}

	var responseData OrderResponse
	err = json.Unmarshal(bodyBytes, &responseData)
	if err != nil {
		return OrderResponse{}, err
	}
	fmt.Println("responseData", responseData)
	return responseData, nil
}

type CreateOrderError struct {
	ErrorMsg string
	Response []byte
}

func (e *CreateOrderError) Error() string {
	return e.ErrorMsg
}
