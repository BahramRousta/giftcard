package adaptor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (g *GiftCard) ConfirmOrder(orderId string) (map[string]any, error) {
	url := g.BaseUrl + "/order/confirm"
	method := "PUT"
	client := &http.Client{}
	fmt.Println(url)
	payload := map[string]any{
		"orderId": orderId,
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

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var responseData map[string]any

	err = json.Unmarshal(bodyBytes, &responseData)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusOK {
		return responseData, nil
	}
	return nil, &ConfirmOrderError{ErrorMsg: "failed to create order", Response: responseData}
}

type ConfirmOrderError struct {
	ErrorMsg string
	Response map[string]interface{}
}

func (e *ConfirmOrderError) Error() string {
	return e.ErrorMsg
}
