package adaptor

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (g *GiftCard) RetrieveOrder(orderId string) (map[string]any, error) {
	url := g.BaseUrl + fmt.Sprintf("/order/get?%s", orderId)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
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
	fmt.Println(responseData)

	if res.StatusCode == http.StatusOK {
		return responseData, nil
	}
	return nil, &RetrieveOrderError{ErrorMsg: "failed to retrieve order", Response: responseData}
}

type RetrieveOrderError struct {
	ErrorMsg string
	Response map[string]interface{}
}

func (e *RetrieveOrderError) Error() string {
	return e.ErrorMsg
}
