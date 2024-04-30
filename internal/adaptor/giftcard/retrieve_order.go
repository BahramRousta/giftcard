package adaptor

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// RetrieveOrder use to get latest status of order
func (g *GiftCard) RetrieveOrder(orderId string) (map[string]any, error) {
	url := g.BaseUrl + fmt.Sprintf("/order/get?orderId=%s", orderId)
	method := "GET"
	fmt.Println(url)
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
		return nil, errors.New("error while sending request")
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("error while process response")
	}

	if res.StatusCode == http.StatusForbidden {
		return nil, &ForbiddenErr{ErrMsg: "Forbidden to access end point."}
	}

	var responseData map[string]any

	err = json.Unmarshal(bodyBytes, &responseData)
	if err != nil {
		return nil, errors.New("error while process response 1")
	}
	fmt.Println("responseData", responseData)

	if res.StatusCode == http.StatusOK {
		return responseData, nil
	}
	return responseData, &RetrieveOrderError{ErrorMsg: "failed to retrieve order", Response: responseData}
}

type RetrieveOrderError struct {
	ErrorMsg string
	Response map[string]interface{}
}

func (e *RetrieveOrderError) Error() string {
	return e.ErrorMsg
}

type ForbiddenErr struct {
	ErrMsg string
}

func (e *ForbiddenErr) Error() string {
	return e.ErrMsg
}
