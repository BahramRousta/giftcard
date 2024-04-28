package adaptor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (g *GiftCard) CreateOrder(productList []map[string]any) (map[string]any, error) {
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
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	token, err := g.Auth()
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	fmt.Println("status", res.StatusCode)

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var responseData map[string]any

	err = json.Unmarshal(bodyBytes, &responseData)
	if err != nil {
		return nil, err
	}
	fmt.Println(responseData)

	if res.StatusCode == http.StatusOK {
		return responseData, nil
	}
	return nil, &OrderError{ErrorMsg: "failed to create order", Response: responseData}
}

type OrderError struct {
	ErrorMsg string
	Response map[string]interface{}
}

func (e *OrderError) Error() string {
	return e.ErrorMsg
}
